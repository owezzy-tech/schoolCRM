# Lesson planning prototype — browser test notes

The prototype in this folder was tested in headless Chrome 156 with viewport
overrides, touch emulation, and Chrome's Slow 3G preset.

## Coverage

- Widths: 375 × 812, 768 × 1024, 1024 × 768, 1440 × 900.
- Touch emulation (five points) at 375, 768, and 1024; `(pointer: coarse)` and
  `(hover: none)` verified.
- Network: 400 ms latency, 500 kbps, cache disabled, cold load per view.
- States: overview, schools, school detail, lesson-plan list, plan detail,
  version view, review queue, open mobile nav, open avatar menu, dark mode.
- Roles: Super admin, School admin, Teacher, HOD.

## Results

- 40 views (10 states × 4 widths): no page-level overflow, no clipped rows,
  no console errors, no touch targets below 44 px at 375/768/1024.
- Slow 3G: 39 of 40 cold loads at 2.16–2.20 s (median 2.17 s) for the 110.7 KB
  build. The delivered file is 112.9 KB and makes no external requests.
- Desktop (1440) keeps Material-density 40 px controls and 43 px nav rows; the
  44 px floor applies to coarse pointers.

## Fixes in this revision

1. Page title no longer clipped in the 375 px top bar.
2. Overview metrics form a 2 × 2 grid instead of four stacked rows (559 → 303 px).
3. Live/Latest cards use a right chevron to read as links.
4. Touch controls meet 44 px (buttons, role select, avatar, icon buttons).
5. Role select is 32 px tall inside its 40 px wrapper on desktop.
6. Top-bar icon buttons keep 44 × 44 beside long titles.
7. Timeline rows no longer overflow; links are 28 px tall (44 px on touch).
8. Top bar stays on one line at 768 and no longer overflows the page.
9. Tables below 820 px become labelled cards instead of clipped columns.
10. Review metadata no longer strands a separator on a wrapped line.

Before/after captures and the full audit are kept alongside the design source.
