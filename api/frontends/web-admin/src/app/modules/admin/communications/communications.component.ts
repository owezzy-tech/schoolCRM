import { ChangeDetectionStrategy, Component, signal } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';

interface Message {
    id: string;
    recipient: string;
    initials: string;
    channelIcon: string;
    subject: string;
    timestamp: string;
    status: string;
}

@Component({
    selector: 'app-communications',
    standalone: true,
    imports: [MatButtonModule, MatIconModule],
    changeDetection: ChangeDetectionStrategy.OnPush,
    templateUrl: './communications.component.html',
})
export class CommunicationsComponent {
    readonly activeTab = signal<string>('All');
    readonly tabs = ['All', 'Email', 'SMS', 'Calls', 'Notes'];

    readonly messages: readonly Message[] = [
        { id: '1', recipient: 'Sofia Martinez', initials: 'SM', channelIcon: 'heroicons_outline:envelope', subject: 'Application Status Update', timestamp: '10:30 AM', status: 'Delivered' },
        { id: '2', recipient: 'James Okoro', initials: 'JO', channelIcon: 'heroicons_outline:chat-bubble-left', subject: 'Reminder: Missing Transcript', timestamp: 'Yesterday', status: 'Sent' },
        { id: '3', recipient: 'Priya Patel', initials: 'PP', channelIcon: 'heroicons_outline:phone', subject: 'Interview Confirmation', timestamp: 'May 28', status: 'Read' },
        { id: '4', recipient: 'Liam Chen', initials: 'LC', channelIcon: 'heroicons_outline:pencil-square', subject: 'Counselor evaluation note added', timestamp: 'May 27', status: 'Read' },
        { id: '5', recipient: 'Aisha Bello', initials: 'AB', channelIcon: 'heroicons_outline:envelope', subject: 'Welcome to the Fall 2026 cohort', timestamp: 'May 25', status: 'Failed' },
        { id: '6', recipient: 'Noah Williams', initials: 'NW', channelIcon: 'heroicons_outline:chat-bubble-left', subject: 'Campus Tour Details', timestamp: 'May 24', status: 'Delivered' },
        { id: '7', recipient: 'Emma Davis', initials: 'ED', channelIcon: 'heroicons_outline:phone', subject: 'Follow-up regarding financial aid', timestamp: 'May 22', status: 'Sent' },
        { id: '8', recipient: 'Lucas Smith', initials: 'LS', channelIcon: 'heroicons_outline:envelope', subject: 'Housing application open', timestamp: 'May 20', status: 'Read' },
    ];
}
