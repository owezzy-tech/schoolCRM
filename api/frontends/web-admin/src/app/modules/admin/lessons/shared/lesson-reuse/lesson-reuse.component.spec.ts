import { HttpErrorResponse } from '@angular/common/http';
import { provideZonelessChangeDetection } from '@angular/core';
import { TestBed } from '@angular/core/testing';
import { provideRouter } from '@angular/router';
import { LessonsService } from 'app/core/lessons/lessons.service';
import { of, throwError } from 'rxjs';
import { describe, expect, it, vi } from 'vitest';
import { LessonReuseComponent } from './lesson-reuse.component';

describe('exact-version lesson reuse', () => {
    it('retries an ambiguous reply with the same command and links to its resulting draft', async () => {
        const reuse = vi
            .fn()
            .mockReturnValueOnce(
                throwError(() => new HttpErrorResponse({ status: 503 }))
            )
            .mockReturnValueOnce(of({ id: 'new-plan' }));
        TestBed.configureTestingModule({
            imports: [LessonReuseComponent],
            providers: [
                provideRouter([]),
                provideZonelessChangeDetection(),
                { provide: LessonsService, useValue: { reuse } },
            ],
        });
        const fixture = TestBed.createComponent(LessonReuseComponent);
        fixture.componentRef.setInput('planID', 'source-plan');
        fixture.componentRef.setInput('version', 2);
        fixture.componentRef.setInput('allowed', true);
        await fixture.whenStable();
        fixture.nativeElement.querySelector('button').click();
        await fixture.whenStable();
        expect(fixture.nativeElement.querySelector('input').readOnly).toBe(
            true
        );
        fixture.nativeElement.querySelector('button').click();
        await fixture.whenStable();
        expect(reuse.mock.calls[0]).toEqual(reuse.mock.calls[1]);
        expect(reuse.mock.calls[0].slice(0, 2)).toEqual(['source-plan', 2]);
        expect(
            fixture.nativeElement.querySelector('a').getAttribute('href')
        ).toBe('/lessons/new-plan');
        expect(fixture.nativeElement.textContent).toContain(
            'new HOD review and dean approval'
        );
    });
    it('hides reuse and refuses a direct call without teaching authority', async () => {
        const reuse = vi.fn();
        TestBed.configureTestingModule({
            imports: [LessonReuseComponent],
            providers: [
                provideRouter([]),
                provideZonelessChangeDetection(),
                { provide: LessonsService, useValue: { reuse } },
            ],
        });
        const fixture = TestBed.createComponent(LessonReuseComponent);
        fixture.componentRef.setInput('planID', 'source-plan');
        fixture.componentRef.setInput('version', 1);
        await fixture.whenStable();
        fixture.componentInstance.reuse();
        expect(fixture.nativeElement.querySelector('button')).toBeNull();
        expect(reuse).not.toHaveBeenCalled();
    });
});
