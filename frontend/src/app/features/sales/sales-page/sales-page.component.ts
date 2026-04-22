import { CurrencyPipe, DatePipe, DecimalPipe } from '@angular/common';
import { ChangeDetectionStrategy, Component, computed, inject, signal } from '@angular/core';
import { takeUntilDestroyed, toSignal } from '@angular/core/rxjs-interop';
import { NonNullableFormBuilder, ReactiveFormsModule } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { startWith } from 'rxjs';

import { Sale, SaleFormValue } from '../../../core/models/sale.model';
import { FarmDataService } from '../../../core/services/farm-data.service';
import { EmptyStateComponent } from '../../../shared/components/empty-state/empty-state.component';
import { ModalComponent } from '../../../shared/components/modal/modal.component';
import { SectionCardComponent } from '../../../shared/components/section-card/section-card.component';
import { SaleFormComponent } from '../sale-form/sale-form.component';

@Component({
  selector: 'app-sales-page',
  standalone: true,
  imports: [
    ReactiveFormsModule,
    CurrencyPipe,
    DatePipe,
    DecimalPipe,
    EmptyStateComponent,
    ModalComponent,
    SectionCardComponent,
    SaleFormComponent,
  ],
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <div class="space-y-5">
      <section class="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <p class="text-sm font-medium text-emerald-700">Sales tracking</p>
          <h2 class="mt-1 text-2xl font-semibold tracking-tight text-stone-950">
            Durian sales records
          </h2>
          <p class="mt-2 max-w-2xl text-sm leading-6 text-stone-500">
            Record grade, weight, and price per kilogram with totals calculated automatically.
          </p>
        </div>
        <button class="btn-primary" type="button" (click)="openCreate()">
          <svg class="size-4" viewBox="0 0 24 24" fill="none" aria-hidden="true">
            <path d="M12 5v14M5 12h14" stroke="currentColor" stroke-width="2" stroke-linecap="round" />
          </svg>
          Add sale
        </button>
      </section>

      <section class="surface p-4">
        <form class="grid gap-3 sm:grid-cols-2 xl:grid-cols-5" [formGroup]="filters">
          <label class="block space-y-2">
            <span class="form-label">Start date</span>
            <input class="form-field" type="date" formControlName="startDate" />
          </label>
          <label class="block space-y-2">
            <span class="form-label">End date</span>
            <input class="form-field" type="date" formControlName="endDate" />
          </label>
          <label class="block space-y-2">
            <span class="form-label">Grade</span>
            <select class="form-field" formControlName="grade">
              <option value="">All grades</option>
              <option value="A">Grade A</option>
              <option value="B">Grade B</option>
            </select>
          </label>
          <label class="block space-y-2">
            <span class="form-label">Plot</span>
            <select class="form-field" formControlName="plot">
              <option value="">All plots</option>
              @for (plot of plots(); track plot.id) {
                <option [value]="plot.name">{{ plot.name }}</option>
              }
            </select>
          </label>
          <div class="flex items-end">
            <button class="btn-secondary w-full" type="button" (click)="clearFilters()">Clear filters</button>
          </div>
        </form>
      </section>

      <app-section-card
        title="Sales list"
        [subtitle]="filteredSales().length + ' records shown'"
      >
        @if (filteredSales().length) {
          <div class="hidden overflow-x-auto md:block">
            <table class="min-w-full border-separate border-spacing-0 text-sm">
              <thead>
                <tr class="table-header">
                  <th class="rounded-l-xl px-4 py-3">Date</th>
                  <th class="px-4 py-3">Plot</th>
                  <th class="px-4 py-3">Grade</th>
                  <th class="px-4 py-3 text-right">Weight</th>
                  <th class="px-4 py-3 text-right">Price/kg</th>
                  <th class="px-4 py-3 text-right">Total</th>
                  <th class="rounded-r-xl px-4 py-3 text-right">Actions</th>
                </tr>
              </thead>
              <tbody>
                @for (sale of filteredSales(); track sale.id) {
                  <tr class="border-b border-stone-100">
                    <td class="px-4 py-4 text-stone-600">{{ sale.date | date: 'mediumDate' }}</td>
                    <td class="px-4 py-4 font-medium text-stone-950">{{ sale.plot }}</td>
                    <td class="px-4 py-4">
                      <span class="rounded-full bg-emerald-50 px-2.5 py-1 text-xs font-semibold text-emerald-800">
                        Grade {{ sale.grade }}
                      </span>
                    </td>
                    <td class="px-4 py-4 text-right text-stone-700">{{ sale.weightKg | number: '1.0-1' }} kg</td>
                    <td class="px-4 py-4 text-right text-stone-700">{{ sale.pricePerKg | currency: 'THB' : 'symbol' : '1.0-0' }}</td>
                    <td class="px-4 py-4 text-right font-semibold text-stone-950">{{ sale.totalPrice | currency: 'THB' : 'symbol' : '1.0-0' }}</td>
                    <td class="px-4 py-4">
                      <div class="flex justify-end gap-2">
                        <button class="btn-secondary px-3 py-2" type="button" (click)="openEdit(sale)">Edit</button>
                        <button class="btn-danger" type="button" (click)="deleteSale(sale)">Delete</button>
                      </div>
                    </td>
                  </tr>
                }
              </tbody>
            </table>
          </div>

          <div class="space-y-3 md:hidden">
            @for (sale of filteredSales(); track sale.id) {
              <article class="rounded-xl border border-stone-200 bg-white p-4 shadow-sm">
                <div class="flex items-start justify-between gap-3">
                  <div>
                    <p class="font-semibold text-stone-950">{{ sale.plot }}</p>
                    <p class="mt-1 text-sm text-stone-500">{{ sale.date | date: 'mediumDate' }}</p>
                  </div>
                  <span class="rounded-full bg-emerald-50 px-2.5 py-1 text-xs font-semibold text-emerald-800">
                    Grade {{ sale.grade }}
                  </span>
                </div>
                <dl class="mt-4 grid grid-cols-2 gap-3 text-sm">
                  <div>
                    <dt class="text-stone-500">Weight</dt>
                    <dd class="mt-1 font-medium text-stone-950">{{ sale.weightKg | number: '1.0-1' }} kg</dd>
                  </div>
                  <div>
                    <dt class="text-stone-500">Price/kg</dt>
                    <dd class="mt-1 font-medium text-stone-950">{{ sale.pricePerKg | currency: 'THB' : 'symbol' : '1.0-0' }}</dd>
                  </div>
                  <div class="col-span-2">
                    <dt class="text-stone-500">Total</dt>
                    <dd class="mt-1 text-lg font-semibold text-stone-950">{{ sale.totalPrice | currency: 'THB' : 'symbol' : '1.0-0' }}</dd>
                  </div>
                </dl>
                <div class="mt-4 grid grid-cols-2 gap-2">
                  <button class="btn-secondary" type="button" (click)="openEdit(sale)">Edit</button>
                  <button class="btn-danger" type="button" (click)="deleteSale(sale)">Delete</button>
                </div>
              </article>
            }
          </div>
        } @else {
          <app-empty-state title="No sales found" message="Try a different filter or add a new sale record.">
            <button class="btn-primary mt-5" type="button" (click)="openCreate()">Add sale</button>
          </app-empty-state>
        }
      </app-section-card>
    </div>

    @if (dialogOpen()) {
      <app-modal [title]="editingSale() ? 'Edit sale' : 'Add sale'" (closed)="closeDialog()">
        <app-sale-form
          [sale]="editingSale()"
          [plots]="plots()"
          (saved)="saveSale($event)"
          (cancel)="closeDialog()"
        />
      </app-modal>
    }
  `,
})
export class SalesPageComponent {
  private readonly fb = inject(NonNullableFormBuilder);
  private readonly farmData = inject(FarmDataService);
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);

  readonly plots = this.farmData.plots;
  readonly dialogOpen = signal(false);
  readonly editingSale = signal<Sale | null>(null);

  readonly filters = this.fb.group({
    startDate: [''],
    endDate: [''],
    grade: [''],
    plot: [''],
  });

  private readonly filterValue = toSignal(
    this.filters.valueChanges.pipe(startWith(this.filters.getRawValue())),
    { initialValue: this.filters.getRawValue() },
  );

  readonly filteredSales = computed(() => {
    const filters = this.filterValue();
    return this.farmData.sortedSales().filter((sale) => {
      const afterStart = !filters.startDate || sale.date >= filters.startDate;
      const beforeEnd = !filters.endDate || sale.date <= filters.endDate;
      const gradeMatches = !filters.grade || sale.grade === filters.grade;
      const plotMatches = !filters.plot || sale.plot === filters.plot;

      return afterStart && beforeEnd && gradeMatches && plotMatches;
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
    this.editingSale.set(null);
    this.dialogOpen.set(true);
  }

  openEdit(sale: Sale): void {
    this.editingSale.set(sale);
    this.dialogOpen.set(true);
  }

  closeDialog(): void {
    this.dialogOpen.set(false);
    this.editingSale.set(null);
  }

  saveSale(value: SaleFormValue): void {
    if (value.id) {
      this.farmData.updateSale(value);
    } else {
      this.farmData.addSale(value);
    }

    this.closeDialog();
  }

  deleteSale(sale: Sale): void {
    if (window.confirm(`Delete sale from ${sale.plot} on ${sale.date}?`)) {
      this.farmData.deleteSale(sale.id);
    }
  }

  clearFilters(): void {
    this.filters.reset({
      startDate: '',
      endDate: '',
      grade: '',
      plot: '',
    });
  }
}
