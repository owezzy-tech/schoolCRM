import { ChangeDetectionStrategy, Component } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';

interface Inquiry {
    id: string;
    from: string;
    program: string;
    source: string;
    priority: 'Low' | 'Medium' | 'High';
    status: 'New' | 'In progress' | 'Closed';
    assigned: string;
    submitted: string;
}

@Component({
    selector: 'app-inquiries',
    standalone: true,
    imports: [MatButtonModule, MatIconModule],
    changeDetection: ChangeDetectionStrategy.OnPush,
    templateUrl: './inquiries.component.html',
})
export class InquiriesComponent {
    readonly filters = ['Status', 'Priority', 'Program'];

    readonly inquiries: readonly Inquiry[] = [
        { id: 'INQ-2048', from: 'Maya Tran', program: 'B.Sc. CS', source: 'Web form', priority: 'High', status: 'New', assigned: 'Avery', submitted: '15m ago' },
        { id: 'INQ-2047', from: 'Daniel Kim', program: 'M.B.A.', source: 'Open house', priority: 'Medium', status: 'In progress', assigned: 'Priya', submitted: '1h ago' },
        { id: 'INQ-2046', from: 'Zara Khan', program: 'B.A. Econ', source: 'Counselor', priority: 'Low', status: 'New', assigned: '—', submitted: '3h ago' },
        { id: 'INQ-2045', from: 'Ethan Lee', program: 'B.Sc. Bio', source: 'Web form', priority: 'High', status: 'In progress', assigned: 'Diego', submitted: '6h ago' },
        { id: 'INQ-2044', from: 'Mira Joshi', program: 'M.A. Edu', source: 'Referral', priority: 'Medium', status: 'Closed', assigned: 'Avery', submitted: 'Yesterday' },
    ];
}
