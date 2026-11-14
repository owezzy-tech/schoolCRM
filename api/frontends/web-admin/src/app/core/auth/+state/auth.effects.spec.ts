import { provideHttpClient } from '@angular/common/http';
import {
    HttpTestingController,
    provideHttpClientTesting,
} from '@angular/common/http/testing';
import { TestBed } from '@angular/core/testing';
import { provideMockActions } from '@ngrx/effects/testing';
import { Action } from '@ngrx/store';
import { TokenStorageService } from 'app/core/auth/token-storage.service';
import { UserService } from 'app/core/user/user.service';
import { firstValueFrom, Subject } from 'rxjs';
import { afterEach, beforeEach, describe, expect, it } from 'vitest';

import { AuthApiActions, AuthPageActions } from './auth.actions';
import { AuthEffects } from './auth.effects';

describe('AuthEffects sign-in failure messages', () => {
    let actions$: Subject<Action>;
    let http: HttpTestingController;
    let effects: AuthEffects;

    beforeEach(() => {
        actions$ = new Subject<Action>();
        TestBed.configureTestingModule({
            providers: [
                AuthEffects,
                provideMockActions(() => actions$),
                provideHttpClient(),
                provideHttpClientTesting(),
                { provide: TokenStorageService, useValue: { set: () => {} } },
                { provide: UserService, useValue: {} },
            ],
        });
        http = TestBed.inject(HttpTestingController);
        effects = TestBed.inject(AuthEffects);
    });

    afterEach(() => http.verify());

    async function failWith(
        status: number,
        body: unknown,
        statusText = 'Error'
    ): Promise<Action> {
        const result = firstValueFrom(effects.signIn$);
        actions$.next(
            AuthPageActions.signIn({
                email: 'person@schoolcrm.invalid',
                password: 'not-the-password',
            })
        );
        const req = http.expectOne('/v1/auth/login');
        if (status === 0) {
            req.error(new ProgressEvent('error'), {
                status: 0,
                statusText: '',
            });
        } else {
            req.flush(body, { status, statusText });
        }
        return result;
    }

    it('shows the client-safe JSON:API detail for a rejected sign-in', async () => {
        const action = await failWith(401, {
            jsonapi: { version: '1.1' },
            errors: [{ status: '401', detail: 'Invalid email or password.' }],
        });
        expect(action).toEqual(
            AuthApiActions.signInFailure({
                error: 'Invalid email or password.',
            })
        );
    });

    it('explains a network failure instead of showing an empty error', async () => {
        const action = await failWith(0, null);
        expect(action).toEqual(
            AuthApiActions.signInFailure({
                error: "We couldn't reach SchoolCRM. Check your connection and try again.",
            })
        );
    });

    it('does not surface server or proxy error bodies', async () => {
        const action = await failWith(
            502,
            {
                errors: [
                    {
                        detail: 'upstream connect error: dial tcp 10.43.0.5:6000',
                    },
                ],
            },
            'Bad Gateway'
        );
        expect(action).toEqual(
            AuthApiActions.signInFailure({
                error: 'Sign-in is temporarily unavailable. Please try again in a moment.',
            })
        );
    });

    it('asks the user to wait when sign-in is rate limited', async () => {
        const action = await failWith(429, null, 'Too Many Requests');
        expect(action).toEqual(
            AuthApiActions.signInFailure({
                error: 'Too many sign-in attempts. Wait a minute and try again.',
            })
        );
    });
});
