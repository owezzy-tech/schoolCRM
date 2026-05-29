import { HttpClient } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { Actions, createEffect, ofType } from '@ngrx/effects';
import { AuthUtils } from 'app/core/auth/auth.utils';
import { User } from 'app/core/user/user.types';
import { UserService } from 'app/core/user/user.service';
import { of } from 'rxjs';
import { catchError, map, switchMap, tap } from 'rxjs/operators';
import { AuthApiActions, AuthPageActions } from './auth.actions';

interface SignInResponse {
    accessToken: string;
    user: User;
}

interface TokenSignInResponse {
    accessToken?: string;
    user: User;
}

interface AuthErrorResponse {
    error?: {
        message?: string;
    };
}

@Injectable()
export class AuthEffects {
    private readonly actions$ = inject(Actions);
    private readonly httpClient = inject(HttpClient);
    private readonly userService = inject(UserService);

    signIn$ = createEffect(() =>
        this.actions$.pipe(
            ofType(AuthPageActions.signIn),
            switchMap(({ email, password }) =>
                this.httpClient
                    .post<SignInResponse>('api/auth/sign-in', {
                        email,
                        password,
                    })
                    .pipe(
                        map(({ accessToken, user }) => {
                            localStorage.setItem('accessToken', accessToken);
                            this.userService.user = user;
                            return AuthApiActions.signInSuccess({
                                accessToken,
                                user,
                            });
                        }),
                        catchError((error: AuthErrorResponse) =>
                            of(
                                AuthApiActions.signInFailure({
                                    error:
                                        error.error?.message ??
                                        'Sign in failed.',
                                })
                            )
                        )
                    )
            )
        )
    );

    signInUsingToken$ = createEffect(() =>
        this.actions$.pipe(
            ofType(AuthPageActions.signInUsingToken),
            switchMap(() => {
                const token = localStorage.getItem('accessToken') ?? '';
                if (!token || AuthUtils.isTokenExpired(token)) {
                    return of(AuthApiActions.signInUsingTokenFailure());
                }
                return this.httpClient
                    .post<TokenSignInResponse>(
                        'api/auth/sign-in-with-token',
                        { accessToken: token }
                    )
                    .pipe(
                        map((response) => {
                            const newToken = response.accessToken ?? token;
                            localStorage.setItem('accessToken', newToken);
                            this.userService.user = response.user;
                            return AuthApiActions.signInUsingTokenSuccess({
                                accessToken: newToken,
                                user: response.user,
                            });
                        }),
                        catchError(() =>
                            of(AuthApiActions.signInUsingTokenFailure())
                        )
                    );
            })
        )
    );

    signOut$ = createEffect(() =>
        this.actions$.pipe(
            ofType(AuthPageActions.signOut),
            tap(() => localStorage.removeItem('accessToken')),
            map(() => AuthApiActions.signOutSuccess())
        )
    );
}
