import {
    ApplicationStatus,
    ApplicationType,
    LeadScoreBand,
} from 'app/core/admissions/admissions.types';

export const CAMPAIGN_STATUSES = [
    'DRAFT',
    'ACTIVE',
    'PAUSED',
    'COMPLETED',
] as const;

export const CAMPAIGN_CHANNELS = ['EMAIL', 'SMS'] as const;

export type CampaignStatus = (typeof CAMPAIGN_STATUSES)[number];
export type CampaignChannel = (typeof CAMPAIGN_CHANNELS)[number];

export interface CampaignSegmentFilter {
    applicationTypes: ApplicationType[];
    applicationStatuses: ApplicationStatus[];
    academicTerms: string[];
    programs: string[];
    eventAttendance: 'ANY' | 'REGISTERED' | 'ATTENDED' | 'NO_SHOW';
    leadScoreBands: LeadScoreBand[];
    recruiters: string[];
    territories: string[];
}

export interface CampaignMetrics {
    audienceSize: number;
    sent: number;
    delivered: number;
    opened: number;
    clicked: number;
    bounced: number;
    replied: number;
}

export interface CampaignAuditEvent {
    actor: string;
    action: string;
    timestamp: string;
}

export interface Campaign {
    id: string;
    name: string;
    status: CampaignStatus;
    channel: CampaignChannel;
    schedule: string;
    audienceName: string;
    templateName: string;
    messagePreview: string;
    segment: CampaignSegmentFilter;
    metrics: CampaignMetrics;
    auditTrail: CampaignAuditEvent[];
    dateCreated: string;
    dateUpdated: string;
}
