export interface NavItem {
  label: string;
  path: string;
  icon: 'dashboard' | 'sales' | 'plots' | 'maintenance' | 'calendar';
}

export const NAV_ITEMS: NavItem[] = [
  {
    label: 'Dashboard',
    path: '/dashboard',
    icon: 'dashboard',
  },
  {
    label: 'Sales',
    path: '/sales',
    icon: 'sales',
  },
  {
    label: 'Plots',
    path: '/plots',
    icon: 'plots',
  },
  {
    label: 'Care Logs',
    path: '/maintenance',
    icon: 'maintenance',
  },
  {
    label: 'Calendar',
    path: '/calendar',
    icon: 'calendar',
  },
];

export interface QuickAction {
  label: string;
  path: string;
  action: 'sale' | 'plot' | 'maintenance';
}

export const QUICK_ACTIONS: QuickAction[] = [
  {
    label: 'Sale',
    path: '/sales',
    action: 'sale',
  },
  {
    label: 'Plot',
    path: '/plots',
    action: 'plot',
  },
  {
    label: 'Care Log',
    path: '/maintenance',
    action: 'maintenance',
  },
];
