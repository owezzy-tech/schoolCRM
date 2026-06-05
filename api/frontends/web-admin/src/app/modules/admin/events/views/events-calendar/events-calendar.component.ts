import {
    ChangeDetectionStrategy,
    Component,
    computed,
    inject,
    ViewEncapsulation,
} from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';
import { FullCalendarModule } from '@fullcalendar/angular';
import { CalendarOptions, EventClickArg } from '@fullcalendar/core';
import dayGridPlugin from '@fullcalendar/daygrid';
import timeGridPlugin from '@fullcalendar/timegrid';
import listPlugin from '@fullcalendar/list';
import interactionPlugin from '@fullcalendar/interaction';
import { EventsComponent } from '../../events.component';

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
    private readonly shell = inject(EventsComponent);

    readonly calendarOptions = computed<CalendarOptions>(() => ({
        plugins: [dayGridPlugin, timeGridPlugin, listPlugin, interactionPlugin],
        initialView: 'dayGridMonth',
        headerToolbar: {
            left: 'prev,next today',
            center: 'title',
            right: 'dayGridMonth,timeGridWeek,listWeek'
        },
        events: this.shell.events().map((event) => ({
            id: event.id,
            title: event.title,
            start: event.start,
            end: event.end,
            extendedProps: {
                type: event.type,
                status: event.status,
            },
            className: [`event-type-${event.type}`],
        })),
        eventClick: this.handleEventClick.bind(this),
        height: 800,
        editable: false,
        eventDisplay: 'block'
    }));

    handleEventClick(clickInfo: EventClickArg) {
        this.router.navigate(['/events', clickInfo.event.id]);
    }
}
