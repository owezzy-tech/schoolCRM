import { HttpClient } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { Actions, ofType } from '@ngrx/effects';
import { Store } from '@ngrx/store';
import { AuthUtils } from 'app/core/auth/auth.utils';
import { TokenStorageService } from 'app/core/auth/token-storage.service';
import { Observable, of } from 'rxjs';
import { map, switchMap, take } from 'rxjs/operators';
import { AuthApiActions, AuthPageActions, selectAuthenticated } from './+state';

@Injectable({ providedIn: 'root' })
export class AuthService {
    private readonly store = inject(Store);
    private readonly actions$ = inject(Actions);
    private readonly httpClient = inject(HttpClient);
    private readonly tokenStorage = inject(TokenStorageService);

    set accessToken(token: string) {
        this.tokenStorage.set(token);
    }

    get accessToken(): string {
        return this.tokenStorage.get() ?? '';
    }

    forgotPassword(email: string): Observable<unknown> {
        return this.httpClient.post('api/auth/forgot-password', email);
    }

    resetPassword(password: string): Observable<unknown> {
        return this.httpClient.post('api/auth/reset-password', password);
    }

    signIn(credentials: { email: string; password: string }): Observable<void> {
        this.store.dispatch(AuthPageActions.signIn(credentials));
        return of(undefined);
    }

    signInUsingToken(): Observable<boolean> {
        this.store.dispatch(AuthPageActions.signInUsingToken());
        return this.actions$.pipe(
            ofType(
                AuthApiActions.signInUsingTokenSuccess,
                AuthApiActions.signInUsingTokenFailure
            ),
            take(1),
            map((action) =>
                action.type === AuthApiActions.signInUsingTokenSuccess.type
            )
        );
    }

    signOut(): Observable<boolean> {
        this.store.dispatch(AuthPageActions.signOut());
        return of(true);
    }

    signUp(user: {
        name: string;
        email: string;
        password: string;
        company: string;
    }): Observable<unknown> {
        return this.httpClient.post('api/auth/sign-up', user);
    }

    unlockSession(credentials: {
        email: string;
        password: string;
    }): Observable<unknown> {
        return this.httpClient.post('api/auth/unlock-session', credentials);
    }

    check(): Observable<boolean> {
        return this.store.select(selectAuthenticated).pipe(
            take(1),
            switchMap((authenticated) => {
                if (authenticated) {
                    return of(true);
                }
                const token = this.accessToken;
                if (!token || AuthUtils.isTokenExpired(token)) {
                    return of(false);
                }
                return this.signInUsingToken();
            })
        );
    }
}
