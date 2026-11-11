/**
 * Forked from ngx-mat-select-search 9.0.0 (https://github.com/bithost-gmbh/ngx-mat-select-search,
 * tag 9.0.0, commit 9fa3f73) and renamed for SchoolCRM. Copyright (c) 2018 Bithost GmbH.
 * MIT licence: see LICENSE in this directory.
 */

import { InjectionToken, Provider } from '@angular/core';
import { ElimuSelectSearchComponent } from './elimu-select-search.component';

/** List of inputs of ElimuSelectSearchComponent that can be configured with a global default. */
export const configurableDefaultOptions = [
    'ariaLabel',
    'clearSearchInput',
    'closeIcon',
    'closeSvgIcon',
    'disableInitialFocus',
    'disableScrollToActiveOnOptionsChanged',
    'enableClearOnEscapePressed',
    'hideClearSearchButton',
    'noEntriesFoundLabel',
    'placeholderLabel',
    'preventHomeEndKeyPropagation',
    'searching',
] as const;

export type ConfigurableDefaultOptions =
    (typeof configurableDefaultOptions)[number];

export const ELIMU_SELECT_SEARCH_DEFAULT_OPTIONS =
    new InjectionToken<ElimuSelectSearchOptions>(
        'elimu-select-search-default-options'
    );

/** Global configurable options for ElimuSelectSearch. */
export type ElimuSelectSearchOptions = Readonly<
    Partial<Pick<ElimuSelectSearchComponent, ConfigurableDefaultOptions>>
>;

/**
 * Provides default options for the `ElimuSelectSearchComponent`. e.g.
 *
 * ```typescript
 * providers: [
 *   provideElimuSelectSearchOptions({
 *     closeIcon: 'delete',
 *     noEntriesFoundLabel: 'No options found'
 *   })
 * ]
 * ```
 *
 * See the corresponding inputs of `ElimuSelectSearchComponent` for documentation.
 */
export function provideElimuSelectSearchOptions(
    options: ElimuSelectSearchOptions
): Provider {
    return {
        provide: ELIMU_SELECT_SEARCH_DEFAULT_OPTIONS,
        useValue: options,
    };
}
