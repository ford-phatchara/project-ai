# Durian Farm Manager

Frontend-only Angular application for managing durian farm sales, plots, care logs,
and maintenance calendar activity.

## Stack

- Angular 21 with standalone components
- Angular Router with `src/app/app.routes.ts`
- TypeScript
- Tailwind CSS 4
- Chart.js
- Local mock data and signal-backed frontend services

## Run Locally

```bash
npm install
npm start
```

Open `http://127.0.0.1:4200/`.

## Useful Scripts

```bash
npm run build
npm test -- --watch=false
```

## App Routes

- `/login`
- `/dashboard`
- `/sales`
- `/plots`
- `/maintenance`
- `/calendar`

The login screen is UI-only. Submitting the form navigates to `/dashboard`.
