import {
    ChangeDetectionStrategy,
    Component,
    DestroyRef,
    computed,
    inject,
    signal,
} from '@angular/core';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { RouterLink } from '@angular/router';
import { AdmissionsService } from 'app/core/admissions/admissions.service';
import {
    AcademicTerm,
    Application,
    ApplicationFee,
    ApplicationFeeStatus,
    ApplicationKCSEResult,
    ApplicationStatus,
    ApplicationType,
    Program,
} from 'app/core/admissions/admissions.types';
import { jsonApiErrorMessage } from 'app/core/api/json-api';
import { finalize, forkJoin } from 'rxjs';

interface ApplicationRow {
    id: string;
    applicant: string;
    program: string;
    programmeCode: string;
    term: string;
    applicationType: ApplicationType;
    status: ApplicationStatus;
    progress: number;
    reviewer: string;
    score: number | null;
    kcseResult?: ApplicationKCSEResult;
    fee: ApplicationFee;
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

@Component({
    selector: 'app-applications',
    imports: [FormsModule, RouterLink, MatButtonModule, MatIconModule],
    changeDetection: ChangeDetectionStrategy.OnPush,
    templateUrl: './applications.component.html',
})
export class ApplicationsComponent {
    private readonly admissionsService = inject(AdmissionsService);
    private readonly destroyRef = inject(DestroyRef);

    readonly filters = [
        'Status',
        'Programme',
        'Application type',
        'KCSE mean grade',
        'Reviewer',
        'Fee status',
    ];

    readonly applications = signal<Application[]>([]);
    readonly errorMessage = signal<string | null>(null);
    readonly isLoading = signal(false);
    readonly programs = signal<Program[]>([]);
    readonly searchTerm = signal('');
    readonly terms = signal<AcademicTerm[]>([]);
    readonly totalApplications = signal(0);

    readonly rows = computed(() => {
        const search = this.searchTerm().trim().toLowerCase();
        const rows = this.applications().map((application) =>
            this.toRow(application)
        );

        if (!search) {
            return rows;
        }

        return rows.filter((row) =>
            [
                row.id,
                row.applicant,
                row.program,
                row.programmeCode,
                row.term,
                row.status,
                row.applicationType,
            ]
                .join(' ')
                .toLowerCase()
                .includes(search)
        );
    });

    constructor() {
        this.loadApplications();
    }

    updateSearchTerm(value: string): void {
        this.searchTerm.set(value);
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

    formatApplicationStatus(status: ApplicationStatus): string {
        return status.replaceAll('_', ' ');
    }

    formatApplicationType(type: ApplicationType): string {
        return type.replaceAll('_', ' ');
    }

    formatAmount(fee: ApplicationFee): string {
        if (fee.status === 'NOT_REQUIRED') {
            return 'No fee';
        }

        return new Intl.NumberFormat('en-KE', {
            style: 'currency',
            currency: fee.currency,
        }).format(fee.amountCents / 100);
    }

    private loadApplications(): void {
        this.isLoading.set(true);
        this.errorMessage.set(null);

        forkJoin({
            applications: this.admissionsService.queryApplications({
                page: 1,
                rows: 50,
                orderBy: 'date_updated,DESC',
            }),
            programs: this.admissionsService.queryPrograms({
                rows: 100,
                active: true,
            }),
            terms: this.admissionsService.queryAcademicTerms({
                rows: 100,
                active: true,
            }),
        })
            .pipe(
                finalize(() => this.isLoading.set(false)),
                takeUntilDestroyed(this.destroyRef)
            )
            .subscribe({
                next: ({ applications, programs, terms }) => {
                    this.applications.set(applications.items);
                    this.programs.set(programs.items);
                    this.terms.set(terms.items);
                    this.totalApplications.set(applications.total);
                },
                error: (error) => {
                    this.errorMessage.set(
                        jsonApiErrorMessage(
                            error,
                            'Applications could not be loaded.'
                        )
                    );
                },
            });
    }

    private toRow(application: Application): ApplicationRow {
        const program = this.programs().find(
            (item) => item.id === application.programID
        );
        const term = this.terms().find(
            (item) => item.id === application.academicTermID
        );

        return {
            id: application.id,
            applicant: application.constituentID,
            program: program?.name ?? 'Programme pending',
            programmeCode: program?.code ?? application.programID,
            term: term?.name ?? application.academicTermID,
            applicationType: application.applicationType,
            status: application.status,
            progress: STATUS_PROGRESS[application.status],
            reviewer: application.assignedReviewerID ?? 'Unassigned',
            score: application.kcseResult?.meanPoints ?? null,
            kcseResult: application.kcseResult,
            fee: this.feeFor(application),
        };
    }

    private feeFor(application: Application): ApplicationFee {
        const status = this.feeStatusFor(application.status);

        return {
            id: `fee-${application.id}`,
            applicationID: application.id,
            amountCents: status === 'NOT_REQUIRED' ? 0 : 15000,
            currency: 'KES',
            status,
            provider: status === 'NOT_REQUIRED' ? 'not_required' : 'manual',
            dueAt: application.submittedAt,
            auditTrail: [],
        };
    }

    private feeStatusFor(status: ApplicationStatus): ApplicationFeeStatus {
        if (status === 'DRAFT' || status === 'WITHDRAWN') {
            return 'NOT_REQUIRED';
        }

        return 'PENDING';
    }
}
