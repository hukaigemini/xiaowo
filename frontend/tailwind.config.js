/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        // 优优的温馨色调系统 - 与静态模板保持一致
        bg: '#FFFBF0', // Warm Cream - 温馨奶油色背景
        primary: {
          DEFAULT: '#FF9F76', // Soft Coral - 珊瑚色
          hover: '#FF8C5A',
          50: '#FFF5F0',
          100: '#FFE6DC',
          200: '#FFD0C7',
          300: '#FFB3A0',
          400: '#FF9977',
          500: '#FF9F76',
          600: '#FF8C5A',
          700: '#E6784A',
          800: '#CC693D',
          900: '#B35A30',
        },
        secondary: {
          DEFAULT: '#A3C9A8', // Sage Green - 鼠尾草绿
          50: '#F0F5F0',
          100: '#DCE7DC',
          200: '#C5D5C5',
          300: '#A3C9A8',
          400: '#8BB58F',
          500: '#7BA478',
          600: '#6B9161',
          700: '#5B7C4A',
          800: '#4B6733',
          900: '#3B521C',
        },
        surface: '#F7EEDD', // Sand - 沙色表面
        text: {
          main: '#4A3B32', // Deep Coffee - 深咖啡色文字
          muted: '#8D7B68', // Mocha - 摩卡色文字
          light: '#FFFBF0', // Cream - 奶油色文字
        },
        dark: {
          bg: '#2C241F', // Dark Coffee - 影院模式深色
          surface: '#3E3228', // Darker wood tone
        },
        border: {
          DEFAULT: '#E6D5B8', // 边框色
        }
      },
      borderRadius: {
        'xl': '12px',
        '2xl': '16px',
        '3xl': '24px',
        '4xl': '32px'
      },
      boxShadow: {
        'warm': '0 10px 40px -10px rgba(255, 159, 118, 0.25)',
        'button-warm': '0 8px 20px -6px rgba(255, 159, 118, 0.4)',
        'button-hover': '0 10px 25px -8px rgba(255, 159, 118, 0.6)',
        'card': '0 20px 60px -15px rgba(74, 59, 50, 0.12)',
        'inner-warm': 'inset 0 2px 4px 0 rgba(74, 59, 50, 0.08)',
        'surface': '0 4px 12px -2px rgba(74, 59, 50, 0.08)',
      },
      fontFamily: {
        sans: ['Inter', 'system-ui', 'sans-serif'],
      },
      spacing: {
        '18': '4.5rem',
        '88': '22rem',
      },
      animation: {
        'fade-in': 'fadeIn 0.5s ease-in-out',
        'slide-up': 'slideUp 0.3s ease-out',
        'bounce-subtle': 'bounceSubtle 2s infinite',
        'pulse-warm': 'pulseWarm 2s infinite',
      },
      keyframes: {
        fadeIn: {
          '0%': { opacity: '0' },
          '100%': { opacity: '1' },
        },
        slideUp: {
          '0%': { transform: 'translateY(10px)', opacity: '0' },
          '100%': { transform: 'translateY(0)', opacity: '1' },
        },
        bounceSubtle: {
          '0%, 100%': { transform: 'translateY(0)' },
          '50%': { transform: 'translateY(-2px)' },
        },
        pulseWarm: {
          '0%, 100%': { opacity: '1' },
          '50%': { opacity: '0.5' },
        },
      },
    },
  },
  plugins: [],
}