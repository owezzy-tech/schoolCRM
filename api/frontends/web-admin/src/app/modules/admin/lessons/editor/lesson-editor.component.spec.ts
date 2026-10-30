import { provideZonelessChangeDetection } from '@angular/core';
import { TestBed } from '@angular/core/testing';
import { provideRouter } from '@angular/router';
import { RouterTestingHarness } from '@angular/router/testing';
import { LessonsService } from 'app/core/lessons/lessons.service';
import { LessonPlan, LessonVersion } from 'app/core/lessons/lessons.types';
import { of } from 'rxjs';
import { describe, expect, it } from 'vitest';
import { LessonEditorComponent } from './lesson-editor.component';

describe('Editing generated snapshots', () => {
    it('preserves citations and generation provenance when a teacher edits the content', async () => {
        const content = {
            schemaVersion: 1,
            framework: 'kenya-cbc',
            stage: 'grade-1',
            subject: 'english',
            curriculumRevision: '2024',
            durationMinutes: 30,
            objectives: ['Original objective'],
            prerequisites: ['Recognise objects'],
            materials: ['Pictures'],
            activities: [
                { minutes: 30, title: 'Practice', detail: 'Name objects' },
            ],
            assessment: 'Name three objects',
            differentiation: 'Use visual prompts',
            citations: [
                {
                    sourceID: 'source-1',
                    page: 17,
                    revision: '2024',
                    title: 'English',
                    sourceURL: 'https://kicd.ac.ke/english.pdf',
                },
            ],
            generation: {
                provider: 'deepseek',
                modelID: 'deepseek-flash',
                completionID: 'completion-1',
            },
            reusedFrom: { planID: 'original', version: 1 },
        };
        const plan: LessonPlan = {
            id: 'plan-1',
            schoolID: 'school-1',
            departmentID: 'dept-1',
            authorID: 'teacher',
            title: 'Objects',
            status: 'draft',
            currentVersion: 2,
            publishedVersion: 1,
            dateCreated: '',
            dateUpdated: '',
        };
        const version: LessonVersion = {
            id: 'plan-1:2',
            planID: plan.id,
            version: 2,
            title: plan.title,
            content,
            changeSummary: '',
            authorID: plan.authorID,
            status: 'draft',
            reviewerID: null,
            reviewedAt: null,
            approverID: null,
            approvedAt: null,
            publishedAt: null,
            reviewFeedback: '',
            approvalFeedback: '',
            dateCreated: '',
        };
        await TestBed.configureTestingModule({
            providers: [
                provideZonelessChangeDetection(),
                provideRouter([
                    {
                        path: 'lessons/:planId/edit',
                        component: LessonEditorComponent,
                    },
                ]),
                {
                    provide: LessonsService,
                    useValue: {
                        plan: () => of(plan),
                        versions: () => of([version]),
                    },
                },
            ],
        }).compileComponents();
        const harness = await RouterTestingHarness.create();
        const component = await harness.navigateByUrl(
            '/lessons/plan-1/edit',
            LessonEditorComponent
        );
        await harness.fixture.whenStable();
        component.form.controls.objectives.at(0).setValue('Edited objective');
        component.form.controls.prerequisites.setValue(
            'Recognise objects\nListen carefully'
        );
        component.activities.at(0).controls.minutes.setValue(40);
        const draft = component.toDraft();
        expect(draft.content.objectives).toEqual(['Edited objective']);
        expect(draft.content.prerequisites).toEqual([
            'Recognise objects',
            'Listen carefully',
        ]);
        expect(draft.content).toMatchObject({
            durationMinutes: 40,
            citations: content.citations,
            generation: content.generation,
            reusedFrom: content.reusedFrom,
            framework: 'kenya-cbc',
        });
        expect(content.durationMinutes).toBe(30);
        expect(content.objectives).toEqual(['Original objective']);
    });
});
