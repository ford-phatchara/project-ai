import { ChangeDetectionStrategy, Component, effect, inject, input, output } from '@angular/core';
import { NonNullableFormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';

import { Plot } from '../../../core/models/plot.model';
import { DurianGrade, Sale, SaleFormValue } from '../../../core/models/sale.model';

@Component({
  selector: 'app-sale-form',
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

      <div class="grid gap-4 sm:grid-cols-3">
        <label class="block space-y-2">
          <span class="form-label">Grade</span>
          <select class="form-field" formControlName="grade">
            <option value="A">Grade A</option>
            <option value="B">Grade B</option>
          </select>
        </label>

        <label class="block space-y-2">
          <span class="form-label">Weight (kg)</span>
          <input class="form-field" type="number" min="0.1" step="0.1" formControlName="weightKg" />
          @if (form.controls.weightKg.invalid && form.controls.weightKg.touched) {
            <span class="text-sm text-red-600">Weight must be greater than 0.</span>
          }
        </label>

        <label class="block space-y-2">
          <span class="form-label">Price per kg</span>
          <input class="form-field" type="number" min="0.1" step="0.1" formControlName="pricePerKg" />
          @if (form.controls.pricePerKg.invalid && form.controls.pricePerKg.touched) {
            <span class="text-sm text-red-600">Price must be greater than 0.</span>
          }
        </label>
      </div>

      <div class="rounded-xl border border-emerald-100 bg-emerald-50 px-4 py-3">
        <p class="text-sm font-medium text-emerald-900">Auto-calculated total</p>
        <p class="mt-1 text-2xl font-semibold text-emerald-950">
          {{ totalPrice().toLocaleString('th-TH', { style: 'currency', currency: 'THB' }) }}
        </p>
      </div>

      <div class="flex flex-col-reverse gap-3 sm:flex-row sm:justify-end">
        <button type="button" class="btn-secondary" (click)="cancel.emit()">Cancel</button>
        <button type="submit" class="btn-primary">Save sale</button>
      </div>
    </form>
  `,
})
export class SaleFormComponent {
  private readonly fb = inject(NonNullableFormBuilder);

  readonly sale = input<Sale | null>(null);
  readonly plots = input<Plot[]>([]);
  readonly saved = output<SaleFormValue>();
  readonly cancel = output<void>();

  readonly form = this.fb.group({
    id: [''],
    date: [this.today(), Validators.required],
    plot: ['', Validators.required],
    grade: ['A' as DurianGrade, Validators.required],
    weightKg: [0, [Validators.required, Validators.min(0.1)]],
    pricePerKg: [0, [Validators.required, Validators.min(0.1)]],
  });

  constructor() {
    effect(() => {
      const sale = this.sale();
      const firstPlot = this.plots()[0]?.name ?? '';

      this.form.reset({
        id: sale?.id ?? '',
        date: sale?.date ?? this.today(),
        plot: sale?.plot ?? firstPlot,
        grade: sale?.grade ?? 'A',
        weightKg: sale?.weightKg ?? 0,
        pricePerKg: sale?.pricePerKg ?? 0,
      });
    });
  }

  totalPrice(): number {
    const { weightKg, pricePerKg } = this.form.getRawValue();
    return Number(weightKg) * Number(pricePerKg);
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
      grade: value.grade,
      weightKg: Number(value.weightKg),
      pricePerKg: Number(value.pricePerKg),
      totalPrice: Number(value.weightKg) * Number(value.pricePerKg),
    });
  }

  private today(): string {
    return new Date().toISOString().slice(0, 10);
  }
}
