"use client"

import { useState } from "react"
import { Leaf, Eye, EyeOff } from "lucide-react"
import Link from "next/link"

export default function LoginPage() {
  const [showPass, setShowPass] = useState(false)
  const [form, setForm] = useState({ username: "", password: "" })
  const [loading, setLoading] = useState(false)

  async function handleLogin(e: React.FormEvent) {
    e.preventDefault()
    setLoading(true)
    await new Promise((r) => setTimeout(r, 800))
    setLoading(false)
    window.location.href = "/"
  }

  return (
    <div className="min-h-screen bg-background flex flex-col items-center justify-center px-5 py-12">
      {/* Logo */}
      <div className="flex flex-col items-center gap-3 mb-10">
        <div className="w-16 h-16 rounded-2xl bg-primary flex items-center justify-center shadow-md">
          <Leaf className="w-9 h-9 text-primary-foreground" />
        </div>
        <div className="text-center">
          <h1 className="text-2xl font-bold text-foreground">Durian Farm</h1>
          <p className="text-sm text-muted-foreground">Manager</p>
        </div>
      </div>

      <div className="w-full max-w-sm bg-card rounded-3xl border border-border shadow-sm p-6">
        <h2 className="text-lg font-bold text-foreground mb-1">Welcome back</h2>
        <p className="text-sm text-muted-foreground mb-6">Sign in to manage your farm</p>

        <form onSubmit={handleLogin} className="flex flex-col gap-4">
          <div className="flex flex-col gap-1.5">
            <label className="text-xs font-semibold text-muted-foreground uppercase tracking-wide">
              Username
            </label>
            <input
              type="text"
              autoComplete="username"
              value={form.username}
              onChange={(e) => setForm({ ...form, username: e.target.value })}
              placeholder="e.g. farm_manager"
              className="farm-input"
            />
          </div>

          <div className="flex flex-col gap-1.5">
            <label className="text-xs font-semibold text-muted-foreground uppercase tracking-wide">
              Password
            </label>
            <div className="relative">
              <input
                type={showPass ? "text" : "password"}
                autoComplete="current-password"
                value={form.password}
                onChange={(e) => setForm({ ...form, password: e.target.value })}
                placeholder="••••••••"
                className="farm-input pr-12"
              />
              <button
                type="button"
                onClick={() => setShowPass((v) => !v)}
                className="absolute right-3 top-1/2 -translate-y-1/2 w-9 h-9 flex items-center justify-center text-muted-foreground"
                aria-label={showPass ? "Hide password" : "Show password"}
              >
                {showPass ? <EyeOff className="w-5 h-5" /> : <Eye className="w-5 h-5" />}
              </button>
            </div>
          </div>

          <button
            type="submit"
            disabled={loading || !form.username || !form.password}
            className="w-full py-4 mt-2 rounded-xl bg-primary text-primary-foreground font-semibold text-base shadow-sm active:scale-95 transition-all disabled:opacity-50 flex items-center justify-center gap-2"
          >
            {loading ? (
              <span className="w-5 h-5 border-2 border-primary-foreground/30 border-t-primary-foreground rounded-full animate-spin" />
            ) : (
              "Sign In"
            )}
          </button>
        </form>

        <p className="text-center text-xs text-muted-foreground mt-5">
          Demo credentials: any username + password
        </p>
      </div>

      <Link href="/" className="mt-6 text-sm text-primary font-medium underline-offset-4 hover:underline">
        Skip to dashboard →
      </Link>
    </div>
  )
}
