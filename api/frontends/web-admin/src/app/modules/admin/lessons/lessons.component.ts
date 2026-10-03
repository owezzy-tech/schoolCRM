import { DatePipe } from '@angular/common';
import {
    ChangeDetectionStrategy,
    Component,
    computed,
    inject,
    signal,
} from '@angular/core';
import { toSignal } from '@angular/core/rxjs-interop';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatButtonToggleModule } from '@angular/material/button-toggle';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatSelectModule } from '@angular/material/select';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { jsonApiErrorMessage, JsonApiErrorResponse } from 'app/core/api/json-api';
import { LessonsService } from 'app/core/lessons/lessons.service';
import { LessonPlan } from 'app/core/lessons/lessons.types';
import { SchoolAccessService } from 'app/core/school-access/school-access.service';
import {
    Capability,
    Membership,
} from 'app/core/school-access/school-access.types';
import { UserService } from 'app/core/user/user.service';
import { catchError, EMPTY, forkJoin, map, of, switchMap } from 'rxjs';

import { LessonStatusComponent } from './shared/lesson-status.component';

/** A department the viewer holds lesson capabilities in. */
export interface LessonScope {
    key: string;
    schoolID: string;
    departmentID: string;
    label: string;
    capabilities: Capability[];
}

export type LessonFilter = 'all' | 'mine' | 'review' | 'approval';

const LESSON_CAPABILITIES: Capability[] = [
    'teach',
    'review_lessons',
    'approve_lessons',
];

@Component({
    selector: 'app-lessons',
    standalone: true,
    imports: [
        DatePipe,
        FormsModule,
        MatButtonModule,
        MatButtonToggleModule,
        MatFormFieldModule,
        MatIconModule,
        MatSelectModule,
        RouterLink,
        LessonStatusComponent,
    ],
    changeDetection: ChangeDetectionStrategy.OnPush,
    templateUrl: './lessons.component.html',
})
export class LessonsComponent {
    private readonly lessons = inject(LessonsService);
    private readonly schoolAccess = inject(SchoolAccessService);
    private readonly route = inject(ActivatedRoute);
    private readonly router = inject(Router);

    readonly viewer = toSignal(inject(UserService).user$);
    readonly scopes = signal<LessonScope[] | null>(null);
    readonly plans = signal<LessonPlan[] | null>(null);
    readonly errorMessage = signal<string | null>(null);
    readonly filter = signal<LessonFilter>('all');
    private readonly selectedKey = toSignal(
        this.route.queryParamMap.pipe(
            map((params) =>
                params.get('school') && params.get('department')
                    ? `${params.get('school')}:${params.get('department')}`
                    : null
            )
        ),
        { initialValue: null }
    );

    readonly scope = computed(() => {
        const scopes = this.scopes() ?? [];
        return scopes.find((s) => s.key === this.selectedKey()) ?? scopes[0] ?? null;
    });
    readonly can = (capability: Capability) =>
        this.scope()?.capabilities.includes(capability) ?? false;

    readonly visiblePlans = computed(() => {
        const plans = this.plans() ?? [];
        const viewerID = this.viewer()?.id;
        switch (this.filter()) {
            case 'mine':
                return plans.filter((plan) => plan.authorID === viewerID);
            case 'review':
                return plans.filter((plan) => plan.status === 'hod_review');
            case 'approval':
                return plans.filter((plan) => plan.status === 'dean_approval');
            default:
                return plans;
        }
    });

    constructor() {
        this.loadScopes();
    }

    selectScope(key: string): void {
        const scope = this.scopes()?.find((s) => s.key === key);
        if (!scope) return;
        this.router.navigate([], {
            relativeTo: this.route,
            queryParams: { school: scope.schoolID, department: scope.departmentID },
        });
        this.loadPlans(scope);
    }

    loadPlans(scope: LessonScope | null = this.scope()): void {
        if (!scope) return;
        this.plans.set(null);
        this.errorMessage.set(null);
        this.lessons
            .plans(scope.schoolID, scope.departmentID)
            .pipe(catchError((error) => this.fail(error, 'Unable to load lesson plans.')))
            .subscribe((plans) => this.plans.set(plans));
    }

    private loadScopes(): void {
        forkJoin([this.schoolAccess.ownMemberships(), this.schoolAccess.schools()])
            .pipe(
                switchMap(([memberships, schools]) => {
                    const lessonMemberships = memberships.filter(
                        (m) => m.departmentID && LESSON_CAPABILITIES.includes(m.capability)
                    );
                    const schoolIDs = [...new Set(lessonMemberships.map((m) => m.schoolID))];
                    if (!schoolIDs.length) return of([]);
                    return forkJoin(schoolIDs.map((id) => this.schoolAccess.departments(id))).pipe(
                        map((departments) =>
                            toScopes(lessonMemberships, schools, departments.flat())
                        )
                    );
                }),
                catchError((error) => this.fail(error, 'Unable to load your departments.'))
            )
            .subscribe((scopes) => {
                this.scopes.set(scopes);
                this.loadPlans();
            });
    }

    private fail(error: unknown, fallback: string) {
        this.errorMessage.set(jsonApiErrorMessage(error as JsonApiErrorResponse, fallback));
        this.plans.set([]);
        return EMPTY;
    }
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
                departments.find((d) => d.id === departmentID)?.name ?? 'Department'
            }`,
            capabilities: [],
        };
        scope.capabilities.push(membership.capability);
        byKey.set(key, scope);
    }
    return [...byKey.values()].sort((a, b) => a.label.localeCompare(b.label));
}
