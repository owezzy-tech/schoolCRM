import { ChangeDetectionStrategy, Component, OnInit, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterLink, ActivatedRoute, Router } from '@angular/router';
import { ReactiveFormsModule, UntypedFormBuilder, UntypedFormGroup, Validators } from '@angular/forms';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { MatRadioModule } from '@angular/material/radio';
import { MatDatepickerModule } from '@angular/material/datepicker';
import { MatNativeDateModule } from '@angular/material/core';
import { MatSlideToggleModule } from '@angular/material/slide-toggle';
import {
    AdmissionsEvent,
    AdmissionsEventRequest,
    EventStatus,
    EventType,
} from 'app/core/admissions/admissions.types';
import { AdmissionsService } from 'app/core/admissions/admissions.service';
import { jsonApiErrorMessage } from 'app/core/api/json-api';
import { catchError, of } from 'rxjs';

@Component({
    selector: 'app-event-form',
    standalone: true,
    imports: [
        CommonModule,
        RouterLink,
        ReactiveFormsModule,
        MatIconModule,
        MatButtonModule,
        MatFormFieldModule,
        MatInputModule,
        MatSelectModule,
        MatRadioModule,
        MatDatepickerModule,
        MatNativeDateModule,
        MatSlideToggleModule
    ],
    changeDetection: ChangeDetectionStrategy.OnPush,
    templateUrl: './event-form.component.html'
})
export class EventFormComponent implements OnInit {
    private readonly admissionsService = inject(AdmissionsService);
    private readonly fb = inject(UntypedFormBuilder);
    private readonly route = inject(ActivatedRoute);
    private readonly router = inject(Router);

    readonly isEditMode = signal(false);
    readonly loading = signal(false);
    readonly saving = signal(false);

    eventForm!: UntypedFormGroup;
    errorMessage = '';

    ngOnInit() {
        this.eventForm = this.fb.group({
            title: ['', Validators.required],
            type: ['open-day', Validators.required],
            status: ['draft', Validators.required],
            description: ['', Validators.required],
            startDate: [new Date(), Validators.required],
            startTime: ['09:00', Validators.required],
            endDate: [new Date(), Validators.required],
            endTime: ['10:00', Validators.required],
            isVirtual: [false, Validators.required],
            location: ['', Validators.required],
            capacity: [100, [Validators.required, Validators.min(0)]],
            registrationDeadline: [null],
            autoConfirmationEnabled: [true],
            autoReminderEnabled: [true],
        });

        const id = this.route.snapshot.paramMap.get('id');
        if (id) {
            this.isEditMode.set(true);

            this.loadEvent(id);
        }
    }

    loadEvent(eventID: string): void {
        this.loading.set(true);
        this.errorMessage = '';
        this.admissionsService
            .getEvent(eventID)
            .pipe(
                catchError((error) => {
                    this.loading.set(false);
                    this.errorMessage = jsonApiErrorMessage(
                        error,
                        'Unable to load event.'
                    );

                    return of(undefined);
                })
            )
            .subscribe((event) => {
                this.loading.set(false);

                if (event) {
                    this.patchForm(event);
                }
            });
    }

    patchForm(event: AdmissionsEvent): void {
        const start = new Date(event.start);
        const end = new Date(event.end);

        this.eventForm.patchValue({
            title: event.title,
            type: event.type,
            status: event.status,
            description: event.description,
            startDate: start,
            startTime: this.formatTime(start),
            endDate: end,
            endTime: this.formatTime(end),
            isVirtual: event.isVirtual,
            location: event.location,
            capacity: event.capacity,
            registrationDeadline: event.registrationDeadline
                ? new Date(event.registrationDeadline)
                : null,
            autoConfirmationEnabled: event.autoConfirmationEnabled,
            autoReminderEnabled: event.autoReminderEnabled,
        });
    }

    save() {
        if (this.eventForm.invalid) {
            this.eventForm.markAllAsTouched();

            return;
        }

        const request = this.toRequest();
        const id = this.route.snapshot.paramMap.get('id');
        const save$ = this.isEditMode() && id
            ? this.admissionsService.updateEvent(id, request)
            : this.admissionsService.createEvent(request);

        this.saving.set(true);
        this.errorMessage = '';
        save$
            .pipe(
                catchError((error) => {
                    this.saving.set(false);
                    this.errorMessage = jsonApiErrorMessage(
                        error,
                        'Unable to save event.'
                    );

                    return of(undefined);
                })
            )
            .subscribe((event) => {
                this.saving.set(false);

                if (event) {
                    this.router.navigate(['/events', event.id]);
                }
            });
    }

    private toRequest(): AdmissionsEventRequest {
        const value = this.eventForm.getRawValue();

        return {
            title: value.title.trim(),
            type: value.type as EventType,
            status: value.status as EventStatus,
            description: value.description.trim(),
            start: this.toIsoString(value.startDate, value.startTime),
            end: this.toIsoString(value.endDate, value.endTime),
            location: value.location.trim(),
            isVirtual: Boolean(value.isVirtual),
            capacity: Number(value.capacity),
            registrationDeadline: value.registrationDeadline
                ? this.toRegistrationDeadlineIsoString(value.registrationDeadline)
                : null,
            autoConfirmationEnabled: Boolean(value.autoConfirmationEnabled),
            autoReminderEnabled: Boolean(value.autoReminderEnabled),
        };
    }

    private toIsoString(dateValue: Date, timeValue: string): string {
        const date = new Date(dateValue);
        const [hours, minutes] = (timeValue || '00:00').split(':').map(Number);

        date.setHours(hours || 0, minutes || 0, 0, 0);

        return date.toISOString();
    }

    private formatTime(date: Date): string {
        return `${String(date.getHours()).padStart(2, '0')}:${String(
            date.getMinutes()
        ).padStart(2, '0')}`;
    }

    private toRegistrationDeadlineIsoString(dateValue: Date): string {
        const deadline = new Date(dateValue);

        deadline.setHours(23, 59, 59, 999);

        return deadline.toISOString();
    }
}
