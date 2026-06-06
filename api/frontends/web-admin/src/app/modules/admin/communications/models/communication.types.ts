export type {
    CommunicationChannel,
    CommunicationDirection,
    CommunicationRecord,
    CommunicationStatus,
} from 'app/core/admissions/admissions.types';

export interface CommunicationSummary {
    label: string;
    value: string;
    icon: string;
    tone: 'blue' | 'green' | 'amber' | 'red';
}
