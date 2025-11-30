import { vi } from 'vitest'

export const useTranslation = () => ({
  t: (key) => key,
  i18n: {
    language: 'en',
    changeLanguage: vi.fn(),
  },
})

vi.mock('react-i18next', () => ({
  useTranslation,
  initReactI18next: {
    type: '3rdParty',
    init: vi.fn(),
  },
}))
