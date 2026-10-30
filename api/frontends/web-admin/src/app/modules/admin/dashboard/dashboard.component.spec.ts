import { DATE_PIPE_DEFAULT_OPTIONS } from '@angular/common';
import { Component, provideZonelessChangeDetection } from '@angular/core';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { MatIconModule } from '@angular/material/icon';
import { provideRouter } from '@angular/router';
import { AdmissionsService } from 'app/core/admissions/admissions.service';
import { AdmissionsEvent } from 'app/core/admissions/admissions.types';
import { UserService } from 'app/core/user/user.service';
import { User } from 'app/core/user/user.types';
import { of, ReplaySubject, throwError } from 'rxjs';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { DashboardComponent } from './dashboard.component';

@Component({ selector: 'mat-icon', template: '', inputs: ['svgIcon'] })
class IconStubComponent {}

describe('Dashboard actual user and events', () => {
    let fixture: ComponentFixture<DashboardComponent>;
    let user: ReplaySubject<User>;
    let queryEvents: ReturnType<typeof vi.fn>;

    beforeEach(async () => {
        user = new ReplaySubject<User>(1);
        queryEvents = vi.fn(() => of({ items: [], total: 0 }));
        await TestBed.configureTestingModule({
            imports: [DashboardComponent],
            providers: [
                provideZonelessChangeDetection(),
                provideRouter([]),
                {
                    provide: DATE_PIPE_DEFAULT_OPTIONS,
                    useValue: { timezone: 'UTC' },
                },
                {
                    provide: UserService,
                    useValue: { user$: user.asObservable() },
                },
                {
                    provide: AdmissionsService,
                    useValue: {
                        queryApplications: () => of({ items: [], total: 0 }),
                        queryEvents,
                    },
                },
            ],
        })
            .overrideComponent(DashboardComponent, {
                remove: { imports: [MatIconModule] },
                add: { imports: [IconStubComponent] },
            })
            .compileComponents();
    });

    async function render(): Promise<HTMLElement> {
        fixture = TestBed.createComponent(DashboardComponent);
        await fixture.whenStable();
        return fixture.nativeElement;
    }

    it('renders and updates the signed-in name without a template identity', async () => {
        user.next({
            id: 'operator',
            name: 'Owen Adirah',
            email: 'operator@schoolcrm.invalid',
        });
        const page = await render();
        expect(page.querySelector('h1')?.textContent).toContain(
            'Welcome back, Owen Adirah'
        );
        expect(page.textContent).not.toContain('Avery');
        user.next({
            id: 'operator',
            name: 'Updated operator',
            email: 'operator@schoolcrm.invalid',
        });
        await fixture.whenStable();
        expect(page.querySelector('h1')?.textContent).toContain(
            'Updated operator'
        );
    });

    it('shows an honest empty state when no upcoming events exist', async () => {
        const page = await render();
        expect(page.querySelector('h1')?.textContent?.trim()).toBe(
            'Welcome back'
        );
        expect(page.textContent).toContain('No upcoming events.');
        expect(page.textContent).not.toContain('Fall Open House');
        expect(page.textContent).not.toContain('Virtual Info Session');
        expect(queryEvents).toHaveBeenCalledWith({
            rows: 3,
            status: 'upcoming',
            orderBy: 'start_time,ASC',
        });
    });

    it('renders the actual event title, date and location from the API', async () => {
        const event: AdmissionsEvent = {
            id: 'event-real',
            title: 'Nairobi admissions briefing',
            type: 'open-day',
            status: 'upcoming',
            description: '',
            start: '2026-11-20T09:00:00Z',
            end: '2026-11-20T10:00:00Z',
            location: 'Languages hall',
            isVirtual: false,
            capacity: 30,
            registeredCount: 0,
            checkedInCount: 0,
            registrations: [],
            autoConfirmationEnabled: false,
            autoReminderEnabled: false,
            dateCreated: '2026-10-09T00:00:00Z',
            dateUpdated: '2026-10-09T00:00:00Z',
        };
        queryEvents.mockReturnValue(of({ items: [event], total: 1 }));
        const page = await render();
        const item = page.querySelector('li');
        expect(item?.textContent).toContain(event.title);
        expect(item?.textContent).toContain('Nov 20');
        expect(item?.textContent).toContain(event.location);
        expect(page.textContent).not.toContain('No upcoming events.');
    });

    it('shows a load failure rather than an empty result and recovers on retry', async () => {
        queryEvents.mockReturnValue(throwError(() => new Error('Unavailable')));
        const page = await render();
        expect(page.textContent).toContain('Unable to load dashboard data.');
        expect(page.textContent).not.toContain('No upcoming events.');
        queryEvents.mockReturnValue(of({ items: [], total: 0 }));
        page.querySelector<HTMLButtonElement>('button')?.click();
        await fixture.whenStable();
        expect(page.textContent).toContain('No upcoming events.');
        expect(page.textContent).not.toContain(
            'Unable to load dashboard data.'
        );
    });
});
