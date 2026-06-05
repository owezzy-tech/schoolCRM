import {
    ChangeDetectionStrategy,
    Component,
    computed,
    inject,
    signal,
} from '@angular/core';
import { toSignal } from '@angular/core/rxjs-interop';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { AdmissionsService } from 'app/core/admissions/admissions.service';
import { Inquiry } from 'app/core/admissions/admissions.types';
import { jsonApiErrorMessage, PaginatedResult } from 'app/core/api/json-api';
import { EMPTY, catchError } from 'rxjs';

interface InquiryRow {
    id: string;
    from: string;
    program: string;
    source: string;
    priority: 'Low' | 'Medium' | 'High';
    status: 'New' | 'In progress' | 'Closed';
    assigned: string;
    submitted: string;
}

@Component({
    selector: 'app-inquiries',
    standalone: true,
    imports: [MatButtonModule, MatIconModule],
    changeDetection: ChangeDetectionStrategy.OnPush,
    templateUrl: './inquiries.component.html',
})
export class InquiriesComponent {
    private readonly admissionsService = inject(AdmissionsService);

    readonly filters = ['Status', 'Priority', 'Program'];
    readonly loading = signal(true);
    readonly errorMessage = signal<string | null>(null);
    readonly result = toSignal(this.admissionsService.inquiries$, {
        initialValue: emptyResult<Inquiry>(),
    });

    readonly inquiries = computed(() =>
        this.result().items.map((inquiry) => this.toRow(inquiry))
    );

    constructor() {
        this.loadInquiries();
    }

    loadInquiries(): void {
        this.loading.set(true);
        this.errorMessage.set(null);

        this.admissionsService
            .queryInquiries({
                page: 1,
                rows: 50,
                orderBy: 'date_created,DESC',
            })
            .pipe(
                catchError((error: unknown) => {
                    this.errorMessage.set(
                        jsonApiErrorMessage(error, 'Unable to load inquiries.')
                    );
                    this.loading.set(false);

                    return EMPTY;
                })
            )
            .subscribe(() => this.loading.set(false));
    }

    private toRow(inquiry: Inquiry): InquiryRow {
        return {
            id: inquiry.id,
            from: `${inquiry.firstName} ${inquiry.lastName}`,
            program: inquiry.programOfInterest ?? 'Unspecified',
            source: this.toTitleCase(inquiry.source),
            priority: this.priorityForStatus(inquiry.status),
            status: this.statusLabel(inquiry.status),
            assigned: 'Admissions queue',
            submitted: this.formatDate(inquiry.dateCreated),
        };
    }

    private priorityForStatus(status: string): InquiryRow['priority'] {
        if (status === 'NEW') {
            return 'High';
        }

        if (status === 'CONTACTED') {
            return 'Medium';
        }

        return 'Low';
    }

    private statusLabel(status: string): InquiryRow['status'] {
        const labels: Record<string, InquiryRow['status']> = {
            NEW: 'New',
            CONTACTED: 'In progress',
            CONVERTED: 'Closed',
            CLOSED: 'Closed',
        };

        return labels[status] ?? 'New';
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
