import { ChangeDetectionStrategy, Component, input } from '@angular/core';
import { RouterLink } from '@angular/router';

@Component({
  selector: 'app-topbar',
  standalone: true,
  imports: [RouterLink],
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <header
      class="sticky top-0 z-20 border-b border-stone-200/80 bg-white/85 px-4 py-3 backdrop-blur lg:px-8"
    >
      <div class="flex items-center justify-between gap-4">
        <div class="flex min-w-0 items-center gap-3">
          <a
            routerLink="/dashboard"
            class="flex size-10 shrink-0 items-center justify-center rounded-xl bg-emerald-700 text-white shadow-sm lg:hidden"
            aria-label="Dashboard"
          >
            <svg class="size-5" viewBox="0 0 24 24" fill="none" aria-hidden="true">
              <path
                d="M12 3.5c4.25 2.1 6.5 5.03 6.5 8.47 0 4.57-3.4 8.53-6.5 8.53s-6.5-3.96-6.5-8.53C5.5 8.53 7.75 5.6 12 3.5Z"
                stroke="currentColor"
                stroke-width="1.6"
                stroke-linejoin="round"
              />
              <path d="M12 7.5v9" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" />
            </svg>
          </a>
          <div class="min-w-0">
            <p class="text-xs font-semibold uppercase tracking-wide text-emerald-700 lg:hidden">
              Durian Farm Manager
            </p>
            <h1 class="truncate text-lg font-semibold text-stone-950 sm:text-xl">
              {{ title() }}
            </h1>
          </div>
        </div>

        <div class="flex items-center gap-2">
          <span
            class="hidden rounded-full border border-emerald-100 bg-emerald-50 px-3 py-1.5 text-sm font-medium text-emerald-800 sm:inline-flex"
          >
            Live mock data
          </span>
          <a
            routerLink="/login"
            class="inline-flex size-10 items-center justify-center rounded-full border border-stone-200 bg-white text-stone-600 transition hover:bg-stone-50 hover:text-stone-950"
            aria-label="Log out"
          >
            <svg class="size-5" viewBox="0 0 24 24" fill="none" aria-hidden="true">
              <path
                d="M15 7.5V6a2 2 0 0 0-2-2H6.5a2 2 0 0 0-2 2v12a2 2 0 0 0 2 2H13a2 2 0 0 0 2-2v-1.5M10 12h9m0 0-3-3m3 3-3 3"
                stroke="currentColor"
                stroke-width="1.7"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
            </svg>
          </a>
        </div>
      </div>
    </header>
  `,
})
export class TopbarComponent {
  readonly title = input.required<string>();
}
