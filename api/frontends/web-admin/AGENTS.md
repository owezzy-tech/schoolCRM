# PROJECT KNOWLEDGE BASE

**Generated:** 2026-06-07
**Commit:** 271bc65
**Branch:** develop

## OVERVIEW

Angular 21 admin dashboard built on Fuse Template v21. Standalone components, Tailwind + SCSS + Angular Material, RxJS service-based state, JSON:API REST integration, and a dev mock layer for template/demo fixtures.

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
| Admissions REST client | `src/app/core/admissions/` | JSON:API unwrap helpers, ReplaySubject streams, query/create/update methods |
| Admissions admin screens | `src/app/modules/admin/` | Applications, events, campaigns, communications, records, audit, and settings use live `AdmissionsService` APIs |
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
npm test           # Angular test runner backed by Vitest specs
```

## NOTES

- **Monorepo context**: Backend uses gRPC (port 6001) but frontend uses REST. No ConnectRPC.
- **Mock API**: `src/app/mock-api/` remains for dev/template fixtures. Admissions admin workflows now call live `/v1/admissions/*` JSON:API endpoints through `AdmissionsService`; do not add new admissions admin behavior only to mocks.
- **Admissions settings**: Custom fields, application form templates, lead scoring rules, programs, academic terms, and import batches are loaded from live Admissions APIs. Assignment rules/candidates remain local preview data until backend endpoints exist.
- **Build budgets**: 3MB warn / 5MB error (initial), 75KB / 90KB (component styles)
- **Auth tokens**: Stored in localStorage, validated via `AuthUtils.isTokenExpired()`
- **API responses**: Backend REST APIs return JSON:API v1.1 envelopes. Read resource payloads from `data.attributes`, collection payloads from `data` plus `meta`, and error messages from `errors[0].detail`.
- **Angular Material native input fix**: In `src/@fuse/styles/overrides/angular-material.scss`, keep MDC native inputs visually hidden for selection controls. The checkbox fix (`.mdc-checkbox__native-control { opacity: 0 !important; }`) should also cover radio, slide-toggle/switch, and slider selectors (`.mdc-radio__native-control`, `.mdc-switch__native-control`, `.mdc-slider__input`, `input[type='range']`) when native inputs render visibly.
- **Dark mode**: CSS class-based (`.dark` selector on body)
- **Tailwind breakpoints**: Non-standard — sm:600, md:960, lg:1280, xl:1440

## BEADS WORKFLOW

- This frontend is part of the root SchoolCRM monorepo Beads database.
- Always run `bd` commands from `/Users/owen_adirah/GolandProjects/schoolCRM`.
- Do **not** run `bd init` in `api/frontends/web-admin` and do not create a nested `.beads/` directory here.
- Create or claim a root Beads issue before implementing frontend work:

```bash
bd ready
bd show <id>
bd update <id> --claim
```

- Use inline Beads commands only; never use `bd edit` because it opens an interactive editor.
- Close the issue only after frontend validation passes (`npx -p node@22 npm run build`, `npx -p node@22 npm test`, and targeted LSP diagnostics).
