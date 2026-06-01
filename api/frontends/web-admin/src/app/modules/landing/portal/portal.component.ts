import {
    ChangeDetectionStrategy,
    Component,
    inject,
    signal,
} from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { Router, RouterLink } from '@angular/router';
import { jsonApiErrorMessage } from 'app/core/api/json-api';
import { PortalAuthService } from 'app/core/portal/portal-auth.service';

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
    imports: [
        FormsModule,
        MatButtonModule,
        MatFormFieldModule,
        MatIconModule,
        MatInputModule,
        RouterLink,
    ],
    templateUrl: './portal.component.html',
    changeDetection: ChangeDetectionStrategy.OnPush,
    host: { class: 'flex w-full flex-auto flex-col' },
})
export class PortalComponent {
    private readonly portalAuthService = inject(PortalAuthService);
    private readonly router = inject(Router);

    readonly portalSession = this.portalAuthService.session;
    readonly email = signal('');
    readonly requestingAccess = signal(false);
    readonly accessError = signal('');

    readonly cards: readonly PortalCard[] = [
        {
            title: 'Apply online',
            description: 'A guided application in under 30 minutes.',
            icon: 'heroicons_outline:pencil-square',
        },
        {
            title: 'Visit campus',
            description: 'Open days, tours, and virtual events.',
            icon: 'heroicons_outline:map-pin',
        },
        {
            title: 'Ask anything',
            description: 'Talk to admissions counselors directly.',
            icon: 'heroicons_outline:chat-bubble-left-right',
        },
        {
            title: 'Scholarships',
            description: 'Need-based and merit awards available.',
            icon: 'heroicons_outline:academic-cap',
        },
    ];

    readonly deadlines: readonly Deadline[] = [
        { date: 'Sep 15', title: 'Early decision deadline' },
        { date: 'Jan 10', title: 'Regular decision deadline' },
        { date: 'Mar 31', title: 'Financial aid forms due' },
    ];

    requestPortalAccess(): void {
        const email = this.email().trim();
        if (!email || this.requestingAccess()) {
            return;
        }

        this.requestingAccess.set(true);
        this.accessError.set('');

        this.portalAuthService.requestAccess(email).subscribe({
            next: () => {
                this.requestingAccess.set(false);
                void this.router.navigate(['/portal/status']);
            },
            error: (error) => {
                this.requestingAccess.set(false);
                this.accessError.set(
                    jsonApiErrorMessage(
                        error,
                        'We could not verify a submitted application for that email.'
                    )
                );
            },
        });
    }
}
