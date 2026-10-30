import {
    AdmissionsEvent,
    AdmissionsEventRegistration,
    EventStatus,
    EventType,
} from 'app/core/admissions/admissions.types';

export type { EventStatus, EventType };
export type EventRegistration = AdmissionsEventRegistration;
export type EventItem = Omit<AdmissionsEvent, 'dateCreated' | 'dateUpdated'> & {
    dateCreated?: string;
    dateUpdated?: string;
};
