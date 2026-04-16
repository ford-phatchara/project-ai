import { cn } from "@/lib/utils"

export function GradeBadge({ grade }: { grade: "A" | "B" }) {
  return (
    <span
      className={cn(
        "inline-flex items-center justify-center rounded-lg px-2.5 py-1 text-xs font-bold min-w-[2rem]",
        grade === "A"
          ? "bg-primary/10 text-primary"
          : "bg-accent/20 text-accent-foreground"
      )}
    >
      {grade}
    </span>
  )
}
