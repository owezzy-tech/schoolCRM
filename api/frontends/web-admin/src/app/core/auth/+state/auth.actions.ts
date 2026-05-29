import { createActionGroup, emptyProps, props } from '@ngrx/store';
import { User } from 'app/core/user/user.types';

export interface AuthSuccessPayload {
    accessToken: string;
    user: User;
}

export const AuthPageActions = createActionGroup({
    source: 'Auth/Page',
    events: {
        'Sign In': props<{ email: string; password: string }>(),
        'Sign Out': emptyProps(),
        'Sign In Using Token': emptyProps(),
    },
});

export const AuthApiActions = createActionGroup({
    source: 'Auth/API',
    events: {
        'Sign In Success': props<AuthSuccessPayload>(),
        'Sign In Failure': props<{ error: string }>(),

        'Sign In Using Token Success': props<AuthSuccessPayload>(),
        'Sign In Using Token Failure': props<{ error: string }>(),

        'Sign Out Success': emptyProps(),
    },
});
