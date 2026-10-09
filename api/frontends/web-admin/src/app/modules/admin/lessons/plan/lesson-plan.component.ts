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
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatRadioModule } from '@angular/material/radio';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { ActivatedRoute, RouterLink } from '@angular/router';
import {
    jsonApiErrorMessage,
    JsonApiErrorResponse,
} from 'app/core/api/json-api';
import { LessonsService } from 'app/core/lessons/lessons.service';
import {
    lessonActions,
    LessonDecision,
    LessonPlan,
    LessonVersion,
    SEPARATION_OF_DUTIES_MESSAGE,
} from 'app/core/lessons/lessons.types';
import { SchoolAccessService } from 'app/core/school-access/school-access.service';
import {
    capabilitiesFor,
    Membership,
} from 'app/core/school-access/school-access.types';
import { UserService } from 'app/core/user/user.service';
import { forkJoin, Observable } from 'rxjs';
import { LessonReuseComponent } from '../shared/lesson-reuse/lesson-reuse.component';

import { LessonContentComponent } from '../shared/lesson-content.component';
import { LessonStatusComponent } from '../shared/lesson-status.component';

@Component({
    selector: 'app-lesson-plan',
    standalone: true,
    imports: [
        DatePipe,
        FormsModule,
        MatButtonModule,
        MatFormFieldModule,
        MatInputModule,
        MatRadioModule,
        MatSnackBarModule,
        RouterLink,
        LessonContentComponent,
        LessonReuseComponent,
        LessonStatusComponent,
    ],
    changeDetection: ChangeDetectionStrategy.OnPush,
    templateUrl: './lesson-plan.component.html',
})
export class LessonPlanComponent {
    private readonly lessons = inject(LessonsService);
    private readonly schoolAccess = inject(SchoolAccessService);
    private readonly snackBar = inject(MatSnackBar);
    readonly planID =
        inject(ActivatedRoute).snapshot.paramMap.get('planId') ?? '';

    readonly separationMessage = SEPARATION_OF_DUTIES_MESSAGE;
    readonly viewer = toSignal(inject(UserService).user$);
    readonly plan = signal<LessonPlan | null>(null);
    readonly versions = signal<LessonVersion[]>([]);
    readonly memberships = signal<Membership[]>([]);
    readonly loading = signal(true);
    readonly busy = signal(false);
    readonly conflict = signal(false);
    readonly errorMessage = signal<string | null>(null);
    readonly decision = signal<LessonDecision>('approve');
    readonly feedback = signal('');

    readonly current = computed(() => {
        const plan = this.plan();
        return (
            this.versions().find((v) => v.version === plan?.currentVersion) ??
            null
        );
    });
    readonly available = computed(() => {
        const plan = this.plan();
        const current = this.current();
        const viewer = this.viewer();
        if (!plan || !current || !viewer)
            return { actions: [], separationOfDuties: false };
        return lessonActions(
            plan,
            current,
            viewer.id,
            capabilitiesFor(
                this.memberships(),
                plan.schoolID,
                plan.departmentID
            )
        );
    });
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
    readonly has = (action: string) =>
        this.available().actions.some((available) => available === action);
    readonly feedbackRequired = computed(
        () => this.decision() === 'request_changes' && !this.feedback().trim()
    );

    constructor() {
        this.load();
    }

    load(): void {
        this.loading.set(true);
        this.conflict.set(false);
        this.errorMessage.set(null);
        forkJoin([
            this.lessons.plan(this.planID),
            this.lessons.versions(this.planID),
            this.schoolAccess.ownMemberships(),
        ]).subscribe({
            next: ([plan, versions, memberships]) => {
                this.plan.set(plan);
                this.versions.set(versions);
                this.memberships.set(memberships);
                this.loading.set(false);
            },
            error: (error: HttpErrorResponse) => {
                this.loading.set(false);
                this.errorMessage.set(
                    error.status === 404
                        ? 'This lesson plan does not exist or is not visible to you.'
                        : jsonApiErrorMessage(
                              error as JsonApiErrorResponse,
                              'Unable to load the lesson plan.'
                          )
                );
            },
        });
    }

    submit(): void {
        this.run(
            (plan, version) => this.lessons.submit(plan, version),
            'Submitted for HOD review.'
        );
    }

    publish(): void {
        this.run(
            (plan, version) => this.lessons.publish(plan, version),
            'Published. This version is now live.'
        );
    }

    decide(): void {
        if (this.feedbackRequired()) return;
        const decision = this.decision();
        const feedback = this.feedback().trim();
        this.run(
            (plan, version) =>
                this.has('review')
                    ? this.lessons.review(plan, version, decision, feedback)
                    : this.lessons.approve(plan, version, decision, feedback),
            decision === 'approve' ? 'Decision recorded.' : 'Changes requested.'
        );
    }

    private run(
        command: (planID: string, version: number) => Observable<LessonVersion>,
        success: string
    ): void {
        const current = this.current();
        if (!current) return;
        this.busy.set(true);
        command(this.planID, current.version).subscribe({
            next: () => {
                this.busy.set(false);
                this.feedback.set('');
                this.decision.set('approve');
                this.snackBar.open(success, undefined, { duration: 4000 });
                this.load();
            },
            error: (error: HttpErrorResponse) => {
                this.busy.set(false);
                if (error.status === 409) {
                    this.conflict.set(true);
                    this.snackBar.open(
                        'That version is no longer current.',
                        undefined,
                        { duration: 5000 }
                    );
                    return;
                }
                this.errorMessage.set(
                    jsonApiErrorMessage(
                        error as JsonApiErrorResponse,
                        'The action could not be completed.'
                    )
                );
            },
        });
    }
}
