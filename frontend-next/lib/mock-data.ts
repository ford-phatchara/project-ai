export type Grade = "A" | "B"

export interface Sale {
  id: string
  date: string
  plot: string
  grade: Grade
  weightKg: number
  pricePerKg: number
  total: number
  buyer?: string
}

export interface Plot {
  id: string
  name: string
  sizeMeter: number
  trees: number
  notes: string
}

export interface MaintenanceLog {
  id: string
  date: string
  plot: string
  activity: "Watering" | "Fertilizing" | "Pruning" | "Pest Control" | "Harvesting"
  durationMin?: number
  notes: string
}

export const PLOTS: Plot[] = [
  { id: "P1", name: "Plot A — Hillside", sizeMeter: 2400, trees: 48, notes: "Older trees, best for Grade A. Needs extra irrigation in dry season." },
  { id: "P2", name: "Plot B — Valley", sizeMeter: 1800, trees: 36, notes: "Younger trees planted 2021. Good drainage." },
  { id: "P3", name: "Plot C — Eastern Ridge", sizeMeter: 3200, trees: 64, notes: "Mixed grades. Watch for root rot near stream." },
  { id: "P4", name: "Plot D — North Field", sizeMeter: 1200, trees: 24, notes: "Newly planted 2023. No harvest yet." },
]

export const SALES: Sale[] = [
  { id: "S001", date: "2024-07-15", plot: "P1", grade: "A", weightKg: 320, pricePerKg: 28, total: 8960, buyer: "Ahmad Traders" },
  { id: "S002", date: "2024-07-18", plot: "P3", grade: "B", weightKg: 210, pricePerKg: 14, total: 2940, buyer: "Local Market" },
  { id: "S003", date: "2024-07-22", plot: "P1", grade: "A", weightKg: 415, pricePerKg: 30, total: 12450, buyer: "KL Export Co." },
  { id: "S004", date: "2024-07-28", plot: "P2", grade: "A", weightKg: 180, pricePerKg: 27, total: 4860, buyer: "Ahmad Traders" },
  { id: "S005", date: "2024-08-03", plot: "P3", grade: "B", weightKg: 290, pricePerKg: 13, total: 3770, buyer: "Local Market" },
  { id: "S006", date: "2024-08-10", plot: "P1", grade: "A", weightKg: 500, pricePerKg: 32, total: 16000, buyer: "KL Export Co." },
  { id: "S007", date: "2024-08-17", plot: "P2", grade: "B", weightKg: 140, pricePerKg: 14, total: 1960, buyer: "Local Market" },
  { id: "S008", date: "2024-08-25", plot: "P3", grade: "A", weightKg: 260, pricePerKg: 29, total: 7540, buyer: "Ahmad Traders" },
  { id: "S009", date: "2024-09-05", plot: "P1", grade: "A", weightKg: 380, pricePerKg: 31, total: 11780, buyer: "KL Export Co." },
  { id: "S010", date: "2024-09-12", plot: "P2", grade: "A", weightKg: 220, pricePerKg: 28, total: 6160, buyer: "Ahmad Traders" },
  { id: "S011", date: "2024-09-20", plot: "P3", grade: "B", weightKg: 310, pricePerKg: 13, total: 4030, buyer: "Local Market" },
  { id: "S012", date: "2024-10-02", plot: "P1", grade: "A", weightKg: 450, pricePerKg: 33, total: 14850, buyer: "KL Export Co." },
]

export const MAINTENANCE_LOGS: MaintenanceLog[] = [
  { id: "M001", date: "2024-07-10", plot: "P1", activity: "Watering", durationMin: 45, notes: "Soil dry, applied extra 20 min." },
  { id: "M002", date: "2024-07-11", plot: "P2", activity: "Fertilizing", notes: "Applied NPK 15-15-15 blend." },
  { id: "M003", date: "2024-07-13", plot: "P3", activity: "Pest Control", durationMin: 60, notes: "Spotted fruit borers, applied pesticide." },
  { id: "M004", date: "2024-07-15", plot: "P1", activity: "Watering", durationMin: 40, notes: "Normal schedule." },
  { id: "M005", date: "2024-07-18", plot: "P2", activity: "Pruning", durationMin: 90, notes: "Removed dead branches, improved airflow." },
  { id: "M006", date: "2024-07-22", plot: "P3", activity: "Watering", durationMin: 50, notes: "After harvesting session." },
  { id: "M007", date: "2024-07-25", plot: "P4", activity: "Watering", durationMin: 30, notes: "Young trees need daily watering." },
  { id: "M008", date: "2024-07-28", plot: "P1", activity: "Fertilizing", notes: "Pre-harvest boosting fertilizer applied." },
  { id: "M009", date: "2024-08-01", plot: "P2", activity: "Watering", durationMin: 40, notes: "Normal schedule." },
  { id: "M010", date: "2024-08-05", plot: "P3", activity: "Pruning", durationMin: 75, notes: "Heavy pruning — overgrown canopy." },
  { id: "M011", date: "2024-08-10", plot: "P1", activity: "Watering", durationMin: 45, notes: "Dry season, increased frequency." },
  { id: "M012", date: "2024-08-14", plot: "P4", activity: "Fertilizing", notes: "First fertilizer application for young trees." },
  { id: "M013", date: "2024-08-20", plot: "P2", activity: "Pest Control", durationMin: 60, notes: "Preventive spray." },
  { id: "M014", date: "2024-08-25", plot: "P3", activity: "Watering", durationMin: 50, notes: "Post-rain check." },
  { id: "M015", date: "2024-09-03", plot: "P1", activity: "Fertilizing", notes: "Monthly NPK application." },
]

// Monthly revenue for chart (Jan–Dec 2024)
export const MONTHLY_REVENUE = [
  { month: "Jan", revenue: 0 },
  { month: "Feb", revenue: 0 },
  { month: "Mar", revenue: 0 },
  { month: "Apr", revenue: 0 },
  { month: "May", revenue: 4200 },
  { month: "Jun", revenue: 8900 },
  { month: "Jul", revenue: 29150 },
  { month: "Aug", revenue: 39270 },
  { month: "Sep", revenue: 21970 },
  { month: "Oct", revenue: 14850 },
  { month: "Nov", revenue: 0 },
  { month: "Dec", revenue: 0 },
]

export function getPlotName(id: string): string {
  return PLOTS.find((p) => p.id === id)?.name ?? id
}
