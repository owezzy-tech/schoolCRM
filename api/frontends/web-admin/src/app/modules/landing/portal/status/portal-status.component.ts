import { TitleCasePipe } from '@angular/common';
import { ChangeDetectionStrategy, Component, inject } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { RouterLink } from '@angular/router';
import { PortalAuthService } from 'app/core/portal/portal-auth.service';
import { FilePondComponent } from 'app/shared/components/file-upload/file-pond.component';
import { FilePondFile } from 'filepond';

type PortalDocumentStatus =
    | 'verified'
    | 'received'
    | 'missing'
    | 'uploading'
    | 'pending-review';

interface PortalDocument {
    id: string;
    label: string;
    description: string;
    status: PortalDocumentStatus;
    required: boolean;
    serverId?: string;
    fileName?: string;
}

@Component({
    selector: 'app-portal-status',
    standalone: true,
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
                            {{
                                portalSession()?.applicationID ?? 'APP-3018'
                            }}
                            &middot; Data Science MSc &middot; Fall 2026
                        </h2>
                    </div>
                    <div
                        class="inline-flex items-center rounded-full bg-yellow-100 px-3 py-1 text-sm font-medium text-yellow-900"
                    >
                        Under Review
                    </div>
                </div>

                <!-- Overall progress card -->
                <div
                    class="bg-card mt-6 flex flex-col rounded-2xl border p-6 shadow-sm"
                >
                    <h3 class="text-default mb-4 text-lg font-semibold">
                        Overall progress
                    </h3>
                    <mat-progress-bar
                        mode="determinate"
                        [value]="62"
                        class="mb-2"
                    ></mat-progress-bar>
                    <p class="text-secondary text-sm font-medium">
                        62% complete &middot; 3 of 5 stages done.
                    </p>
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
                            <!-- Vertical dashed line -->
                            <div
                                class="absolute bottom-4 left-[11px] top-4 border-l-2 border-dashed border-gray-300"
                            ></div>

                            @for (item of timeline; track item.label) {
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
                                    to the mock upload service and move into
                                    pending review.
                                </p>
                            </div>
                            <div
                                class="rounded-full bg-primary-50 px-3 py-1 text-xs font-semibold text-primary dark:bg-primary-900/20"
                            >
                                {{ uploadedDocumentCount }} /
                                {{ documents.length }} received
                            </div>
                        </div>

                        <div class="mt-4 flex flex-col divide-y">
                            @for (doc of documents; track doc.label) {
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
                                            [class.bg-green-100]="
                                                doc.status === 'verified'
                                            "
                                            [class.text-green-800]="
                                                doc.status === 'verified'
                                            "
                                            [class.bg-blue-100]="
                                                doc.status === 'received'
                                            "
                                            [class.text-blue-800]="
                                                doc.status === 'received'
                                            "
                                            [class.bg-red-100]="
                                                doc.status === 'missing'
                                            "
                                            [class.text-red-800]="
                                                doc.status === 'missing'
                                            "
                                            [class.bg-amber-100]="
                                                doc.status === 'uploading' ||
                                                doc.status === 'pending-review'
                                            "
                                            [class.text-amber-800]="
                                                doc.status === 'uploading' ||
                                                doc.status === 'pending-review'
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

                <!-- Counselor card -->
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
                            Reach your counselor
                        </h3>
                        <p class="text-secondary mt-0.5">
                            Maya Schultz &middot; maya&#64;northbrook.edu
                            &middot; Replies within 1 business day.
                        </p>
                    </div>
                </div>
            </div>
        </div>
    `,
    changeDetection: ChangeDetectionStrategy.OnPush,
    host: { class: 'flex w-full flex-auto flex-col' },
})
export class PortalStatusComponent {
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

    readonly timeline = [
        { label: 'Submitted', date: 'May 12, 2026', done: true },
        { label: 'Documents received', date: 'May 14, 2026', done: true },
        { label: 'Under review', date: 'May 18, 2026', done: true },
        { label: 'Interview', date: 'Scheduled for Jun 2', done: false },
        { label: 'Decision', date: 'Expected late June', done: false },
    ];

    documents: PortalDocument[] = [
        {
            id: 'personal-statement',
            label: 'Personal statement',
            description: 'PDF or image copy of your admissions essay.',
            status: 'verified',
            required: true,
            fileName: 'personal-statement.pdf',
        },
        {
            id: 'transcripts',
            label: 'Transcripts',
            description:
                'Official undergraduate transcript from your institution.',
            status: 'verified',
            required: true,
            fileName: 'northbrook-transcript.pdf',
        },
        {
            id: 'recommendation-1',
            label: 'Recommendation letter 1',
            description:
                'First recommender letter received and awaiting review.',
            status: 'received',
            required: true,
            fileName: 'recommendation-letter-1.pdf',
        },
        {
            id: 'recommendation-2',
            label: 'Recommendation letter 2',
            description:
                'Upload the second recommender letter when it is ready.',
            status: 'missing',
            required: true,
        },
        {
            id: 'test-scores',
            label: 'Test scores',
            description: 'Optional GRE/GMAT score report.',
            status: 'verified',
            required: false,
            fileName: 'gre-score-report.pdf',
        },
    ];

    get uploadedDocumentCount(): number {
        return this.documents.filter(
            (document) => document.status !== 'missing'
        ).length;
    }

    canUploadDocument(document: PortalDocument): boolean {
        return (
            this.hasPortalAccess() &&
            ['missing', 'received', 'uploading'].includes(document.status)
        );
    }

    hasPortalAccess(): boolean {
        return this.portalAuthService.hasValidSession();
    }

    uploadLabelFor(document: PortalDocument): string {
        if (document.status === 'received') {
            return 'Reupload replacement or <span class="filepond--label-action">browse</span>';
        }

        return `Drop ${document.label.toLowerCase()} or <span class="filepond--label-action">browse</span>`;
    }

    markUploadStarted(documentId: string): void {
        this.updateDocument(documentId, { status: 'uploading' });
    }

    markUploaded(documentId: string, file: FilePondFile): void {
        this.updateDocument(documentId, {
            status: 'pending-review',
            serverId: file.serverId,
            fileName: file.filename,
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

    private updateDocument(
        documentId: string,
        changes: Partial<PortalDocument>
    ): void {
        this.documents = this.documents.map((document) =>
            document.id === documentId ? { ...document, ...changes } : document
        );
    }
}
