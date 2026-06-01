import { Routes } from '@angular/router';
import { ApplicationsComponent } from './applications.component';

export default [
    {
        path: '',
        component: ApplicationsComponent,
    },
    {
        path: 'review',
        loadComponent: () =>
            import('./review/staff-review-workspace.component').then(
                (c) => c.StaffReviewWorkspaceComponent
            ),
    },
    {
        path: ':id',
        loadComponent: () =>
            import('./detail/application-detail.component').then(
                (c) => c.ApplicationDetailComponent
            ),
    },
] as Routes;
