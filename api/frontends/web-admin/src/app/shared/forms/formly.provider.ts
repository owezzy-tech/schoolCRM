import { EnvironmentProviders, makeEnvironmentProviders } from '@angular/core';
import { FormlyFieldConfig, provideFormlyCore } from '@ngx-formly/core';
import { withFormlyMaterial } from '@ngx-formly/material';

/** Messages Formly shows under a field once it is touched or its form is submitted. */
export const FORMLY_VALIDATION_MESSAGES = [
    { name: 'required', message: 'This field is required.' },
    { name: 'email', message: 'Enter a valid email address.' },
    {
        name: 'minLength',
        message: (_: unknown, field: FormlyFieldConfig) =>
            `Enter at least ${field.props?.minLength} characters.`,
    },
    {
        name: 'maxLength',
        message: (_: unknown, field: FormlyFieldConfig) =>
            `Enter no more than ${field.props?.maxLength} characters.`,
    },
    {
        name: 'min',
        message: (_: unknown, field: FormlyFieldConfig) =>
            `Enter ${field.props?.min} or more.`,
    },
    {
        name: 'max',
        message: (_: unknown, field: FormlyFieldConfig) =>
            `Enter ${field.props?.max} or less.`,
    },
    { name: 'pattern', message: 'The value is not in the expected format.' },
];

/**
 * Formly for dynamic forms, rendered with Angular Material field types.
 *
 * Add it to the `providers` of the lazy route that renders Formly forms, not to app.config.ts:
 * the Material field types cost about 140 kB, which would otherwise land in the initial bundle
 * of every page. Extra field types go in the same route through provideFormlyConfig().
 */
export function provideAppFormly(): EnvironmentProviders {
    return makeEnvironmentProviders([
        provideFormlyCore([
            ...withFormlyMaterial(),
            { validationMessages: FORMLY_VALIDATION_MESSAGES },
        ]),
    ]);
}
