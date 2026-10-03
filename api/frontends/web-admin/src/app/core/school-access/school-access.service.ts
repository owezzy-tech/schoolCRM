import { HttpClient } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import {
    JsonApiCollectionDocument,
    JsonApiDocument,
    unwrapJsonApiCollection,
    unwrapJsonApiResource,
} from 'app/core/api/json-api';
import { map, Observable } from 'rxjs';

import {
    Department,
    GrantRequest,
    Membership,
    School,
} from './school-access.types';

/** Client for the school access API (docs/school-access-api.md). */
@Injectable({ providedIn: 'root' })
export class SchoolAccessService {
    private readonly httpClient = inject(HttpClient);

    /** The caller's active memberships; used only to decide which actions to offer. */
    ownMemberships(): Observable<Membership[]> {
        return this.list<Membership>('/v1/me/school-memberships');
    }

    schools(): Observable<School[]> {
        return this.list<School>('/v1/schools');
    }

    createSchool(name: string): Observable<School> {
        return this.httpClient
            .post<JsonApiDocument<School>>('/v1/schools', { name })
            .pipe(map(unwrapJsonApiResource));
    }

    departments(schoolID: string): Observable<Department[]> {
        return this.list<Department>(`/v1/schools/${schoolID}/departments`);
    }

    createDepartment(schoolID: string, name: string): Observable<Department> {
        return this.httpClient
            .post<JsonApiDocument<Department>>(
                `/v1/schools/${schoolID}/departments`,
                { name }
            )
            .pipe(map(unwrapJsonApiResource));
    }

    memberships(schoolID: string): Observable<Membership[]> {
        return this.list<Membership>(`/v1/schools/${schoolID}/memberships`);
    }

    grant(schoolID: string, request: GrantRequest): Observable<Membership> {
        return this.httpClient
            .post<JsonApiDocument<Membership>>(
                `/v1/schools/${schoolID}/memberships`,
                request
            )
            .pipe(map(unwrapJsonApiResource));
    }

    revoke(schoolID: string, membershipID: string): Observable<void> {
        return this.httpClient.delete<void>(
            `/v1/schools/${schoolID}/memberships/${membershipID}`
        );
    }

    private list<T extends { id: string }>(url: string): Observable<T[]> {
        return this.httpClient
            .get<JsonApiCollectionDocument<T>>(url)
            .pipe(map((document) => unwrapJsonApiCollection(document).items));
    }
}
