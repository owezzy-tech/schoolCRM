import { AsyncPipe } from '@angular/common';
import { ChangeDetectionStrategy, Component, inject } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { AdmissionsService } from 'app/core/admissions/admissions.service';
import {
    AdmissionsDocument,
    AdmissionsDocumentVerificationRequest,
    ChecklistItem,
    DOCUMENT_STATUSES,
    DocumentStatus,
} from 'app/core/admissions/admissions.types';
import { jsonApiErrorMessage } from 'app/core/api/json-api';
import { catchError, forkJoin, of } from 'rxjs';

interface DocumentReviewForm {
    documentID: string;
    status: AdmissionsDocumentVerificationRequest['status'];
    reviewerNotes: string;
}

interface DocumentUploadForm {
    checklistItemID: string;
    fileName: string;
    contentType: string;
    sizeBytes: number;
    storageKey: string;
}

const emptyReviewForm: DocumentReviewForm = {
    documentID: '',
    status: 'ACCEPTED',
    reviewerNotes: '',
};

const emptyUploadForm: DocumentUploadForm = {
    checklistItemID: '',
    fileName: '',
    contentType: 'application/pdf',
    sizeBytes: 1,
    storageKey: '',
};

@Component({
    selector: 'app-admissions-documents',
    imports: [
        AsyncPipe,
        FormsModule,
        MatButtonModule,
        MatFormFieldModule,
        MatIconModule,
        MatInputModule,
        MatSelectModule,
    ],
    templateUrl: './admissions-documents.component.html',
    changeDetection: ChangeDetectionStrategy.OnPush,
})
export class AdmissionsDocumentsComponent {
    private readonly admissionsService = inject(AdmissionsService);

    readonly checklistItems$ = this.admissionsService.checklistItems$;
    readonly documents$ = this.admissionsService.documents$;
    readonly documentStatuses = DOCUMENT_STATUSES;
    readonly reviewStatuses: DocumentReviewForm['status'][] = [
        'ACCEPTED',
        'REJECTED',
        'WAIVED',
    ];

    applicationID = '';
    errorMessage = '';
    saving = false;
    reviewForm: DocumentReviewForm = { ...emptyReviewForm };
    uploadForm: DocumentUploadForm = { ...emptyUploadForm };

    loadApplicationDocuments(): void {
        if (!this.applicationID.trim()) {
            this.errorMessage = 'Enter an application ID to load documents.';
            return;
        }

        this.errorMessage = '';
        forkJoin([
            this.admissionsService.queryChecklistItems(this.applicationID, {
                page: 1,
                rows: 100,
                orderBy: 'display_order,ASC',
            }),
            this.admissionsService.queryDocuments(this.applicationID, {
                page: 1,
                rows: 100,
                orderBy: 'uploaded_at,DESC',
            }),
        ])
            .pipe(
                catchError((error) => {
                    this.errorMessage = jsonApiErrorMessage(
                        error,
                        'Unable to load document checklist.'
                    );
                    return of(undefined);
                })
            )
            .subscribe();
    }

    selectChecklistItem(item: ChecklistItem): void {
        this.uploadForm = {
            ...emptyUploadForm,
            checklistItemID: item.id,
            fileName: `${item.itemKey}.pdf`,
            storageKey: `admissions/applications/${item.applicationID}/${item.itemKey}.pdf`,
        };
    }

    selectDocument(document: AdmissionsDocument): void {
        this.reviewForm = {
            documentID: document.id,
            status: this.toReviewStatus(document.status),
            reviewerNotes: document.reviewerNotes ?? '',
        };
    }

    saveUpload(): void {
        if (!this.applicationID.trim() || !this.uploadForm.checklistItemID) {
            this.errorMessage =
                'Select a checklist item before saving metadata.';
            return;
        }

        this.saving = true;
        this.errorMessage = '';
        this.admissionsService
            .createDocument(this.applicationID, {
                checklistItemID: this.uploadForm.checklistItemID,
                fileName: this.uploadForm.fileName,
                contentType: this.uploadForm.contentType,
                sizeBytes: Number(this.uploadForm.sizeBytes),
                storageKey: this.uploadForm.storageKey,
            })
            .pipe(
                catchError((error) => {
                    this.errorMessage = jsonApiErrorMessage(
                        error,
                        'Unable to save document metadata.'
                    );
                    return of(undefined);
                })
            )
            .subscribe((document) => {
                this.saving = false;
                if (document) {
                    this.uploadForm = { ...emptyUploadForm };
                    this.loadApplicationDocuments();
                }
            });
    }

    saveReview(): void {
        if (!this.reviewForm.documentID) {
            this.errorMessage = 'Select a document before saving review.';
            return;
        }

        this.saving = true;
        this.errorMessage = '';
        this.admissionsService
            .verifyDocument(this.reviewForm.documentID, {
                status: this.reviewForm.status,
                reviewerNotes: this.reviewForm.reviewerNotes || null,
            })
            .pipe(
                catchError((error) => {
                    this.errorMessage = jsonApiErrorMessage(
                        error,
                        'Unable to save document review.'
                    );
                    return of(undefined);
                })
            )
            .subscribe((document) => {
                this.saving = false;
                if (document) {
                    this.reviewForm = { ...emptyReviewForm };
                    this.loadApplicationDocuments();
                }
            });
    }

    requestDownload(document: AdmissionsDocument): void {
        this.errorMessage = '';
        this.admissionsService
            .downloadDocument(document.id)
            .pipe(
                catchError((error) => {
                    this.errorMessage = jsonApiErrorMessage(
                        error,
                        'Unable to audit document download.'
                    );
                    return of(undefined);
                })
            )
            .subscribe();
    }

    statusClass(status: DocumentStatus): string {
        switch (status) {
            case 'ACCEPTED':
                return 'bg-green-100 text-green-700';
            case 'REJECTED':
                return 'bg-red-100 text-red-700';
            case 'WAIVED':
                return 'bg-purple-100 text-purple-700';
            case 'PENDING_REVIEW':
            case 'UPLOADED':
                return 'bg-amber-100 text-amber-700';
            default:
                return 'bg-slate-100 text-secondary';
        }
    }

    private toReviewStatus(
        status: DocumentStatus
    ): AdmissionsDocumentVerificationRequest['status'] {
        if (status === 'REJECTED' || status === 'WAIVED') {
            return status;
        }

        return 'ACCEPTED';
    }
}
