import { TestBed } from '@angular/core/testing';
import { AdmissionsService } from 'app/core/admissions/admissions.service';
import { Campaign } from 'app/core/admissions/admissions.types';
import { PaginatedResult } from 'app/core/api/json-api';
import { of, throwError } from 'rxjs';
import { beforeEach, describe, expect, it, vi } from 'vitest';

import { CampaignsComponent } from './campaigns.component';

describe('CampaignsComponent', () => {
    let queryCampaigns: ReturnType<typeof vi.fn>;

    beforeEach(async () => {
        queryCampaigns = vi.fn(() => of(campaignsResult([openHouseCampaign])));

        await TestBed.configureTestingModule({
            imports: [CampaignsComponent],
            providers: [
                {
                    provide: AdmissionsService,
                    useValue: {
                        queryCampaigns,
                    },
                },
            ],
        })
            .overrideComponent(CampaignsComponent, {
                set: {
                    template: '',
                },
            })
            .compileComponents();
    });

    it('loads campaigns from the admissions API using the live query contract', () => {
        TestBed.createComponent(CampaignsComponent);

        expect(queryCampaigns).toHaveBeenCalledWith({
            page: 1,
            rows: 50,
            orderBy: 'starts_at,DESC',
        });
    });

    it('renders campaign metrics from API data', () => {
        const fixture = TestBed.createComponent(CampaignsComponent);
        const component = fixture.componentInstance;

        expect(component.campaigns()).toEqual([openHouseCampaign]);
        expect(component.selectedCampaign()).toEqual(openHouseCampaign);
        expect(component.totalAudience()).toBe(2450);
        expect(component.activeCampaigns()).toBe(1);
        expect(component.averageOpenRate()).toBe(44);
    });

    it('surfaces JSON:API errors when campaign loading fails', () => {
        queryCampaigns.mockReturnValue(
            throwError(() => ({
                error: {
                    errors: [
                        {
                            detail: 'Campaign service unavailable.',
                        },
                    ],
                },
            }))
        );

        const fixture = TestBed.createComponent(CampaignsComponent);
        const component = fixture.componentInstance;

        expect(component.loading()).toBe(false);
        expect(component.errorMessage()).toBe('Campaign service unavailable.');
        expect(component.campaigns()).toEqual([]);
    });
});

function campaignsResult(items: Campaign[]): PaginatedResult<Campaign> {
    return {
        items,
        total: items.length,
        page: 1,
        rowsPerPage: 50,
    };
}

const openHouseCampaign: Campaign = {
    id: 'campaign-open-house',
    name: 'Nairobi Open House Invite',
    status: 'ACTIVE',
    channel: 'EMAIL',
    audienceName: 'Warm prospects near campus',
    templateName: 'Open house invitation',
    messagePreview: 'Join our admissions team for a guided campus experience.',
    segment: {
        applicationTypes: ['KUCCPS_PLACEMENT'],
        applicationStatuses: ['SUBMITTED'],
        academicTerms: ['September 2026'],
        programs: ['B.Sc. Computer Science'],
        eventAttendance: 'ANY',
        leadScoreBands: ['WARM'],
        recruiters: ['Owen Adirah'],
        territories: ['Nairobi'],
    },
    metrics: {
        audienceSize: 2450,
        sent: 2450,
        delivered: 2402,
        opened: 1081,
        clicked: 294,
        bounced: 48,
        replied: 0,
    },
    startsAt: '2026-05-01T09:00:00Z',
    endsAt: '2026-06-06T09:00:00Z',
    auditTrail: [
        {
            id: 'campaign-audit-1',
            campaignID: 'campaign-open-house',
            actorName: 'Owen Adirah',
            action: 'Activated campaign',
            occurredAt: '2026-05-01T09:15:00Z',
            dateCreated: '2026-05-01T09:15:00Z',
        },
    ],
    dateCreated: '2026-04-24T09:00:00Z',
    dateUpdated: '2026-05-18T09:00:00Z',
};
