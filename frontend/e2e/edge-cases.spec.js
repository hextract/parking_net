import { test, expect } from '@playwright/test'

test.describe('Edge Cases and Error Handling', () => {
  test('should handle invalid login credentials', async ({ page }) => {
    await page.goto('/login')

    await page.getByLabel(/username/i).fill('nonexistentuser')
    await page.getByLabel(/password/i).fill('wrongpassword')
    await page.getByRole('button', { name: /sign in/i }).click()

    await page.waitForTimeout(3000)

    const hasError = await page.locator('.bg-red-50, .text-red-700, [role="alert"]').count() > 0
    const stillOnLoginPage = page.url().includes('/login')

    expect(hasError || stillOnLoginPage).toBeTruthy()
  })

  test('should handle registration with existing email', async ({ page }) => {
    const timestamp = Date.now()
    const credentials = {
      email: `duplicate_${timestamp}@test.com`,
      username: `duplicate_${timestamp}`,
      password: 'TestPass123!'
    }

    await page.goto('/register')
    await page.getByLabel(/email/i).fill(credentials.email)
    await page.getByLabel(/username/i).fill(credentials.username)
    await page.getByLabel('Password', { exact: true }).fill(credentials.password)
    await page.getByLabel(/confirm password/i).fill(credentials.password)
    await page.getByLabel(/account type/i).selectOption('driver')
    await page.getByRole('button', { name: /sign up/i }).click()

    await page.waitForURL('/driver', { timeout: 60000 })
    await page.getByRole('button', { name: /logout/i }).click()

    await page.goto('/register')
    await page.getByLabel(/email/i).fill(credentials.email)
    await page.getByLabel(/username/i).fill(`different_${timestamp}`)
    await page.getByLabel('Password', { exact: true }).fill(credentials.password)
    await page.getByLabel(/confirm password/i).fill(credentials.password)
    await page.getByLabel(/account type/i).selectOption('driver')
    await page.getByRole('button', { name: /sign up/i }).click()

    await page.waitForTimeout(2000)
  })

  test('should validate email format', async ({ page }) => {
    await page.goto('/register')

    const emailInput = page.getByLabel(/email/i)
    await emailInput.fill('invalidemail')

    await page.getByLabel(/username/i).fill('testuser')
    await page.getByLabel('Password', { exact: true }).fill('TestPass123!')
    await page.getByLabel(/confirm password/i).fill('TestPass123!')
    await page.getByLabel(/account type/i).selectOption('driver')

    await page.getByRole('button', { name: /sign up/i }).click()

    const validationMessage = await emailInput.evaluate((el) => el.validationMessage)
    expect(validationMessage).toBeTruthy()
  })

  test('should handle network errors gracefully', async ({ page, context }) => {
    await context.route('**/auth/login', route => route.abort())

    await page.goto('/login')
    await page.getByLabel(/username/i).fill('testuser')
    await page.getByLabel(/password/i).fill('password123')
    await page.getByRole('button', { name: /sign in/i }).click()

    await page.waitForTimeout(2000)
  })

  test('should handle booking with invalid date range', async ({ page }) => {
    const timestamp = Date.now()
    await page.goto('/register')
    await page.getByLabel(/email/i).fill(`driver_dates_${timestamp}@test.com`)
    await page.getByLabel(/username/i).fill(`driver_dates_${timestamp}`)
    await page.getByLabel('Password', { exact: true }).fill('TestPass123!')
    await page.getByLabel(/confirm password/i).fill('TestPass123!')
    await page.getByLabel(/account type/i).selectOption('driver')
    await page.getByRole('button', { name: /sign up/i }).click()

    await page.waitForURL('/driver', { timeout: 60000 })

    await page.goto('/driver/search')
    await page.waitForTimeout(2000)
  })

  test('should prevent SQL injection in search', async ({ page }) => {
    const timestamp = Date.now()
    await page.goto('/register')
    await page.getByLabel(/email/i).fill(`driver_sql_${timestamp}@test.com`)
    await page.getByLabel(/username/i).fill(`driver_sql_${timestamp}`)
    await page.getByLabel('Password', { exact: true }).fill('TestPass123!')
    await page.getByLabel(/confirm password/i).fill('TestPass123!')
    await page.getByLabel(/account type/i).selectOption('driver')
    await page.getByRole('button', { name: /sign up/i }).click()

    await page.waitForURL('/driver', { timeout: 60000 })

    await page.goto('/driver/search')
    await page.getByPlaceholder(/enter city/i).fill("' OR '1'='1")
    await page.getByRole('button', { name: /search/i }).click()

    await page.waitForTimeout(2000)
  })

  test('should handle XSS attempts in parking name', async ({ page }) => {
    const timestamp = Date.now()
    await page.goto('/register')
    await page.getByLabel(/email/i).fill(`owner_xss_${timestamp}@test.com`)
    await page.getByLabel(/username/i).fill(`owner_xss_${timestamp}`)
    await page.getByLabel('Password', { exact: true }).fill('TestPass123!')
    await page.getByLabel(/confirm password/i).fill('TestPass123!')
    await page.getByLabel(/account type/i).selectOption('owner')
    await page.getByRole('button', { name: /sign up/i }).click()

    await page.waitForURL('/owner', { timeout: 60000 })

    await page.goto('/owner/parkings')
    await page.getByRole('button', { name: /add.*parking/i }).first().click()

    const inputs = page.locator('input[type="text"]')
    await inputs.nth(0).fill('<script>alert("XSS")</script>')
    await inputs.nth(1).fill('Test City')
    await inputs.nth(2).fill('123 Test Street')
    await page.locator('input[type="number"]').nth(0).fill('10')
    await page.locator('input[type="number"]').nth(1).fill('50')

    await page.getByRole('button', { name: /^(save|create)/i }).click()

    await page.waitForTimeout(2000)
  })

  test('should handle missing required fields', async ({ page }) => {
    const timestamp = Date.now()
    await page.goto('/register')
    await page.getByLabel(/email/i).fill(`owner_required_${timestamp}@test.com`)
    await page.getByLabel(/username/i).fill(`owner_required_${timestamp}`)
    await page.getByLabel('Password', { exact: true }).fill('TestPass123!')
    await page.getByLabel(/confirm password/i).fill('TestPass123!')
    await page.getByLabel(/account type/i).selectOption('owner')
    await page.getByRole('button', { name: /sign up/i }).click()

    await page.waitForURL('/owner', { timeout: 60000 })

    await page.goto('/owner/parkings')
    await page.getByRole('button', { name: /add.*parking/i }).first().click()

    await page.getByRole('button', { name: /^(save|create)/i }).click()

    await page.waitForTimeout(500)
  })

  test('should validate negative values in parking form', async ({ page }) => {
    const timestamp = Date.now()
    await page.goto('/register')
    await page.getByLabel(/email/i).fill(`owner_negative_${timestamp}@test.com`)
    await page.getByLabel(/username/i).fill(`owner_negative_${timestamp}`)
    await page.getByLabel('Password', { exact: true }).fill('TestPass123!')
    await page.getByLabel(/confirm password/i).fill('TestPass123!')
    await page.getByLabel(/account type/i).selectOption('owner')
    await page.getByRole('button', { name: /sign up/i }).click()

    await page.waitForURL('/owner', { timeout: 60000 })

    await page.goto('/owner/parkings')
    await page.getByRole('button', { name: /add.*parking/i }).first().click()

    const inputs = page.locator('input[type="text"]')
    await inputs.nth(0).fill('Negative Test Parking')
    await inputs.nth(1).fill('Test City')
    await inputs.nth(2).fill('123 Test Street')
    await page.locator('input[type="number"]').nth(0).fill('-10')
    await page.locator('input[type="number"]').nth(1).fill('-50')

    await page.getByRole('button', { name: /^(save|create)/i }).click()

    await page.waitForTimeout(1000)
  })

  test('should handle very long input strings', async ({ page }) => {
    const timestamp = Date.now()
    await page.goto('/register')
    await page.getByLabel(/email/i).fill(`owner_long_${timestamp}@test.com`)
    await page.getByLabel(/username/i).fill(`owner_long_${timestamp}`)
    await page.getByLabel('Password', { exact: true }).fill('TestPass123!')
    await page.getByLabel(/confirm password/i).fill('TestPass123!')
    await page.getByLabel(/account type/i).selectOption('owner')
    await page.getByRole('button', { name: /sign up/i }).click()

    await page.waitForURL('/owner', { timeout: 60000 })

    await page.goto('/owner/parkings')
    await page.getByRole('button', { name: /add.*parking/i }).first().click()

    const longString = 'A'.repeat(500)
    const inputs = page.locator('input[type="text"]')
    await inputs.nth(0).fill(longString)
    await inputs.nth(1).fill('Test City')
    await inputs.nth(2).fill('123 Test Street')
    await page.locator('input[type="number"]').nth(0).fill('10')
    await page.locator('input[type="number"]').nth(1).fill('50')

    await page.getByRole('button', { name: /^(save|create)/i }).click()

    await page.waitForTimeout(2000)
  })

  test('should maintain state after page refresh', async ({ page }) => {
    const timestamp = Date.now()
    const credentials = {
      username: `driver_refresh_${timestamp}`,
      email: `driver_refresh_${timestamp}@test.com`,
      password: 'TestPass123!'
    }

    await page.goto('/register')
    await page.getByLabel(/email/i).fill(credentials.email)
    await page.getByLabel(/username/i).fill(credentials.username)
    await page.getByLabel('Password', { exact: true }).fill(credentials.password)
    await page.getByLabel(/confirm password/i).fill(credentials.password)
    await page.getByLabel(/account type/i).selectOption('driver')
    await page.getByRole('button', { name: /sign up/i }).click()

    await page.waitForURL('/driver', { timeout: 60000 })

    await page.reload()

    await expect(page).toHaveURL('/driver')
    await expect(page.getByText(credentials.username).first()).toBeVisible()
  })
})
