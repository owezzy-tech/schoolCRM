import { SchoolAccessService } from 'app/core/school-access/school-access.service';
import {
    Capability,
    Membership,
} from 'app/core/school-access/school-access.types';
import { forkJoin, map, of, switchMap } from 'rxjs';

/** A department the viewer holds lesson capabilities in. */
export interface LessonScope {
    key: string;
    schoolID: string;
    departmentID: string;
    label: string;
    capabilities: Capability[];
}

const LESSON_CAPABILITIES: Capability[] = [
    'teach',
    'review_lessons',
    'approve_lessons',
];

/** Current server memberships own the available lesson departments. */
export function loadLessonScopes(access: SchoolAccessService) {
    return forkJoin([access.ownMemberships(), access.schools()]).pipe(
        switchMap(([memberships, schools]) => {
            const active = memberships.filter(
                (m) =>
                    m.active &&
                    m.departmentID &&
                    LESSON_CAPABILITIES.includes(m.capability)
            );
            const schoolIDs = [...new Set(active.map((m) => m.schoolID))];
            if (!schoolIDs.length) return of([] as LessonScope[]);
            return forkJoin(schoolIDs.map((id) => access.departments(id))).pipe(
                map((departments) =>
                    toScopes(active, schools, departments.flat())
                )
            );
        })
    );
}

function toScopes(
    memberships: Membership[],
    schools: { id: string; name: string }[],
    departments: { id: string; name: string }[]
): LessonScope[] {
    const byKey = new Map<string, LessonScope>();
    for (const membership of memberships) {
        const departmentID = membership.departmentID as string;
        const key = `${membership.schoolID}:${departmentID}`;
        const scope = byKey.get(key) ?? {
            key,
            schoolID: membership.schoolID,
            departmentID,
            label: `${schools.find((s) => s.id === membership.schoolID)?.name ?? 'School'} · ${
                departments.find((d) => d.id === departmentID)?.name ??
                'Department'
            }`,
            capabilities: [],
        };
        scope.capabilities.push(membership.capability);
        byKey.set(key, scope);
    }
    return [...byKey.values()].sort((a, b) => a.label.localeCompare(b.label));
}
