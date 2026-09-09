import {
    ChangeDetectionStrategy,
    Component,
    computed,
    inject,
    OnInit,
    signal,
} from '@angular/core';
import { SlicePipe } from '@angular/common';
import { RouterLink } from '@angular/router';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { AdmissionsService } from 'app/core/admissions/admissions.service';
import {
    Application,
} from 'app/core/admissions/admissions.types';
import { forkJoin } from 'rxjs';

interface KpiCard {
    label: string;
    value: number;
}

interface UpcomingEvent {
    title: string;
    date: string;
    location: string;
}

@Component({
    selector: 'app-dashboard',
    standalone: true,
    imports: [
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

    readonly userName = 'Avery';
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

    readonly upcomingEvents: readonly UpcomingEvent[] = [
        { title: 'Fall Open House', date: 'Sat, Jun 6 · 10:00 AM', location: 'Main Quad' },
        { title: 'Virtual Info Session — Engineering', date: 'Tue, Jun 9 · 6:00 PM', location: 'Online' },
        { title: 'Counselor Meetup', date: 'Thu, Jun 11 · 2:00 PM', location: 'Admissions Hall' },
    ];

    ngOnInit(): void {
        forkJoin({
            all: this.admissionsService.queryApplications({ rows: 1 }),
            submitted: this.admissionsService.queryApplications({ rows: 1, status: 'SUBMITTED' }),
            admitted: this.admissionsService.queryApplications({ rows: 1, status: 'ADMITTED' }),
            enrolled: this.admissionsService.queryApplications({ rows: 1, status: 'ENROLLED' }),
            recent: this.admissionsService.queryApplications({ rows: 5, orderBy: 'date_created,DESC' }),
        }).subscribe({
            next: (results) => {
                this.kpiCounts.set({
                    ACTIVE: results.all.total,
                    SUBMITTED: results.submitted.total,
                    ADMITTED: results.admitted.total,
                    ENROLLED: results.enrolled.total,
                });
                this.recentApplications.set(results.recent.items);
                this.loading.set(false);
            },
            error: () => {
                this.error.set('Unable to load dashboard data. Please try again.');
                this.loading.set(false);
            },
        });
    }
}
