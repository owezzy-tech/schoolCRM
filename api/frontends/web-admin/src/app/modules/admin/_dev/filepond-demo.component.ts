import { ChangeDetectionStrategy, Component } from '@angular/core';
import { FilePondComponent } from 'app/shared/components/file-upload/file-pond.component';

@Component({
    selector: 'app-filepond-demo',
    standalone: true,
    imports: [FilePondComponent],
    changeDetection: ChangeDetectionStrategy.OnPush,
    template: `
        <div class="flex flex-col flex-auto min-w-0 p-6 sm:p-10">
            <div class="text-3xl font-bold tracking-tight">FilePond Demo</div>
            <div class="mt-2 text-secondary">
                Phase 1 sanity check — drop a file to exercise the mock
                <code>api/common/file-upload</code> endpoint.
            </div>

            <div class="mt-8 max-w-2xl">
                <div class="text-lg font-semibold mb-3">Single image upload</div>
                <app-file-pond
                    [acceptedFileTypes]="['image/*']"
                    [maxFileSize]="'5MB'"
                    [server]="serverConfig"
                />
            </div>

            <div class="mt-10 max-w-2xl">
                <div class="text-lg font-semibold mb-3">
                    Multiple document upload
                </div>
                <app-file-pond
                    [allowMultiple]="true"
                    [maxFiles]="5"
                    [acceptedFileTypes]="['application/pdf', 'image/*']"
                    [maxFileSize]="'10MB'"
                    [server]="serverConfig"
                />
            </div>
        </div>
    `,
})
export class FilePondDemoComponent {
    readonly serverConfig = {
        url: '/api/common/file-upload',
        process: {
            url: '',
            method: 'POST' as const,
        },
        revert: {
            url: '',
            method: 'DELETE' as const,
        },
    };
}
