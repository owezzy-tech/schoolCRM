import {
    ChangeDetectionStrategy,
    Component,
    computed,
    inject,
    signal,
} from '@angular/core';
import { toSignal } from '@angular/core/rxjs-interop';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import {
    ActivatedRoute,
    NavigationEnd,
    Router,
    RouterLink,
    RouterOutlet,
} from '@angular/router';
import { AdmissionsService } from 'app/core/admissions/admissions.service';
import { AdmissionsEvent } from 'app/core/admissions/admissions.types';
import { jsonApiErrorMessage } from 'app/core/api/json-api';
import { catchError, of } from 'rxjs';
import { filter } from 'rxjs/operators';

@Component({
    selector: 'app-events',
    standalone: true,
    imports: [RouterOutlet, RouterLink, MatButtonModule, MatIconModule],
    changeDetection: ChangeDetectionStrategy.OnPush,
    templateUrl: './events.component.html',
})
export class EventsComponent {
    private readonly admissionsService = inject(AdmissionsService);
    private readonly router = inject(Router);
    private readonly route = inject(ActivatedRoute);

    readonly currentView = signal<'table' | 'calendar' | 'kanban'>('table');
    readonly loading = signal(false);
    readonly typeFilter = signal<string | null>(null);
    readonly eventsResult = toSignal(this.admissionsService.events$, {
        initialValue: { items: [], total: 0, page: 1, rowsPerPage: 50 },
    });
    readonly events = computed<AdmissionsEvent[]>(() => {
        const selectedType = this.typeFilter();

        return this.eventsResult().items.filter(
            (event) => !selectedType || event.type === selectedType
        );
    });
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

    errorMessage = '';

    constructor() {
        this.router.events
            .pipe(filter((event) => event instanceof NavigationEnd))
            .subscribe(() => {
                const url = this.router.url;
                if (url.includes('/calendar')) {
                    this.currentView.set('calendar');
                } else if (url.includes('/kanban')) {
                    this.currentView.set('kanban');
                } else {
                    this.currentView.set('table');
                }
            });

        this.loadEvents();
    }

    loadEvents(): void {
        this.loading.set(true);
        this.errorMessage = '';
        this.admissionsService
            .queryEvents({
                page: 1,
                rows: 50,
                orderBy: 'start_time,ASC',
            })
            .pipe(
                catchError((error) => {
                    this.loading.set(false);
                    this.errorMessage = jsonApiErrorMessage(
                        error,
                        'Unable to load events.'
                    );

                    return of(undefined);
                })
            )
            .subscribe(() => {
                this.loading.set(false);
            });
    }

    onViewChange(view: 'table' | 'calendar' | 'kanban'): void {
        this.currentView.set(view);
        this.router.navigate([view], { relativeTo: this.route });
    }

    toggleTypeFilter(type: string): void {
        this.typeFilter.update((current) => (current === type ? null : type));
    }
}
