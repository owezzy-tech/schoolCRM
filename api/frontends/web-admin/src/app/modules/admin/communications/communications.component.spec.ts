import { TestBed } from '@angular/core/testing';
import { AdmissionsService } from 'app/core/admissions/admissions.service';
import { CommunicationRecord } from 'app/core/admissions/admissions.types';
import { PaginatedResult } from 'app/core/api/json-api';
import { of, throwError } from 'rxjs';
import { beforeEach, describe, expect, it, vi } from 'vitest';

import { CommunicationsComponent } from './communications.component';

describe('CommunicationsComponent', () => {
    let queryCommunications: ReturnType<typeof vi.fn>;

    beforeEach(async () => {
        queryCommunications = vi.fn(() =>
            of(communicationsResult([emailCommunication, smsCommunication]))
        );

        await TestBed.configureTestingModule({
            imports: [CommunicationsComponent],
            providers: [
                {
                    provide: AdmissionsService,
                    useValue: {
                        queryCommunications,
                    },
                },
            ],
        })
            .overrideComponent(CommunicationsComponent, {
                set: {
                    template: '',
                },
            })
            .compileComponents();
    });

    it('loads communications from the admissions API using the live query contract', () => {
        TestBed.createComponent(CommunicationsComponent);

        expect(queryCommunications).toHaveBeenCalledWith({
            page: 1,
            rows: 50,
            orderBy: 'occurred_at,DESC',
            channel: undefined,
            direction: undefined,
            status: undefined,
        });
    });

    it('reloads API data with selected filters', () => {
        const fixture = TestBed.createComponent(CommunicationsComponent);
        const component = fixture.componentInstance;

        component.setChannel('SMS');

        expect(queryCommunications).toHaveBeenLastCalledWith({
            page: 1,
            rows: 50,
            orderBy: 'occurred_at,DESC',
            channel: 'SMS',
            direction: undefined,
            status: undefined,
        });
    });

    it('derives dashboard summaries from API communication records', () => {
        const fixture = TestBed.createComponent(CommunicationsComponent);
        const component = fixture.componentInstance;

        expect(component.records()).toEqual([emailCommunication, smsCommunication]);
        expect(component.summaries()).toEqual([
            expect.objectContaining({ label: 'Total communications', value: '2' }),
            expect.objectContaining({ label: 'Provider-tracked', value: '2' }),
            expect.objectContaining({ label: 'Linked applications', value: '1' }),
            expect.objectContaining({ label: 'Needs attention', value: '1' }),
        ]);
    });

    it('surfaces JSON:API errors when communication loading fails', () => {
        queryCommunications.mockReturnValue(
            throwError(() => ({
                error: {
                    errors: [
                        {
                            detail: 'Communication service unavailable.',
                        },
                    ],
                },
            }))
        );

        const fixture = TestBed.createComponent(CommunicationsComponent);
        const component = fixture.componentInstance;

        expect(component.loading()).toBe(false);
        expect(component.errorMessage()).toBe(
            'Communication service unavailable.'
        );
        expect(component.records()).toEqual([]);
    });
});

function communicationsResult(
    items: CommunicationRecord[]
): PaginatedResult<CommunicationRecord> {
    return {
        items,
        total: items.length,
        page: 1,
        rowsPerPage: 50,
    };
}

const emailCommunication: CommunicationRecord = {
    id: 'communication-email',
    externalMessageID: 'MSG-8492',
    channel: 'EMAIL',
    direction: 'OUTBOUND',
    constituentID: 'constituent-sofia',
    applicationID: 'application-3024',
    campaignID: 'campaign-open-house',
    recipientSender: 'Sofia Martinez',
    recipientInitials: 'SM',
    subject: 'Application status update',
    preview: 'Your application has moved into review.',
    status: 'OPENED',
    provider: 'SendGrid',
    ownerName: 'Owen Adirah',
    occurredAt: '2026-06-06T09:00:00Z',
    dateCreated: '2026-06-06T09:00:00Z',
    dateUpdated: '2026-06-06T09:00:00Z',
};

const smsCommunication: CommunicationRecord = {
    id: 'communication-sms',
    externalMessageID: 'MSG-8491',
    channel: 'SMS',
    direction: 'OUTBOUND',
    constituentID: 'constituent-james',
    recipientSender: 'James Okoro',
    recipientInitials: 'JO',
    subject: 'Missing transcript reminder',
    preview: 'Please upload your corrected transcript.',
    status: 'FAILED',
    provider: 'Celcom Africa',
    ownerName: 'Admissions Automation',
    outcome: 'Retry scheduled',
    occurredAt: '2026-06-06T08:00:00Z',
    dateCreated: '2026-06-06T08:00:00Z',
    dateUpdated: '2026-06-06T08:00:00Z',
};
