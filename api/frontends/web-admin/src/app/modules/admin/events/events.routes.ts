import { Routes } from '@angular/router';
import { EventsComponent } from './events.component';

export default [
    {
        path: '',
        component: EventsComponent,
        children: [
            {
                path: '',
                pathMatch: 'full',
                redirectTo: 'table'
            },
            {
                path: 'table',
                loadComponent: () => import('./views/events-table/events-table.component').then(m => m.EventsTableComponent)
            },
            {
                path: 'calendar',
                loadComponent: () => import('./views/events-calendar/events-calendar.component').then(m => m.EventsCalendarComponent)
            },
            {
                path: 'kanban',
                loadComponent: () => import('./views/events-kanban/events-kanban.component').then(m => m.EventsKanbanComponent)
            }
        ]
    },
    {
        path: 'new',
        loadComponent: () => import('./form/event-form.component').then(m => m.EventFormComponent)
    },
    {
        path: ':id',
        loadComponent: () => import('./detail/event-detail.component').then(m => m.EventDetailComponent)
    },
    {
        path: ':id/edit',
        loadComponent: () => import('./form/event-form.component').then(m => m.EventFormComponent)
    },
    {
        path: ':id/checkin',
        loadComponent: () => import('./checkin/event-checkin.component').then(m => m.EventCheckinComponent)
    }
] as Routes;
