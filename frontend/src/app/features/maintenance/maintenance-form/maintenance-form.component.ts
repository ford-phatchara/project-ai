import { ChangeDetectionStrategy, Component, effect, inject, input, output } from '@angular/core';
import { NonNullableFormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';

import {
  MaintenanceLog,
  MaintenanceLogFormValue,
} from '../../../core/models/maintenance-log.model';
import { Plot } from '../../../core/models/plot.model';

@Component({
  selector: 'app-maintenance-form',
  standalone: true,
  imports: [ReactiveFormsModule],
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <form class="space-y-5" [formGroup]="form" (ngSubmit)="submit()">
      <input type="hidden" formControlName="id" />

      <div class="grid gap-4 sm:grid-cols-2">
        <label class="block space-y-2">
          <span class="form-label">Date</span>
          <input class="form-field" type="date" formControlName="date" />
          @if (form.controls.date.invalid && form.controls.date.touched) {
            <span class="text-sm text-red-600">Date is required.</span>
          }
        </label>

        <label class="block space-y-2">
          <span class="form-label">Plot</span>
          <select class="form-field" formControlName="plot">
            <option value="" disabled>Select plot</option>
            @for (plot of plots(); track plot.id) {
              <option [value]="plot.name">{{ plot.name }}</option>
            }
          </select>
          @if (form.controls.plot.invalid && form.controls.plot.touched) {
            <span class="text-sm text-red-600">Plot is required.</span>
          }
        </label>
      </div>

      <label class="block space-y-2">
        <span class="form-label">Watering duration (minutes)</span>
        <input class="form-field" type="number" min="0" step="5" formControlName="wateringMinutes" />
        @if (form.controls.wateringMinutes.invalid && form.controls.wateringMinutes.touched) {
          <span class="text-sm text-red-600">Watering duration is required.</span>
        }
      </label>

      <label class="flex items-center justify-between gap-4 rounded-xl border border-stone-200 bg-white px-4 py-3 shadow-sm">
        <span>
          <span class="block text-sm font-medium text-stone-800">Fertilizing</span>
          <span class="block text-sm text-stone-500">Mark when fertilizer was applied.</span>
        </span>
        <input
          type="checkbox"
          class="size-5 rounded border-stone-300 text-emerald-700 focus:ring-emerald-200"
          formControlName="fertilizing"
        />
      </label>

      <label class="block space-y-2">
        <span class="form-label">Notes</span>
        <textarea class="form-field min-h-28 resize-y" formControlName="notes" placeholder="Care details, weather, pests, soil moisture, or follow-up work"></textarea>
      </label>

      <div class="flex flex-col-reverse gap-3 sm:flex-row sm:justify-end">
        <button type="button" class="btn-secondary" (click)="cancel.emit()">Cancel</button>
        <button type="submit" class="btn-primary">Save care log</button>
      </div>
    </form>
  `,
})
export class MaintenanceFormComponent {
  private readonly fb = inject(NonNullableFormBuilder);

  readonly log = input<MaintenanceLog | null>(null);
  readonly plots = input<Plot[]>([]);
  readonly saved = output<MaintenanceLogFormValue>();
  readonly cancel = output<void>();

  readonly form = this.fb.group({
    id: [''],
    date: [this.today(), Validators.required],
    plot: ['', Validators.required],
    wateringMinutes: [0, [Validators.required, Validators.min(0)]],
    fertilizing: [false],
    notes: [''],
  });

  constructor() {
    effect(() => {
      const log = this.log();
      const firstPlot = this.plots()[0]?.name ?? '';

      this.form.reset({
        id: log?.id ?? '',
        date: log?.date ?? this.today(),
        plot: log?.plot ?? firstPlot,
        wateringMinutes: log?.wateringMinutes ?? 0,
        fertilizing: log?.fertilizing ?? false,
        notes: log?.notes ?? '',
      });
    });
  }

  submit(): void {
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }

    const value = this.form.getRawValue();
    this.saved.emit({
      id: value.id || undefined,
      date: value.date,
      plot: value.plot,
      wateringMinutes: Number(value.wateringMinutes),
      fertilizing: value.fertilizing,
      notes: value.notes,
    });
  }

  private today(): string {
    return new Date().toISOString().slice(0, 10);
  }
}
