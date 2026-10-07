import { HttpErrorResponse } from '@angular/common/http';
import { provideZonelessChangeDetection } from '@angular/core';
import { TestBed } from '@angular/core/testing';
import { provideRouter } from '@angular/router';
import { RouterTestingHarness } from '@angular/router/testing';
import { LessonsService } from 'app/core/lessons/lessons.service';
import { LessonPlan, LessonVersion } from 'app/core/lessons/lessons.types';
import { of, throwError } from 'rxjs';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';

import { LessonExportComponent } from './lesson-export.component';

const plan: LessonPlan = {
    id: 'plan-1',
    schoolID: 'school-1',
    departmentID: 'languages',
    authorID: 'teacher',
    title: 'Latest title',
    status: 'draft',
    currentVersion: 2,
    publishedVersion: 1,
    dateCreated: '2026-10-01T08:00:00Z',
    dateUpdated: '2026-10-02T08:00:00Z',
};
const published: LessonVersion = {
    id: 'version-1',
    planID: plan.id,
    version: 1,
    title: 'Original title',
    content: {
        objectives: ['Listen carefully'],
        activities: [{ minutes: 0, title: 'Observe' }],
        citations: [
            {
                source_url: 'https://kicd.ac.ke/original.pdf',
                page: 17,
                revision: '2024',
            },
        ],
        differentiation: 'Use visual prompts',
        verified: false,
        extra: '<script>unsafe()</script>',
    },
    changeSummary: 'Original snapshot',
    authorID: 'teacher',
    status: 'published',
    reviewerID: 'hod',
    reviewedAt: '2026-10-01T09:00:00Z',
    approverID: 'dean',
    approvedAt: '2026-10-01T10:00:00Z',
    publishedAt: '2026-10-01T11:00:00Z',
    reviewFeedback: 'Reviewed',
    approvalFeedback: 'Approved',
    dateCreated: '2026-10-01T08:00:00Z',
};
const draft: LessonVersion = {
    ...published,
    id: 'version-2',
    version: 2,
    title: 'Latest draft',
    status: 'draft',
    publishedAt: null,
    reviewedAt: null,
    approvedAt: null,
};

describe('Lesson PDF export', () => {
    let planRequest: ReturnType<typeof vi.fn>;
    let harness: RouterTestingHarness;

    beforeEach(async () => {
        planRequest = vi.fn(() => of(plan));
        await TestBed.configureTestingModule({
            providers: [
                provideZonelessChangeDetection(),
                provideRouter([
                    {
                        path: 'lesson-export/:planId/:version',
                        component: LessonExportComponent,
                    },
                ]),
                {
                    provide: LessonsService,
                    useValue: {
                        plan: planRequest,
                        versions: () => of([draft, published]),
                    },
                },
            ],
        }).compileComponents();
        harness = await RouterTestingHarness.create();
    });
    afterEach(() => vi.restoreAllMocks());

    it('exports the selected historical snapshot with every stored field and escaped text', async () => {
        await harness.navigateByUrl(
            '/lesson-export/plan-1/1',
            LessonExportComponent
        );
        await harness.fixture.whenStable();
        const article = harness.routeNativeElement?.querySelector('article');
        expect(article?.textContent).toContain('Original title');
        expect(article?.textContent).not.toContain('Latest title');
        expect(article?.textContent).toContain('Live publication');
        expect(article?.textContent).toContain(
            'https://kicd.ac.ke/original.pdf'
        );
        expect(article?.textContent).toContain('17');
        expect(article?.textContent).toContain('Use visual prompts');
        expect(article?.textContent).toContain('false');
        expect(article?.textContent).toContain('<script>unsafe()</script>');
        expect(article?.querySelector('script')).toBeNull();
    });

    it('reloads the selected version when route parameters change and labels drafts', async () => {
        await harness.navigateByUrl(
            '/lesson-export/plan-1/1',
            LessonExportComponent
        );
        const component = await harness.navigateByUrl(
            '/lesson-export/plan-1/2',
            LessonExportComponent
        );
        await harness.fixture.whenStable();
        expect(component.version()?.version).toBe(2);
        expect(harness.routeNativeElement?.textContent).toContain(
            'This version has not been published.'
        );
        expect(
            harness.routeNativeElement?.querySelector('article')?.textContent
        ).not.toContain('Live publication');
    });

    it('clears the printable snapshot and never prints after authority is revoked', async () => {
        const component = await harness.navigateByUrl(
            '/lesson-export/plan-1/1',
            LessonExportComponent
        );
        const print = vi.spyOn(window, 'print').mockImplementation(() => {});
        planRequest.mockReturnValue(
            throwError(
                () =>
                    new HttpErrorResponse({
                        status: 403,
                        error: { errors: [{ detail: 'Permission revoked' }] },
                    })
            )
        );
        component.print();
        await harness.fixture.whenStable();
        expect(planRequest).toHaveBeenCalledTimes(2);
        expect(component.version()).toBeNull();
        expect(harness.routeNativeElement?.querySelector('article')).toBeNull();
        expect(harness.routeNativeElement?.textContent).toContain(
            'Permission revoked'
        );
        expect(print).not.toHaveBeenCalled();
    });

    it('prints only after fresh authorised reads have rendered the chosen version', async () => {
        const component = await harness.navigateByUrl(
            '/lesson-export/plan-1/1',
            LessonExportComponent
        );
        const print = vi.spyOn(window, 'print').mockImplementation(() => {
            expect(
                harness.routeNativeElement?.querySelector('h1')?.textContent
            ).toBe('Original title');
        });
        component.print();
        await harness.fixture.whenStable();
        expect(planRequest).toHaveBeenCalledTimes(2);
        expect(print).toHaveBeenCalledOnce();
    });

    it('cannot print a version omitted by the authorised API', async () => {
        const component = await harness.navigateByUrl(
            '/lesson-export/plan-1/99',
            LessonExportComponent
        );
        await harness.fixture.whenStable();
        expect(component.version()).toBeNull();
        expect(harness.routeNativeElement?.querySelector('article')).toBeNull();
        expect(
            harness.routeNativeElement?.querySelector('button')?.disabled
        ).toBe(true);
    });
});
