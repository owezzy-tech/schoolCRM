import { ChangeDetectionStrategy, Component, inject, signal } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import {
    ActivatedRoute,
    NavigationEnd,
    Router,
    RouterLink,
    RouterOutlet,
} from '@angular/router';
import { filter } from 'rxjs/operators';

@Component({
    selector: 'app-events',
    standalone: true,
    imports: [
        RouterOutlet,
        RouterLink,
        MatButtonModule,
        MatIconModule,
    ],
    changeDetection: ChangeDetectionStrategy.OnPush,
    templateUrl: './events.component.html',
})
export class EventsComponent {
    private readonly router = inject(Router);
    private readonly route = inject(ActivatedRoute);

    readonly currentView = signal<'table' | 'calendar' | 'kanban'>('table');
    
    // Filter state
    readonly typeFilter = signal<string | null>(null);
    readonly eventTypeFilters = [
        {
            type: 'open-day',
            label: 'Open Day',
            dotClass: 'bg-purple-500',
        },
        {
            type: 'webinar',
            label: 'Webinar',
            dotClass: 'bg-teal-500',
        },
        {
            type: 'info-session',
            label: 'Info Session',
            dotClass: 'bg-green-500',
        },
        {
            type: 'campus-tour',
            label: 'Campus Tour',
            dotClass: 'bg-amber-500',
        },
        {
            type: 'fair',
            label: 'Fair',
            dotClass: 'bg-red-500',
        },
    ] as const;

    constructor() {
        this.router.events.pipe(
            filter(event => event instanceof NavigationEnd)
        ).subscribe(() => {
            const url = this.router.url;
            if (url.includes('/calendar')) {
                this.currentView.set('calendar');
            } else if (url.includes('/kanban')) {
                this.currentView.set('kanban');
            } else {
                this.currentView.set('table');
            }
        });
    }

    onViewChange(view: 'table' | 'calendar' | 'kanban'): void {
        this.currentView.set(view);
        this.router.navigate([view], { relativeTo: this.route });
    }

    toggleTypeFilter(type: string): void {
        this.typeFilter.update(current => current === type ? null : type);
    }
}
