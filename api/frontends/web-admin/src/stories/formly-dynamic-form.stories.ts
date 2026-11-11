import { JsonPipe } from '@angular/common';
import { Component, input } from '@angular/core';
import { FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { FormlyFieldConfig, FormlyForm } from '@ngx-formly/core';
import type { Meta, StoryObj } from '@storybook/angular';
import { applicationConfig } from '@storybook/angular';
import { provideAppFormly } from 'app/shared/forms/formly.provider';

/** Story host: the same wiring a feature component uses. */
@Component({
    selector: 'story-formly-form',
    imports: [ReactiveFormsModule, FormlyForm, MatButtonModule, JsonPipe],
    template: `
        <form [formGroup]="form" (ngSubmit)="submit()" style="max-width: 560px">
            <formly-form [form]="form" [fields]="fields()" [model]="model" />
            <button mat-flat-button color="primary" type="submit">
                Submit
            </button>
        </form>
        @if (submitted) {
            <pre style="margin-top: 24px; font-size: 13px">{{
                submitted | json
            }}</pre>
        }
    `,
})
class StoryFormlyFormComponent {
    fields = input.required<FormlyFieldConfig[]>();
    form = new FormGroup({});
    model = {};
    submitted: unknown = null;

    submit(): void {
        this.submitted = this.form.valid ? this.form.value : null;
    }
}

const enquiryFields: FormlyFieldConfig[] = [
    {
        fieldGroupClassName: 'grid grid-cols-1 sm:grid-cols-2 gap-x-4',
        fieldGroup: [
            {
                key: 'guardianName',
                type: 'input',
                props: { label: 'Guardian name', required: true, minLength: 2 },
            },
            {
                key: 'email',
                type: 'input',
                props: {
                    label: 'Email',
                    type: 'email',
                    required: true,
                    pattern: /^\S+@\S+\.\S+$/,
                },
            },
        ],
    },
    {
        key: 'grade',
        type: 'select',
        props: {
            label: 'Grade applying for',
            required: true,
            options: [
                'PP1',
                'PP2',
                'Grade 1',
                'Grade 2',
                'Grade 3',
                'Grade 4',
            ].map((g) => ({ label: g, value: g })),
        },
    },
    {
        key: 'siblings',
        type: 'radio',
        props: {
            label: 'Siblings already at the school?',
            options: [
                { label: 'Yes', value: true },
                { label: 'No', value: false },
            ],
        },
    },
    {
        key: 'siblingNames',
        type: 'input',
        props: { label: 'Sibling names' },
        // Shown only when the guardian answers yes: the dynamic part of a Formly form.
        expressions: { hide: 'model.siblings !== true' },
    },
    {
        key: 'notes',
        type: 'textarea',
        props: { label: 'Anything we should know?', rows: 3, maxLength: 500 },
    },
    {
        key: 'consent',
        type: 'checkbox',
        defaultValue: false,
        // Formly's Material checkbox is tri-state by default; consent must be a plain yes/no.
        props: {
            label: 'I agree to be contacted about this enquiry',
            required: true,
            indeterminate: false,
        },
        validators: { validation: [Validators.requiredTrue] },
    },
];

const meta: Meta<StoryFormlyFormComponent> = {
    title: 'Forms/Formly dynamic form',
    component: StoryFormlyFormComponent,
    decorators: [applicationConfig({ providers: [provideAppFormly()] })],
    parameters: {
        layout: 'padded',
        docs: {
            description: {
                component:
                    'Forms described as `FormlyFieldConfig[]` and rendered with Angular Material. `provideAppFormly()` (added to the providers of the lazy route that uses Formly) supplies the Material field types and shared validation messages. Errors appear once a field is touched or the form is submitted.',
            },
        },
    },
};

export default meta;
type Story = StoryObj<StoryFormlyFormComponent>;

export const AdmissionsEnquiry: Story = {
    name: 'Admissions enquiry',
    args: { fields: enquiryFields },
};
