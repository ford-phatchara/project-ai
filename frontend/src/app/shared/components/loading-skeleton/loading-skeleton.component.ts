import { ChangeDetectionStrategy, Component, input } from '@angular/core';

@Component({
  selector: 'app-loading-skeleton',
  standalone: true,
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <div class="space-y-3" aria-label="Loading content">
      @for (line of linesArray(); track line) {
        <div
          class="h-4 animate-pulse rounded-full bg-stone-200"
          [class.w-11/12]="line % 3 === 0"
          [class.w-3/4]="line % 3 === 1"
          [class.w-full]="line % 3 === 2"
        ></div>
      }
    </div>
  `,
})
export class LoadingSkeletonComponent {
  readonly lines = input(4);

  linesArray(): number[] {
    return Array.from({ length: this.lines() }, (_, index) => index);
  }
}
