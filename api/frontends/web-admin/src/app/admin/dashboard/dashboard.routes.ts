import { Route } from '@angular/router';
import { MainComponent } from './main/main.component';
import { Dashboard2Component } from './dashboard2/dashboard2.component';
import { DashboardComponent as StudentDashboard } from 'app/student/dashboard/dashboard.component';
import { DashboardComponent } from 'app/teacher/dashboard/dashboard.component';
import { SmartDemoPageComponent } from '@shared/components/smart-demo-page/smart-demo-page.component';
import { Page404Component } from 'app/authentication/page404/page404.component';
export const DASHBOARD_ROUTE: Route[] = [
  {
    path: '',
    redirectTo: 'main',
    pathMatch: 'full',
  },
  {
    path: 'main',
    component: MainComponent,
  },
  {
    path: 'dashboard2',
    component: Dashboard2Component,
  },
  {
    path: 'teacher-dashboard',
    component: DashboardComponent,
  },
  {
    path: 'student-dashboard',
    component: StudentDashboard,
  },
  {
    path: 'library-dashboard',
    component: SmartDemoPageComponent,
    data: { title: 'Library Dashboard', section: 'Admin Dashboard', kind: 'dashboard', icon: 'local_library' },
  },
  {
    path: 'transport-dashboard',
    component: SmartDemoPageComponent,
    data: { title: 'Transport Dashboard', section: 'Admin Dashboard', kind: 'dashboard', icon: 'directions_bus' },
  },
  { path: '**', component: Page404Component },
];
