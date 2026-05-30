export const SCHOOL_ROLES = [
    'SUPER_ADMIN',
    'SCHOOL_ADMIN',
    'TEACHER',
    'STUDENT',
    'PARENT',
] as const;

export type SchoolRole = (typeof SCHOOL_ROLES)[number];

export interface AdminUser {
    id: string;
    name: string;
    email: string;
    roles: SchoolRole[];
    department: string;
    enabled: boolean;
    dateCreated: string;
    dateUpdated: string;
}

export interface UserQuery {
    page?: number;
    rows?: number;
    orderBy?: string;
    name?: string;
    email?: string;
}

export interface UpdateUserRolesRequest {
    roles: SchoolRole[];
}

export interface UpdateUserRequest {
    name?: string;
    email?: string;
    department?: string;
    password?: string;
    passwordConfirm?: string;
    enabled?: boolean;
}
