import { ChangeDetectionStrategy, Component } from '@angular/core';
import { FormControl } from '@angular/forms';
import {
    FieldType,
    FieldTypeConfig,
    FormlyFieldProps,
    FormlyValidationMessage,
    provideFormlyConfig,
} from '@ngx-formly/core';
import {
    SearchableSelectComponent,
    SelectOption,
} from './searchable-select.component';

interface SearchableSelectProps extends FormlyFieldProps {
    options: SelectOption[];
    multiple?: boolean;
    searchPlaceholder?: string;
    noEntriesLabel?: string;
}

/**
 * Formly field type `searchable-select`. The component draws its own mat-form-field,
 * so the type registers without Formly's form-field wrapper.
 */
@Component({
    selector: 'formly-field-searchable-select',
    imports: [SearchableSelectComponent, FormlyValidationMessage],
    changeDetection: ChangeDetectionStrategy.OnPush,
    template: `
        <app-searchable-select
            [label]="props.label ?? ''"
            [control]="control"
            [options]="props.options"
            [multiple]="props.multiple ?? false"
            [required]="props.required ?? false"
            [hint]="props.description ?? ''"
            [searchPlaceholder]="props.searchPlaceholder ?? 'Search'"
            [noEntriesLabel]="props.noEntriesLabel ?? 'No matches'"
        >
            <formly-validation-message error [field]="field" />
        </app-searchable-select>
    `,
})
export class FormlySearchableSelectType extends FieldType<
    FieldTypeConfig<SearchableSelectProps>
> {
    protected get control(): FormControl {
        return this.formControl as FormControl;
    }
}

/** Registers `type: 'searchable-select'`. Use alongside provideAppFormly() in the same route. */
export function provideSearchableSelectFormly() {
    return provideFormlyConfig({
        types: [{ name: 'searchable-select', component: FormlySearchableSelectType }],
    });
}
