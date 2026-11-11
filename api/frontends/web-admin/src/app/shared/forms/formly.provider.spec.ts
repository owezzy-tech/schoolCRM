import { Component } from '@angular/core';
import { TestBed } from '@angular/core/testing';
import { FormGroup, ReactiveFormsModule } from '@angular/forms';
import { provideNoopAnimations } from '@angular/platform-browser/animations';
import { FormlyFieldConfig, FormlyForm } from '@ngx-formly/core';
import { describe, expect, it } from 'vitest';

import { provideAppFormly } from './formly.provider';

@Component({
    imports: [ReactiveFormsModule, FormlyForm],
    template: `<form [formGroup]="form"><formly-form [form]="form" [fields]="fields" [model]="model" /></form>`,
})
class HostComponent {
    form = new FormGroup({});
    model = { email: '' };
    fields: FormlyFieldConfig[] = [
        { key: 'email', type: 'input', props: { label: 'Email', required: true, type: 'email' } },
    ];
}

function render() {
    TestBed.configureTestingModule({
        imports: [HostComponent],
        providers: [provideNoopAnimations(), provideAppFormly()],
    });
    const fixture = TestBed.createComponent(HostComponent);
    fixture.detectChanges();
    return fixture;
}

describe('provideAppFormly', () => {
    it('renders Material fields with a persistent label', () => {
        const el: HTMLElement = render().nativeElement;

        expect(el.querySelector('mat-form-field input')).not.toBeNull();
        expect(el.querySelector('mat-label')?.textContent).toContain('Email');
    });

    it('shows no error until the field is touched, then the shared message', () => {
        const fixture = render();
        const el: HTMLElement = fixture.nativeElement;
        expect(el.querySelector('mat-error')).toBeNull();

        fixture.componentInstance.form.get('email')!.markAsTouched();
        fixture.detectChanges();

        expect(el.querySelector('mat-error')?.textContent).toContain('This field is required.');
    });
});
