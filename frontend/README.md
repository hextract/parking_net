# ParkingNet Frontend

Modern React.js frontend for the ParkingNet parking booking system with multi-language support (English/Russian).

## Features

### For Drivers
- 🔍 **Search Parking** - Find available parking spaces by city, name, or type
- 📅 **Create Bookings** - Book parking spaces with date range selection
- 📋 **Manage Bookings** - View and cancel your bookings
- 💰 **Cost Calculation** - Automatic cost calculation based on duration

### For Parking Owners
- 🏢 **Manage Parkings** - Create, update, and delete your parking places
- 📊 **View Bookings** - See all bookings for your parking places
- ✅ **Booking Management** - Confirm or cancel booking requests
- 📈 **Dashboard** - Overview of your parking business

### Common Features
- 🔐 **Authentication** - Secure login and registration
- 🎨 **Modern UI** - Beautiful, responsive design with Tailwind CSS
- 🚀 **Fast Performance** - Built with Vite for lightning-fast development
- 📱 **Mobile Responsive** - Works perfectly on all devices
- 🔄 **Real-time Updates** - Dynamic data loading and state management

## Technology Stack

- **React 18.3** - Modern React with hooks
- **React Router 6** - Client-side routing
- **Vite 5** - Next-generation frontend tooling
- **Tailwind CSS 3** - Utility-first CSS framework
- **Axios** - HTTP client for API calls
- **Lucide React** - Beautiful icon library
- **date-fns** - Modern date utility library

## Project Structure

```
frontend/
├── public/              # Static assets
├── src/
│   ├── components/      # Reusable components
│   │   ├── Layout.jsx
│   │   ├── Navbar.jsx
│   │   ├── ProtectedRoute.jsx
│   │   └── LoadingSpinner.jsx
│   ├── context/         # React context providers
│   │   └── AuthContext.jsx
│   ├── pages/           # Page components
│   │   ├── auth/        # Authentication pages
│   │   │   ├── LoginPage.jsx
│   │   │   └── RegisterPage.jsx
│   │   ├── driver/      # Driver role pages
│   │   │   ├── DriverDashboard.jsx
│   │   │   ├── SearchParking.jsx
│   │   │   └── MyBookings.jsx
│   │   └── owner/       # Owner role pages
│   │       ├── OwnerDashboard.jsx
│   │       ├── MyParkings.jsx
│   │       └── ParkingBookings.jsx
│   ├── services/        # API service layer
│   │   ├── api.js
│   │   ├── authService.js
│   │   ├── parkingService.js
│   │   └── bookingService.js
│   ├── config/          # Configuration files
│   │   ├── api.js       # API endpoints configuration
│   │   └── env.js       # Environment variables
│   ├── App.jsx          # Root component
│   ├── main.jsx         # Application entry point
│   └── index.css        # Global styles
├── Dockerfile           # Docker configuration
├── nginx.conf           # Nginx configuration for production
├── package.json         # Dependencies and scripts
├── vite.config.js       # Vite configuration
└── tailwind.config.js   # Tailwind CSS configuration
```

## Prerequisites

- Node.js 18+ and npm
- Backend services running (Auth, Parking, Booking)

## Architecture

The frontend communicates with backend services through the **nginx gateway**:

```
Frontend (Port 3000) → Nginx Gateway (Port 80) → Backend Services
                                                  ├─ Auth (8800)
                                                  ├─ Parking (8888)
                                                  └─ Booking (8880)
```

All API calls go through nginx at `/auth`, `/parking`, and `/booking` endpoints.

## Quick Deployment Examples

### Same Server (Development)
```bash
Backend: http://localhost (nginx on port 80)
Frontend: http://localhost:3000

cd frontend
echo "VITE_API_BASE_URL=http://localhost" > .env
docker-compose up -d
```

### Separate Servers (Production)
```bash
Backend Server: 192.168.1.100
Frontend Server: 192.168.1.101

# On frontend server
docker build --build-arg VITE_API_BASE_URL=http://192.168.1.100 -t parking-frontend .
docker run -d -p 3000:80 parking-frontend
```

### With Domain Names
```bash
Backend: https://api.yourcompany.com
Frontend: https://app.yourcompany.com

docker build --build-arg VITE_API_BASE_URL=https://api.yourcompany.com -t parking-frontend .
docker run -d -p 3000:80 parking-frontend
```

## Getting Started

### 1. Install Dependencies

```bash
cd frontend
npm install
```

### 2. Environment Configuration

The frontend reads configuration from the root `.env` file (using ports defined there).

**Default Configuration (Local Development):**
Works out-of-the-box with ports from `../.env-example`:
- Auth Service: `http://localhost:8800`
- Parking Service: `http://localhost:8888`
- Booking Service: `http://localhost:8880`
- Keycloak: `http://localhost:8080`
- Jaeger: `http://localhost:16686`
- Prometheus: `http://localhost:9090`
- Grafana: `http://localhost:3000`

**Custom Configuration:**
To override defaults, set environment variables at build time:

```bash
# Change base host for all services
export VITE_BASE_HOST=192.168.1.100

# Or override specific ports
export VITE_AUTH_REST_PORT=8801
export VITE_PARKING_REST_PORT=8889

# Or provide full URLs
export VITE_API_BASE_URL=https://api.production.com
export VITE_AUTH_SERVICE_URL=https://auth.production.com

npm run build
```

The application automatically constructs URLs from the root `.env` configuration.

### 3. Start Development Server

```bash
npm run dev
```

The application will be available at `http://localhost:3000`

### 4. Build for Production

```bash
npm run build
```

Built files will be in the `dist/` directory.

### 5. Run Tests

```bash
npm test
```

Watch mode:
```bash
npm run test:watch
```

With UI:
```bash
npm run test:ui
```

Coverage report:
```bash
npm run test:coverage
```

## Docker Deployment

### Standalone Deployment (Separate Servers)

**Scenario: Backend at 192.168.1.100, Frontend at 192.168.1.101**

On frontend server:
```bash
docker build -t parking-frontend \
  --build-arg VITE_API_BASE_URL=http://192.168.1.100 \
  .

docker run -d -p 3000:80 --name parking-frontend parking-frontend
```

Access frontend at: http://192.168.1.101:3000

### Using Docker Compose (Same Network)

If backend and frontend are on the same Docker network:

```bash
cd frontend
docker-compose up -d
```

The frontend will connect to `parking-network` and use http://nginx-gateway

### Local Development

```bash
echo "VITE_API_BASE_URL=http://localhost" > .env
npm install
npm run dev
```

## API Integration

The frontend communicates with backend through nginx gateway:

### API Endpoints

All requests go through `VITE_API_BASE_URL`:

**Auth Endpoints:**
- `POST {BASE_URL}/auth/login`
- `POST {BASE_URL}/auth/register`
- `POST {BASE_URL}/auth/change-password`

**Parking Endpoints:**
- `GET {BASE_URL}/parking`
- `POST {BASE_URL}/parking`
- `GET {BASE_URL}/parking/{id}`
- `PUT {BASE_URL}/parking/{id}`
- `DELETE {BASE_URL}/parking/{id}`

**Booking Endpoints:**
- `GET {BASE_URL}/booking`
- `POST {BASE_URL}/booking`
- `GET {BASE_URL}/booking/{id}`
- `PUT {BASE_URL}/booking/{id}`
- `DELETE {BASE_URL}/booking/{id}`

### CORS Configuration

Backend nginx must allow requests from frontend origin:
```nginx
add_header Access-Control-Allow-Origin "http://your-frontend-server:3000";
```

## Authentication

The app uses JWT token-based authentication:

1. User logs in or registers
2. Backend returns a JWT token
3. Token is stored in localStorage
4. Token is sent with each API request in the `api_key` header
5. If token is invalid (401), user is redirected to login

## Role-Based Access Control

The application has two user roles:

### Driver Role
- Access to: Dashboard, Search Parking, My Bookings
- Can: Search parkings, create bookings, cancel own bookings
- Cannot: Create or manage parking places

### Owner Role
- Access to: Dashboard, My Parkings, Parking Bookings, Admin Panel
- Can: Create parkings, manage parkings, view all bookings, confirm/cancel bookings
- Cannot: Create bookings

### Admin Panel
- Access to: Both drivers and owners (authenticated users)
- Features:
  - Quick links to monitoring tools (Jaeger, Prometheus, Grafana)
  - Backend service status and endpoints
  - System health overview

Routes are protected using the `ProtectedRoute` component.

## Development

### Available Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run preview` - Preview production build
- `npm run lint` - Run ESLint

### Code Style

- Use functional components with hooks
- Follow React best practices
- Use Tailwind CSS utility classes
- Keep components small and focused
- Extract reusable logic into custom hooks
- Use async/await for API calls

### Component Patterns

**Loading States:**
```jsx
{loading ? <LoadingSpinner /> : <Content />}
```

**Error Handling:**
```jsx
{error && (
  <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
    {error}
  </div>
)}
```

**Protected Routes:**
```jsx
<Route
  path="/driver"
  element={
    <ProtectedRoute requiredRole="driver">
      <DriverDashboard />
    </ProtectedRoute>
  }
/>
```

## Styling

The project uses Tailwind CSS with custom configurations:

### Custom Colors
```js
primary: {
  500: '#3b82f6',
  600: '#2563eb',
  700: '#1d4ed8',
}
```

### Custom Components (in index.css)
- `.btn-primary` - Primary action buttons
- `.btn-secondary` - Secondary action buttons
- `.btn-danger` - Destructive action buttons
- `.input-field` - Form input fields
- `.card` - Card containers
- `.badge-*` - Status badges

## Language Support

The application supports **English** and **Russian**:

- Toggle language using the button in navigation (EN/RU)
- Language preference is saved to localStorage
- Auto-detects browser language on first visit

To add more languages:
1. Add translation file in `src/i18n/locales/{lang}.json`
2. Update `src/i18n/config.js`

## Troubleshooting

### API Connection Issues

If you can't connect to the backend:

1. Check nginx gateway is running:
   ```bash
   curl http://localhost/auth/metrics
   ```

2. Verify `VITE_API_BASE_URL` points to nginx gateway

3. Check browser console for CORS errors

4. Ensure nginx allows requests from frontend origin

### Build Issues

If build fails:

1. Clear node_modules and reinstall:
   ```bash
   rm -rf node_modules package-lock.json
   npm install
   ```

2. Clear Vite cache:
   ```bash
   rm -rf .vite
   ```

3. Check Node.js version (requires 18+):
   ```bash
   node --version
   ```

### Login Issues

If login doesn't work:

1. Check Auth service is running
2. Verify credentials in Keycloak
3. Check browser localStorage for token
4. Clear browser cache and localStorage
5. Check backend logs for errors

## Performance Optimization

The application includes several optimizations:

- **Code Splitting** - Automatic route-based code splitting with React Router
- **Lazy Loading** - Images and components loaded on demand
- **Caching** - API responses cached where appropriate
- **Compression** - Gzip compression in Nginx
- **CDN-Ready** - Static assets with cache headers

## Security Considerations

- **XSS Protection** - React's built-in XSS protection
- **CSRF** - Token-based authentication eliminates CSRF risks
- **Secure Headers** - Added in Nginx configuration
- **Input Validation** - Client-side validation for all forms
- **Sanitization** - All user inputs are sanitized

## Browser Support

- Chrome/Edge (latest)
- Firefox (latest)
- Safari (latest)
- Mobile browsers (iOS Safari, Chrome Mobile)

## Testing

The frontend includes comprehensive test suite using **Vitest** and **React Testing Library**.

### Running Tests

```bash
cd frontend
npm install
npm test
```

### Test Coverage

**✓ 47 tests passed across 10 test suites:**

- ✓ App component (1 test)
- ✓ LoadingSpinner component (4 tests)
- ✓ LanguageSwitcher component (4 tests)
- ✓ ProtectedRoute component (4 tests)
- ✓ AdminPage component (8 tests)
- ✓ API configuration (6 tests)
- ✓ AuthContext (3 tests)
- ✓ AuthService (5 tests)
- ✓ Translation files (9 tests)
- ✓ Date formatting (3 tests)

### What's Tested

- Component rendering and props
- User interactions and events
- Authentication flow and state management
- Protected routes and role-based access
- Language switching (EN/RU)
- API endpoint configuration
- Service methods and localStorage
- Translation completeness

### Test Commands

```bash
npm test              # Run all tests
npm run test:watch    # Watch mode
npm run test:ui       # Visual test runner
npm run test:coverage # Coverage report
```

## Contributing

When contributing to the frontend:

1. Follow the existing code style
2. Write meaningful component and variable names
3. Test on multiple screen sizes
4. Write tests for new features
5. Ensure all tests pass
6. Ensure no console errors or warnings
7. Update documentation as needed

## Future Enhancements

Potential improvements:

- [ ] Add user profile management
- [ ] Implement real-time notifications
- [ ] Add parking place ratings and reviews
- [ ] Integrate maps for parking location
- [ ] Add payment integration
- [ ] Implement advanced search filters
- [ ] Add booking history analytics
- [x] Support for multiple languages (EN/RU)
- [ ] Dark mode support
- [ ] Progressive Web App (PWA) features

## License

This frontend is part of the ParkingNet project. See the main project README for license information.

## Support

For issues or questions:

1. Check this documentation
2. Review the main project README
3. Check backend API documentation
4. Review browser console for errors
5. Check Docker container logs if using containers

## Resources

- [React Documentation](https://react.dev)
- [Vite Documentation](https://vitejs.dev)
- [Vitest Documentation](https://vitest.dev)
- [React Testing Library](https://testing-library.com/react)
- [Tailwind CSS Documentation](https://tailwindcss.com)
- [React Router Documentation](https://reactrouter.com)
- [Axios Documentation](https://axios-http.com)
