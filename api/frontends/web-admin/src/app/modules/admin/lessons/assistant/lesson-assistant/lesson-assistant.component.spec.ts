import { provideZonelessChangeDetection } from '@angular/core';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { provideRouter, Router } from '@angular/router';
import { AuthService } from 'app/core/auth/auth.service';
import { LessonAssistantService } from 'app/core/lessons/lesson-assistant.service';
import {
    AssistantEvent,
    AssistantStreamError,
    GenerationRequest,
    LessonThread,
} from 'app/core/lessons/lesson-assistant.types';
import { LessonsService } from 'app/core/lessons/lessons.service';
import { SchoolAccessService } from 'app/core/school-access/school-access.service';
import { Observable, of, Subject } from 'rxjs';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { LessonAssistantComponent } from './lesson-assistant.component';

const request: GenerationRequest = {
    request_id: 'saved-request',
    school_id: 'school',
    department_id: 'dept',
    framework: 'kenya-cbc',
    stage: 'grade-1',
    subject: 'english',
    revision: '2024',
    topic: 'Recover my vocabulary lesson',
    duration_minutes: 30,
};
describe('lesson assistant screen', () => {
    let capability = 'teach';
    let events: Subject<AssistantEvent>;
    let generate: ReturnType<typeof vi.fn>;
    let threadRead: Observable<LessonThread>;
    let signOut: ReturnType<typeof vi.fn>;
    let fixture: ComponentFixture<LessonAssistantComponent>;
    beforeEach(async () => {
        capability = 'teach';
        events = new Subject();
        generate = vi.fn(() => events);
        signOut = vi.fn();
        threadRead = of({
            id: request.request_id,
            request,
            status: 'pending',
            result: null,
        });
        await TestBed.configureTestingModule({
            imports: [LessonAssistantComponent],
            providers: [
                provideRouter([]),
                provideZonelessChangeDetection(),
                { provide: AuthService, useValue: { signOut } },
                {
                    provide: SchoolAccessService,
                    useValue: {
                        ownMemberships: () =>
                            of([
                                {
                                    id: 'm',
                                    schoolID: 'school',
                                    departmentID: 'dept',
                                    capability,
                                    active: true,
                                },
                            ]),
                        schools: () =>
                            of([{ id: 'school', name: 'Test school' }]),
                        departments: () =>
                            of([{ id: 'dept', name: 'English' }]),
                    },
                },
                {
                    provide: LessonAssistantService,
                    useValue: {
                        threads: () =>
                            of([
                                {
                                    id: request.request_id,
                                    request,
                                    status: 'pending',
                                    result: null,
                                },
                            ]),
                        thread: () => threadRead,
                        generate,
                    },
                },
                {
                    provide: LessonsService,
                    useValue: {
                        versions: () =>
                            of([
                                {
                                    version: 1,
                                    title: 'Saved vocabulary',
                                    content: {
                                        objectives: ['Name items at school'],
                                        activities: [
                                            {
                                                title: 'Practice',
                                                minutes: 30,
                                                detail: 'Read together',
                                            },
                                        ],
                                        assessment: 'Name five items',
                                        citations: [
                                            {
                                                sourceID: 'source',
                                                title: 'Official curriculum',
                                                sourceURL:
                                                    'https://kicd.ac.ke/source.pdf',
                                                page: 17,
                                                revision: '2024',
                                                sourceSHA256: 'a'.repeat(64),
                                            },
                                        ],
                                    },
                                },
                            ]),
                    },
                },
            ],
        }).compileComponents();
        vi.spyOn(TestBed.inject(Router), 'navigate').mockResolvedValue(true);
    });
    async function render() {
        fixture = TestBed.createComponent(LessonAssistantComponent);
        await fixture.whenStable();
        return fixture.componentInstance;
    }
    function button(text: string): HTMLButtonElement {
        return [...fixture.nativeElement.querySelectorAll('button')].find(
            (b: HTMLButtonElement) => b.textContent?.includes(text)
        );
    }
    it('recovers a persisted unfinished request after refresh and retries its exact input', async () => {
        const component = await render();
        button(request.topic).click();
        await fixture.whenStable();
        button('Reconnect this request').click();
        expect(generate).toHaveBeenCalledWith(request);
        events.error(new TypeError('Failed to fetch secret-host:5432'));
        await fixture.whenStable();
        expect(fixture.nativeElement.textContent).toContain(
            'Reconnect to recover this same request'
        );
        expect(fixture.nativeElement.textContent).not.toContain('secret-host');
        events = new Subject();
        generate.mockReturnValue(events);
        button('Reconnect this request').click();
        expect(generate.mock.calls[1][0]).toEqual(request);
        expect(component.request()?.request_id).toBe(request.request_id);
    });
    it('renders exact saved content and evidence while requiring human decisions', async () => {
        const component = await render();
        component.openRequest(request.request_id);
        component.reconnect();
        events.next({
            type: 'citation',
            citations: [
                {
                    sourceID: 'unused',
                    title: 'Retrieved but unused',
                    sourceURL: 'https://kicd.ac.ke/unused.pdf',
                    page: 18,
                    revision: '2024',
                    passage: 'Unused retrieved passage',
                },
            ],
        });
        await fixture.whenStable();
        expect(fixture.nativeElement.textContent).toContain(
            'not all may be cited'
        );
        expect(
            fixture.nativeElement
                .querySelector('[aria-live="polite"]')
                .hasAttribute('aria-busy')
        ).toBe(false);
        events.next({
            type: 'structured-result',
            result: { planID: 'plan', title: 'Saved vocabulary', version: 1 },
        });
        events.next({
            type: 'approval-request',
            planID: 'plan',
            version: 1,
            action: 'submit-for-review',
        });
        events.next({ type: 'completed' });
        events.complete();
        await fixture.whenStable();
        expect(fixture.nativeElement.textContent).toContain(
            'Name items at school'
        );
        expect(fixture.nativeElement.textContent).not.toContain(
            'Unused retrieved passage'
        );
        expect(fixture.nativeElement.textContent).toContain('PDF page 17');
        expect(fixture.nativeElement.textContent).toContain(
            'assistant cannot approve or publish'
        );
        expect(
            fixture.nativeElement.querySelector(
                'a[href="/lessons/plan/versions/1"]'
            )
        ).not.toBeNull();
        expect(button('Reconnect this request')).toBeUndefined();
    });
    it('does not offer generation or private history to a reviewer without teaching access', async () => {
        capability = 'review_lessons';
        await render();
        expect(fixture.nativeElement.textContent).toContain(
            'Teaching access required'
        );
        expect(fixture.nativeElement.querySelector('form')).toBeNull();
        expect(generate).not.toHaveBeenCalled();
    });
    it('shows a terminal failure and leaves the same request recoverable', async () => {
        const component = await render();
        component.openRequest(request.request_id);
        component.reconnect();
        events.next({
            type: 'terminal-error',
            status: 422,
            detail: 'No approved evidence supports this topic',
        });
        events.complete();
        await fixture.whenStable();
        expect(fixture.nativeElement.textContent).toContain(
            'No approved evidence supports this topic'
        );
        expect(button('Reconnect this request').disabled).toBe(false);
        expect(component.completed()).toBe(false);
    });
    it('shows only client-safe stream protocol failures', async () => {
        const component = await render();
        component.openRequest(request.request_id);
        component.reconnect();
        events.error(
            new AssistantStreamError('The assistant returned an invalid event.')
        );
        await fixture.whenStable();
        expect(fixture.nativeElement.textContent).toContain(
            'The assistant returned an invalid event.'
        );
    });
    it('sends a rejected bearer to sign-in once instead of looping reconnects', async () => {
        const component = await render();
        component.openRequest(request.request_id);
        component.reconnect();
        const navigate = vi.mocked(TestBed.inject(Router).navigate);
        navigate.mockClear();
        events.next({
            type: 'terminal-error',
            status: 401,
            detail: 'Authentication failed',
        });
        expect(signOut).toHaveBeenCalledOnce();
        expect(navigate).toHaveBeenCalledWith(['/sign-in'], {
            queryParams: { redirectURL: expect.any(String) },
        });
        expect(component.busy()).toBe(false);
        expect(events.observed).toBe(false);
    });
    it('ignores a late history read after a new request starts', async () => {
        const late = new Subject<LessonThread>();
        threadRead = late;
        const component = await render();
        component.openRequest('older-request');
        component.topic.set('A brand new topic');
        component.start();
        expect(late.observed).toBe(false);
        expect(component.request()?.topic).toBe('A brand new topic');
        expect(generate).toHaveBeenCalledOnce();
    });
});
