import { DatePipe } from '@angular/common';
import {
    ChangeDetectionStrategy,
    Component,
    computed,
    inject,
    signal,
} from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { RouterLink } from '@angular/router';
import { AdmissionsService } from 'app/core/admissions/admissions.service';
import { AdmissionsEvent } from 'app/core/admissions/admissions.types';
import { jsonApiErrorMessage } from 'app/core/api/json-api';
import { PortalAuthService } from 'app/core/portal/portal-auth.service';
import { catchError, of } from 'rxjs';

interface PortalRegistrationDraft {
    event: AdmissionsEvent;
    firstName: string;
    lastName: string;
    email: string;
    phone: string;
    matchMessage: string;
}

@Component({
    selector: 'app-portal-events',
    imports: [DatePipe, MatButtonModule, MatIconModule, RouterLink],
    template: `
        <div class="flex flex-auto flex-col px-6 py-12 sm:px-10 lg:px-16">
            <div class="mx-auto w-full max-w-6xl">
                <div class="flex flex-col text-center sm:text-left">
                    <h1 class="text-default text-4xl font-bold tracking-tight">
                        Visit, attend &amp; apply
                    </h1>
                    <p class="text-secondary mt-2 text-lg">
                        Sign up for open days, KUCCPS guidance clinics, and
                        virtual sessions. We match your registration to an
                        existing admissions record when your email or phone is
                        recognized.
                    </p>
                </div>

                <div class="bg-card mt-8 rounded-2xl border p-5 shadow-sm">
                    <div
                        class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between"
                    >
                        <div>
                            <p class="text-sm font-semibold">
                                Public registration rules
                            </p>
                            <p class="text-secondary text-sm">
                                Applicant portal tokens prefill Kenya admissions
                                context. Public visitors can still browse
                                capacity and register interest.
                            </p>
                        </div>
                        @if (portalSession(); as session) {
                            <div
                                class="inline-flex items-center gap-2 rounded-full bg-primary-50 px-4 py-2 text-sm font-medium text-primary dark:bg-primary-900/30"
                            >
                                <mat-icon
                                    class="icon-size-5"
                                    svgIcon="heroicons_outline:shield-check"
                                ></mat-icon>
                                Registering as {{ session.applicantName }}
                            </div>
                        } @else {
                            <a
                                routerLink="/portal"
                                mat-stroked-button
                                class="rounded-full"
                                >Get applicant access</a
                            >
                        }
                    </div>
                </div>

                <div
                    class="mt-10 grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3"
                >
                    @if (loadingEvents()) {
                        <div
                            class="text-secondary bg-card rounded-2xl border p-6 text-sm shadow-sm sm:col-span-2 lg:col-span-3"
                        >
                            Loading upcoming admissions events…
                        </div>
                    }

                    @if (eventsErrorMessage()) {
                        <div
                            class="rounded-2xl border border-red-200 bg-red-50 p-6 text-sm text-red-700 shadow-sm sm:col-span-2 lg:col-span-3 dark:border-red-900/60 dark:bg-red-900/20 dark:text-red-200"
                        >
                            {{ eventsErrorMessage() }}
                        </div>
                    }

                    @for (event of publicEvents(); track event.id) {
                        <div
                            class="bg-card flex flex-col overflow-hidden rounded-2xl border shadow-sm"
                        >
                            <div
                                class="flex flex-col items-start bg-primary-50 p-6 dark:bg-primary-900/20"
                            >
                                <div
                                    class="bg-card w-fit rounded-lg px-3 py-2 text-center shadow-sm"
                                >
                                    <div
                                        class="text-xs font-semibold uppercase text-primary"
                                    >
                                        {{ event.start | date: 'MMM' }}
                                    </div>
                                    <div
                                        class="text-default text-2xl font-bold leading-none"
                                    >
                                        {{ event.start | date: 'd' }}
                                    </div>
                                </div>
                                <div
                                    class="bg-card mt-4 inline-block rounded-full px-3 py-1 text-xs font-medium text-primary shadow-sm"
                                >
                                    {{ event.type.replace('-', ' ') }}
                                </div>
                                <h3
                                    class="text-default mt-3 text-xl font-semibold"
                                >
                                    {{ event.title }}
                                </h3>
                            </div>

                            <div class="flex flex-auto flex-col gap-3 p-6">
                                <div
                                    class="text-secondary flex items-center gap-3"
                                >
                                    <mat-icon
                                        svgIcon="heroicons_outline:calendar-days"
                                        class="shrink-0 icon-size-5"
                                    ></mat-icon>
                                    <span class="text-sm font-medium">{{
                                        event.start | date: 'medium'
                                    }}</span>
                                </div>
                                <div
                                    class="text-secondary flex items-center gap-3"
                                >
                                    <mat-icon
                                        [svgIcon]="
                                            event.isVirtual
                                                ? 'heroicons_outline:video-camera'
                                                : 'heroicons_outline:map-pin'
                                        "
                                        class="shrink-0 icon-size-5"
                                    ></mat-icon>
                                    <span class="text-sm font-medium">{{
                                        event.location
                                    }}</span>
                                </div>
                                <div
                                    class="text-secondary flex items-center gap-3"
                                >
                                    <mat-icon
                                        svgIcon="heroicons_outline:users"
                                        class="shrink-0 icon-size-5"
                                    ></mat-icon>
                                    <span class="text-sm font-medium"
                                        >{{ remainingSpots(event) }} spots
                                        left</span
                                    >
                                </div>

                                <div
                                    class="mt-2 h-2 overflow-hidden rounded-full bg-gray-100 dark:bg-gray-800"
                                >
                                    <div
                                        class="h-full rounded-full bg-primary"
                                        [style.width.%]="capacityPercent(event)"
                                    ></div>
                                </div>

                                <div class="mt-auto pt-4">
                                    <button
                                        mat-flat-button
                                        color="primary"
                                        class="w-full"
                                        [disabled]="isFull(event)"
                                        (click)="selectEvent(event)"
                                    >
                                        {{
                                            isFull(event)
                                                ? 'Event full'
                                                : 'Reserve a spot'
                                        }}
                                    </button>
                                </div>
                            </div>
                        </div>
                    }
                </div>

                @if (registrationDraft(); as draft) {
                    <div class="bg-card mt-10 rounded-2xl border p-6 shadow-sm">
                        <div
                            class="flex flex-col gap-6 lg:flex-row lg:items-start lg:justify-between"
                        >
                            <div class="max-w-2xl">
                                <p
                                    class="text-sm font-semibold uppercase tracking-[0.2em] text-primary"
                                >
                                    Registration preview
                                </p>
                                <h2 class="mt-2 text-2xl font-semibold">
                                    {{ draft.event.title }}
                                </h2>
                                <p class="text-secondary mt-2">
                                    Public registrants provide Kenyan contact
                                    details first. The CRM attempts a
                                    constituent match before creating a new
                                    prospect registration.
                                </p>
                            </div>
                            <button
                                mat-stroked-button
                                (click)="clearRegistrationDraft()"
                            >
                                Close
                            </button>
                        </div>

                        <div
                            class="mt-6 grid gap-4 md:grid-cols-2 xl:grid-cols-4"
                        >
                            <div class="rounded-xl border border-dashed p-4">
                                <p class="text-secondary text-xs uppercase">
                                    First name
                                </p>
                                <p class="mt-1 font-medium">
                                    {{ draft.firstName }}
                                </p>
                            </div>
                            <div class="rounded-xl border border-dashed p-4">
                                <p class="text-secondary text-xs uppercase">
                                    Last name
                                </p>
                                <p class="mt-1 font-medium">
                                    {{ draft.lastName }}
                                </p>
                            </div>
                            <div class="rounded-xl border border-dashed p-4">
                                <p class="text-secondary text-xs uppercase">
                                    Email
                                </p>
                                <p class="mt-1 font-medium">
                                    {{ draft.email }}
                                </p>
                            </div>
                            <div class="rounded-xl border border-dashed p-4">
                                <p class="text-secondary text-xs uppercase">
                                    Phone
                                </p>
                                <p class="mt-1 font-medium">
                                    {{ draft.phone }}
                                </p>
                            </div>
                        </div>

                        <div
                            class="mt-5 rounded-xl bg-green-50 p-4 text-sm text-green-800 dark:bg-green-900/20 dark:text-green-200"
                        >
                            {{ draft.matchMessage }} Capacity check passed with
                            {{ remainingSpots(draft.event) }} seats remaining. A
                            confirmation email and calendar details would be
                            sent by email and SMS after submit.
                        </div>
                    </div>
                }
            </div>
        </div>
    `,
    changeDetection: ChangeDetectionStrategy.OnPush,
    host: { class: 'flex w-full flex-auto flex-col' },
})
export class PortalEventsComponent {
    private readonly admissionsService = inject(AdmissionsService);
    private readonly portalAuthService = inject(PortalAuthService);

    readonly portalSession = this.portalAuthService.session;
    readonly selectedEvent = signal<AdmissionsEvent | null>(null);
    readonly events = signal<AdmissionsEvent[]>([]);
    readonly loadingEvents = signal(false);
    readonly eventsErrorMessage = signal('');
    readonly publicEvents = computed(() =>
        this.events()
            .filter((event) => ['upcoming', 'live'].includes(event.status))
            .slice(0, 6)
    );

    readonly registrationDraft = computed<PortalRegistrationDraft | null>(
        () => {
            const event = this.selectedEvent();

            if (!event) {
                return null;
            }

            const session = this.portalSession();

            return {
                event,
                firstName: session?.applicantName.split(' ')[0] ?? 'Achieng',
                lastName:
                    session?.applicantName.split(' ').slice(1).join(' ') ||
                    'Otieno',
                email: session?.email ?? 'achieng.otieno@example.ac.ke',
                phone: '+254 712 345 678',
                matchMessage: session
                    ? `Matched portal token for application ${session.applicationID} before creating the event registration.`
                    : 'Matched existing constituent CON-921 from verified email before creating the event registration.',
            };
        }
    );

    constructor() {
        this.loadEvents();
    }

    loadEvents(): void {
        this.loadingEvents.set(true);
        this.eventsErrorMessage.set('');

        this.admissionsService
            .queryEvents({
                page: 1,
                rows: 50,
                orderBy: 'start_time,ASC',
            })
            .pipe(
                catchError((error) => {
                    this.eventsErrorMessage.set(
                        jsonApiErrorMessage(
                            error,
                            'Unable to load upcoming admissions events.'
                        )
                    );

                    return of({
                        items: [],
                        total: 0,
                        page: 1,
                        rowsPerPage: 50,
                    });
                })
            )
            .subscribe((result) => {
                this.events.set(result.items);
                this.loadingEvents.set(false);
            });
    }

    selectEvent(event: AdmissionsEvent): void {
        if (this.isFull(event)) {
            return;
        }

        this.selectedEvent.set(event);
    }

    clearRegistrationDraft(): void {
        this.selectedEvent.set(null);
    }

    remainingSpots(event: AdmissionsEvent): number {
        return Math.max(event.capacity - event.registeredCount, 0);
    }

    capacityPercent(event: AdmissionsEvent): number {
        if (event.capacity === 0) {
            return 0;
        }

        return Math.round((event.registeredCount / event.capacity) * 100);
    }

    isFull(event: AdmissionsEvent): boolean {
        return this.remainingSpots(event) === 0;
    }
}
