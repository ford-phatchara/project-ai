"use client"

import { useState } from "react"
import { Plus, Pencil, Trash2, Trees } from "lucide-react"
import { PageHeader } from "@/components/page-header"
import { PLOTS, type Plot } from "@/lib/mock-data"
import { cn } from "@/lib/utils"

const EMPTY_PLOT: Omit<Plot, "id"> = { name: "", sizeMeter: 0, trees: 0, notes: "" }

export default function PlotsPage() {
  const [plots, setPlots] = useState<Plot[]>(PLOTS)
  const [modalOpen, setModalOpen] = useState(false)
  const [editPlot, setEditPlot] = useState<Plot | null>(null)
  const [deleteId, setDeleteId] = useState<string | null>(null)
  const [form, setForm] = useState(EMPTY_PLOT)

  function openAdd() {
    setEditPlot(null)
    setForm(EMPTY_PLOT)
    setModalOpen(true)
  }

  function openEdit(plot: Plot) {
    setEditPlot(plot)
    setForm({ name: plot.name, sizeMeter: plot.sizeMeter, trees: plot.trees, notes: plot.notes })
    setModalOpen(true)
  }

  function handleSave() {
    if (editPlot) {
      setPlots((prev) => prev.map((p) => p.id === editPlot.id ? { ...editPlot, ...form } : p))
    } else {
      setPlots((prev) => [...prev, { ...form, id: `P${Date.now()}` }])
    }
    setModalOpen(false)
  }

  function handleDelete(id: string) {
    setPlots((prev) => prev.filter((p) => p.id !== id))
    setDeleteId(null)
  }

  return (
    <div>
      <PageHeader
        title="Plots"
        subtitle={`${plots.length} farm plots`}
        action={
          <button
            onClick={openAdd}
            className="flex items-center gap-2 bg-primary text-primary-foreground rounded-xl px-4 py-3 text-sm font-semibold shadow-sm active:scale-95 transition-transform"
          >
            <Plus className="w-5 h-5" />
            Add Plot
          </button>
        }
      />

      <div className="px-4 lg:px-8 flex flex-col gap-3">
        {plots.length === 0 ? (
          <div className="bg-card rounded-2xl border border-border shadow-sm p-10 flex flex-col items-center gap-3 text-center">
            <div className="w-14 h-14 rounded-2xl bg-muted flex items-center justify-center">
              <Trees className="w-7 h-7 text-muted-foreground" />
            </div>
            <p className="font-semibold text-foreground">No plots yet</p>
            <p className="text-sm text-muted-foreground">Add your first farm plot to get started.</p>
            <button onClick={openAdd} className="mt-2 bg-primary text-primary-foreground rounded-xl px-5 py-3 text-sm font-semibold">
              Add First Plot
            </button>
          </div>
        ) : (
          plots.map((plot) => (
            <div key={plot.id} className="bg-card rounded-2xl border border-border shadow-sm p-4">
              <div className="flex items-start justify-between gap-2">
                <div className="flex gap-3 items-start flex-1 min-w-0">
                  <div className="w-11 h-11 rounded-xl bg-primary/10 flex items-center justify-center shrink-0">
                    <Trees className="w-6 h-6 text-primary" />
                  </div>
                  <div className="flex-1 min-w-0">
                    <p className="font-semibold text-foreground text-sm leading-tight">{plot.name}</p>
                    <div className="flex flex-wrap gap-3 mt-1.5">
                      <span className="text-xs text-muted-foreground">{plot.sizeMeter.toLocaleString()} m²</span>
                      <span className="text-xs text-muted-foreground">{plot.trees} trees</span>
                    </div>
                    {plot.notes && (
                      <p className="text-xs text-muted-foreground mt-2 leading-relaxed line-clamp-2">{plot.notes}</p>
                    )}
                  </div>
                </div>
                <div className="flex gap-1 shrink-0">
                  <button
                    onClick={() => openEdit(plot)}
                    className="w-10 h-10 rounded-xl bg-muted flex items-center justify-center active:scale-95 transition-transform"
                    aria-label="Edit plot"
                  >
                    <Pencil className="w-4 h-4 text-muted-foreground" />
                  </button>
                  <button
                    onClick={() => setDeleteId(plot.id)}
                    className="w-10 h-10 rounded-xl bg-destructive/10 flex items-center justify-center active:scale-95 transition-transform"
                    aria-label="Delete plot"
                  >
                    <Trash2 className="w-4 h-4 text-destructive" />
                  </button>
                </div>
              </div>
            </div>
          ))
        )}
      </div>

      {/* Modal */}
      {modalOpen && (
        <div className="fixed inset-0 z-50 flex items-end sm:items-center justify-center bg-black/50 backdrop-blur-sm" onClick={() => setModalOpen(false)}>
          <div
            className="bg-card w-full max-w-lg rounded-t-3xl sm:rounded-3xl shadow-2xl p-6 max-h-[90vh] overflow-y-auto"
            onClick={(e) => e.stopPropagation()}
          >
            <h2 className="text-lg font-bold text-foreground mb-5">{editPlot ? "Edit Plot" : "Add New Plot"}</h2>
            <div className="flex flex-col gap-4">
              <FormField label="Plot Name">
                <input
                  type="text"
                  value={form.name}
                  onChange={(e) => setForm({ ...form, name: e.target.value })}
                  placeholder="e.g. Plot A — Hillside"
                  className="farm-input"
                />
              </FormField>

              <div className="grid grid-cols-2 gap-3">
                <FormField label="Size (m²)">
                  <input
                    type="number"
                    inputMode="numeric"
                    min={0}
                    value={form.sizeMeter || ""}
                    onChange={(e) => setForm({ ...form, sizeMeter: parseInt(e.target.value) || 0 })}
                    placeholder="e.g. 2400"
                    className="farm-input"
                  />
                </FormField>
                <FormField label="No. of Trees">
                  <input
                    type="number"
                    inputMode="numeric"
                    min={0}
                    value={form.trees || ""}
                    onChange={(e) => setForm({ ...form, trees: parseInt(e.target.value) || 0 })}
                    placeholder="e.g. 48"
                    className="farm-input"
                  />
                </FormField>
              </div>

              <FormField label="Notes">
                <textarea
                  value={form.notes}
                  onChange={(e) => setForm({ ...form, notes: e.target.value })}
                  placeholder="Additional notes about this plot..."
                  rows={3}
                  className={cn("farm-input resize-none")}
                />
              </FormField>
            </div>

            <div className="flex gap-3 mt-6">
              <button onClick={() => setModalOpen(false)} className="flex-1 py-4 rounded-xl border border-border text-foreground font-semibold text-sm active:scale-95 transition-transform">
                Cancel
              </button>
              <button
                onClick={handleSave}
                disabled={!form.name}
                className="flex-1 py-4 rounded-xl bg-primary text-primary-foreground font-semibold text-sm shadow-sm active:scale-95 transition-transform disabled:opacity-40"
              >
                {editPlot ? "Save Changes" : "Add Plot"}
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Delete Confirm */}
      {deleteId && (
        <div className="fixed inset-0 z-50 flex items-end sm:items-center justify-center bg-black/50 backdrop-blur-sm">
          <div className="bg-card w-full max-w-sm rounded-t-3xl sm:rounded-3xl shadow-2xl p-6">
            <h2 className="text-lg font-bold text-foreground mb-2">Delete Plot?</h2>
            <p className="text-sm text-muted-foreground mb-6">All data linked to this plot will remain. This cannot be undone.</p>
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
