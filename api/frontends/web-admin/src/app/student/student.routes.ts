import { Route } from '@angular/router';
import { DashboardComponent } from './dashboard/dashboard.component';
import { SmartDemoPageComponent } from '@shared/components/smart-demo-page/smart-demo-page.component';

export const STUDENT_ROUTE: Route[] = [
  {
    path: 'dashboard',
    component: DashboardComponent,
  },
  {
    path: '**',
    component: SmartDemoPageComponent,
    data: { section: 'Student', kind: 'workspace', icon: 'school' },
  },
];
