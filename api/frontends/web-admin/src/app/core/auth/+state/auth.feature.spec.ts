import { describe, expect, it } from 'vitest';

import { AuthApiActions, AuthPageActions } from './auth.actions';
import { authReducer } from './auth.feature';

describe('authReducer', () => {
    it('stores the authenticated user and token on sign-in success', () => {
        const loadingState = authReducer(
            undefined,
            AuthPageActions.signIn({
                email: 'admin@example.com',
                password: 'gophers',
            })
        );

        const state = authReducer(
            loadingState,
            AuthApiActions.signInSuccess({
                accessToken: 'access-token',
                user: {
                    id: '5cf37266-3473-4006-984f-9325122678b7',
                    name: 'Admin Gopher',
                    email: 'admin@example.com',
                    roles: ['SCHOOL_ADMIN'],
                },
            })
        );

        expect(state.authenticated).toBe(true);
        expect(state.accessToken).toBe('access-token');
        expect(state.user?.roles).toEqual(['SCHOOL_ADMIN']);
        expect(state.loading).toBe(false);
        expect(state.error).toBeNull();
    });

    it('clears auth state when token validation fails', () => {
        const authenticatedState = authReducer(
            undefined,
            AuthApiActions.signInSuccess({
                accessToken: 'access-token',
                user: {
                    id: '5cf37266-3473-4006-984f-9325122678b7',
                    name: 'Admin Gopher',
                    email: 'admin@example.com',
                    roles: ['SCHOOL_ADMIN'],
                },
            })
        );

        const state = authReducer(
            authenticatedState,
            AuthApiActions.signInUsingTokenFailure()
        );

        expect(state.authenticated).toBe(false);
        expect(state.accessToken).toBeNull();
        expect(state.user).toBeNull();
        expect(state.loading).toBe(false);
    });
});
