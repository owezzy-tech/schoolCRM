import { AsyncPipe, JsonPipe } from '@angular/common';
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
    CUSTOM_FIELD_DATA_TYPES,
    CUSTOM_FIELD_OWNERS,
    CustomFieldDataType,
    CustomFieldDefinition,
    CustomFieldDefinitionRequest,
    CustomFieldOwner,
    EVENT_ATTENDANCE_STATUSES,
    IMPORT_BATCH_STATUSES,
    IMPORT_FILE_TYPES,
    IMPORT_SOURCES,
    IMPORT_TARGETS,
    ImportBatch,
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

interface CustomFieldDefinitionForm {
    id: string;
    owner: CustomFieldOwner;
    fieldKey: string;
    label: string;
    description: string;
    dataType: CustomFieldDataType;
    required: boolean;
    optionsText: string;
    validation: string;
    searchable: boolean;
    reportable: boolean;
    importable: boolean;
    exportable: boolean;
    displayOrder: number;
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

interface ImportTemplateColumn {
    key: string;
    label: string;
    required: boolean;
    owner: 'Core' | 'Custom Field';
    example: string;
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
        fieldName: 'kcse_index_number',
        fieldType: 'text',
        required: true,
        displayOrder: 1,
        validation: 'kenya-kcse-index',
    },
    {
        fieldName: 'programme_id',
        fieldType: 'select',
        required: true,
        displayOrder: 2,
    },
];

const defaultChecklistItems: ApplicationChecklistTemplateItem[] = [
    {
        itemKey: 'kcse_result_slip',
        documentName: 'KCSE result slip',
        description: 'Official KCSE result slip or KNEC verification evidence.',
        required: true,
        displayOrder: 1,
    },
];

const emptyTemplateForm: ApplicationFormTemplateForm = {
    id: '',
    programID: '',
    academicTermID: '',
    applicationType: 'KUCCPS_PLACEMENT',
    name: '',
    description: '',
    requiredFieldsText: JSON.stringify(defaultRequiredFields, null, 2),
    checklistItemsText: JSON.stringify(defaultChecklistItems, null, 2),
    priority: 100,
    active: true,
};

const emptyCustomFieldForm: CustomFieldDefinitionForm = {
    id: '',
    owner: 'CONSTITUENT',
    fieldKey: '',
    label: '',
    description: '',
    dataType: 'TEXT',
    required: false,
    optionsText: '',
    validation: '',
    searchable: true,
    reportable: true,
    importable: true,
    exportable: true,
    displayOrder: 100,
    active: true,
};

const CUSTOM_FIELD_DEFINITIONS: CustomFieldDefinition[] = [
    {
        id: 'field-kuccps-placement-id',
        owner: 'APPLICATION',
        fieldKey: 'kuccps_placement_id',
        label: 'KUCCPS placement ID',
        description:
            'KUCCPS placement reference captured before admissions review.',
        dataType: 'TEXT',
        required: false,
        options: [],
        validation: 'max:40',
        searchable: true,
        reportable: true,
        importable: true,
        exportable: true,
        displayOrder: 5,
        active: true,
        dateCreated: '2026-06-01T08:00:00Z',
        dateUpdated: '2026-06-01T08:00:00Z',
    },
    {
        id: 'field-kcse-mean-grade',
        owner: 'CONSTITUENT',
        fieldKey: 'kcse_mean_grade',
        label: 'KCSE mean grade',
        description:
            'Normalized KCSE mean grade for application eligibility reporting.',
        dataType: 'SELECT',
        required: false,
        options: [
            'A',
            'A-',
            'B+',
            'B',
            'B-',
            'C+',
            'C',
            'C-',
            'D+',
            'D',
            'D-',
            'E',
        ],
        searchable: true,
        reportable: true,
        importable: true,
        exportable: true,
        displayOrder: 10,
        active: true,
        dateCreated: '2026-06-01T08:00:00Z',
        dateUpdated: '2026-06-01T08:00:00Z',
    },
    {
        id: 'field-kuccps-cluster-points',
        owner: 'APPLICATION',
        fieldKey: 'kuccps_cluster_points',
        label: 'KUCCPS cluster points',
        description:
            'Public cluster-point approximation used until private KNEC PI data is available.',
        dataType: 'NUMBER',
        required: false,
        options: [],
        validation: 'min:0|max:48',
        searchable: true,
        reportable: true,
        importable: true,
        exportable: true,
        displayOrder: 20,
        active: true,
        dateCreated: '2026-06-01T08:00:00Z',
        dateUpdated: '2026-06-01T08:00:00Z',
    },
    {
        id: 'field-alumni-connection',
        owner: 'CONSTITUENT',
        fieldKey: 'alumni_connection',
        label: 'Alumni connection',
        description:
            'Optional context for family or mentor relationship reporting.',
        dataType: 'BOOLEAN',
        required: false,
        options: [],
        searchable: false,
        reportable: true,
        importable: true,
        exportable: true,
        displayOrder: 30,
        active: false,
        dateCreated: '2026-06-01T08:00:00Z',
        dateUpdated: '2026-06-01T08:00:00Z',
    },
];

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
    'B.Sc. Computer Science',
    'Bachelor of Commerce',
    'Diploma in ICT',
    'M.A. Education',
    'Certificate in Data Support',
] as const;

const TERMS = ['September 2026', 'January 2027', 'May 2027'] as const;

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
                values: ['KUCCPS_PLACEMENT', 'SELF_SPONSORED_UNDERGRAD'],
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

const IMPORT_BATCHES: ImportBatch[] = [
    {
        id: 'import-2026-0601-constituents',
        source: 'MANUAL_UPLOAD',
        fileType: 'CSV',
        target: 'CONSTITUENTS',
        status: 'VALIDATION_FAILED',
        fileName: 'fall-2026-open-house-leads.csv',
        storageKey: 'admissions/imports/fall-2026-open-house-leads.csv',
        uploadedByID: 'user-maya-schultz',
        totalRows: 420,
        validRows: 392,
        invalidRows: 18,
        duplicateRows: 10,
        fieldMapping: {
            first_name: 'First Name',
            last_name: 'Last Name',
            primary_email: 'Email',
            scholarship_level: 'Scholarship Level',
        },
        invalidReportKey:
            'admissions/imports/reports/fall-2026-open-house-leads-invalid.csv',
        validationSummary:
            '18 rows need correction: missing email, invalid phone, or unknown program code.',
        dateCreated: '2026-06-01T10:30:00Z',
        dateUpdated: '2026-06-01T10:34:00Z',
    },
    {
        id: 'import-2026-0531-applications',
        source: 'SIS_EXPORT',
        fileType: 'XLSX',
        target: 'APPLICATIONS',
        status: 'COMPLETED',
        fileName: 'spring-2027-applications.xlsx',
        uploadedByID: 'user-priya-reviewer',
        totalRows: 128,
        validRows: 124,
        invalidRows: 0,
        duplicateRows: 4,
        fieldMapping: {
            constituent_id: 'Student ID',
            program_id: 'Program Code',
            academic_term_id: 'Term Code',
        },
        committedAt: '2026-05-31T16:45:00Z',
        validationSummary:
            '124 applications committed. 4 duplicate active applications skipped.',
        dateCreated: '2026-05-31T16:35:00Z',
        dateUpdated: '2026-05-31T16:45:00Z',
    },
];

const IMPORT_TEMPLATE_COLUMNS: ImportTemplateColumn[] = [
    {
        key: 'first_name',
        label: 'First name',
        required: true,
        owner: 'Core',
        example: 'Aisha',
    },
    {
        key: 'last_name',
        label: 'Last name',
        required: true,
        owner: 'Core',
        example: 'Bello',
    },
    {
        key: 'primary_email',
        label: 'Primary email',
        required: true,
        owner: 'Core',
        example: 'aisha@example.edu',
    },
    {
        key: 'program_id',
        label: 'Programme ID',
        required: true,
        owner: 'Core',
        example: 'BSC-CS',
    },
    {
        key: 'kcse_index_number',
        label: 'KCSE index number',
        required: true,
        owner: 'Core',
        example: '204003024',
    },
    {
        key: 'kuccps_placement_id',
        label: 'KUCCPS placement ID',
        required: false,
        owner: 'Custom Field',
        example: 'KUCCPS-2026-003024',
    },
    {
        key: 'kcse_mean_grade',
        label: 'KCSE mean grade',
        required: false,
        owner: 'Custom Field',
        example: 'A-',
    },
];

@Component({
    selector: 'app-admissions-settings',
    imports: [
        AsyncPipe,
        JsonPipe,
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
    readonly customFieldOwners = CUSTOM_FIELD_OWNERS;
    readonly customFieldDataTypes = CUSTOM_FIELD_DATA_TYPES;
    readonly customFieldDefinitions = CUSTOM_FIELD_DEFINITIONS;
    readonly importBatchStatuses = IMPORT_BATCH_STATUSES;
    readonly importSources = IMPORT_SOURCES;
    readonly importFileTypes = IMPORT_FILE_TYPES;
    readonly importTargets = IMPORT_TARGETS;
    readonly importBatches = IMPORT_BATCHES;
    readonly importTemplateColumns = IMPORT_TEMPLATE_COLUMNS;
    readonly eventAttendanceStatuses = EVENT_ATTENDANCE_STATUSES;
    readonly assignmentRules = ASSIGNMENT_RULES;
    readonly assignmentCandidates = ASSIGNMENT_CANDIDATES;
    readonly assignees = ASSIGNEES;
    readonly rules$ = this.admissionsService.rules$;
    readonly templates$ = this.admissionsService.templates$;
    readonly customFieldDefinitions$ =
        this.admissionsService.customFieldDefinitions$;

    errorMessage = '';
    saving = false;
    form: LeadScoreRuleForm = { ...emptyForm };
    templateSaving = false;
    templateForm: ApplicationFormTemplateForm = { ...emptyTemplateForm };
    customFieldForm: CustomFieldDefinitionForm = { ...emptyCustomFieldForm };
    assignmentForm: AssignmentRuleForm = { ...emptyAssignmentForm };

    constructor() {
        this.loadRules();
        this.loadTemplates();
    }

    editCustomField(definition: CustomFieldDefinition): void {
        this.customFieldForm = {
            id: definition.id,
            owner: definition.owner,
            fieldKey: definition.fieldKey,
            label: definition.label,
            description: definition.description ?? '',
            dataType: definition.dataType,
            required: definition.required,
            optionsText: definition.options.join(', '),
            validation: definition.validation ?? '',
            searchable: definition.searchable,
            reportable: definition.reportable,
            importable: definition.importable,
            exportable: definition.exportable,
            displayOrder: definition.displayOrder,
            active: definition.active,
        };
    }

    resetCustomFieldForm(): void {
        this.customFieldForm = { ...emptyCustomFieldForm };
    }

    customFieldSummary(definition: CustomFieldDefinition): string {
        const flags = [
            definition.searchable ? 'Search' : '',
            definition.reportable ? 'Report' : '',
            definition.importable ? 'Import' : '',
            definition.exportable ? 'Export' : '',
        ].filter(Boolean);

        return `${definition.dataType} · ${definition.required ? 'Required' : 'Optional'} · ${flags.join(' / ')}`;
    }

    customFieldOptionsLabel(definition: CustomFieldDefinition): string {
        return definition.options.length > 0
            ? definition.options.join(', ')
            : 'No fixed options';
    }

    toCustomFieldRequest(
        form: CustomFieldDefinitionForm
    ): CustomFieldDefinitionRequest {
        return {
            owner: form.owner,
            fieldKey: form.fieldKey.trim(),
            label: form.label.trim(),
            description: form.description.trim() || null,
            dataType: form.dataType,
            required: form.required,
            options: form.optionsText
                .split(',')
                .map((value) => value.trim())
                .filter((value) => value.length > 0),
            validation: form.validation.trim() || null,
            searchable: form.searchable,
            reportable: form.reportable,
            importable: form.importable,
            exportable: form.exportable,
            displayOrder: Number(form.displayOrder),
            active: form.active,
        };
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

    importBatchProgress(batch: ImportBatch): number {
        if (batch.totalRows === 0) {
            return 0;
        }

        return Math.round((batch.validRows / batch.totalRows) * 100);
    }

    importBatchStatusLabel(batch: ImportBatch): string {
        return batch.status.replaceAll('_', ' ');
    }

    importBatchSummary(batch: ImportBatch): string {
        return `${batch.validRows} valid · ${batch.invalidRows} invalid · ${batch.duplicateRows} duplicates`;
    }

    importMappingSummary(batch: ImportBatch): string {
        return Object.entries(batch.fieldMapping)
            .map(([field, column]) => `${column} → ${field}`)
            .join(' · ');
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
