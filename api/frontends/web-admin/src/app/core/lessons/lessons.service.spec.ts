import { provideHttpClient } from '@angular/common/http';
import {
    HttpTestingController,
    provideHttpClientTesting,
} from '@angular/common/http/testing';
import { TestBed } from '@angular/core/testing';
import { afterEach, beforeEach, describe, expect, it } from 'vitest';

import { LessonsService } from './lessons.service';

describe('LessonsService', () => {
    let service: LessonsService;
    let http: HttpTestingController;

    beforeEach(() => {
        TestBed.configureTestingModule({
            providers: [provideHttpClient(), provideHttpClientTesting()],
        });
        service = TestBed.inject(LessonsService);
        http = TestBed.inject(HttpTestingController);
    });

    afterEach(() => http.verify());

    it('unwraps the department plan collection', () => {
        let titles: string[] = [];
        service.plans('school-1', 'sciences').subscribe((plans) => {
            titles = plans.map((plan) => `${plan.id}:${plan.title}`);
        });

        http.expectOne('/v1/schools/school-1/departments/sciences/lessons').flush({
            jsonapi: { version: '1.1' },
            data: [{ id: 'plan-1', type: 'lessonplan', attributes: { title: 'Forces' } }],
            meta: { total: 1 },
        });

        expect(titles).toEqual(['plan-1:Forces']);
    });

    it('revises from the base version it was edited from', () => {
        const draft = {
            title: 'Forces v2',
            changeSummary: 'Adds Hooke',
            content: { objectives: [], materials: [], activities: [], assessment: '' },
        };
        service.revise('plan-1', 1, draft).subscribe();

        const request = http.expectOne('/v1/lessons/plan-1/versions');
        expect(request.request.method).toBe('POST');
        expect(request.request.body).toEqual({ baseVersion: 1, ...draft });
        request.flush({ data: { id: 'plan-1:2', type: 'lessonversion', attributes: { version: 2 } } });
    });

    it('sends decisions to the exact version', () => {
        service.approve('plan-1', 3, 'request_changes', 'Add assessment').subscribe();

        const request = http.expectOne('/v1/lessons/plan-1/versions/3/approval');
        expect(request.request.body).toEqual({
            decision: 'request_changes',
            feedback: 'Add assessment',
        });
        request.flush({ data: { id: 'plan-1:3', type: 'lessonversion', attributes: {} } });
    });

    it('submits and publishes without a body', () => {
        service.submit('plan-1', 1).subscribe();
        service.publish('plan-1', 1).subscribe();

        for (const step of ['submit', 'publish']) {
            const request = http.expectOne(`/v1/lessons/plan-1/versions/1/${step}`);
            expect(request.request.body).toBeNull();
            request.flush({ data: { id: 'plan-1:1', type: 'lessonversion', attributes: {} } });
        }
    });
});
