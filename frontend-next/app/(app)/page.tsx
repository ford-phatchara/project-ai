"use client"

import { DollarSign, Weight, TrendingUp, Leaf } from "lucide-react"
import { StatCard } from "@/components/stat-card"
import { PageHeader } from "@/components/page-header"
import { GradeBadge } from "@/components/grade-badge"
import { SALES, MONTHLY_REVENUE, getPlotName } from "@/lib/mock-data"
import {
  LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer,
  PieChart, Pie, Cell, Legend,
} from "recharts"

const totalRevenue = SALES.reduce((sum, s) => sum + s.total, 0)
const totalWeight = SALES.reduce((sum, s) => sum + s.weightKg, 0)
const avgPrice = totalRevenue / totalWeight

const gradeA = SALES.filter((s) => s.grade === "A").reduce((sum, s) => sum + s.total, 0)
const gradeB = SALES.filter((s) => s.grade === "B").reduce((sum, s) => sum + s.total, 0)
const PIE_DATA = [
  { name: "Grade A", value: gradeA },
  { name: "Grade B", value: gradeB },
]
const PIE_COLORS = ["#3d7a45", "#c9a227"]

const recentSales = [...SALES].sort((a, b) => b.date.localeCompare(a.date)).slice(0, 5)

function formatRM(val: number) {
  return `RM ${val.toLocaleString()}`
}

export default function DashboardPage() {
  return (
    <div>
      <PageHeader
        title="Dashboard"
        subtitle={`Season 2024 overview`}
      />

      <div className="px-4 lg:px-8 flex flex-col gap-4">
        {/* Stat Cards */}
        <div className="grid grid-cols-1 gap-3 sm:grid-cols-3">
          <StatCard
            label="Total Revenue"
            value={`RM ${(totalRevenue / 1000).toFixed(1)}k`}
            sub="2024 season"
            icon={<DollarSign className="w-4 h-4 text-primary" />}
            accent
          />
          <StatCard
            label="Total Weight Sold"
            value={`${totalWeight.toLocaleString()} kg`}
            sub="All grades"
            icon={<Weight className="w-4 h-4 text-muted-foreground" />}
          />
          <StatCard
            label="Avg. Price / kg"
            value={`RM ${avgPrice.toFixed(2)}`}
            sub="Blended average"
            icon={<TrendingUp className="w-4 h-4 text-muted-foreground" />}
          />
        </div>

        {/* Monthly Revenue Chart */}
        <div className="bg-card rounded-2xl border border-border shadow-sm p-4">
          <p className="text-sm font-semibold text-foreground mb-4">Monthly Revenue (RM)</p>
          <ResponsiveContainer width="100%" height={180}>
            <LineChart data={MONTHLY_REVENUE} margin={{ top: 4, right: 8, left: -20, bottom: 0 }}>
              <CartesianGrid strokeDasharray="3 3" stroke="var(--border)" />
              <XAxis dataKey="month" tick={{ fontSize: 11 }} stroke="var(--muted-foreground)" />
              <YAxis tick={{ fontSize: 11 }} stroke="var(--muted-foreground)" tickFormatter={(v) => `${v / 1000}k`} />
              <Tooltip
                formatter={(val: number) => [`RM ${val.toLocaleString()}`, "Revenue"]}
                contentStyle={{ borderRadius: "0.75rem", border: "1px solid var(--border)", background: "var(--card)", color: "var(--foreground)" }}
              />
              <Line type="monotone" dataKey="revenue" stroke="var(--color-primary)" strokeWidth={2.5} dot={{ r: 3, fill: "var(--color-primary)" }} />
            </LineChart>
          </ResponsiveContainer>
        </div>

        {/* Grade Breakdown */}
        <div className="bg-card rounded-2xl border border-border shadow-sm p-4">
          <p className="text-sm font-semibold text-foreground mb-2">Revenue by Grade</p>
          <div className="flex items-center justify-center">
            <ResponsiveContainer width="100%" height={180}>
              <PieChart>
                <Pie
                  data={PIE_DATA}
                  cx="50%"
                  cy="50%"
                  innerRadius={50}
                  outerRadius={75}
                  paddingAngle={4}
                  dataKey="value"
                >
                  {PIE_DATA.map((entry, index) => (
                    <Cell key={entry.name} fill={PIE_COLORS[index % PIE_COLORS.length]} />
                  ))}
                </Pie>
                <Legend iconType="circle" iconSize={10} formatter={(val) => <span className="text-xs text-foreground">{val}</span>} />
                <Tooltip
                  formatter={(val: number) => [`RM ${val.toLocaleString()}`, ""]}
                  contentStyle={{ borderRadius: "0.75rem", border: "1px solid var(--border)", background: "var(--card)", color: "var(--foreground)" }}
                />
              </PieChart>
            </ResponsiveContainer>
          </div>
        </div>

        {/* Recent Sales */}
        <div className="bg-card rounded-2xl border border-border shadow-sm p-4">
          <div className="flex items-center justify-between mb-3">
            <p className="text-sm font-semibold text-foreground">Recent Sales</p>
            <a href="/sales" className="text-xs text-primary font-medium">View all</a>
          </div>
          <div className="flex flex-col gap-2">
            {recentSales.map((sale) => (
              <div key={sale.id} className="flex items-center justify-between py-3 border-b border-border last:border-0">
                <div className="flex items-center gap-3">
                  <div className="w-10 h-10 rounded-xl bg-muted flex items-center justify-center shrink-0">
                    <Leaf className="w-5 h-5 text-primary" />
                  </div>
                  <div>
                    <p className="text-sm font-medium text-foreground">{getPlotName(sale.plot)}</p>
                    <p className="text-xs text-muted-foreground">{sale.date} · {sale.weightKg} kg</p>
                  </div>
                </div>
                <div className="flex items-center gap-2 shrink-0">
                  <GradeBadge grade={sale.grade} />
                  <p className="text-sm font-semibold text-foreground">{formatRM(sale.total)}</p>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  )
}
