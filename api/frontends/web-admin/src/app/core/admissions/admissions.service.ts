import { HttpClient, HttpParams } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import {
    AdmissionsDocument,
    AdmissionsDocumentQuery,
    AdmissionsDocumentRequest,
    AdmissionsDocumentVerificationRequest,
    ApplicationFormTemplate,
    ApplicationFormTemplateQuery,
    ApplicationFormTemplateRequest,
    ChecklistItem,
    ChecklistItemQuery,
    ChecklistItemRequest,
    Inquiry,
    InquiryRequest,
    LeadScore,
    LeadScoreQuery,
    LeadScoreRule,
    LeadScoreRuleQuery,
    LeadScoreRuleRequest,
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
    private readonly checklistItemsSubject = new ReplaySubject<
        PaginatedResult<ChecklistItem>
    >(1);
    private readonly documentsSubject = new ReplaySubject<
        PaginatedResult<AdmissionsDocument>
    >(1);

    readonly scores$ = this.scoresSubject.asObservable();
    readonly rules$ = this.rulesSubject.asObservable();
    readonly templates$ = this.templatesSubject.asObservable();
    readonly checklistItems$ = this.checklistItemsSubject.asObservable();
    readonly documents$ = this.documentsSubject.asObservable();

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

    private queryParams(
        query:
            | LeadScoreQuery
            | LeadScoreRuleQuery
            | ApplicationFormTemplateQuery
            | ChecklistItemQuery
            | AdmissionsDocumentQuery
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
