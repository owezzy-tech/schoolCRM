import { HttpErrorResponse } from '@angular/common/http';
import { TestBed } from '@angular/core/testing';
import { MatSnackBar } from '@angular/material/snack-bar';
import { ActivatedRoute, convertToParamMap } from '@angular/router';
import { LessonsService } from 'app/core/lessons/lessons.service';
import { LessonPlan, LessonVersion } from 'app/core/lessons/lessons.types';
import { SchoolAccessService } from 'app/core/school-access/school-access.service';
import { UserService } from 'app/core/user/user.service';
import { of, throwError } from 'rxjs';
import { beforeEach, describe, expect, it, vi } from 'vitest';

import { versionRelation } from '../version/lesson-version.component';
import { LessonPlanComponent } from './lesson-plan.component';

const plan: LessonPlan = {
    id: 'plan-1',
    schoolID: 'school-1',
    departmentID: 'sciences',
    authorID: 'grace',
    title: 'Forces',
    status: 'hod_review',
    currentVersion: 2,
    publishedVersion: 1,
    dateCreated: '2026-10-01T08:00:00Z',
    dateUpdated: '2026-10-02T08:00:00Z',
};

const current: LessonVersion = {
    id: 'plan-1:2',
    planID: 'plan-1',
    version: 2,
    title: 'Forces',
    content: {},
    changeSummary: '',
    authorID: 'grace',
    status: 'hod_review',
    reviewerID: null,
    reviewedAt: null,
    approverID: null,
    approvedAt: null,
    publishedAt: null,
    reviewFeedback: '',
    approvalFeedback: '',
    dateCreated: '2026-10-02T08:00:00Z',
};

describe('LessonPlanComponent', () => {
    let review: ReturnType<typeof vi.fn>;

    beforeEach(async () => {
        review = vi.fn(() => of(current));
        await TestBed.configureTestingModule({
            imports: [LessonPlanComponent],
            providers: [
                {
                    provide: LessonsService,
                    useValue: {
                        plan: () => of(plan),
                        versions: () => of([current]),
                        review,
                    },
                },
                {
                    provide: SchoolAccessService,
                    useValue: {
                        ownMemberships: () =>
                            of([
                                {
                                    id: 'm-1',
                                    schoolID: 'school-1',
                                    departmentID: 'sciences',
                                    userID: 'peter',
                                    capability: 'review_lessons',
                                    active: true,
                                    dateUpdated: '2026-10-01T08:00:00Z',
                                },
                            ]),
                    },
                },
                { provide: UserService, useValue: { user$: of({ id: 'peter', name: 'Peter', email: 'p@x' }) } },
                { provide: MatSnackBar, useValue: { open: vi.fn() } },
                {
                    provide: ActivatedRoute,
                    useValue: { snapshot: { paramMap: convertToParamMap({ planId: 'plan-1' }) } },
                },
            ],
        })
            .overrideComponent(LessonPlanComponent, { set: { template: '' } })
            .compileComponents();
    });

    it('offers HOD review on the current version to a different reviewer', () => {
        const component = TestBed.createComponent(LessonPlanComponent).componentInstance;

        expect(component.current()?.version).toBe(2);
        expect(component.has('review')).toBe(true);
        expect(component.has('approve')).toBe(false);
    });

    it('requires feedback before requesting changes', () => {
        const component = TestBed.createComponent(LessonPlanComponent).componentInstance;
        component.decision.set('request_changes');

        component.decide();

        expect(component.feedbackRequired()).toBe(true);
        expect(review).not.toHaveBeenCalled();
    });

    it('shows the conflict banner when the version is no longer current', () => {
        review.mockReturnValue(throwError(() => new HttpErrorResponse({ status: 409 })));
        const component = TestBed.createComponent(LessonPlanComponent).componentInstance;

        component.decide();

        expect(review).toHaveBeenCalledWith('plan-1', 2, 'approve', '');
        expect(component.conflict()).toBe(true);
        expect(component.errorMessage()).toBeNull();
    });
});

describe('versionRelation', () => {
    it('distinguishes live, latest and superseded versions', () => {
        const draftOverLive = { ...plan, currentVersion: 3, publishedVersion: 2 };

        expect(versionRelation(draftOverLive, 2)).toBe('live');
        expect(versionRelation(draftOverLive, 3)).toBe('latest');
        expect(versionRelation(draftOverLive, 1)).toBe('superseded');
        expect(versionRelation({ ...plan, currentVersion: 2, publishedVersion: 2 }, 2)).toBe(
            'live-and-latest'
        );
    });
});
