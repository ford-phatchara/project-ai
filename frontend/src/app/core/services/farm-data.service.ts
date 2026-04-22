import { computed, Injectable, signal } from '@angular/core';

import {
  MaintenanceLog,
  MaintenanceLogFormValue,
} from '../models/maintenance-log.model';
import { Plot, PlotFormValue } from '../models/plot.model';
import { Sale, SaleFormValue } from '../models/sale.model';
import {
  MOCK_MAINTENANCE_LOGS,
  MOCK_PLOTS,
  MOCK_SALES,
} from '../mock-data/farm.mock';

@Injectable({
  providedIn: 'root',
})
export class FarmDataService {
  private readonly plotsSignal = signal<Plot[]>([...MOCK_PLOTS]);
  private readonly salesSignal = signal<Sale[]>([...MOCK_SALES]);
  private readonly maintenanceLogsSignal = signal<MaintenanceLog[]>([
    ...MOCK_MAINTENANCE_LOGS,
  ]);

  readonly plots = this.plotsSignal.asReadonly();
  readonly sales = this.salesSignal.asReadonly();
  readonly maintenanceLogs = this.maintenanceLogsSignal.asReadonly();

  readonly sortedSales = computed(() =>
    [...this.sales()].sort((a, b) => b.date.localeCompare(a.date)),
  );

  readonly sortedMaintenanceLogs = computed(() =>
    [...this.maintenanceLogs()].sort((a, b) => b.date.localeCompare(a.date)),
  );

  addSale(value: SaleFormValue): void {
    const sale = this.normalizeSale({
      ...value,
      id: this.createId('sale'),
    });
    this.salesSignal.update((sales) => [sale, ...sales]);
  }

  updateSale(value: SaleFormValue): void {
    if (!value.id) {
      return;
    }

    const sale = this.normalizeSale(value as Sale);
    this.salesSignal.update((sales) =>
      sales.map((item) => (item.id === sale.id ? sale : item)),
    );
  }

  deleteSale(id: string): void {
    this.salesSignal.update((sales) => sales.filter((sale) => sale.id !== id));
  }

  addPlot(value: PlotFormValue): void {
    const plot: Plot = {
      id: this.createId('plot'),
      name: value.name.trim(),
      size: Number(value.size),
      notes: value.notes.trim(),
    };
    this.plotsSignal.update((plots) => [plot, ...plots]);
  }

  updatePlot(value: PlotFormValue): void {
    if (!value.id) {
      return;
    }

    const plot: Plot = {
      id: value.id,
      name: value.name.trim(),
      size: Number(value.size),
      notes: value.notes.trim(),
    };
    this.plotsSignal.update((plots) =>
      plots.map((item) => (item.id === plot.id ? plot : item)),
    );
  }

  deletePlot(id: string): void {
    const plot = this.plots().find((item) => item.id === id);
    this.plotsSignal.update((plots) => plots.filter((item) => item.id !== id));

    if (plot) {
      this.salesSignal.update((sales) =>
        sales.filter((sale) => sale.plot !== plot.name),
      );
      this.maintenanceLogsSignal.update((logs) =>
        logs.filter((log) => log.plot !== plot.name),
      );
    }
  }

  addMaintenanceLog(value: MaintenanceLogFormValue): void {
    const log: MaintenanceLog = {
      id: this.createId('care'),
      date: value.date,
      plot: value.plot,
      wateringMinutes: Number(value.wateringMinutes),
      fertilizing: Boolean(value.fertilizing),
      notes: value.notes.trim(),
    };
    this.maintenanceLogsSignal.update((logs) => [log, ...logs]);
  }

  updateMaintenanceLog(value: MaintenanceLogFormValue): void {
    if (!value.id) {
      return;
    }

    const log: MaintenanceLog = {
      id: value.id,
      date: value.date,
      plot: value.plot,
      wateringMinutes: Number(value.wateringMinutes),
      fertilizing: Boolean(value.fertilizing),
      notes: value.notes.trim(),
    };
    this.maintenanceLogsSignal.update((logs) =>
      logs.map((item) => (item.id === log.id ? log : item)),
    );
  }

  deleteMaintenanceLog(id: string): void {
    this.maintenanceLogsSignal.update((logs) =>
      logs.filter((log) => log.id !== id),
    );
  }

  private normalizeSale(value: SaleFormValue & { id: string }): Sale {
    const weightKg = Number(value.weightKg);
    const pricePerKg = Number(value.pricePerKg);

    return {
      id: value.id,
      date: value.date,
      plot: value.plot,
      grade: value.grade,
      weightKg,
      pricePerKg,
      totalPrice: weightKg * pricePerKg,
    };
  }

  private createId(prefix: string): string {
    const id =
      globalThis.crypto?.randomUUID?.() ??
      Math.random().toString(36).slice(2, 11);

    return `${prefix}-${id}`;
  }
}
