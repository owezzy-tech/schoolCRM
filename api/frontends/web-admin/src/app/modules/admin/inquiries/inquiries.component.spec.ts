import { TestBed } from '@angular/core/testing';
import { AdmissionsService } from 'app/core/admissions/admissions.service';
import { Inquiry } from 'app/core/admissions/admissions.types';
import { PaginatedResult } from 'app/core/api/json-api';
import { of, ReplaySubject, throwError } from 'rxjs';
import { beforeEach, describe, expect, it, vi } from 'vitest';

import { InquiriesComponent } from './inquiries.component';

describe('InquiriesComponent', () => {
    let inquiriesSubject: ReplaySubject<PaginatedResult<Inquiry>>;
    let queryInquiries: ReturnType<typeof vi.fn>;

    beforeEach(async () => {
        inquiriesSubject = new ReplaySubject<PaginatedResult<Inquiry>>(1);
        queryInquiries = vi.fn(() => of(inquiriesResult([achiengInquiry])));

        await TestBed.configureTestingModule({
            imports: [InquiriesComponent],
            providers: [
                {
                    provide: AdmissionsService,
                    useValue: {
                        inquiries$: inquiriesSubject.asObservable(),
                        queryInquiries,
                    },
                },
            ],
        })
            .overrideComponent(InquiriesComponent, {
                set: {
                    template: '',
                },
            })
            .compileComponents();
    });

    it('loads inquiries from the admissions API', () => {
        TestBed.createComponent(InquiriesComponent);

        expect(queryInquiries).toHaveBeenCalledWith({
            page: 1,
            rows: 50,
            orderBy: 'date_created,DESC',
        });
    });

    it('derives inquiry table rows from the service stream', () => {
        const fixture = TestBed.createComponent(InquiriesComponent);
        const component = fixture.componentInstance;

        inquiriesSubject.next(inquiriesResult([achiengInquiry]));

        expect(component.inquiries()[0]).toMatchObject({
            id: 'inquiry-achieng',
            from: 'Achieng Otieno',
            source: 'Portal',
            priority: 'High',
            status: 'New',
        });
    });

    it('surfaces JSON:API errors when inquiry loading fails', () => {
        queryInquiries.mockReturnValue(
            throwError(() => ({
                error: {
                    errors: [
                        {
                            detail: 'Inquiry service unavailable.',
                        },
                    ],
                },
            }))
        );

        const fixture = TestBed.createComponent(InquiriesComponent);
        const component = fixture.componentInstance;

        expect(component.loading()).toBe(false);
        expect(component.errorMessage()).toBe('Inquiry service unavailable.');
    });
});

function inquiriesResult(items: Inquiry[]): PaginatedResult<Inquiry> {
    return {
        items,
        total: items.length,
        page: 1,
        rowsPerPage: 50,
    };
}

const achiengInquiry: Inquiry = {
    id: 'inquiry-achieng',
    constituentID: 'constituent-achieng',
    firstName: 'Achieng',
    lastName: 'Otieno',
    dateOfBirth: '2003-02-21T00:00:00Z',
    primaryEmail: 'achieng.otieno@example.com',
    primaryPhone: '+254711100003',
    source: 'portal',
    message: 'Interested in September intake scholarships.',
    status: 'NEW',
    dateCreated: '2026-05-18T11:00:00Z',
    dateUpdated: '2026-05-18T11:00:00Z',
};
