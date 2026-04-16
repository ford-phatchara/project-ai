"use client"

import { useState } from "react"
import { ChevronLeft, ChevronRight, Droplets, Sprout, Scissors, Bug, Leaf } from "lucide-react"
import { PageHeader } from "@/components/page-header"
import { MAINTENANCE_LOGS, getPlotName, type MaintenanceLog } from "@/lib/mock-data"
import { cn } from "@/lib/utils"

type Activity = MaintenanceLog["activity"]

const ACTIVITY_META: Record<Activity, { icon: React.ElementType; color: string; bg: string; dot: string }> = {
  Watering:       { icon: Droplets, color: "text-blue-700", bg: "bg-blue-50", dot: "bg-blue-500" },
  Fertilizing:    { icon: Sprout,   color: "text-green-700", bg: "bg-green-50", dot: "bg-green-500" },
  Pruning:        { icon: Scissors, color: "text-orange-700", bg: "bg-orange-50", dot: "bg-orange-500" },
  "Pest Control": { icon: Bug,      color: "text-red-700", bg: "bg-red-50", dot: "bg-red-500" },
  Harvesting:     { icon: Leaf,     color: "text-primary", bg: "bg-primary/10", dot: "bg-primary" },
}

const DAYS = ["Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"]
const MONTHS = ["January","February","March","April","May","June","July","August","September","October","November","December"]

function getDaysInMonth(year: number, month: number) {
  return new Date(year, month + 1, 0).getDate()
}
function getFirstDayOfMonth(year: number, month: number) {
  return new Date(year, month, 1).getDay()
}

export default function CalendarPage() {
  const today = new Date()
  const [year, setYear] = useState(2024)
  const [month, setMonth] = useState(6) // July 2024 — where mock data starts
  const [selectedDate, setSelectedDate] = useState<string | null>(null)

  const daysInMonth = getDaysInMonth(year, month)
  const firstDay = getFirstDayOfMonth(year, month)

  // Map date string -> logs
  const logsByDate: Record<string, MaintenanceLog[]> = {}
  MAINTENANCE_LOGS.forEach((l) => {
    if (!logsByDate[l.date]) logsByDate[l.date] = []
    logsByDate[l.date].push(l)
  })

  function prevMonth() {
    if (month === 0) { setMonth(11); setYear((y) => y - 1) }
    else setMonth((m) => m - 1)
    setSelectedDate(null)
  }
  function nextMonth() {
    if (month === 11) { setMonth(0); setYear((y) => y + 1) }
    else setMonth((m) => m + 1)
    setSelectedDate(null)
  }

  function formatDate(day: number) {
    return `${year}-${String(month + 1).padStart(2, "0")}-${String(day).padStart(2, "0")}`
  }

  const selectedLogs = selectedDate ? (logsByDate[selectedDate] ?? []) : []

  // Month-level summary
  const monthLogs = MAINTENANCE_LOGS.filter((l) => {
    const [y, m] = l.date.split("-").map(Number)
    return y === year && m === month + 1
  })

  return (
    <div>
      <PageHeader title="Calendar" subtitle="Farm activity schedule" />

      <div className="px-4 lg:px-8 flex flex-col gap-4">
        {/* Month Navigator */}
        <div className="bg-card rounded-2xl border border-border shadow-sm p-4">
          <div className="flex items-center justify-between mb-4">
            <button
              onClick={prevMonth}
              className="w-11 h-11 rounded-xl bg-muted flex items-center justify-center active:scale-95 transition-transform"
              aria-label="Previous month"
            >
              <ChevronLeft className="w-5 h-5 text-foreground" />
            </button>
            <h2 className="text-base font-bold text-foreground">
              {MONTHS[month]} {year}
            </h2>
            <button
              onClick={nextMonth}
              className="w-11 h-11 rounded-xl bg-muted flex items-center justify-center active:scale-95 transition-transform"
              aria-label="Next month"
            >
              <ChevronRight className="w-5 h-5 text-foreground" />
            </button>
          </div>

          {/* Day Headers */}
          <div className="grid grid-cols-7 mb-2">
            {DAYS.map((d) => (
              <div key={d} className="text-center text-[11px] font-semibold text-muted-foreground py-1">
                {d}
              </div>
            ))}
          </div>

          {/* Calendar Grid */}
          <div className="grid grid-cols-7 gap-y-1">
            {Array.from({ length: firstDay }).map((_, i) => (
              <div key={`empty-${i}`} />
            ))}
            {Array.from({ length: daysInMonth }).map((_, i) => {
              const day = i + 1
              const dateStr = formatDate(day)
              const hasLogs = !!logsByDate[dateStr]
              const logsForDay = logsByDate[dateStr] ?? []
              const isToday = dateStr === today.toISOString().slice(0, 10)
              const isSelected = dateStr === selectedDate

              return (
                <button
                  key={day}
                  onClick={() => setSelectedDate(isSelected ? null : dateStr)}
                  className={cn(
                    "relative flex flex-col items-center justify-start py-2 rounded-xl min-h-[3rem] transition-colors",
                    isSelected
                      ? "bg-primary text-primary-foreground"
                      : isToday
                      ? "bg-accent/30 text-foreground"
                      : hasLogs
                      ? "bg-muted text-foreground"
                      : "text-foreground hover:bg-muted/50"
                  )}
                >
                  <span className={cn("text-sm font-medium", isSelected && "text-primary-foreground")}>
                    {day}
                  </span>
                  {hasLogs && !isSelected && (
                    <div className="flex gap-0.5 mt-0.5 flex-wrap justify-center max-w-[2rem]">
                      {logsForDay.slice(0, 3).map((l, idx) => (
                        <span key={idx} className={cn("w-1.5 h-1.5 rounded-full", ACTIVITY_META[l.activity].dot)} />
                      ))}
                    </div>
                  )}
                </button>
              )
            })}
          </div>
        </div>

        {/* Legend */}
        <div className="bg-card rounded-2xl border border-border shadow-sm p-4">
          <p className="text-xs font-semibold text-muted-foreground uppercase tracking-wide mb-3">Activity Legend</p>
          <div className="flex flex-wrap gap-3">
            {Object.entries(ACTIVITY_META).map(([name, meta]) => (
              <div key={name} className="flex items-center gap-1.5">
                <span className={cn("w-2.5 h-2.5 rounded-full shrink-0", meta.dot)} />
                <span className="text-xs text-muted-foreground">{name}</span>
              </div>
            ))}
          </div>
        </div>

        {/* Selected Day Detail */}
        {selectedDate && (
          <div className="bg-card rounded-2xl border border-border shadow-sm p-4">
            <p className="text-sm font-semibold text-foreground mb-3">
              {selectedDate}
            </p>
            {selectedLogs.length === 0 ? (
              <p className="text-sm text-muted-foreground">No activities logged for this day.</p>
            ) : (
              <div className="flex flex-col gap-3">
                {selectedLogs.map((log) => {
                  const meta = ACTIVITY_META[log.activity]
                  const Icon = meta.icon
                  return (
                    <div key={log.id} className={cn("flex items-start gap-3 p-3 rounded-xl", meta.bg)}>
                      <div className="w-9 h-9 rounded-lg bg-white/60 flex items-center justify-center shrink-0">
                        <Icon className={cn("w-5 h-5", meta.color)} />
                      </div>
                      <div>
                        <p className={cn("text-sm font-semibold", meta.color)}>{log.activity}</p>
                        <p className="text-xs text-muted-foreground">{getPlotName(log.plot)}{log.durationMin ? ` · ${log.durationMin} min` : ""}</p>
                        {log.notes && <p className="text-xs text-muted-foreground mt-1">{log.notes}</p>}
                      </div>
                    </div>
                  )
                })}
              </div>
            )}
          </div>
        )}

        {/* Month Summary */}
        <div className="bg-card rounded-2xl border border-border shadow-sm p-4">
          <p className="text-sm font-semibold text-foreground mb-3">{MONTHS[month]} Summary</p>
          {monthLogs.length === 0 ? (
            <p className="text-sm text-muted-foreground">No activities recorded this month.</p>
          ) : (
            <div className="flex flex-col gap-2">
              {(Object.keys(ACTIVITY_META) as Activity[]).map((act) => {
                const count = monthLogs.filter((l) => l.activity === act).length
                if (!count) return null
                const meta = ACTIVITY_META[act]
                const Icon = meta.icon
                return (
                  <div key={act} className="flex items-center gap-3">
                    <div className={cn("w-8 h-8 rounded-lg flex items-center justify-center", meta.bg)}>
                      <Icon className={cn("w-4 h-4", meta.color)} />
                    </div>
                    <span className="text-sm text-foreground flex-1">{act}</span>
                    <span className="text-sm font-semibold text-foreground">{count}×</span>
                  </div>
                )
              })}
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
