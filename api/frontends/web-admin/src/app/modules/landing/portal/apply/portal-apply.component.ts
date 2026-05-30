import { ChangeDetectionStrategy, Component } from '@angular/core';

@Component({
    selector: 'app-portal-apply',
    standalone: true,
    template: `
        <div class="flex flex-col flex-auto items-center justify-center py-20 px-6">
            <div class="max-w-2xl w-full text-center">
                <h1 class="text-4xl font-bold text-default mb-4">Apply online</h1>
                <p class="text-lg text-secondary mb-10">A guided application in under 30 minutes.</p>
                
                <div class="p-12 rounded-2xl bg-card shadow-sm border flex flex-col items-center">
                    <span class="text-xl font-medium text-default mb-2">Application portal coming soon</span>
                    <span class="text-secondary">We are preparing the 2026 application forms.</span>
                </div>
            </div>
        </div>
    `,
    changeDetection: ChangeDetectionStrategy.OnPush,
    host: { class: 'flex w-full flex-auto flex-col' },
})
export class PortalApplyComponent {}
