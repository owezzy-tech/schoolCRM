import { HttpClient, HttpHeaders } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { Actions, createEffect, ofType } from '@ngrx/effects';
import { AuthUtils } from 'app/core/auth/auth.utils';
import { TokenStorageService } from 'app/core/auth/token-storage.service';
import { User } from 'app/core/user/user.types';
import { UserService } from 'app/core/user/user.service';
import { of } from 'rxjs';
import { catchError, map, switchMap, tap } from 'rxjs/operators';
import { AuthApiActions, AuthPageActions } from './auth.actions';

interface JsonApiResource<T> {
    attributes: T;
}

interface JsonApiDocument<T> {
    data: JsonApiResource<T>;
}

interface SignInResponse {
    accessToken: string;
    tokenType: 'Bearer';
    expiresAt: string;
    expiresIn: number;
    user: User;
}

interface AuthenticateResponse {
    userID: string;
    claims: {
        roles?: string[];
    };
    user: User;
}

interface JsonApiError {
    detail?: string;
    title?: string;
}

interface AuthErrorResponse {
    error?: {
        message?: string;
        errors?: JsonApiError[];
    };
}

@Injectable()
export class AuthEffects {
    private readonly actions$ = inject(Actions);
    private readonly httpClient = inject(HttpClient);
    private readonly tokenStorage = inject(TokenStorageService);
    private readonly userService = inject(UserService);

    signIn$ = createEffect(() =>
        this.actions$.pipe(
            ofType(AuthPageActions.signIn),
            switchMap(({ email, password }) =>
                this.httpClient
                    .post<JsonApiDocument<SignInResponse>>('/v1/auth/login', {
                        email,
                        password,
                    })
                    .pipe(
                        map(({ data }) => {
                            const { accessToken, user } = data.attributes;
                            this.tokenStorage.set(accessToken);
                            this.userService.user = user;
                            return AuthApiActions.signInSuccess({
                                accessToken,
                                user,
                            });
                        }),
                        catchError((error: AuthErrorResponse) =>
                            of(
                                AuthApiActions.signInFailure({
                                    error: authErrorMessage(error),
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
                const token = this.tokenStorage.get() ?? '';
                if (!token || AuthUtils.isTokenExpired(token)) {
                    return of(
                        AuthApiActions.signInUsingTokenFailure({
                            error: 'Your session has expired. Please sign in again.',
                        })
                    );
                }

                const headers = new HttpHeaders({
                    Authorization: `Bearer ${token}`,
                });

                return this.httpClient
                    .get<JsonApiDocument<AuthenticateResponse>>('/v1/auth/authenticate', {
                        headers,
                    })
                    .pipe(
                        map(({ data }) => {
                            const response = data.attributes;
                            const user = {
                                ...response.user,
                                roles:
                                    response.user.roles ??
                                    response.claims.roles ??
                                    [],
                            };
                            this.tokenStorage.set(token);
                            this.userService.user = user;
                            return AuthApiActions.signInUsingTokenSuccess({
                                accessToken: token,
                                user,
                            });
                        }),
                        catchError(() =>
                            of(
                                AuthApiActions.signInUsingTokenFailure({
                                    error: 'Unable to restore your session. Please sign in again.',
                                })
                            )
                        )
                    );
            })
        )
    );

    signOut$ = createEffect(() =>
        this.actions$.pipe(
            ofType(AuthPageActions.signOut),
            tap(() => this.tokenStorage.clear()),
            map(() => AuthApiActions.signOutSuccess())
        )
    );
}

function authErrorMessage(error: AuthErrorResponse): string {
    return (
        error.error?.errors?.[0]?.detail ??
        error.error?.errors?.[0]?.title ??
        error.error?.message ??
        'Sign in failed.'
    );
}
