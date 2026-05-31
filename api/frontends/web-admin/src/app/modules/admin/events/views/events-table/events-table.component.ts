import { ChangeDetectionStrategy, Component, signal, inject } from '@angular/core';
import { CommonModule, DatePipe } from '@angular/common';
import { RouterLink } from '@angular/router';
import { MatTableModule } from '@angular/material/table';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { MatMenuModule } from '@angular/material/menu';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MOCK_EVENTS } from '../../data/events.mock';
import { EventsComponent } from '../../events.component';

@Component({
    selector: 'app-events-table',
    standalone: true,
    imports: [
        CommonModule,
        RouterLink,
        MatTableModule,
        MatIconModule,
        MatButtonModule,
        MatMenuModule,
        MatProgressBarModule,
        DatePipe
    ],
    changeDetection: ChangeDetectionStrategy.OnPush,
    templateUrl: './events-table.component.html',
})
export class EventsTableComponent {
    private readonly shell = inject(EventsComponent);
    
    readonly displayedColumns = [
        'datetime',
        'title',
        'type',
        'status',
        'location',
        'capacity',
        'actions'
    ];
    
    readonly events = signal(MOCK_EVENTS);
    
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
    
    getStatusClasses(status: string): string {
        switch (status) {
            case 'upcoming': return 'bg-blue-100 text-blue-700 dark:bg-blue-900 dark:text-blue-200';
            case 'live': return 'bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-200';
            case 'completed': return 'bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-300';
            case 'draft': return 'bg-amber-100 text-amber-700 dark:bg-amber-900 dark:text-amber-200';
            case 'cancelled': return 'bg-red-100 text-red-700 dark:bg-red-900 dark:text-red-200';
            default: return 'bg-gray-100 text-gray-700';
        }
    }
    
    getCapacityPercentage(registered: number, capacity: number): number {
        if (capacity === 0) return 0;
        return Math.round((registered / capacity) * 100);
    }
}
