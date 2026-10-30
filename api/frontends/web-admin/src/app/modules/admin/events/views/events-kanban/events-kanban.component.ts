import { ChangeDetectionStrategy, Component, computed, inject } from '@angular/core';
import { CommonModule, DatePipe } from '@angular/common';
import { RouterLink } from '@angular/router';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { DragDropModule, CdkDragDrop, moveItemInArray, transferArrayItem } from '@angular/cdk/drag-drop';
import { EventItem } from '../../models/event.types';
import { EventsComponent } from '../../events.component';

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
    private readonly shell = inject(EventsComponent);

    readonly columns = computed(() => {
        const events = this.shell.events();

        return [
            {
                id: 'draft',
                title: 'Draft',
                events: events.filter((event) => event.status === 'draft'),
            },
            {
                id: 'upcoming',
                title: 'Upcoming',
                events: events.filter((event) => event.status === 'upcoming'),
            },
            {
                id: 'live',
                title: 'Live',
                events: events.filter((event) => event.status === 'live'),
            },
            {
                id: 'completed',
                title: 'Completed',
                events: events.filter((event) => event.status === 'completed'),
            },
        ];
    });

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
