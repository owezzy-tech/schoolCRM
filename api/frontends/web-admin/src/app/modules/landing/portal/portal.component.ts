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
            description: 'Complete KUCCPS or self-sponsored applications online.',
            icon: 'heroicons_outline:pencil-square',
        },
        {
            title: 'Track intakes',
            description: 'Follow placement, document review, and admission status.',
            icon: 'heroicons_outline:map-pin',
        },
        {
            title: 'Ask anything',
            description: 'Reach the admissions office by email, phone, or WhatsApp.',
            icon: 'heroicons_outline:chat-bubble-left-right',
        },
        {
            title: 'Pay locally',
            description: 'Application fees are tracked in KES with M-Pesa support.',
            icon: 'heroicons_outline:academic-cap',
        },
    ];

    readonly deadlines: readonly Deadline[] = [
        { date: 'Sep 1', title: 'KUCCPS placement confirmation opens' },
        { date: 'Jan 15', title: '2026 main intake application deadline' },
        { date: 'Mar 31', title: 'KCSE and document verification closes' },
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
