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
    LEAD_SCORE_CRITERION_FIELDS,
    LEAD_SCORE_CRITERION_OPERATORS,
    LIFECYCLE_STAGES,
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
    readonly lifecycleStages = LIFECYCLE_STAGES;
    readonly applicationTypes = APPLICATION_TYPES;
    readonly applicationStatuses = APPLICATION_STATUSES;
    readonly rules$ = this.admissionsService.rules$;
    readonly templates$ = this.admissionsService.templates$;

    errorMessage = '';
    saving = false;
    form: LeadScoreRuleForm = { ...emptyForm };
    templateSaving = false;
    templateForm: ApplicationFormTemplateForm = { ...emptyTemplateForm };

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
