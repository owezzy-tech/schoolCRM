import { describe, expect, it } from 'vitest';

import {
    lessonActions,
    lessonCitations,
    LessonPlan,
    LessonVersion,
    toLessonContent,
} from './lessons.types';

const plan: LessonPlan = {
    id: 'plan-1',
    schoolID: 'school-1',
    departmentID: 'sciences',
    authorID: 'grace',
    title: 'Forces',
    status: 'draft',
    currentVersion: 1,
    publishedVersion: null,
    dateCreated: '2026-10-01T08:00:00Z',
    dateUpdated: '2026-10-01T08:00:00Z',
};

function version(overrides: Partial<LessonVersion>): LessonVersion {
    return {
        id: 'plan-1:1',
        planID: 'plan-1',
        version: 1,
        title: 'Forces',
        content: {},
        changeSummary: '',
        authorID: 'grace',
        status: 'draft',
        reviewerID: null,
        reviewedAt: null,
        approverID: null,
        approvedAt: null,
        publishedAt: null,
        reviewFeedback: '',
        approvalFeedback: '',
        dateCreated: '2026-10-01T08:00:00Z',
        ...overrides,
    };
}

describe('lessonActions', () => {
    it('lets the teaching author edit and submit a draft, but not publish it', () => {
        expect(
            lessonActions(plan, version({}), 'grace', ['teach']).actions
        ).toEqual(['edit', 'submit']);
    });

    it('offers publication only for an approved version', () => {
        expect(
            lessonActions(plan, version({ status: 'approved' }), 'grace', [
                'teach',
            ]).actions
        ).toEqual(['edit', 'publish']);
    });

    it('offers no workflow actions to a teacher who is not the author', () => {
        expect(
            lessonActions(plan, version({}), 'peter', ['teach']).actions
        ).toEqual([]);
    });

    it('offers HOD review to a different reviewer', () => {
        const result = lessonActions(
            plan,
            version({ status: 'hod_review' }),
            'peter',
            ['review_lessons']
        );
        expect(result).toEqual({
            actions: ['review'],
            separationOfDuties: false,
        });
    });

    it('blocks the author from reviewing their own version', () => {
        const result = lessonActions(
            plan,
            version({ status: 'hod_review' }),
            'grace',
            ['teach', 'review_lessons']
        );
        expect(result.actions).not.toContain('review');
        expect(result.separationOfDuties).toBe(true);
    });

    it('blocks the HOD reviewer from also approving as dean', () => {
        const result = lessonActions(
            plan,
            version({ status: 'dean_approval', reviewerID: 'peter' }),
            'peter',
            ['review_lessons', 'approve_lessons']
        );
        expect(result).toEqual({ actions: [], separationOfDuties: true });
    });

    it('offers dean approval to a third person', () => {
        expect(
            lessonActions(
                plan,
                version({ status: 'dean_approval', reviewerID: 'peter' }),
                'mary',
                ['approve_lessons']
            ).actions
        ).toEqual(['approve']);
    });
});

describe('toLessonContent', () => {
    it('retains the approved prerequisite and differentiation fields', () => {
        const parsed = toLessonContent({
            prerequisites: ['Identify objects'],
            differentiation: 'Use pictures',
        });
        expect(parsed.prerequisites).toEqual(['Identify objects']);
        expect(parsed.differentiation).toBe('Use pictures');
    });
    it('does not expose unsafe links as curriculum citations', () => {
        const citation = {
            sourceID: 'source-1',
            title: 'English',
            page: 17,
            revision: '2024',
            sourceURL: 'https://kicd.ac.ke/english.pdf',
        };
        expect(lessonCitations({ citations: [citation] })).toEqual([citation]);
        expect(
            lessonCitations({
                citations: [{ ...citation, sourceURL: 'javascript:alert(1)' }],
            })
        ).toEqual([]);
    });
    it('reads the structured shape', () => {
        expect(
            toLessonContent({
                objectives: ['Define force'],
                materials: ['Springs'],
                activities: [
                    { minutes: 10, title: 'Starter', detail: 'Recall' },
                ],
                assessment: 'Exit ticket',
            })
        ).toEqual({
            objectives: ['Define force'],
            materials: ['Springs'],
            activities: [{ minutes: 10, title: 'Starter', detail: 'Recall' }],
            assessment: 'Exit ticket',
        });
    });

    it('renders unknown shapes as empty sections', () => {
        expect(
            toLessonContent({ objectives: 'not a list', activities: [null] })
        ).toEqual({
            objectives: [],
            materials: [],
            activities: [],
            assessment: '',
        });
        expect(toLessonContent(null).objectives).toEqual([]);
    });
});
