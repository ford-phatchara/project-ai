import { DatePipe, DecimalPipe } from '@angular/common';
import { ChangeDetectionStrategy, Component, computed, inject, signal } from '@angular/core';
import { takeUntilDestroyed, toSignal } from '@angular/core/rxjs-interop';
import { NonNullableFormBuilder, ReactiveFormsModule } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { startWith } from 'rxjs';

import {
  MaintenanceLog,
  MaintenanceLogFormValue,
} from '../../../core/models/maintenance-log.model';
import { FarmDataService } from '../../../core/services/farm-data.service';
import { EmptyStateComponent } from '../../../shared/components/empty-state/empty-state.component';
import { ModalComponent } from '../../../shared/components/modal/modal.component';
import { SectionCardComponent } from '../../../shared/components/section-card/section-card.component';
import { MaintenanceFormComponent } from '../maintenance-form/maintenance-form.component';

@Component({
  selector: 'app-maintenance-page',
  standalone: true,
  imports: [
    ReactiveFormsModule,
    DatePipe,
    DecimalPipe,
    EmptyStateComponent,
    ModalComponent,
    SectionCardComponent,
    MaintenanceFormComponent,
  ],
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <div class="space-y-5">
      <section class="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <p class="text-sm font-medium text-emerald-700">Farm maintenance</p>
          <h2 class="mt-1 text-2xl font-semibold tracking-tight text-stone-950">
            Care logs
          </h2>
          <p class="mt-2 max-w-2xl text-sm leading-6 text-stone-500">
            Track watering, fertilizing, and notes by plot and date.
          </p>
        </div>
        <button class="btn-primary" type="button" (click)="openCreate()">
          <svg class="size-4" viewBox="0 0 24 24" fill="none" aria-hidden="true">
            <path d="M12 5v14M5 12h14" stroke="currentColor" stroke-width="2" stroke-linecap="round" />
          </svg>
          Add care log
        </button>
      </section>

      <section class="surface p-4">
        <form class="grid gap-3 sm:grid-cols-[1fr_1fr_auto]" [formGroup]="filters">
          <label class="block space-y-2">
            <span class="form-label">Plot</span>
            <select class="form-field" formControlName="plot">
              <option value="">All plots</option>
              @for (plot of plots(); track plot.id) {
                <option [value]="plot.name">{{ plot.name }}</option>
              }
            </select>
          </label>

          <label class="block space-y-2">
            <span class="form-label">Date</span>
            <input class="form-field" type="date" formControlName="date" />
          </label>

          <div class="flex items-end">
            <button class="btn-secondary w-full sm:w-auto" type="button" (click)="clearFilters()">Clear filters</button>
          </div>
        </form>
      </section>

      <app-section-card
        title="Care activity"
        [subtitle]="filteredLogs().length + ' logs shown'"
      >
        @if (filteredLogs().length) {
          <div class="hidden overflow-x-auto md:block">
            <table class="min-w-full border-separate border-spacing-0 text-sm">
              <thead>
                <tr class="table-header">
                  <th class="rounded-l-xl px-4 py-3">Date</th>
                  <th class="px-4 py-3">Plot</th>
                  <th class="px-4 py-3 text-right">Watering</th>
                  <th class="px-4 py-3">Fertilizing</th>
                  <th class="px-4 py-3">Notes</th>
                  <th class="rounded-r-xl px-4 py-3 text-right">Actions</th>
                </tr>
              </thead>
              <tbody>
                @for (log of filteredLogs(); track log.id) {
                  <tr>
                    <td class="px-4 py-4 text-stone-600">{{ log.date | date: 'mediumDate' }}</td>
                    <td class="px-4 py-4 font-medium text-stone-950">{{ log.plot }}</td>
                    <td class="px-4 py-4 text-right text-stone-700">{{ log.wateringMinutes | number: '1.0-0' }} min</td>
                    <td class="px-4 py-4">
                      <span
                        class="rounded-full px-2.5 py-1 text-xs font-semibold"
                        [class.bg-emerald-50]="log.fertilizing"
                        [class.text-emerald-800]="log.fertilizing"
                        [class.bg-stone-100]="!log.fertilizing"
                        [class.text-stone-600]="!log.fertilizing"
                      >
                        {{ log.fertilizing ? 'Yes' : 'No' }}
                      </span>
                    </td>
                    <td class="max-w-sm px-4 py-4 text-stone-600">{{ log.notes }}</td>
                    <td class="px-4 py-4">
                      <div class="flex justify-end gap-2">
                        <button class="btn-secondary px-3 py-2" type="button" (click)="openEdit(log)">Edit</button>
                        <button class="btn-danger" type="button" (click)="deleteLog(log)">Delete</button>
                      </div>
                    </td>
                  </tr>
                }
              </tbody>
            </table>
          </div>

          <div class="space-y-3 md:hidden">
            @for (log of filteredLogs(); track log.id) {
              <article class="rounded-xl border border-stone-200 bg-white p-4 shadow-sm">
                <div class="flex items-start justify-between gap-3">
                  <div>
                    <p class="font-semibold text-stone-950">{{ log.plot }}</p>
                    <p class="mt-1 text-sm text-stone-500">{{ log.date | date: 'mediumDate' }}</p>
                  </div>
                  <span
                    class="rounded-full px-2.5 py-1 text-xs font-semibold"
                    [class.bg-emerald-50]="log.fertilizing"
                    [class.text-emerald-800]="log.fertilizing"
                    [class.bg-stone-100]="!log.fertilizing"
                    [class.text-stone-600]="!log.fertilizing"
                  >
                    {{ log.fertilizing ? 'Fertilized' : 'Water only' }}
                  </span>
                </div>

                <dl class="mt-4 grid grid-cols-2 gap-3 text-sm">
                  <div>
                    <dt class="text-stone-500">Watering</dt>
                    <dd class="mt-1 font-medium text-stone-950">{{ log.wateringMinutes | number: '1.0-0' }} min</dd>
                  </div>
                  <div>
                    <dt class="text-stone-500">Fertilizing</dt>
                    <dd class="mt-1 font-medium text-stone-950">{{ log.fertilizing ? 'Yes' : 'No' }}</dd>
                  </div>
                </dl>

                <p class="mt-4 text-sm leading-6 text-stone-600">{{ log.notes }}</p>

                <div class="mt-4 grid grid-cols-2 gap-2">
                  <button class="btn-secondary" type="button" (click)="openEdit(log)">Edit</button>
                  <button class="btn-danger" type="button" (click)="deleteLog(log)">Delete</button>
                </div>
              </article>
            }
          </div>
        } @else {
          <app-empty-state title="No care logs found" message="Try a different filter or add a new maintenance record.">
            <button class="btn-primary mt-5" type="button" (click)="openCreate()">Add care log</button>
          </app-empty-state>
        }
      </app-section-card>
    </div>

    @if (dialogOpen()) {
      <app-modal [title]="editingLog() ? 'Edit care log' : 'Add care log'" (closed)="closeDialog()">
        <app-maintenance-form
          [log]="editingLog()"
          [plots]="plots()"
          (saved)="saveLog($event)"
          (cancel)="closeDialog()"
        />
      </app-modal>
    }
  `,
})
export class MaintenancePageComponent {
  private readonly fb = inject(NonNullableFormBuilder);
  private readonly farmData = inject(FarmDataService);
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);

  readonly plots = this.farmData.plots;
  readonly dialogOpen = signal(false);
  readonly editingLog = signal<MaintenanceLog | null>(null);

  readonly filters = this.fb.group({
    plot: [''],
    date: [''],
  });

  private readonly filterValue = toSignal(
    this.filters.valueChanges.pipe(startWith(this.filters.getRawValue())),
    { initialValue: this.filters.getRawValue() },
  );

  readonly filteredLogs = computed(() => {
    const filters = this.filterValue();
    return this.farmData.sortedMaintenanceLogs().filter((log) => {
      const plotMatches = !filters.plot || log.plot === filters.plot;
      const dateMatches = !filters.date || log.date === filters.date;

      return plotMatches && dateMatches;
    });
  });

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
    this.editingLog.set(null);
    this.dialogOpen.set(true);
  }

  openEdit(log: MaintenanceLog): void {
    this.editingLog.set(log);
    this.dialogOpen.set(true);
  }

  closeDialog(): void {
    this.dialogOpen.set(false);
    this.editingLog.set(null);
  }

  saveLog(value: MaintenanceLogFormValue): void {
    if (value.id) {
      this.farmData.updateMaintenanceLog(value);
    } else {
      this.farmData.addMaintenanceLog(value);
    }

    this.closeDialog();
  }

  deleteLog(log: MaintenanceLog): void {
    if (window.confirm(`Delete care log for ${log.plot} on ${log.date}?`)) {
      this.farmData.deleteMaintenanceLog(log.id);
    }
  }

  clearFilters(): void {
    this.filters.reset({
      plot: '',
      date: '',
    });
  }
}
