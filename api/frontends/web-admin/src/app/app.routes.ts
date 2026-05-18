import { Route } from '@angular/router';
import { MainLayoutComponent } from './layout/app-layout/main-layout/main-layout.component';
import { AuthGuard } from '@core/guard/auth.guard';
import { AuthLayoutComponent } from './layout/app-layout/auth-layout/auth-layout.component';
import { Page404Component } from './authentication/page404/page404.component';
import { Role } from '@core';
import { SmartDemoPageComponent } from '@shared/components/smart-demo-page/smart-demo-page.component';

export const APP_ROUTE: Route[] = [
  {
    path: '',
    component: MainLayoutComponent,
    canActivate: [AuthGuard],
    children: [
      { path: '', redirectTo: '/authentication/signin', pathMatch: 'full' },

      {
        path: 'admin',
        canActivate: [AuthGuard],
        data: {
          role: [Role.SuperAdmin, Role.SchoolAdmin],
        },
        loadChildren: () =>
          import('./admin/admin.routes').then((m) => m.ADMIN_ROUTE),
      },
      {
        path: 'teacher',
        canActivate: [AuthGuard],
        data: {
          role: [Role.Teacher],
        },
        loadChildren: () =>
          import('./teacher/teacher.routes').then((m) => m.TEACHER_ROUTE),
      },
      {
        path: 'student',
        canActivate: [AuthGuard],
        data: {
          role: [Role.Student],
        },
        loadChildren: () =>
          import('./student/student.routes').then((m) => m.STUDENT_ROUTE),
      },
      {
        path: 'extra-pages',
        loadChildren: () =>
          import('./extra-pages/extra-pages.routes').then(
            (m) => m.EXTRA_PAGES_ROUTE
          ),
      },
      {
        path: 'multilevel',
        loadChildren: () =>
          import('./multilevel/multilevel.routes').then(
            (m) => m.MULTILEVEL_ROUTE
          ),
      },
      {
        path: 'calendar',
        component: SmartDemoPageComponent,
        data: { title: 'Calendar', section: 'Apps', kind: 'calendar', icon: 'event_note' },
      },
      {
        path: 'task',
        component: SmartDemoPageComponent,
        data: { title: 'Tasks', section: 'Apps', kind: 'list', icon: 'fact_check' },
      },
      {
        path: 'contacts',
        component: SmartDemoPageComponent,
        data: { title: 'Contacts', section: 'Apps', kind: 'list', icon: 'contacts' },
      },
      ...['email', 'apps', 'widget', 'ui', 'forms', 'tables', 'charts', 'timeline', 'icons'].map(
        (path) => ({
          path,
          children: [
            {
              path: '**',
              component: SmartDemoPageComponent,
              data: {
                section: path.charAt(0).toUpperCase() + path.slice(1),
                kind: 'workspace',
                icon: 'dashboard_customize',
              },
            },
          ],
        })
      ),
    ],
  },
  {
    path: 'authentication',
    component: AuthLayoutComponent,
    loadChildren: () =>
      import('./authentication/auth.routes').then((m) => m.AUTH_ROUTE),
  },
  { path: '**', component: Page404Component },
];
