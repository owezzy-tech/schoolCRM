# @fuse — Internal UI Framework

## OVERVIEW

Fuse admin template framework layer. Provides layout system, themed components, services, animations, directives, and Tailwind plugins.

## STRUCTURE

```
@fuse/
├── animations/         # Reusable Angular animations (fade, slide, shake, etc.)
├── components/         # Alert, Card, Drawer, Fullscreen, Highlight, LoadingBar, Masonry, Navigation
├── directives/         # Scrollbar (PerfectScrollbar wrapper)
├── lib/mock-api/       # MockApiInterceptor engine (intercepts HttpClient calls)
├── services/           # Config, Confirmation (Material Dialog), Loading, Media Watcher, Platform, Splash Screen, Utils
├── styles/             # SCSS: Material overrides, Tailwind base layers, theme variables
├── validators/         # Custom form validators: mustMatch, isEmptyInputValue
├── version/            # Template version tracking
└── fuse.provider.ts    # provideFuse() — root DI setup for all Fuse services
```

## WHERE TO LOOK

| Task | Location | Notes |
|------|----------|-------|
| Add animation | `animations/` | Export from module, use in component `@trigger` |
| Confirmation dialog | `services/confirmation/` | Injects MatDialog, returns observable |
| Loading state | `services/loading/` | Auto-tracked per URL via interceptor |
| Theme/scheme config | `services/config/` | FuseConfigService manages layout, scheme, theme |
| Responsive breakpoints | `services/media-watcher/` | Wraps BreakpointObserver with custom breakpoints |
| Navigation component | `components/navigation/` | Vertical + horizontal variants with appearances |
| Mock API handler | `lib/mock-api/` | Extend FuseMockApiHandler, register in app mock-api/index.ts |
| Tailwind plugins | `tailwind.config.js` (root) | icon-size, theming, custom utilities |

## CONVENTIONS

- **Do NOT modify @fuse/** unless extending the template framework itself
- Components are standalone with `changeDetection: OnPush`
- Services use `providedIn: 'root'` or are provided via `provideFuse()`
- Animations exported as const arrays (e.g., `fuseAnimations`)
- SCSS in `styles/` uses `@use` with Material theming API

## ANTI-PATTERNS

- Don't import @fuse internals directly — use the public barrel exports
- Don't duplicate Fuse service logic in app code (use DI)
- Don't override Fuse SCSS with `!important` — extend via Tailwind or theme variables
