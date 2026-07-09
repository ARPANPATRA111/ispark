# ISPARC — IIPS Skill, Personality Advancement & Refinement Cell

Official student development and academic portal of the International Institute of Professional Studies (IIPS), DAVV, Indore.

## Repository Architecture

This project is structured as a monorepo containing both the SvelteKit web application and the Go Fiber API backend:

```
root/
├── web/                      # SvelteKit Frontend (Svelte 5 + TypeScript + Tailwind v4)
│   ├── src/
│   ├── static/
│   ├── package.json
│   └── vite.config.ts
│
├── api/                      # Go Fiber Backend API (Future Implementation)
│   └── .gitkeep
│
├── docker-compose.yml        # Docker Multi-container Orchestration
└── README.md                 # Main Documentation
```

## Getting Started

### Prerequisites
- Node.js (v20+)
- npm (v10+)
- Docker & Docker Compose (optional, for deployment)

### Development Setup (Frontend)

1. Navigate to the `web` directory:
   ```bash
   cd web
   ```

2. Install dependencies:
   ```bash
   npm install
   ```

3. Start the Vite development server:
   ```bash
   npm run dev
   ```

4. Open [http://localhost:5173](http://localhost:5173) in your browser.

### Production Build (Frontend)
To bundle the frontend application for production:
```bash
cd web
npm run build
npm run preview
```

## Backend Team Integration Notes
- The `/api` folder is reserved for the Go Fiber server.
- The root `docker-compose.yml` configures ports:
  - Frontend: `3000`
  - Backend API: `8080`
  - PostgreSQL Database: `5432`
- Ensure CORS configurations in the Go backend allow requests from port `3000` and development port `5173`.
