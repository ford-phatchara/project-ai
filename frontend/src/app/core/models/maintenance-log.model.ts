export interface MaintenanceLog {
  id: string;
  date: string;
  plot: string;
  wateringMinutes: number;
  fertilizing: boolean;
  notes: string;
}

export type MaintenanceLogFormValue = Omit<MaintenanceLog, 'id'> &
  Partial<Pick<MaintenanceLog, 'id'>>;
