import {
    ChangeDetectionStrategy,
    Component,
    computed,
    signal,
} from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { RouterLink } from '@angular/router';
import { COMMUNICATIONS } from './data/communications.mock';
import {
    CommunicationChannel,
    CommunicationDirection,
    CommunicationRecord,
    CommunicationStatus,
    CommunicationSummary,
} from './models/communication.types';

interface CommunicationTab {
    label: string;
    channel: CommunicationChannel | 'ALL';
}

@Component({
    selector: 'app-communications',
    imports: [MatButtonModule, MatIconModule, RouterLink],
    changeDetection: ChangeDetectionStrategy.OnPush,
    templateUrl: './communications.component.html',
})
export class CommunicationsComponent {
    readonly activeChannel = signal<CommunicationChannel | 'ALL'>('ALL');
    readonly directionFilter = signal<CommunicationDirection | 'ALL'>('ALL');
    readonly statusFilter = signal<CommunicationStatus | 'ALL'>('ALL');
    readonly records = signal<CommunicationRecord[]>(COMMUNICATIONS);

    readonly tabs: CommunicationTab[] = [
        { label: 'All Messages', channel: 'ALL' },
        { label: 'Email', channel: 'EMAIL' },
        { label: 'SMS', channel: 'SMS' },
        { label: 'Phone Calls', channel: 'PHONE_CALL' },
        { label: 'Notifications', channel: 'NOTIFICATION' },
    ];

    readonly directionOptions: (CommunicationDirection | 'ALL')[] = [
        'ALL',
        'INBOUND',
        'OUTBOUND',
    ];

    readonly statusOptions: (CommunicationStatus | 'ALL')[] = [
        'ALL',
        'SENT',
        'DELIVERED',
        'FAILED',
        'OPENED',
        'BOUNCED',
        'REPLIED',
        'LOGGED',
    ];

    readonly filteredRecords = computed(() =>
        this.records().filter((record) => {
            const channelMatches =
                this.activeChannel() === 'ALL' ||
                record.channel === this.activeChannel();
            const directionMatches =
                this.directionFilter() === 'ALL' ||
                record.direction === this.directionFilter();
            const statusMatches =
                this.statusFilter() === 'ALL' ||
                record.status === this.statusFilter();

            return channelMatches && directionMatches && statusMatches;
        })
    );

    readonly summaries = computed<CommunicationSummary[]>(() => {
        const records = this.records();
        const providerTracked = records.filter((record) => record.provider).length;
        const failed = records.filter((record) =>
            ['FAILED', 'BOUNCED'].includes(record.status)
        ).length;
        const linkedApplications = records.filter(
            (record) => record.applicationId
        ).length;

        return [
            {
                label: 'Total communications',
                value: String(records.length),
                icon: 'heroicons_outline:chat-bubble-left-right',
                tone: 'blue',
            },
            {
                label: 'Provider-tracked',
                value: String(providerTracked),
                icon: 'heroicons_outline:signal',
                tone: 'green',
            },
            {
                label: 'Linked applications',
                value: String(linkedApplications),
                icon: 'heroicons_outline:document-text',
                tone: 'amber',
            },
            {
                label: 'Needs attention',
                value: String(failed),
                icon: 'heroicons_outline:exclamation-triangle',
                tone: 'red',
            },
        ];
    });

    readonly phoneLogDraft = {
        constituent: 'Priya Patel',
        application: 'APP-3019',
        duration: '12 min',
        outcome: 'Interview confirmed',
        notes: 'Confirmed virtual interview schedule and scholarship document expectations.',
    };

    setChannel(channel: CommunicationChannel | 'ALL'): void {
        this.activeChannel.set(channel);
    }

    setDirectionFilter(event: Event): void {
        const value = (event.target as HTMLSelectElement).value;

        if (this.isDirectionFilter(value)) {
            this.directionFilter.set(value);
        }
    }

    setStatusFilter(event: Event): void {
        const value = (event.target as HTMLSelectElement).value;

        if (this.isStatusFilter(value)) {
            this.statusFilter.set(value);
        }
    }

    statusLabel(status: CommunicationStatus): string {
        return status
            .split('_')
            .map((part) => part[0] + part.slice(1).toLowerCase())
            .join(' ');
    }

    channelLabel(channel: CommunicationChannel): string {
        const labels: Record<CommunicationChannel, string> = {
            EMAIL: 'Email',
            SMS: 'SMS',
            PHONE_CALL: 'Phone call',
            NOTIFICATION: 'Notification',
        };

        return labels[channel];
    }

    channelIcon(channel: CommunicationChannel): string {
        const icons: Record<CommunicationChannel, string> = {
            EMAIL: 'heroicons_outline:envelope',
            SMS: 'heroicons_outline:chat-bubble-left',
            PHONE_CALL: 'heroicons_outline:phone',
            NOTIFICATION: 'heroicons_outline:bell-alert',
        };

        return icons[channel];
    }

    private isDirectionFilter(
        value: string
    ): value is CommunicationDirection | 'ALL' {
        return this.directionOptions.includes(
            value as CommunicationDirection | 'ALL'
        );
    }

    private isStatusFilter(value: string): value is CommunicationStatus | 'ALL' {
        return this.statusOptions.includes(value as CommunicationStatus | 'ALL');
    }
}
