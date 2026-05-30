import { ChangeDetectionStrategy, Component, inject, signal } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { RouterLink } from '@angular/router';
import { FuseAlertComponent, FuseAlertType } from '@fuse/components/alert';
import { AdmissionsService } from 'app/core/admissions/admissions.service';
import { InquiryRequest } from 'app/core/admissions/admissions.types';
import { jsonApiErrorMessage } from 'app/core/api/json-api';
import { finalize } from 'rxjs';

@Component({
    selector: 'app-inquiry',
    imports: [
        RouterLink,
        ReactiveFormsModule,
        FuseAlertComponent,
        MatButtonModule,
        MatFormFieldModule,
        MatIconModule,
        MatInputModule,
        MatProgressSpinnerModule,
    ],
    templateUrl: './inquiry.component.html',
    changeDetection: ChangeDetectionStrategy.OnPush,
})
export class InquiryComponent {
    private readonly admissionsService = inject(AdmissionsService);
    private readonly formBuilder = inject(FormBuilder);

    readonly alert = signal<{ type: FuseAlertType; message: string } | null>(
        null
    );
    readonly submitting = signal(false);
    readonly submitted = signal(false);
    readonly inquiryForm = this.formBuilder.nonNullable.group({
        firstName: ['', Validators.required],
        lastName: ['', Validators.required],
        dateOfBirth: ['', Validators.required],
        primaryEmail: ['', [Validators.required, Validators.email]],
        primaryPhone: ['', Validators.required],
        message: [''],
    });

    submitInquiry(): void {
        if (this.inquiryForm.invalid) {
            this.inquiryForm.markAllAsTouched();
            return;
        }

        this.submitting.set(true);
        this.alert.set(null);

        this.admissionsService
            .createInquiry(this.inquiryRequest())
            .pipe(finalize(() => this.submitting.set(false)))
            .subscribe({
                next: () => {
                    this.submitted.set(true);
                    this.inquiryForm.reset();
                    this.alert.set({
                        type: 'success',
                        message:
                            'Your inquiry has been received. Our admissions team will follow up shortly.',
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
        const message = value.message.trim();

        return {
            firstName: value.firstName.trim(),
            lastName: value.lastName.trim(),
            dateOfBirth: new Date(value.dateOfBirth).toISOString(),
            primaryEmail: value.primaryEmail.trim(),
            primaryPhone: value.primaryPhone.trim(),
            programOfInterest: null,
            termOfInterest: null,
            source: 'PUBLIC_INQUIRY_FORM',
            utmSource: this.queryParam('utm_source'),
            utmMedium: this.queryParam('utm_medium'),
            utmCampaign: this.queryParam('utm_campaign'),
            message: message || null,
        };
    }

    private queryParam(name: string): string | null {
        const value = new URLSearchParams(window.location.search).get(name);

        return value?.trim() || null;
    }
}
