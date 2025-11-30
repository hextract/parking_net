# Testing Guide - ParkingNet Frontend

Complete testing documentation for the ParkingNet frontend application.

## Table of Contents

1. [Testing Strategy](#testing-strategy)
2. [Unit Tests (Vitest)](#unit-tests-vitest)
3. [Integration Tests (Playwright)](#integration-tests-playwright)
4. [Running Tests](#running-tests)
5. [CI/CD Integration](#cicd-integration)
6. [Writing Tests](#writing-tests)

---

## Testing Strategy

Our testing approach follows the testing pyramid:

```
        /\
       /  \      E2E Tests (Playwright)
      /____\     - Critical user flows
     /      \    - Role-based access
    /  INTE \   - Complete journeys
   /__________\
  /            \ Unit Tests (Vitest)
 /   UNIT      \ - Components
/________________\ - Services
                  - Utils
```

### Test Types

| Type | Tool | Coverage | Purpose |
|------|------|----------|---------|
| **Unit** | Vitest | Components, Services, Utils | Fast feedback, isolated logic |
| **E2E** | Playwright | Complete user flows | Confidence in production |

---

## Unit Tests (Vitest)

### Location
```
src/
├── components/__tests__/
├── services/__tests__/
├── config/__tests__/
├── context/__tests__/
├── utils/__tests__/
└── __tests__/
```

### Running Unit Tests

```bash
# Run all unit tests
npm test

# Watch mode
npm run test:watch

# With UI
npm run test:ui

# With coverage
npm run test:coverage
```

### Example Unit Test

```javascript
import { render, screen } from '@testing-library/react'
import { describe, it, expect } from 'vitest'
import LoadingSpinner from '../LoadingSpinner'

describe('LoadingSpinner', () => {
  it('renders correctly', () => {
    render(<LoadingSpinner />)
    expect(screen.getByRole('status')).toBeInTheDocument()
  })
})
```

---

## Integration Tests (Playwright)

### Location
```
e2e/
├── auth.spec.js              # Authentication flows
├── driver-flow.spec.js       # Driver user journey
├── owner-flow.spec.js        # Owner user journey
├── navigation-i18n.spec.js   # Navigation & i18n
└── README.md                 # Detailed E2E docs
```

### Prerequisites

1. **Install Playwright** (first time):
   ```bash
   npm install
   npx playwright install chromium
   ```

2. **Start backend services**:
   ```bash
   cd /Users/kirillandreev/Desktop/Lessons/parking_net
   colima start
   colima ssh <<'EOF'
   cd /Users/kirillandreev/Desktop/Lessons/parking_net
   docker compose up -d
   EOF
   ```

3. **Verify services are running**:
   ```bash
   colima ssh -- docker ps
   ```

### Running E2E Tests

```bash
# Run all E2E tests
npm run test:e2e

# Interactive UI mode
npm run test:e2e:ui

# Headed mode (see browser)
npm run test:e2e:headed

# Specific test file
npx playwright test e2e/auth.spec.js

# Debug mode
npx playwright test --debug

# Specific test by name
npx playwright test -g "should login successfully"
```

### Test Coverage

#### Authentication (9 tests)
- ✅ Login page display
- ✅ Form validation
- ✅ Registration flow
- ✅ Password validation
- ✅ Auto-login after registration

#### Driver Flow (12 tests)
- ✅ Dashboard access
- ✅ Navigation links
- ✅ Search parking
- ✅ My bookings
- ✅ Role-based access control
- ✅ Logout/Login cycle

#### Owner Flow (11 tests)
- ✅ Dashboard access
- ✅ Navigation (correct tabs)
- ✅ Create parking
- ✅ Form validation
- ✅ Display parkings list
- ✅ Role-based access control

#### Navigation & i18n (15 tests)
- ✅ Language switching
- ✅ Language persistence
- ✅ Page navigation
- ✅ Admin panel
- ✅ Protected routes
- ✅ Responsive design

#### RBAC (2 tests)
- ✅ Driver cannot access owner routes
- ✅ Owner cannot access driver routes

**Total: 49 E2E tests**

---

## Running Tests

### Quick Commands

```bash
# All tests (unit + E2E)
npm test && npm run test:e2e

# Unit tests only
npm test

# E2E tests only
npm run test:e2e

# Watch unit tests
npm run test:watch

# Interactive E2E
npm run test:e2e:ui
```

### Test Reports

#### Unit Tests
```bash
npm run test:coverage
open coverage/index.html
```

#### E2E Tests
```bash
npx playwright show-report
```

---

## CI/CD Integration

### GitHub Actions

E2E tests run automatically on:
- Push to `main` or `develop`
- Pull requests to `main` or `develop`

Workflow file: `.github/workflows/e2e-tests.yml`

### Manual CI Run

```bash
CI=true npm run test:e2e
```

Features in CI:
- ✅ Automatic retries (2x)
- ✅ Artifact uploads (reports, screenshots)
- ✅ 30-day retention
- ✅ Fail-fast on `.only` tests

---

## Writing Tests

### Unit Test Example

```javascript
import { render, screen, fireEvent } from '@testing-library/react'
import { describe, it, expect, vi } from 'vitest'
import MyComponent from '../MyComponent'

describe('MyComponent', () => {
  it('should handle click', () => {
    const handleClick = vi.fn()
    render(<MyComponent onClick={handleClick} />)

    fireEvent.click(screen.getByRole('button'))

    expect(handleClick).toHaveBeenCalledOnce()
  })
})
```

### E2E Test Example

```javascript
import { test, expect } from '@playwright/test'

test.describe('Feature Name', () => {
  test('should do something', async ({ page }) => {
    await page.goto('/path')

    await page.getByLabel('Username').fill('testuser')
    await page.getByRole('button', { name: 'Submit' }).click()

    await expect(page).toHaveURL('/success')
  })
})
```

### Best Practices

#### Unit Tests
1. ✅ Test behavior, not implementation
2. ✅ Use meaningful test names
3. ✅ Mock external dependencies
4. ✅ Keep tests fast (<100ms)
5. ✅ One assertion per test (when possible)

#### E2E Tests
1. ✅ Use semantic selectors (role, label, text)
2. ✅ Create unique test data (timestamps)
3. ✅ Wait for network requests
4. ✅ Clean up state
5. ✅ Test critical paths only

---

## Debugging

### Unit Tests

```bash
# Debug in VS Code
# Set breakpoint, then F5

# Console log
console.log('Debug:', value)
```

### E2E Tests

```bash
# Playwright Inspector
npx playwright test --debug

# Headed mode
npm run test:e2e:headed

# Slow motion
npx playwright test --headed --slow-mo=1000

# Trace viewer
npx playwright show-trace trace.zip
```

### Common Issues

#### "Cannot find module"
```bash
npm install
```

#### "Port 5173 already in use"
```bash
# Kill existing process
lsof -ti:5173 | xargs kill -9
```

#### "Navigation timeout"
```bash
# Check backend is running
colima ssh -- docker ps

# Restart Keycloak
colima ssh -- docker restart keycloak
```

---

## Coverage Goals

| Type | Current | Goal |
|------|---------|------|
| Unit Tests | ~60% | 80% |
| E2E Tests | Critical paths | All user flows |

### Coverage Reports

Generate and view:
```bash
npm run test:coverage
open coverage/index.html
```

---

## Maintenance

### When to Update Tests

- ✅ New features added
- ✅ Bug fixes implemented
- ✅ UI/UX changes
- ✅ API changes
- ✅ Business logic changes

### Test Health Checks

Run weekly:
```bash
npm test
npm run test:e2e
```

Review:
- ❌ Flaky tests (fix or remove)
- ❌ Slow tests (optimize)
- ❌ Skipped tests (implement or remove)

---

## Resources

- [Vitest Documentation](https://vitest.dev/)
- [Playwright Documentation](https://playwright.dev/)
- [Testing Library](https://testing-library.com/)
- [E2E Best Practices](https://playwright.dev/docs/best-practices)

---

## Support

For testing questions:
1. Check this documentation
2. Review test examples in `/e2e/`
3. Check Playwright/Vitest docs
4. Contact the development team

---

**Last Updated**: November 2025
**Test Suite Version**: 1.0.0
**Total Tests**: 49 E2E + Unit tests
