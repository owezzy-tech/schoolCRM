import { ChangeDetectionStrategy, Component } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';

interface EventCard {
    month: string;
    day: string;
    type: string;
    title: string;
    dateTime: string;
    location: string;
    spots: string;
}

@Component({
    selector: 'app-portal-events',
    standalone: true,
    imports: [
        MatButtonModule,
        MatIconModule
    ],
    template: `
        <div class="flex flex-col flex-auto py-12 px-6 sm:px-10 lg:px-16">
            <div class="mx-auto w-full max-w-6xl">
                <div class="flex flex-col text-center sm:text-left">
                    <h1 class="text-4xl font-bold tracking-tight text-default">Visit &amp; explore</h1>
                    <p class="mt-2 text-lg text-secondary">Sign up for upcoming events to learn more about Northbrook.</p>
                </div>

                <div class="mt-10 grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
                    @for (event of events; track event.title) {
                        <div class="flex flex-col overflow-hidden rounded-2xl border bg-card shadow-sm">
                            <!-- TOP section -->
                            <div class="bg-primary-50 p-6 dark:bg-primary-900/20 flex flex-col items-start">
                                <div class="w-fit rounded-lg bg-card px-3 py-2 text-center shadow-sm">
                                    <div class="text-xs font-semibold uppercase text-primary">{{event.month}}</div>
                                    <div class="text-2xl font-bold leading-none text-default">{{event.day}}</div>
                                </div>
                                <div class="mt-4 inline-block rounded-full bg-card px-3 py-1 text-xs font-medium text-primary shadow-sm">
                                    {{event.type}}
                                </div>
                                <h3 class="mt-3 text-xl font-semibold text-default">{{event.title}}</h3>
                            </div>
                            
                            <!-- BOTTOM section -->
                            <div class="flex flex-col gap-3 p-6 flex-auto">
                                <div class="flex items-center gap-3 text-secondary">
                                    <mat-icon svgIcon="heroicons_outline:calendar-days" class="icon-size-5 shrink-0"></mat-icon>
                                    <span class="text-sm font-medium">{{event.dateTime}}</span>
                                </div>
                                <div class="flex items-center gap-3 text-secondary">
                                    <mat-icon svgIcon="heroicons_outline:map-pin" class="icon-size-5 shrink-0"></mat-icon>
                                    <span class="text-sm font-medium">{{event.location}}</span>
                                </div>
                                <div class="flex items-center gap-3 text-secondary">
                                    <mat-icon svgIcon="heroicons_outline:users" class="icon-size-5 shrink-0"></mat-icon>
                                    <span class="text-sm font-medium">{{event.spots}}</span>
                                </div>
                                
                                <div class="mt-auto pt-4">
                                    <button mat-flat-button color="primary" class="w-full">Reserve a spot</button>
                                </div>
                            </div>
                        </div>
                    }
                </div>
            </div>
        </div>
    `,
    changeDetection: ChangeDetectionStrategy.OnPush,
    host: { class: 'flex w-full flex-auto flex-col' },
})
export class PortalEventsComponent {
    readonly events: ReadonlyArray<EventCard> = [
        { month: 'JUN', day: '8', type: 'Open Day', title: 'Spring Open Day', dateTime: '6/8/2026, 10:00 AM', location: 'Main Campus', spots: '88 spots left' },
        { month: 'JUN', day: '12', type: 'Webinar', title: 'Engineering Webinar', dateTime: '6/12/2026, 5:00 PM', location: 'Online', spots: '358 spots left' },
        { month: 'JUN', day: '14', type: 'Info Session', title: 'MBA Info Session', dateTime: '6/14/2026, 2:00 PM', location: 'Online', spots: '52 spots left' },
        { month: 'JUN', day: '21', type: 'Campus Tour', title: 'Architecture Studio Tour', dateTime: '6/21/2026, 11:00 AM', location: 'School of Design', spots: '19 spots left' },
        { month: 'MAY', day: '30', type: 'Webinar', title: 'Scholarship Q&A', dateTime: '5/30/2026, 4:00 PM', location: 'Online', spots: '213 spots left' },
    ];
}
