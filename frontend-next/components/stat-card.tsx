import { cn } from "@/lib/utils"

interface StatCardProps {
  label: string
  value: string
  sub?: string
  icon?: React.ReactNode
  accent?: boolean
}

export function StatCard({ label, value, sub, icon, accent }: StatCardProps) {
  return (
    <div
      className={cn(
        "rounded-2xl p-4 shadow-sm border",
        accent
          ? "bg-primary text-primary-foreground border-primary"
          : "bg-card text-card-foreground border-border"
      )}
    >
      <div className="flex items-start justify-between gap-2">
        <p className={cn("text-xs font-medium uppercase tracking-wide", accent ? "text-primary-foreground/70" : "text-muted-foreground")}>
          {label}
        </p>
        {icon && (
          <div className={cn("w-8 h-8 rounded-xl flex items-center justify-center shrink-0", accent ? "bg-primary-foreground/15" : "bg-muted")}>
            {icon}
          </div>
        )}
      </div>
      <p className={cn("text-3xl font-bold mt-2 leading-none", accent ? "text-primary-foreground" : "text-foreground")}>
        {value}
      </p>
      {sub && (
        <p className={cn("text-xs mt-1.5", accent ? "text-primary-foreground/60" : "text-muted-foreground")}>
          {sub}
        </p>
      )}
    </div>
  )
}
