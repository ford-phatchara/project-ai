"use client"

import { useState } from "react"
import { Plus, Pencil, Trash2, Filter, ChevronDown } from "lucide-react"
import { PageHeader } from "@/components/page-header"
import { GradeBadge } from "@/components/grade-badge"
import { SALES, PLOTS, getPlotName, type Sale, type Grade } from "@/lib/mock-data"
import { cn } from "@/lib/utils"

const EMPTY_SALE: Omit<Sale, "id" | "total"> = {
  date: new Date().toISOString().slice(0, 10),
  plot: "P1",
  grade: "A",
  weightKg: 0,
  pricePerKg: 0,
}

export default function SalesPage() {
  const [sales, setSales] = useState<Sale[]>(SALES)
  const [filterGrade, setFilterGrade] = useState<Grade | "All">("All")
  const [filterPlot, setFilterPlot] = useState<string>("All")
  const [modalOpen, setModalOpen] = useState(false)
  const [editSale, setEditSale] = useState<Sale | null>(null)
  const [deleteId, setDeleteId] = useState<string | null>(null)
  const [form, setForm] = useState(EMPTY_SALE)
  const [showFilters, setShowFilters] = useState(false)

  const filtered = sales.filter((s) => {
    if (filterGrade !== "All" && s.grade !== filterGrade) return false
    if (filterPlot !== "All" && s.plot !== filterPlot) return false
    return true
  })

  function openAdd() {
    setEditSale(null)
    setForm(EMPTY_SALE)
    setModalOpen(true)
  }

  function openEdit(sale: Sale) {
    setEditSale(sale)
    setForm({ date: sale.date, plot: sale.plot, grade: sale.grade, weightKg: sale.weightKg, pricePerKg: sale.pricePerKg, buyer: sale.buyer })
    setModalOpen(true)
  }

  function handleSave() {
    const total = form.weightKg * form.pricePerKg
    if (editSale) {
      setSales((prev) => prev.map((s) => s.id === editSale.id ? { ...editSale, ...form, total } : s))
    } else {
      const newSale: Sale = { ...form, id: `S${Date.now()}`, total }
      setSales((prev) => [newSale, ...prev])
    }
    setModalOpen(false)
  }

  function handleDelete(id: string) {
    setSales((prev) => prev.filter((s) => s.id !== id))
    setDeleteId(null)
  }

  return (
    <div>
      <PageHeader
        title="Sales"
        subtitle={`${filtered.length} records`}
        action={
          <button
            onClick={openAdd}
            className="flex items-center gap-2 bg-primary text-primary-foreground rounded-xl px-4 py-3 text-sm font-semibold shadow-sm active:scale-95 transition-transform"
          >
            <Plus className="w-5 h-5" />
            Add Sale
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
            <div className="flex items-center gap-2">
              <Filter className="w-4 h-4 text-muted-foreground" />
              Filter Sales
            </div>
            <ChevronDown className={cn("w-4 h-4 text-muted-foreground transition-transform", showFilters && "rotate-180")} />
          </button>
          {showFilters && (
            <div className="px-4 pb-4 flex flex-col gap-3 border-t border-border pt-3">
              <div>
                <label className="text-xs font-medium text-muted-foreground uppercase tracking-wide">Grade</label>
                <div className="flex gap-2 mt-2">
                  {["All", "A", "B"].map((g) => (
                    <button
                      key={g}
                      onClick={() => setFilterGrade(g as Grade | "All")}
                      className={cn(
                        "flex-1 py-2.5 rounded-xl text-sm font-medium border transition-colors",
                        filterGrade === g
                          ? "bg-primary text-primary-foreground border-primary"
                          : "bg-muted text-muted-foreground border-border"
                      )}
                    >
                      {g === "All" ? "All Grades" : `Grade ${g}`}
                    </button>
                  ))}
                </div>
              </div>
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
            </div>
          )}
        </div>

        {/* Sales List */}
        {filtered.length === 0 ? (
          <div className="bg-card rounded-2xl border border-border shadow-sm p-10 flex flex-col items-center gap-3 text-center">
            <div className="w-14 h-14 rounded-2xl bg-muted flex items-center justify-center">
              <Filter className="w-7 h-7 text-muted-foreground" />
            </div>
            <p className="font-semibold text-foreground">No records found</p>
            <p className="text-sm text-muted-foreground">Try adjusting your filters or add a new sale.</p>
          </div>
        ) : (
          <div className="flex flex-col gap-3">
            {filtered.sort((a, b) => b.date.localeCompare(a.date)).map((sale) => (
              <div key={sale.id} className="bg-card rounded-2xl border border-border shadow-sm p-4">
                <div className="flex items-start justify-between gap-2">
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2 flex-wrap">
                      <GradeBadge grade={sale.grade} />
                      <p className="text-sm font-semibold text-foreground truncate">{getPlotName(sale.plot)}</p>
                    </div>
                    <p className="text-xs text-muted-foreground mt-1">{sale.date}</p>
                    <div className="flex flex-wrap gap-3 mt-2">
                      <span className="text-xs text-muted-foreground">{sale.weightKg} kg</span>
                      <span className="text-xs text-muted-foreground">RM {sale.pricePerKg}/kg</span>
                      {sale.buyer && <span className="text-xs text-muted-foreground">{sale.buyer}</span>}
                    </div>
                  </div>
                  <div className="flex flex-col items-end gap-3 shrink-0">
                    <p className="text-base font-bold text-primary">RM {sale.total.toLocaleString()}</p>
                    <div className="flex gap-1">
                      <button
                        onClick={() => openEdit(sale)}
                        className="w-9 h-9 rounded-xl bg-muted flex items-center justify-center active:scale-95 transition-transform"
                        aria-label="Edit sale"
                      >
                        <Pencil className="w-4 h-4 text-muted-foreground" />
                      </button>
                      <button
                        onClick={() => setDeleteId(sale.id)}
                        className="w-9 h-9 rounded-xl bg-destructive/10 flex items-center justify-center active:scale-95 transition-transform"
                        aria-label="Delete sale"
                      >
                        <Trash2 className="w-4 h-4 text-destructive" />
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      {/* Add/Edit Modal */}
      {modalOpen && (
        <div className="fixed inset-0 z-50 flex items-end sm:items-center justify-center bg-black/50 backdrop-blur-sm" onClick={() => setModalOpen(false)}>
          <div
            className="bg-card w-full max-w-lg rounded-t-3xl sm:rounded-3xl shadow-2xl p-6 max-h-[90vh] overflow-y-auto"
            onClick={(e) => e.stopPropagation()}
          >
            <h2 className="text-lg font-bold text-foreground mb-5">{editSale ? "Edit Sale" : "Add New Sale"}</h2>
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

              <FormField label="Grade">
                <div className="flex gap-3">
                  {(["A", "B"] as Grade[]).map((g) => (
                    <button
                      key={g}
                      type="button"
                      onClick={() => setForm({ ...form, grade: g })}
                      className={cn(
                        "flex-1 py-3.5 rounded-xl text-sm font-bold border-2 transition-colors",
                        form.grade === g
                          ? "bg-primary text-primary-foreground border-primary"
                          : "bg-muted text-muted-foreground border-border"
                      )}
                    >
                      Grade {g}
                    </button>
                  ))}
                </div>
              </FormField>

              <FormField label="Weight (kg)">
                <input
                  type="number"
                  inputMode="decimal"
                  min={0}
                  value={form.weightKg || ""}
                  onChange={(e) => setForm({ ...form, weightKg: parseFloat(e.target.value) || 0 })}
                  placeholder="e.g. 250"
                  className="farm-input"
                />
              </FormField>

              <FormField label="Price per kg (RM)">
                <input
                  type="number"
                  inputMode="decimal"
                  min={0}
                  value={form.pricePerKg || ""}
                  onChange={(e) => setForm({ ...form, pricePerKg: parseFloat(e.target.value) || 0 })}
                  placeholder="e.g. 28"
                  className="farm-input"
                />
              </FormField>

              {form.weightKg > 0 && form.pricePerKg > 0 && (
                <div className="bg-primary/8 rounded-xl px-4 py-3 flex justify-between items-center">
                  <span className="text-sm text-muted-foreground font-medium">Total</span>
                  <span className="text-lg font-bold text-primary">RM {(form.weightKg * form.pricePerKg).toLocaleString()}</span>
                </div>
              )}

              <FormField label="Buyer (optional)">
                <input
                  type="text"
                  value={form.buyer ?? ""}
                  onChange={(e) => setForm({ ...form, buyer: e.target.value })}
                  placeholder="e.g. Ahmad Traders"
                  className="farm-input"
                />
              </FormField>
            </div>

            <div className="flex gap-3 mt-6 mb-20">
              <button
                onClick={() => setModalOpen(false)}
                className="flex-1 py-4 rounded-xl border border-border text-foreground font-semibold text-sm active:scale-95 transition-transform"
              >
                Cancel
              </button>
              <button
                onClick={handleSave}
                disabled={!form.weightKg || !form.pricePerKg}
                className="flex-1 py-4 rounded-xl bg-primary text-primary-foreground font-semibold text-sm shadow-sm active:scale-95 transition-transform disabled:opacity-40"
              >
                {editSale ? "Save Changes" : "Add Sale"}
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Delete Confirm */}
      {deleteId && (
        <div className="fixed inset-0 z-50 flex items-end sm:items-center justify-center bg-black/50 backdrop-blur-sm">
          <div className="bg-card w-full max-w-sm rounded-t-3xl sm:rounded-3xl shadow-2xl p-6">
            <h2 className="text-lg font-bold text-foreground mb-2">Delete Sale?</h2>
            <p className="text-sm text-muted-foreground mb-6">This action cannot be undone.</p>
            <div className="flex gap-3">
              <button onClick={() => setDeleteId(null)} className="flex-1 py-4 rounded-xl border border-border text-foreground font-semibold text-sm">
                Cancel
              </button>
              <button onClick={() => handleDelete(deleteId)} className="flex-1 py-4 rounded-xl bg-destructive text-destructive-foreground font-semibold text-sm active:scale-95">
                Delete
              </button>
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
