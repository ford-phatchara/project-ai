import { ChangeDetectionStrategy, Component, DestroyRef, inject, signal } from '@angular/core';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { NavigationEnd, Router, RouterOutlet } from '@angular/router';
import { filter } from 'rxjs';

import { BottomNavComponent } from '../bottom-nav/bottom-nav.component';
import { QUICK_ACTIONS, QuickAction } from '../navigation';
import { SidebarComponent } from '../sidebar/sidebar.component';
import { TopbarComponent } from '../topbar/topbar.component';
import { FabComponent } from '../../shared/components/fab/fab.component';

const ROUTE_TITLES: Record<string, string> = {
  '/dashboard': 'Dashboard',
  '/sales': 'Sales Management',
  '/plots': 'Plot Management',
  '/maintenance': 'Maintenance / Care Logs',
  '/calendar': 'Calendar View',
};

@Component({
  selector: 'app-shell',
  standalone: true,
  imports: [
    RouterOutlet,
    SidebarComponent,
    BottomNavComponent,
    TopbarComponent,
    FabComponent,
  ],
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <div class="min-h-screen">
      <app-sidebar />

      <div class="min-h-screen lg:pl-72">
        <app-topbar [title]="title()" />

        <main class="mx-auto w-full max-w-7xl px-4 py-5 pb-28 sm:px-6 lg:px-8 lg:pb-10">
          <router-outlet />
        </main>
      </div>

      <app-bottom-nav />
      <app-fab [items]="quickActions" (selected)="openQuickAction($event)" />
    </div>
  `,
})
export class AppShellComponent {
  private readonly router = inject(Router);
  private readonly destroyRef = inject(DestroyRef);

  readonly title = signal(this.getTitle(this.router.url));
  readonly quickActions = QUICK_ACTIONS;

  constructor() {
    this.router.events
      .pipe(
        filter((event): event is NavigationEnd => event instanceof NavigationEnd),
        takeUntilDestroyed(this.destroyRef),
      )
      .subscribe((event) => {
        this.title.set(this.getTitle(event.urlAfterRedirects));
      });
  }

  openQuickAction(action: QuickAction): void {
    void this.router.navigate([action.path], {
      queryParams: { action: 'add' },
    });
  }

  private getTitle(url: string): string {
    const path = url.split('?')[0] || '/dashboard';
    return ROUTE_TITLES[path] ?? 'Durian Farm Manager';
  }
}
