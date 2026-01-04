import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import HomePage from '@/views/HomePage.vue'

describe('HomePage', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('renders home page correctly', () => {
    const wrapper = mount(HomePage)
    expect(wrapper.text()).toContain('小窝同步观影')
  })

  it('shows welcome message', () => {
    const wrapper = mount(HomePage)
    expect(wrapper.find('[data-testid="welcome-message"]').exists()).toBe(true)
  })

  it('has create room button', () => {
    const wrapper = mount(HomePage)
    expect(wrapper.find('[data-testid="create-room-btn"]').exists()).toBe(true)
  })
})
