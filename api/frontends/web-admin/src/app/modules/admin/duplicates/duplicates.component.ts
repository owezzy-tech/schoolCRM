import { ChangeDetectionStrategy, Component } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';

interface ConstituentRecord {
    name: string;
    email: string;
    program: string;
    createdAt: string;
    activity: string;
}

interface DuplicatePair {
    id: string;
    score: number;
    left: ConstituentRecord;
    right: ConstituentRecord;
}

@Component({
    selector: 'app-duplicates',
    standalone: true,
    imports: [MatButtonModule, MatIconModule],
    changeDetection: ChangeDetectionStrategy.OnPush,
    templateUrl: './duplicates.component.html',
})
export class DuplicatesComponent {
    readonly pairs: readonly DuplicatePair[] = [
        {
            id: 'DUP-001',
            score: 97,
            left: { name: 'Sofia Martinez', email: 'sofia.m@example.com', program: 'B.Sc. CS', createdAt: 'Jan 12, 2026', activity: '12 events' },
            right: { name: 'Sofia Martínez', email: 'sofia.martinez@example.com', program: 'B.Sc. Computer Science', createdAt: 'Mar 4, 2026', activity: '3 events' },
        },
        {
            id: 'DUP-002',
            score: 89,
            left: { name: 'James Okoro', email: 'james.o@example.com', program: 'M.A. Education', createdAt: 'Feb 2, 2026', activity: '8 events' },
            right: { name: 'J. Okoro', email: 'jokoro@example.com', program: 'M.A. Edu', createdAt: 'Apr 18, 2026', activity: '2 events' },
        },
        {
            id: 'DUP-003',
            score: 82,
            left: { name: 'Aisha Bello', email: 'aisha.b@example.com', program: 'M.B.A.', createdAt: 'Jan 30, 2026', activity: '5 events' },
            right: { name: 'Ayisha Bello', email: 'a.bello@example.com', program: 'MBA', createdAt: 'May 1, 2026', activity: '1 event' },
        },
    ];
}
