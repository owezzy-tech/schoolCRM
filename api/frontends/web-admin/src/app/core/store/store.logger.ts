import { isDevMode } from '@angular/core';
import {
    Action,
    ActionReducer,
    ActionReducerMap,
    MetaReducer,
} from '@ngrx/store';

import {
    getRouterSelectors,
    routerReducer,
    RouterReducerState,
} from '@ngrx/router-store';
import { authFeature, authFeatureKey, AuthState } from '../auth/+state';

export interface AppState {
    [authFeatureKey]: AuthState;
    router: RouterReducerState;
}

/**
 * Our state is composed of a map of action reducer functions.
 * These reducer functions are called with each dispatched action
 * and the current or initial state and return a new immutable state.
 */
export const rootReducers: ActionReducerMap<AppState> = {
    [authFeature.name]: authFeature.reducer,
    router: routerReducer,
};

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

export const { selectRouteData } = getRouterSelectors();
