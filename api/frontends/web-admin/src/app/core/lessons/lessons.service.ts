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
    LessonDecision,
    LessonDraft,
    LessonPlan,
    LessonVersion,
} from './lessons.types';

/** Client for the lesson publishing API (docs/lesson-publishing-api.md). */
@Injectable({ providedIn: 'root' })
export class LessonsService {
    private readonly httpClient = inject(HttpClient);

    plans(schoolID: string, departmentID: string): Observable<LessonPlan[]> {
        return this.httpClient
            .get<
                JsonApiCollectionDocument<LessonPlan>
            >(this.departmentPath(schoolID, departmentID))
            .pipe(map((document) => unwrapJsonApiCollection(document).items));
    }

    create(
        schoolID: string,
        departmentID: string,
        draft: LessonDraft
    ): Observable<LessonPlan> {
        return this.httpClient
            .post<
                JsonApiDocument<LessonPlan>
            >(this.departmentPath(schoolID, departmentID), draft)
            .pipe(map(unwrapJsonApiResource));
    }

    plan(planID: string): Observable<LessonPlan> {
        return this.httpClient
            .get<JsonApiDocument<LessonPlan>>(`/v1/lessons/${planID}`)
            .pipe(map(unwrapJsonApiResource));
    }

    /** Versions newest first. */
    versions(planID: string): Observable<LessonVersion[]> {
        return this.httpClient
            .get<
                JsonApiCollectionDocument<LessonVersion>
            >(`/v1/lessons/${planID}/versions`)
            .pipe(map((document) => unwrapJsonApiCollection(document).items));
    }

    /** Saves a new draft version; the API answers 409 when baseVersion is stale. */
    revise(
        planID: string,
        baseVersion: number,
        draft: LessonDraft
    ): Observable<LessonVersion> {
        return this.httpClient
            .post<
                JsonApiDocument<LessonVersion>
            >(`/v1/lessons/${planID}/versions`, { baseVersion, ...draft })
            .pipe(map(unwrapJsonApiResource));
    }

    /** Stable request ID makes an ambiguous creation reply safe to retry. */
    reuse(
        planID: string,
        version: number,
        requestID: string,
        title: string
    ): Observable<LessonPlan> {
        return this.httpClient
            .post<
                JsonApiDocument<LessonPlan>
            >(`/v1/lessons/${planID}/reuse`, { version, requestID, title })
            .pipe(map(unwrapJsonApiResource));
    }

    submit(planID: string, version: number): Observable<LessonVersion> {
        return this.command(planID, version, 'submit');
    }

    review(
        planID: string,
        version: number,
        decision: LessonDecision,
        feedback: string
    ): Observable<LessonVersion> {
        return this.command(planID, version, 'review', { decision, feedback });
    }

    approve(
        planID: string,
        version: number,
        decision: LessonDecision,
        feedback: string
    ): Observable<LessonVersion> {
        return this.command(planID, version, 'approval', {
            decision,
            feedback,
        });
    }

    publish(planID: string, version: number): Observable<LessonVersion> {
        return this.command(planID, version, 'publish');
    }

    private command(
        planID: string,
        version: number,
        step: 'submit' | 'review' | 'approval' | 'publish',
        body: { decision: LessonDecision; feedback: string } | null = null
    ): Observable<LessonVersion> {
        return this.httpClient
            .post<
                JsonApiDocument<LessonVersion>
            >(`/v1/lessons/${planID}/versions/${version}/${step}`, body)
            .pipe(map(unwrapJsonApiResource));
    }

    private departmentPath(schoolID: string, departmentID: string): string {
        return `/v1/schools/${schoolID}/departments/${departmentID}/lessons`;
    }
}
