import {
    ChangeDetectionStrategy,
    Component,
    computed,
    inject,
    signal,
} from '@angular/core';
import { toSignal } from '@angular/core/rxjs-interop';
import { MatIconModule } from '@angular/material/icon';
import { AuditService } from 'app/core/audit/audit.service';
import { AuditEvent } from 'app/core/audit/audit.types';
import { jsonApiErrorMessage, PaginatedResult } from 'app/core/api/json-api';
import { EMPTY, catchError } from 'rxjs';

interface AuditEventRow {
    id: string;
    timestamp: string;
    actor: string;
    action: string;
    entity: string;
    ip: string;
}

@Component({
    selector: 'app-audit',
    standalone: true,
    imports: [MatIconModule],
    changeDetection: ChangeDetectionStrategy.OnPush,
    templateUrl: './audit.component.html',
})
export class AuditComponent {
    private readonly auditService = inject(AuditService);

    readonly filters = ['Date range', 'Actor', 'Action'];
    readonly loading = signal(true);
    readonly errorMessage = signal<string | null>(null);
    readonly result = toSignal(this.auditService.auditEvents$, {
        initialValue: emptyResult<AuditEvent>(),
    });

    readonly events = computed(() =>
        this.result().items.map((event) => this.toRow(event))
    );

    constructor() {
        this.loadAuditEvents();
    }

    loadAuditEvents(): void {
        this.loading.set(true);
        this.errorMessage.set(null);

        this.auditService
            .queryAuditEvents({
                page: 1,
                rows: 50,
                orderBy: 'timestamp,DESC',
            })
            .pipe(
                catchError((error: unknown) => {
                    this.errorMessage.set(
                        jsonApiErrorMessage(error, 'Unable to load audit log.')
                    );
                    this.loading.set(false);

                    return EMPTY;
                })
            )
            .subscribe(() => this.loading.set(false));
    }

    private toRow(event: AuditEvent): AuditEventRow {
        return {
            id: event.id,
            timestamp: this.formatDate(event.timestamp),
            actor: event.actor_id,
            action: this.toTitleCase(event.action),
            entity: `${this.toTitleCase(event.obj_domain)} (${event.obj_name})`,
            ip: this.ipFromData(event.data),
        };
    }

    private ipFromData(data: string): string {
        if (!data) {
            return '—';
        }

        try {
            const parsed = JSON.parse(data) as { ip?: unknown };

            return typeof parsed.ip === 'string' ? parsed.ip : '—';
        } catch {
            return '—';
        }
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
