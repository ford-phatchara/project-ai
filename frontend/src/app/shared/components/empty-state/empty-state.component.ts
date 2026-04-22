import { ChangeDetectionStrategy, Component, input } from '@angular/core';

@Component({
  selector: 'app-empty-state',
  standalone: true,
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <div
      class="flex min-h-48 flex-col items-center justify-center rounded-xl border border-dashed border-stone-300 bg-white/70 px-6 py-10 text-center"
    >
      <div class="flex size-12 items-center justify-center rounded-xl bg-emerald-50 text-emerald-700">
        <svg class="size-6" viewBox="0 0 24 24" fill="none" aria-hidden="true">
          <path
            d="M6.75 5.75h10.5M6.75 9.25h10.5M6.75 12.75h6.5M5 19.25h14A2.25 2.25 0 0 0 21.25 17V5A2.25 2.25 0 0 0 19 2.75H5A2.25 2.25 0 0 0 2.75 5v12A2.25 2.25 0 0 0 5 19.25Z"
            stroke="currentColor"
            stroke-width="1.6"
            stroke-linecap="round"
            stroke-linejoin="round"
          />
        </svg>
      </div>
      <h3 class="mt-4 text-base font-semibold text-stone-950">{{ title() }}</h3>
      <p class="mt-2 max-w-sm text-sm leading-6 text-stone-500">{{ message() }}</p>
      <ng-content />
    </div>
  `,
})
export class EmptyStateComponent {
  readonly title = input.required<string>();
  readonly message = input.required<string>();
}
