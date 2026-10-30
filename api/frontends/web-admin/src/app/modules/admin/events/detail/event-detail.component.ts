import {
    ChangeDetectionStrategy,
    Component,
    OnInit,
    computed,
    inject,
    signal,
} from '@angular/core';
import { CommonModule, DatePipe } from '@angular/common';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatTabsModule } from '@angular/material/tabs';
import { MatMenuModule } from '@angular/material/menu';
import { AdmissionsService } from 'app/core/admissions/admissions.service';
import { jsonApiErrorMessage } from 'app/core/api/json-api';
import { catchError, of } from 'rxjs';
import { EventItem, EventRegistration } from '../models/event.types';

@Component({
    selector: 'app-event-detail',
    standalone: true,
    imports: [
        CommonModule,
        RouterLink,
        MatIconModule,
        MatButtonModule,
        MatProgressBarModule,
        MatTabsModule,
        MatMenuModule,
        DatePipe
    ],
    changeDetection: ChangeDetectionStrategy.OnPush,
    templateUrl: './event-detail.component.html'
})
export class EventDetailComponent implements OnInit {
    private readonly admissionsService = inject(AdmissionsService);
    private readonly route = inject(ActivatedRoute);

    readonly event = signal<EventItem | null>(null);
    readonly registrations = computed(() => this.event()?.registrations ?? []);
    readonly remainingCapacity = computed(() => {
        const event = this.event();

        if (!event) {
            return 0;
        }

        return Math.max(event.capacity - event.registeredCount, 0);
    });
    readonly capacityEnforced = computed(() => this.remainingCapacity() === 0);
    readonly loading = signal(false);

    errorMessage = '';

    ngOnInit(): void {
        const id = this.route.snapshot.paramMap.get('id');
        if (!id) {
            this.errorMessage = 'Event ID is missing.';

            return;
        }

        this.loadEvent(id);
    }

    loadEvent(eventID: string): void {
        this.loading.set(true);
        this.errorMessage = '';
        this.admissionsService
            .getEvent(eventID)
            .pipe(
                catchError((error) => {
                    this.loading.set(false);
                    this.errorMessage = jsonApiErrorMessage(
                        error,
                        'Unable to load event details.'
                    );

                    return of(undefined);
                })
            )
            .subscribe((event) => {
                this.loading.set(false);
                if (event) {
                    this.event.set(event);
                }
            });
    }

    getMatchStatusClasses(status: EventRegistration['matchStatus']): string {
        switch (status) {
            case 'matched':
                return 'bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-200';
            case 'new-prospect':
                return 'bg-blue-100 text-blue-700 dark:bg-blue-900 dark:text-blue-200';
            case 'needs-review':
                return 'bg-amber-100 text-amber-700 dark:bg-amber-900 dark:text-amber-200';
            default:
                return 'bg-gray-100 text-gray-700';
        }
    }

    getRegistrationStatusClasses(
        status: EventRegistration['status']
    ): string {
        switch (status) {
            case 'checked-in':
                return 'bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-200';
            case 'registered':
                return 'bg-blue-100 text-blue-700 dark:bg-blue-900 dark:text-blue-200';
            case 'cancelled':
                return 'bg-red-100 text-red-700 dark:bg-red-900 dark:text-red-200';
            default:
                return 'bg-gray-100 text-gray-700';
        }
    }

    getTypeColor(type: string): string {
        switch (type) {
            case 'open-day': return 'bg-purple-500';
            case 'webinar': return 'bg-teal-500';
            case 'info-session': return 'bg-green-500';
            case 'campus-tour': return 'bg-amber-500';
            case 'fair': return 'bg-red-500';
            default: return 'bg-gray-500';
        }
    }
    
    getStatusClasses(status: string): string {
        switch (status) {
            case 'upcoming': return 'bg-blue-100 text-blue-700 dark:bg-blue-900 dark:text-blue-200';
            case 'live': return 'bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-200';
            case 'completed': return 'bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-300';
            case 'draft': return 'bg-amber-100 text-amber-700 dark:bg-amber-900 dark:text-amber-200';
            case 'cancelled': return 'bg-red-100 text-red-700 dark:bg-red-900 dark:text-red-200';
            default: return 'bg-gray-100 text-gray-700';
        }
    }
    
    getCapacityPercentage(registered: number, capacity: number): number {
        if (capacity === 0) return 0;
        return Math.round((registered / capacity) * 100);
    }
}
