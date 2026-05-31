/* eslint-disable */
import { DateTime } from 'luxon';

/* Get the current instant */
const now = DateTime.now();

export const notifications = [
    {
        id: '493190c9-5b61-4912-afe5-78c21f1044d7',
        icon: 'heroicons_mini:document-check',
        title: 'Application submitted',
        description:
            '<strong>Sofia Martinez</strong> submitted application <em>APP-3024</em> for Fall 2026 admission.',
        time: now.minus({ minute: 12 }).toISO(),
        read: false,
        link: '/applications/APP-3024',
        useRouter: true,
    },
    {
        id: '6e3e97e5-effc-4fb7-b730-52a151f0b641',
        icon: 'heroicons_mini:user-plus',
        title: 'Review assignment',
        description:
            '<strong>Admissions Review</strong> assigned Liam Chen\'s application to your queue for first review.',
        time: now.minus({ minute: 38 }).toISO(),
        read: false,
        link: '/applications',
        useRouter: true,
    },
    {
        id: 'b91ccb58-b06c-413b-b389-87010e03a120',
        icon: 'heroicons_mini:exclamation-triangle',
        title: 'Document needs attention',
        description:
            '<strong>James Okoro</strong> has a rejected transcript. Applicant notification is ready for follow-up.',
        time: now.minus({ hour: 2 }).toISO(),
        read: false,
        link: '/applications/APP-3023',
        useRouter: true,
    },
    {
        id: '541416c9-84a7-408a-8d74-27a43c38d797',
        icon: 'heroicons_mini:check-circle',
        title: 'Transcript verified',
        description:
            'Your official transcript has been verified. The checklist for <em>APP-3024</em> is now 80% complete.',
        time: now.minus({ hour: 4 }).toISO(),
        read: false,
        link: '/portal/status',
        useRouter: true,
    },
    {
        id: 'ef7b95a7-8e8b-4616-9619-130d9533add9',
        icon: 'heroicons_mini:clock',
        title: 'Decision pending',
        description:
            '<strong>Amara Ndlovu</strong> is ready for final decision review after all required documents were accepted.',
        time: now.minus({ hour: 7 }).toISO(),
        read: true,
        link: '/applications',
        useRouter: true,
    },
    {
        id: 'eb8aa470-635e-461d-88e1-23d9ea2a5665',
        icon: 'heroicons_mini:sparkles',
        title: 'Decision posted',
        description:
            'Congratulations! Your admission decision for <em>APP-3018</em> has been posted in the applicant portal.',
        time: now.minus({ hour: 9 }).toISO(),
        read: true,
        link: '/portal/status',
        useRouter: true,
    },
    {
        id: 'b85c2338-cc98-4140-bbf8-c226ce4e395e',
        icon: 'heroicons_mini:arrow-path',
        title: 'SIS sync completed',
        description:
            'The nightly SIS sync matched <strong>27 application updates</strong> and flagged 3 records for review.',
        time: now.minus({ day: 1 }).toISO(),
        read: true,
        link: '/reports',
        useRouter: true,
    },
    {
        id: '8f8e1bf9-4661-4939-9e43-390957b60f42',
        icon: 'heroicons_mini:users',
        title: 'Duplicate detected',
        description:
            '<strong>Sofia Martinez</strong> has a 97% duplicate match with an existing constituent record.',
        time: now.minus({ day: 2 }).toISO(),
        read: true,
        link: '/duplicates',
        useRouter: true,
    },
    {
        id: '30af917b-7a6a-45d1-822f-9e7ad7f8bf69',
        icon: 'heroicons_mini:megaphone',
        title: 'Campaign follow-up sent',
        description:
            '<strong>Missing Documents Reminder</strong> reached 342 applicants with pending checklist items.',
        time: now.minus({ day: 3 }).toISO(),
        read: true,
        link: '/campaigns',
        useRouter: true,
    },
    {
        id: '1c55dc97-d902-46b6-99b2-12812b480f8c',
        icon: 'heroicons_mini:calendar-days',
        title: 'Event registration milestone',
        description:
            '<strong>Spring Open Day</strong> reached 845 registrations. Capacity review is recommended.',
        time: now.minus({ day: 4 }).toISO(),
        read: true,
        link: '/events',
        useRouter: true,
    },
];
