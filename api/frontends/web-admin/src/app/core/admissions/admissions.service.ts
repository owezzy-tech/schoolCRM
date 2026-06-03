import { HttpClient, HttpParams } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import {
    AcademicTerm,
    AcademicTermQuery,
    AdmissionsDocument,
    AdmissionsDocumentQuery,
    AdmissionsDocumentRequest,
    AdmissionsDocumentVerificationRequest,
    Application,
    ApplicationFormTemplate,
    ApplicationFormTemplateQuery,
    ApplicationFormTemplateRequest,
    ApplicationQuery,
    ApplicationRequest,
    ApplicationTransitionRequest,
    ChecklistItem,
    ChecklistItemQuery,
    ChecklistItemRequest,
    CustomFieldDefinition,
    CustomFieldDefinitionQuery,
    CustomFieldDefinitionRequest,
    CustomFieldValue,
    CustomFieldValueRequest,
    ImportBatch,
    ImportBatchQuery,
    ImportBatchRequest,
    ImportInvalidRow,
    ImportInvalidRowQuery,
    ImportInvalidRowsRequest,
    Inquiry,
    InquiryRequest,
    LeadScore,
    LeadScoreQuery,
    LeadScoreRule,
    LeadScoreRuleQuery,
    LeadScoreRuleRequest,
    Program,
    ProgramQuery,
} from 'app/core/admissions/admissions.types';
import {
    JsonApiCollectionDocument,
    JsonApiDocument,
    PaginatedResult,
    unwrapJsonApiCollection,
    unwrapJsonApiResource,
} from 'app/core/api/json-api';
import { map, Observable, ReplaySubject, tap } from 'rxjs';

@Injectable({ providedIn: 'root' })
export class AdmissionsService {
    private readonly httpClient = inject(HttpClient);
    private readonly scoresSubject = new ReplaySubject<
        PaginatedResult<LeadScore>
    >(1);
    private readonly rulesSubject = new ReplaySubject<
        PaginatedResult<LeadScoreRule>
    >(1);
    private readonly templatesSubject = new ReplaySubject<
        PaginatedResult<ApplicationFormTemplate>
    >(1);
    private readonly applicantApplicationsSubject = new ReplaySubject<
        PaginatedResult<Application>
    >(1);
    private readonly applicantProgramsSubject = new ReplaySubject<
        PaginatedResult<Program>
    >(1);
    private readonly applicantAcademicTermsSubject = new ReplaySubject<
        PaginatedResult<AcademicTerm>
    >(1);
    private readonly applicantTemplatesSubject = new ReplaySubject<
        PaginatedResult<ApplicationFormTemplate>
    >(1);
    private readonly checklistItemsSubject = new ReplaySubject<
        PaginatedResult<ChecklistItem>
    >(1);
    private readonly documentsSubject = new ReplaySubject<
        PaginatedResult<AdmissionsDocument>
    >(1);
    private readonly customFieldDefinitionsSubject = new ReplaySubject<
        PaginatedResult<CustomFieldDefinition>
    >(1);
    private readonly importBatchesSubject = new ReplaySubject<
        PaginatedResult<ImportBatch>
    >(1);
    private readonly importInvalidRowsSubject = new ReplaySubject<
        PaginatedResult<ImportInvalidRow>
    >(1);

    readonly scores$ = this.scoresSubject.asObservable();
    readonly rules$ = this.rulesSubject.asObservable();
    readonly templates$ = this.templatesSubject.asObservable();
    readonly applicantApplications$ =
        this.applicantApplicationsSubject.asObservable();
    readonly applicantPrograms$ = this.applicantProgramsSubject.asObservable();
    readonly applicantAcademicTerms$ =
        this.applicantAcademicTermsSubject.asObservable();
    readonly applicantTemplates$ =
        this.applicantTemplatesSubject.asObservable();
    readonly checklistItems$ = this.checklistItemsSubject.asObservable();
    readonly documents$ = this.documentsSubject.asObservable();
    readonly customFieldDefinitions$ =
        this.customFieldDefinitionsSubject.asObservable();
    readonly importBatches$ = this.importBatchesSubject.asObservable();
    readonly importInvalidRows$ = this.importInvalidRowsSubject.asObservable();

    queryLeadScores(
        query: LeadScoreQuery = {}
    ): Observable<PaginatedResult<LeadScore>> {
        return this.httpClient
            .get<
                JsonApiCollectionDocument<LeadScore>
            >('/v1/admissions/lead-scores', { params: this.queryParams(query) })
            .pipe(
                map(unwrapJsonApiCollection),
                tap((result) => this.scoresSubject.next(result))
            );
    }

    queryLeadScoreRules(
        query: LeadScoreRuleQuery = {}
    ): Observable<PaginatedResult<LeadScoreRule>> {
        return this.httpClient
            .get<
                JsonApiCollectionDocument<LeadScoreRule>
            >('/v1/admissions/lead-score-rules', { params: this.queryParams(query) })
            .pipe(
                map(unwrapJsonApiCollection),
                tap((result) => this.rulesSubject.next(result))
            );
    }

    createLeadScoreRule(
        request: LeadScoreRuleRequest
    ): Observable<LeadScoreRule> {
        return this.httpClient
            .post<
                JsonApiDocument<LeadScoreRule>
            >('/v1/admissions/lead-score-rules', request)
            .pipe(map(unwrapJsonApiResource));
    }

    updateLeadScoreRule(
        ruleID: string,
        request: LeadScoreRuleRequest
    ): Observable<LeadScoreRule> {
        return this.httpClient
            .put<
                JsonApiDocument<LeadScoreRule>
            >(`/v1/admissions/lead-score-rules/${ruleID}`, request)
            .pipe(map(unwrapJsonApiResource));
    }

    createInquiry(request: InquiryRequest): Observable<Inquiry> {
        return this.httpClient
            .post<JsonApiDocument<Inquiry>>('/v1/admissions/inquiries', request)
            .pipe(map(unwrapJsonApiResource));
    }

    queryApplicationFormTemplates(
        query: ApplicationFormTemplateQuery = {}
    ): Observable<PaginatedResult<ApplicationFormTemplate>> {
        return this.httpClient
            .get<
                JsonApiCollectionDocument<ApplicationFormTemplate>
            >('/v1/admissions/application-form-templates', { params: this.queryParams(query) })
            .pipe(
                map(unwrapJsonApiCollection),
                tap((result) => this.templatesSubject.next(result))
            );
    }

    createApplicationFormTemplate(
        request: ApplicationFormTemplateRequest
    ): Observable<ApplicationFormTemplate> {
        return this.httpClient
            .post<
                JsonApiDocument<ApplicationFormTemplate>
            >('/v1/admissions/application-form-templates', request)
            .pipe(map(unwrapJsonApiResource));
    }

    updateApplicationFormTemplate(
        templateID: string,
        request: ApplicationFormTemplateRequest
    ): Observable<ApplicationFormTemplate> {
        return this.httpClient
            .put<
                JsonApiDocument<ApplicationFormTemplate>
            >(`/v1/admissions/application-form-templates/${templateID}`, request)
            .pipe(map(unwrapJsonApiResource));
    }

    queryApplicantPrograms(
        query: ProgramQuery = {}
    ): Observable<PaginatedResult<Program>> {
        return this.httpClient
            .get<
                JsonApiCollectionDocument<Program>
            >('/v1/admissions/applicant/programs', { params: this.queryParams(query) })
            .pipe(
                map(unwrapJsonApiCollection),
                tap((result) => this.applicantProgramsSubject.next(result))
            );
    }

    queryApplicantAcademicTerms(
        query: AcademicTermQuery = {}
    ): Observable<PaginatedResult<AcademicTerm>> {
        return this.httpClient
            .get<
                JsonApiCollectionDocument<AcademicTerm>
            >('/v1/admissions/applicant/academic-terms', { params: this.queryParams(query) })
            .pipe(
                map(unwrapJsonApiCollection),
                tap((result) => this.applicantAcademicTermsSubject.next(result))
            );
    }

    queryApplicantApplicationFormTemplates(
        query: ApplicationFormTemplateQuery = {}
    ): Observable<PaginatedResult<ApplicationFormTemplate>> {
        return this.httpClient
            .get<
                JsonApiCollectionDocument<ApplicationFormTemplate>
            >('/v1/admissions/applicant/application-form-templates', { params: this.queryParams(query) })
            .pipe(
                map(unwrapJsonApiCollection),
                tap((result) => this.applicantTemplatesSubject.next(result))
            );
    }

    queryApplicantApplications(
        query: ApplicationQuery = {}
    ): Observable<PaginatedResult<Application>> {
        return this.httpClient
            .get<
                JsonApiCollectionDocument<Application>
            >('/v1/admissions/applicant/applications', { params: this.queryParams(query) })
            .pipe(
                map(unwrapJsonApiCollection),
                tap((result) => this.applicantApplicationsSubject.next(result))
            );
    }

    getApplicantApplication(applicationID: string): Observable<Application> {
        return this.httpClient
            .get<
                JsonApiDocument<Application>
            >(`/v1/admissions/applicant/applications/${applicationID}`)
            .pipe(map(unwrapJsonApiResource));
    }

    createApplicantApplication(
        request: ApplicationRequest
    ): Observable<Application> {
        return this.httpClient
            .post<
                JsonApiDocument<Application>
            >('/v1/admissions/applicant/applications', request)
            .pipe(map(unwrapJsonApiResource));
    }

    updateApplicantApplication(
        applicationID: string,
        request: ApplicationRequest
    ): Observable<Application> {
        return this.httpClient
            .put<
                JsonApiDocument<Application>
            >(`/v1/admissions/applicant/applications/${applicationID}`, request)
            .pipe(map(unwrapJsonApiResource));
    }

    transitionApplicantApplication(
        applicationID: string,
        request: ApplicationTransitionRequest
    ): Observable<Application> {
        return this.httpClient
            .post<
                JsonApiDocument<Application>
            >(`/v1/admissions/applicant/applications/${applicationID}/transitions`, request)
            .pipe(map(unwrapJsonApiResource));
    }

    queryCustomFieldDefinitions(
        query: CustomFieldDefinitionQuery = {}
    ): Observable<PaginatedResult<CustomFieldDefinition>> {
        return this.httpClient
            .get<
                JsonApiCollectionDocument<CustomFieldDefinition>
            >('/v1/admissions/custom-field-definitions', { params: this.queryParams(query) })
            .pipe(
                map(unwrapJsonApiCollection),
                tap((result) => this.customFieldDefinitionsSubject.next(result))
            );
    }

    createCustomFieldDefinition(
        request: CustomFieldDefinitionRequest
    ): Observable<CustomFieldDefinition> {
        return this.httpClient
            .post<
                JsonApiDocument<CustomFieldDefinition>
            >('/v1/admissions/custom-field-definitions', request)
            .pipe(map(unwrapJsonApiResource));
    }

    updateCustomFieldDefinition(
        definitionID: string,
        request: CustomFieldDefinitionRequest
    ): Observable<CustomFieldDefinition> {
        return this.httpClient
            .put<
                JsonApiDocument<CustomFieldDefinition>
            >(`/v1/admissions/custom-field-definitions/${definitionID}`, request)
            .pipe(map(unwrapJsonApiResource));
    }

    setCustomFieldValue(
        request: CustomFieldValueRequest
    ): Observable<CustomFieldValue> {
        return this.httpClient
            .put<
                JsonApiDocument<CustomFieldValue>
            >('/v1/admissions/custom-field-values', request)
            .pipe(map(unwrapJsonApiResource));
    }

    queryChecklistItems(
        applicationID: string,
        query: ChecklistItemQuery = {}
    ): Observable<PaginatedResult<ChecklistItem>> {
        return this.httpClient
            .get<
                JsonApiCollectionDocument<ChecklistItem>
            >(`/v1/admissions/applications/${applicationID}/checklist-items`, { params: this.queryParams(query) })
            .pipe(
                map(unwrapJsonApiCollection),
                tap((result) => this.checklistItemsSubject.next(result))
            );
    }

    queryApplicantChecklistItems(
        applicationID: string,
        query: ChecklistItemQuery = {}
    ): Observable<PaginatedResult<ChecklistItem>> {
        return this.httpClient
            .get<
                JsonApiCollectionDocument<ChecklistItem>
            >(`/v1/admissions/applicant/applications/${applicationID}/checklist-items`, { params: this.queryParams(query) })
            .pipe(
                map(unwrapJsonApiCollection),
                tap((result) => this.checklistItemsSubject.next(result))
            );
    }

    createChecklistItem(
        applicationID: string,
        request: ChecklistItemRequest
    ): Observable<ChecklistItem> {
        return this.httpClient
            .post<
                JsonApiDocument<ChecklistItem>
            >(`/v1/admissions/applications/${applicationID}/checklist-items`, request)
            .pipe(map(unwrapJsonApiResource));
    }

    updateChecklistItem(
        applicationID: string,
        checklistItemID: string,
        request: ChecklistItemRequest
    ): Observable<ChecklistItem> {
        return this.httpClient
            .put<
                JsonApiDocument<ChecklistItem>
            >(`/v1/admissions/applications/${applicationID}/checklist-items/${checklistItemID}`, request)
            .pipe(map(unwrapJsonApiResource));
    }

    queryDocuments(
        applicationID: string,
        query: AdmissionsDocumentQuery = {}
    ): Observable<PaginatedResult<AdmissionsDocument>> {
        return this.httpClient
            .get<
                JsonApiCollectionDocument<AdmissionsDocument>
            >(`/v1/admissions/applications/${applicationID}/documents`, { params: this.queryParams(query) })
            .pipe(
                map(unwrapJsonApiCollection),
                tap((result) => this.documentsSubject.next(result))
            );
    }

    createDocument(
        applicationID: string,
        request: AdmissionsDocumentRequest
    ): Observable<AdmissionsDocument> {
        return this.httpClient
            .post<
                JsonApiDocument<AdmissionsDocument>
            >(`/v1/admissions/applications/${applicationID}/documents`, request)
            .pipe(map(unwrapJsonApiResource));
    }

    queryApplicantDocuments(
        applicationID: string,
        query: AdmissionsDocumentQuery = {}
    ): Observable<PaginatedResult<AdmissionsDocument>> {
        return this.httpClient
            .get<
                JsonApiCollectionDocument<AdmissionsDocument>
            >(`/v1/admissions/applicant/applications/${applicationID}/documents`, { params: this.queryParams(query) })
            .pipe(
                map(unwrapJsonApiCollection),
                tap((result) => this.documentsSubject.next(result))
            );
    }

    createApplicantDocument(
        applicationID: string,
        request: AdmissionsDocumentRequest
    ): Observable<AdmissionsDocument> {
        return this.httpClient
            .post<
                JsonApiDocument<AdmissionsDocument>
            >(`/v1/admissions/applicant/applications/${applicationID}/documents`, request)
            .pipe(map(unwrapJsonApiResource));
    }

    verifyDocument(
        documentID: string,
        request: AdmissionsDocumentVerificationRequest
    ): Observable<AdmissionsDocument> {
        return this.httpClient
            .put<
                JsonApiDocument<AdmissionsDocument>
            >(`/v1/admissions/documents/${documentID}/verification`, request)
            .pipe(map(unwrapJsonApiResource));
    }

    downloadDocument(documentID: string): Observable<AdmissionsDocument> {
        return this.httpClient
            .post<
                JsonApiDocument<AdmissionsDocument>
            >(`/v1/admissions/documents/${documentID}/download`, {})
            .pipe(map(unwrapJsonApiResource));
    }

    queryImportBatches(
        query: ImportBatchQuery = {}
    ): Observable<PaginatedResult<ImportBatch>> {
        return this.httpClient
            .get<
                JsonApiCollectionDocument<ImportBatch>
            >('/v1/admissions/import-batches', { params: this.queryParams(query) })
            .pipe(
                map(unwrapJsonApiCollection),
                tap((result) => this.importBatchesSubject.next(result))
            );
    }

    createImportBatch(request: ImportBatchRequest): Observable<ImportBatch> {
        return this.httpClient
            .post<
                JsonApiDocument<ImportBatch>
            >('/v1/admissions/import-batches', request)
            .pipe(map(unwrapJsonApiResource));
    }

    queryImportInvalidRows(
        importBatchID: string,
        query: ImportInvalidRowQuery = {}
    ): Observable<PaginatedResult<ImportInvalidRow>> {
        return this.httpClient
            .get<
                JsonApiCollectionDocument<ImportInvalidRow>
            >(`/v1/admissions/import-batches/${importBatchID}/invalid-rows`, { params: this.queryParams(query) })
            .pipe(
                map(unwrapJsonApiCollection),
                tap((result) => this.importInvalidRowsSubject.next(result))
            );
    }

    createImportInvalidRows(
        importBatchID: string,
        request: ImportInvalidRowsRequest
    ): Observable<PaginatedResult<ImportInvalidRow>> {
        return this.httpClient
            .post<
                JsonApiCollectionDocument<ImportInvalidRow>
            >(`/v1/admissions/import-batches/${importBatchID}/invalid-rows`, request)
            .pipe(map(unwrapJsonApiCollection));
    }

    downloadImportInvalidRows(
        importBatchID: string
    ): Observable<PaginatedResult<ImportInvalidRow>> {
        return this.httpClient
            .post<
                JsonApiCollectionDocument<ImportInvalidRow>
            >(`/v1/admissions/import-batches/${importBatchID}/invalid-rows/download`, {})
            .pipe(map(unwrapJsonApiCollection));
    }

    private queryParams(
        query:
            | LeadScoreQuery
            | LeadScoreRuleQuery
            | ProgramQuery
            | AcademicTermQuery
            | ApplicationFormTemplateQuery
            | ApplicationQuery
            | CustomFieldDefinitionQuery
            | ChecklistItemQuery
            | AdmissionsDocumentQuery
            | ImportBatchQuery
            | ImportInvalidRowQuery
    ): HttpParams {
        let params = new HttpParams();

        for (const [key, value] of Object.entries(query)) {
            if (value !== undefined && value !== '') {
                params = params.set(key, String(value));
            }
        }

        return params;
    }
}
