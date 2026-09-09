import { Component } from '@angular/core';
import { TestBed } from '@angular/core/testing';
import {
    ActivatedRoute,
    NavigationEnd,
    NavigationExtras,
    Router,
    Event as RouterEvent,
    RouterOutlet,
} from '@angular/router';
import { AdmissionsService } from 'app/core/admissions/admissions.service';
import { AdmissionsEvent } from 'app/core/admissions/admissions.types';
import { PaginatedResult } from 'app/core/api/json-api';
import { of, ReplaySubject, Subject, throwError } from 'rxjs';
import { beforeEach, describe, expect, it, vi } from 'vitest';

import { EventsComponent } from './events.component';

@Component({
    selector: 'router-outlet',
    template: '',
})
class RouterOutletStubComponent {}

describe('EventsComponent', () => {
    let eventsSubject: ReplaySubject<PaginatedResult<AdmissionsEvent>>;
    let routerEvents: Subject<RouterEvent>;
    let queryEvents: ReturnType<typeof vi.fn>;
    let routerNavigate: Router['navigate'];
    let routerNavigateCalls: Array<{
        commands: readonly unknown[];
        extras?: NavigationExtras;
    }>;
    let routerStub: Pick<Router, 'events' | 'navigate' | 'url'>;
    let routerUrl: string;

    beforeEach(async () => {
        eventsSubject = new ReplaySubject<PaginatedResult<AdmissionsEvent>>(1);
        routerEvents = new Subject<RouterEvent>();
        queryEvents = vi.fn(() =>
            of(eventsResult([openDayEvent, webinarEvent]))
        );
        routerNavigateCalls = [];
        routerNavigate = (commands, extras) => {
            routerNavigateCalls.push({ commands, extras });

            return Promise.resolve(true);
        };
        routerUrl = '/admin/events/table';
        routerStub = {
            events: routerEvents.asObservable(),
            navigate: routerNavigate,
            get url() {
                return routerUrl;
            },
        };

        await TestBed.configureTestingModule({
            imports: [EventsComponent],
            providers: [
                {
                    provide: AdmissionsService,
                    useValue: {
                        events$: eventsSubject.asObservable(),
                        queryEvents,
                    },
                },
                {
                    provide: Router,
                    useValue: routerStub,
                },
                {
                    provide: ActivatedRoute,
                    useValue: {},
                },
            ],
        })
            .overrideComponent(EventsComponent, {
                remove: {
                    imports: [RouterOutlet],
                },
                add: {
                    imports: [RouterOutletStubComponent],
                    template: '',
                },
            })
            .compileComponents();
    });

    it('loads events from the admissions API using the live event query contract', () => {
        TestBed.createComponent(EventsComponent);

        expect(queryEvents).toHaveBeenCalledWith({
            page: 1,
            rows: 50,
            orderBy: 'start_time,ASC',
        });
    });

    it('derives filtered events from the live service stream', () => {
        const fixture = TestBed.createComponent(EventsComponent);
        const component = fixture.componentInstance;

        eventsSubject.next(eventsResult([openDayEvent, webinarEvent]));
        component.toggleTypeFilter('webinar');

        expect(component.events()).toEqual([webinarEvent]);
    });

    it('surfaces JSON:API errors when event loading fails', () => {
        queryEvents.mockReturnValue(
            throwError(() => ({
                error: {
                    errors: [
                        {
                            detail: 'Events service unavailable.',
                        },
                    ],
                },
            }))
        );

        const fixture = TestBed.createComponent(EventsComponent);
        const component = fixture.componentInstance;

        expect(component.loading()).toBe(false);
        expect(component.errorMessage).toBe('Events service unavailable.');
    });

    it('tracks the active view from router navigation and route buttons', () => {
        const fixture = TestBed.createComponent(EventsComponent);
        const component = fixture.componentInstance;

        routerUrl = '/admin/events/calendar';
        routerEvents.next(new NavigationEnd(1, routerUrl, routerUrl));

        expect(component.currentView()).toBe('calendar');

        component.onViewChange('kanban');

        expect(component.currentView()).toBe('kanban');
        expect(routerNavigateCalls).toEqual([
            {
                commands: ['kanban'],
                extras: {
                    relativeTo: TestBed.inject(ActivatedRoute),
                },
            },
        ]);
    });
});

function eventsResult(
    items: AdmissionsEvent[]
): PaginatedResult<AdmissionsEvent> {
    return {
        items,
        total: items.length,
        page: 1,
        rowsPerPage: 50,
    };
}

const openDayEvent: AdmissionsEvent = {
    id: 'event-open-day',
    title: 'Nairobi Open Day',
    type: 'open-day',
    status: 'upcoming',
    description: 'Meet admissions advisors on campus.',
    start: '2026-07-15T09:00:00Z',
    end: '2026-07-15T12:00:00Z',
    location: 'Main campus',
    isVirtual: false,
    capacity: 120,
    registeredCount: 24,
    checkedInCount: 0,
    autoConfirmationEnabled: true,
    autoReminderEnabled: true,
    registrations: [],
    dateCreated: '2026-06-01T00:00:00Z',
    dateUpdated: '2026-06-01T00:00:00Z',
};

const webinarEvent: AdmissionsEvent = {
    id: 'event-webinar',
    title: 'KUCCPS Guidance Webinar',
    type: 'webinar',
    status: 'live',
    description: 'Online admissions guidance.',
    start: '2026-07-20T09:00:00Z',
    end: '2026-07-20T10:00:00Z',
    location: 'Virtual',
    isVirtual: true,
    capacity: 200,
    registeredCount: 80,
    checkedInCount: 12,
    autoConfirmationEnabled: true,
    autoReminderEnabled: false,
    registrations: [],
    dateCreated: '2026-06-01T00:00:00Z',
    dateUpdated: '2026-06-01T00:00:00Z',
};
