import { Component } from '@angular/core';
import { TestBed } from '@angular/core/testing';
import { FormControl, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { provideNoopAnimations } from '@angular/platform-browser/animations';
import { FormlyFieldConfig, FormlyForm } from '@ngx-formly/core';
import { provideAppFormly } from 'app/shared/forms/formly.provider';
import { afterEach, describe, expect, it } from 'vitest';

import { filterOptions, SearchableSelectComponent, SelectOption } from './searchable-select.component';
import { provideSearchableSelectFormly } from './searchable-select.formly';

const teachers: SelectOption<string>[] = [
    { label: 'Amina Otieno', value: 't1' },
    { label: 'José Mwangi', value: 't2' },
    { label: 'Grace Wanjiru', value: 't3' },
];

describe('filterOptions', () => {
    it('matches case- and accent-insensitively', () => {
        expect(filterOptions(teachers, 'jose').map((o) => o.value)).toEqual(['t2']);
        expect(filterOptions(teachers, '  WANJ ').map((o) => o.value)).toEqual(['t3']);
    });

    it('returns every option for an empty query', () => {
        expect(filterOptions(teachers, '')).toHaveLength(3);
    });
});

@Component({
    imports: [SearchableSelectComponent],
    template: `<app-searchable-select label="Class teacher" [control]="control" [options]="options" />`,
})
class ReactiveHostComponent {
    control = new FormControl<string | null>(null, Validators.required);
    options = teachers;
}

function overlayOptions(): string[] {
    // The first mat-option hosts the search box (class contains-mat-select-search); skip it.
    return Array.from(
        document.querySelectorAll('.cdk-overlay-container mat-option:not(.contains-mat-select-search)')
    ).map((o) => o.textContent?.trim() ?? '');
}

describe('SearchableSelectComponent', () => {
    afterEach(() => document.querySelector('.cdk-overlay-container')?.replaceChildren());

    it('filters the open panel as the user types and writes the chosen value', async () => {
        TestBed.configureTestingModule({ imports: [ReactiveHostComponent], providers: [provideNoopAnimations()] });
        const fixture = TestBed.createComponent(ReactiveHostComponent);
        fixture.detectChanges();
        const el: HTMLElement = fixture.nativeElement;

        expect(el.querySelector('mat-label')?.textContent).toContain('Class teacher');
        expect(el.querySelector('.mat-mdc-form-field-required-marker')).not.toBeNull();

        (el.querySelector('.mat-mdc-select-trigger') as HTMLElement).click();
        fixture.detectChanges();
        await fixture.whenStable();
        expect(overlayOptions()).toEqual(['Amina Otieno', 'José Mwangi', 'Grace Wanjiru']);

        const search = document.querySelector('.cdk-overlay-container input.mat-select-search-input:not(.mat-select-search-hidden)') as HTMLInputElement;
        search.value = 'grace';
        search.dispatchEvent(new Event('input'));
        fixture.detectChanges();
        await fixture.whenStable();
        expect(overlayOptions()).toEqual(['Grace Wanjiru']);

        const grace = Array.from(document.querySelectorAll<HTMLElement>('.cdk-overlay-container mat-option'))
            .find((o) => o.textContent?.includes('Grace'))!;
        grace.click();
        fixture.detectChanges();
        expect(fixture.componentInstance.control.value).toBe('t3');
    });
});

@Component({
    imports: [ReactiveFormsModule, FormlyForm],
    template: `<form [formGroup]="form"><formly-form [form]="form" [fields]="fields" [model]="model" /></form>`,
})
class FormlyHostComponent {
    form = new FormGroup({});
    model = {};
    fields: FormlyFieldConfig[] = [
        { key: 'teacher', type: 'searchable-select', props: { label: 'Class teacher', required: true, options: teachers } },
    ];
}

describe('searchable-select Formly type', () => {
    it('renders with the required marker and shows the shared message once touched', () => {
        TestBed.configureTestingModule({
            imports: [FormlyHostComponent],
            providers: [provideNoopAnimations(), provideAppFormly(), provideSearchableSelectFormly()],
        });
        const fixture = TestBed.createComponent(FormlyHostComponent);
        fixture.detectChanges();
        const el: HTMLElement = fixture.nativeElement;

        expect(el.querySelectorAll('mat-form-field')).toHaveLength(1);
        expect(el.querySelector('.mat-mdc-form-field-required-marker')).not.toBeNull();
        expect(el.querySelector('mat-error')).toBeNull();

        fixture.componentInstance.form.get('teacher')!.markAsTouched();
        fixture.detectChanges();

        expect(el.querySelector('mat-error')?.textContent).toContain('This field is required.');
    });
});
