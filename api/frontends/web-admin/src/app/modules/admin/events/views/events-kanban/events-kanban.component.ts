import { ChangeDetectionStrategy, Component, signal } from '@angular/core';
import { CommonModule, DatePipe } from '@angular/common';
import { RouterLink } from '@angular/router';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { DragDropModule, CdkDragDrop, moveItemInArray, transferArrayItem } from '@angular/cdk/drag-drop';
import { MOCK_EVENTS } from '../../data/events.mock';
import { EventItem } from '../../models/event.types';

@Component({
    selector: 'app-events-kanban',
    standalone: true,
    imports: [
        CommonModule,
        RouterLink,
        MatIconModule,
        MatButtonModule,
        MatProgressBarModule,
        DragDropModule,
        DatePipe
    ],
    changeDetection: ChangeDetectionStrategy.OnPush,
    templateUrl: './events-kanban.component.html',
    styleUrls: ['./events-kanban.component.scss']
})
export class EventsKanbanComponent {
    
    readonly columns = signal([
        { id: 'draft', title: 'Draft', events: MOCK_EVENTS.filter(e => e.status === 'draft') },
        { id: 'upcoming', title: 'Upcoming', events: MOCK_EVENTS.filter(e => e.status === 'upcoming') },
        { id: 'live', title: 'Live', events: MOCK_EVENTS.filter(e => e.status === 'live') },
        { id: 'completed', title: 'Completed', events: MOCK_EVENTS.filter(e => e.status === 'completed') }
    ]);

    drop(event: CdkDragDrop<EventItem[]>) {
        if (event.previousContainer === event.container) {
            moveItemInArray(event.container.data, event.previousIndex, event.currentIndex);
        } else {
            transferArrayItem(
                event.previousContainer.data,
                event.container.data,
                event.previousIndex,
                event.currentIndex,
            );
        }
    }
    
    getTypeColor(type: string): string {
        switch (type) {
            case 'open-day': return 'bg-purple-500';
            case 'webinar': return 'bg-teal-500';
            case 'info-session': return 'bg-green-500';
            case 'campus-tour': return 'bg-amber-500';
            case 'fair': return 'bg-red-500';
            default: return 'bg-gray-500';
        }
    }

    getCapacityPercentage(registered: number, capacity: number): number {
        if (capacity === 0) return 0;
        return Math.round((registered / capacity) * 100);
    }
}
