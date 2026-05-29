import { isDevMode } from '@angular/core';
import { Action, ActionReducer, MetaReducer } from '@ngrx/store';

import { authFeatureKey, AuthState } from './auth.feature';

export interface AppState {
    [authFeatureKey]: AuthState;
}

export function logger(
    reducer: ActionReducer<AppState>
): ActionReducer<AppState> {
    return (state: AppState | undefined, action: Action): AppState => {
        const result = reducer(state, action);

        console.groupCollapsed(action.type);
        console.log('prev state', state);
        console.log('action', action);
        console.log('next state', result);
        console.groupEnd();

        return result;
    };
}

export const metaReducers: MetaReducer<AppState>[] = isDevMode()
    ? [logger]
    : [];
