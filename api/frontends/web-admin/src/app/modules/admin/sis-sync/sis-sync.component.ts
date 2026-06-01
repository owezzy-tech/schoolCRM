import {
    ChangeDetectionStrategy,
    Component,
    computed,
    signal,
} from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatChipsModule } from '@angular/material/chips';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatSelectModule } from '@angular/material/select';
import {
    SIS_SYNC_EVENT_STATUSES,
    SIS_SYNC_JOB_STATUSES,
    SisSyncEvent,
    SisSyncEventStatus,
    SisSyncJob,
    SisSyncJobStatus,
} from 'app/core/admissions/admissions.types';

@Component({
    selector: 'app-sis-sync',
    imports: [
        FormsModule,
        MatButtonModule,
        MatChipsModule,
        MatFormFieldModule,
        MatIconModule,
        MatSelectModule,
    ],
    changeDetection: ChangeDetectionStrategy.OnPush,
    templateUrl: './sis-sync.component.html',
})
export class SisSyncComponent {
    readonly jobStatuses = SIS_SYNC_JOB_STATUSES;
    readonly eventStatuses = SIS_SYNC_EVENT_STATUSES;
    readonly selectedJobStatus = signal<SisSyncJobStatus | ''>('');
    readonly selectedEventStatus = signal<SisSyncEventStatus | ''>('');

    readonly jobs: readonly SisSyncJob[] = [
        {
            id: 'sync-job-20260601-nightly',
            name: 'Nightly SIS reconciliation',
            status: 'SUCCEEDED',
            direction: 'INBOUND',
            startedAt: 'Jun 1, 2026 01:00',
            completedAt: 'Jun 1, 2026 01:18',
            recordsPulled: 1842,
            recordsPushed: 0,
            eventsRequeued: 3,
            retryable: false,
            owner: 'System scheduler',
            dateCreated: 'Jun 1, 2026 01:00',
            dateUpdated: 'Jun 1, 2026 01:18',
        },
        {
            id: 'sync-job-20260601-retry',
            name: 'Failed event reconciliation',
            status: 'RUNNING',
            direction: 'OUTBOUND',
            startedAt: 'Jun 1, 2026 12:35',
            recordsPulled: 0,
            recordsPushed: 12,
            eventsRequeued: 0,
            retryable: false,
            owner: 'Admissions operations',
            dateCreated: 'Jun 1, 2026 12:35',
            dateUpdated: 'Jun 1, 2026 12:42',
        },
        {
            id: 'sync-job-20260531-nightly',
            name: 'Nightly SIS reconciliation',
            status: 'RETRY_READY',
            direction: 'INBOUND',
            startedAt: 'May 31, 2026 01:00',
            completedAt: 'May 31, 2026 01:07',
            recordsPulled: 613,
            recordsPushed: 0,
            eventsRequeued: 8,
            failureReason: 'SIS enrollment endpoint timed out after page 4.',
            retryable: true,
            owner: 'System scheduler',
            dateCreated: 'May 31, 2026 01:00',
            dateUpdated: 'May 31, 2026 01:07',
        },
    ];

    readonly events: readonly SisSyncEvent[] = [
        this.event(
            'evt-app-submit-3024',
            'APPLICATION_SUBMISSION',
            'SUCCEEDED',
            'OUTBOUND',
            'application',
            'APP-3024',
            'Application submission accepted by SIS.'
        ),
        this.event(
            'evt-decision-3020',
            'APPLICATION_DECISION',
            'PROCESSING',
            'OUTBOUND',
            'application',
            'APP-3020',
            'Decision payload is being pushed to SIS.',
            { attempts: 1 }
        ),
        this.event(
            'evt-document-3023',
            'DOCUMENT_STATUS',
            'RETRY_READY',
            'OUTBOUND',
            'document',
            'DOC-8841',
            'Document status failed and is ready for retry.',
            {
                attempts: 3,
                failureReason: 'SIS rejected checklist code transcript-final.',
                nextRetryAt: 'Jun 1, 2026 13:00',
            }
        ),
        this.event(
            'evt-terms-nightly',
            'BATCH_TERMS_PULL',
            'SUCCEEDED',
            'INBOUND',
            'academicTerm',
            'TERM-2026-FALL',
            'Pulled SIS-owned academic terms during nightly batch.'
        ),
        this.event(
            'evt-enrollment-intent-3019',
            'ENROLLMENT_INTENT',
            'QUEUED',
            'OUTBOUND',
            'application',
            'APP-3019',
            'Enrollment intent queued after applicant confirmation.'
        ),
    ];

    readonly filteredJobs = computed(() => {
        const status = this.selectedJobStatus();
        return status
            ? this.jobs.filter((job) => job.status === status)
            : this.jobs;
    });

    readonly filteredEvents = computed(() => {
        const status = this.selectedEventStatus();
        return status
            ? this.events.filter((event) => event.status === status)
            : this.events;
    });

    readonly failedEventCount = computed(
        () =>
            this.events.filter(
                (event) =>
                    event.status === 'FAILED' || event.status === 'RETRY_READY'
            ).length
    );

    readonly activeJobCount = computed(
        () =>
            this.jobs.filter(
                (job) => job.status === 'RUNNING' || job.status === 'QUEUED'
            ).length
    );

    readonly lastBatchJob = computed(() => this.jobs[0]);

    setJobStatus(status: SisSyncJobStatus | ''): void {
        this.selectedJobStatus.set(status);
    }

    setEventStatus(status: SisSyncEventStatus | ''): void {
        this.selectedEventStatus.set(status);
    }

    statusClass(status: SisSyncJobStatus | SisSyncEventStatus): string {
        switch (status) {
            case 'SUCCEEDED':
                return 'bg-green-100 text-green-700 ring-green-200';
            case 'RUNNING':
            case 'PROCESSING':
                return 'bg-blue-100 text-blue-700 ring-blue-200';
            case 'QUEUED':
                return 'bg-slate-100 text-slate-700 ring-slate-200';
            case 'FAILED':
                return 'bg-red-100 text-red-700 ring-red-200';
            case 'RETRY_READY':
                return 'bg-amber-100 text-amber-700 ring-amber-200';
        }
    }

    formatLabel(value: string): string {
        return value.replaceAll('_', ' ');
    }

    private event(
        id: string,
        eventType: SisSyncEvent['eventType'],
        status: SisSyncEventStatus,
        direction: SisSyncEvent['direction'],
        resourceType: string,
        resourceID: string,
        auditMessage: string,
        overrides: Partial<SisSyncEvent> = {}
    ): SisSyncEvent {
        return {
            id,
            eventType,
            status,
            direction,
            resourceType,
            resourceID,
            payloadHash: `sha256:${id}`,
            attempts: 0,
            auditMessage,
            dateCreated: 'Jun 1, 2026 12:30',
            dateUpdated: 'Jun 1, 2026 12:42',
            ...overrides,
        };
    }
}
