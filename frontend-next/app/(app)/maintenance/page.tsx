"use client"

import { useState } from "react"
import { Plus, Pencil, Trash2, Droplets, Sprout, Scissors, Bug, Leaf, ChevronDown } from "lucide-react"
import { PageHeader } from "@/components/page-header"
import { MAINTENANCE_LOGS, PLOTS, getPlotName, type MaintenanceLog } from "@/lib/mock-data"
import { cn } from "@/lib/utils"

type Activity = MaintenanceLog["activity"]

const ACTIVITIES: Activity[] = ["Watering", "Fertilizing", "Pruning", "Pest Control", "Harvesting"]

const ACTIVITY_META: Record<Activity, { icon: React.ElementType; color: string; bg: string }> = {
  Watering:     { icon: Droplets, color: "text-blue-600", bg: "bg-blue-50" },
  Fertilizing:  { icon: Sprout,   color: "text-green-600", bg: "bg-green-50" },
  Pruning:      { icon: Scissors, color: "text-orange-600", bg: "bg-orange-50" },
  "Pest Control": { icon: Bug,    color: "text-red-600", bg: "bg-red-50" },
  Harvesting:   { icon: Leaf,     color: "text-primary", bg: "bg-primary/10" },
}

const EMPTY_LOG: Omit<MaintenanceLog, "id"> = {
  date: new Date().toISOString().slice(0, 10),
  plot: "P1",
  activity: "Watering",
  durationMin: undefined,
  notes: "",
}

export default function MaintenancePage() {
  const [logs, setLogs] = useState<MaintenanceLog[]>(MAINTENANCE_LOGS)
  const [filterPlot, setFilterPlot] = useState("All")
  const [filterActivity, setFilterActivity] = useState<Activity | "All">("All")
  const [modalOpen, setModalOpen] = useState(false)
  const [editLog, setEditLog] = useState<MaintenanceLog | null>(null)
  const [deleteId, setDeleteId] = useState<string | null>(null)
  const [form, setForm] = useState<Omit<MaintenanceLog, "id">>(EMPTY_LOG)
  const [showFilters, setShowFilters] = useState(false)

  const filtered = logs.filter((l) => {
    if (filterPlot !== "All" && l.plot !== filterPlot) return false
    if (filterActivity !== "All" && l.activity !== filterActivity) return false
    return true
  })

  function openAdd() {
    setEditLog(null)
    setForm(EMPTY_LOG)
    setModalOpen(true)
  }

  function openEdit(log: MaintenanceLog) {
    setEditLog(log)
    setForm({ date: log.date, plot: log.plot, activity: log.activity, durationMin: log.durationMin, notes: log.notes })
    setModalOpen(true)
  }

  function handleSave() {
    if (editLog) {
      setLogs((prev) => prev.map((l) => l.id === editLog.id ? { ...editLog, ...form } : l))
    } else {
      setLogs((prev) => [{ ...form, id: `M${Date.now()}` }, ...prev])
    }
    setModalOpen(false)
  }

  function handleDelete(id: string) {
    setLogs((prev) => prev.filter((l) => l.id !== id))
    setDeleteId(null)
  }

  return (
    <div>
      <PageHeader
        title="Maintenance"
        subtitle={`${filtered.length} care logs`}
        action={
          <button
            onClick={openAdd}
            className="flex items-center gap-2 bg-primary text-primary-foreground rounded-xl px-4 py-3 text-sm font-semibold shadow-sm active:scale-95 transition-transform"
          >
            <Plus className="w-5 h-5" />
            Log Activity
          </button>
        }
      />

      <div className="px-4 lg:px-8 flex flex-col gap-4">
        {/* Filters */}
        <div className="bg-card rounded-2xl border border-border shadow-sm overflow-hidden">
          <button
            onClick={() => setShowFilters((v) => !v)}
            className="flex items-center justify-between w-full px-4 py-3 text-sm font-medium text-foreground"
          >
            <span className="text-sm font-medium text-foreground">Filter Logs</span>
            <ChevronDown className={cn("w-4 h-4 text-muted-foreground transition-transform", showFilters && "rotate-180")} />
          </button>
          {showFilters && (
            <div className="px-4 pb-4 flex flex-col gap-3 border-t border-border pt-3">
              <div>
                <label className="text-xs font-medium text-muted-foreground uppercase tracking-wide">Plot</label>
                <select
                  value={filterPlot}
                  onChange={(e) => setFilterPlot(e.target.value)}
                  className="w-full mt-2 border border-input bg-background rounded-xl px-3 py-3 text-sm text-foreground"
                >
                  <option value="All">All Plots</option>
                  {PLOTS.map((p) => <option key={p.id} value={p.id}>{p.name}</option>)}
                </select>
              </div>
              <div>
                <label className="text-xs font-medium text-muted-foreground uppercase tracking-wide mb-2 block">Activity</label>
                <div className="flex flex-wrap gap-2">
                  {(["All", ...ACTIVITIES] as (Activity | "All")[]).map((a) => (
                    <button
                      key={a}
                      onClick={() => setFilterActivity(a)}
                      className={cn(
                        "px-3 py-2 rounded-xl text-xs font-medium border transition-colors",
                        filterActivity === a
                          ? "bg-primary text-primary-foreground border-primary"
                          : "bg-muted text-muted-foreground border-border"
                      )}
                    >
                      {a}
                    </button>
                  ))}
                </div>
              </div>
            </div>
          )}
        </div>

        {/* Log List */}
        {filtered.length === 0 ? (
          <div className="bg-card rounded-2xl border border-border shadow-sm p-10 flex flex-col items-center gap-3 text-center">
            <div className="w-14 h-14 rounded-2xl bg-muted flex items-center justify-center">
              <Leaf className="w-7 h-7 text-muted-foreground" />
            </div>
            <p className="font-semibold text-foreground">No logs found</p>
            <p className="text-sm text-muted-foreground">Try adjusting filters or log a new activity.</p>
          </div>
        ) : (
          <div className="flex flex-col gap-3">
            {filtered.sort((a, b) => b.date.localeCompare(a.date)).map((log) => {
              const meta = ACTIVITY_META[log.activity]
              const Icon = meta.icon
              return (
                <div key={log.id} className="bg-card rounded-2xl border border-border shadow-sm p-4">
                  <div className="flex items-start justify-between gap-2">
                    <div className="flex gap-3 items-start flex-1 min-w-0">
                      <div className={cn("w-11 h-11 rounded-xl flex items-center justify-center shrink-0", meta.bg)}>
                        <Icon className={cn("w-5 h-5", meta.color)} />
                      </div>
                      <div className="flex-1 min-w-0">
                        <div className="flex items-center gap-2 flex-wrap">
                          <p className="text-sm font-semibold text-foreground">{log.activity}</p>
                          <span className="text-xs bg-muted text-muted-foreground rounded-lg px-2 py-0.5">{getPlotName(log.plot)}</span>
                        </div>
                        <p className="text-xs text-muted-foreground mt-0.5">{log.date}{log.durationMin ? ` · ${log.durationMin} min` : ""}</p>
                        {log.notes && <p className="text-xs text-muted-foreground mt-1.5 leading-relaxed line-clamp-2">{log.notes}</p>}
                      </div>
                    </div>
                    <div className="flex gap-1 shrink-0">
                      <button
                        onClick={() => openEdit(log)}
                        className="w-10 h-10 rounded-xl bg-muted flex items-center justify-center active:scale-95 transition-transform"
                        aria-label="Edit log"
                      >
                        <Pencil className="w-4 h-4 text-muted-foreground" />
                      </button>
                      <button
                        onClick={() => setDeleteId(log.id)}
                        className="w-10 h-10 rounded-xl bg-destructive/10 flex items-center justify-center active:scale-95 transition-transform"
                        aria-label="Delete log"
                      >
                        <Trash2 className="w-4 h-4 text-destructive" />
                      </button>
                    </div>
                  </div>
                </div>
              )
            })}
          </div>
        )}
      </div>

      {/* Modal */}
      {modalOpen && (
        <div className="fixed inset-0 z-50 flex items-end sm:items-center justify-center bg-black/50 backdrop-blur-sm" onClick={() => setModalOpen(false)}>
          <div
            className="bg-card w-full max-w-lg rounded-t-3xl sm:rounded-3xl shadow-2xl p-6 max-h-[90vh] overflow-y-auto"
            onClick={(e) => e.stopPropagation()}
          >
            <h2 className="text-lg font-bold text-foreground mb-5">{editLog ? "Edit Log" : "Log Activity"}</h2>
            <div className="flex flex-col gap-4">
              <FormField label="Date">
                <input
                  type="date"
                  value={form.date}
                  onChange={(e) => setForm({ ...form, date: e.target.value })}
                  className="farm-input"
                />
              </FormField>

              <FormField label="Plot">
                <select
                  value={form.plot}
                  onChange={(e) => setForm({ ...form, plot: e.target.value })}
                  className="farm-input"
                >
                  {PLOTS.map((p) => <option key={p.id} value={p.id}>{p.name}</option>)}
                </select>
              </FormField>

              <FormField label="Activity">
                <div className="grid grid-cols-2 gap-2">
                  {ACTIVITIES.map((a) => {
                    const Icon = ACTIVITY_META[a].icon
                    return (
                      <button
                        key={a}
                        type="button"
                        onClick={() => setForm({ ...form, activity: a })}
                        className={cn(
                          "flex items-center gap-2 px-3 py-3 rounded-xl text-sm font-medium border-2 transition-colors",
                          form.activity === a
                            ? "bg-primary text-primary-foreground border-primary"
                            : "bg-muted text-muted-foreground border-border"
                        )}
                      >
                        <Icon className="w-4 h-4 shrink-0" />
                        {a}
                      </button>
                    )
                  })}
                </div>
              </FormField>

              {(form.activity === "Watering" || form.activity === "Pruning" || form.activity === "Pest Control") && (
                <FormField label="Duration (minutes)">
                  <input
                    type="number"
                    inputMode="numeric"
                    min={0}
                    value={form.durationMin ?? ""}
                    onChange={(e) => setForm({ ...form, durationMin: parseInt(e.target.value) || undefined })}
                    placeholder="e.g. 45"
                    className="farm-input"
                  />
                </FormField>
              )}

              <FormField label="Notes">
                <textarea
                  value={form.notes}
                  onChange={(e) => setForm({ ...form, notes: e.target.value })}
                  placeholder="What was done? Any observations?"
                  rows={3}
                  className="farm-input resize-none"
                />
              </FormField>
            </div>

            <div className="flex gap-3 mt-6">
              <button onClick={() => setModalOpen(false)} className="flex-1 py-4 rounded-xl border border-border text-foreground font-semibold text-sm active:scale-95 transition-transform">
                Cancel
              </button>
              <button
                onClick={handleSave}
                className="flex-1 py-4 rounded-xl bg-primary text-primary-foreground font-semibold text-sm shadow-sm active:scale-95 transition-transform"
              >
                {editLog ? "Save Changes" : "Log Activity"}
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Delete Confirm */}
      {deleteId && (
        <div className="fixed inset-0 z-50 flex items-end sm:items-center justify-center bg-black/50 backdrop-blur-sm">
          <div className="bg-card w-full max-w-sm rounded-t-3xl sm:rounded-3xl shadow-2xl p-6">
            <h2 className="text-lg font-bold text-foreground mb-2">Delete Log?</h2>
            <p className="text-sm text-muted-foreground mb-6">This action cannot be undone.</p>
            <div className="flex gap-3">
              <button onClick={() => setDeleteId(null)} className="flex-1 py-4 rounded-xl border border-border text-foreground font-semibold text-sm">Cancel</button>
              <button onClick={() => handleDelete(deleteId)} className="flex-1 py-4 rounded-xl bg-destructive text-destructive-foreground font-semibold text-sm active:scale-95">Delete</button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}

function FormField({ label, children }: { label: string; children: React.ReactNode }) {
  return (
    <div className="flex flex-col gap-1.5">
      <label className="text-xs font-semibold text-muted-foreground uppercase tracking-wide">{label}</label>
      {children}
    </div>
  )
}
