import { DatePipe, DecimalPipe } from '@angular/common';
import { ChangeDetectionStrategy, Component, computed, inject, signal } from '@angular/core';

import { MaintenanceLog } from '../../../core/models/maintenance-log.model';
import { FarmDataService } from '../../../core/services/farm-data.service';
import { EmptyStateComponent } from '../../../shared/components/empty-state/empty-state.component';
import { SectionCardComponent } from '../../../shared/components/section-card/section-card.component';

interface CalendarDay {
  date: string;
  dayNumber: number;
  inMonth: boolean;
  logs: MaintenanceLog[];
}

type CalendarMode = 'month' | 'day';

@Component({
  selector: 'app-calendar-page',
  standalone: true,
  imports: [DatePipe, DecimalPipe, EmptyStateComponent, SectionCardComponent],
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <div class="space-y-5">
      <section class="flex flex-col gap-4 md:flex-row md:items-end md:justify-between">
        <div>
          <p class="text-sm font-medium text-emerald-700">Maintenance calendar</p>
          <h2 class="mt-1 text-2xl font-semibold tracking-tight text-stone-950">
            Care schedule
          </h2>
          <p class="mt-2 max-w-2xl text-sm leading-6 text-stone-500">
            See watering and fertilizing activity by month or by selected day.
          </p>
        </div>

        <div class="grid grid-cols-2 rounded-xl border border-stone-200 bg-white p-1 shadow-sm">
          <button
            type="button"
            class="rounded-lg px-4 py-2 text-sm font-semibold transition"
            [class.bg-emerald-700]="mode() === 'month'"
            [class.text-white]="mode() === 'month'"
            [class.text-stone-600]="mode() !== 'month'"
            (click)="mode.set('month')"
          >
            Month
          </button>
          <button
            type="button"
            class="rounded-lg px-4 py-2 text-sm font-semibold transition"
            [class.bg-emerald-700]="mode() === 'day'"
            [class.text-white]="mode() === 'day'"
            [class.text-stone-600]="mode() !== 'day'"
            (click)="mode.set('day')"
          >
            Day
          </button>
        </div>
      </section>

      <app-section-card [title]="monthTitle()" subtitle="Maintenance events from care logs">
        <div section-action class="flex gap-2">
          <button class="btn-secondary px-3 py-2" type="button" aria-label="Previous month" (click)="previousMonth()">
            <svg class="size-4" viewBox="0 0 24 24" fill="none" aria-hidden="true">
              <path d="m15 6-6 6 6 6" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" />
            </svg>
          </button>
          <button class="btn-secondary px-3 py-2" type="button" aria-label="Next month" (click)="nextMonth()">
            <svg class="size-4" viewBox="0 0 24 24" fill="none" aria-hidden="true">
              <path d="m9 6 6 6-6 6" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" />
            </svg>
          </button>
        </div>

        @if (mode() === 'month') {
          <div class="grid grid-cols-7 gap-1 text-center text-xs font-semibold uppercase tracking-wide text-stone-500">
            @for (day of weekDays; track day) {
              <div class="py-2">{{ day }}</div>
            }
          </div>

          <div class="mt-1 grid grid-cols-7 gap-1">
            @for (day of monthDays(); track day.date) {
              <button
                type="button"
                class="min-h-24 rounded-xl border p-2 text-left transition hover:border-emerald-200 hover:bg-emerald-50/60"
                [class.border-stone-200]="day.inMonth && selectedDate() !== day.date"
                [class.border-transparent]="!day.inMonth"
                [class.bg-white]="day.inMonth && selectedDate() !== day.date"
                [class.bg-stone-50]="!day.inMonth"
                [class.opacity-50]="!day.inMonth"
                [class.border-emerald-500]="selectedDate() === day.date"
                [class.bg-emerald-50]="selectedDate() === day.date"
                (click)="selectDay(day)"
              >
                <span class="text-sm font-semibold text-stone-950">{{ day.dayNumber }}</span>
                @if (day.logs.length) {
                  <div class="mt-2 space-y-1">
                    @for (log of day.logs.slice(0, 2); track log.id) {
                      <div class="truncate rounded-lg bg-emerald-100 px-2 py-1 text-xs font-medium text-emerald-900">
                        {{ log.plot }}
                      </div>
                    }
                    @if (day.logs.length > 2) {
                      <p class="px-1 text-xs font-medium text-emerald-800">
                        +{{ day.logs.length - 2 }} more
                      </p>
                    }
                  </div>
                }
              </button>
            }
          </div>
        } @else {
          <div class="rounded-xl border border-stone-200 bg-stone-50 p-4">
            <p class="text-sm font-medium text-stone-500">Selected day</p>
            <p class="mt-1 text-xl font-semibold text-stone-950">
              {{ selectedDate() | date: 'fullDate' }}
            </p>
          </div>

          <div class="mt-4">
            @if (selectedLogs().length) {
              <div class="space-y-3">
                @for (log of selectedLogs(); track log.id) {
                  <article class="rounded-xl border border-stone-200 bg-white p-4 shadow-sm">
                    <div class="flex items-start justify-between gap-3">
                      <div>
                        <p class="font-semibold text-stone-950">{{ log.plot }}</p>
                        <p class="mt-1 text-sm text-stone-500">{{ log.wateringMinutes | number: '1.0-0' }} minutes watering</p>
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
                    <p class="mt-3 text-sm leading-6 text-stone-600">{{ log.notes }}</p>
                  </article>
                }
              </div>
            } @else {
              <app-empty-state title="No activity on this day" message="Choose a highlighted day to inspect its care logs." />
            }
          </div>
        }
      </app-section-card>
    </div>
  `,
})
export class CalendarPageComponent {
  private readonly farmData = inject(FarmDataService);

  readonly weekDays = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];
  readonly mode = signal<CalendarMode>('month');
  readonly cursor = signal(new Date(2026, 3, 1));
  readonly selectedDate = signal('2026-04-22');

  readonly eventsByDate = computed(() => {
    const events = new Map<string, MaintenanceLog[]>();

    for (const log of this.farmData.maintenanceLogs()) {
      const logs = events.get(log.date) ?? [];
      events.set(log.date, [...logs, log]);
    }

    return events;
  });

  readonly monthTitle = computed(() =>
    this.cursor().toLocaleDateString('en-US', {
      month: 'long',
      year: 'numeric',
    }),
  );

  readonly monthDays = computed(() => {
    const current = this.cursor();
    const firstOfMonth = new Date(current.getFullYear(), current.getMonth(), 1);
    const start = new Date(firstOfMonth);
    start.setDate(firstOfMonth.getDate() - firstOfMonth.getDay());

    return Array.from({ length: 42 }, (_, index): CalendarDay => {
      const date = new Date(start);
      date.setDate(start.getDate() + index);
      const key = this.formatDate(date);

      return {
        date: key,
        dayNumber: date.getDate(),
        inMonth: date.getMonth() === current.getMonth(),
        logs: this.eventsByDate().get(key) ?? [],
      };
    });
  });

  readonly selectedLogs = computed(() =>
    this.eventsByDate().get(this.selectedDate()) ?? [],
  );

  selectDay(day: CalendarDay): void {
    this.selectedDate.set(day.date);
    this.mode.set('day');
  }

  previousMonth(): void {
    this.cursor.update(
      (date) => new Date(date.getFullYear(), date.getMonth() - 1, 1),
    );
  }

  nextMonth(): void {
    this.cursor.update(
      (date) => new Date(date.getFullYear(), date.getMonth() + 1, 1),
    );
  }

  private formatDate(date: Date): string {
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');

    return `${year}-${month}-${day}`;
  }
}
