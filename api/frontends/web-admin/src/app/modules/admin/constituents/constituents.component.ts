import { ChangeDetectionStrategy, Component } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatMenuModule } from '@angular/material/menu';

interface Constituent {
    name: string;
    email: string;
    stage: string;
    program: string;
    term: string;
    source: string;
    score: number;
    owner: string;
    lastActivity: string;
}

@Component({
    selector: 'app-constituents',
    standalone: true,
    imports: [MatButtonModule, MatIconModule, MatMenuModule],
    changeDetection: ChangeDetectionStrategy.OnPush,
    templateUrl: './constituents.component.html',
})
export class ConstituentsComponent {
    readonly filters = ['Stage', 'Program', 'Source', 'Owner'];

    readonly constituents: readonly Constituent[] = [
        { name: 'Sofia Martinez', email: 'sofia.m@example.com', stage: 'Applicant', program: 'B.Sc. CS', term: 'Fall 2026', source: 'Web form', score: 86, owner: 'Avery', lastActivity: '2h ago' },
        { name: 'James Okoro', email: 'james.o@example.com', stage: 'Lead', program: 'M.A. Edu', term: 'Spring 2027', source: 'Open house', score: 71, owner: 'Priya', lastActivity: '5h ago' },
        { name: 'Liam Chen', email: 'liam.c@example.com', stage: 'Admitted', program: 'B.Sc. Bio', term: 'Fall 2026', source: 'Referral', score: 92, owner: 'Avery', lastActivity: 'Yesterday' },
        { name: 'Aisha Bello', email: 'aisha.b@example.com', stage: 'Applicant', program: 'M.B.A.', term: 'Fall 2026', source: 'Web form', score: 78, owner: 'Diego', lastActivity: '2d ago' },
        { name: 'Priya Patel', email: 'priya.p@example.com', stage: 'Inquiry', program: 'B.A. Econ', term: 'Fall 2026', source: 'Counselor', score: 64, owner: 'Avery', lastActivity: '3d ago' },
        { name: 'Noah Williams', email: 'noah.w@example.com', stage: 'Enrolled', program: 'B.Sc. CS', term: 'Fall 2026', source: 'Web form', score: 95, owner: 'Diego', lastActivity: '4d ago' },
    ];
}
