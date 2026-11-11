import { JsonPipe } from '@angular/common';
import { Component, input } from '@angular/core';
import { FormControl, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { FormlyFieldConfig, FormlyForm } from '@ngx-formly/core';
import type { Meta, StoryObj } from '@storybook/angular';
import { applicationConfig, moduleMetadata } from '@storybook/angular';
import {
    provideSearchableSelectFormly,
    SearchableSelectComponent,
    SelectOption,
} from 'app/shared/components/searchable-select';
import { provideAppFormly } from 'app/shared/forms/formly.provider';

const teachers: SelectOption<string>[] = [
    'Amina Otieno', 'Brian Kiprotich', 'Catherine Njeri', 'David Ochieng', 'Esther Achieng',
    'Faith Mutua', 'Grace Wanjiru', 'Hassan Abdi', 'Irene Chebet', 'José Mwangi',
    'Kevin Kamau', 'Lucy Akinyi', 'Mercy Nyambura', 'Nelson Omondi', 'Purity Wairimu',
].map((name, i) => ({ label: name, value: `t${i + 1}` }));

@Component({
    selector: 'story-searchable-select',
    imports: [SearchableSelectComponent, JsonPipe],
    template: `
        <div style="max-width: 420px">
            <app-searchable-select
                [label]="label()"
                [control]="control"
                [options]="options()"
                [multiple]="multiple()"
                hint="Type to filter the list"
            >
                <span error>Choose a class teacher.</span>
            </app-searchable-select>
            <pre style="font-size: 13px">value: {{ control.value | json }}</pre>
        </div>
    `,
})
class StorySearchableSelectComponent {
    label = input('Class teacher');
    options = input<SelectOption<string>[]>(teachers);
    multiple = input(false);
    control = new FormControl<string | string[] | null>(null, Validators.required);
}

const meta: Meta<StorySearchableSelectComponent> = {
    title: 'Forms/Searchable select',
    component: StorySearchableSelectComponent,
    parameters: {
        layout: 'padded',
        docs: {
            description: {
                component:
                    '`<app-searchable-select>` from `app/shared/components/searchable-select`: a Material select with an ngx-mat-select-search box, bound to any `FormControl`. Filtering ignores case and accents ("jose" finds "José"). In Formly forms use `type: \'searchable-select\'` after `provideSearchableSelectFormly()`.',
            },
        },
    },
    args: { label: 'Class teacher', options: teachers, multiple: false },
};

export default meta;
type Story = StoryObj<StorySearchableSelectComponent>;

export const Single: Story = {};

export const Multiple: Story = {
    args: { label: 'Teachers on the review panel', multiple: true },
};

@Component({
    selector: 'story-searchable-select-formly',
    imports: [ReactiveFormsModule, FormlyForm, JsonPipe],
    template: `
        <form [formGroup]="form" style="max-width: 420px">
            <formly-form [form]="form" [fields]="fields" [model]="model" />
        </form>
        <pre style="font-size: 13px">{{ model | json }}</pre>
    `,
})
class StorySearchableSelectFormlyComponent {
    form = new FormGroup({});
    model = {};
    fields: FormlyFieldConfig[] = [
        {
            key: 'classTeacher',
            type: 'searchable-select',
            props: { label: 'Class teacher', required: true, options: teachers },
        },
        {
            key: 'subjectTeachers',
            type: 'searchable-select',
            props: {
                label: 'Subject teachers',
                multiple: true,
                options: teachers,
                description: 'Select every teacher who covers this class',
            },
        },
    ];
}

export const InFormly: Story = {
    name: 'In a Formly form',
    parameters: { controls: { disable: true } },
    render: () => ({ template: '<story-searchable-select-formly />' }),
    decorators: [
        moduleMetadata({ imports: [StorySearchableSelectFormlyComponent] }),
        applicationConfig({ providers: [provideAppFormly(), provideSearchableSelectFormly()] }),
    ],
};
