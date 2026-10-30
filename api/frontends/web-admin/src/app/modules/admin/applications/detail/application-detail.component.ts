import {
    ChangeDetectionStrategy,
    Component,
    DestroyRef,
    computed,
    inject,
    signal,
} from '@angular/core';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatTabsModule } from '@angular/material/tabs';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { AdmissionsService } from 'app/core/admissions/admissions.service';
import {
    AcademicTerm,
    Application,
    ApplicationFee,
    ApplicationFeeStatus,
    ApplicationKCSEResult,
    ApplicationStatus,
    ApplicationTransition,
    KUCCPSPlacement,
    NotificationChannel,
    NotificationPreferences,
    Program,
} from 'app/core/admissions/admissions.types';
import { jsonApiErrorMessage } from 'app/core/api/json-api';
import { finalize, forkJoin } from 'rxjs';
import { ApplicationDocumentsComponent } from '../components/documents/application-documents.component';

interface TimelineEvent {
    icon: string;
    title: string;
    date: string;
}

interface Note {
    authorInitials: string;
    authorName: string;
    timestamp: string;
    body: string;
}

interface Insight {
    title: string;
    description: string;
}

interface ContactPreferenceRow {
    channel: NotificationChannel;
    label: string;
    destination: string;
    optIn: boolean;
    provider: string;
    consentNote: string;
}

const EMPTY_KUCCPS_PLACEMENT: KUCCPSPlacement = {
    placementID: 'Not captured',
    institutionCode: '—',
    programmeCode: '—',
    programmeName: 'No KUCCPS placement attached',
    placementYear: new Date().getFullYear(),
    weightedPointsNote:
        'KUCCPS placement details will appear here when the admissions API returns them.',
};

const EMPTY_KCSE_RESULT: ApplicationKCSEResult = {
    indexNumber: 'Not captured',
    examYear: new Date().getFullYear(),
    meanGrade: '—',
    meanPoints: 0,
    subjects: [],
};

@Component({
    selector: 'app-application-detail',
    imports: [
        RouterLink,
        MatButtonModule,
        MatIconModule,
        MatTabsModule,
        ApplicationDocumentsComponent,
    ],
    changeDetection: ChangeDetectionStrategy.OnPush,
    templateUrl: './application-detail.component.html',
})
export class ApplicationDetailComponent {
    private readonly admissionsService = inject(AdmissionsService);
    private readonly destroyRef = inject(DestroyRef);
    private readonly route = inject(ActivatedRoute);

    readonly applicationId = this.route.snapshot.paramMap.get('id') ?? '';
    readonly application = signal<Application | null>(null);
    readonly errorMessage = signal<string | null>(null);
    readonly isLoading = signal(false);
    readonly programs = signal<Program[]>([]);
    readonly terms = signal<AcademicTerm[]>([]);
    readonly transitions = signal<ApplicationTransition[]>([]);

    readonly applicantName = computed(
        () => this.application()?.constituentID ?? 'Applicant pending'
    );
    readonly status = computed(() =>
        this.formatApplicationStatus(this.application()?.status)
    );
    readonly applicationType = computed(() =>
        this.application()?.applicationType
            ? this.formatApplicationType(this.application()!.applicationType)
            : 'Application type pending'
    );
    readonly kuccpsPlacement = computed(
        () => this.application()?.kuccpsPlacement ?? EMPTY_KUCCPS_PLACEMENT
    );
    readonly kcseResult = computed(
        () => this.application()?.kcseResult ?? EMPTY_KCSE_RESULT
    );
    readonly program = computed(() =>
        this.programs().find(
            (program) => program.id === this.application()?.programID
        )
    );
    readonly term = computed(() =>
        this.terms().find(
            (term) => term.id === this.application()?.academicTermID
        )
    );
    readonly fee = computed<ApplicationFee>(() => {
        const application = this.application();
        const status = this.feeStatusFor(application?.status);

        return {
            id: `fee-${application?.id ?? this.applicationId}`,
            applicationID: application?.id ?? this.applicationId,
            amountCents: status === 'NOT_REQUIRED' ? 0 : 15000,
            currency: 'KES',
            status,
            provider: status === 'NOT_REQUIRED' ? 'not_required' : 'manual',
            dueAt: application?.submittedAt,
            auditTrail: this.transitions().map((transition) => ({
                actor: transition.actorID,
                action: `${this.formatApplicationStatus(transition.fromStatus)} → ${this.formatApplicationStatus(transition.toStatus)}`,
                reason:
                    transition.reason ?? transition.note ?? 'Status updated',
                timestamp: this.formatDate(transition.dateCreated),
            })),
        };
    });
    readonly summaryCards = computed(() => [
        {
            label: 'Programme',
            value: this.program()?.name ?? 'Programme pending',
        },
        { label: 'Application type', value: this.applicationType() },
        { label: 'Term', value: this.term()?.name ?? 'Term pending' },
        {
            label: 'Submitted',
            value: this.application()?.submittedAt
                ? this.formatDate(this.application()!.submittedAt!)
                : 'Not submitted',
        },
        {
            label: 'Application Fee',
            value: `${this.formatFeeStatus(this.fee().status)} · ${this.formatAmount(this.fee())}`,
        },
    ]);
    readonly applicantInfo = computed(() => {
        const application = this.application();
        return [
            {
                label: 'Constituent ID',
                value: application?.constituentID ?? '—',
            },
            {
                label: 'Application ID',
                value: application?.id ?? this.applicationId,
            },
            { label: 'Programme ID', value: application?.programID ?? '—' },
            {
                label: 'Academic term ID',
                value: application?.academicTermID ?? '—',
            },
            {
                label: 'Reviewer',
                value: application?.assignedReviewerID ?? 'Unassigned',
            },
            {
                label: 'Created',
                value: application?.dateCreated
                    ? this.formatDate(application.dateCreated)
                    : '—',
            },
            {
                label: 'Updated',
                value: application?.dateUpdated
                    ? this.formatDate(application.dateUpdated)
                    : '—',
            },
            {
                label: 'Data source',
                value: 'Admissions application API',
            },
        ];
    });
    readonly notificationPreferences: NotificationPreferences = {
        smsOptIn: true,
        whatsAppOptIn: false,
        emailOptIn: true,
        priority: ['SMS', 'WHATSAPP', 'EMAIL'],
    };
    readonly contactPreferenceRows = computed<ContactPreferenceRow[]>(() => [
        {
            channel: 'SMS',
            label: 'SMS',
            destination: 'Applicant phone from constituent profile',
            optIn: this.notificationPreferences.smsOptIn,
            provider: 'Celcom Africa SMS',
            consentNote: 'Constituent preference API is not embedded yet.',
        },
        {
            channel: 'WHATSAPP',
            label: 'WhatsApp',
            destination: 'Applicant WhatsApp from constituent profile',
            optIn: this.notificationPreferences.whatsAppOptIn,
            provider: 'WhatsApp Cloud',
            consentNote: 'Off until the applicant explicitly opts in.',
        },
        {
            channel: 'EMAIL',
            label: 'Email',
            destination: 'Applicant email from constituent profile',
            optIn: this.notificationPreferences.emailOptIn,
            provider: 'SchoolCRM mailer',
            consentNote: 'Kept on for long-form admissions notices.',
        },
    ]);
    readonly notes = computed<Note[]>(() => {
        const notes = this.transitions()
            .filter((transition) => transition.note || transition.reason)
            .map((transition) => ({
                authorInitials: this.initialsFor(transition.actorID),
                authorName: transition.actorID,
                timestamp: this.formatDate(transition.dateCreated),
                body:
                    transition.note ??
                    transition.reason ??
                    'Status transition recorded.',
            }));

        if (notes.length > 0) {
            return notes;
        }

        return [
            {
                authorInitials: 'API',
                authorName: 'Admissions API',
                timestamp: 'Live data',
                body: 'No reviewer notes have been recorded in application transitions yet.',
            },
        ];
    });
    readonly timelineEvents = computed<TimelineEvent[]>(() => {
        const transitions = this.transitions();
        if (transitions.length > 0) {
            return transitions.map((transition) => ({
                icon: 'heroicons_outline:arrow-path',
                title: `${this.formatApplicationStatus(transition.fromStatus)} → ${this.formatApplicationStatus(transition.toStatus)}`,
                date: this.formatDate(transition.dateCreated),
            }));
        }

        const application = this.application();
        if (!application) {
            return [];
        }

        return [
            {
                icon: 'heroicons_outline:paper-airplane',
                title: this.formatApplicationStatus(application.status),
                date: this.formatDate(application.dateUpdated),
            },
        ];
    });
    readonly insights = computed<Insight[]>(() => {
        const application = this.application();
        const kcse = application?.kcseResult;

        return [
            {
                title: 'Recommended next action',
                description: application
                    ? `Application is ${this.formatApplicationStatus(application.status).toLowerCase()}; use transitions and document reviews to move it forward.`
                    : 'Load the application to calculate next action.',
            },
            {
                title: 'KCSE signal',
                description: kcse
                    ? `${kcse.meanGrade} mean grade with ${kcse.meanPoints} aggregate points.`
                    : 'KCSE result has not been captured on this application.',
            },
            {
                title: 'Reviewer assignment',
                description:
                    application?.assignedReviewerID ??
                    'No reviewer assigned yet.',
            },
        ];
    });

    constructor() {
        this.loadApplication();
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
            return 'No fee';
        }

        return new Intl.NumberFormat('en-KE', {
            style: 'currency',
            currency: fee.currency,
        }).format(fee.amountCents / 100);
    }

    preferenceIcon(channel: NotificationChannel): string {
        const icons: Record<NotificationChannel, string> = {
            SMS: 'heroicons_outline:chat-bubble-left',
            WHATSAPP: 'heroicons_outline:device-phone-mobile',
            EMAIL: 'heroicons_outline:envelope',
        };

        return icons[channel];
    }

    private loadApplication(): void {
        if (!this.applicationId) {
            this.errorMessage.set('No application ID was supplied.');
            return;
        }

        this.isLoading.set(true);
        this.errorMessage.set(null);

        forkJoin({
            application: this.admissionsService.getApplication(
                this.applicationId
            ),
            transitions: this.admissionsService.queryApplicationTransitions(
                this.applicationId,
                {
                    page: 1,
                    rows: 100,
                    orderBy: 'date_created,DESC',
                }
            ),
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
                next: ({ application, transitions, programs, terms }) => {
                    this.application.set(application);
                    this.transitions.set(transitions.items);
                    this.programs.set(programs.items);
                    this.terms.set(terms.items);
                },
                error: (error) => {
                    this.errorMessage.set(
                        jsonApiErrorMessage(
                            error,
                            'Application detail could not be loaded.'
                        )
                    );
                },
            });
    }

    private feeStatusFor(
        status: ApplicationStatus | undefined
    ): ApplicationFeeStatus {
        if (!status || status === 'DRAFT' || status === 'WITHDRAWN') {
            return 'NOT_REQUIRED';
        }

        return 'PENDING';
    }

    private formatApplicationStatus(
        status: ApplicationStatus | undefined
    ): string {
        return status ? status.replaceAll('_', ' ') : 'Status pending';
    }

    private formatApplicationType(type: string): string {
        return type.replaceAll('_', ' ');
    }

    private formatDate(value: string): string {
        return new Intl.DateTimeFormat('en-KE', {
            month: 'short',
            day: 'numeric',
            year: 'numeric',
        }).format(new Date(value));
    }

    private initialsFor(value: string): string {
        const clean = value.replaceAll('-', ' ').trim();
        const parts = clean.split(/\s+/).filter(Boolean);

        return (parts[0]?.[0] ?? 'A') + (parts[1]?.[0] ?? 'P');
    }
}
