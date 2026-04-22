export type DurianGrade = 'A' | 'B';

export interface Sale {
  id: string;
  date: string;
  plot: string;
  grade: DurianGrade;
  weightKg: number;
  pricePerKg: number;
  totalPrice: number;
}

export type SaleFormValue = Omit<Sale, 'id'> & Partial<Pick<Sale, 'id'>>;
