import {
    ChangeDetectionStrategy,
    Component,
    OnInit,
    computed,
    inject,
    signal,
} from '@angular/core';
import {
    FormsModule,
    ReactiveFormsModule,
    UntypedFormBuilder,
    UntypedFormGroup,
} from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { jsonApiErrorMessage } from 'app/core/api/json-api';
import { AdmissionsService } from 'app/core/admissions/admissions.service';
import {
    AcademicTerm,
    Application,
    ApplicationFee,
    ApplicationFeeStatus,
    ApplicationFormTemplate,
    ApplicationKCSEResult,
    ApplicationRequest,
    ApplicationType,
    KUCCPSPlacement,
    Program,
} from 'app/core/admissions/admissions.types';
import { PortalAuthService } from 'app/core/portal/portal-auth.service';
import { FilePondComponent } from 'app/shared/components/file-upload/file-pond.component';
import { forkJoin } from 'rxjs';

interface PortalApplyStep {
    number: number;
    label: string;
    heading: string;
    description: string;
}

interface DocumentUploadRequirement {
    label: string;
    description: string;
    labelIdle: string;
}

@Component({
    selector: 'app-portal-apply',
    standalone: true,
    imports: [
        FormsModule,
        ReactiveFormsModule,
        MatButtonModule,
        MatFormFieldModule,
        MatIconModule,
        MatInputModule,
        MatProgressBarModule,
        FilePondComponent,
    ],
    template: `
        <div class="flex flex-auto flex-col px-6 py-12 sm:px-10 lg:px-16">
            <div class="mx-auto w-full max-w-6xl">
                <div class="mb-8 flex flex-col">
                    <span
                        class="text-secondary mb-2 text-sm font-semibold uppercase tracking-wider"
                        >Step {{ currentStep().number }} of {{ steps.length }} -
                        {{ currentStep().heading }}</span
                    >
                    <h1
                        class="text-default text-3xl font-bold tracking-tight sm:text-4xl"
                    >
                        Apply for the 2026 main intake
                    </h1>
                </div>

                <div class="mb-10 w-full">
                    <mat-progress-bar
                        mode="determinate"
                        [value]="progressValue()"
                        [attr.aria-label]="
                            'Application progress: step ' +
                            currentStep().number +
                            ' of ' +
                            steps.length +
                            ', ' +
                            currentStep().label.toLowerCase()
                        "
                    ></mat-progress-bar>
                </div>

                <div class="grid grid-cols-1 gap-8 lg:grid-cols-[1fr_320px]">
                    <!-- Left Column - Form Card -->
                    <div class="bg-card rounded-2xl border p-8 shadow-sm">
                        <form [formGroup]="form" class="flex flex-col">
                            <div class="mb-6 rounded-2xl bg-primary-50 p-5 dark:bg-primary-900/20">
                                <h2 class="text-default text-lg font-semibold">
                                    {{ currentStep().heading }}
                                </h2>
                                <p class="text-secondary mt-1 text-sm leading-relaxed">
                                    {{ currentStep().description }}
                                </p>
                            </div>

                            @switch (currentStep().number) {
                                @case (1) {
                                    <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
                                        <mat-form-field appearance="outline" class="w-full">
                                            <mat-label>First name</mat-label>
                                            <input matInput formControlName="firstName" autocomplete="given-name" />
                                        </mat-form-field>

                                        <mat-form-field appearance="outline" class="w-full">
                                            <mat-label>Last name</mat-label>
                                            <input matInput formControlName="lastName" autocomplete="family-name" />
                                        </mat-form-field>

                                        <mat-form-field appearance="outline" class="w-full">
                                            <mat-label>Email</mat-label>
                                            <input matInput type="email" formControlName="email" autocomplete="email" />
                                        </mat-form-field>

                                        <mat-form-field appearance="outline" class="w-full">
                                            <mat-label>Phone</mat-label>
                                            <input matInput type="tel" formControlName="phone" autocomplete="tel" />
                                        </mat-form-field>

                                        <mat-form-field appearance="outline" class="w-full">
                                            <mat-label>Date of birth</mat-label>
                                            <input matInput type="date" formControlName="dateOfBirth" autocomplete="bday" />
                                        </mat-form-field>

                                        <mat-form-field appearance="outline" class="w-full">
                                            <mat-label>Password</mat-label>
                                            <input matInput type="password" formControlName="password" autocomplete="new-password" aria-describedby="portal-password-help" />
                                        </mat-form-field>

                                        <mat-form-field appearance="outline" class="w-full">
                                            <mat-label>Confirm password</mat-label>
                                            <input matInput type="password" formControlName="confirmPassword" autocomplete="new-password" />
                                        </mat-form-field>
                                    </div>

                                    <p id="portal-password-help" class="text-secondary mt-2 text-sm">
                                        Use a memorable password for your applicant account. You will use it alongside the email on your KUCCPS or self-sponsored application.
                                    </p>
                                }

                                @case (2) {
                                    <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
                                        <mat-form-field appearance="outline" class="w-full">
                                            <mat-label>Preferred programme</mat-label>
                                            <input matInput formControlName="programme" list="portal-programme-options" placeholder="Bachelor of Commerce" />
                                            <datalist id="portal-programme-options">
                                                @for (program of programs(); track program.id) {
                                                    <option [value]="program.name">{{ program.code }}</option>
                                                }
                                            </datalist>
                                        </mat-form-field>

                                        <mat-form-field appearance="outline" class="w-full">
                                            <mat-label>Intake</mat-label>
                                            <input matInput formControlName="intake" list="portal-intake-options" placeholder="2026 Main Intake" />
                                            <datalist id="portal-intake-options">
                                                @for (term of academicTerms(); track term.id) {
                                                    <option [value]="term.name">{{ term.code }}</option>
                                                }
                                            </datalist>
                                        </mat-form-field>

                                        <mat-form-field appearance="outline" class="w-full">
                                            <mat-label>Sponsorship route</mat-label>
                                            <input matInput formControlName="sponsorshipRoute" placeholder="KUCCPS or self-sponsored" />
                                        </mat-form-field>

                                        <mat-form-field appearance="outline" class="w-full">
                                            <mat-label>Campus preference</mat-label>
                                            <input matInput formControlName="campusPreference" placeholder="Main campus" />
                                        </mat-form-field>
                                    </div>
                                }

                                @case (3) {
                                    <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
                                        <mat-form-field appearance="outline" class="w-full">
                                            <mat-label>KCSE index number</mat-label>
                                            <input matInput formControlName="kcseIndexNumber" placeholder="12345678901/2025" />
                                        </mat-form-field>

                                        <mat-form-field appearance="outline" class="w-full">
                                            <mat-label>KCSE year</mat-label>
                                            <input matInput formControlName="kcseYear" placeholder="2025" />
                                        </mat-form-field>

                                        <mat-form-field appearance="outline" class="w-full">
                                            <mat-label>KCSE mean grade</mat-label>
                                            <input matInput formControlName="kcseMeanGrade" placeholder="B+" />
                                        </mat-form-field>

                                        <mat-form-field appearance="outline" class="w-full">
                                            <mat-label>Subject highlights</mat-label>
                                            <input matInput formControlName="kcseSubjectHighlights" placeholder="Maths B+, English A-" />
                                        </mat-form-field>
                                    </div>
                                }

                                @case (4) {
                                    <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
                                        <mat-form-field appearance="outline" class="w-full">
                                            <mat-label>KUCCPS placement number</mat-label>
                                            <input matInput formControlName="kuccpsPlacementNumber" placeholder="KUCCPS-2026-0001" />
                                        </mat-form-field>

                                        <mat-form-field appearance="outline" class="w-full">
                                            <mat-label>Application route</mat-label>
                                            <input matInput formControlName="applicationRoute" placeholder="KUCCPS placement" />
                                        </mat-form-field>

                                        <mat-form-field appearance="outline" class="w-full">
                                            <mat-label>County</mat-label>
                                            <input matInput formControlName="county" placeholder="Nairobi" />
                                        </mat-form-field>

                                        <mat-form-field appearance="outline" class="w-full">
                                            <mat-label>National ID or passport</mat-label>
                                            <input matInput formControlName="identityNumber" placeholder="12345678" />
                                        </mat-form-field>
                                    </div>
                                }

                                @case (5) {
                                    <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
                                        @for (document of documentUploadRequirements; track document.label) {
                                            <div class="rounded-2xl border border-dashed bg-gray-50/70 p-4 dark:bg-gray-900/20">
                                                <div class="text-default font-semibold">
                                                    {{ document.label }}
                                                </div>
                                                <p class="text-secondary mb-3 mt-1 text-sm">
                                                    {{ document.description }}
                                                </p>
                                                <app-file-pond
                                                    [acceptedFileTypes]="acceptedDocumentTypes"
                                                    [maxFileSize]="'10MB'"
                                                    [server]="serverConfig"
                                                    [labelIdle]="document.labelIdle"
                                                />
                                            </div>
                                        }
                                    </div>
                                }

                                @case (6) {
                                    <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
                                        <mat-form-field appearance="outline" class="w-full">
                                            <mat-label>M-Pesa confirmation code</mat-label>
                                            <input matInput formControlName="mpesaConfirmationCode" placeholder="QF123ABC45" />
                                        </mat-form-field>

                                        <mat-form-field appearance="outline" class="w-full">
                                            <mat-label>Review notes</mat-label>
                                            <textarea matInput formControlName="reviewNotes" rows="4" placeholder="Anything admissions should know before review"></textarea>
                                        </mat-form-field>
                                    </div>

                                    <p class="text-secondary mt-2 text-sm">
                                        By submitting, you confirm the information is accurate and may be verified through KUCCPS, KNEC, IPRS, and finance office records.
                                    </p>
                                }
                            }

                            @if (apiError()) {
                                <div class="mt-6 rounded-2xl border border-red-200 bg-red-50 px-4 py-3 text-sm font-medium text-red-700" role="alert">
                                    {{ apiError() }}
                                </div>
                            }

                            @if (apiMessage()) {
                                <div class="mt-6 rounded-2xl border border-green-200 bg-green-50 px-4 py-3 text-sm font-medium text-green-700" role="status">
                                    {{ apiMessage() }}
                                </div>
                            }

                            <div class="mt-8 flex items-center justify-between">
                                <button
                                    mat-stroked-button
                                    type="button"
                                    [disabled]="isFirstStep()"
                                    (click)="goToPreviousStep()"
                                >
                                    <mat-icon
                                        svgIcon="heroicons_outline:arrow-left"
                                        class="icon-size-5"
                                        aria-hidden="true"
                                    ></mat-icon>
                                    <span class="ml-2">Back</span>
                                </button>
                                <button
                                    mat-flat-button
                                    color="primary"
                                    type="button"
                                    [disabled]="isSaving() || isLoadingOptions()"
                                    (click)="goToNextStep()"
                                >
                                    <span class="mr-2">{{ nextButtonLabel() }}</span>
                                    <mat-icon
                                        svgIcon="heroicons_outline:arrow-right"
                                        class="icon-size-5"
                                        aria-hidden="true"
                                    ></mat-icon>
                                </button>
                            </div>
                        </form>
                    </div>

                    <!-- Right Column - Sidebar -->
                    <div class="flex flex-col gap-6">
                        <!-- Progress Card -->
                        <div class="bg-card rounded-2xl border p-6 shadow-sm">
                            <h2
                                class="text-secondary mb-6 text-sm font-bold uppercase tracking-wider"
                            >
                                Progress
                            </h2>
                            <div class="flex flex-col gap-4">
                                @for (step of steps; track step.number) {
                                    <button
                                        type="button"
                                        class="flex items-center gap-4 rounded-xl p-2 text-left transition hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-primary dark:hover:bg-gray-800"
                                        [attr.aria-current]="
                                            step.number === currentStep().number
                                                ? 'step'
                                                : null
                                        "
                                        (click)="goToStep(step.number)"
                                    >
                                        <div
                                            class="flex h-8 w-8 items-center justify-center rounded-full text-sm font-bold"
                                            [class]="
                                                step.number === currentStep().number
                                                    ? 'bg-primary text-on-primary'
                                                    : 'text-secondary border border-gray-300'
                                            "
                                        >
                                            {{ step.number }}
                                        </div>
                                        <span
                                            class="text-base"
                                            [class]="
                                                step.number === currentStep().number
                                                    ? 'text-default font-bold'
                                                    : 'text-secondary font-medium'
                                            "
                                        >
                                            {{ step.label }}
                                        </span>
                                    </button>
                                }
                            </div>
                        </div>

                        <!-- Application Fee Card -->
                        <div class="bg-card rounded-2xl border p-6 shadow-sm">
                            <div class="flex items-start justify-between gap-3">
                                <div>
                                    <h2
                                        class="text-secondary text-sm font-bold uppercase tracking-wider"
                                    >
                                        Application fee
                                    </h2>
                                    <p class="text-secondary mt-2 text-sm">
                                        Fee status is tracked in Kenya shillings
                                        and reconciled through M-Pesa or finance
                                        office waivers.
                                    </p>
                                </div>
                                <span
                                    class="rounded-full px-2.5 py-1 text-xs font-semibold"
                                    [class]="feeStatusClass(reviewFee.status)"
                                    [attr.aria-label]="
                                        'Application fee status: ' +
                                        formatFeeStatus(reviewFee.status)
                                    "
                                >
                                    {{ formatFeeStatus(reviewFee.status) }}
                                </span>
                            </div>
                            <div
                                class="mt-5 rounded-2xl bg-gray-50 p-4 dark:bg-gray-800"
                            >
                                <div
                                    class="text-secondary text-xs font-semibold uppercase tracking-wide"
                                >
                                    Due at review
                                </div>
                                <div
                                    class="text-default mt-1 text-2xl font-bold"
                                >
                                    {{ formatAmount(reviewFee) }}
                                </div>
                                <div class="text-secondary mt-1 text-xs">
                                     Payment channel: {{ formatPaymentProvider(reviewFee.provider) }} ·
                                    status updates after reconciliation
                                 </div>
                             </div>
                         </div>

                        <!-- Need Help Card -->
                        <div
                            class="rounded-2xl bg-gray-50 p-6 dark:bg-gray-800"
                        >
                            <h2 class="text-default mb-2 text-lg font-semibold">
                                Need help?
                            </h2>
                            <p class="text-secondary leading-relaxed">
                                Email admissions&#64;schoolcrm.ac.ke or WhatsApp
                                +254 700 123 456 on weekdays 8am-5pm EAT.
                            </p>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    `,
    changeDetection: ChangeDetectionStrategy.OnPush,
    host: { class: 'flex w-full flex-auto flex-col' },
})
export class PortalApplyComponent implements OnInit {
    private readonly admissionsService = inject(AdmissionsService);
    private readonly portalAuthService = inject(PortalAuthService);
    private readonly formBuilder = inject(UntypedFormBuilder);

    readonly reviewFee: ApplicationFee = {
        id: 'fee-application-draft',
        applicationID: 'APP-DRAFT',
        amountCents: 15000,
        currency: 'KES',
        status: 'PENDING',
        provider: 'manual',
        dueAt: 'Before final submission',
        auditTrail: [],
    };

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

    readonly documentUploadRequirements: DocumentUploadRequirement[] = [
        {
            label: 'KCSE result document',
            description: 'Upload a PDF or image of your KCSE result slip or certificate.',
            labelIdle:
                'Drop KCSE result document or <span class="filepond--label-action">browse</span>',
        },
        {
            label: 'ID or passport document',
            description: 'Upload your Kenyan national ID, birth certificate, or passport.',
            labelIdle:
                'Drop ID or passport document or <span class="filepond--label-action">browse</span>',
        },
        {
            label: 'KUCCPS placement letter',
            description: 'Upload your placement letter if you are applying through KUCCPS.',
            labelIdle:
                'Drop KUCCPS placement letter or <span class="filepond--label-action">browse</span>',
        },
        {
            label: 'M-Pesa receipt',
            description: 'Upload an M-Pesa confirmation message or finance office receipt if available.',
            labelIdle:
                'Drop M-Pesa receipt or <span class="filepond--label-action">browse</span>',
        },
    ];

    readonly steps: PortalApplyStep[] = [
        {
            number: 1,
            label: 'Account',
            heading: 'Applicant account',
            description:
                'Create your secure applicant profile with the email and phone number admissions staff will use for Kenya intake updates.',
        },
        {
            number: 2,
            label: 'Programme',
            heading: 'Programme selection',
            description:
                'Choose your preferred programme, intake, and sponsorship route before adding academic records.',
        },
        {
            number: 3,
            label: 'KCSE details',
            heading: 'KCSE details',
            description:
                'Capture your KCSE index number, examination year, mean grade, and subject highlights for verification.',
        },
        {
            number: 4,
            label: 'Placement',
            heading: 'KUCCPS or self-sponsored placement',
            description:
                'Confirm KUCCPS placement details or mark the application as self-sponsored for direct admissions review.',
        },
        {
            number: 5,
            label: 'Documents',
            heading: 'Supporting documents',
            description:
                'Upload KCSE results, ID or passport, placement letter where applicable, and any finance office receipts.',
        },
        {
            number: 6,
            label: 'Review',
            heading: 'Review and submit',
            description:
                'Check your details, confirm the M-Pesa or waiver status, and submit your application for admissions review.',
        },
    ];

    readonly activeStepIndex = signal(0);
    readonly isLoadingOptions = signal(false);
    readonly isSaving = signal(false);
    readonly apiError = signal('');
    readonly apiMessage = signal('');
    readonly programs = signal<Program[]>([]);
    readonly academicTerms = signal<AcademicTerm[]>([]);
    readonly applicationFormTemplates = signal<ApplicationFormTemplate[]>([]);
    readonly application = signal<Application | null>(null);
    readonly portalSession = this.portalAuthService.session;
    readonly currentStep = computed(() => this.steps[this.activeStepIndex()]);
    readonly progressValue = computed(
        () => (this.currentStep().number / this.steps.length) * 100
    );
    readonly isFirstStep = computed(() => this.activeStepIndex() === 0);
    readonly isLastStep = computed(
        () => this.activeStepIndex() === this.steps.length - 1
    );
    readonly nextButtonLabel = computed(() =>
        this.isLastStep() ? 'Submit application' : 'Continue'
    );

    form: UntypedFormGroup;

    constructor() {
        this.form = this.formBuilder.group({
            firstName: [''],
            lastName: [''],
            email: [''],
            phone: [''],
            dateOfBirth: [''],
            password: [''],
            confirmPassword: [''],
            programme: [''],
            intake: ['2026 Main Intake'],
            sponsorshipRoute: [''],
            campusPreference: [''],
            kcseIndexNumber: [''],
            kcseYear: [''],
            kcseMeanGrade: [''],
            kcseSubjectHighlights: [''],
            kuccpsPlacementNumber: [''],
            applicationRoute: [''],
            county: [''],
            identityNumber: [''],
            mpesaConfirmationCode: [''],
            reviewNotes: [''],
        });
    }

    ngOnInit(): void {
        if (this.portalAuthService.hasValidSession()) {
            this.loadApplicantOptions();
        }
        this.loadExistingApplication();
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

    goToStep(stepNumber: number): void {
        const stepIndex = this.steps.findIndex((step) => step.number === stepNumber);

        if (stepIndex >= 0) {
            this.activeStepIndex.set(stepIndex);
        }
    }

    goToNextStep(): void {
        if (this.currentStep().number === 1 && !this.portalAuthService.hasValidSession()) {
            this.createApplicantSession(() => {
                this.loadApplicantOptions();
                this.advanceStep();
            });
            return;
        }

        if (this.isLastStep()) {
            this.submitApplication();
            return;
        }

        if (this.shouldPersistCurrentStep()) {
            this.saveDraft(() => this.advanceStep());
            return;
        }

        this.advanceStep();
    }

    goToPreviousStep(): void {
        this.activeStepIndex.update((stepIndex) => Math.max(stepIndex - 1, 0));
    }

    private loadApplicantOptions(): void {
        if (!this.portalAuthService.hasValidSession()) {
            return;
        }

        this.isLoadingOptions.set(true);
        this.apiError.set('');

        forkJoin({
            programs: this.admissionsService.queryApplicantPrograms({ rows: 100 }),
            academicTerms: this.admissionsService.queryApplicantAcademicTerms({ rows: 100 }),
            templates: this.admissionsService.queryApplicantApplicationFormTemplates({ rows: 100 }),
        }).subscribe({
            next: ({ programs, academicTerms, templates }) => {
                this.programs.set(programs.items);
                this.academicTerms.set(academicTerms.items);
                this.applicationFormTemplates.set(templates.items);
                this.prefillApplicantOptions();
                const application = this.application();
                if (application) {
                    this.patchApplication(application);
                }
                this.isLoadingOptions.set(false);
            },
            error: (error) => {
                this.isLoadingOptions.set(false);
                this.apiError.set(
                    jsonApiErrorMessage(
                        error,
                        'Sign in from the applicant portal home page before loading live programme and intake options.'
                    )
                );
            },
        });
    }

    private loadExistingApplication(): void {
        const applicationID = this.portalSession()?.applicationID;
        if (!applicationID) {
            return;
        }

        this.admissionsService.getApplicantApplication(applicationID).subscribe({
            next: (application) => {
                this.application.set(application);
                this.patchApplication(application);
            },
            error: (error) => {
                this.apiError.set(
                    jsonApiErrorMessage(
                        error,
                        'We could not load your current applicant application.'
                    )
                );
            },
        });
    }

    private saveDraft(onSuccess: () => void): void {
        const request = this.buildApplicationRequest();
        if (!request || this.isSaving()) {
            return;
        }

        this.isSaving.set(true);
        this.apiError.set('');
        this.apiMessage.set('');

        const existingApplication = this.application();
        const saveRequest = existingApplication
            ? this.admissionsService.updateApplicantApplication(
                  existingApplication.id,
                  request
              )
            : this.admissionsService.createApplicantApplication(request);

        saveRequest.subscribe({
            next: (application) => {
                this.application.set(application);
                this.portalAuthService.updateApplicationID(application.id);
                this.isSaving.set(false);
                this.apiMessage.set('Draft saved to admissions.');
                onSuccess();
            },
            error: (error) => {
                this.isSaving.set(false);
                this.apiError.set(
                    jsonApiErrorMessage(
                        error,
                        'We could not save your application draft.'
                    )
                );
            },
        });
    }

    private submitApplication(): void {
        this.saveDraft(() => {
            const application = this.application();
            if (!application || this.isSaving()) {
                return;
            }

            this.isSaving.set(true);
            this.apiError.set('');
            this.apiMessage.set('');

            this.admissionsService
                .transitionApplicantApplication(application.id, {
                    toStatus: 'SUBMITTED',
                    reason: 'Applicant submitted application from portal',
                    note: this.formString('reviewNotes') || null,
                    metadata: {
                        mpesaConfirmationCode:
                            this.formString('mpesaConfirmationCode') || null,
                    },
                })
                .subscribe({
                    next: (updated) => {
                        this.application.set(updated);
                        this.isSaving.set(false);
                        this.apiMessage.set(
                            'Application submitted for admissions review.'
                        );
                    },
                    error: (error) => {
                        this.isSaving.set(false);
                        this.apiError.set(
                            jsonApiErrorMessage(
                                error,
                                'We could not submit your application.'
                            )
                        );
                    },
                });
        });
    }

    private advanceStep(): void {
        this.activeStepIndex.update((stepIndex) =>
            Math.min(stepIndex + 1, this.steps.length - 1)
        );
    }

    private createApplicantSession(onSuccess: () => void): void {
        if (this.isSaving()) {
            return;
        }

        const dateOfBirth = this.formString('dateOfBirth');
        if (
            !this.formString('firstName') ||
            !this.formString('lastName') ||
            !this.formString('email') ||
            !this.formString('phone') ||
            !dateOfBirth ||
            !this.formString('password') ||
            !this.formString('confirmPassword')
        ) {
            this.apiError.set('Complete your applicant account details before continuing.');
            return;
        }

        this.isSaving.set(true);
        this.apiError.set('');
        this.apiMessage.set('');

        this.portalAuthService
            .onboardApplicant({
                firstName: this.formString('firstName'),
                lastName: this.formString('lastName'),
                email: this.formString('email'),
                phone: this.formString('phone'),
                password: this.formString('password'),
                confirmPassword: this.formString('confirmPassword'),
                dateOfBirth: new Date(`${dateOfBirth}T00:00:00Z`).toISOString(),
            })
            .subscribe({
                next: () => {
                    this.isSaving.set(false);
                    this.apiMessage.set('Applicant account created.');
                    onSuccess();
                },
                error: (error) => {
                    this.isSaving.set(false);
                    this.apiError.set(
                        jsonApiErrorMessage(
                            error,
                            'We could not create your applicant account.'
                        )
                    );
                },
            });
    }

    private shouldPersistCurrentStep(): boolean {
        return this.currentStep().number >= 2;
    }

    private buildApplicationRequest(): ApplicationRequest | null {
        const constituentID = this.portalSession()?.constituentID;
        if (!constituentID) {
            this.apiError.set(
                'Sign in from the applicant portal home page before saving your live application.'
            );
            return null;
        }

        const program = this.selectedProgram();
        const academicTerm = this.selectedAcademicTerm();
        if (!program || !academicTerm) {
            this.apiError.set(
                'Choose a valid programme and intake before saving your application.'
            );
            return null;
        }

        const applicationType = this.selectedApplicationType();

        return {
            constituentID,
            programID: program.id,
            academicTermID: academicTerm.id,
            applicationType,
            kuccpsPlacement: this.buildKuccpsPlacement(program, applicationType),
            kcseResult: this.buildKcseResult(),
            assignedReviewerID: null,
        };
    }

    private prefillApplicantOptions(): void {
        if (!this.formString('programme') && this.programs().length > 0) {
            this.form.patchValue({ programme: this.programs()[0].name });
        }

        if (!this.formString('intake') && this.academicTerms().length > 0) {
            this.form.patchValue({ intake: this.academicTerms()[0].name });
        }

        if (!this.formString('sponsorshipRoute')) {
            this.form.patchValue({ sponsorshipRoute: 'KUCCPS placement' });
        }
    }

    private patchApplication(application: Application): void {
        this.form.patchValue({
            programme:
                this.programs().find((program) => program.id === application.programID)
                    ?.name ?? application.programID,
            intake:
                this.academicTerms().find(
                    (term) => term.id === application.academicTermID
                )?.name ?? application.academicTermID,
            sponsorshipRoute:
                application.applicationType === 'KUCCPS_PLACEMENT'
                    ? 'KUCCPS placement'
                    : 'Self-sponsored',
            kcseIndexNumber: application.kcseResult?.indexNumber ?? '',
            kcseYear: application.kcseResult?.examYear ?? '',
            kcseMeanGrade: application.kcseResult?.meanGrade ?? '',
            kcseSubjectHighlights:
                application.kcseResult?.subjects
                    .map((subject) => `${subject.subjectCode} ${subject.grade}`)
                    .join(', ') ?? '',
            kuccpsPlacementNumber:
                application.kuccpsPlacement?.placementID ?? '',
        });
    }

    private selectedProgram(): Program | undefined {
        const value = this.formString('programme').toLowerCase();

        return this.programs().find(
            (program) =>
                program.id.toLowerCase() === value ||
                program.name.toLowerCase() === value ||
                program.code.toLowerCase() === value
        );
    }

    private selectedAcademicTerm(): AcademicTerm | undefined {
        const value = this.formString('intake').toLowerCase();

        return this.academicTerms().find(
            (term) =>
                term.id.toLowerCase() === value ||
                term.name.toLowerCase() === value ||
                term.code.toLowerCase() === value
        );
    }

    private selectedApplicationType(): ApplicationType {
        const route = `${this.formString('sponsorshipRoute')} ${this.formString('applicationRoute')}`.toLowerCase();

        return route.includes('kuccps') || Boolean(this.formString('kuccpsPlacementNumber'))
            ? 'KUCCPS_PLACEMENT'
            : 'SELF_SPONSORED_UNDERGRAD';
    }

    private buildKuccpsPlacement(
        program: Program,
        applicationType: ApplicationType
    ): KUCCPSPlacement | null {
        if (applicationType !== 'KUCCPS_PLACEMENT') {
            return null;
        }

        const placementID = this.formString('kuccpsPlacementNumber');
        if (!placementID) {
            return null;
        }

        return {
            placementID,
            institutionCode: 'STU',
            programmeCode: program.code,
            programmeName: program.name,
            placementYear: this.numericFormValue('kcseYear') || new Date().getFullYear(),
            weightedPointsNote: this.formString('applicationRoute') || undefined,
        };
    }

    private buildKcseResult(): ApplicationKCSEResult | null {
        const indexNumber = this.formString('kcseIndexNumber');
        const examYear = this.numericFormValue('kcseYear');
        const meanGrade = this.formString('kcseMeanGrade').toUpperCase();

        if (!indexNumber || !examYear || !meanGrade) {
            return null;
        }

        return {
            indexNumber,
            examYear,
            subjects: this.parseSubjectHighlights(),
            meanGrade,
            meanPoints: this.gradePoints(meanGrade),
        };
    }

    private parseSubjectHighlights(): ApplicationKCSEResult['subjects'] {
        return this.formString('kcseSubjectHighlights')
            .split(',')
            .map((part) => part.trim())
            .filter(Boolean)
            .map((part) => {
                const [subjectCode, grade = ''] = part.split(/\s+/);
                const normalizedGrade = grade.toUpperCase();

                return {
                    subjectCode: subjectCode.toUpperCase(),
                    grade: normalizedGrade,
                    points: this.gradePoints(normalizedGrade),
                };
            });
    }

    private gradePoints(grade: string): number {
        switch (grade.toUpperCase()) {
            case 'A':
                return 12;
            case 'A-':
                return 11;
            case 'B+':
                return 10;
            case 'B':
                return 9;
            case 'B-':
                return 8;
            case 'C+':
                return 7;
            case 'C':
                return 6;
            case 'C-':
                return 5;
            case 'D+':
                return 4;
            case 'D':
                return 3;
            case 'D-':
                return 2;
            case 'E':
                return 1;
            default:
                return 0;
        }
    }

    private formString(controlName: string): string {
        const value = this.form.get(controlName)?.value;

        return typeof value === 'string' ? value.trim() : String(value ?? '').trim();
    }

    private numericFormValue(controlName: string): number {
        const value = Number(this.formString(controlName));

        return Number.isFinite(value) ? value : 0;
    }
}
