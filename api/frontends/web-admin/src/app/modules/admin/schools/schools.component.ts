import { DatePipe } from '@angular/common';
import { HttpErrorResponse } from '@angular/common/http';
import {
    ChangeDetectionStrategy,
    Component,
    computed,
    inject,
    signal,
} from '@angular/core';
import { toSignal } from '@angular/core/rxjs-interop';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { RouterLink } from '@angular/router';
import { jsonApiErrorMessage, JsonApiErrorResponse } from 'app/core/api/json-api';
import { SchoolAccessService } from 'app/core/school-access/school-access.service';
import { School } from 'app/core/school-access/school-access.types';
import { UserService } from 'app/core/user/user.service';

@Component({
    selector: 'app-schools',
    standalone: true,
    imports: [DatePipe, FormsModule, MatButtonModule, MatFormFieldModule, MatInputModule, RouterLink],
    changeDetection: ChangeDetectionStrategy.OnPush,
    templateUrl: './schools.component.html',
})
export class SchoolsComponent {
    private readonly schoolAccess = inject(SchoolAccessService);
    private readonly viewer = toSignal(inject(UserService).user$);

    /** Only SUPER_ADMIN may create schools; the server enforces this too. */
    readonly canCreate = computed(() => this.viewer()?.roles?.includes('SUPER_ADMIN') ?? false);
    readonly schools = signal<School[] | null>(null);
    readonly errorMessage = signal<string | null>(null);
    readonly name = signal('');
    readonly saving = signal(false);

    constructor() {
        this.load();
    }

    load(): void {
        this.schoolAccess.schools().subscribe({
            next: (schools) => this.schools.set(schools),
            error: (error: HttpErrorResponse) => {
                this.schools.set([]);
                this.errorMessage.set(jsonApiErrorMessage(error as JsonApiErrorResponse, 'Unable to load schools.'));
            },
        });
    }

    create(): void {
        const name = this.name().trim();
        if (!name) return;
        this.saving.set(true);
        this.errorMessage.set(null);
        this.schoolAccess.createSchool(name).subscribe({
            next: (school) => {
                this.saving.set(false);
                this.name.set('');
                this.schools.update((schools) =>
                    [...(schools ?? []), school].sort((a, b) => a.name.localeCompare(b.name))
                );
            },
            error: (error: HttpErrorResponse) => {
                this.saving.set(false);
                this.errorMessage.set(jsonApiErrorMessage(error as JsonApiErrorResponse, 'Unable to create the school.'));
            },
        });
    }
}
