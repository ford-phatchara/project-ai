import { ChangeDetectionStrategy, Component } from '@angular/core';
import { RouterLink, RouterLinkActive } from '@angular/router';

import { NAV_ITEMS } from '../navigation';

@Component({
  selector: 'app-bottom-nav',
  standalone: true,
  imports: [RouterLink, RouterLinkActive],
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <nav
      class="fixed inset-x-0 bottom-0 z-30 border-t border-stone-200 bg-white/95 px-2 pb-[max(env(safe-area-inset-bottom),0.4rem)] pt-2 shadow-[0_-12px_30px_rgba(28,25,23,0.08)] backdrop-blur lg:hidden"
      aria-label="Primary"
    >
      <div class="mx-auto grid max-w-xl grid-cols-5 gap-1">
        @for (item of navItems; track item.path) {
          <a
            [routerLink]="item.path"
            routerLinkActive="text-emerald-800"
            class="flex min-h-14 flex-col items-center justify-center gap-1 rounded-xl px-1 text-xs font-medium text-stone-500 transition hover:bg-stone-50 hover:text-stone-900"
          >
            <span class="flex size-7 items-center justify-center rounded-lg bg-stone-50 text-stone-500">
              <svg class="size-4" viewBox="0 0 24 24" fill="none" aria-hidden="true">
                @switch (item.icon) {
                  @case ('dashboard') {
                    <path d="M4 13h6V4H4v9Zm10 7h6V4h-6v16ZM4 20h6v-3H4v3Z" stroke="currentColor" stroke-width="1.7" stroke-linejoin="round" />
                  }
                  @case ('sales') {
                    <path d="M5 19V6.5A1.5 1.5 0 0 1 6.5 5h11A1.5 1.5 0 0 1 19 6.5V19l-2.5-1.5L14 19l-2-1.5L9.5 19 7 17.5 5 19Z" stroke="currentColor" stroke-width="1.7" stroke-linejoin="round" />
                  }
                  @case ('plots') {
                    <path d="M4 7.5 12 4l8 3.5-8 3.5L4 7.5Zm0 5 8 3.5 8-3.5M4 17.5 12 21l8-3.5" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round" />
                  }
                  @case ('maintenance') {
                    <path d="M12 4.5c3 1.65 4.5 3.87 4.5 6.1 0 3.13-2.08 5.4-4.5 5.4s-4.5-2.27-4.5-5.4c0-2.23 1.5-4.45 4.5-6.1Z" stroke="currentColor" stroke-width="1.7" stroke-linejoin="round" />
                  }
                  @case ('calendar') {
                    <path d="M7 4v3M17 4v3M5.5 8.5h13M7 20h10a2 2 0 0 0 2-2V7.5a2 2 0 0 0-2-2H7a2 2 0 0 0-2 2V18a2 2 0 0 0 2 2Z" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" />
                  }
                }
              </svg>
            </span>
            <span class="max-w-full truncate">{{ item.label }}</span>
          </a>
        }
      </div>
    </nav>
  `,
})
export class BottomNavComponent {
  readonly navItems = NAV_ITEMS;
}
