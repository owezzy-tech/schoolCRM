import {
    ChangeDetectionStrategy,
    Component,
    inject,
    signal,
} from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { FuseAlertComponent, FuseAlertType } from '@fuse/components/alert';
import { AdmissionsService } from 'app/core/admissions/admissions.service';
import { InquiryRequest } from 'app/core/admissions/admissions.types';
import { jsonApiErrorMessage } from 'app/core/api/json-api';
import { finalize } from 'rxjs';

@Component({
    selector: 'app-portal-inquiry',
    imports: [
        FuseAlertComponent,
        MatButtonModule,
        MatFormFieldModule,
        MatIconModule,
        MatInputModule,
        MatProgressSpinnerModule,
        ReactiveFormsModule,
    ],
    template: `
        <div class="flex flex-auto flex-col px-6 py-12 sm:px-10 lg:px-16">
            <div class="mx-auto w-full max-w-6xl text-center">
                <h1
                    class="text-default text-3xl font-bold tracking-tight sm:text-4xl"
                >
                    Ask admissions
                </h1>
                <p class="text-secondary mb-8 mt-2 text-lg">
                    Send a question about KUCCPS placement, self-sponsored
                    intake, KCSE documents, or M-Pesa fee status. We will get
                    back to you within one business day.
                </p>

                <div
                    class="bg-card mx-auto max-w-3xl rounded-2xl border p-8 text-left shadow-sm"
                >
                    @if (alert(); as currentAlert) {
                        <fuse-alert
                            class="mb-6"
                            [appearance]="'outline'"
                            [showIcon]="true"
                            [type]="currentAlert.type"
                        >
                            {{ currentAlert.message }}
                        </fuse-alert>
                    }

                    <form
                        [formGroup]="inquiryForm"
                        class="flex flex-col gap-4"
                        (ngSubmit)="submitInquiry()"
                    >
                        <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
                            <mat-form-field appearance="outline" class="w-full">
                                <mat-label>First name</mat-label>
                                <input
                                    matInput
                                    autocomplete="given-name"
                                    formControlName="firstName"
                                />
                                @if (
                                    inquiryForm.controls.firstName.hasError(
                                        'required'
                                    )
                                ) {
                                    <mat-error
                                        >First name is required</mat-error
                                    >
                                }
                            </mat-form-field>

                            <mat-form-field appearance="outline" class="w-full">
                                <mat-label>Last name</mat-label>
                                <input
                                    matInput
                                    autocomplete="family-name"
                                    formControlName="lastName"
                                />
                                @if (
                                    inquiryForm.controls.lastName.hasError(
                                        'required'
                                    )
                                ) {
                                    <mat-error>Last name is required</mat-error>
                                }
                            </mat-form-field>

                            <mat-form-field appearance="outline" class="w-full">
                                <mat-label>Date of birth</mat-label>
                                <input
                                    matInput
                                    type="date"
                                    formControlName="dateOfBirth"
                                />
                                @if (
                                    inquiryForm.controls.dateOfBirth.hasError(
                                        'required'
                                    )
                                ) {
                                    <mat-error
                                        >Date of birth is required</mat-error
                                    >
                                }
                            </mat-form-field>

                            <mat-form-field appearance="outline" class="w-full">
                                <mat-label>Email</mat-label>
                                <input
                                    matInput
                                    type="email"
                                    autocomplete="email"
                                    formControlName="primaryEmail"
                                />
                                @if (
                                    inquiryForm.controls.primaryEmail.hasError(
                                        'required'
                                    )
                                ) {
                                    <mat-error>Email is required</mat-error>
                                }
                                @if (
                                    inquiryForm.controls.primaryEmail.hasError(
                                        'email'
                                    )
                                ) {
                                    <mat-error>Enter a valid email</mat-error>
                                }
                            </mat-form-field>

                            <mat-form-field appearance="outline" class="w-full">
                                <mat-label>Phone</mat-label>
                                <input
                                    matInput
                                    type="tel"
                                    autocomplete="tel"
                                    formControlName="primaryPhone"
                                />
                                @if (
                                    inquiryForm.controls.primaryPhone.hasError(
                                        'required'
                                    )
                                ) {
                                    <mat-error>Phone is required</mat-error>
                                }
                            </mat-form-field>

                            <mat-form-field appearance="outline" class="w-full">
                                <mat-label>Programme of interest</mat-label>
                                <input
                                    matInput
                                    formControlName="programOfInterest"
                                    placeholder="e.g. Bachelor of Commerce"
                                />
                            </mat-form-field>

                            <mat-form-field appearance="outline" class="w-full">
                                <mat-label>Intake term</mat-label>
                                <input
                                    matInput
                                    formControlName="termOfInterest"
                                    placeholder="2026 Main Intake"
                                />
                            </mat-form-field>
                        </div>

                        <mat-form-field appearance="outline" class="w-full">
                            <mat-label>Your question</mat-label>
                            <textarea
                                matInput
                                rows="6"
                                formControlName="message"
                                placeholder="Tell us what you'd like to know..."
                            ></textarea>
                            @if (
                                inquiryForm.controls.message.hasError(
                                    'required'
                                )
                            ) {
                                <mat-error>Your question is required</mat-error>
                            }
                        </mat-form-field>

                        <p class="text-secondary text-sm">
                            By submitting, you agree to our privacy policy.
                        </p>

                        <div class="mt-4 flex justify-end">
                            <button
                                mat-flat-button
                                color="primary"
                                type="submit"
                                [disabled]="submitting()"
                            >
                                <mat-icon
                                    svgIcon="heroicons_outline:paper-airplane"
                                    class="mr-2 icon-size-5"
                                ></mat-icon>
                                @if (!submitting()) {
                                    <span>Send inquiry</span>
                                } @else {
                                    <mat-progress-spinner
                                        [diameter]="24"
                                        [mode]="'indeterminate'"
                                    ></mat-progress-spinner>
                                }
                            </button>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    `,
    changeDetection: ChangeDetectionStrategy.OnPush,
    host: { class: 'flex w-full flex-auto flex-col' },
})
export class PortalInquiryComponent {
    private readonly admissionsService = inject(AdmissionsService);
    private readonly formBuilder = inject(FormBuilder);

    readonly alert = signal<{ type: FuseAlertType; message: string } | null>(
        null
    );
    readonly submitting = signal(false);
    readonly inquiryForm = this.formBuilder.nonNullable.group({
        firstName: ['', Validators.required],
        lastName: ['', Validators.required],
        dateOfBirth: ['', Validators.required],
        primaryEmail: ['', [Validators.required, Validators.email]],
        primaryPhone: ['', Validators.required],
        programOfInterest: [''],
        termOfInterest: [''],
        message: ['', Validators.required],
    });

    submitInquiry(): void {
        if (this.inquiryForm.invalid) {
            this.inquiryForm.markAllAsTouched();
            return;
        }

        this.alert.set(null);
        this.submitting.set(true);

        this.admissionsService
            .createInquiry(this.inquiryRequest())
            .pipe(finalize(() => this.submitting.set(false)))
            .subscribe({
                next: () => {
                    this.inquiryForm.reset();
                    this.alert.set({
                        type: 'success',
                        message:
                            'Your inquiry has been sent. Our admissions team will reply within one business day.',
                    });
                },
                error: (error) => {
                    this.alert.set({
                        type: 'error',
                        message: jsonApiErrorMessage(
                            error,
                            'Unable to submit your inquiry. Please try again.'
                        ),
                    });
                },
            });
    }

    private inquiryRequest(): InquiryRequest {
        const value = this.inquiryForm.getRawValue();

        return {
            firstName: value.firstName.trim(),
            lastName: value.lastName.trim(),
            dateOfBirth: new Date(value.dateOfBirth).toISOString(),
            primaryEmail: value.primaryEmail.trim(),
            primaryPhone: value.primaryPhone.trim(),
            programOfInterest: this.optionalString(value.programOfInterest),
            termOfInterest: this.optionalString(value.termOfInterest),
            source: 'PORTAL_INQUIRY_FORM',
            utmSource: this.queryParam('utm_source'),
            utmMedium: this.queryParam('utm_medium'),
            utmCampaign: this.queryParam('utm_campaign'),
            message: value.message.trim(),
        };
    }

    private optionalString(value: string): string | null {
        return value.trim() || null;
    }

    private queryParam(name: string): string | null {
        const value = new URLSearchParams(window.location.search).get(name);

        return value?.trim() || null;
    }
}
