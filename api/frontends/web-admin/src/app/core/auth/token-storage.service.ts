import { isPlatformBrowser } from '@angular/common';
import { inject, Injectable, PLATFORM_ID } from '@angular/core';

@Injectable({ providedIn: 'root' })
export class TokenStorageService {
    private readonly platformId = inject(PLATFORM_ID);
    private readonly key = 'accessToken';

    get(): string | null {
        if (!isPlatformBrowser(this.platformId)) {
            return null;
        }

        return localStorage.getItem(this.key);
    }

    set(token: string): void {
        if (!isPlatformBrowser(this.platformId)) {
            return;
        }

        localStorage.setItem(this.key, token);
    }

    clear(): void {
        if (!isPlatformBrowser(this.platformId)) {
            return;
        }

        localStorage.removeItem(this.key);
    }
}
