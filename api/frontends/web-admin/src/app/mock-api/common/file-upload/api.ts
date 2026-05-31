import { Injectable } from '@angular/core';
import { FuseMockApiService, FuseMockApiUtils } from '@fuse/lib/mock-api';

@Injectable({ providedIn: 'root' })
export class FileUploadMockApi {
    /**
     * Constructor
     */
    constructor(private _fuseMockApiService: FuseMockApiService) {
        this.registerHandlers();
    }

    // -----------------------------------------------------------------------------------------------------
    // @ Public methods
    // -----------------------------------------------------------------------------------------------------

    /**
     * Register Mock API handlers
     *
     * FilePond server contract:
     * - POST: returns text/plain server-side file ID
     * - DELETE: revert a previously uploaded file by ID (request body = id)
     * See https://pqina.nl/filepond/docs/api/server/
     */
    registerHandlers(): void {
        this._fuseMockApiService.onPost('api/common/file-upload').reply(() => {
            const serverId = FuseMockApiUtils.guid();
            return [200, serverId];
        });

        this._fuseMockApiService
            .onDelete('api/common/file-upload')
            .reply(() => [200, '']);
    }
}
