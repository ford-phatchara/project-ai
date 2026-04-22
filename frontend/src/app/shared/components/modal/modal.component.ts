import { ChangeDetectionStrategy, Component, input, output } from '@angular/core';

@Component({
  selector: 'app-modal',
  standalone: true,
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <div
      class="fixed inset-0 z-50 flex items-end justify-center bg-stone-950/35 px-4 py-4 backdrop-blur-sm sm:items-center"
      role="presentation"
      (click)="closed.emit()"
    >
      <section
        class="max-h-[92vh] w-full max-w-2xl overflow-hidden rounded-xl border border-stone-200 bg-white shadow-2xl shadow-stone-900/20"
        role="dialog"
        aria-modal="true"
        [attr.aria-label]="title()"
        (click)="$event.stopPropagation()"
      >
        <header class="flex items-center justify-between border-b border-stone-100 px-5 py-4">
          <div>
            <p class="text-xs font-semibold uppercase tracking-wide text-emerald-700">
              Durian Farm Manager
            </p>
            <h2 class="mt-1 text-lg font-semibold text-stone-950">{{ title() }}</h2>
          </div>

          <button
            type="button"
            class="rounded-full p-2 text-stone-500 transition hover:bg-stone-100 hover:text-stone-900 focus:outline-none focus:ring-4 focus:ring-stone-200"
            aria-label="Close dialog"
            (click)="closed.emit()"
          >
            <svg class="size-5" viewBox="0 0 24 24" fill="none" aria-hidden="true">
              <path
                d="M6 6l12 12M18 6 6 18"
                stroke="currentColor"
                stroke-width="1.8"
                stroke-linecap="round"
              />
            </svg>
          </button>
        </header>

        <div class="max-h-[calc(92vh-5.5rem)] overflow-y-auto px-5 py-5">
          <ng-content />
        </div>
      </section>
    </div>
  `,
})
export class ModalComponent {
  readonly title = input.required<string>();
  readonly closed = output<void>();
}
