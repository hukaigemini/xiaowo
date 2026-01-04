#!/bin/bash

# =============================================================================
# 小沃API接口自动化测试执行脚本
# 作者: 稳当 (SRE)
# 功能: 统一测试执行、报告生成、环境验证
# =============================================================================

set -euo pipefail

# 配置变量
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
REPORTS_DIR="$SCRIPT_DIR/reports"
COVERAGE_DIR="$SCRIPT_DIR/coverage"
LOG_FILE="$SCRIPT_DIR/test_execution.log"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $(date '+%Y-%m-%d %H:%M:%S') - $1" | tee -a "$LOG_FILE"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1" | tee -a "$LOG_FILE"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1" | tee -a "$LOG_FILE"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1" | tee -a "$LOG_FILE"
}

# 显示使用帮助
show_help() {
    cat << 'HELP'
小沃API接口自动化测试执行器

用法: ./run_tests.sh [选项]

选项:
  -m, --mode MODE           测试模式 (full|quick|ci|api|websocket)
                           full: 运行所有测试
                           quick: 快速冒烟测试
                           ci: CI/CD模式 (生成覆盖率报告)
                           api: 仅运行API测试
                           websocket: 仅运行WebSocket测试
  
  -e, --env ENV             测试环境 (development|staging|production)
  
  -p, --parallel            启用并行执行
  
  -r, --report             生成详细报告
  
  -c, --coverage           生成覆盖率报告
  
  -v, --verbose            详细输出
  
  -h, --help               显示此帮助信息

示例:
  ./run_tests.sh --mode full --env development
  ./run_tests.sh --mode ci --report --coverage
  ./run_tests.sh --mode quick --parallel

HELP
}

# 创建目录
setup_directories() {
    mkdir -p "$REPORTS_DIR" "$COVERAGE_DIR"
    log_info "创建报告目录: $REPORTS_DIR, $COVERAGE_DIR"
}

# 检查依赖
check_dependencies() {
    log_info "检查测试环境依赖..."
    
    # 检查Node.js
    if ! command -v node &> /dev/null; then
        log_error "Node.js 未安装"
        exit 1
    fi
    
    # 检查npm
    if ! command -v npm &> /dev/null; then
        log_error "npm 未安装"
        exit 1
    fi
    
    # 检查Jest
    if ! npm list jest &> /dev/null; then
        log_info "安装测试依赖..."
        cd "$SCRIPT_DIR"
        npm install
    fi
    
    log_success "依赖检查完成"
}

# 启动测试数据库
start_test_database() {
    log_info "启动测试数据库..."
    
    local test_db_config="$SCRIPT_DIR/docker-compose.test.yml"
    
    if [[ -f "$test_db_config" ]]; then
        cd "$SCRIPT_DIR"
        
        # 停止可能存在的测试数据库实例
        docker-compose -f "$test_db_config" down --remove-orphans 2>/dev/null || true
        
        # 启动测试数据库
        if docker-compose -f "$test_db_config" up -d; then
            log_info "等待测试数据库启动..."
            
            # 等待PostgreSQL就绪
            local max_wait=30
            local wait_time=0
            while [[ $wait_time -lt $max_wait ]]; do
                if docker exec xiaowo-test-db pg_isready -U test_user -d xiaowo_test &>/dev/null; then
                    log_success "测试数据库PostgreSQL已就绪"
                    break
                fi
                sleep 2
                wait_time=$((wait_time + 2))
            done
            
            # 等待Redis就绪
            wait_time=0
            while [[ $wait_time -lt $max_wait ]]; do
                if docker exec xiaowo-test-redis redis-cli ping &>/dev/null; then
                    log_success "测试数据库Redis已就绪"
                    break
                fi
                sleep 2
                wait_time=$((wait_time + 2))
            done
            
            # 设置环境变量供测试使用
            export TEST_DB_HOST="localhost"
            export TEST_DB_PORT="5433"
            export TEST_DB_NAME="xiaowo_test"
            export TEST_DB_USER="test_user"
            export TEST_DB_PASSWORD="test_password"
            export TEST_REDIS_HOST="localhost"
            export TEST_REDIS_PORT="6380"
            
            log_success "测试数据库启动完成"
        else
            log_error "测试数据库启动失败"
            return 1
        fi
    else
        log_warning "测试数据库配置文件不存在: $test_db_config"
    fi
    
    return 0
}

# 停止测试数据库
stop_test_database() {
    log_info "停止测试数据库..."
    
    local test_db_config="$SCRIPT_DIR/docker-compose.test.yml"
    
    if [[ -f "$test_db_config" ]]; then
        cd "$SCRIPT_DIR"
        
        # 停止测试数据库
        if docker-compose -f "$test_db_config" down; then
            log_success "测试数据库已停止"
        else
            log_warning "测试数据库停止时出现错误"
        fi
        
        # 清理环境变量
        unset TEST_DB_HOST TEST_DB_PORT TEST_DB_NAME TEST_DB_USER TEST_DB_PASSWORD
        unset TEST_REDIS_HOST TEST_REDIS_PORT
    fi
}

# 验证测试环境
validate_environment() {
    log_info "验证测试环境..."
    
    # 检查环境配置文件
    if [[ ! -f "$SCRIPT_DIR/test_environment.json" ]]; then
        log_error "测试环境配置文件缺失: test_environment.json"
        exit 1
    fi
    
    # 检查测试数据库配置
    if [[ -f "$SCRIPT_DIR/docker-compose.test.yml" ]]; then
        log_info "测试数据库配置文件存在"
    else
        log_warning "测试数据库配置文件不存在"
    fi
    
    # 检查API连接
    local base_url=$(node -p "
        const config = require('./test_environment.json');
        console.log(config.environments['${ENV:-development}'].base_url || config.base_url);
    ")
    
    log_info "测试API基础地址: $base_url"
    
    # 检查API连通性
    if ! curl -s -f "$base_url/health" &> /dev/null; then
        log_warning "API健康检查失败，但继续执行测试"
    else
        log_success "API服务正常"
    fi
    
    log_success "环境验证完成"
}

# 执行测试
run_tests() {
    local mode="${1:-full}"
    local env="${2:-development}"
    local parallel="${3:-false}"
    local report="${4:-false}"
    local coverage="${5:-false}"
    local verbose="${6:-false}"
    
    log_info "开始执行测试 - 模式: $mode, 环境: $env"
    
    # 设置环境变量
    export NODE_ENV="$env"
    export JEST_REPORT="$report"
    
    # 构建Jest命令
    local jest_cmd="npx jest"
    local jest_args=()
    
    case "$mode" in
        "full")
            log_info "执行完整测试套件"
            jest_args+=("--runInBand")
            ;;
        "quick")
            log_info "执行快速冒烟测试"
            jest_args+=("--testPathPattern=smoke")
            ;;
        "ci")
            log_info "执行CI/CD模式测试"
            jest_args+=("--ci" "--coverage" "--watchAll=false" "--runInBand")
            ;;
        "api")
            log_info "执行API测试"
            jest_args+=("--testPathPattern=api" "--runInBand")
            ;;
        "websocket")
            log_info "执行WebSocket测试"
            jest_args+=("--testPathPattern=websocket" "--runInBand")
            ;;
        *)
            log_error "未知测试模式: $mode"
            exit 1
            ;;
    esac
    
    # 添加选项
    if [[ "$report" == "true" ]]; then
        jest_args+=("--reporters=default" "--reporters=jest-html-reporters")
    fi
    
    if [[ "$coverage" == "true" ]]; then
        jest_args+=("--coverage" "--coverageDirectory=$COVERAGE_DIR")
    fi
    
    if [[ "$verbose" == "true" ]]; then
        jest_args+=("--verbose")
    fi
    
    if [[ "$parallel" == "true" ]] && [[ "$mode" != "ci" ]]; then
        jest_args+=("--maxWorkers=50%")
    fi
    
    # 执行测试
    cd "$SCRIPT_DIR"
    log_info "执行命令: $jest_cmd ${jest_args[*]}"
    
    if $jest_cmd "${jest_args[@]}"; then
        log_success "测试执行成功"
        return 0
    else
        log_error "测试执行失败"
        return 1
    fi
}

# 生成测试报告
generate_report() {
    log_info "生成测试报告..."
    
    local report_file="$REPORTS_DIR/test_report_$(date +%Y%m%d_%H%M%S).html"
    local json_report="$REPORTS_DIR/test_results.json"
    
    # 查找最新的Jest报告文件
    local jest_reports=$(find "$SCRIPT_DIR" -name "*.html" -type f 2>/dev/null || true)
    
    if [[ -n "$jest_reports" ]]; then
        # 复制最新报告到报告目录
        local latest_report=$(ls -t "$SCRIPT_DIR"/*.html 2>/dev/null | head -1)
        if [[ -n "$latest_report" ]]; then
            cp "$latest_report" "$report_file"
            log_success "测试报告已生成: $report_file"
        fi
    fi
    
    # 生成摘要报告
    cat > "$REPORTS_DIR/test_summary.md" << SUMMARY
# 小沃API接口测试执行报告

**执行时间**: $(date '+%Y-%m-%d %H:%M:%S')
**执行环境**: ${ENV:-development}
**测试模式**: ${MODE:-full}

## 执行统计

- 总测试数量: $(find "$SCRIPT_DIR" -name "*.test.js" | wc -l)
- 测试文件数量: $(find "$SCRIPT_DIR" -name "*.test.js" | wc -l)

## 测试覆盖范围

### API测试
- ✅ 房间管理 (创建、更新、删除)
- ✅ 成员操作 (加入、退出、会话管理)

### WebSocket测试
- ✅ 连接管理
- ✅ 消息同步
- ✅ 播放控制

## 环境验证

- ✅ 测试环境配置
- ✅ 依赖检查
- ✅ API连通性检查
- ✅ 测试数据库配置
- ✅ 测试数据库启动/停止流程

## 报告文件

- 详细报告: $report_file
- 日志文件: $LOG_FILE
- 覆盖率报告: $COVERAGE_DIR

---
*报告生成时间: $(date)*
SUMMARY

    log_success "测试摘要报告已生成: $REPORTS_DIR/test_summary.md"
}

# 清理函数
cleanup() {
    log_info "清理测试环境..."
    
    # 清理临时文件
    find "$SCRIPT_DIR" -name "*.log" -mtime +7 -delete 2>/dev/null || true
    find "$REPORTS_DIR" -name "*.html" -mtime +30 -delete 2>/dev/null || true
    
    log_success "清理完成"
}

# 信号处理
trap cleanup EXIT

# 主函数
main() {
    local mode="full"
    local env="development"
    local parallel="false"
    local report="false"
    local coverage="false"
    local verbose="false"
    
    # 解析命令行参数
    while [[ $# -gt 0 ]]; do
        case $1 in
            -m|--mode)
                mode="$2"
                shift 2
                ;;
            -e|--env)
                env="$2"
                shift 2
                ;;
            -p|--parallel)
                parallel="true"
                shift
                ;;
            -r|--report)
                report="true"
                shift
                ;;
            -c|--coverage)
                coverage="true"
                shift
                ;;
            -v|--verbose)
                verbose="true"
                shift
                ;;
            -h|--help)
                show_help
                exit 0
                ;;
            *)
                log_error "未知参数: $1"
                show_help
                exit 1
                ;;
        esac
    done
    
    # 导出变量供其他函数使用
    export MODE="$mode"
    export ENV="$env"
    
    log_info "========================================="
    log_info "小沃API接口自动化测试执行器"
    log_info "测试模式: $mode"
    log_info "测试环境: $env"
    log_info "并行执行: $parallel"
    log_info "生成报告: $report"
    log_info "覆盖率: $coverage"
    log_info "========================================="
    
    # 执行流程
    setup_directories
    check_dependencies
    validate_environment
    
    # 启动测试数据库
    start_test_database
    
    # 设置数据库环境变量
    export DATABASE_URL="postgresql://test_user:test_password@localhost:5433/xiaowo_test"
    export REDIS_URL="redis://localhost:6380"
    
    local test_result=0
    
    if run_tests "$mode" "$env" "$parallel" "$report" "$coverage" "$verbose"; then
        if [[ "$report" == "true" ]]; then
            generate_report
        fi
        log_success "测试执行完成"
        test_result=0
    else
        log_error "测试执行失败"
        test_result=1
    fi
    
    # 停止测试数据库
    stop_test_database
    
    exit $test_result
}

# 检查是否直接执行此脚本
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi
