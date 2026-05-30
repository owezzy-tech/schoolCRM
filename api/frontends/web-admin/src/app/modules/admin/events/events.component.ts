import { ChangeDetectionStrategy, Component } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';

interface Event {
    id: string;
    title: string;
    datetime: string;
    location: string;
    attendeeCount: number;
}

@Component({
    selector: 'app-events',
    standalone: true,
    imports: [MatButtonModule, MatIconModule],
    changeDetection: ChangeDetectionStrategy.OnPush,
    templateUrl: './events.component.html',
})
export class EventsComponent {
    readonly events: readonly Event[] = [
        { id: '1', title: 'Fall Open House', datetime: 'Sat, Jun 6 · 10:00 AM - 2:00 PM', location: 'Main Quad', attendeeCount: 450 },
        { id: '2', title: 'Virtual Info Session — Engineering', datetime: 'Tue, Jun 9 · 6:00 PM - 7:00 PM', location: 'Online (Zoom)', attendeeCount: 125 },
        { id: '3', title: 'Counselor Meetup', datetime: 'Thu, Jun 11 · 2:00 PM - 4:00 PM', location: 'Admissions Hall', attendeeCount: 45 },
        { id: '4', title: 'Campus Tour — Science Facilities', datetime: 'Fri, Jun 12 · 1:00 PM - 3:00 PM', location: 'Science Center', attendeeCount: 30 },
        { id: '5', title: 'Financial Aid Q&A', datetime: 'Mon, Jun 15 · 5:00 PM - 6:00 PM', location: 'Online (Zoom)', attendeeCount: 85 },
        { id: '6', title: 'Accepted Students Day', datetime: 'Sat, Jun 20 · 9:00 AM - 3:00 PM', location: 'Student Union', attendeeCount: 600 },
    ];
}
