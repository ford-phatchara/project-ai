import { ChangeDetectionStrategy, Component, input } from '@angular/core';

@Component({
  selector: 'app-chart-card',
  standalone: true,
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <section class="surface p-4 sm:p-5">
      <div class="mb-4 flex items-start justify-between gap-4">
        <div>
          <h2 class="text-base font-semibold text-stone-950">{{ title() }}</h2>
          <p class="mt-1 text-sm text-stone-500">{{ subtitle() }}</p>
        </div>
        <ng-content select="[card-action]" />
      </div>
      <ng-content />
    </section>
  `,
})
export class ChartCardComponent {
  readonly title = input.required<string>();
  readonly subtitle = input<string>('');
}
