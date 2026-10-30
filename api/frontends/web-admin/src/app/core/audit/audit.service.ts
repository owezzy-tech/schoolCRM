import { HttpClient, HttpParams } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import {
    JsonApiCollectionDocument,
    PaginatedResult,
    unwrapJsonApiCollection,
} from 'app/core/api/json-api';
import { Observable, ReplaySubject, map, tap } from 'rxjs';

import { AuditEvent, AuditEventQuery } from './audit.types';

@Injectable({ providedIn: 'root' })
export class AuditService {
    private readonly httpClient = inject(HttpClient);
    private readonly auditEventsSubject = new ReplaySubject<
        PaginatedResult<AuditEvent>
    >(1);

    readonly auditEvents$ = this.auditEventsSubject.asObservable();

    queryAuditEvents(
        query: AuditEventQuery = {}
    ): Observable<PaginatedResult<AuditEvent>> {
        return this.httpClient
            .get<JsonApiCollectionDocument<AuditEvent>>('/v1/audits', {
                params: this.queryParams(query),
            })
            .pipe(
                map(unwrapJsonApiCollection),
                tap((result) => this.auditEventsSubject.next(result))
            );
    }

    private queryParams(query: AuditEventQuery): HttpParams {
        let params = new HttpParams();

        for (const [key, value] of Object.entries(query)) {
            if (value !== undefined && value !== '') {
                params = params.set(key, String(value));
            }
        }

        return params;
    }
}
