import { TestBed } from '@angular/core/testing';
import { AuditService } from 'app/core/audit/audit.service';
import { AuditEvent } from 'app/core/audit/audit.types';
import { PaginatedResult } from 'app/core/api/json-api';
import { of, ReplaySubject, throwError } from 'rxjs';
import { beforeEach, describe, expect, it, vi } from 'vitest';

import { AuditComponent } from './audit.component';

describe('AuditComponent', () => {
    let auditEventsSubject: ReplaySubject<PaginatedResult<AuditEvent>>;
    let queryAuditEvents: ReturnType<typeof vi.fn>;

    beforeEach(async () => {
        auditEventsSubject = new ReplaySubject<PaginatedResult<AuditEvent>>(1);
        queryAuditEvents = vi.fn(() => of(auditEventsResult([auditEvent])));

        await TestBed.configureTestingModule({
            imports: [AuditComponent],
            providers: [
                {
                    provide: AuditService,
                    useValue: {
                        auditEvents$: auditEventsSubject.asObservable(),
                        queryAuditEvents,
                    },
                },
            ],
        })
            .overrideComponent(AuditComponent, {
                set: {
                    template: '',
                },
            })
            .compileComponents();
    });

    it('loads audit events from the audit API', () => {
        TestBed.createComponent(AuditComponent);

        expect(queryAuditEvents).toHaveBeenCalledWith({
            page: 1,
            rows: 50,
            orderBy: 'timestamp,DESC',
        });
    });

    it('derives audit table rows from the service stream', () => {
        const fixture = TestBed.createComponent(AuditComponent);
        const component = fixture.componentInstance;

        auditEventsSubject.next(auditEventsResult([auditEvent]));

        expect(component.events()[0]).toMatchObject({
            actor: 'user-admin',
            action: 'Sync Push',
            entity: 'Admissions Communication (Campus tour reminder)',
            ip: '192.168.1.105',
        });
    });

    it('surfaces JSON:API errors when audit loading fails', () => {
        queryAuditEvents.mockReturnValue(
            throwError(() => ({
                error: {
                    errors: [
                        {
                            detail: 'Audit service unavailable.',
                        },
                    ],
                },
            }))
        );

        const fixture = TestBed.createComponent(AuditComponent);
        const component = fixture.componentInstance;

        expect(component.loading()).toBe(false);
        expect(component.errorMessage()).toBe('Audit service unavailable.');
    });
});

function auditEventsResult(items: AuditEvent[]): PaginatedResult<AuditEvent> {
    return {
        items,
        total: items.length,
        page: 1,
        rowsPerPage: 50,
    };
}

const auditEvent: AuditEvent = {
    id: 'audit-communication-failure',
    obj_id: 'communication-sms-failed',
    obj_domain: 'admissions_communication',
    obj_name: 'Campus tour reminder',
    actor_id: 'user-admin',
    action: 'SYNC_PUSH',
    data: '{"ip":"192.168.1.105"}',
    message: 'SMS provider rate limit created retry-ready sync event.',
    timestamp: '2026-05-27T17:46:00Z',
};
