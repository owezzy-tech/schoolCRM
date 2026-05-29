# PROJECT KNOWLEDGE BASE

**Generated:** 2026-05-29
**Commit:** d7e9bab
**Branch:** feature/api-auth

## OVERVIEW

Angular 21 admin dashboard built on Fuse Template v21. Standalone components, Tailwind + SCSS + Angular Material, RxJS service-based state, REST API with mock layer for dev.

## STRUCTURE

```
web-admin-ng/
├── src/
│   ├── @fuse/              # Internal UI framework (components, services, animations, validators)
│   ├── app/
│   │   ├── core/           # Singleton services: auth, user, navigation, icons, transloco
│   │   ├── layout/         # Layout shell + common widgets (notifications, messages, shortcuts)
│   │   ├── modules/        # Feature routes: auth/, admin/, landing/
│   │   ├── mock-api/       # Dev-only mock data interceptors
│   │   ├── app.config.ts   # Provider configuration (DI root)
│   │   ├── app.routes.ts   # Top-level routing with lazy loading
│   │   └── app.resolvers.ts # Initial data resolver (forkJoin)
│   └── styles/             # Global SCSS + Tailwind entry
├── public/                 # Static assets (images, icons)
├── angular.json            # CLI config: project "fuse", SCSS, budgets
├── tailwind.config.js      # Custom breakpoints, themes, Fuse plugins
└── .prettierrc             # 4-space, single quotes, organize-imports + tailwindcss plugins
```

## WHERE TO LOOK

| Task | Location | Notes |
|------|----------|-------|
| Add feature module | `src/app/modules/` | Standalone component + lazy route in `app.routes.ts` |
| Auth flow | `src/app/core/auth/` | AuthService, guards (CanActivateFn), interceptor |
| Add HTTP interceptor | `src/app/core/auth/auth.provider.ts` | Chain via `withInterceptors()` |
| Layout/theme config | `src/app/app.config.ts` → `provideFuse()` | Scheme, theme, layout type |
| Custom validators | `src/@fuse/validators/` | `mustMatch`, `isEmptyInputValue` |
| Mock API data | `src/app/mock-api/` | Register in `mock-api/index.ts` |
| Global styles | `src/styles/` | Tailwind layers + Material theme overrides |
| Notifications/messages | `src/app/layout/common/` | Services use ReplaySubject |
| i18n translations | `src/assets/i18n/` | Transloco, languages: en, tr |

## CONVENTIONS

- **4-space indent** everywhere (Prettier + EditorConfig enforced)
- **Single quotes** for TS
- **Standalone components only** — no NgModules
- **UntypedFormBuilder** used for forms (legacy pattern from Fuse template)
- **ReplaySubject** for service state (not signals, not NgRx)
- **Prefix**: `app-` for components
- **Style**: SCSS scoped per component
- **Tailwind classes sorted** via prettier-plugin-tailwindcss
- **Import order** managed by prettier-plugin-organize-imports
- **Lazy loading**: All feature routes via `loadChildren`

## ANTI-PATTERNS (THIS PROJECT)

- No NgModules — don't create `.module.ts` files
- No NgRx — state lives in services with RxJS subjects
- No ESLint configured — rely on Prettier + TypeScript compiler
- No environment files — API URLs are relative (`api/...`)
- Don't use `@ts-ignore` or `as any`

## COMMANDS

```bash
npm start          # Dev server (:4200)
npm run build      # Production build → dist/fuse/
npm test           # Karma + Jasmine (no tests exist yet)
```

## NOTES

- **Monorepo context**: Backend uses gRPC (port 6001) but frontend uses REST. No ConnectRPC.
- **Mock API**: All `api/*` calls intercepted in dev by `src/app/mock-api/`. Disable by removing from `provideFuse()`.
- **Build budgets**: 3MB warn / 5MB error (initial), 75KB / 90KB (component styles)
- **Auth tokens**: Stored in localStorage, validated via `AuthUtils.isTokenExpired()`
- **Dark mode**: CSS class-based (`.dark` selector on body)
- **Tailwind breakpoints**: Non-standard — sm:600, md:960, lg:1280, xl:1440
