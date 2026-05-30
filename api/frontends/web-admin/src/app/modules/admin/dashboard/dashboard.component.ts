import { ChangeDetectionStrategy, Component } from '@angular/core';
import { RouterLink } from '@angular/router';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';

interface KpiCard {
    label: string;
    value: string;
    delta: string;
    trend: 'up' | 'down';
}

interface RecentApplication {
    id: string;
    applicant: string;
    program: string;
    status: string;
    submitted: string;
}

interface UpcomingEvent {
    title: string;
    date: string;
    location: string;
}

@Component({
    selector: 'app-dashboard',
    standalone: true,
    imports: [RouterLink, MatButtonModule, MatIconModule],
    changeDetection: ChangeDetectionStrategy.OnPush,
    templateUrl: './dashboard.component.html',
})
export class DashboardComponent {
    readonly userName = 'Avery';

    readonly kpis: readonly KpiCard[] = [
        { label: 'Inquiries', value: '1,284', delta: '+12.4%', trend: 'up' },
        { label: 'Active applications', value: '486', delta: '+3.8%', trend: 'up' },
        { label: 'Admitted', value: '152', delta: '-1.2%', trend: 'down' },
        { label: 'Enrolled', value: '98', delta: '+5.6%', trend: 'up' },
    ];

    readonly recentApplications: readonly RecentApplication[] = [
        { id: 'APP-3024', applicant: 'Sofia Martinez', program: 'B.Sc. Computer Science', status: 'In review', submitted: '2h ago' },
        { id: 'APP-3023', applicant: 'James Okoro', program: 'M.A. Education', status: 'Documents pending', submitted: '5h ago' },
        { id: 'APP-3022', applicant: 'Priya Patel', program: 'B.A. Economics', status: 'Submitted', submitted: 'Yesterday' },
        { id: 'APP-3021', applicant: 'Liam Chen', program: 'B.Sc. Biology', status: 'Interview', submitted: 'Yesterday' },
        { id: 'APP-3020', applicant: 'Aisha Bello', program: 'M.B.A.', status: 'Decision pending', submitted: '2d ago' },
    ];

    readonly upcomingEvents: readonly UpcomingEvent[] = [
        { title: 'Fall Open House', date: 'Sat, Jun 6 · 10:00 AM', location: 'Main Quad' },
        { title: 'Virtual Info Session — Engineering', date: 'Tue, Jun 9 · 6:00 PM', location: 'Online' },
        { title: 'Counselor Meetup', date: 'Thu, Jun 11 · 2:00 PM', location: 'Admissions Hall' },
    ];
}
