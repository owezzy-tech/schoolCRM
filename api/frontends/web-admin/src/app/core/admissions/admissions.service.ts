import { HttpClient, HttpParams } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import {
    JsonApiCollectionDocument,
    JsonApiDocument,
    PaginatedResult,
    unwrapJsonApiCollection,
    unwrapJsonApiResource,
} from 'app/core/api/json-api';
import {
    LeadScore,
    LeadScoreQuery,
    LeadScoreRule,
    LeadScoreRuleQuery,
    LeadScoreRuleRequest,
} from 'app/core/admissions/admissions.types';
import { map, Observable, ReplaySubject, tap } from 'rxjs';

@Injectable({ providedIn: 'root' })
export class AdmissionsService {
    private readonly httpClient = inject(HttpClient);
    private readonly scoresSubject = new ReplaySubject<PaginatedResult<LeadScore>>(
        1
    );
    private readonly rulesSubject = new ReplaySubject<
        PaginatedResult<LeadScoreRule>
    >(1);

    readonly scores$ = this.scoresSubject.asObservable();
    readonly rules$ = this.rulesSubject.asObservable();

    queryLeadScores(
        query: LeadScoreQuery = {}
    ): Observable<PaginatedResult<LeadScore>> {
        return this.httpClient
            .get<JsonApiCollectionDocument<LeadScore>>(
                '/v1/admissions/lead-scores',
                { params: this.queryParams(query) }
            )
            .pipe(
                map(unwrapJsonApiCollection),
                tap((result) => this.scoresSubject.next(result))
            );
    }

    queryLeadScoreRules(
        query: LeadScoreRuleQuery = {}
    ): Observable<PaginatedResult<LeadScoreRule>> {
        return this.httpClient
            .get<JsonApiCollectionDocument<LeadScoreRule>>(
                '/v1/admissions/lead-score-rules',
                { params: this.queryParams(query) }
            )
            .pipe(
                map(unwrapJsonApiCollection),
                tap((result) => this.rulesSubject.next(result))
            );
    }

    createLeadScoreRule(
        request: LeadScoreRuleRequest
    ): Observable<LeadScoreRule> {
        return this.httpClient
            .post<JsonApiDocument<LeadScoreRule>>(
                '/v1/admissions/lead-score-rules',
                request
            )
            .pipe(map(unwrapJsonApiResource));
    }

    updateLeadScoreRule(
        ruleID: string,
        request: LeadScoreRuleRequest
    ): Observable<LeadScoreRule> {
        return this.httpClient
            .put<JsonApiDocument<LeadScoreRule>>(
                `/v1/admissions/lead-score-rules/${ruleID}`,
                request
            )
            .pipe(map(unwrapJsonApiResource));
    }

    private queryParams(
        query: LeadScoreQuery | LeadScoreRuleQuery
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
