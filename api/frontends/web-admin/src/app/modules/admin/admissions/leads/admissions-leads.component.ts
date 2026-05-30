import { ChangeDetectionStrategy, Component, computed, inject } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { RouterLink } from '@angular/router';
import { MatButtonModule } from '@angular/material/button';
import { MatChipsModule } from '@angular/material/chips';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatSelectModule } from '@angular/material/select';
import { toSignal } from '@angular/core/rxjs-interop';
import { AdmissionsService } from 'app/core/admissions/admissions.service';
import { LEAD_SCORE_BANDS, LeadScore, LeadScoreBand } from 'app/core/admissions/admissions.types';
import { jsonApiErrorMessage } from 'app/core/api/json-api';
import { catchError, of } from 'rxjs';

interface LeadColumn {
    band: LeadScoreBand;
    title: string;
    description: string;
    scores: LeadScore[];
}

@Component({
    selector: 'app-admissions-leads',
    imports: [
        FormsModule,
        RouterLink,
        MatButtonModule,
        MatChipsModule,
        MatFormFieldModule,
        MatIconModule,
        MatSelectModule,
    ],
    templateUrl: './admissions-leads.component.html',
    changeDetection: ChangeDetectionStrategy.OnPush,
})
export class AdmissionsLeadsComponent {
    private readonly admissionsService = inject(AdmissionsService);

    readonly bands = LEAD_SCORE_BANDS;
    readonly scores$ = this.admissionsService.scores$;
    readonly scores = toSignal(this.scores$, {
        initialValue: { items: [], total: 0, page: 1, rowsPerPage: 25 },
    });
    readonly columns = computed<LeadColumn[]>(() =>
        this.bands.map((band) => ({
            band,
            title: this.bandTitle(band),
            description: this.bandDescription(band),
            scores: this.scores().items.filter((score) => score.band === band),
        }))
    );

    errorMessage = '';
    selectedBand: LeadScoreBand | '' = '';

    constructor() {
        this.loadScores();
    }

    loadScores(): void {
        this.errorMessage = '';
        this.admissionsService
            .queryLeadScores({
                page: 1,
                rows: 50,
                orderBy: 'total_score,DESC',
                band: this.selectedBand || undefined,
            })
            .pipe(
                catchError((error) => {
                    this.errorMessage = jsonApiErrorMessage(
                        error,
                        'Unable to load lead scores.'
                    );
                    return of(undefined);
                })
            )
            .subscribe();
    }

    bandClass(band: LeadScoreBand): string {
        switch (band) {
            case 'READY_TO_APPLY':
                return 'border-purple-500 bg-purple-50 text-purple-800';
            case 'HOT':
                return 'border-red-500 bg-red-50 text-red-800';
            case 'WARM':
                return 'border-orange-500 bg-orange-50 text-orange-800';
            case 'COLD':
                return 'border-blue-500 bg-blue-50 text-blue-800';
        }
    }

    private bandTitle(band: LeadScoreBand): string {
        return band.replaceAll('_', ' ');
    }

    private bandDescription(band: LeadScoreBand): string {
        switch (band) {
            case 'READY_TO_APPLY':
                return '76+ points';
            case 'HOT':
                return '51–75 points';
            case 'WARM':
                return '26–50 points';
            case 'COLD':
                return '0–25 points';
        }
    }
}
