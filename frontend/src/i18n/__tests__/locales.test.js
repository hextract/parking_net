import { describe, it, expect } from 'vitest'
import en from '../locales/en.json'
import ru from '../locales/ru.json'

describe('Translation Files', () => {
  it('has English translations', () => {
    expect(en).toBeDefined()
    expect(en.app).toBeDefined()
    expect(en.auth).toBeDefined()
    expect(en.nav).toBeDefined()
  })

  it('has Russian translations', () => {
    expect(ru).toBeDefined()
    expect(ru.app).toBeDefined()
    expect(ru.auth).toBeDefined()
    expect(ru.nav).toBeDefined()
  })

  it('has matching keys structure between EN and RU', () => {
    const enKeys = Object.keys(en).sort()
    const ruKeys = Object.keys(ru).sort()
    expect(enKeys).toEqual(ruKeys)
  })

  it('has app translations in both languages', () => {
    expect(en.app.name).toBeDefined()
    expect(ru.app.name).toBeDefined()
  })

  it('has auth translations in both languages', () => {
    expect(en.auth.login).toBeDefined()
    expect(ru.auth.login).toBeDefined()
    expect(en.auth.register).toBeDefined()
    expect(ru.auth.register).toBeDefined()
  })

  it('has navigation translations in both languages', () => {
    expect(en.nav.dashboard).toBeDefined()
    expect(ru.nav.dashboard).toBeDefined()
  })

  it('has parking translations in both languages', () => {
    expect(en.parking.name).toBeDefined()
    expect(ru.parking.name).toBeDefined()
  })

  it('has booking translations in both languages', () => {
    expect(en.booking.bookNow).toBeDefined()
    expect(ru.booking.bookNow).toBeDefined()
  })

  it('has message translations in both languages', () => {
    expect(en.messages.loading).toBeDefined()
    expect(ru.messages.loading).toBeDefined()
  })
})
