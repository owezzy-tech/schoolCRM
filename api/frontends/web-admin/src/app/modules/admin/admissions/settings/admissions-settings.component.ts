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

    errorMessage = '';
    saving = false;
    form: LeadScoreRuleForm = { ...emptyForm };

    constructor() {
        this.loadRules();
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
}
