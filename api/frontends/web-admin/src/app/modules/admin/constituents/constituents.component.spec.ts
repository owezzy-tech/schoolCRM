import { TestBed } from '@angular/core/testing';
import { AdmissionsService } from 'app/core/admissions/admissions.service';
import { Constituent } from 'app/core/admissions/admissions.types';
import { PaginatedResult } from 'app/core/api/json-api';
import { of, ReplaySubject, throwError } from 'rxjs';
import { beforeEach, describe, expect, it, vi } from 'vitest';

import { ConstituentsComponent } from './constituents.component';

describe('ConstituentsComponent', () => {
    let constituentsSubject: ReplaySubject<PaginatedResult<Constituent>>;
    let queryConstituents: ReturnType<typeof vi.fn>;

    beforeEach(async () => {
        constituentsSubject = new ReplaySubject<PaginatedResult<Constituent>>(1);
        queryConstituents = vi.fn(() =>
            of(constituentsResult([sofiaConstituent]))
        );

        await TestBed.configureTestingModule({
            imports: [ConstituentsComponent],
            providers: [
                {
                    provide: AdmissionsService,
                    useValue: {
                        constituents$: constituentsSubject.asObservable(),
                        queryConstituents,
                    },
                },
            ],
        })
            .overrideComponent(ConstituentsComponent, {
                set: {
                    template: '',
                },
            })
            .compileComponents();
    });

    it('loads constituents from the admissions API', () => {
        TestBed.createComponent(ConstituentsComponent);

        expect(queryConstituents).toHaveBeenCalledWith({
            page: 1,
            rows: 50,
            orderBy: 'date_created,DESC',
        });
    });

    it('derives constituent table rows from the service stream', () => {
        const fixture = TestBed.createComponent(ConstituentsComponent);
        const component = fixture.componentInstance;

        constituentsSubject.next(constituentsResult([sofiaConstituent]));

        expect(component.constituents()[0]).toMatchObject({
            name: 'Sofia Martinez',
            email: 'sofia.martinez@example.com',
            stage: 'Applicant',
        });
    });

    it('surfaces JSON:API errors when constituent loading fails', () => {
        queryConstituents.mockReturnValue(
            throwError(() => ({
                error: {
                    errors: [
                        {
                            detail: 'Constituent service unavailable.',
                        },
                    ],
                },
            }))
        );

        const fixture = TestBed.createComponent(ConstituentsComponent);
        const component = fixture.componentInstance;

        expect(component.loading()).toBe(false);
        expect(component.errorMessage()).toBe(
            'Constituent service unavailable.'
        );
    });
});

function constituentsResult(
    items: Constituent[]
): PaginatedResult<Constituent> {
    return {
        items,
        total: items.length,
        page: 1,
        rowsPerPage: 50,
    };
}

const sofiaConstituent: Constituent = {
    id: 'constituent-sofia',
    firstName: 'Sofia',
    lastName: 'Martinez',
    dateOfBirth: '2002-04-12T00:00:00Z',
    primaryEmail: 'sofia.martinez@example.com',
    primaryPhone: '+254711100001',
    externalSISID: 'ADM-CON-921',
    lifecycleStage: 'APPLICANT',
    duplicateStatus: 'ACTIVE',
    notificationPreferences: {
        smsOptIn: true,
        whatsAppOptIn: true,
        emailOptIn: true,
        priority: ['EMAIL', 'SMS'],
    },
    sisSyncedAt: '2026-05-28T08:00:00Z',
    dateCreated: '2026-05-20T09:00:00Z',
    dateUpdated: '2026-05-28T08:00:00Z',
};
