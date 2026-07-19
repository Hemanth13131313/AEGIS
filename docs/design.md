# AEGIS — Design System

**Related docs:** `prd.md`, `architecture.md`, `rules.md`, `phases.md`, `memory.md`
**Scope:** the AEGIS SOC dashboard (React/TypeScript UI) — policy management, detection triage, trace explorer, compliance reporting, red-team results.

**[ASSUMPTION]** — no visual brand assets (logos, existing palettes, Figma files) were provided in source material. This design system is an industry-informed default for a security operations product, built to be internally consistent and easy to restyle later if real brand assets arrive.

---

## 1. Design Philosophy

AEGIS's UI is a **high-signal, low-noise instrument panel** for security analysts, not a marketing surface. Every design decision optimizes for:

1. **Fast triage** — an analyst should be able to tell severity and category at a glance, without reading.
2. **Evidence-first** — every claim the UI makes (a detection, a policy block) is one click from its underlying evidence (trace, confidence score, matched rule).
3. **Calm under alert fatigue** — high-severity states are visually distinct but the overall palette stays restrained so that "everything looks urgent" never becomes true.
4. **Operator trust** — nothing is hidden behind unnecessary animation or ambiguity; state changes are immediate and legible.

## 2. Brand Identity

- **Name presentation:** "AEGIS" in wordmark, capitalized, paired with a small shield-derived glyph (abstract, not a literal shield clipart) used as the app icon/favicon only — not as decoration throughout the UI.
- **Personality:** precise, technical, quietly confident. Think "control tower," not "cybersecurity stock photo."
- **Voice in UI copy:** direct and specific ("Blocked by policy `tool_allowlist` v12" not "Something was blocked for security reasons").

## 3. Visual Language

- Dark-mode-first (SOC tools are typically used in low-light NOC/SOC environments), with a fully-supported light mode.
- Flat design with restrained elevation (subtle shadows only for overlays/modals — no skeuomorphism, no heavy gradients).
- Data density is embraced, not avoided: tables and traces can be dense, but typographic hierarchy keeps them scannable.
- Severity is communicated through a consistent, limited color+icon vocabulary (never color alone — see §12 Accessibility).

## 4. Color Palette

**[ASSUMPTION]** exact hex values below are a reasonable default dark-first palette; swap freely if real brand colors are provided later — component code should reference the token names (§8), not hex values directly.

### Core neutrals (dark mode base)
| Token | Hex | Usage |
|---|---|---|
| `color-bg-canvas` | `#0B0E14` | App background |
| `color-bg-surface` | `#131826` | Cards, panels |
| `color-bg-surface-raised` | `#1B2233` | Modals, popovers |
| `color-border-subtle` | `#2A3245` | Dividers, table borders |
| `color-text-primary` | `#E6E9F0` | Primary text |
| `color-text-secondary` | `#9AA4B8` | Secondary/muted text |
| `color-text-disabled` | `#5B6478` | Disabled text |

### Brand / accent
| Token | Hex | Usage |
|---|---|---|
| `color-accent-primary` | `#4C8DFF` | Primary actions, links, focus rings |
| `color-accent-primary-hover` | `#6FA3FF` | Hover state |

### Severity system (color + icon, never color alone)
| Token | Hex | Meaning | Icon convention |
|---|---|---|---|
| `color-severity-critical` | `#FF4D5E` | Critical detection / active block | Filled octagon |
| `color-severity-high` | `#FF9640` | High-risk detection | Filled triangle |
| `color-severity-medium` | `#F2C94C` | Medium / needs-review | Filled diamond |
| `color-severity-low` | `#4C8DFF` | Low / informational | Filled circle |
| `color-severity-safe` | `#3ECF8E` | Allowed / passed | Checkmark circle |

### Light mode
Light mode inverts the neutral scale (`color-bg-canvas` → `#F7F8FA`, `color-text-primary` → `#12131A`, etc.) while **severity tokens stay fixed** across both modes for consistent meaning.

## 5. Typography

| Token | Font | Usage |
|---|---|---|
| `font-family-ui` | Inter (system-ui fallback) | All UI text |
| `font-family-mono` | JetBrains Mono (ui-monospace fallback) | Trace payloads, policy/Rego code, IDs, log lines |

| Token | Size / Line-height | Usage |
|---|---|---|
| `text-xs` | 12px / 16px | Metadata, timestamps, table captions |
| `text-sm` | 13px / 18px | Table body, secondary labels |
| `text-base` | 14px / 20px | Default body text |
| `text-md` | 16px / 24px | Section headers within a panel |
| `text-lg` | 20px / 28px | Page titles |
| `text-xl` | 24px / 32px | Dashboard hero numbers (e.g., KPI counters) |

Font weights: 400 (body), 500 (emphasis/labels), 600 (headers). Avoid weight 700+ except in KPI hero numbers.

## 6. Spacing System

4px base unit, exposed as tokens:

| Token | Value |
|---|---|
| `space-1` | 4px |
| `space-2` | 8px |
| `space-3` | 12px |
| `space-4` | 16px |
| `space-5` | 24px |
| `space-6` | 32px |
| `space-7` | 48px |
| `space-8` | 64px |

Component internal padding defaults to `space-3`–`space-4`; section/page gutters default to `space-5`–`space-6`.

## 7. Grid System

- 12-column responsive grid, max content width `1440px`, gutters `space-5`.
- Dashboard layout: fixed left navigation rail (240px) + main content area; main content area uses the 12-column grid internally for panel layout.
- Breakpoints: `sm` 640px, `md` 768px, `lg` 1024px, `xl` 1280px, `2xl` 1536px (Tailwind defaults — reused for consistency with `frontend-design` conventions).

## 8. Design Tokens (summary reference)

All tokens above (`color-*`, `text-*`, `space-*`, `font-family-*`) are defined as CSS custom properties (or a Tailwind theme extension) — components must reference tokens, never raw hex/px values, so a future re-brand only requires a token-file change.

```css
:root {
  --color-bg-canvas: #0B0E14;
  --color-bg-surface: #131826;
  --color-accent-primary: #4C8DFF;
  --color-severity-critical: #FF4D5E;
  --color-severity-high: #FF9640;
  --color-severity-medium: #F2C94C;
  --color-severity-low: #4C8DFF;
  --color-severity-safe: #3ECF8E;
  --font-family-ui: 'Inter', system-ui, sans-serif;
  --font-family-mono: 'JetBrains Mono', ui-monospace, monospace;
  --space-4: 16px;
}
```

## 9. Component Library

### Buttons
- **Primary** — filled `color-accent-primary`, white text, used for the single primary action per view (e.g., "Deploy Policy").
- **Secondary** — outlined, `color-border-subtle` border, `color-text-primary` text.
- **Destructive** — filled `color-severity-critical`, reserved for irreversible actions (e.g., "Delete Policy Version").
- **Icon buttons** — 32px square hit target minimum, always paired with a tooltip and `aria-label`.

### Inputs
- Text inputs: `color-bg-surface` background, `color-border-subtle` border, `color-accent-primary` focus ring (2px, visible, never `outline: none` without replacement).
- Rego/code inputs use the monospace font and a syntax-highlighted code editor (e.g., CodeMirror/Monaco) with lint markers surfaced inline for `opa test` failures.

### Cards
- Used for KPI summaries and grouped policy/detection info; `color-bg-surface`, 8px corner radius, `space-4` internal padding, subtle 1px border rather than heavy shadow.

### Tables
- Used for policy lists, detection feeds, red-team run history. Dense by default (row height ~36px), with a "comfortable" density toggle.
- Severity column always leftmost, using the color+icon severity system (§4).
- Sticky header on scroll; sortable columns indicated with a static sort-direction icon (no icon shown only on hover — must be visible for accessibility).

### Navigation
- Left rail: primary sections (Overview, Policies, Detections, Trace Explorer, Red-Team, Compliance, Settings), icon + label, collapsible to icon-only.
- Top bar: org/app/environment scope switcher (reflecting the policy hierarchy from `architecture.md §9`), search, user menu.

### Forms
- Policy-scope forms explicitly show the resolved hierarchy (org → app → model → environment) as breadcrumb-style scope selector before editing fields, so the analyst always knows what level they're editing.
- Inline validation on blur, not only on submit; error text in `color-severity-high` with an icon, never red-only.

### Modals
- Used sparingly — for confirmation of destructive actions and for the "simulate policy change" flow (`prd.md §14 US-09`).
- Max width 560px for confirmation modals; up to 960px for the policy-simulation modal (needs to show before/after traffic impact).

### Animations & Micro-interactions
- Durations: 120ms for hover/focus transitions, 200ms for panel expand/collapse, 250ms for modal enter/exit. No animation exceeds 300ms — this is an operational tool, not a marketing site.
- Severity badges use a one-time subtle pulse (not looping) when a new critical detection arrives in a live feed, to draw attention without becoming distracting over a full shift.
- Respect `prefers-reduced-motion`: all non-essential animation is disabled when set.

## 10. Dark Mode

Dark mode is the default and primary-designed experience (§3). Light mode is a first-class supported alternative using the same token structure (§4) with inverted neutrals and fixed severity colors. No component may hard-code a light- or dark-specific color outside the token system.

## 11. Responsive Design

- Primary target is desktop/large-monitor SOC use (≥1280px) — this is where the product will actually be used day-to-day.
- Tablet (≥768px) support: left rail collapses to icon-only by default; trace explorer switches from side-by-side to stacked panels.
- Mobile (<768px): read-only summary views only (KPIs, alert list) — editing policy or replaying full traces is explicitly out of scope for mobile per `prd.md §18 Out-of-Scope`.

## 12. Accessibility

- WCAG 2.1 AA baseline (see `rules.md §5` for the enforceable engineering rules this design system must satisfy).
- Severity is always color **+ icon + text label**, never color alone.
- Focus states are always visible (2px `color-accent-primary` ring); never removed for aesthetic reasons.
- All charts/visualizations have an accompanying data-table view toggle for screen-reader and low-vision users.

## 13. Empty States

- Empty states always explain **why** the view is empty and **what action resolves it** — never a bare "No data."
  - Example (Detections feed, new tenant): "No detections yet. AEGIS is monitoring traffic — once policies are applied and traffic flows through the gateway, detections will appear here. [Apply your first policy →]"
- Illustration: simple, single-color line icon consistent with §14 icon style — never a large decorative illustration that outweighs the actionable text.

## 14. Loading States

- Skeleton screens (matching the target layout's shape) for table/list views — never a full-page spinner once initial app shell has loaded.
- Streaming data (e.g., live detection feed) uses a subtle top-of-list "new items" indicator rather than auto-scrolling the analyst's view out from under them.
- Long-running operations (compliance report generation, red-team run) show determinate progress where possible; indeterminate spinner only when duration is genuinely unknown, always paired with descriptive text ("Generating evidence bundle — this can take a few minutes").

## 15. Error States

- Distinguish visually and textually between the three error categories defined in `rules.md §8`: policy block (informational, `color-severity-medium`), upstream error (`color-severity-high`, retry affordance), internal error (`color-severity-critical`, "contact support/check status" affordance).
- Inline form errors are field-adjacent, never only in a toast, so they persist while the user fixes the issue.
- A global toast system is used only for transient, non-blocking notifications (e.g., "Policy deployed successfully").

## 16. Illustration Style

- Minimal, single-color (using `color-text-secondary`) line-art icons/illustrations only, consistent stroke width (2px at 24px viewport). No gradients, no photographic imagery, no literal "hacker in a hoodie" security clichés.

## 17. Icon Style

- Icon set: a single consistent outline icon library (e.g., Lucide/Feather-style — 24px grid, 2px stroke) used everywhere; no mixing of filled and outline styles except for the severity-system glyphs (§4), which are intentionally filled to stand out from the rest of the outline-based UI.

## 18. Image Guidelines

- The product has no photographic imagery in-app (it's an operational tool, not marketing). Marketing/landing-page assets (if built later) are explicitly out of scope for this design system and should be treated as a separate brand exercise.

## 19. Overall UI Consistency Rules

- Every screen has exactly one primary action, styled as the Primary button.
- Severity vocabulary (§4, §9 tables) is used identically across every feature area — Detections, Red-Team results, and Compliance findings all reuse the same five-level severity scale, never a feature-specific variant.
- Any new component must be built from existing tokens (§8) and existing primitives (§9) before a new primitive is introduced — check this design system first, consistent with the "check skills/docs before building" discipline in the broader engineering process.
- Copy tone follows §2: specific, technical, no hype language ("AI-powered," "next-gen") in the product UI itself.
