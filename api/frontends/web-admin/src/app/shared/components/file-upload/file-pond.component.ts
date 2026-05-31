import {
    ChangeDetectionStrategy,
    Component,
    DestroyRef,
    ElementRef,
    EventEmitter,
    Input,
    NgZone,
    Output,
    PLATFORM_ID,
    ViewChild,
    afterNextRender,
    inject,
} from '@angular/core';
import { isPlatformBrowser } from '@angular/common';
import * as FilePondLib from 'filepond';
import FilePondPluginImagePreview from 'filepond-plugin-image-preview';
import FilePondPluginImageExifOrientation from 'filepond-plugin-image-exif-orientation';
import FilePondPluginFileValidateType from 'filepond-plugin-file-validate-type';
import FilePondPluginFileValidateSize from 'filepond-plugin-file-validate-size';

// Register plugins once at module load (browser side only — guard below for SSR safety).
let _pluginsRegistered = false;
function registerPluginsOnce(): void {
    if (_pluginsRegistered) {
        return;
    }
    FilePondLib.registerPlugin(
        FilePondPluginImagePreview,
        FilePondPluginImageExifOrientation,
        FilePondPluginFileValidateType,
        FilePondPluginFileValidateSize
    );
    _pluginsRegistered = true;
}

/**
 * Thin standalone Angular wrapper around the FilePond core library.
 *
 * Why a custom wrapper?
 *   - `ngx-filepond` peer-deps cap at Angular 16; we run Angular 21.
 *   - Lets us mount FilePond inside Angular's lifecycle with SSR-safe rendering
 *     via `afterNextRender()` and event callbacks wrapped in NgZone so signals
 *     and change detection stay in sync.
 */
@Component({
    selector: 'app-file-pond',
    standalone: true,
    changeDetection: ChangeDetectionStrategy.OnPush,
    template: `<input type="file" #fileInput class="filepond" />`,
    styles: [
        `
            :host {
                display: block;
            }
        `,
    ],
})
export class FilePondComponent {
    @ViewChild('fileInput', { static: true })
    fileInputRef!: ElementRef<HTMLInputElement>;

    /** Accepted MIME types — e.g. `['image/*', 'application/pdf']`. */
    @Input() acceptedFileTypes: string[] | null = null;

    /** Maximum number of files allowed. */
    @Input() maxFiles = 1;

    /** Maximum single file size as a human-readable string, e.g. `'10MB'`. */
    @Input() maxFileSize: string | null = null;

    /** Allow multiple file uploads in one drop. */
    @Input() allowMultiple = false;

    /** Server configuration object matching FilePond's server contract. */
    @Input() server: FilePondLib.FilePondServerConfigProps['server'] | null = null;

    /** Placeholder label rendered inside the drop zone. */
    @Input() labelIdle =
        `Drag & drop files or <span class="filepond--label-action">browse</span>`;

    @Output() filesAdded = new EventEmitter<FilePondLib.FilePondFile[]>();
    @Output() fileUploaded = new EventEmitter<FilePondLib.FilePondFile>();
    @Output() fileReverted = new EventEmitter<FilePondLib.FilePondFile>();
    @Output() uploadError = new EventEmitter<{
        error: FilePondLib.FilePondErrorDescription | null;
        file: FilePondLib.FilePondFile | null;
    }>();

    private readonly _zone = inject(NgZone);
    private readonly _destroyRef = inject(DestroyRef);
    private readonly _platformId = inject(PLATFORM_ID);
    private _pond: FilePondLib.FilePond | null = null;

    constructor() {
        // `afterNextRender` only runs in the browser, making this SSR-safe.
        afterNextRender(() => {
            this._mountFilePond();
        });

        this._destroyRef.onDestroy(() => {
            if (this._pond) {
                this._pond.destroy();
                this._pond = null;
            }
        });
    }

    private _mountFilePond(): void {
        if (!isPlatformBrowser(this._platformId)) {
            return;
        }

        registerPluginsOnce();

        // Plugin-contributed options (acceptedFileTypes, maxFileSize) are not
        // on FilePondOptionProps; widen to a record so we can pass them through.
        const options: Record<string, unknown> = {
            allowMultiple: this.allowMultiple,
            maxFiles: this.maxFiles,
            labelIdle: this.labelIdle,
            credits: false,
        };

        if (this.acceptedFileTypes) {
            options['acceptedFileTypes'] = this.acceptedFileTypes;
        }
        if (this.maxFileSize) {
            options['maxFileSize'] = this.maxFileSize;
        }
        if (this.server) {
            options['server'] = this.server;
        }

        this._pond = FilePondLib.create(
            this.fileInputRef.nativeElement,
            options
        );

        // Re-enter the Angular zone for every callback so signals / outputs
        // trigger change detection.
        this._pond.on('addfile', (error, file) => {
            this._zone.run(() => {
                if (error) {
                    this.uploadError.emit({ error, file });
                    return;
                }
                this.filesAdded.emit([file]);
            });
        });

        this._pond.on('processfile', (error, file) => {
            this._zone.run(() => {
                if (error) {
                    this.uploadError.emit({ error, file });
                    return;
                }
                this.fileUploaded.emit(file);
            });
        });

        this._pond.on('processfilerevert', (file) => {
            this._zone.run(() => {
                this.fileReverted.emit(file);
            });
        });
    }
}
