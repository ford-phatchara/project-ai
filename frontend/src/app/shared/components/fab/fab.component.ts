import { ChangeDetectionStrategy, Component, input, output, signal } from '@angular/core';

import { QuickAction } from '../../../layout/navigation';

@Component({
  selector: 'app-fab',
  standalone: true,
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <div class="fixed bottom-20 right-5 z-40 lg:hidden">
      @if (open()) {
        <div class="mb-3 flex flex-col items-end gap-2">
          @for (item of items(); track item.action) {
            <button
              type="button"
              class="inline-flex items-center gap-2 rounded-xl border border-stone-200 bg-white px-3.5 py-2 text-sm font-semibold text-stone-800 shadow-lg shadow-stone-900/10 transition hover:border-emerald-200 hover:text-emerald-800"
              (click)="choose(item)"
            >
              <span class="flex size-7 items-center justify-center rounded-lg bg-emerald-50 text-emerald-700">
                <svg class="size-4" viewBox="0 0 24 24" fill="none" aria-hidden="true">
                  <path
                    d="M12 5v14M5 12h14"
                    stroke="currentColor"
                    stroke-width="1.8"
                    stroke-linecap="round"
                  />
                </svg>
              </span>
              {{ item.label }}
            </button>
          }
        </div>
      }

      <button
        type="button"
        class="flex size-14 items-center justify-center rounded-full bg-emerald-700 text-white shadow-xl shadow-emerald-950/20 transition hover:bg-emerald-800 focus:outline-none focus:ring-4 focus:ring-emerald-200"
        aria-label="Quick add"
        (click)="open.update((value) => !value)"
      >
        <svg
          class="size-6 transition"
          [class.rotate-45]="open()"
          viewBox="0 0 24 24"
          fill="none"
          aria-hidden="true"
        >
          <path
            d="M12 5v14M5 12h14"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
          />
        </svg>
      </button>
    </div>
  `,
})
export class FabComponent {
  readonly items = input<QuickAction[]>([]);
  readonly selected = output<QuickAction>();
  readonly open = signal(false);

  choose(item: QuickAction): void {
    this.open.set(false);
    this.selected.emit(item);
  }
}
