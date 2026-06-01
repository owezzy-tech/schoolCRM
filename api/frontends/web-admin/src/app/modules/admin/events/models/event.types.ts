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
    constituentId?: string;
    constituentName: string;
    email: string;
    phone?: string;
    status: 'registered' | 'checked-in' | 'cancelled';
    registeredAt: string;
    matchStatus: 'matched' | 'new-prospect' | 'needs-review';
    source: 'portal' | 'staff' | 'campaign';
    checkedInAt?: string;
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
    registrationDeadline: string;
    autoConfirmationEnabled: boolean;
    autoReminderEnabled: boolean;
    registrations: EventRegistration[];
}
