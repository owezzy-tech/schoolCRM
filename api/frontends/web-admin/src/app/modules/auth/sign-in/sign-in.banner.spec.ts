import { provideHttpClient } from '@angular/common/http';
import { provideHttpClientTesting } from '@angular/common/http/testing';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { provideNoopAnimations } from '@angular/platform-browser/animations';
import { provideRouter, Router } from '@angular/router';
import { Actions } from '@ngrx/effects';
import { Action, Store } from '@ngrx/store';
import { AuthApiActions } from 'app/core/auth/+state';
import { Subject } from 'rxjs';
import { beforeEach, describe, expect, it, vi } from 'vitest';

import { AuthSignInComponent } from './sign-in.component';

describe('AuthSignInComponent result banner', () => {
    let fixture: ComponentFixture<AuthSignInComponent>;
    let actions$: Subject<Action>;
    let navigateByUrl: ReturnType<typeof vi.spyOn>;

    beforeEach(async () => {
        actions$ = new Subject<Action>();
        await TestBed.configureTestingModule({
            imports: [AuthSignInComponent],
            providers: [
                provideRouter([]),
                provideNoopAnimations(),
                provideHttpClient(),
                provideHttpClientTesting(),
                { provide: Store, useValue: { dispatch: vi.fn() } },
                { provide: Actions, useValue: actions$ },
            ],
        }).compileComponents();
        navigateByUrl = vi
            .spyOn(TestBed.inject(Router), 'navigateByUrl')
            .mockResolvedValue(true);
        fixture = TestBed.createComponent(AuthSignInComponent);
        fixture.detectChanges();
    });

    function banner(): HTMLElement | null {
        return fixture.nativeElement.querySelector(
            '[data-testid="sign-in-banner"]'
        );
    }

    it('shows no banner before the user submits', () => {
        expect(banner()).toBeNull();
    });

    it('shows a red alert with the server message and keeps the email', () => {
        const form = fixture.componentInstance.signInForm;
        form.patchValue({
            email: 'person@schoolcrm.invalid',
            password: 'wrong-password',
        });

        actions$.next(
            AuthApiActions.signInFailure({
                error: 'Invalid email or password.',
            })
        );
        fixture.detectChanges();

        const el = banner();
        expect(el.getAttribute('role')).toBe('alert');
        expect(el.getAttribute('aria-live')).toBe('assertive');
        expect(el.classList).toContain('fuse-alert-type-error');
        expect(el.classList).toContain('fuse-alert-appearance-soft');
        expect(el.textContent).toContain("Couldn't sign you in");
        expect(el.textContent).toContain('Invalid email or password.');
        expect(form.get('email').value).toBe('person@schoolcrm.invalid');
        expect(form.get('password').value).toBe('');
        expect(form.enabled).toBe(true);
    });

    it('shows a green status banner on success and redirects', () => {
        actions$.next(
            AuthApiActions.signInSuccess({
                accessToken: 'token',
                user: {
                    id: '5cf37266-3473-4006-984f-9325122678b7',
                    name: 'Admin',
                    email: 'admin@schoolcrm.invalid',
                    roles: ['SCHOOL_ADMIN'],
                },
            })
        );
        fixture.detectChanges();

        const el = banner();
        expect(el.getAttribute('role')).toBe('status');
        expect(el.getAttribute('aria-live')).toBe('polite');
        expect(el.classList).toContain('fuse-alert-type-success');
        expect(el.textContent).toContain('Signed in');
        expect(navigateByUrl).toHaveBeenCalledWith('/signed-in-redirect');
    });
});
