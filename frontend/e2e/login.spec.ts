import { test, expect } from '@playwright/test'

test.describe('Login Flow', () => {
  test('should display login page', async ({ page }) => {
    await page.goto('/')
    
    // Should redirect to login
    await expect(page).toHaveURL('/login')
    
    // Check for login form elements
    await expect(page.getByRole('heading', { name: /3X-UI Modern/i })).toBeVisible()
    await expect(page.getByLabel(/username/i)).toBeVisible()
    await expect(page.getByLabel(/password/i)).toBeVisible()
    await expect(page.getByRole('button', { name: /login/i })).toBeVisible()
  })

  test('should show validation errors for empty fields', async ({ page }) => {
    await page.goto('/login')
    
    // Click login without filling fields
    await page.getByRole('button', { name: /login/i }).click()
    
    // Browser native validation should prevent submission
    const username = page.getByLabel(/username/i)
    await expect(username).toHaveAttribute('required')
  })

  test('should handle login with invalid credentials', async ({ page }) => {
    await page.goto('/login')
    
    // Fill in invalid credentials
    await page.getByLabel(/username/i).fill('wronguser')
    await page.getByLabel(/password/i).fill('wrongpass')
    await page.getByRole('button', { name: /login/i }).click()
    
    // Should show error message (adjust selector based on your implementation)
    // await expect(page.getByText(/invalid credentials/i)).toBeVisible()
  })
})
