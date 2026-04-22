import { ChangeDetectionStrategy, Component, input } from '@angular/core';

@Component({
  selector: 'app-stat-card',
  standalone: true,
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <article class="surface p-4 sm:p-5">
      <div class="flex items-start justify-between gap-4">
        <div>
          <p class="text-sm font-medium text-stone-500">{{ label() }}</p>
          <p class="mt-2 text-2xl font-semibold tracking-tight text-stone-950">
            {{ value() }}
          </p>
        </div>
        <div
          class="flex size-11 items-center justify-center rounded-xl bg-emerald-50 text-emerald-700 ring-1 ring-emerald-100"
        >
          <ng-content />
        </div>
      </div>
      <p class="mt-4 text-sm text-stone-500">{{ caption() }}</p>
    </article>
  `,
})
export class StatCardComponent {
  readonly label = input.required<string>();
  readonly value = input.required<string>();
  readonly caption = input<string>('');
}
