import { AsyncPipe, DatePipe, SlicePipe } from '@angular/common';
import {
    ChangeDetectionStrategy,
    Component,
    computed,
    inject,
    OnInit,
    signal,
} from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { RouterLink } from '@angular/router';
import { AdmissionsService } from 'app/core/admissions/admissions.service';
import {
    AdmissionsEvent,
    Application,
} from 'app/core/admissions/admissions.types';
import { UserService } from 'app/core/user/user.service';
import { forkJoin } from 'rxjs';

interface KpiCard {
    label: string;
    value: number;
}

@Component({
    selector: 'app-dashboard',
    standalone: true,
    imports: [
        AsyncPipe,
        DatePipe,
        SlicePipe,
        RouterLink,
        MatButtonModule,
        MatIconModule,
        MatProgressSpinnerModule,
    ],
    changeDetection: ChangeDetectionStrategy.OnPush,
    templateUrl: './dashboard.component.html',
})
export class DashboardComponent implements OnInit {
    private readonly admissionsService = inject(AdmissionsService);

    readonly user$ = inject(UserService).user$;
    readonly loading = signal(true);
    readonly error = signal<string | null>(null);

    readonly recentApplications = signal<Application[]>([]);
    readonly kpiCounts = signal<Record<string, number>>({});

    readonly kpis = computed<KpiCard[]>(() => {
        const counts = this.kpiCounts();
        return [
            { label: 'Active applications', value: counts['ACTIVE'] ?? 0 },
            { label: 'Submitted', value: counts['SUBMITTED'] ?? 0 },
            { label: 'Admitted', value: counts['ADMITTED'] ?? 0 },
            { label: 'Enrolled', value: counts['ENROLLED'] ?? 0 },
        ];
    });

    readonly upcomingEvents = signal<AdmissionsEvent[]>([]);

    ngOnInit(): void {
        this.loading.set(true);
        this.error.set(null);
        forkJoin({
            all: this.admissionsService.queryApplications({ rows: 1 }),
            submitted: this.admissionsService.queryApplications({
                rows: 1,
                status: 'SUBMITTED',
            }),
            admitted: this.admissionsService.queryApplications({
                rows: 1,
                status: 'ADMITTED',
            }),
            enrolled: this.admissionsService.queryApplications({
                rows: 1,
                status: 'ENROLLED',
            }),
            recent: this.admissionsService.queryApplications({
                rows: 5,
                orderBy: 'date_created,DESC',
            }),
            events: this.admissionsService.queryEvents({
                rows: 3,
                status: 'upcoming',
                orderBy: 'start_time,ASC',
            }),
        }).subscribe({
            next: (results) => {
                this.kpiCounts.set({
                    ACTIVE: results.all.total,
                    SUBMITTED: results.submitted.total,
                    ADMITTED: results.admitted.total,
                    ENROLLED: results.enrolled.total,
                });
                this.recentApplications.set(results.recent.items);
                this.upcomingEvents.set(results.events.items);
                this.loading.set(false);
            },
            error: () => {
                this.error.set(
                    'Unable to load dashboard data. Please try again.'
                );
                this.loading.set(false);
            },
        });
    }
}
