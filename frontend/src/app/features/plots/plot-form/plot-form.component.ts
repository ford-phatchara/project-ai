import { ChangeDetectionStrategy, Component, effect, inject, input, output } from '@angular/core';
import { NonNullableFormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';

import { Plot, PlotFormValue } from '../../../core/models/plot.model';

@Component({
  selector: 'app-plot-form',
  standalone: true,
  imports: [ReactiveFormsModule],
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <form class="space-y-5" [formGroup]="form" (ngSubmit)="submit()">
      <input type="hidden" formControlName="id" />

      <label class="block space-y-2">
        <span class="form-label">Plot name</span>
        <input class="form-field" type="text" formControlName="name" placeholder="North Ridge" />
        @if (form.controls.name.invalid && form.controls.name.touched) {
          <span class="text-sm text-red-600">Name is required.</span>
        }
      </label>

      <label class="block space-y-2">
        <span class="form-label">Size (rai)</span>
        <input class="form-field" type="number" min="0.1" step="0.1" formControlName="size" />
        @if (form.controls.size.invalid && form.controls.size.touched) {
          <span class="text-sm text-red-600">Size must be greater than 0.</span>
        }
      </label>

      <label class="block space-y-2">
        <span class="form-label">Notes</span>
        <textarea class="form-field min-h-28 resize-y" formControlName="notes" placeholder="Soil, irrigation, cultivar, or care notes"></textarea>
      </label>

      <div class="flex flex-col-reverse gap-3 sm:flex-row sm:justify-end">
        <button type="button" class="btn-secondary" (click)="cancel.emit()">Cancel</button>
        <button type="submit" class="btn-primary">Save plot</button>
      </div>
    </form>
  `,
})
export class PlotFormComponent {
  private readonly fb = inject(NonNullableFormBuilder);

  readonly plot = input<Plot | null>(null);
  readonly saved = output<PlotFormValue>();
  readonly cancel = output<void>();

  readonly form = this.fb.group({
    id: [''],
    name: ['', Validators.required],
    size: [0, [Validators.required, Validators.min(0.1)]],
    notes: [''],
  });

  constructor() {
    effect(() => {
      const plot = this.plot();
      this.form.reset({
        id: plot?.id ?? '',
        name: plot?.name ?? '',
        size: plot?.size ?? 0,
        notes: plot?.notes ?? '',
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
      name: value.name,
      size: Number(value.size),
      notes: value.notes,
    });
  }
}
