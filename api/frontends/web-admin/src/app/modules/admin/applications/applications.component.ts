import { ChangeDetectionStrategy, Component } from '@angular/core';
import { RouterLink } from '@angular/router';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';

interface Application {
    id: string;
    applicant: string;
    program: string;
    term: string;
    status: string;
    progress: number;
    reviewer: string;
    score: number | null;
}

@Component({
    selector: 'app-applications',
    standalone: true,
    imports: [RouterLink, MatButtonModule, MatIconModule],
    changeDetection: ChangeDetectionStrategy.OnPush,
    templateUrl: './applications.component.html',
})
export class ApplicationsComponent {
    readonly filters = ['Status', 'Program', 'Term', 'Reviewer'];

    readonly applications: readonly Application[] = [
        { id: 'APP-3024', applicant: 'Sofia Martinez', program: 'B.Sc. CS', term: 'Fall 2026', status: 'In review', progress: 80, reviewer: 'Avery', score: 86 },
        { id: 'APP-3023', applicant: 'James Okoro', program: 'M.A. Edu', term: 'Spring 2027', status: 'Documents pending', progress: 45, reviewer: 'Priya', score: null },
        { id: 'APP-3022', applicant: 'Priya Patel', program: 'B.A. Econ', term: 'Fall 2026', status: 'Submitted', progress: 60, reviewer: 'Diego', score: 71 },
        { id: 'APP-3021', applicant: 'Liam Chen', program: 'B.Sc. Bio', term: 'Fall 2026', status: 'Interview', progress: 90, reviewer: 'Avery', score: 92 },
        { id: 'APP-3020', applicant: 'Aisha Bello', program: 'M.B.A.', term: 'Fall 2026', status: 'Decision pending', progress: 95, reviewer: 'Diego', score: 78 },
        { id: 'APP-3019', applicant: 'Noah Williams', program: 'B.Sc. CS', term: 'Fall 2026', status: 'Admitted', progress: 100, reviewer: 'Avery', score: 95 },
    ];
}
