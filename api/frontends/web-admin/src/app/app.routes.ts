import { Route } from '@angular/router';
import { initialDataResolver } from 'app/app.resolvers';
import { AuthGuard } from 'app/core/auth/guards/auth.guard';
import { NoAuthGuard } from 'app/core/auth/guards/noAuth.guard';
import { LayoutComponent } from 'app/layout/layout.component';

// @formatter:off
/* eslint-disable max-len */
/* eslint-disable @typescript-eslint/explicit-function-return-type */
export const appRoutes: Route[] = [
    // Default landing → Dashboard
    { path: '', pathMatch: 'full', redirectTo: 'dashboard' },

    // Post sign-in redirect → Dashboard
    { path: 'signed-in-redirect', pathMatch: 'full', redirectTo: 'dashboard' },

    // Auth routes for guests
    {
        path: '',
        canActivate: [NoAuthGuard],
        canActivateChild: [NoAuthGuard],
        component: LayoutComponent,
        data: {
            layout: 'empty',
        },
        children: [
            {
                path: 'confirmation-required',
                loadChildren: () =>
                    import(
                        'app/modules/auth/confirmation-required/confirmation-required.routes'
                    ),
            },
            {
                path: 'forgot-password',
                loadChildren: () =>
                    import(
                        'app/modules/auth/forgot-password/forgot-password.routes'
                    ),
            },
            {
                path: 'reset-password',
                loadChildren: () =>
                    import(
                        'app/modules/auth/reset-password/reset-password.routes'
                    ),
            },
            {
                path: 'sign-in',
                loadChildren: () =>
                    import('app/modules/auth/sign-in/sign-in.routes'),
            },
            {
                path: 'sign-up',
                loadChildren: () =>
                    import('app/modules/auth/sign-up/sign-up.routes'),
            },
        ],
    },

    // Auth routes for authenticated users
    {
        path: '',
        canActivate: [AuthGuard],
        canActivateChild: [AuthGuard],
        component: LayoutComponent,
        data: {
            layout: 'empty',
        },
        children: [
            {
                path: 'sign-out',
                loadChildren: () =>
                    import('app/modules/auth/sign-out/sign-out.routes'),
            },
            {
                path: 'unlock-session',
                loadChildren: () =>
                    import(
                        'app/modules/auth/unlock-session/unlock-session.routes'
                    ),
            },
        ],
    },

    // Public landing routes (no auth)
    {
        path: '',
        component: LayoutComponent,
        data: {
            layout: 'empty',
        },
        children: [
            {
                path: 'home',
                loadChildren: () =>
                    import('app/modules/landing/home/home.routes'),
            },
            {
                path: 'inquiry',
                loadChildren: () =>
                    import('app/modules/landing/inquiry/inquiry.routes'),
            },
            {
                path: 'portal',
                loadChildren: () =>
                    import('app/modules/landing/portal/portal.routes'),
            },
        ],
    },

    // Admin routes (authenticated, with sidebar layout)
    {
        path: '',
        canActivate: [AuthGuard],
        canActivateChild: [AuthGuard],
        component: LayoutComponent,
        resolve: {
            initialData: initialDataResolver,
        },
        children: [
            // Workspace
            {
                path: 'dashboard',
                loadChildren: () =>
                    import('app/modules/admin/dashboard/dashboard.routes'),
            },
            {
                path: 'constituents',
                loadChildren: () =>
                    import(
                        'app/modules/admin/constituents/constituents.routes'
                    ),
            },
            {
                path: 'duplicates',
                loadChildren: () =>
                    import('app/modules/admin/duplicates/duplicates.routes'),
            },
            {
                path: 'inquiries',
                loadChildren: () =>
                    import('app/modules/admin/inquiries/inquiries.routes'),
            },

            // Admissions
            {
                path: 'applications',
                loadChildren: () =>
                    import(
                        'app/modules/admin/applications/applications.routes'
                    ),
            },
            {
                path: 'leads',
                loadChildren: () =>
                    import('app/modules/admin/leads/admissions-leads.routes'),
            },

            // Engagement
            {
                path: 'communications',
                loadChildren: () =>
                    import(
                        'app/modules/admin/communications/communications.routes'
                    ),
            },
            {
                path: 'campaigns',
                loadChildren: () =>
                    import('app/modules/admin/campaigns/campaigns.routes'),
            },
            {
                path: 'events',
                loadChildren: () =>
                    import('app/modules/admin/events/events.routes'),
            },

            // Insights
            {
                path: 'reports',
                loadChildren: () =>
                    import('app/modules/admin/reports/reports.routes'),
            },

            // Admin
            {
                path: 'users',
                loadChildren: () =>
                    import('app/modules/admin/users/admin-users.routes'),
            },
            {
                path: 'audit',
                loadChildren: () =>
                    import('app/modules/admin/audit/audit.routes'),
            },
            {
                path: 'settings',
                loadChildren: () =>
                    import('app/modules/admin/settings-global/settings.routes'),
            },
            {
                path: 'admissions-settings',
                loadChildren: () =>
                    import(
                        'app/modules/admin/settings/admissions-settings.routes'
                    ),
            },

            {
                path: '404-not-found',
                pathMatch: 'full',
                loadChildren: () =>
                    import(
                        'app/modules/admin/pages/error/error-404/error-404.routes'
                    ),
            },
            { path: '**', redirectTo: '404-not-found' },
        ],
    },
];
