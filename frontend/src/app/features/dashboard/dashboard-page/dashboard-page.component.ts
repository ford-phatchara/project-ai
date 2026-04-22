import {
  AfterViewInit,
  ChangeDetectionStrategy,
  Component,
  ElementRef,
  OnDestroy,
  ViewChild,
  computed,
  effect,
  inject,
  signal,
} from '@angular/core';
import { CurrencyPipe, DatePipe, DecimalPipe } from '@angular/common';
import { Chart, registerables } from 'chart.js';

import { FarmDataService } from '../../../core/services/farm-data.service';
import { ChartCardComponent } from '../../../shared/components/chart-card/chart-card.component';
import { EmptyStateComponent } from '../../../shared/components/empty-state/empty-state.component';
import { LoadingSkeletonComponent } from '../../../shared/components/loading-skeleton/loading-skeleton.component';
import { SectionCardComponent } from '../../../shared/components/section-card/section-card.component';
import { StatCardComponent } from '../../../shared/components/stat-card/stat-card.component';

Chart.register(...registerables);

@Component({
  selector: 'app-dashboard-page',
  standalone: true,
  imports: [
    CurrencyPipe,
    DatePipe,
    DecimalPipe,
    StatCardComponent,
    ChartCardComponent,
    EmptyStateComponent,
    LoadingSkeletonComponent,
    SectionCardComponent,
  ],
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <div class="space-y-5">
      <section class="grid gap-4 sm:grid-cols-2 xl:grid-cols-3">
        @if (isLoading()) {
          @for (item of [1, 2, 3]; track item) {
            <article class="surface p-5">
              <app-loading-skeleton [lines]="4" />
            </article>
          }
        } @else {
          <app-stat-card
            label="Total yearly revenue"
            [value]="(totalRevenue() | currency: 'THB' : 'code' : '1.0-0') ?? 'THB 0'"
            caption="Revenue from all 2026 sales"
          >
            <svg class="size-5" viewBox="0 0 24 24" fill="none" aria-hidden="true">
              <path d="M4 18.5h16M6 15l4-4 3 3 5-7" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" />
            </svg>
          </app-stat-card>

          <app-stat-card
            label="Total weight sold"
            [value]="(totalWeight() | number: '1.0-0') + ' kg'"
            caption="Across Grade A and Grade B fruit"
          >
            <svg class="size-5" viewBox="0 0 24 24" fill="none" aria-hidden="true">
              <path d="M8 10.5h8M9.5 7h5M7 20h10a2 2 0 0 0 1.96-2.39l-1.7-8.5A2 2 0 0 0 15.3 7.5H8.7a2 2 0 0 0-1.96 1.61l-1.7 8.5A2 2 0 0 0 7 20Z" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" />
            </svg>
          </app-stat-card>

          <app-stat-card
            label="Average price per kg"
            [value]="(averagePrice() | currency: 'THB' : 'code' : '1.0-0') ?? 'THB 0'"
            caption="Weighted by kilograms sold"
          >
            <svg class="size-5" viewBox="0 0 24 24" fill="none" aria-hidden="true">
              <path d="M12 4v16M8 7.5h5.5a2.5 2.5 0 0 1 0 5H10.5a2.5 2.5 0 0 0 0 5H16" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" />
            </svg>
          </app-stat-card>
        }
      </section>

      <section class="grid gap-5 xl:grid-cols-[1.5fr_1fr]">
        <app-chart-card title="Monthly revenue" subtitle="Line chart from mock sales data">
          @if (hasSales()) {
            <div class="h-72">
              <canvas #revenueCanvas aria-label="Monthly revenue line chart"></canvas>
            </div>
          } @else {
            <app-empty-state title="No revenue yet" message="Add a sale to populate the revenue chart." />
          }
        </app-chart-card>

        <app-chart-card title="Grade distribution" subtitle="Grade A vs Grade B by weight">
          @if (hasSales()) {
            <div class="mx-auto h-72 max-w-sm">
              <canvas #gradeCanvas aria-label="Grade distribution pie chart"></canvas>
            </div>
          } @else {
            <app-empty-state title="No grades yet" message="Add sales with Grade A or Grade B fruit." />
          }
        </app-chart-card>
      </section>

      <app-section-card title="Recent sales" subtitle="Latest farm gate and wholesale records">
        @if (recentSales().length) {
          <div class="divide-y divide-stone-100">
            @for (sale of recentSales(); track sale.id) {
              <div class="flex flex-col gap-3 py-3 first:pt-0 last:pb-0 sm:flex-row sm:items-center sm:justify-between">
                <div>
                  <p class="font-medium text-stone-950">{{ sale.plot }} - Grade {{ sale.grade }}</p>
                  <p class="mt-1 text-sm text-stone-500">
                    {{ sale.date | date: 'mediumDate' }} - {{ sale.weightKg | number: '1.0-1' }} kg
                  </p>
                </div>
                <p class="text-base font-semibold text-stone-950">
                  {{ sale.totalPrice | currency: 'THB' : 'symbol' : '1.0-0' }}
                </p>
              </div>
            }
          </div>
        } @else {
          <app-empty-state title="No sales recorded" message="Sales you create will appear here." />
        }
      </app-section-card>
    </div>
  `,
})
export class DashboardPageComponent implements AfterViewInit, OnDestroy {
  @ViewChild('revenueCanvas') private revenueCanvas?: ElementRef<HTMLCanvasElement>;
  @ViewChild('gradeCanvas') private gradeCanvas?: ElementRef<HTMLCanvasElement>;

  private readonly farmData = inject(FarmDataService);
  private revenueChart?: Chart;
  private gradeChart?: Chart;

  readonly isLoading = signal(true);
  readonly sales = this.farmData.sortedSales;
  readonly hasSales = computed(() => this.sales().length > 0);
  readonly recentSales = computed(() => this.sales().slice(0, 5));

  readonly totalRevenue = computed(() =>
    this.sales().reduce((total, sale) => total + sale.totalPrice, 0),
  );

  readonly totalWeight = computed(() =>
    this.sales().reduce((total, sale) => total + sale.weightKg, 0),
  );

  readonly averagePrice = computed(() => {
    const totalWeight = this.totalWeight();
    return totalWeight ? this.totalRevenue() / totalWeight : 0;
  });

  readonly monthlyRevenue = computed(() => {
    const months = Array.from({ length: 12 }, (_, index) => ({
      label: new Date(2026, index, 1).toLocaleString('en-US', { month: 'short' }),
      value: 0,
    }));

    for (const sale of this.sales()) {
      const date = new Date(`${sale.date}T00:00:00`);
      if (date.getFullYear() === 2026) {
        months[date.getMonth()].value += sale.totalPrice;
      }
    }

    return months;
  });

  readonly gradeDistribution = computed(() => {
    const totals = { A: 0, B: 0 };

    for (const sale of this.sales()) {
      totals[sale.grade] += sale.weightKg;
    }

    return totals;
  });

  constructor() {
    window.setTimeout(() => this.isLoading.set(false), 450);

    effect(() => {
      this.monthlyRevenue();
      this.gradeDistribution();
      queueMicrotask(() => this.renderCharts());
    });
  }

  ngAfterViewInit(): void {
    this.renderCharts();
  }

  ngOnDestroy(): void {
    this.revenueChart?.destroy();
    this.gradeChart?.destroy();
  }

  private renderCharts(): void {
    if (!this.revenueCanvas || !this.gradeCanvas) {
      return;
    }

    const monthlyRevenue = this.monthlyRevenue();
    const gradeDistribution = this.gradeDistribution();

    if (this.revenueChart) {
      this.revenueChart.data.labels = monthlyRevenue.map((item) => item.label);
      this.revenueChart.data.datasets[0].data = monthlyRevenue.map((item) => item.value);
      this.revenueChart.update();
    } else {
      this.revenueChart = new Chart(this.revenueCanvas.nativeElement, {
        type: 'line',
        data: {
          labels: monthlyRevenue.map((item) => item.label),
          datasets: [
            {
              label: 'Revenue',
              data: monthlyRevenue.map((item) => item.value),
              borderColor: '#047857',
              backgroundColor: 'rgba(4, 120, 87, 0.12)',
              fill: true,
              tension: 0.35,
              pointBackgroundColor: '#047857',
              pointRadius: 4,
            },
          ],
        },
        options: {
          responsive: true,
          maintainAspectRatio: false,
          plugins: {
            legend: { display: false },
          },
          scales: {
            x: { grid: { display: false } },
            y: {
              ticks: {
                callback: (value) => `THB ${Number(value).toLocaleString('th-TH')}`,
              },
            },
          },
        },
      });
    }

    if (this.gradeChart) {
      this.gradeChart.data.datasets[0].data = [
        gradeDistribution.A,
        gradeDistribution.B,
      ];
      this.gradeChart.update();
    } else {
      this.gradeChart = new Chart(this.gradeCanvas.nativeElement, {
        type: 'pie',
        data: {
          labels: ['Grade A', 'Grade B'],
          datasets: [
            {
              data: [gradeDistribution.A, gradeDistribution.B],
              backgroundColor: ['#047857', '#f59e0b'],
              borderColor: '#ffffff',
              borderWidth: 4,
            },
          ],
        },
        options: {
          responsive: true,
          maintainAspectRatio: false,
          plugins: {
            legend: {
              position: 'bottom',
              labels: {
                boxWidth: 12,
                usePointStyle: true,
              },
            },
          },
        },
      });
    }
  }
}
