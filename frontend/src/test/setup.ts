import { expect, afterEach } from 'vitest'
import { cleanup } from '@testing-library/vue'
import * as matchers from '@testing-library/jest-dom/matchers'

// extends Vitest's expect with jest-dom matchers
expect.extend(matchers)

// runs a cleanup after each test case
afterEach(() => {
  cleanup()
})
