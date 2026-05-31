import { ChangeDetectionStrategy, Component, ViewChild, signal, inject, ViewEncapsulation } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';
import { FullCalendarModule } from '@fullcalendar/angular';
import { CalendarOptions, EventClickArg } from '@fullcalendar/core';
import dayGridPlugin from '@fullcalendar/daygrid';
import timeGridPlugin from '@fullcalendar/timegrid';
import listPlugin from '@fullcalendar/list';
import interactionPlugin from '@fullcalendar/interaction';
import { MOCK_EVENTS } from '../../data/events.mock';

@Component({
    selector: 'app-events-calendar',
    standalone: true,
    imports: [CommonModule, FullCalendarModule],
    changeDetection: ChangeDetectionStrategy.OnPush,
    encapsulation: ViewEncapsulation.None,
    templateUrl: './events-calendar.component.html',
})
export class EventsCalendarComponent {
    private readonly router = inject(Router);

    readonly calendarOptions = signal<CalendarOptions>({
        plugins: [dayGridPlugin, timeGridPlugin, listPlugin, interactionPlugin],
        initialView: 'dayGridMonth',
        initialDate: '2026-06-01', // Lock to June 2026 for mock data
        headerToolbar: {
            left: 'prev,next today',
            center: 'title',
            right: 'dayGridMonth,timeGridWeek,listWeek'
        },
        events: MOCK_EVENTS.map(e => ({
            id: e.id,
            title: e.title,
            start: e.start,
            end: e.end,
            extendedProps: {
                type: e.type,
                status: e.status
            },
            className: [`event-type-${e.type}`]
        })),
        eventClick: this.handleEventClick.bind(this),
        height: 800,
        editable: true,
        eventDisplay: 'block'
    });

    handleEventClick(clickInfo: EventClickArg) {
        this.router.navigate(['/events', clickInfo.event.id]);
    }
}
