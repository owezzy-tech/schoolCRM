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
    KENYA_ADAPTERS,
    KenyaAdapter,
    SIS_SYNC_EVENT_STATUSES,
    SIS_SYNC_JOB_STATUSES,
    SisSyncEvent,
    SisSyncEventStatus,
    SisSyncJob,
    SisSyncJobStatus,
    SisSyncOperation,
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
    readonly adapters = KENYA_ADAPTERS;
    readonly jobStatuses = SIS_SYNC_JOB_STATUSES;
    readonly eventStatuses = SIS_SYNC_EVENT_STATUSES;
    readonly selectedJobStatus = signal<SisSyncJobStatus | ''>('');
    readonly selectedEventStatus = signal<SisSyncEventStatus | ''>('');
    readonly selectedAdapter = signal<KenyaAdapter | ''>('');

    readonly jobs: readonly SisSyncJob[] = [
        {
            id: 'sync-job-kuccps-placement-pull',
            name: 'KUCCPS placement pull',
            adapter: 'KUCCPS',
            operation: 'PLACEMENT_PULL',
            status: 'SUCCEEDED',
            direction: 'INBOUND',
            startedAt: 'Jun 1, 2026 01:00',
            completedAt: 'Jun 1, 2026 01:18',
            recordsPulled: 482,
            recordsPushed: 0,
            eventsRequeued: 0,
            attempts: 1,
            maxAttempts: 3,
            externalRef: 'KUCCPS-2026-PLACEMENTS',
            retryable: false,
            owner: 'System scheduler',
            dateCreated: 'Jun 1, 2026 01:00',
            dateUpdated: 'Jun 1, 2026 01:18',
        },
        {
            id: 'sync-job-knec-verification',
            name: 'KNEC result verification',
            adapter: 'KNEC',
            operation: 'RESULT_VERIFICATION',
            status: 'FAILED',
            direction: 'INBOUND',
            startedAt: 'Jun 1, 2026 09:15',
            completedAt: 'Jun 1, 2026 09:18',
            recordsPulled: 0,
            recordsPushed: 0,
            eventsRequeued: 7,
            attempts: 3,
            maxAttempts: 3,
            errorCode: 'KNEC_SERVICE_UNAVAILABLE',
            errorDetail: 'KNEC verification gateway returned 503.',
            lastErrorAt: 'Jun 1, 2026 09:18',
            failureReason: 'KNEC result verification is isolated for retry.',
            retryable: true,
            owner: 'Admissions operations',
            dateCreated: 'Jun 1, 2026 09:15',
            dateUpdated: 'Jun 1, 2026 09:18',
        },
        {
            id: 'sync-job-iprs-identity',
            name: 'IPRS identity verification',
            adapter: 'IPRS',
            operation: 'IDENTITY_VERIFICATION',
            status: 'RETRY_READY',
            direction: 'OUTBOUND',
            startedAt: 'Jun 1, 2026 10:30',
            completedAt: 'Jun 1, 2026 10:34',
            recordsPulled: 0,
            recordsPushed: 36,
            eventsRequeued: 12,
            attempts: 2,
            maxAttempts: 3,
            nextRetryAt: 'Jun 1, 2026 11:00',
            errorCode: 'IPRS_TIMEOUT',
            errorDetail: 'Identity endpoint timed out after 10 seconds.',
            lastErrorAt: 'Jun 1, 2026 10:34',
            failureReason: 'IPRS timeout does not block KUCCPS or KNEC queues.',
            retryable: true,
            owner: 'Admissions operations',
            dateCreated: 'Jun 1, 2026 10:30',
            dateUpdated: 'Jun 1, 2026 10:34',
        },
        {
            id: 'sync-job-mpesa-daraja',
            name: 'M-Pesa Daraja payment status',
            adapter: 'MPESA_DARAJA',
            operation: 'TRANSACTION_QUERY',
            status: 'RUNNING',
            direction: 'OUTBOUND',
            startedAt: 'Jun 1, 2026 12:35',
            recordsPulled: 0,
            recordsPushed: 21,
            eventsRequeued: 0,
            attempts: 1,
            maxAttempts: 3,
            retryable: false,
            owner: 'Finance operations',
            dateCreated: 'Jun 1, 2026 12:35',
            dateUpdated: 'Jun 1, 2026 12:42',
        },
        {
            id: 'sync-job-celcom-sms',
            name: 'Celcom Africa SMS delivery',
            adapter: 'CELCOM_AFRICA_SMS',
            operation: 'DELIVERY_REPORT_PULL',
            status: 'QUEUED',
            direction: 'INBOUND',
            recordsPulled: 0,
            recordsPushed: 0,
            eventsRequeued: 0,
            attempts: 0,
            maxAttempts: 3,
            retryable: false,
            owner: 'Communications',
            dateCreated: 'Jun 1, 2026 12:45',
            dateUpdated: 'Jun 1, 2026 12:45',
        },
        {
            id: 'sync-job-whatsapp-cloud',
            name: 'WhatsApp Cloud webhooks',
            adapter: 'WHATSAPP_CLOUD',
            operation: 'WA_WEBHOOK_INBOUND',
            status: 'SUCCEEDED',
            direction: 'INBOUND',
            startedAt: 'Jun 1, 2026 11:00',
            completedAt: 'Jun 1, 2026 11:02',
            recordsPulled: 18,
            recordsPushed: 0,
            eventsRequeued: 1,
            attempts: 1,
            maxAttempts: 3,
            retryable: false,
            owner: 'Communications',
            dateCreated: 'Jun 1, 2026 11:00',
            dateUpdated: 'Jun 1, 2026 11:02',
        },
    ];

    readonly events: readonly SisSyncEvent[] = [
        this.event(
            'evt-kuccps-placement-pull',
            'KUCCPS',
            'PLACEMENT_PULL',
            'KUCCPS_PLACEMENT_PULL',
            'SUCCEEDED',
            'INBOUND',
            'placement',
            'KUCCPS-PLC-3024',
            'Pulled KUCCPS placement row for applicant matching.',
            {
                jobID: 'sync-job-kuccps-placement-pull',
                externalRef: 'KUCCPS-PLC-3024',
            }
        ),
        this.event(
            'evt-kuccps-placement-confirm',
            'KUCCPS',
            'PLACEMENT_CONFIRM',
            'KUCCPS_PLACEMENT_CONFIRM',
            'RETRY_READY',
            'OUTBOUND',
            'placement',
            'KUCCPS-PLC-3025',
            'Confirmation callback is ready for isolated retry.',
            {
                attempts: 2,
                jobID: 'sync-job-kuccps-placement-pull',
                externalRef: 'KUCCPS-PLC-3025',
                nextRetryAt: 'Jun 1, 2026 13:10',
                errorCode: 'KUCCPS_CALLBACK_TIMEOUT',
                errorDetail: 'Placement confirmation endpoint timed out.',
                lastErrorAt: 'Jun 1, 2026 12:58',
            }
        ),
        this.event(
            'evt-knec-result-verification',
            'KNEC',
            'RESULT_VERIFICATION',
            'KNEC_RESULT_VERIFICATION',
            'FAILED',
            'INBOUND',
            'kcseResult',
            'KCSE-2025-397001001',
            'KNEC result verification failed after maximum attempts.',
            {
                attempts: 3,
                jobID: 'sync-job-knec-verification',
                externalRef: 'KCSE-2025-397001001',
                errorCode: 'KNEC_SERVICE_UNAVAILABLE',
                errorDetail: 'KNEC verification gateway returned 503.',
                lastErrorAt: 'Jun 1, 2026 09:18',
                failureReason:
                    'Manual review required if service remains unavailable.',
            }
        ),
        this.event(
            'evt-iprs-identity-verification',
            'IPRS',
            'IDENTITY_VERIFICATION',
            'IPRS_IDENTITY_VERIFICATION',
            'QUEUED',
            'OUTBOUND',
            'constituent',
            'ID-23456789',
            'National ID verification queued after applicant submission.',
            { jobID: 'sync-job-iprs-identity', externalRef: '23456789' }
        ),
        this.event(
            'evt-mpesa-stk-push',
            'MPESA_DARAJA',
            'STK_PUSH',
            'MPESA_STK_PUSH',
            'SUCCEEDED',
            'OUTBOUND',
            'applicationFee',
            'FEE-3024',
            'STK push accepted by Daraja.',
            {
                jobID: 'sync-job-mpesa-daraja',
                externalReceiptID: 'ws_CO_01062026124511234',
            }
        ),
        this.event(
            'evt-mpesa-c2b-callback',
            'MPESA_DARAJA',
            'C2B_CALLBACK',
            'MPESA_C2B_CALLBACK',
            'FAILED',
            'INBOUND',
            'applicationFee',
            'FEE-3018',
            'C2B callback signature validation failed.',
            {
                attempts: 3,
                jobID: 'sync-job-mpesa-daraja',
                externalRef: 'OEZ12345',
                externalReceiptID: 'QF12AB34CD',
                errorCode: 'MPESA_CALLBACK_SIGNATURE_INVALID',
                errorDetail: 'Callback signature did not match shared secret.',
                lastErrorAt: 'Jun 1, 2026 12:40',
            }
        ),
        this.event(
            'evt-celcom-sms-outbound',
            'CELCOM_AFRICA_SMS',
            'SMS_OUTBOUND',
            'SMS_OUTBOUND',
            'PROCESSING',
            'OUTBOUND',
            'notification',
            'SMS-3024',
            'Offer notification SMS is being delivered via Celcom Africa.',
            {
                jobID: 'sync-job-celcom-sms',
                externalRef: 'CELCOM-MSG-8821',
                attempts: 1,
            }
        ),
        this.event(
            'evt-celcom-delivery-report',
            'CELCOM_AFRICA_SMS',
            'DELIVERY_REPORT_PULL',
            'SMS_DELIVERY_REPORT',
            'SUCCEEDED',
            'INBOUND',
            'notification',
            'SMS-3022',
            'Delivery report confirmed applicant receipt.',
            { jobID: 'sync-job-celcom-sms', externalRef: 'CELCOM-MSG-8819' }
        ),
        this.event(
            'evt-whatsapp-message-send',
            'WHATSAPP_CLOUD',
            'WA_MESSAGE_SEND',
            'WHATSAPP_MESSAGE_SEND',
            'SUCCEEDED',
            'OUTBOUND',
            'notification',
            'WA-3024',
            'WhatsApp offer reminder accepted by Cloud API.',
            {
                jobID: 'sync-job-whatsapp-cloud',
                externalReceiptID: 'wamid.HBgMMjU0...',
            }
        ),
        this.event(
            'evt-whatsapp-webhook-inbound',
            'WHATSAPP_CLOUD',
            'WA_WEBHOOK_INBOUND',
            'WHATSAPP_WEBHOOK_INBOUND',
            'RETRY_READY',
            'INBOUND',
            'notification',
            'WA-3019',
            'Inbound webhook replay is ready after parsing error.',
            {
                attempts: 2,
                jobID: 'sync-job-whatsapp-cloud',
                nextRetryAt: 'Jun 1, 2026 13:20',
                errorCode: 'WHATSAPP_WEBHOOK_PARSE_ERROR',
                errorDetail: 'Webhook payload missed message status field.',
                lastErrorAt: 'Jun 1, 2026 12:52',
            }
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
        const adapter = this.selectedAdapter();
        return this.events.filter(
            (event) =>
                (!status || event.status === status) &&
                (!adapter || event.adapter === adapter)
        );
    });

    readonly adapterHealth = computed(() =>
        this.adapters.map((adapter) => {
            const events = this.events.filter(
                (event) => event.adapter === adapter
            );
            const failed = events.filter(
                (event) =>
                    event.status === 'FAILED' || event.status === 'RETRY_READY'
            ).length;

            return { adapter, failed, total: events.length };
        })
    );

    readonly adapterAggregates = computed(() =>
        this.adapters.map((adapter) => {
            const jobs = this.jobs.filter((job) => job.adapter === adapter);
            const events = this.events.filter(
                (event) => event.adapter === adapter
            );
            return {
                adapter,
                jobs: jobs.length,
                running: jobs.filter(
                    (job) => job.status === 'RUNNING' || job.status === 'QUEUED'
                ).length,
                failed: events.filter(
                    (event) =>
                        event.status === 'FAILED' ||
                        event.status === 'RETRY_READY'
                ).length,
                lastJobAt: jobs[0]?.dateUpdated ?? '—',
            };
        })
    );

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

    setAdapter(adapter: KenyaAdapter | ''): void {
        this.selectedAdapter.set(adapter);
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
        adapter: KenyaAdapter,
        operation: SisSyncOperation,
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
            adapter,
            operation,
            eventType,
            status,
            direction,
            resourceType,
            resourceID,
            payloadHash: `sha256:${id}`,
            attempts: 0,
            maxAttempts: 3,
            auditMessage,
            dateCreated: 'Jun 1, 2026 12:30',
            dateUpdated: 'Jun 1, 2026 12:42',
            ...overrides,
        };
    }
}
