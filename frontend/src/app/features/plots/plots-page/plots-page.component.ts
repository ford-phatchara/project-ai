import { DecimalPipe } from '@angular/common';
import { ChangeDetectionStrategy, Component, computed, inject, signal } from '@angular/core';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { ActivatedRoute, Router } from '@angular/router';

import { Plot, PlotFormValue } from '../../../core/models/plot.model';
import { FarmDataService } from '../../../core/services/farm-data.service';
import { EmptyStateComponent } from '../../../shared/components/empty-state/empty-state.component';
import { ModalComponent } from '../../../shared/components/modal/modal.component';
import { PlotFormComponent } from '../plot-form/plot-form.component';

@Component({
  selector: 'app-plots-page',
  standalone: true,
  imports: [DecimalPipe, EmptyStateComponent, ModalComponent, PlotFormComponent],
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <div class="space-y-5">
      <section class="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <p class="text-sm font-medium text-emerald-700">Plot management</p>
          <h2 class="mt-1 text-2xl font-semibold tracking-tight text-stone-950">
            Farm plots
          </h2>
          <p class="mt-2 max-w-2xl text-sm leading-6 text-stone-500">
            Keep track of plot size, growing notes, and active orchard blocks.
          </p>
        </div>
        <button class="btn-primary" type="button" (click)="openCreate()">
          <svg class="size-4" viewBox="0 0 24 24" fill="none" aria-hidden="true">
            <path d="M12 5v14M5 12h14" stroke="currentColor" stroke-width="2" stroke-linecap="round" />
          </svg>
          Add plot
        </button>
      </section>

      @if (plots().length) {
        <section class="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
          @for (plot of plots(); track plot.id) {
            <article class="surface p-5">
              <div class="flex items-start justify-between gap-4">
                <div>
                  <p class="text-sm font-medium text-emerald-700">Active plot</p>
                  <h3 class="mt-1 text-lg font-semibold text-stone-950">{{ plot.name }}</h3>
                </div>
                <div class="rounded-xl bg-amber-50 px-3 py-2 text-right ring-1 ring-amber-100">
                  <p class="text-lg font-semibold text-amber-900">{{ plot.size | number: '1.0-1' }}</p>
                  <p class="text-xs font-medium uppercase tracking-wide text-amber-700">rai</p>
                </div>
              </div>

              <p class="mt-4 min-h-16 text-sm leading-6 text-stone-600">{{ plot.notes || 'No notes added.' }}</p>

              <div class="mt-5 grid grid-cols-2 gap-2">
                <button class="btn-secondary" type="button" (click)="openEdit(plot)">Edit</button>
                <button class="btn-danger" type="button" (click)="deletePlot(plot)">Delete</button>
              </div>
            </article>
          }
        </section>
      } @else {
        <app-empty-state
          title="No plots yet"
          message="Create your first farm plot before adding sales or care logs."
        >
          <button class="btn-primary mt-5" type="button" (click)="openCreate()">Add plot</button>
        </app-empty-state>
      }
    </div>

    @if (dialogOpen()) {
      <app-modal [title]="editingPlot() ? 'Edit plot' : 'Add plot'" (closed)="closeDialog()">
        <app-plot-form
          [plot]="editingPlot()"
          (saved)="savePlot($event)"
          (cancel)="closeDialog()"
        />
      </app-modal>
    }
  `,
})
export class PlotsPageComponent {
  private readonly farmData = inject(FarmDataService);
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);

  readonly plots = computed(() =>
    [...this.farmData.plots()].sort((a, b) => a.name.localeCompare(b.name)),
  );
  readonly dialogOpen = signal(false);
  readonly editingPlot = signal<Plot | null>(null);

  constructor() {
    this.route.queryParamMap.pipe(takeUntilDestroyed()).subscribe((params) => {
      if (params.get('action') === 'add') {
        this.openCreate();
        void this.router.navigate([], {
          relativeTo: this.route,
          queryParams: {},
          replaceUrl: true,
        });
      }
    });
  }

  openCreate(): void {
    this.editingPlot.set(null);
    this.dialogOpen.set(true);
  }

  openEdit(plot: Plot): void {
    this.editingPlot.set(plot);
    this.dialogOpen.set(true);
  }

  closeDialog(): void {
    this.dialogOpen.set(false);
    this.editingPlot.set(null);
  }

  savePlot(value: PlotFormValue): void {
    if (value.id) {
      this.farmData.updatePlot(value);
    } else {
      this.farmData.addPlot(value);
    }

    this.closeDialog();
  }

  deletePlot(plot: Plot): void {
    const confirmed = window.confirm(
      `Delete ${plot.name}? Related sales and care logs will also be removed.`,
    );

    if (confirmed) {
      this.farmData.deletePlot(plot.id);
    }
  }
}
