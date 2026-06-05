import {
    ChangeDetectionStrategy,
    Component,
    computed,
    inject,
    signal,
} from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatMenuModule } from '@angular/material/menu';
import { toSignal } from '@angular/core/rxjs-interop';
import { AdmissionsService } from 'app/core/admissions/admissions.service';
import { Constituent } from 'app/core/admissions/admissions.types';
import { jsonApiErrorMessage, PaginatedResult } from 'app/core/api/json-api';
import { EMPTY, catchError } from 'rxjs';

interface ConstituentRow {
    name: string;
    email: string;
    stage: string;
    program: string;
    term: string;
    source: string;
    score: number;
    owner: string;
    lastActivity: string;
}

@Component({
    selector: 'app-constituents',
    standalone: true,
    imports: [MatButtonModule, MatIconModule, MatMenuModule],
    changeDetection: ChangeDetectionStrategy.OnPush,
    templateUrl: './constituents.component.html',
})
export class ConstituentsComponent {
    private readonly admissionsService = inject(AdmissionsService);

    readonly filters = ['Stage', 'Program', 'Source', 'Owner'];
    readonly loading = signal(true);
    readonly errorMessage = signal<string | null>(null);
    readonly result = toSignal(this.admissionsService.constituents$, {
        initialValue: emptyResult<Constituent>(),
    });

    readonly constituents = computed(() =>
        this.result().items.map((constituent) => this.toRow(constituent))
    );

    constructor() {
        this.loadConstituents();
    }

    loadConstituents(): void {
        this.loading.set(true);
        this.errorMessage.set(null);

        this.admissionsService
            .queryConstituents({
                page: 1,
                rows: 50,
                orderBy: 'date_created,DESC',
            })
            .pipe(
                catchError((error: unknown) => {
                    this.errorMessage.set(
                        jsonApiErrorMessage(
                            error,
                            'Unable to load constituents.'
                        )
                    );
                    this.loading.set(false);

                    return EMPTY;
                })
            )
            .subscribe(() => this.loading.set(false));
    }

    private toRow(constituent: Constituent): ConstituentRow {
        return {
            name: `${constituent.firstName} ${constituent.lastName}`,
            email: constituent.primaryEmail,
            stage: this.toTitleCase(constituent.lifecycleStage),
            program: constituent.externalSISID ?? 'Admissions record',
            term: constituent.sisSyncedAt ? 'Synced' : 'Unmatched',
            source: constituent.nationalID ?? constituent.upi ?? 'Direct entry',
            score: this.stageScore(constituent.lifecycleStage),
            owner: constituent.notificationPreferences.priority[0] ?? 'Admissions',
            lastActivity: this.formatDate(constituent.dateUpdated),
        };
    }

    private stageScore(stage: Constituent['lifecycleStage']): number {
        const scores: Record<Constituent['lifecycleStage'], number> = {
            PROSPECT: 35,
            INQUIRY: 50,
            APPLICANT: 75,
            ADMITTED: 88,
            ENROLLED: 96,
            ALUMNI: 70,
        };

        return scores[stage];
    }

    private toTitleCase(value: string): string {
        return value
            .toLowerCase()
            .replace(/_/g, ' ')
            .replace(/\b\w/g, (character) => character.toUpperCase());
    }

    private formatDate(value: string): string {
        return new Intl.DateTimeFormat('en-US', {
            month: 'short',
            day: 'numeric',
            hour: 'numeric',
            minute: '2-digit',
        }).format(new Date(value));
    }
}

function emptyResult<T>(): PaginatedResult<T> {
    return {
        items: [],
        total: 0,
        page: 1,
        rowsPerPage: 50,
    };
}
