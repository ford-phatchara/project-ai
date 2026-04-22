import { ChangeDetectionStrategy, Component, input } from '@angular/core';

@Component({
  selector: 'app-section-card',
  standalone: true,
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <section class="surface overflow-hidden">
      <header class="border-b border-stone-100 px-4 py-4 sm:px-5">
        <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <h2 class="text-base font-semibold text-stone-950">{{ title() }}</h2>
            <p class="mt-1 text-sm text-stone-500">{{ subtitle() }}</p>
          </div>
          <ng-content select="[section-action]" />
        </div>
      </header>
      <div class="p-4 sm:p-5">
        <ng-content />
      </div>
    </section>
  `,
})
export class SectionCardComponent {
  readonly title = input.required<string>();
  readonly subtitle = input<string>('');
}
