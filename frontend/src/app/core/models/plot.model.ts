export interface Plot {
  id: string;
  name: string;
  size: number;
  notes: string;
}

export type PlotFormValue = Omit<Plot, 'id'> & Partial<Pick<Plot, 'id'>>;
