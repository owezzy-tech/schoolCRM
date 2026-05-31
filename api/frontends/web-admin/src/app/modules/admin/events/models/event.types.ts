export type EventType =
    | 'open-day'
    | 'webinar'
    | 'info-session'
    | 'campus-tour'
    | 'fair';

export type EventStatus =
    | 'draft'
    | 'upcoming'
    | 'live'
    | 'completed'
    | 'cancelled';

export interface EventRegistration {
    id: string;
    constituentId: string;
    constituentName: string;
    status: 'registered' | 'checked-in' | 'cancelled';
    registeredAt: string;
}

export interface EventItem {
    id: string;
    title: string;
    type: EventType;
    status: EventStatus;
    description: string;
    start: string;
    end: string;
    location: string;
    isVirtual: boolean;
    capacity: number;
    registeredCount: number;
    checkedInCount: number;
}
