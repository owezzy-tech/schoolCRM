import { DatePipe } from '@angular/common';
import { HttpErrorResponse } from '@angular/common/http';
import {
    ChangeDetectionStrategy,
    Component,
    computed,
    inject,
    signal,
} from '@angular/core';
import { toSignal } from '@angular/core/rxjs-interop';
import { MatButtonModule } from '@angular/material/button';
import { ActivatedRoute, RouterLink } from '@angular/router';
import {
    jsonApiErrorMessage,
    JsonApiErrorResponse,
} from 'app/core/api/json-api';
import { LessonsService } from 'app/core/lessons/lessons.service';
import { LessonPlan, LessonVersion } from 'app/core/lessons/lessons.types';
import { SchoolAccessService } from 'app/core/school-access/school-access.service';
import {
    capabilitiesFor,
    Membership,
} from 'app/core/school-access/school-access.types';
import { catchError, forkJoin, map, of } from 'rxjs';
import { LessonReuseComponent } from '../shared/lesson-reuse/lesson-reuse.component';

import { LessonContentComponent } from '../shared/lesson-content.component';
import { LessonStatusComponent } from '../shared/lesson-status.component';

export type VersionRelation =
    | 'live'
    | 'latest'
    | 'live-and-latest'
    | 'superseded';

/** How a version relates to the plan's live publication and latest version. */
export function versionRelation(
    plan: LessonPlan,
    version: number
): VersionRelation {
    const live = plan.publishedVersion === version;
    const latest = plan.currentVersion === version;
    if (live && latest) return 'live-and-latest';
    if (live) return 'live';
    if (latest) return 'latest';
    return 'superseded';
}

@Component({
    selector: 'app-lesson-version',
    standalone: true,
    imports: [
        LessonReuseComponent,
        DatePipe,
        MatButtonModule,
        RouterLink,
        LessonContentComponent,
        LessonStatusComponent,
    ],
    changeDetection: ChangeDetectionStrategy.OnPush,
    templateUrl: './lesson-version.component.html',
})
export class LessonVersionComponent {
    private readonly lessons = inject(LessonsService);
    private readonly route = inject(ActivatedRoute);
    readonly planID = this.route.snapshot.paramMap.get('planId') ?? '';
    readonly number = toSignal(
        this.route.paramMap.pipe(
            map((params) => Number(params.get('version')))
        ),
        { initialValue: Number(this.route.snapshot.paramMap.get('version')) }
    );

    readonly memberships = signal<Membership[]>([]);
    readonly canReuse = computed(() => {
        const plan = this.plan();
        return (
            !!plan &&
            capabilitiesFor(
                this.memberships(),
                plan.schoolID,
                plan.departmentID
            ).includes('teach')
        );
    });
    readonly plan = signal<LessonPlan | null>(null);
    readonly versions = signal<LessonVersion[]>([]);
    readonly loading = signal(true);
    readonly errorMessage = signal<string | null>(null);

    readonly version = computed(
        () => this.versions().find((v) => v.version === this.number()) ?? null
    );
    readonly relation = computed(() => {
        const plan = this.plan();
        return plan ? versionRelation(plan, this.number()) : null;
    });

    constructor() {
        forkJoin([
            this.lessons.plan(this.planID),
            this.lessons.versions(this.planID),
            inject(SchoolAccessService)
                .ownMemberships()
                .pipe(catchError(() => of([]))),
        ]).subscribe({
            next: ([plan, versions, memberships]) => {
                this.memberships.set(memberships);
                this.plan.set(plan);
                this.versions.set(versions);
                this.loading.set(false);
            },
            error: (error: HttpErrorResponse) => {
                this.loading.set(false);
                this.errorMessage.set(
                    error.status === 404
                        ? 'This lesson plan does not exist or is not visible to you.'
                        : jsonApiErrorMessage(
                              error as JsonApiErrorResponse,
                              'Unable to load this version.'
                          )
                );
            },
        });
    }
}
