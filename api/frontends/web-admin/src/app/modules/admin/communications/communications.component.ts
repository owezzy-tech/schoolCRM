import {
    ChangeDetectionStrategy,
    Component,
    computed,
    inject,
    signal,
} from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { RouterLink } from '@angular/router';
import { AdmissionsService } from 'app/core/admissions/admissions.service';
import { jsonApiErrorMessage } from 'app/core/api/json-api';
import { catchError, of } from 'rxjs';
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
    private readonly admissionsService = inject(AdmissionsService);

    readonly activeChannel = signal<CommunicationChannel | 'ALL'>('ALL');
    readonly directionFilter = signal<CommunicationDirection | 'ALL'>('ALL');
    readonly statusFilter = signal<CommunicationStatus | 'ALL'>('ALL');
    readonly loading = signal(false);
    readonly errorMessage = signal('');
    readonly records = signal<CommunicationRecord[]>([]);

    readonly tabs: CommunicationTab[] = [
        { label: 'All Messages', channel: 'ALL' },
        { label: 'SMS', channel: 'SMS' },
        { label: 'WhatsApp', channel: 'WHATSAPP' },
        { label: 'Email', channel: 'EMAIL' },
        { label: 'Phone calls', channel: 'PHONE_CALL' },
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
            (record) => record.applicationID
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

    constructor() {
        this.loadRecords();
    }

    loadRecords(): void {
        const channel = this.activeChannel();
        const direction = this.directionFilter();
        const status = this.statusFilter();

        this.loading.set(true);
        this.errorMessage.set('');

        this.admissionsService
            .queryCommunications({
                page: 1,
                rows: 50,
                orderBy: 'occurred_at,DESC',
                channel: channel === 'ALL' ? undefined : channel,
                direction: direction === 'ALL' ? undefined : direction,
                status: status === 'ALL' ? undefined : status,
            })
            .pipe(
                catchError((error) => {
                    this.errorMessage.set(
                        jsonApiErrorMessage(
                            error,
                            'Unable to load communications.'
                        )
                    );
                    this.records.set([]);

                    return of(undefined);
                })
            )
            .subscribe((result) => {
                this.loading.set(false);

                if (!result) {
                    return;
                }

                this.records.set(result.items);
            });
    }

    setChannel(channel: CommunicationChannel | 'ALL'): void {
        this.activeChannel.set(channel);
        this.loadRecords();
    }

    setDirectionFilter(event: Event): void {
        const value = (event.target as HTMLSelectElement).value;

        if (this.isDirectionFilter(value)) {
            this.directionFilter.set(value);
            this.loadRecords();
        }
    }

    setStatusFilter(event: Event): void {
        const value = (event.target as HTMLSelectElement).value;

        if (this.isStatusFilter(value)) {
            this.statusFilter.set(value);
            this.loadRecords();
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
            SMS: 'SMS',
            WHATSAPP: 'WhatsApp',
            EMAIL: 'Email',
            PHONE_CALL: 'Phone call',
            NOTIFICATION: 'Notification',
        };

        return labels[channel];
    }

    channelIcon(channel: CommunicationChannel): string {
        const icons: Record<CommunicationChannel, string> = {
            SMS: 'heroicons_outline:chat-bubble-left',
            WHATSAPP: 'heroicons_outline:device-phone-mobile',
            EMAIL: 'heroicons_outline:envelope',
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
