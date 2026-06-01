import { NgClass } from '@angular/common';
import { ChangeDetectionStrategy, Component } from '@angular/core';
import {
    FormsModule,
    ReactiveFormsModule,
    UntypedFormBuilder,
    UntypedFormGroup,
} from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import {
    ApplicationFee,
    ApplicationFeeStatus,
} from 'app/core/admissions/admissions.types';

@Component({
    selector: 'app-portal-apply',
    standalone: true,
    imports: [
        FormsModule,
        ReactiveFormsModule,
        MatButtonModule,
        MatFormFieldModule,
        MatIconModule,
        MatInputModule,
        MatProgressBarModule,
        NgClass,
    ],
    template: `
        <div class="flex flex-auto flex-col px-6 py-12 sm:px-10 lg:px-16">
            <div class="mx-auto w-full max-w-6xl">
                <div class="mb-8 flex flex-col">
                    <span
                        class="text-secondary mb-2 text-sm font-semibold uppercase tracking-wider"
                        >Step 1 of 6 - Account</span
                    >
                    <h1
                        class="text-default text-3xl font-bold tracking-tight sm:text-4xl"
                    >
                        Apply for Fall 2026
                    </h1>
                </div>

                <div class="mb-10 w-full">
                    <mat-progress-bar
                        mode="determinate"
                        [value]="16"
                        aria-label="Application progress: step 1 of 6, account information"
                    ></mat-progress-bar>
                </div>

                <div class="grid grid-cols-1 gap-8 lg:grid-cols-[1fr_320px]">
                    <!-- Left Column - Form Card -->
                    <div class="bg-card rounded-2xl border p-8 shadow-sm">
                        <form [formGroup]="form" class="flex flex-col">
                            <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
                                <mat-form-field
                                    appearance="outline"
                                    class="w-full"
                                >
                                    <mat-label>First name</mat-label>
                                    <input
                                        matInput
                                        formControlName="firstName"
                                        autocomplete="given-name"
                                    />
                                </mat-form-field>

                                <mat-form-field
                                    appearance="outline"
                                    class="w-full"
                                >
                                    <mat-label>Last name</mat-label>
                                    <input
                                        matInput
                                        formControlName="lastName"
                                        autocomplete="family-name"
                                    />
                                </mat-form-field>

                                <mat-form-field
                                    appearance="outline"
                                    class="w-full"
                                >
                                    <mat-label>Email</mat-label>
                                    <input
                                        matInput
                                        type="email"
                                        formControlName="email"
                                        autocomplete="email"
                                    />
                                </mat-form-field>

                                <mat-form-field
                                    appearance="outline"
                                    class="w-full"
                                >
                                    <mat-label>Phone</mat-label>
                                    <input
                                        matInput
                                        type="tel"
                                        formControlName="phone"
                                        autocomplete="tel"
                                    />
                                </mat-form-field>

                                <mat-form-field
                                    appearance="outline"
                                    class="w-full"
                                >
                                    <mat-label>Password</mat-label>
                                    <input
                                        matInput
                                        type="password"
                                        formControlName="password"
                                        autocomplete="new-password"
                                        aria-describedby="portal-password-help"
                                    />
                                </mat-form-field>

                                <mat-form-field
                                    appearance="outline"
                                    class="w-full"
                                >
                                    <mat-label>Confirm password</mat-label>
                                    <input
                                        matInput
                                        type="password"
                                        formControlName="confirmPassword"
                                        autocomplete="new-password"
                                    />
                                </mat-form-field>
                            </div>

                            <p
                                id="portal-password-help"
                                class="text-secondary mt-2 text-sm"
                            >
                                Use a memorable password for the applicant
                                account; stronger password rules will be
                                enforced when account APIs are integrated.
                            </p>

                            <div class="mt-8 flex items-center justify-between">
                                <button mat-stroked-button type="button">
                                    <mat-icon
                                        svgIcon="heroicons_outline:arrow-left"
                                        class="icon-size-5"
                                        aria-hidden="true"
                                    ></mat-icon>
                                    <span class="ml-2">Back</span>
                                </button>
                                <button
                                    mat-flat-button
                                    color="primary"
                                    type="button"
                                >
                                    <span class="mr-2">Continue</span>
                                    <mat-icon
                                        svgIcon="heroicons_outline:arrow-right"
                                        class="icon-size-5"
                                        aria-hidden="true"
                                    ></mat-icon>
                                </button>
                            </div>
                        </form>
                    </div>

                    <!-- Right Column - Sidebar -->
                    <div class="flex flex-col gap-6">
                        <!-- Progress Card -->
                        <div class="bg-card rounded-2xl border p-6 shadow-sm">
                            <h2
                                class="text-secondary mb-6 text-sm font-bold uppercase tracking-wider"
                            >
                                Progress
                            </h2>
                            <div class="flex flex-col gap-4">
                                @for (step of steps; track step.number) {
                                    <div class="flex items-center gap-4">
                                        <div
                                            class="flex h-8 w-8 items-center justify-center rounded-full text-sm font-bold"
                                            [ngClass]="
                                                step.active
                                                    ? 'bg-primary text-on-primary'
                                                    : 'text-secondary border border-gray-300'
                                            "
                                        >
                                            {{ step.number }}
                                        </div>
                                        <span
                                            class="text-base"
                                            [ngClass]="
                                                step.active
                                                    ? 'text-default font-bold'
                                                    : 'text-secondary font-medium'
                                            "
                                        >
                                            {{ step.label }}
                                        </span>
                                    </div>
                                }
                            </div>
                        </div>

                        <!-- Application Fee Card -->
                        <div class="bg-card rounded-2xl border p-6 shadow-sm">
                            <div class="flex items-start justify-between gap-3">
                                <div>
                                    <h2
                                        class="text-secondary text-sm font-bold uppercase tracking-wider"
                                    >
                                        Application fee
                                    </h2>
                                    <p class="text-secondary mt-2 text-sm">
                                        Fee status is tracked separately from
                                        event payments and discounts.
                                    </p>
                                </div>
                                <span
                                    class="rounded-full px-2.5 py-1 text-xs font-semibold"
                                    [class]="feeStatusClass(reviewFee.status)"
                                    [attr.aria-label]="
                                        'Application fee status: ' +
                                        formatFeeStatus(reviewFee.status)
                                    "
                                >
                                    {{ formatFeeStatus(reviewFee.status) }}
                                </span>
                            </div>
                            <div
                                class="mt-5 rounded-2xl bg-gray-50 p-4 dark:bg-gray-800"
                            >
                                <div
                                    class="text-secondary text-xs font-semibold uppercase tracking-wide"
                                >
                                    Due at review
                                </div>
                                <div
                                    class="text-default mt-1 text-2xl font-bold"
                                >
                                    {{ formatAmount(reviewFee) }}
                                </div>
                                <div class="text-secondary mt-1 text-xs">
                                    Provider seam: {{ reviewFee.provider }} ·
                                    status updates later via API
                                </div>
                            </div>
                        </div>

                        <!-- Need Help Card -->
                        <div
                            class="rounded-2xl bg-gray-50 p-6 dark:bg-gray-800"
                        >
                            <h2 class="text-default mb-2 text-lg font-semibold">
                                Need help?
                            </h2>
                            <p class="text-secondary leading-relaxed">
                                Email admissions&#64;northbrook.edu or chat with
                                a counselor weekdays 9am-6pm ET.
                            </p>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    `,
    changeDetection: ChangeDetectionStrategy.OnPush,
    host: { class: 'flex w-full flex-auto flex-col' },
})
export class PortalApplyComponent {
    readonly reviewFee: ApplicationFee = {
        id: 'fee-application-draft',
        applicationID: 'APP-DRAFT',
        amountCents: 15000,
        currency: 'USD',
        status: 'PENDING',
        provider: 'stripe',
        dueAt: 'Before final submission',
        auditTrail: [],
    };

    readonly steps = [
        { number: 1, label: 'Account', active: true },
        { number: 2, label: 'Program', active: false },
        { number: 3, label: 'Background', active: false },
        { number: 4, label: 'Essay', active: false },
        { number: 5, label: 'Documents', active: false },
        { number: 6, label: 'Review', active: false },
    ];

    form: UntypedFormGroup;

    constructor(private _formBuilder: UntypedFormBuilder) {
        this.form = this._formBuilder.group({
            firstName: [''],
            lastName: [''],
            email: [''],
            phone: [''],
            password: [''],
            confirmPassword: [''],
        });
    }

    feeStatusClass(status: ApplicationFeeStatus): string {
        switch (status) {
            case 'PAID':
                return 'bg-green-100 text-green-700';
            case 'PENDING':
                return 'bg-amber-100 text-amber-700';
            case 'FAILED':
                return 'bg-red-100 text-red-700';
            case 'WAIVED':
                return 'bg-purple-100 text-purple-700';
            case 'REFUNDED':
                return 'bg-blue-100 text-blue-700';
            case 'NOT_REQUIRED':
                return 'bg-slate-100 text-secondary';
        }
    }

    formatFeeStatus(status: ApplicationFeeStatus): string {
        return status.replaceAll('_', ' ');
    }

    formatAmount(fee: ApplicationFee): string {
        if (fee.status === 'NOT_REQUIRED') {
            return 'No fee';
        }

        return new Intl.NumberFormat('en-US', {
            style: 'currency',
            currency: fee.currency,
        }).format(fee.amountCents / 100);
    }
}
