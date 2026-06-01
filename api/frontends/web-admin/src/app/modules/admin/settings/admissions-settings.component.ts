import { AsyncPipe } from '@angular/common';
import { ChangeDetectionStrategy, Component, inject } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatChipsModule } from '@angular/material/chips';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { MatSlideToggleModule } from '@angular/material/slide-toggle';
import { AdmissionsService } from 'app/core/admissions/admissions.service';
import {
    APPLICATION_STATUSES,
    APPLICATION_TYPES,
    ApplicationChecklistTemplateItem,
    ApplicationFormField,
    ApplicationFormTemplate,
    ApplicationFormTemplateRequest,
    ApplicationType,
    AssignmentCandidate,
    EVENT_ATTENDANCE_STATUSES,
    LEAD_ASSIGNMENT_CRITERION_FIELDS,
    LEAD_ASSIGNMENT_CRITERION_OPERATORS,
    LEAD_SCORE_CRITERION_FIELDS,
    LEAD_SCORE_CRITERION_OPERATORS,
    LIFECYCLE_STAGES,
    LeadAssignmentCriterionField,
    LeadAssignmentRule,
    LeadScoreCriterionField,
    LeadScoreCriterionOperator,
    LeadScoreRule,
    LeadScoreRuleRequest,
} from 'app/core/admissions/admissions.types';
import { jsonApiErrorMessage } from 'app/core/api/json-api';
import { catchError, of } from 'rxjs';

interface LeadScoreRuleForm {
    id: string;
    name: string;
    description: string;
    field: LeadScoreCriterionField;
    operator: LeadScoreCriterionOperator;
    valuesText: string;
    points: number;
    priority: number;
    active: boolean;
}

interface ApplicationFormTemplateForm {
    id: string;
    programID: string;
    academicTermID: string;
    applicationType: ApplicationType;
    name: string;
    description: string;
    requiredFieldsText: string;
    checklistItemsText: string;
    priority: number;
    active: boolean;
}

interface AssignmentRuleForm {
    id: string;
    name: string;
    description: string;
    field: LeadAssignmentCriterionField;
    operator: 'EQ' | 'IN' | 'BETWEEN';
    valuesText: string;
    assignee: string;
    priority: number;
    active: boolean;
}

const emptyForm: LeadScoreRuleForm = {
    id: '',
    name: '',
    description: '',
    field: 'lifecycle_stage',
    operator: 'EQ',
    valuesText: 'INQUIRY',
    points: 10,
    priority: 100,
    active: true,
};

const defaultRequiredFields: ApplicationFormField[] = [
    {
        fieldName: 'personal_statement',
        fieldType: 'textarea',
        required: true,
        displayOrder: 1,
    },
];

const defaultChecklistItems: ApplicationChecklistTemplateItem[] = [
    {
        itemKey: 'transcript',
        documentName: 'Transcript',
        description: 'Official transcript for this application type.',
        required: true,
        displayOrder: 1,
    },
];

const emptyTemplateForm: ApplicationFormTemplateForm = {
    id: '',
    programID: '',
    academicTermID: '',
    applicationType: 'FRESHMAN',
    name: '',
    description: '',
    requiredFieldsText: JSON.stringify(defaultRequiredFields, null, 2),
    checklistItemsText: JSON.stringify(defaultChecklistItems, null, 2),
    priority: 100,
    active: true,
};

const emptyAssignmentForm: AssignmentRuleForm = {
    id: '',
    name: '',
    description: '',
    field: 'territory',
    operator: 'IN',
    valuesText: 'Northeast, Mid-Atlantic',
    assignee: 'Maya Schultz',
    priority: 100,
    active: true,
};

const TERRITORIES = [
    'Northeast',
    'Mid-Atlantic',
    'Midwest',
    'Southwest',
] as const;

const PROGRAMS = [
    'Computer Science',
    'Business Analytics',
    'B.Sc. CS',
    'M.B.A.',
    'B.A. Econ',
] as const;

const TERMS = ['Fall 2026', 'Spring 2027', 'Spring 2026'] as const;

const LEAD_SOURCES = [
    'Web form',
    'Open house',
    'Counselor',
    'Referral',
] as const;

const ASSIGNEES = [
    'Maya Schultz',
    'Andre Park',
    'Avery',
    'Priya',
    'Diego',
] as const;

const ASSIGNMENT_RULES: LeadAssignmentRule[] = [
    {
        id: 'assign-northeast-cs',
        name: 'Northeast CS recruiter route',
        description:
            'Prioritizes regional Computer Science prospects for the recruiter covering Northeast and Mid-Atlantic territories.',
        criteria: [
            {
                field: 'territory',
                operator: 'IN',
                values: ['Northeast', 'Mid-Atlantic'],
            },
            {
                field: 'program_interest',
                operator: 'IN',
                values: ['Computer Science', 'B.Sc. CS'],
            },
        ],
        assignee: 'Maya Schultz',
        assigneeRole: 'RECRUITER',
        active: true,
        priority: 10,
        lastMatchedCount: 128,
        auditTrail: [
            {
                actor: 'Owen Adirah',
                action: 'Activated assignment rule',
                reason: 'Regional recruiter ownership for Fall 2026 CS demand.',
                timestamp: 'Jun 1, 2026 09:20',
            },
        ],
    },
    {
        id: 'assign-open-house-hot-leads',
        name: 'Open house attendee follow-up',
        description:
            'Routes attended event prospects to the recruiter queue for same-day follow-up.',
        criteria: [
            {
                field: 'event_attendance',
                operator: 'EQ',
                values: ['ATTENDED'],
            },
            {
                field: 'lead_source',
                operator: 'IN',
                values: ['Open house', 'Referral'],
            },
        ],
        assignee: 'Andre Park',
        assigneeRole: 'RECRUITER',
        active: true,
        priority: 20,
        lastMatchedCount: 64,
        auditTrail: [
            {
                actor: 'System',
                action: 'Auto-assigned matching leads',
                reason: 'Event attendance rule evaluated after nightly sync.',
                timestamp: 'Jun 1, 2026 07:00',
            },
        ],
    },
    {
        id: 'assign-alpha-review',
        name: 'Alpha split application review',
        description:
            'Balances submitted applications by applicant surname range when no higher priority program rule matches.',
        criteria: [
            {
                field: 'alpha_range',
                operator: 'BETWEEN',
                values: ['A', 'M'],
            },
            {
                field: 'application_type',
                operator: 'IN',
                values: ['FRESHMAN', 'TRANSFER'],
            },
        ],
        assignee: 'Avery',
        assigneeRole: 'APPLICATION_REVIEWER',
        active: true,
        priority: 90,
        lastMatchedCount: 216,
        auditTrail: [
            {
                actor: 'Priya',
                action: 'Adjusted priority',
                reason: 'Keep alpha split after territory and program rules.',
                timestamp: 'May 30, 2026 16:45',
            },
        ],
    },
];

const ASSIGNMENT_CANDIDATES: AssignmentCandidate[] = [
    {
        id: 'APP-3024',
        applicant: 'Sofia Martinez',
        program: 'B.Sc. CS',
        term: 'Fall 2026',
        source: 'Open house',
        territory: 'Northeast',
        eventAttendance: 'ATTENDED',
        alphaKey: 'M',
        currentAssignee: 'Avery',
        recommendedAssignee: 'Maya Schultz',
        matchedRule: 'Northeast CS recruiter route',
        assignmentMode: 'MANUAL_OVERRIDE',
        overrideReason:
            'Keep with original reviewer until interview notes close.',
        overrideActor: 'Priya',
        overrideTimestamp: 'Jun 1, 2026 10:12',
    },
    {
        id: 'APP-3023',
        applicant: 'James Okoro',
        program: 'M.A. Edu',
        term: 'Spring 2027',
        source: 'Counselor',
        territory: 'Midwest',
        eventAttendance: 'REGISTERED',
        alphaKey: 'O',
        currentAssignee: 'Priya',
        recommendedAssignee: 'Priya',
        matchedRule: 'Alpha split application review',
        assignmentMode: 'AUTO',
    },
    {
        id: 'LEAD-1442',
        applicant: 'Aisha Bello',
        program: 'M.B.A.',
        term: 'Fall 2026',
        source: 'Referral',
        territory: 'Southwest',
        eventAttendance: 'ATTENDED',
        alphaKey: 'B',
        currentAssignee: 'Andre Park',
        recommendedAssignee: 'Andre Park',
        matchedRule: 'Open house attendee follow-up',
        assignmentMode: 'AUTO',
    },
];

@Component({
    selector: 'app-admissions-settings',
    imports: [
        AsyncPipe,
        FormsModule,
        MatButtonModule,
        MatChipsModule,
        MatFormFieldModule,
        MatIconModule,
        MatInputModule,
        MatSelectModule,
        MatSlideToggleModule,
    ],
    templateUrl: './admissions-settings.component.html',
    changeDetection: ChangeDetectionStrategy.OnPush,
})
export class AdmissionsSettingsComponent {
    private readonly admissionsService = inject(AdmissionsService);

    readonly fields = LEAD_SCORE_CRITERION_FIELDS;
    readonly operators = LEAD_SCORE_CRITERION_OPERATORS;
    readonly assignmentFields = LEAD_ASSIGNMENT_CRITERION_FIELDS;
    readonly assignmentOperators = LEAD_ASSIGNMENT_CRITERION_OPERATORS;
    readonly lifecycleStages = LIFECYCLE_STAGES;
    readonly applicationTypes = APPLICATION_TYPES;
    readonly applicationStatuses = APPLICATION_STATUSES;
    readonly eventAttendanceStatuses = EVENT_ATTENDANCE_STATUSES;
    readonly assignmentRules = ASSIGNMENT_RULES;
    readonly assignmentCandidates = ASSIGNMENT_CANDIDATES;
    readonly assignees = ASSIGNEES;
    readonly rules$ = this.admissionsService.rules$;
    readonly templates$ = this.admissionsService.templates$;

    errorMessage = '';
    saving = false;
    form: LeadScoreRuleForm = { ...emptyForm };
    templateSaving = false;
    templateForm: ApplicationFormTemplateForm = { ...emptyTemplateForm };
    assignmentForm: AssignmentRuleForm = { ...emptyAssignmentForm };

    constructor() {
        this.loadRules();
        this.loadTemplates();
    }

    loadRules(): void {
        this.errorMessage = '';
        this.admissionsService
            .queryLeadScoreRules({
                page: 1,
                rows: 50,
                orderBy: 'priority,ASC',
            })
            .pipe(
                catchError((error) => {
                    this.errorMessage = jsonApiErrorMessage(
                        error,
                        'Unable to load lead scoring rules.'
                    );
                    return of(undefined);
                })
            )
            .subscribe();
    }

    loadTemplates(): void {
        this.errorMessage = '';
        this.admissionsService
            .queryApplicationFormTemplates({
                page: 1,
                rows: 50,
                orderBy: 'priority,ASC',
            })
            .pipe(
                catchError((error) => {
                    this.errorMessage = jsonApiErrorMessage(
                        error,
                        'Unable to load application form templates.'
                    );
                    return of(undefined);
                })
            )
            .subscribe();
    }

    editRule(rule: LeadScoreRule): void {
        const criterion = rule.criteria[0];
        this.form = {
            id: rule.id,
            name: rule.name,
            description: rule.description ?? '',
            field: criterion?.field ?? 'lifecycle_stage',
            operator: criterion?.operator ?? 'EQ',
            valuesText: criterion?.values.join(', ') ?? '',
            points: rule.points,
            priority: rule.priority,
            active: rule.active,
        };
    }

    resetForm(): void {
        this.form = { ...emptyForm };
    }

    editTemplate(template: ApplicationFormTemplate): void {
        this.templateForm = {
            id: template.id,
            programID: template.programID,
            academicTermID: template.academicTermID,
            applicationType: template.applicationType,
            name: template.name,
            description: template.description ?? '',
            requiredFieldsText: JSON.stringify(
                template.requiredFields,
                null,
                2
            ),
            checklistItemsText: JSON.stringify(
                template.checklistItems,
                null,
                2
            ),
            priority: template.priority,
            active: template.active,
        };
    }

    resetTemplateForm(): void {
        this.templateForm = { ...emptyTemplateForm };
    }

    editAssignmentRule(rule: LeadAssignmentRule): void {
        const criterion = rule.criteria[0];
        this.assignmentForm = {
            id: rule.id,
            name: rule.name,
            description: rule.description,
            field: criterion?.field ?? 'territory',
            operator: criterion?.operator ?? 'IN',
            valuesText: criterion?.values.join(', ') ?? '',
            assignee: rule.assignee,
            priority: rule.priority,
            active: rule.active,
        };
    }

    resetAssignmentForm(): void {
        this.assignmentForm = { ...emptyAssignmentForm };
    }

    saveRule(): void {
        this.errorMessage = '';
        this.saving = true;

        const request = this.toRequest(this.form);
        const save$ = this.form.id
            ? this.admissionsService.updateLeadScoreRule(this.form.id, request)
            : this.admissionsService.createLeadScoreRule(request);

        save$
            .pipe(
                catchError((error) => {
                    this.errorMessage = jsonApiErrorMessage(
                        error,
                        'Unable to save lead scoring rule.'
                    );
                    return of(undefined);
                })
            )
            .subscribe(() => {
                this.saving = false;
                this.resetForm();
                this.loadRules();
            });
    }

    saveTemplate(): void {
        this.errorMessage = '';
        this.templateSaving = true;

        let request: ApplicationFormTemplateRequest;
        try {
            request = this.toTemplateRequest(this.templateForm);
        } catch (error) {
            this.errorMessage =
                'Required fields and checklist items must be valid JSON arrays.';
            this.templateSaving = false;
            return;
        }

        const save$ = this.templateForm.id
            ? this.admissionsService.updateApplicationFormTemplate(
                  this.templateForm.id,
                  request
              )
            : this.admissionsService.createApplicationFormTemplate(request);

        save$
            .pipe(
                catchError((error) => {
                    this.errorMessage = jsonApiErrorMessage(
                        error,
                        'Unable to save application form template.'
                    );
                    return of(undefined);
                })
            )
            .subscribe(() => {
                this.templateSaving = false;
                this.resetTemplateForm();
                this.loadTemplates();
            });
    }

    criterionSummary(rule: LeadScoreRule): string {
        return rule.criteria
            .map(
                (criterion) =>
                    `${criterion.field} ${criterion.operator} ${criterion.values.join(', ')}`
            )
            .join(' · ');
    }

    suggestedValues(): readonly string[] {
        switch (this.form.field) {
            case 'lifecycle_stage':
                return this.lifecycleStages;
            case 'application_type':
                return this.applicationTypes;
            case 'application_status':
                return this.applicationStatuses;
            case 'program_id':
            case 'academic_term_id':
                return [];
        }
    }

    templateSummary(template: ApplicationFormTemplate): string {
        return `${template.requiredFields.length} fields · ${template.checklistItems.length} checklist items · v${template.version}`;
    }

    assignmentCriterionSummary(rule: LeadAssignmentRule): string {
        return rule.criteria
            .map(
                (criterion) =>
                    `${criterion.field} ${criterion.operator} ${criterion.values.join(', ')}`
            )
            .join(' · ');
    }

    assignmentSuggestedValues(): readonly string[] {
        switch (this.assignmentForm.field) {
            case 'territory':
                return TERRITORIES;
            case 'program_interest':
                return PROGRAMS;
            case 'application_type':
                return this.applicationTypes;
            case 'academic_term':
                return TERMS;
            case 'lead_source':
                return LEAD_SOURCES;
            case 'event_attendance':
                return this.eventAttendanceStatuses;
            case 'alpha_range':
                return ['A, M', 'N, Z'];
        }
    }

    private toRequest(form: LeadScoreRuleForm): LeadScoreRuleRequest {
        return {
            name: form.name,
            description: form.description.trim() || null,
            criteria: [
                {
                    field: form.field,
                    operator: form.operator,
                    values: form.valuesText
                        .split(',')
                        .map((value) => value.trim())
                        .filter((value) => value.length > 0),
                },
            ],
            points: Number(form.points),
            priority: Number(form.priority),
            active: form.active,
        };
    }

    private toTemplateRequest(
        form: ApplicationFormTemplateForm
    ): ApplicationFormTemplateRequest {
        return {
            programID: form.programID.trim(),
            academicTermID: form.academicTermID.trim(),
            applicationType: form.applicationType,
            name: form.name.trim(),
            description: form.description.trim() || null,
            requiredFields: JSON.parse(
                form.requiredFieldsText
            ) as ApplicationFormField[],
            checklistItems: JSON.parse(
                form.checklistItemsText
            ) as ApplicationChecklistTemplateItem[],
            priority: Number(form.priority),
            active: form.active,
        };
    }
}
