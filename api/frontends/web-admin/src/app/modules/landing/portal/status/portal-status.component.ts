import { TitleCasePipe } from '@angular/common';
import {
    ChangeDetectionStrategy,
    Component,
    DestroyRef,
    OnInit,
    computed,
    inject,
    signal,
} from '@angular/core';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { RouterLink } from '@angular/router';
import { AdmissionsService } from 'app/core/admissions/admissions.service';
import {
    AcademicTerm,
    AdmissionsDocument,
    Application,
    ApplicationFee,
    ApplicationFeeStatus,
    ApplicationStatus,
    ChecklistItem,
    Program,
} from 'app/core/admissions/admissions.types';
import { jsonApiErrorMessage } from 'app/core/api/json-api';
import { PortalAuthService } from 'app/core/portal/portal-auth.service';
import { FilePondComponent } from 'app/shared/components/file-upload/file-pond.component';
import { FilePondFile } from 'filepond';
import { Observable, finalize, forkJoin, map, of, switchMap } from 'rxjs';

type PortalDocumentStatus =
    | 'verified'
    | 'received'
    | 'missing'
    | 'uploading'
    | 'pending-review'
    | 'rejected';

interface PortalDocument {
    id: string;
    checklistItemID: string;
    label: string;
    description: string;
    status: PortalDocumentStatus;
    required: boolean;
    serverId?: string;
    fileName?: string;
}

interface TimelineItem {
    label: string;
    date: string;
    done: boolean;
}

const STATUS_PROGRESS: Record<ApplicationStatus, number> = {
    DRAFT: 10,
    SUBMITTED: 30,
    AWAITING_DOCUMENTS: 45,
    READY_FOR_REVIEW: 60,
    IN_REVIEW: 70,
    DECISION_PENDING: 85,
    ADMITTED: 100,
    DENIED: 100,
    WAITLISTED: 90,
    DEFERRED: 80,
    WITHDRAWN: 100,
    ENROLLED: 100,
};

const TIMELINE_STAGES: Array<{
    status: ApplicationStatus;
    label: string;
}> = [
    { status: 'SUBMITTED', label: 'Submitted' },
    { status: 'AWAITING_DOCUMENTS', label: 'Documents requested' },
    { status: 'READY_FOR_REVIEW', label: 'Ready for review' },
    { status: 'IN_REVIEW', label: 'Under review' },
    { status: 'DECISION_PENDING', label: 'Decision pending' },
];

@Component({
    selector: 'app-portal-status',
    imports: [
        FilePondComponent,
        MatButtonModule,
        MatIconModule,
        MatProgressBarModule,
        RouterLink,
        TitleCasePipe,
    ],
    template: `
        <div class="flex flex-auto flex-col px-6 py-12 sm:px-10 lg:px-16">
            <div class="mx-auto w-full max-w-6xl">
                @if (!hasPortalAccess()) {
                    <div
                        class="mb-6 rounded-2xl border border-amber-200 bg-amber-50 p-5 text-amber-900 dark:border-amber-700 dark:bg-amber-900/20 dark:text-amber-100"
                    >
                        <div
                            class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between"
                        >
                            <div>
                                <h2 class="text-lg font-semibold">
                                    Portal access required
                                </h2>
                                <p class="mt-1 text-sm">
                                    Use your submitted application email to
                                    create a temporary status token before
                                    uploading documents.
                                </p>
                            </div>
                            <a
                                routerLink="/portal"
                                mat-flat-button
                                color="primary"
                                >Get access</a
                            >
                        </div>
                    </div>
                }

                @if (errorMessage()) {
                    <div
                        class="mb-6 rounded-2xl border border-red-200 bg-red-50 p-5 text-sm text-red-700 dark:border-red-900 dark:bg-red-900/20 dark:text-red-200"
                    >
                        {{ errorMessage() }}
                    </div>
                }

                <!-- Header section -->
                <div
                    class="flex flex-col items-start gap-4 sm:flex-row sm:items-center sm:justify-between"
                >
                    <div class="flex flex-col">
                        <h1
                            class="text-default text-3xl font-bold tracking-tight sm:text-4xl"
                        >
                            Application status
                        </h1>
                        <h2 class="text-secondary mt-1 text-lg">
                            {{ applicationTitle() }}
                        </h2>
                    </div>
                    <div
                        class="inline-flex items-center rounded-full px-3 py-1 text-sm font-medium"
                        [class]="statusBadgeClass(application()?.status)"
                    >
                        {{ statusLabel(application()?.status) }}
                    </div>
                </div>

                <!-- Overall progress card -->
                <div
                    class="bg-card mt-6 flex flex-col rounded-2xl border p-6 shadow-sm"
                >
                    <div class="mb-4 flex items-center justify-between gap-4">
                        <h3 class="text-default text-lg font-semibold">
                            Overall progress
                        </h3>
                        @if (isLoading()) {
                            <span class="text-secondary text-sm"
                                >Loading live status…</span
                            >
                        }
                    </div>
                    <mat-progress-bar
                        mode="determinate"
                        [value]="overallProgress()"
                        class="mb-2"
                    ></mat-progress-bar>
                    <p class="text-secondary text-sm font-medium">
                        {{ overallProgress() }}% complete &middot;
                        {{ completedStageCount() }} of {{ timeline().length }}
                        stages done.
                    </p>
                </div>

                <!-- Fee status card -->
                <div class="bg-card mt-6 rounded-2xl border p-6 shadow-sm">
                    <div
                        class="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between"
                    >
                        <div>
                            <h3 class="text-default text-lg font-semibold">
                                Application fee
                            </h3>
                            <p class="text-secondary mt-1 text-sm">
                                M-Pesa processing is not active in this preview.
                                The portal derives fee readiness from the live
                                application status until the fee API is enabled.
                            </p>
                        </div>
                        <span
                            class="w-fit rounded-full px-3 py-1 text-sm font-semibold"
                            [class]="feeStatusClass(applicationFee().status)"
                        >
                            {{ formatFeeStatus(applicationFee().status) }}
                        </span>
                    </div>

                    <div class="mt-5 grid grid-cols-1 gap-4 sm:grid-cols-3">
                        <div
                            class="rounded-2xl bg-gray-50 p-4 dark:bg-gray-900"
                        >
                            <div
                                class="text-secondary text-xs font-semibold uppercase tracking-wide"
                            >
                                Amount
                            </div>
                            <div class="mt-1 text-2xl font-bold">
                                {{ formatAmount(applicationFee()) }}
                            </div>
                        </div>
                        <div
                            class="rounded-2xl bg-gray-50 p-4 dark:bg-gray-900"
                        >
                            <div
                                class="text-secondary text-xs font-semibold uppercase tracking-wide"
                            >
                                Due date
                            </div>
                            <div class="mt-1 text-2xl font-bold">
                                {{ applicationFee().dueAt || '—' }}
                            </div>
                        </div>
                        <div
                            class="rounded-2xl bg-gray-50 p-4 dark:bg-gray-900"
                        >
                            <div
                                class="text-secondary text-xs font-semibold uppercase tracking-wide"
                            >
                                Payment channel
                            </div>
                            <div class="mt-1 text-2xl font-bold capitalize">
                                {{
                                    formatPaymentProvider(
                                        applicationFee().provider
                                    )
                                }}
                            </div>
                        </div>
                    </div>
                </div>

                <!-- 2-column grid -->
                <div class="mt-6 grid grid-cols-1 gap-6 md:grid-cols-2">
                    <!-- Timeline Card -->
                    <div
                        class="bg-card flex flex-col rounded-2xl border p-6 shadow-sm"
                    >
                        <h3 class="text-default mb-6 text-lg font-semibold">
                            Timeline
                        </h3>
                        <div class="relative flex flex-col gap-6">
                            <div
                                class="absolute bottom-4 left-[11px] top-4 border-l-2 border-dashed border-gray-300"
                            ></div>

                            @for (item of timeline(); track item.label) {
                                <div
                                    class="relative z-10 flex items-start gap-4"
                                >
                                    <div
                                        class="mt-0.5 flex h-6 w-6 shrink-0 items-center justify-center rounded-full"
                                        [class.bg-primary]="item.done"
                                        [class.border-2]="!item.done"
                                        [class.border-gray-300]="!item.done"
                                        [class.bg-card]="!item.done"
                                    ></div>
                                    <div class="flex flex-col">
                                        <span
                                            class="text-default font-medium leading-tight"
                                            >{{ item.label }}</span
                                        >
                                        <span class="text-secondary text-sm">{{
                                            item.date
                                        }}</span>
                                    </div>
                                </div>
                            }
                        </div>
                    </div>

                    <!-- Documents Card -->
                    <div
                        class="bg-card flex flex-col rounded-2xl border p-6 shadow-sm"
                    >
                        <div class="flex items-start justify-between gap-4">
                            <div>
                                <h3 class="text-default text-lg font-semibold">
                                    Documents
                                </h3>
                                <p class="text-secondary mt-1 text-sm">
                                    Upload missing or rejected items. Files go
                                    to the upload service and this portal saves
                                    the returned storage ID as admissions
                                    document metadata.
                                </p>
                            </div>
                            <div
                                class="rounded-full bg-primary-50 px-3 py-1 text-xs font-semibold text-primary dark:bg-primary-900/20"
                            >
                                {{ uploadedDocumentCount() }} /
                                {{ documents().length }} received
                            </div>
                        </div>

                        @if (documents().length === 0 && !isLoading()) {
                            <div
                                class="text-secondary mt-6 rounded-2xl border border-dashed p-4 text-sm"
                            >
                                No checklist items have been assigned to this
                                application yet.
                            </div>
                        }

                        <div class="mt-4 flex flex-col divide-y">
                            @for (doc of documents(); track doc.id) {
                                <div class="py-4">
                                    <div
                                        class="flex items-start justify-between gap-4"
                                    >
                                        <div class="flex items-start gap-3">
                                            <mat-icon
                                                svgIcon="heroicons_outline:document-text"
                                                class="text-secondary mt-0.5 icon-size-5"
                                            ></mat-icon>
                                            <div>
                                                <div
                                                    class="flex items-center gap-2"
                                                >
                                                    <span
                                                        class="text-default font-medium"
                                                        >{{ doc.label }}</span
                                                    >
                                                    @if (doc.required) {
                                                        <span
                                                            class="rounded-full bg-red-50 px-2 py-0.5 text-xs font-semibold text-red-600 dark:bg-red-900/20 dark:text-red-300"
                                                            >Required</span
                                                        >
                                                    }
                                                </div>
                                                <p
                                                    class="text-secondary mt-1 text-sm"
                                                >
                                                    {{ doc.description }}
                                                </p>
                                                @if (doc.fileName) {
                                                    <p
                                                        class="text-secondary mt-1 text-xs"
                                                    >
                                                        {{ doc.fileName }}
                                                        @if (doc.serverId) {
                                                            <span
                                                                >&middot; Upload
                                                                ID
                                                                {{
                                                                    doc.serverId
                                                                }}</span
                                                            >
                                                        }
                                                    </p>
                                                }
                                            </div>
                                        </div>
                                        <div
                                            class="rounded-full px-2.5 py-0.5 text-xs font-medium"
                                            [class]="
                                                documentStatusClass(doc.status)
                                            "
                                        >
                                            {{ doc.status | titlecase }}
                                        </div>
                                    </div>

                                    @if (canUploadDocument(doc)) {
                                        <div
                                            class="mt-4 rounded-2xl border border-dashed bg-gray-50/70 p-3 dark:bg-gray-900/20"
                                        >
                                            <app-file-pond
                                                [acceptedFileTypes]="
                                                    acceptedDocumentTypes
                                                "
                                                [maxFileSize]="'10MB'"
                                                [server]="serverConfig"
                                                [disabled]="
                                                    doc.status === 'uploading'
                                                "
                                                [labelIdle]="
                                                    uploadLabelFor(doc)
                                                "
                                                (uploadStarted)="
                                                    markUploadStarted(doc.id)
                                                "
                                                (fileUploaded)="
                                                    markUploaded(doc.id, $event)
                                                "
                                                (fileReverted)="
                                                    revertUpload(doc.id)
                                                "
                                                (uploadAborted)="
                                                    revertUpload(doc.id)
                                                "
                                                (uploadError)="
                                                    markUploadFailed(doc.id)
                                                "
                                            />
                                        </div>
                                    }

                                    @if (doc.status === 'pending-review') {
                                        <div
                                            class="mt-3 flex items-center gap-2 rounded-xl bg-amber-50 px-3 py-2 text-sm text-amber-800 dark:bg-amber-900/20 dark:text-amber-200"
                                        >
                                            <mat-icon
                                                svgIcon="heroicons_outline:clock"
                                                class="icon-size-4"
                                            ></mat-icon>
                                            This upload is queued for admissions
                                            review.
                                        </div>
                                    }
                                </div>
                            }
                        </div>
                    </div>
                </div>

                <!-- Admissions officer card -->
                <div
                    class="mt-6 flex items-center gap-4 rounded-2xl bg-primary-50 p-6 dark:bg-primary-900/20"
                >
                    <div
                        class="flex h-12 w-12 shrink-0 items-center justify-center rounded-full bg-primary/10"
                    >
                        <mat-icon
                            svgIcon="heroicons_outline:envelope"
                            class="text-primary icon-size-6"
                        ></mat-icon>
                    </div>
                    <div class="flex flex-col">
                        <h3 class="text-default font-semibold">
                            Reach your admissions officer
                        </h3>
                        <p class="text-secondary mt-0.5">
                            Amina Otieno &middot; admissions&#64;schoolcrm.ac.ke
                            &middot; Replies within 1 business day, 8am-5pm EAT.
                        </p>
                    </div>
                </div>
            </div>
        </div>
    `,
    changeDetection: ChangeDetectionStrategy.OnPush,
    host: { class: 'flex w-full flex-auto flex-col' },
})
export class PortalStatusComponent implements OnInit {
    private readonly admissionsService = inject(AdmissionsService);
    private readonly destroyRef = inject(DestroyRef);
    private readonly portalAuthService = inject(PortalAuthService);

    readonly portalSession = this.portalAuthService.session;
    readonly acceptedDocumentTypes = ['application/pdf', 'image/*'];
    readonly serverConfig = {
        url: '/api/common/file-upload',
        process: {
            url: '',
            method: 'POST' as const,
        },
        revert: {
            url: '',
            method: 'DELETE' as const,
        },
    };

    readonly application = signal<Application | null>(null);
    readonly documents = signal<PortalDocument[]>([]);
    readonly errorMessage = signal<string | null>(null);
    readonly isLoading = signal(false);
    readonly terms = signal<AcademicTerm[]>([]);
    readonly programs = signal<Program[]>([]);

    readonly applicationTitle = computed(() => {
        const session = this.portalSession();
        const application = this.application();

        if (!application) {
            return session?.applicationID || 'No active application selected';
        }

        const program = this.programs().find(
            (item) => item.id === application.programID
        );
        const term = this.terms().find(
            (item) => item.id === application.academicTermID
        );

        return [
            application.id,
            program?.name ?? 'Programme pending',
            term?.name ?? 'Intake pending',
        ].join(' · ');
    });

    readonly applicationFee = computed<ApplicationFee>(() => {
        const application = this.application();
        const status = this.feeStatusFor(application?.status);

        return {
            id: `fee-${application?.id ?? 'pending'}`,
            applicationID:
                application?.id ??
                this.portalSession()?.applicationID ??
                'pending',
            amountCents: status === 'NOT_REQUIRED' ? 0 : 15000,
            currency: 'KES',
            status,
            provider: status === 'NOT_REQUIRED' ? 'not_required' : 'manual',
            dueAt: application?.submittedAt
                ? this.formatDate(application.submittedAt)
                : undefined,
            auditTrail: [],
        };
    });

    readonly overallProgress = computed(() => {
        const status = this.application()?.status;

        return status ? STATUS_PROGRESS[status] : 0;
    });

    readonly timeline = computed<TimelineItem[]>(() => {
        const application = this.application();
        if (!application) {
            return TIMELINE_STAGES.map((stage) => ({
                label: stage.label,
                date: 'Waiting for live application data',
                done: false,
            }));
        }

        const currentIndex = this.statusIndex(application.status);

        return TIMELINE_STAGES.map((stage) => {
            const stageIndex = this.statusIndex(stage.status);
            const done = currentIndex >= stageIndex;
            const date = this.timelineDate(application, stage.status, done);

            return { label: stage.label, date, done };
        });
    });

    readonly completedStageCount = computed(
        () => this.timeline().filter((item) => item.done).length
    );

    readonly uploadedDocumentCount = computed(
        () =>
            this.documents().filter(
                (document) => !['missing', 'rejected'].includes(document.status)
            ).length
    );

    ngOnInit(): void {
        this.loadStatus();
    }

    feeStatusClass(status: ApplicationFeeStatus): string {
        switch (status) {
            case 'PAID':
                return 'bg-green-100 text-green-700';
            case 'PENDING':
                return 'bg-amber-100 text-amber-700';
            case 'FAILED':
                return 'bg-red-100 text-red-700';
            case 'WAIVED':
                return 'bg-purple-100 text-purple-700';
            case 'REFUNDED':
                return 'bg-blue-100 text-blue-700';
            case 'NOT_REQUIRED':
                return 'bg-slate-100 text-secondary';
        }
    }

    formatFeeStatus(status: ApplicationFeeStatus): string {
        return status.replaceAll('_', ' ');
    }

    formatAmount(fee: ApplicationFee): string {
        if (fee.status === 'NOT_REQUIRED') {
            return 'No fee due';
        }

        return new Intl.NumberFormat('en-KE', {
            style: 'currency',
            currency: fee.currency,
        }).format(fee.amountCents / 100);
    }

    formatPaymentProvider(provider: ApplicationFee['provider']): string {
        switch (provider) {
            case 'manual':
                return 'M-Pesa';
            case 'not_required':
                return 'Not required';
            case 'stripe':
                return 'Stripe';
            case 'square':
                return 'Square';
        }
    }

    statusBadgeClass(status: ApplicationStatus | undefined): string {
        switch (status) {
            case 'ADMITTED':
            case 'ENROLLED':
                return 'bg-green-100 text-green-900';
            case 'DENIED':
            case 'WITHDRAWN':
                return 'bg-red-100 text-red-900';
            case 'WAITLISTED':
            case 'DEFERRED':
                return 'bg-blue-100 text-blue-900';
            case 'DRAFT':
                return 'bg-slate-100 text-slate-700';
            default:
                return 'bg-yellow-100 text-yellow-900';
        }
    }

    statusLabel(status: ApplicationStatus | undefined): string {
        return status ? status.replaceAll('_', ' ') : 'Waiting for access';
    }

    documentStatusClass(status: PortalDocumentStatus): string {
        switch (status) {
            case 'verified':
                return 'bg-green-100 text-green-800';
            case 'received':
                return 'bg-blue-100 text-blue-800';
            case 'missing':
            case 'rejected':
                return 'bg-red-100 text-red-800';
            case 'uploading':
            case 'pending-review':
                return 'bg-amber-100 text-amber-800';
        }
    }

    canUploadDocument(document: PortalDocument): boolean {
        return (
            this.hasPortalAccess() &&
            ['missing', 'received', 'rejected', 'uploading'].includes(
                document.status
            )
        );
    }

    hasPortalAccess(): boolean {
        return this.portalAuthService.hasValidSession();
    }

    uploadLabelFor(document: PortalDocument): string {
        if (document.status === 'received' || document.status === 'rejected') {
            return 'Reupload replacement or <span class="filepond--label-action">browse</span>';
        }

        return `Drop ${document.label.toLowerCase()} or <span class="filepond--label-action">browse</span>`;
    }

    markUploadStarted(documentId: string): void {
        this.updateDocument(documentId, { status: 'uploading' });
    }

    markUploaded(documentId: string, file: FilePondFile): void {
        const application = this.application();
        const document = this.documents().find(
            (item) => item.id === documentId
        );
        const storageKey = String(file.serverId ?? '');

        if (!application || !document || storageKey.length === 0) {
            this.updateDocument(documentId, { status: 'missing' });
            this.errorMessage.set(
                'The upload finished, but the upload service did not return a storage ID.'
            );
            return;
        }

        this.updateDocument(documentId, {
            status: 'pending-review',
            serverId: storageKey,
            fileName: file.filename,
        });

        this.admissionsService
            .createApplicantDocument(application.id, {
                checklistItemID: document.checklistItemID,
                fileName: file.filename,
                contentType: file.file.type || 'application/octet-stream',
                sizeBytes: file.file.size,
                storageKey,
            })
            .pipe(takeUntilDestroyed(this.destroyRef))
            .subscribe({
                next: (created) => {
                    this.mergeCreatedDocument(created);
                    this.errorMessage.set(null);
                },
                error: (error) => {
                    this.updateDocument(documentId, {
                        status: 'missing',
                        serverId: undefined,
                        fileName: undefined,
                    });
                    this.errorMessage.set(
                        jsonApiErrorMessage(
                            error,
                            'The file uploaded, but admissions could not save the document metadata.'
                        )
                    );
                },
            });
    }

    revertUpload(documentId: string): void {
        this.updateDocument(documentId, {
            status: 'missing',
            serverId: undefined,
            fileName: undefined,
        });
    }

    markUploadFailed(documentId: string): void {
        this.updateDocument(documentId, { status: 'missing' });
    }

    private loadStatus(): void {
        if (!this.hasPortalAccess()) {
            return;
        }

        this.isLoading.set(true);
        this.errorMessage.set(null);

        this.resolveApplicationID()
            .pipe(
                switchMap((applicationID) => {
                    if (!applicationID) {
                        throw new Error(
                            'No application is linked to this portal session yet.'
                        );
                    }

                    return forkJoin({
                        application:
                            this.admissionsService.getApplicantApplication(
                                applicationID
                            ),
                        checklist:
                            this.admissionsService.queryApplicantChecklistItems(
                                applicationID,
                                {
                                    rows: 100,
                                    orderBy: 'display_order,ASC',
                                }
                            ),
                        documents:
                            this.admissionsService.queryApplicantDocuments(
                                applicationID,
                                {
                                    rows: 100,
                                    orderBy: 'uploaded_at,DESC',
                                }
                            ),
                        programs: this.admissionsService.queryApplicantPrograms(
                            {
                                rows: 100,
                                active: true,
                            }
                        ),
                        terms: this.admissionsService.queryApplicantAcademicTerms(
                            {
                                rows: 100,
                                active: true,
                            }
                        ),
                    });
                }),
                finalize(() => this.isLoading.set(false)),
                takeUntilDestroyed(this.destroyRef)
            )
            .subscribe({
                next: ({
                    application,
                    checklist,
                    documents,
                    programs,
                    terms,
                }) => {
                    this.application.set(application);
                    this.programs.set(programs.items);
                    this.terms.set(terms.items);
                    this.documents.set(
                        this.mapPortalDocuments(
                            checklist.items,
                            documents.items
                        )
                    );
                    this.portalAuthService.updateApplicationID(application.id);
                },
                error: (error) => {
                    this.errorMessage.set(
                        jsonApiErrorMessage(
                            error,
                            error instanceof Error
                                ? error.message
                                : 'Admissions status could not be loaded.'
                        )
                    );
                },
            });
    }

    private resolveApplicationID(): Observable<string | null> {
        const sessionApplicationID = this.portalSession()?.applicationID;
        if (sessionApplicationID) {
            return of(sessionApplicationID);
        }

        return this.admissionsService
            .queryApplicantApplications({
                rows: 1,
                orderBy: 'date_updated,DESC',
            })
            .pipe(
                map((result) => {
                    const applicationID = result.items[0]?.id ?? null;
                    if (applicationID) {
                        this.portalAuthService.updateApplicationID(
                            applicationID
                        );
                    }

                    return applicationID;
                })
            );
    }

    private mapPortalDocuments(
        checklistItems: ChecklistItem[],
        documents: AdmissionsDocument[]
    ): PortalDocument[] {
        return [...checklistItems]
            .sort((left, right) => left.displayOrder - right.displayOrder)
            .map((item) => {
                const document = documents.find(
                    (candidate) => candidate.checklistItemID === item.id
                );

                return {
                    id: item.id,
                    checklistItemID: item.id,
                    label: item.documentName,
                    description: item.description ?? item.itemKey,
                    status: this.portalStatusFor(
                        document?.status ?? item.status
                    ),
                    required: item.required,
                    serverId: document?.storageKey,
                    fileName: document?.fileName,
                };
            });
    }

    private mergeCreatedDocument(document: AdmissionsDocument): void {
        this.documents.update((documents) =>
            documents.map((item) =>
                item.checklistItemID === document.checklistItemID
                    ? {
                          ...item,
                          status: this.portalStatusFor(document.status),
                          serverId: document.storageKey,
                          fileName: document.fileName,
                      }
                    : item
            )
        );
    }

    private updateDocument(
        documentId: string,
        changes: Partial<PortalDocument>
    ): void {
        this.documents.update((documents) =>
            documents.map((document) =>
                document.id === documentId
                    ? { ...document, ...changes }
                    : document
            )
        );
    }

    private portalStatusFor(
        status: AdmissionsDocument['status']
    ): PortalDocumentStatus {
        switch (status) {
            case 'ACCEPTED':
            case 'SYNCED_TO_SIS':
                return 'verified';
            case 'UPLOADED':
                return 'received';
            case 'PENDING_REVIEW':
                return 'pending-review';
            case 'REJECTED':
            case 'EXPIRED':
                return 'rejected';
            case 'WAIVED':
                return 'verified';
        }
    }

    private feeStatusFor(
        status: ApplicationStatus | undefined
    ): ApplicationFeeStatus {
        if (!status || status === 'DRAFT' || status === 'WITHDRAWN') {
            return 'NOT_REQUIRED';
        }

        return 'PENDING';
    }

    private statusIndex(status: ApplicationStatus): number {
        const stageIndex = TIMELINE_STAGES.findIndex(
            (stage) => stage.status === status
        );

        if (stageIndex >= 0) {
            return stageIndex;
        }

        if (
            [
                'ADMITTED',
                'DENIED',
                'WAITLISTED',
                'DEFERRED',
                'ENROLLED',
            ].includes(status)
        ) {
            return TIMELINE_STAGES.length;
        }

        return -1;
    }

    private timelineDate(
        application: Application,
        stage: ApplicationStatus,
        done: boolean
    ): string {
        if (!done) {
            return 'Pending';
        }

        if (stage === 'SUBMITTED' && application.submittedAt) {
            return this.formatDate(application.submittedAt);
        }

        return this.formatDate(
            application.dateUpdated || application.dateCreated
        );
    }

    private formatDate(value: string): string {
        return new Intl.DateTimeFormat('en-KE', {
            month: 'short',
            day: 'numeric',
            year: 'numeric',
        }).format(new Date(value));
    }
}
