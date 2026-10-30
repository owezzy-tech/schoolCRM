import { TestBed } from '@angular/core/testing';
import { AdmissionsService } from 'app/core/admissions/admissions.service';
import {
    AcademicTerm,
    ApplicationFormTemplate,
    CustomFieldDefinition,
    ImportBatch,
    LeadScoreRule,
    Program,
} from 'app/core/admissions/admissions.types';
import { PaginatedResult } from 'app/core/api/json-api';
import { of, throwError } from 'rxjs';
import { beforeEach, describe, expect, it, vi } from 'vitest';

import { AdmissionsSettingsComponent } from './admissions-settings.component';

describe('AdmissionsSettingsComponent', () => {
    let queryPrograms: ReturnType<typeof vi.fn>;
    let queryAcademicTerms: ReturnType<typeof vi.fn>;
    let queryCustomFieldDefinitions: ReturnType<typeof vi.fn>;
    let queryImportBatches: ReturnType<typeof vi.fn>;
    let queryLeadScoreRules: ReturnType<typeof vi.fn>;
    let queryApplicationFormTemplates: ReturnType<typeof vi.fn>;
    let createCustomFieldDefinition: ReturnType<typeof vi.fn>;
    let updateCustomFieldDefinition: ReturnType<typeof vi.fn>;

    beforeEach(async () => {
        queryPrograms = vi.fn(() => of(result([commerceProgram])));
        queryAcademicTerms = vi.fn(() => of(result([septemberTerm])));
        queryCustomFieldDefinitions = vi.fn(() => of(result([countyField])));
        queryImportBatches = vi.fn(() => of(result([importBatch])));
        queryLeadScoreRules = vi.fn(() => of(result([leadScoreRule])));
        queryApplicationFormTemplates = vi.fn(() => of(result([formTemplate])));
        createCustomFieldDefinition = vi.fn(() => of(countyField));
        updateCustomFieldDefinition = vi.fn(() => of(countyField));

        await TestBed.configureTestingModule({
            imports: [AdmissionsSettingsComponent],
            providers: [
                {
                    provide: AdmissionsService,
                    useValue: {
                        academicTerms$: of(result([septemberTerm])),
                        customFieldDefinitions$: of(result([countyField])),
                        importBatches$: of(result([importBatch])),
                        programs$: of(result([commerceProgram])),
                        rules$: of(result([leadScoreRule])),
                        templates$: of(result([formTemplate])),
                        createCustomFieldDefinition,
                        queryAcademicTerms,
                        queryApplicationFormTemplates,
                        queryCustomFieldDefinitions,
                        queryImportBatches,
                        queryLeadScoreRules,
                        queryPrograms,
                        updateCustomFieldDefinition,
                    },
                },
            ],
        })
            .overrideComponent(AdmissionsSettingsComponent, {
                set: {
                    template: '',
                },
            })
            .compileComponents();
    });

    it('loads settings options from admissions APIs using live query contracts', () => {
        const fixture = TestBed.createComponent(AdmissionsSettingsComponent);
        const component = fixture.componentInstance;

        expect(queryPrograms).toHaveBeenCalledWith({
            page: 1,
            rows: 100,
            active: true,
            orderBy: 'name,ASC',
        });
        expect(queryAcademicTerms).toHaveBeenCalledWith({
            page: 1,
            rows: 100,
            active: true,
            orderBy: 'start_date,DESC',
        });
        expect(queryCustomFieldDefinitions).toHaveBeenCalledWith({
            page: 1,
            rows: 100,
            orderBy: 'display_order,ASC',
        });
        expect(queryImportBatches).toHaveBeenCalledWith({
            page: 1,
            rows: 25,
            orderBy: 'date_created,DESC',
        });
        expect(component.suggestedValues()).toEqual([
            'PROSPECT',
            'INQUIRY',
            'APPLICANT',
            'ADMITTED',
            'ENROLLED',
            'ALUMNI',
        ]);

        component.form.field = 'program_id';
        expect(component.suggestedValues()).toEqual([
            'Bachelor of Commerce (BCOM-QA)',
        ]);

        component.assignmentForm.field = 'academic_term';
        expect(component.assignmentSuggestedValues()).toEqual([
            'September 2026 Intake (SEP2026-QA)',
        ]);
    });

    it('saves a custom field through the admissions API and reloads definitions', () => {
        const fixture = TestBed.createComponent(AdmissionsSettingsComponent);
        const component = fixture.componentInstance;
        queryCustomFieldDefinitions.mockClear();

        component.customFieldForm = {
            ...component.customFieldForm,
            fieldKey: 'county_code',
            label: 'County of Residence',
            dataType: 'SELECT',
            optionsText: '1, 16, 47',
        };

        component.saveCustomField();

        expect(createCustomFieldDefinition).toHaveBeenCalledWith({
            owner: 'CONSTITUENT',
            fieldKey: 'county_code',
            label: 'County of Residence',
            description: null,
            dataType: 'SELECT',
            required: false,
            options: ['1', '16', '47'],
            validation: null,
            searchable: true,
            reportable: true,
            importable: true,
            exportable: true,
            displayOrder: 100,
            active: true,
        });
        expect(queryCustomFieldDefinitions).toHaveBeenCalledWith({
            page: 1,
            rows: 100,
            orderBy: 'display_order,ASC',
        });
        expect(component.customFieldSaving).toBe(false);
        expect(component.customFieldForm.fieldKey).toBe('');
    });

    it('surfaces JSON:API errors when custom field save fails', () => {
        createCustomFieldDefinition.mockReturnValue(
            throwError(() => ({
                error: {
                    errors: [
                        {
                            detail: 'Field key already exists.',
                        },
                    ],
                },
            }))
        );
        const fixture = TestBed.createComponent(AdmissionsSettingsComponent);
        const component = fixture.componentInstance;
        component.customFieldForm = {
            ...component.customFieldForm,
            fieldKey: 'county_code',
            label: 'County of Residence',
        };

        component.saveCustomField();

        expect(component.errorMessage).toBe('Field key already exists.');
        expect(component.customFieldSaving).toBe(false);
        expect(component.customFieldForm.fieldKey).toBe('county_code');
    });
});

function result<T>(items: T[]): PaginatedResult<T> {
    return {
        items,
        total: items.length,
        page: 1,
        rowsPerPage: 50,
    };
}

const commerceProgram: Program = {
    id: 'program-commerce',
    externalSISID: 'KUCCPS-BCOM-2026-QA',
    name: 'Bachelor of Commerce',
    code: 'BCOM-QA',
    description: 'QA fixture commerce programme',
    degreeLevel: 'BACHELOR',
    active: true,
    dateCreated: '2026-06-01T08:00:00Z',
    dateUpdated: '2026-06-01T08:00:00Z',
};

const septemberTerm: AcademicTerm = {
    id: 'term-september',
    externalSISID: 'INTAKE-SEP-2026-QA',
    name: 'September 2026 Intake',
    code: 'SEP2026-QA',
    termType: 'MAIN',
    startDate: '2026-09-07T00:00:00Z',
    endDate: '2026-12-18T00:00:00Z',
    applicationStartDate: '2026-04-01T00:00:00Z',
    applicationDeadline: '2026-08-21T00:00:00Z',
    active: true,
    dateCreated: '2026-06-01T08:00:00Z',
    dateUpdated: '2026-06-01T08:00:00Z',
};

const countyField: CustomFieldDefinition = {
    id: 'field-county-code',
    owner: 'CONSTITUENT',
    fieldKey: 'county_code',
    label: 'County of Residence',
    description: 'Kenya county code used for admissions reporting.',
    dataType: 'SELECT',
    required: true,
    options: ['1', '16', '47'],
    validation: '^[0-9]{1,2}$',
    searchable: true,
    reportable: true,
    importable: true,
    exportable: true,
    displayOrder: 10,
    active: true,
    dateCreated: '2026-06-01T08:00:00Z',
    dateUpdated: '2026-06-01T08:00:00Z',
};

const importBatch: ImportBatch = {
    id: 'import-constituents-template',
    source: 'MANUAL_UPLOAD',
    fileType: 'CSV',
    target: 'CONSTITUENTS',
    status: 'PREVIEWED',
    fileName: 'admissions-constituents-template.csv',
    storageKey: 'seed/imports/admissions/constituents-template.csv',
    uploadedByID: 'admin-user',
    totalRows: 0,
    validRows: 0,
    invalidRows: 0,
    duplicateRows: 0,
    fieldMapping: {},
    validationSummary: 'Seeded blank constituent import template.',
    dateCreated: '2026-06-01T08:00:00Z',
    dateUpdated: '2026-06-01T08:00:00Z',
};

const leadScoreRule: LeadScoreRule = {
    id: 'lead-rule-applicant',
    name: 'Applicant lifecycle stage',
    description: 'Applicants are warmer than prospects.',
    criteria: [
        {
            field: 'lifecycle_stage',
            operator: 'EQ',
            values: ['APPLICANT'],
        },
    ],
    points: 30,
    active: true,
    priority: 10,
    dateCreated: '2026-06-01T08:00:00Z',
    dateUpdated: '2026-06-01T08:00:00Z',
};

const formTemplate: ApplicationFormTemplate = {
    id: 'template-freshman',
    programID: 'program-commerce',
    academicTermID: 'term-september',
    applicationType: 'KUCCPS_PLACEMENT',
    name: 'KUCCPS placement undergraduate form',
    description: 'Default active form for placed undergraduate applicants.',
    version: 1,
    requiredFields: [],
    checklistItems: [],
    active: true,
    priority: 10,
    dateCreated: '2026-06-01T08:00:00Z',
    dateUpdated: '2026-06-01T08:00:00Z',
};
