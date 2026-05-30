import { ChangeDetectionStrategy, Component } from '@angular/core';

@Component({
    selector: 'app-portal-status',
    standalone: true,
    template: `
        <div class="flex flex-col flex-auto items-center justify-center py-20 px-6">
            <div class="max-w-2xl w-full text-center">
                <h1 class="text-4xl font-bold text-default mb-4">Check status</h1>
                <p class="text-lg text-secondary mb-10">Follow your application status in real time.</p>
                
                <div class="p-12 rounded-2xl bg-card shadow-sm border flex flex-col items-center">
                    <span class="text-xl font-medium text-default mb-2">Status portal coming soon</span>
                    <span class="text-secondary">Please check back later to view your application status.</span>
                </div>
            </div>
        </div>
    `,
    changeDetection: ChangeDetectionStrategy.OnPush,
    host: { class: 'flex w-full flex-auto flex-col' },
})
export class PortalStatusComponent {}
