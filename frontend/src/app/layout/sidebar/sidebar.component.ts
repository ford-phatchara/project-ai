import { NgTemplateOutlet } from '@angular/common';
import { ChangeDetectionStrategy, Component } from '@angular/core';
import { RouterLink, RouterLinkActive } from '@angular/router';

import { NAV_ITEMS } from '../navigation';

@Component({
  selector: 'app-sidebar',
  standalone: true,
  imports: [NgTemplateOutlet, RouterLink, RouterLinkActive],
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <aside
      class="fixed inset-y-0 left-0 z-30 hidden w-72 border-r border-stone-200 bg-white/90 px-4 py-5 shadow-sm shadow-stone-900/5 backdrop-blur lg:block"
    >
      <a routerLink="/dashboard" class="flex items-center gap-3 rounded-xl px-2 py-2">
        <span class="flex size-10 items-center justify-center rounded-xl bg-emerald-700 text-white shadow-sm">
          <svg class="size-5" viewBox="0 0 24 24" fill="none" aria-hidden="true">
            <path
              d="M12 3.5c4.25 2.1 6.5 5.03 6.5 8.47 0 4.57-3.4 8.53-6.5 8.53s-6.5-3.96-6.5-8.53C5.5 8.53 7.75 5.6 12 3.5Z"
              stroke="currentColor"
              stroke-width="1.6"
              stroke-linejoin="round"
            />
            <path
              d="M12 7.5v9"
              stroke="currentColor"
              stroke-width="1.6"
              stroke-linecap="round"
            />
          </svg>
        </span>
        <div>
          <p class="text-sm font-semibold text-stone-950">Durian Farm</p>
          <p class="text-xs text-stone-500">Manager</p>
        </div>
      </a>

      <nav class="mt-8 space-y-1">
        @for (item of navItems; track item.path) {
          <a
            [routerLink]="item.path"
            routerLinkActive="bg-emerald-50 text-emerald-800 ring-1 ring-emerald-100"
            class="group flex items-center gap-3 rounded-xl px-3 py-2.5 text-sm font-medium text-stone-600 transition hover:bg-stone-50 hover:text-stone-950"
          >
            <span class="text-stone-400 transition group-hover:text-emerald-700">
              <ng-container [ngTemplateOutlet]="icon" [ngTemplateOutletContext]="{ name: item.icon }" />
            </span>
            {{ item.label }}
          </a>
        }
      </nav>

      <div class="absolute inset-x-4 bottom-5 rounded-xl border border-emerald-100 bg-emerald-50 p-4">
        <p class="text-sm font-semibold text-emerald-950">2026 harvest season</p>
        <p class="mt-1 text-sm leading-6 text-emerald-800">
          Track fruit sales and care work from one calm workspace.
        </p>
      </div>
    </aside>

    <ng-template #icon let-name="name">
      <svg class="size-5" viewBox="0 0 24 24" fill="none" aria-hidden="true">
        @switch (name) {
          @case ('dashboard') {
            <path d="M4 13h6V4H4v9Zm10 7h6V4h-6v16ZM4 20h6v-3H4v3Z" stroke="currentColor" stroke-width="1.6" stroke-linejoin="round" />
          }
          @case ('sales') {
            <path d="M5 19V6.5A1.5 1.5 0 0 1 6.5 5h11A1.5 1.5 0 0 1 19 6.5V19l-2.5-1.5L14 19l-2-1.5L9.5 19 7 17.5 5 19Z" stroke="currentColor" stroke-width="1.6" stroke-linejoin="round" />
            <path d="M8 9h8M8 12h8M8 15h4" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" />
          }
          @case ('plots') {
            <path d="M4 7.5 12 4l8 3.5-8 3.5L4 7.5Zm0 5 8 3.5 8-3.5M4 17.5 12 21l8-3.5" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round" />
          }
          @case ('maintenance') {
            <path d="M12 4.5c3 1.65 4.5 3.87 4.5 6.1 0 3.13-2.08 5.4-4.5 5.4s-4.5-2.27-4.5-5.4c0-2.23 1.5-4.45 4.5-6.1Z" stroke="currentColor" stroke-width="1.6" stroke-linejoin="round" />
            <path d="M7 19h10M12 9v7" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" />
          }
          @case ('calendar') {
            <path d="M7 4v3M17 4v3M5.5 8.5h13M7 20h10a2 2 0 0 0 2-2V7.5a2 2 0 0 0-2-2H7a2 2 0 0 0-2 2V18a2 2 0 0 0 2 2Z" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" />
          }
        }
      </svg>
    </ng-template>
  `,
})
export class SidebarComponent {
  readonly navItems = NAV_ITEMS;
}
