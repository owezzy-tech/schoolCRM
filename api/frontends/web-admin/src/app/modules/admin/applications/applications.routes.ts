import { Routes } from '@angular/router';
import { ApplicationsComponent } from './applications.component';

export default [
    {
        path: '',
        component: ApplicationsComponent,
    },
    {
        path: ':id',
        loadComponent: () => import('./detail/application-detail.component').then(c => c.ApplicationDetailComponent),
    }
] as Routes;
