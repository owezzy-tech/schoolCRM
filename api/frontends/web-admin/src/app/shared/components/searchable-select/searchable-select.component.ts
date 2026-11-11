import {
    ChangeDetectionStrategy,
    Component,
    computed,
    input,
    signal,
} from '@angular/core';
import { FormControl, FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatSelectModule } from '@angular/material/select';
import { ElimuSelectSearchComponent } from 'app/shared/components/select-search';

export interface SelectOption<T = unknown> {
    label: string;
    value: T;
}

/** Case- and accent-insensitive match, so "jose" finds "José". */
export function filterOptions<T>(
    options: readonly SelectOption<T>[],
    query: string
): SelectOption<T>[] {
    const fold = (s: string) =>
        s
            .normalize('NFD')
            .replace(/\p{Diacritic}/gu, '')
            .toLowerCase();
    const q = fold(query.trim());
    return q ? options.filter((o) => fold(o.label).includes(q)) : [...options];
}

/**
 * Material select with a search box (<elimu-select-search>, our fork of ngx-mat-select-search), bound to any FormControl.
 *
 * ```html
 * <app-searchable-select label="Class" [control]="form.controls.classId" [options]="classes" required>
 *     <span error>Choose a class.</span>
 * </app-searchable-select>
 * ```
 *
 * In Formly forms use the `searchable-select` field type instead (provideSearchableSelectFormly).
 */
@Component({
    selector: 'app-searchable-select',
    imports: [
        FormsModule,
        ReactiveFormsModule,
        MatFormFieldModule,
        MatSelectModule,
        ElimuSelectSearchComponent,
    ],
    changeDetection: ChangeDetectionStrategy.OnPush,
    host: { class: 'block' },
    template: `
        <mat-form-field class="w-full">
            <mat-label>{{ label() }}</mat-label>
            <mat-select
                [formControl]="control()"
                [multiple]="multiple()"
                [required]="required()"
            >
                <mat-option>
                    <elimu-select-search
                        [ngModel]="query()"
                        (ngModelChange)="query.set($event)"
                        [placeholderLabel]="searchPlaceholder()"
                        [noEntriesFoundLabel]="noEntriesLabel()"
                        [ariaLabel]="'Search ' + label()"
                        closeSvgIcon="heroicons_outline:x-mark"
                    />
                </mat-option>
                @for (option of filtered(); track option.value) {
                    <mat-option [value]="option.value">{{
                        option.label
                    }}</mat-option>
                }
            </mat-select>
            @if (hint()) {
                <mat-hint>{{ hint() }}</mat-hint>
            }
            <mat-error><ng-content select="[error]" /></mat-error>
        </mat-form-field>
    `,
})
export class SearchableSelectComponent<T = unknown> {
    readonly label = input.required<string>();
    readonly control = input.required<FormControl<T | T[] | null>>();
    readonly options = input.required<readonly SelectOption<T>[]>();
    readonly multiple = input(false);
    /** Shows the required marker, as on mat-select; the control still owns validation. */
    readonly required = input(false);
    readonly hint = input('');
    readonly searchPlaceholder = input('Search');
    readonly noEntriesLabel = input('No matches');

    protected readonly query = signal('');
    protected readonly filtered = computed(() =>
        filterOptions(this.options(), this.query())
    );
}
