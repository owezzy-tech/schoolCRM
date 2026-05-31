export const COMMUNICATION_CHANNELS = [
    'EMAIL',
    'SMS',
    'PHONE_CALL',
    'NOTIFICATION',
] as const;

export const COMMUNICATION_DIRECTIONS = ['INBOUND', 'OUTBOUND'] as const;

export const COMMUNICATION_STATUSES = [
    'QUEUED',
    'SENT',
    'DELIVERED',
    'FAILED',
    'OPENED',
    'BOUNCED',
    'REPLIED',
    'LOGGED',
] as const;

export type CommunicationChannel = (typeof COMMUNICATION_CHANNELS)[number];
export type CommunicationDirection = (typeof COMMUNICATION_DIRECTIONS)[number];
export type CommunicationStatus = (typeof COMMUNICATION_STATUSES)[number];

export interface CommunicationRecord {
    id: string;
    channel: CommunicationChannel;
    direction: CommunicationDirection;
    recipientSender: string;
    recipientInitials: string;
    constituentId: string;
    constituentName: string;
    applicationId?: string;
    campaignId?: string;
    subject: string;
    preview: string;
    status: CommunicationStatus;
    date: string;
    provider?: string;
    owner: string;
    outcome?: string;
    duration?: string;
}

export interface CommunicationSummary {
    label: string;
    value: string;
    icon: string;
    tone: 'blue' | 'green' | 'amber' | 'red';
}
