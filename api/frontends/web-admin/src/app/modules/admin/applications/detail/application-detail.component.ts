import { ChangeDetectionStrategy, Component, inject } from '@angular/core';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatTabsModule } from '@angular/material/tabs';
import { ApplicationDocumentsComponent } from '../components/documents/application-documents.component';

interface TimelineEvent {
    icon: string;
    title: string;
    date: string;
}

interface Note {
    authorInitials: string;
    authorName: string;
    timestamp: string;
    body: string;
}

interface Insight {
    title: string;
    description: string;
}

@Component({
    selector: 'app-application-detail',
    standalone: true,
    imports: [RouterLink, MatButtonModule, MatIconModule, MatTabsModule, ApplicationDocumentsComponent],
    changeDetection: ChangeDetectionStrategy.OnPush,
    templateUrl: './application-detail.component.html',
})
export class ApplicationDetailComponent {
    private readonly route = inject(ActivatedRoute);
    readonly applicationId = this.route.snapshot.paramMap.get('id');

    readonly applicantName = 'Sofia Martinez';
    readonly status = 'In review';
    
    readonly summaryCards = [
        { label: 'Program', value: 'B.Sc. Computer Science' },
        { label: 'Term', value: 'Fall 2026' },
        { label: 'Submitted', value: '2 hours ago' },
        { label: 'Score', value: '86 / 100' },
    ];

    readonly applicantInfo = [
        { label: 'Name', value: 'Sofia Martinez' },
        { label: 'Email', value: 'sofia.m@example.com' },
        { label: 'Phone', value: '+1 (555) 123-4567' },
        { label: 'DOB', value: 'May 14, 2008' },
        { label: 'Address', value: '123 Tech Lane, San Francisco, CA 94105' },
        { label: 'GPA', value: '3.9' },
        { label: 'Test Scores', value: 'SAT: 1450' },
        { label: 'Essay excerpt', value: '"My passion for computer science began when I built my first robot..."' },
    ];

    readonly notes: Note[] = [
        { authorInitials: 'AM', authorName: 'Avery', timestamp: '1 hour ago', body: 'Strong candidate. Great math scores.' },
        { authorInitials: 'PP', authorName: 'Priya', timestamp: 'Yesterday', body: 'Missing one recommendation letter, but otherwise complete.' },
        { authorInitials: 'SM', authorName: 'System', timestamp: '2 days ago', body: 'Application submitted successfully.' },
    ];

    readonly timelineEvents: TimelineEvent[] = [
        { icon: 'heroicons_outline:check-circle', title: 'Review started', date: 'Today, 9:00 AM' },
        { icon: 'heroicons_outline:document-text', title: 'All documents received', date: 'Yesterday, 2:30 PM' },
        { icon: 'heroicons_outline:arrow-up-tray', title: 'Transcript uploaded', date: 'Yesterday, 10:15 AM' },
        { icon: 'heroicons_outline:envelope', title: 'Recommendation 2 received', date: 'May 28, 4:00 PM' },
        { icon: 'heroicons_outline:envelope', title: 'Recommendation 1 received', date: 'May 27, 11:20 AM' },
        { icon: 'heroicons_outline:paper-airplane', title: 'Application submitted', date: 'May 26, 8:45 PM' },
    ];

    readonly insights: Insight[] = [
        { title: 'Recommended decision', description: 'Strong admit based on GPA and test scores.' },
        { title: 'Risk flags', description: 'None identified.' },
        { title: 'Similar profiles', description: 'Matches 85% of successful admits in the past 3 years.' },
    ];
}
