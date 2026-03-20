# Frontend - EASM Dashboard UI

Modern, responsive React dashboard for managing assets and EASM scan jobs.

## 🎨 Features

### Dashboard Components

- **Statistics Panel**: Real-time asset counts by type and status
- **Asset Management**: 
  - Create new assets with name and type selector
  - Delete assets with confirmation
  - View full asset list with pagination
  - Search assets by name
  - Filter by type and status

- **Scan Control**:
  - Select asset and scan type
  - Start scan jobs asynchronously
  - Monitor scan job status in real-time
  - View detailed scan results as JSON
  - Support for all 10 scan types

### Design Highlights

- **Responsive**: Mobile-first CSS with breakpoints at 980px and 640px
- **Modern Styling**: Gradient theme with warm earth tones
- **Clean UI**: No external CSS frameworks, custom implementations
- **Error Handling**: User feedback for all API operations
- **Loading States**: Visual indicators during async operations

## 📋 Requirements

- **Node.js**: 16+ (with npm 7+)
- **Backend API**: Running at configured URL (default: `http://localhost:8080`)

## 🚀 Quick Start

### 1. Environment Configuration

Copy `.env.example` to `.env.local`:

```bash
cp .env.example .env.local
```

Edit `.env.local`:
```env
VITE_API_URL=http://localhost:8080
```

For production builds:
```env
VITE_API_URL=https://api.example.com
```

### 2. Install Dependencies

```bash
npm install
```

### 3. Development Server

```bash
npm run dev
```

Open browser: `http://localhost:5173`

- Hot Module Replacement (HMR) enabled
- Auto-refresh on code changes
- Development errors shown in browser

### 4. Production Build

```bash
npm run build
```

Output:
- `dist/index.html` - Main HTML
- `dist/assets/index-*.js` - Bundled JavaScript
- `dist/assets/index-*.css` - Bundled CSS

Size:
- JS: ~199 KB (62 KB gzipped)
- CSS: ~4 KB (1.5 KB gzipped)

### 5. Preview Production Build

```bash
npm run preview
```

Opens `http://localhost:4173`

## 📁 Project Structure

```
frontend/
├── src/
│   ├── main.jsx              # React entry point
│   ├── App.jsx               # Dashboard component (600+ LOC)
│   └── styles.css            # Styling system (2400+ LOC)
├── index.html                # HTML template
├── vite.config.js            # Vite build config
├── package.json              # Dependencies & scripts
├── .env.example              # Environment template
├── eslint.config.js          # Linting rules
├── public/                   # Static assets
├── dist/                     # Built output (after npm run build)
├── node_modules/             # Dependencies
└── README.md                 # This file
```

## 🔌 API Integration

### Configuration

The application reads `VITE_API_URL` from environment variables:

```javascript
// src/App.jsx
const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';
```

### Endpoints Used

| Endpoint | Method | Purpose |
|----------|--------|---------|
| `/assets/stats` | GET | Fetch statistics |
| `/assets` | GET | List assets with pagination |
| `/assets` | POST | Create new asset |
| `/assets/{id}` | DELETE | Delete asset |
| `/assets/search` | GET | Search assets by name |
| `/assets/{id}/scan` | POST | Start scan job |
| `/scan-jobs/{id}` | GET | Get job status |
| `/scan-jobs/{id}/results` | GET | Fetch scan results |

### Error Handling

API errors are caught and displayed to user:

```javascript
try {
  const response = await api('/endpoint');
} catch (error) {
  setError(error.message);
  // User sees error message in UI
}
```

## 🎯 Usage Examples

### Create Asset

1. Enter asset name (e.g., `example.com`)
2. Select type (domain, ip, hostname)
3. Click "Create Asset"
4. Stats update automatically

### Run Scan

1. Select asset from "Scan Control" dropdown
2. Choose scan type (dns, whois, all, etc.)
3. Click "Start Scan"
4. Monitor status in "Recent Scans" table
5. Click "View Results" to see JSON output

### Search Assets

1. Type in search box (auto-searches as you type)
2. Results show matching assets
3. Use type/status dropdowns to filter further
4. Navigate with pagination controls

## 🧹 Code Quality

### Frontend Test Status

Current submission does not include frontend unit test files yet.
Validation currently uses lint + production build:

```bash
npm run lint
npm run build
```

### Linting

```bash
npm run lint
```

Checks:
- React-specific rules
- Best practices
- ESLint errors/warnings

### Formatting

Code formatted for readability:
- 2-space indentation
- 80-character line suggestions
- Trailing commas

## 🚢 Deployment

### Docker (with Nginx)

From root directory:

```bash
docker-compose up
```

Frontend served at: `http://localhost:3000`

### Static Hosting

Serve `dist/` folder with any static host:

```bash
# GitHub Pages
npm run build
# Push dist/ to gh-pages branch

# Vercel
npm run build
# Connect to Vercel, auto-deploys dist/

# Netlify
npm run build
# Drag dist/ folder to Netlify
```

### Environment-Specific Builds

```bash
# Development
VITE_API_URL=http://localhost:8080 npm run dev

# Staging
VITE_API_URL=https://staging-api.example.com npm run build

# Production
VITE_API_URL=https://api.example.com npm run build
```

## 🔄 Component Communication

The dashboard uses React hooks for state management:

- `useState` for component state
- `useEffect` for side effects (API calls)
- `useMemo` for derived values
- Props passing for child components

No external state management library needed.

## 📱 Browser Support

- Chrome/Edge 90+
- Firefox 88+
- Safari 14+
- Mobile browsers (iOS Safari, Chrome Mobile)

## 🐛 Troubleshooting

### API Connection Issues

1. Check backend is running: `curl http://localhost:8080/health`
2. Verify `VITE_API_URL` in `.env.local`
3. Check browser console for CORS errors
4. Ensure backend CORS middleware is enabled

### Build Size Too Large

1. Check for unused imports: `npm run lint`
2. Run `npm run build` and check dist/ output
3. Consider code-splitting for large components

### Port Already in Use

```bash
# Change dev port
VITE_PORT=5174 npm run dev

# Change preview port
VITE_PREVIEW_PORT=4174 npm run preview
```

---

**For backend API documentation, see [backend/README.md](../backend/README.md)**

**For full project guide and Day 3 assignments, see [root README.md](../README.md)**
