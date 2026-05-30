import { ChangeDetectionStrategy, Component } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatTabsModule } from '@angular/material/tabs';
import { MatSlideToggleModule } from '@angular/material/slide-toggle';

interface TeamMember {
    id: string;
    name: string;
    email: string;
    role: string;
    lastSeen: string;
}

interface PipelineStage {
    id: string;
    name: string;
    count: number;
}

interface Integration {
    id: string;
    name: string;
    description: string;
    icon: string;
    connected: boolean;
}

interface NotificationSetting {
    id: string;
    label: string;
    description: string;
    enabled: boolean;
}

@Component({
    selector: 'app-settings-global',
    standalone: true,
    imports: [MatButtonModule, MatIconModule, MatTabsModule, MatSlideToggleModule],
    changeDetection: ChangeDetectionStrategy.OnPush,
    templateUrl: './settings.component.html',
})
export class SettingsComponent {
    readonly teamMembers: readonly TeamMember[] = [
        { id: '1', name: 'Avery Johnson', email: 'avery@example.com', role: 'Admin', lastSeen: 'Active now' },
        { id: '2', name: 'Diego Ramirez', email: 'diego@example.com', role: 'Counselor', lastSeen: '2h ago' },
        { id: '3', name: 'Priya Patel', email: 'priya@example.com', role: 'Counselor', lastSeen: 'Yesterday' },
        { id: '4', name: 'Sam Chen', email: 'sam@example.com', role: 'Reviewer', lastSeen: '3 days ago' },
    ];

    readonly pipelineStages: readonly PipelineStage[] = [
        { id: '1', name: 'Inquiry', count: 1284 },
        { id: '2', name: 'Lead', count: 486 },
        { id: '3', name: 'Applicant', count: 342 },
        { id: '4', name: 'In Review', count: 156 },
        { id: '5', name: 'Admitted', count: 82 },
        { id: '6', name: 'Enrolled', count: 45 },
    ];

    readonly integrations: readonly Integration[] = [
        { id: '1', name: 'Google Workspace', description: 'Sync emails and calendar events', icon: 'heroicons_outline:envelope', connected: true },
        { id: '2', name: 'Slack', description: 'Get notifications in Slack channels', icon: 'heroicons_outline:chat-bubble-left-ellipsis', connected: true },
        { id: '3', name: 'Twilio', description: 'Send SMS campaigns and alerts', icon: 'heroicons_outline:device-phone-mobile', connected: false },
        { id: '4', name: 'Stripe', description: 'Process application fees and deposits', icon: 'heroicons_outline:credit-card', connected: false },
    ];

    readonly notificationSettings: readonly NotificationSetting[] = [
        { id: '1', label: 'New Application', description: 'When a candidate submits a new application', enabled: true },
        { id: '2', label: 'Document Uploaded', description: 'When a candidate uploads a required document', enabled: true },
        { id: '3', label: 'Review Completed', description: 'When a reviewer submits an evaluation', enabled: false },
        { id: '4', label: 'Message Received', description: 'When a candidate replies to an email or SMS', enabled: true },
        { id: '5', label: 'Daily Summary', description: 'Receive a daily digest of all activity', enabled: false },
    ];
}
