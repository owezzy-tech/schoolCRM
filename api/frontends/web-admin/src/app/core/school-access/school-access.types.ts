export type Capability =
    | 'manage_members'
    | 'teach'
    | 'review_lessons'
    | 'approve_lessons';

export const CAPABILITY_LABELS: Record<Capability, string> = {
    manage_members: 'Manages members',
    teach: 'Teaches',
    review_lessons: 'Reviews lessons (HOD)',
    approve_lessons: 'Approves lessons (Dean)',
};

export interface School {
    id: string;
    name: string;
    dateCreated: string;
}

export interface Department {
    id: string;
    schoolID: string;
    name: string;
    dateCreated: string;
}

export interface Membership {
    id: string;
    schoolID: string;
    departmentID?: string;
    userID: string;
    capability: Capability;
    active: boolean;
    dateUpdated: string;
}

export interface GrantRequest {
    userID: string;
    capability: Capability;
    departmentID?: string;
}

/** Active capabilities the memberships grant in one school department. */
export function capabilitiesFor(
    memberships: Membership[],
    schoolID: string,
    departmentID: string
): Capability[] {
    return memberships
        .filter(
            (membership) =>
                membership.active &&
                membership.schoolID === schoolID &&
                membership.departmentID === departmentID
        )
        .map((membership) => membership.capability);
}
