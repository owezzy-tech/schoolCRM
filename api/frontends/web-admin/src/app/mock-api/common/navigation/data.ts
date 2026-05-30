/* eslint-disable */
import { FuseNavigationItem } from '@fuse/components/navigation';

export const defaultNavigation: FuseNavigationItem[] = [
    {
        id: 'example',
        title: 'Example',
        type: 'basic',
        icon: 'heroicons_outline:chart-pie',
        link: '/example',
    },
    {
        id: 'calendar',
        title: 'Calendar',
        type: 'basic',
        icon: 'heroicons_outline:calendar',
        link: '/calendar',
    },
    {
        id: 'admin',
        title: 'Administration',
        type: 'collapsable',
        icon: 'heroicons_outline:cog-6-tooth',
        children: [
            {
                id: 'admin.users',
                title: 'User Roles',
                type: 'basic',
                icon: 'heroicons_outline:users',
                link: '/admin/users',
            },
        ],
    },
    {
        id: 'admissions',
        title: 'Admissions',
        type: 'collapsable',
        icon: 'heroicons_outline:academic-cap',
        children: [
            {
                id: 'admissions.leads',
                title: 'Leads',
                type: 'basic',
                icon: 'heroicons_outline:queue-list',
                link: '/admin/admissions/leads',
            },
            {
                id: 'admissions.documents',
                title: 'Documents',
                type: 'basic',
                icon: 'heroicons_outline:document-check',
                link: '/admin/admissions/documents',
            },
            {
                id: 'admissions.settings',
                title: 'Settings',
                type: 'basic',
                icon: 'heroicons_outline:adjustments-horizontal',
                link: '/admin/admissions/settings',
            },
        ],
    },
];
export const compactNavigation: FuseNavigationItem[] = [
    {
        id: 'example',
        title: 'Example',
        type: 'basic',
        icon: 'heroicons_outline:chart-pie',
        link: '/example',
    },
    {
        id: 'calendar',
        title: 'Calendar',
        type: 'basic',
        icon: 'heroicons_outline:calendar',
        link: '/calendar',
    },
    {
        id: 'admin.users',
        title: 'User Roles',
        type: 'basic',
        icon: 'heroicons_outline:users',
        link: '/admin/users',
    },
    {
        id: 'admissions.leads',
        title: 'Leads',
        type: 'basic',
        icon: 'heroicons_outline:queue-list',
        link: '/admin/admissions/leads',
    },
    {
        id: 'admissions.documents',
        title: 'Documents',
        type: 'basic',
        icon: 'heroicons_outline:document-check',
        link: '/admin/admissions/documents',
    },
    {
        id: 'admissions.settings',
        title: 'Settings',
        type: 'basic',
        icon: 'heroicons_outline:adjustments-horizontal',
        link: '/admin/admissions/settings',
    },
];
export const futuristicNavigation: FuseNavigationItem[] = [
    {
        id: 'example',
        title: 'Example',
        type: 'basic',
        icon: 'heroicons_outline:chart-pie',
        link: '/example',
    },
    {
        id: 'calendar',
        title: 'Calendar',
        type: 'basic',
        icon: 'heroicons_outline:calendar',
        link: '/calendar',
    },
    {
        id: 'admin',
        title: 'Administration',
        type: 'collapsable',
        icon: 'heroicons_outline:cog-6-tooth',
        children: [
            {
                id: 'admin.users',
                title: 'User Roles',
                type: 'basic',
                icon: 'heroicons_outline:users',
                link: '/admin/users',
            },
        ],
    },
    {
        id: 'admissions',
        title: 'Admissions',
        type: 'collapsable',
        icon: 'heroicons_outline:academic-cap',
        children: [
            {
                id: 'admissions.leads',
                title: 'Leads',
                type: 'basic',
                icon: 'heroicons_outline:queue-list',
                link: '/admin/admissions/leads',
            },
            {
                id: 'admissions.documents',
                title: 'Documents',
                type: 'basic',
                icon: 'heroicons_outline:document-check',
                link: '/admin/admissions/documents',
            },
            {
                id: 'admissions.settings',
                title: 'Settings',
                type: 'basic',
                icon: 'heroicons_outline:adjustments-horizontal',
                link: '/admin/admissions/settings',
            },
        ],
    },
];
export const horizontalNavigation: FuseNavigationItem[] = [
    {
        id: 'example',
        title: 'Example',
        type: 'basic',
        icon: 'heroicons_outline:chart-pie',
        link: '/example',
    },
    {
        id: 'calendar',
        title: 'Calendar',
        type: 'basic',
        icon: 'heroicons_outline:calendar',
        link: '/calendar',
    },
    {
        id: 'admin.users',
        title: 'User Roles',
        type: 'basic',
        icon: 'heroicons_outline:users',
        link: '/admin/users',
    },
    {
        id: 'admissions.leads',
        title: 'Leads',
        type: 'basic',
        icon: 'heroicons_outline:queue-list',
        link: '/admin/admissions/leads',
    },
    {
        id: 'admissions.documents',
        title: 'Documents',
        type: 'basic',
        icon: 'heroicons_outline:document-check',
        link: '/admin/admissions/documents',
    },
    {
        id: 'admissions.settings',
        title: 'Settings',
        type: 'basic',
        icon: 'heroicons_outline:adjustments-horizontal',
        link: '/admin/admissions/settings',
    },
];
