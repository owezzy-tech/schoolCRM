import { Routes } from '@angular/router';
import { PortalComponent } from './portal.component';
import { PortalLayoutComponent } from './components/portal-layout/portal-layout.component';

export default [
    {
        path: '',
        component: PortalLayoutComponent,
        children: [
            { path: '', component: PortalComponent },
            { path: 'apply', loadComponent: () => import('./apply/portal-apply.component').then(m => m.PortalApplyComponent) },
            { path: 'status', loadComponent: () => import('./status/portal-status.component').then(m => m.PortalStatusComponent) },
            { path: 'events', loadComponent: () => import('./events/portal-events.component').then(m => m.PortalEventsComponent) },
            { path: 'inquiry', loadComponent: () => import('./inquiry/portal-inquiry.component').then(m => m.PortalInquiryComponent) },
        ]
    }
] as Routes;
