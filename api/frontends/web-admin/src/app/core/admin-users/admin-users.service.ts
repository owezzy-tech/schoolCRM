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
    AdminUser,
    UpdateUserRequest,
    UpdateUserRolesRequest,
    UserQuery,
} from 'app/core/admin-users/admin-users.types';
import { map, Observable, ReplaySubject, tap } from 'rxjs';

@Injectable({ providedIn: 'root' })
export class AdminUsersService {
    private readonly httpClient = inject(HttpClient);
    private readonly usersSubject = new ReplaySubject<PaginatedResult<AdminUser>>(
        1
    );

    readonly users$ = this.usersSubject.asObservable();

    query(query: UserQuery = {}): Observable<PaginatedResult<AdminUser>> {
        return this.httpClient
            .get<JsonApiCollectionDocument<AdminUser>>('/v1/users', {
                params: this.queryParams(query),
            })
            .pipe(
                map(unwrapJsonApiCollection),
                tap((result) => this.usersSubject.next(result))
            );
    }

    updateRoles(
        userID: string,
        request: UpdateUserRolesRequest
    ): Observable<AdminUser> {
        return this.httpClient
            .put<JsonApiDocument<AdminUser>>(`/v1/users/role/${userID}`, request)
            .pipe(map(unwrapJsonApiResource));
    }

    update(userID: string, request: UpdateUserRequest): Observable<AdminUser> {
        return this.httpClient
            .put<JsonApiDocument<AdminUser>>(`/v1/users/${userID}`, request)
            .pipe(map(unwrapJsonApiResource));
    }

    delete(userID: string): Observable<void> {
        return this.httpClient.delete<void>(`/v1/users/${userID}`);
    }

    private queryParams(query: UserQuery): HttpParams {
        let params = new HttpParams();
        const entries: [string, string | number | undefined][] = [
            ['page', query.page],
            ['rows', query.rows],
            ['orderBy', query.orderBy],
            ['name', query.name],
            ['email', query.email],
        ];

        for (const [key, value] of entries) {
            if (value !== undefined && value !== '') {
                params = params.set(key, String(value));
            }
        }

        return params;
    }
}
