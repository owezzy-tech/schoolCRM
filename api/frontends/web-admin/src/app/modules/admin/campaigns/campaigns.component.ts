import { ChangeDetectionStrategy, Component } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';

interface CampaignStats {
    sent: string;
    opened: string;
    clicked: string;
    replied: string;
}

interface Campaign {
    id: string;
    title: string;
    status: 'Active' | 'Draft' | 'Completed';
    channel: 'Email' | 'SMS';
    dateRange: string;
    stats: CampaignStats;
}

@Component({
    selector: 'app-campaigns',
    standalone: true,
    imports: [MatButtonModule, MatIconModule],
    changeDetection: ChangeDetectionStrategy.OnPush,
    templateUrl: './campaigns.component.html',
})
export class CampaignsComponent {
    readonly campaigns: readonly Campaign[] = [
        { id: '1', title: 'Fall 2026 Open House Invite', status: 'Active', channel: 'Email', dateRange: 'May 1 - Jun 6, 2026', stats: { sent: '2,450', opened: '45%', clicked: '12%', replied: '-' } },
        { id: '2', title: 'Missing Documents Reminder', status: 'Active', channel: 'SMS', dateRange: 'May 15 - Aug 1, 2026', stats: { sent: '342', opened: '-', clicked: '28%', replied: '8%' } },
        { id: '3', title: 'Financial Aid Deadline', status: 'Draft', channel: 'Email', dateRange: 'Jun 15 - Jul 1, 2026', stats: { sent: '-', opened: '-', clicked: '-', replied: '-' } },
        { id: '4', title: 'Spring 2026 Yield Campaign', status: 'Completed', channel: 'Email', dateRange: 'Mar 1 - May 1, 2026', stats: { sent: '1,200', opened: '62%', clicked: '18%', replied: '4%' } },
        { id: '5', title: 'Counselor Newsletter - May', status: 'Completed', channel: 'Email', dateRange: 'May 1 - May 2, 2026', stats: { sent: '85', opened: '74%', clicked: '35%', replied: '-' } },
        { id: '6', title: 'Interview Scheduling Ping', status: 'Active', channel: 'Email', dateRange: 'May 10 - Jun 30, 2026', stats: { sent: '156', opened: '82%', clicked: '45%', replied: '-' } },
    ];
}
