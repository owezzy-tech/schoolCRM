import { ChangeDetectionStrategy, Component, inject, signal } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatChipsModule } from '@angular/material/chips';
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
        MatChipsModule,
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
