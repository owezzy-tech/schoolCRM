import { ChangeDetectionStrategy, Component, signal } from '@angular/core';
import { FullCalendarModule } from '@fullcalendar/angular';
import { CalendarOptions, EventClickArg } from '@fullcalendar/core';
import dayGridPlugin from '@fullcalendar/daygrid';
import interactionPlugin, { DateClickArg } from '@fullcalendar/interaction';
import listPlugin from '@fullcalendar/list';
import timeGridPlugin from '@fullcalendar/timegrid';

@Component({
    selector: 'app-calendar',
    imports: [FullCalendarModule],
    templateUrl: './calendar.component.html',
    styleUrl: './calendar.component.scss',
    changeDetection: ChangeDetectionStrategy.OnPush,
})
export class CalendarComponent {
    readonly selectedDate = signal('Select a date or event to inspect it.');

    readonly calendarOptions: CalendarOptions = {
        plugins: [
            dayGridPlugin,
            timeGridPlugin,
            listPlugin,
            interactionPlugin,
        ],
        initialDate: '2020-06-04',
        initialView: 'dayGridMonth',
        headerToolbar: {
            left: 'prev,next today',
            center: 'title',
            right: 'dayGridMonth,timeGridWeek,timeGridDay,listWeek',
        },
        buttonText: {
            today: 'Today',
            month: 'Month',
            week: 'Week',
            day: 'Day',
            list: 'List',
        },
        weekends: true,
        selectable: true,
        editable: true,
        dayMaxEvents: true,
        events: '/assets/data/calendar.json',
        dateClick: (arg: DateClickArg) => {
            this.selectedDate.set(`Selected date: ${arg.dateStr}`);
        },
        eventClick: (arg: EventClickArg) => {
            const details = arg.event.extendedProps['details'];
            this.selectedDate.set(
                `${arg.event.title}: ${typeof details === 'string' ? details : 'No details provided.'}`
            );
        },
    };
}
