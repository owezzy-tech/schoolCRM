import { createFeature, createReducer, on } from '@ngrx/store';
import { User } from 'app/core/user/user.types';
import { AuthApiActions, AuthPageActions } from './auth.actions';

export interface AuthState {
    authenticated: boolean;
    accessToken: string | null;
    user: User | null;
    loading: boolean;
    error: string | null;
}

const getInitialAccessToken = (): string | null => {
    if (typeof localStorage === 'undefined') {
        return null;
    }

    return localStorage.getItem('accessToken');
};

const initialState: AuthState = {
    authenticated: false,
    accessToken: getInitialAccessToken(),
    user: null,
    loading: false,
    error: null,
};

export const authFeature = createFeature({
    name: 'auth',
    reducer: createReducer(
        initialState,

        on(AuthPageActions.signIn, (state) => ({
            ...state,
            loading: true,
            error: null,
        })),

        on(AuthApiActions.signInSuccess, (state, { accessToken, user }) => ({
            ...state,
            authenticated: true,
            accessToken,
            user,
            loading: false,
            error: null,
        })),

        on(AuthApiActions.signInFailure, (state, { error }) => ({
            ...state,
            authenticated: false,
            accessToken: null,
            user: null,
            loading: false,
            error,
        })),

        on(AuthPageActions.signInUsingToken, (state) => ({
            ...state,
            loading: true,
            error: null,
        })),

        on(
            AuthApiActions.signInUsingTokenSuccess,
            (state, { accessToken, user }) => ({
                ...state,
                authenticated: true,
                accessToken,
                user,
                loading: false,
                error: null,
            })
        ),

        on(AuthApiActions.signInUsingTokenFailure, (state) => ({
            ...state,
            authenticated: false,
            accessToken: null,
            user: null,
            loading: false,
        })),

        on(AuthApiActions.signOutSuccess, (_state) => ({
            authenticated: false,
            accessToken: null,
            user: null,
            loading: false,
            error: null,
        }))
    ),
});

export const {
    name: authFeatureKey,
    reducer: authReducer,
    selectAuthState,
    selectAuthenticated,
    selectAccessToken,
    selectUser,
    selectLoading,
    selectError,
} = authFeature;
