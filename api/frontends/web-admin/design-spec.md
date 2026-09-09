# Elimu International School — Design Spec

Design system and theme mapping for SchoolCRM. Live reference: `elimu-school-dashboard.html`.
Direction: **calm scholarly precision** — a refined, low-chroma institutional indigo on
neutral cool-slate surfaces, serif display for brand + page titles, sans/mono for data work.

## Tokens

Light is the default; dark is a separate, audited value set. One accent, used for the
active nav marker and the single primary action — never splashed.

| Token               | Light       | Dark        | Role                                 |
|---------------------|-------------|-------------|--------------------------------------|
| `--bg`              | `#eef1f5`   | `#12151b`   | page ground (cool, never warm/cream) |
| `--surface`         | `#ffffff`   | `#1a1e26`   | cards, topbar, sheets                |
| `--surface-2`       | `#f6f8fa`   | `#20252f`   | wells, tab rails                     |
| `--fg`              | `#1b1f27`   | `#e7eaf0`   | primary text (never #000 / #fff)     |
| `--muted`           | `#5f6773`   | `#9aa2af`   | secondary text, captions             |
| `--border`          | `#d9dee6`   | `#2b313c`   | hairlines, dividers                  |
| `--accent`          | `#2f5fa8`   | `#82aef0`   | brand indigo (refined from #2196F3)  |
| `--danger`/`--success`/`--warning` | `#b3261e`/`#256c49`/`#8a6000` | `#f2b8b5`/`#7fd0a9`/`#e8c777` | semantic status — never the brand accent |

Derived tones (never new hex): `--accent-soft` and `--fg-soft` via `color-mix(in oklch, …)`.

## Design system details

### Typography scale
- **Display** (brand moment + page titles): serif stack `var(--font-display)` — Iowan Old Style / Charter / Georgia. Page titles clamp 26–34px, weight 600, letter-spacing −0.01em.
- **Body**: system sans, 15–16px, line-height 1.6, line-length 65–75ch. Data-dense tables reuse the sans at 14px.
- **Mono / numerics**: `ui-monospace` stack with `font-variant-numeric: tabular-nums` for **every** number; eyebrow / caption labels 9–11px with ~0.18em tracking, uppercase.

### Spacing
- 8pt grid (4px base). Card padding 20–24px; section gutter 24px; inter-component gap 12–20px; nav item height 40px. No arbitrary values.

### Radii & elevation
- `--radius: 10px`, `--radius-lg: 14px`. Cards: `--surface` background, 1px `--border`, soft `--shadow`. No glow, no gradient washes, no layered stacking.

### Iconography
- One monoline family, inline SVG, consistent ~2px stroke with 2px radius; 18px in nav, 16px in actions. Never mix filled and outlined icons at the same level, and never emoji as a functional icon.

### Component notes
- **Button**: primary = `--accent` fill + on-accent label; soft = `--accent-soft` bg; ghost = border-only. `:active { transform: scale(0.97) }` over 160ms. One primary action per view.
- **Menu / popover**: `--surface`, 1px `--border`, `--radius`, origin-aware enter, 140–200ms.
- **Table**: mono 11px headers, row hover under `(hover:hover)`, numeric columns tabular and right-aligned.
- **Form**: persistent label (never a placeholder), validate on blur, `--radius` inputs, inline success/error.
- **Status pill**: dot + text label, never color-only — success / warning / danger.

## Mapping to the Angular app

These values are the source of truth for the Fuse theme + Material integration:

- **Fuse theme** (`tailwind.config.js`, `src/app/app.config.ts`): set the default theme
  palette to the refined indigo instead of stock indigo/#2196F3. In
  `tailwind.config.js` put the brand palette in the layout/theme `palettes` block and
  make it the default theme's `primary`. Keep `darkMode: ['selector','.dark']` and apply
  `.dark` on `<body>`; a Fuse theme is a single palette applied in two modes, so your
  light/dark difference here is the accent light/dark value pair above.
- **Material integration** (`src/@fuse/styles/overrides/angular-material.scss`): this repo
  already carries the native-control fix (`opacity: 0 !important`) for checkbox, radio,
  slide-toggle and slider (see the `@ Selection Controls` block) plus the Material v21 token
  bridge, so the `--mdc-*` surfaces resolve to the Fuse theme tokens. No further change is
  needed there; the brand color now flows through the default theme primary wired in
  `tailwind.config.js`.
- **Theme tokens**: set the Material + Fuse theme via these same CSS custom properties so
  light and dark both flow from one place. Put the custom props on `:root`/`.dark` and
  reference them in the Material `--mdc-theme-*` and `--mat-*` variables.

## Motion (applied)

- Only `transform` + `opacity`. Press = `scale(0.97)` with `cubic-bezier(.23,1,.32,1)`, 160ms.
- Dropdowns/popovers: 140–200ms, origin-aware `transform-origin`, never `scale(0)`.
- View switch: fade + 8px translate, 240ms, `cubic-bezier(.32,.72,0,1)`. Modal: 200–240ms.
- Hover only under `@media (hover:hover)`; `prefers-reduced-motion` disables all tr/anim.
- Theme cross-fade: a 200ms `background-color`/`color`/`border-color` transition only.

## Accessibility

- Body text contrast ≥4.5:1 in both modes; status is never color-only — every pill pairs
  a dot with a text label (`.pill`).
- Visible `:focus-visible` outline (2px accent, offset 2px) on every interactive element.
- Controls are keyboard-reachable; `Esc` closes menus/modal; modal focuses first field and
  restores focus on close. Avatars double as an initials fallback (`.avatar-initials`).

## Applied in code — the Fuse theme wiring

These are the exact changes wired into `api/frontends/web-admin` so the design system ships
in the real app. They are already applied; this is the source of truth.

### `tailwind.config.js`

Swap the stock template palette for the refined institutional indigo and make it the default
Fuse theme's `primary`:

```js
const customPalettes = {
    brand: generatePalette('#2f5fa8'),
};

const themes = {
    default: {
        primary: {
            ...customPalettes.brand,
            DEFAULT: customPalettes.brand[600],
        },
        // accent/warn left unchanged ...
    },
    // ...
};
```

`generatePalette('#2f5fa8')` derives the 50–900 ramp at build time. `app.config.ts` already
applies `theme: 'theme-default'`, so no change is needed there.

### `src/@fuse/styles/main.scss`

Add a scheme-aware dark-mode primary. The light-mode brand indigo (`#2f5fa8`) reads at ~2.3:1
on dark surfaces, so dark mode remaps the Fuse primary tokens to a lighter, readable indigo:

```scss
body.dark {
    --fuse-primary: #82aef0;
    --fuse-primary-rgb: 130, 174, 240;
    --fuse-primary-50: #f0f5fd;
    --fuse-primary-100: #dfeafb;
    --fuse-primary-200: #cbddf9;
    --fuse-primary-300: #b4cef6;
    --fuse-primary-400: #9ec0f3;
    --fuse-primary-500: #82aef0;
    --fuse-primary-600: #6a8fc6;
    --fuse-primary-700: #526f9c;
    --fuse-primary-800: #384d6e;
    --fuse-primary-900: #1c293e;
    --fuse-on-primary: #0c1420;
    --fuse-on-primary-rgb: 12, 20, 32;
}
```

`body.dark` (specificity 0,1,1) beats the generated `body, .theme-default` theme vars
(0,1,0), so the override wins regardless of source order. Material inherits it automatically
through the existing `--mat-*` → `--fuse-primary` bridge, so no Material edit is needed.

## Logo system

The Elimu mark is a rounded indigo tile with a small spark (top right) and three bars
ascending left to right — growth, knowledge. The wordmark sets **Elimu** in the serif
display face with a tracked "INTERNATIONAL SCHOOL" caption in the mono face. Brand colors
are baked in (no CSS variables) so the files are portable. All SVGs scale losslessly to
any width; raster favicon sizes are provided for drop-in.

| File | Use |
|------|-----|
| `elimu-logo-primary.svg` | lockup on light surfaces |
| `elimu-logo-reversed.svg` | lockup on dark surfaces (white wordmark + hairline) |
| `elimu-mark.svg` | symbol alone on light |
| `elimu-mark-reversed.svg` | symbol alone on dark |
| `elimu-logo-primary.png` | raster lockup, 960×192 (3×) |
| `elimu-logo-reversed.png` | raster lockup, 960×192 (3×) |
| `elimu-mark.png` | 512×512 app icon / og-image |
| `elimu-mark-180.png` | 180×180 apple-touch-icon |
| `elimu-mark-32.png` | 32×32 browser favicon |

Usage rules: primary lockup on light, reversed on dark — never the primary version on a dark
surface or the reversed version on a light one. Keep clear space = the height of one bar
(~0.5× the tile) around the mark. Min lockup width 160px print / 96px screen; below that use
the mark alone. Never recolor, outline, rotate, or add a drop shadow to the mark.
