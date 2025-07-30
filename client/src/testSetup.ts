 
import * as vitest from 'vitest';
import { cleanup } from '@testing-library/react'
import '@testing-library/jest-dom/vitest'
import { RuleTester } from '@typescript-eslint/rule-tester'
import { afterEach, vi } from 'vitest';


RuleTester.afterAll = vitest.afterAll

// 
afterEach(() => {
  cleanup()
  vi.clearAllMocks()
})