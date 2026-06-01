/* eslint-disable */
import { FuseNavigationItem } from '@fuse/components/navigation';

export const defaultNavigation: FuseNavigationItem[] = [
    {
        id: 'workspace',
        title: 'Workspace',
        type: 'group',
        children: [
            {
                id: 'workspace.dashboard',
                title: 'Dashboard',
                type: 'basic',
                icon: 'heroicons_outline:home',
                link: '/dashboard',
            },
            {
                id: 'workspace.constituents',
                title: 'Constituents',
                type: 'basic',
                icon: 'heroicons_outline:user-group',
                link: '/constituents',
            },
            {
                id: 'workspace.duplicates',
                title: 'Duplicates',
                type: 'basic',
                icon: 'heroicons_outline:document-duplicate',
                link: '/duplicates',
            },
            {
                id: 'workspace.inquiries',
                title: 'Inquiries',
                type: 'basic',
                icon: 'heroicons_outline:inbox',
                link: '/inquiries',
            },
        ],
    },
    {
        id: 'admissions',
        title: 'Admissions',
        type: 'group',
        children: [
            {
                id: 'admissions.applications',
                title: 'Applications',
                type: 'basic',
                icon: 'heroicons_outline:document-text',
                link: '/applications',
            },
            {
                id: 'admissions.reviews',
                title: 'Review Workspace',
                type: 'basic',
                icon: 'heroicons_outline:clipboard-document-check',
                link: '/applications/review',
            },
            {
                id: 'admissions.leads',
                title: 'Leads',
                type: 'basic',
                icon: 'heroicons_outline:queue-list',
                link: '/leads',
            },
        ],
    },
    {
        id: 'engagement',
        title: 'Engagement',
        type: 'group',
        children: [
            {
                id: 'engagement.communications',
                title: 'Communications',
                type: 'basic',
                icon: 'heroicons_outline:chat-bubble-left-right',
                link: '/communications',
            },
            {
                id: 'engagement.campaigns',
                title: 'Campaigns',
                type: 'basic',
                icon: 'heroicons_outline:megaphone',
                link: '/campaigns',
            },
            {
                id: 'engagement.events',
                title: 'Events',
                type: 'basic',
                icon: 'heroicons_outline:calendar-days',
                link: '/events',
            },
        ],
    },
    {
        id: 'insights',
        title: 'Insights',
        type: 'group',
        children: [
            {
                id: 'insights.reports',
                title: 'Reports',
                type: 'basic',
                icon: 'heroicons_outline:chart-bar',
                link: '/reports',
            },
        ],
    },
    {
        id: 'admin',
        title: 'Admin',
        type: 'group',
        children: [
            {
                id: 'admin.users',
                title: 'Users & Roles',
                type: 'basic',
                icon: 'heroicons_outline:users',
                link: '/users',
            },
            {
                id: 'admin.audit',
                title: 'Audit Log',
                type: 'basic',
                icon: 'heroicons_outline:clipboard-document-list',
                link: '/audit',
            },
            {
                id: 'admin.settings',
                title: 'Settings',
                type: 'basic',
                icon: 'heroicons_outline:cog-6-tooth',
                link: '/settings',
            },
        ],
    },
    {
        id: 'portal',
        title: 'Applicant Portal',
        type: 'basic',
        icon: 'heroicons_outline:globe-alt',
        link: '/portal',
    },
];

export const compactNavigation: FuseNavigationItem[] = defaultNavigation;
export const futuristicNavigation: FuseNavigationItem[] = defaultNavigation;
export const horizontalNavigation: FuseNavigationItem[] = defaultNavigation;
