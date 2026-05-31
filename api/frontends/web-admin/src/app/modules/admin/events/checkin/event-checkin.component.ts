import { ChangeDetectionStrategy, Component, signal, inject, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterLink, ActivatedRoute } from '@angular/router';
import { ReactiveFormsModule, UntypedFormControl } from '@angular/forms';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { MatInputModule } from '@angular/material/input';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MOCK_EVENTS } from '../data/events.mock';
import { EventItem, EventRegistration } from '../models/event.types';

@Component({
    selector: 'app-event-checkin',
    standalone: true,
    imports: [
        CommonModule,
        RouterLink,
        ReactiveFormsModule,
        MatIconModule,
        MatButtonModule,
        MatInputModule,
        MatFormFieldModule
    ],
    changeDetection: ChangeDetectionStrategy.OnPush,
    templateUrl: './event-checkin.component.html'
})
export class EventCheckinComponent implements OnInit {
    private readonly route = inject(ActivatedRoute);
    
    readonly event = signal<EventItem | null>(null);
    readonly searchControl = new UntypedFormControl('');
    
    // Mock registrants
    readonly registrants = signal<EventRegistration[]>([
        { id: 'reg_1', constituentId: 'c_1', constituentName: 'Sofia Martinez', status: 'registered', registeredAt: '2026-05-10T10:00:00Z' },
        { id: 'reg_2', constituentId: 'c_2', constituentName: 'James Okoro', status: 'checked-in', registeredAt: '2026-05-12T14:30:00Z' },
        { id: 'reg_3', constituentId: 'c_3', constituentName: 'Liam Chen', status: 'registered', registeredAt: '2026-05-15T09:15:00Z' },
        { id: 'reg_4', constituentId: 'c_4', constituentName: 'Aisha Bello', status: 'registered', registeredAt: '2026-05-16T11:45:00Z' },
        { id: 'reg_5', constituentId: 'c_5', constituentName: 'Priya Patel', status: 'checked-in', registeredAt: '2026-05-18T16:20:00Z' },
    ]);

    ngOnInit() {
        const id = this.route.snapshot.paramMap.get('id');
        const found = MOCK_EVENTS.find(e => e.id === id);
        if (found) {
            this.event.set(found);
        } else {
            this.event.set(MOCK_EVENTS[0]);
        }
    }
    
    get filteredRegistrants() {
        const term = (this.searchControl.value || '').toLowerCase();
        if (!term) return this.registrants();
        return this.registrants().filter(r => r.constituentName.toLowerCase().includes(term));
    }
    
    get checkedInCount() {
        return this.registrants().filter(r => r.status === 'checked-in').length;
    }
    
    toggleCheckin(id: string) {
        this.registrants.update(regs => 
            regs.map(r => {
                if (r.id === id) {
                    return { ...r, status: r.status === 'checked-in' ? 'registered' : 'checked-in' };
                }
                return r;
            })
        );
    }
}
