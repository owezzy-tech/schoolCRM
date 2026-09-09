import { ChangeDetectionStrategy, Component, computed, signal, inject, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterLink, ActivatedRoute } from '@angular/router';
import { ReactiveFormsModule, UntypedFormControl } from '@angular/forms';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { MatInputModule } from '@angular/material/input';
import { MatFormFieldModule } from '@angular/material/form-field';
import { AdmissionsService } from 'app/core/admissions/admissions.service';
import { jsonApiErrorMessage } from 'app/core/api/json-api';
import { catchError, of } from 'rxjs';
import { EventItem, EventRegistration } from '../models/event.types';

@Component({
    selector: 'app-event-checkin',
    standalone: true,
    imports: [
        CommonModule,
        RouterLink,
        ReactiveFormsModule,
        MatIconModule,
        MatButtonModule,
        MatInputModule,
        MatFormFieldModule
    ],
    changeDetection: ChangeDetectionStrategy.OnPush,
    templateUrl: './event-checkin.component.html'
})
export class EventCheckinComponent implements OnInit {
    private readonly admissionsService = inject(AdmissionsService);
    private readonly route = inject(ActivatedRoute);

    readonly event = signal<EventItem | null>(null);
    readonly loading = signal(false);
    readonly searchControl = new UntypedFormControl('');
    readonly registrants = computed(
        () => this.event()?.registrations ?? []
    );

    errorMessage = '';
    savingRegistrationID = '';

    ngOnInit() {
        const id = this.route.snapshot.paramMap.get('id');
        if (!id) {
            this.errorMessage = 'Event ID is missing.';

            return;
        }

        this.loadEvent(id);
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
                        'Unable to load check-in attendees.'
                    );

                    return of(undefined);
                })
            )
            .subscribe((event) => {
                this.loading.set(false);
                if (event) {
                    this.event.set(event);
                }
            });
    }

    get filteredRegistrants() {
        const term = (this.searchControl.value || '').toLowerCase();
        if (!term) return this.registrants();
        return this.registrants().filter((registration) =>
            [
                registration.constituentName,
                registration.email,
                registration.id,
                registration.constituentId ?? '',
            ]
                .join(' ')
                .toLowerCase()
                .includes(term)
        );
    }

    get checkedInCount() {
        return this.registrants().filter(r => r.status === 'checked-in').length;
    }

    toggleCheckin(id: string) {
        const event = this.event();
        const registration = this.registrants().find((item) => item.id === id);

        if (!event || !registration || registration.status === 'checked-in') {
            return;
        }

        this.savingRegistrationID = id;
        this.errorMessage = '';
        this.admissionsService
            .checkInEventRegistration(id, {
                registrationId: id,
            })
            .pipe(
                catchError((error) => {
                    this.savingRegistrationID = '';
                    this.errorMessage = jsonApiErrorMessage(
                        error,
                        'Unable to check in attendee.'
                    );

                    return of(undefined);
                })
            )
            .subscribe((updatedRegistration) => {
                this.savingRegistrationID = '';

                if (!updatedRegistration) {
                    return;
                }

                this.event.update((current) => {
                    if (!current) {
                        return current;
                    }

                    const registrations = current.registrations.map((item) =>
                        item.id === updatedRegistration.id
                            ? updatedRegistration
                            : item
                    );

                    return {
                        ...current,
                        registrations,
                        checkedInCount: registrations.filter(
                            (item) => item.status === 'checked-in'
                        ).length,
                    };
                });
            });
    }
}
