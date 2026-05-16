import { Route } from '@angular/router';
import { SmartDemoPageComponent } from '@shared/components/smart-demo-page/smart-demo-page.component';

export const ADMIN_ROUTE: Route[] = [
  {
    path: 'dashboard',
    loadChildren: () =>
      import('./dashboard/dashboard.routes').then((m) => m.DASHBOARD_ROUTE),
  },
  {
    path: '**',
    component: SmartDemoPageComponent,
    data: { section: 'Admin', kind: 'workspace', icon: 'school' },
  },
];
