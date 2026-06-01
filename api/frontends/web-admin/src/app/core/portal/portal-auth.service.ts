import { isPlatformBrowser } from '@angular/common';
import { HttpClient } from '@angular/common/http';
import { inject, Injectable, PLATFORM_ID, signal } from '@angular/core';
import { JsonApiDocument, unwrapJsonApiResource } from 'app/core/api/json-api';
import { AuthUtils } from 'app/core/auth/auth.utils';
import { map, Observable, tap } from 'rxjs';

export interface ApplicantPortalSession {
    id: string;
    accessToken: string;
    tokenType: 'Bearer';
    expiresAt: string;
    expiresIn: number;
    applicationID: string;
    constituentID: string;
    applicantName: string;
    email: string;
}

@Injectable({ providedIn: 'root' })
export class PortalAuthService {
    private readonly httpClient = inject(HttpClient);
    private readonly platformId = inject(PLATFORM_ID);
    private readonly storageKey = 'applicantPortalSession';
    private readonly sessionSignal = signal<ApplicantPortalSession | null>(
        this.loadSession()
    );

    readonly session = this.sessionSignal.asReadonly();

    requestAccess(email: string): Observable<ApplicantPortalSession> {
        return this.httpClient
            .post<
                JsonApiDocument<ApplicantPortalSession>
            >('/v1/auth/applicant-portal/token', { email })
            .pipe(
                map(unwrapJsonApiResource),
                tap((session) => this.setSession(session))
            );
    }

    hasValidSession(): boolean {
        const session = this.sessionSignal();

        return Boolean(
            session?.accessToken &&
                !AuthUtils.isTokenExpired(session.accessToken)
        );
    }

    clear(): void {
        this.sessionSignal.set(null);

        if (!isPlatformBrowser(this.platformId)) {
            return;
        }

        localStorage.removeItem(this.storageKey);
    }

    private setSession(session: ApplicantPortalSession): void {
        this.sessionSignal.set(session);

        if (!isPlatformBrowser(this.platformId)) {
            return;
        }

        localStorage.setItem(this.storageKey, JSON.stringify(session));
    }

    private loadSession(): ApplicantPortalSession | null {
        if (!isPlatformBrowser(this.platformId)) {
            return null;
        }

        const rawSession = localStorage.getItem(this.storageKey);
        if (!rawSession) {
            return null;
        }

        try {
            const session = JSON.parse(rawSession) as ApplicantPortalSession;

            if (
                !session.accessToken ||
                AuthUtils.isTokenExpired(session.accessToken)
            ) {
                localStorage.removeItem(this.storageKey);
                return null;
            }

            return session;
        } catch {
            localStorage.removeItem(this.storageKey);
            return null;
        }
    }
}
