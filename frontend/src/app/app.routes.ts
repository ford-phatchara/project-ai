import { Routes } from '@angular/router';

export const routes: Routes = [
  {
    path: 'login',
    loadComponent: () =>
      import('./features/auth/login-page/login-page.component').then(
        (m) => m.LoginPageComponent,
      ),
  },
  {
    path: '',
    loadComponent: () =>
      import('./layout/app-shell/app-shell.component').then(
        (m) => m.AppShellComponent,
      ),
    children: [
      {
        path: '',
        pathMatch: 'full',
        redirectTo: 'dashboard',
      },
      {
        path: 'dashboard',
        loadComponent: () =>
          import('./features/dashboard/dashboard-page/dashboard-page.component').then(
            (m) => m.DashboardPageComponent,
          ),
      },
      {
        path: 'sales',
        loadComponent: () =>
          import('./features/sales/sales-page/sales-page.component').then(
            (m) => m.SalesPageComponent,
          ),
      },
      {
        path: 'plots',
        loadComponent: () =>
          import('./features/plots/plots-page/plots-page.component').then(
            (m) => m.PlotsPageComponent,
          ),
      },
      {
        path: 'maintenance',
        loadComponent: () =>
          import(
            './features/maintenance/maintenance-page/maintenance-page.component'
          ).then((m) => m.MaintenancePageComponent),
      },
      {
        path: 'calendar',
        loadComponent: () =>
          import('./features/calendar/calendar-page/calendar-page.component').then(
            (m) => m.CalendarPageComponent,
          ),
      },
    ],
  },
  {
    path: '**',
    redirectTo: 'dashboard',
  },
];
