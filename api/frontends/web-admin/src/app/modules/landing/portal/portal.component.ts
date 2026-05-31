import { ChangeDetectionStrategy, Component } from '@angular/core';
import { RouterLink } from '@angular/router';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';

interface PortalCard {
    title: string;
    description: string;
    icon: string;
}

interface Deadline {
    date: string;
    title: string;
}

@Component({
    selector: 'app-portal',
    standalone: true,
    imports: [RouterLink, MatIconModule, MatButtonModule],
    templateUrl: './portal.component.html',
    changeDetection: ChangeDetectionStrategy.OnPush,
    host: { class: 'flex w-full flex-auto flex-col' },
})
export class PortalComponent {
    readonly cards: readonly PortalCard[] = [
        {
            title: 'Apply online',
            description: 'A guided application in under 30 minutes.',
            icon: 'heroicons_outline:pencil-square'
        },
        {
            title: 'Visit campus',
            description: 'Open days, tours, and virtual events.',
            icon: 'heroicons_outline:map-pin'
        },
        {
            title: 'Ask anything',
            description: 'Talk to admissions counselors directly.',
            icon: 'heroicons_outline:chat-bubble-left-right'
        },
        {
            title: 'Scholarships',
            description: 'Need-based and merit awards available.',
            icon: 'heroicons_outline:academic-cap'
        }
    ];

    readonly deadlines: readonly Deadline[] = [
        { date: 'Sep 15', title: 'Early decision deadline' },
        { date: 'Jan 10', title: 'Regular decision deadline' },
        { date: 'Mar 31', title: 'Financial aid forms due' }
    ];
}
