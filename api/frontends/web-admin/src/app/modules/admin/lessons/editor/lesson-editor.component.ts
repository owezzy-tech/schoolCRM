import { HttpErrorResponse } from '@angular/common/http';
import {
    ChangeDetectionStrategy,
    Component,
    inject,
    signal,
} from '@angular/core';
import {
    FormArray,
    FormBuilder,
    FormGroup,
    FormControl,
    ReactiveFormsModule,
    Validators,
} from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { jsonApiErrorMessage, JsonApiErrorResponse } from 'app/core/api/json-api';
import { LessonsService } from 'app/core/lessons/lessons.service';
import {
    LessonContent,
    LessonDraft,
    LessonPlan,
    toLessonContent,
} from 'app/core/lessons/lessons.types';
import { forkJoin, map, Observable } from 'rxjs';

type ActivityForm = FormGroup<{
    minutes: FormControl<number>;
    title: FormControl<string>;
    detail: FormControl<string>;
}>;

@Component({
    selector: 'app-lesson-editor',
    standalone: true,
    imports: [
        MatButtonModule,
        MatFormFieldModule,
        MatIconModule,
        MatInputModule,
        ReactiveFormsModule,
        RouterLink,
    ],
    changeDetection: ChangeDetectionStrategy.OnPush,
    templateUrl: './lesson-editor.component.html',
})
export class LessonEditorComponent {
    private readonly lessons = inject(LessonsService);
    private readonly route = inject(ActivatedRoute);
    private readonly router = inject(Router);
    private readonly fb = inject(FormBuilder).nonNullable;

    readonly planID = this.route.snapshot.paramMap.get('planId');
    private readonly schoolID = this.route.snapshot.queryParamMap.get('school');
    private readonly departmentID = this.route.snapshot.queryParamMap.get('department');

    /** The plan being revised; null when creating. */
    readonly plan = signal<LessonPlan | null>(null);
    readonly loading = signal(this.planID !== null);
    readonly saving = signal(false);
    readonly conflict = signal(false);
    readonly errorMessage = signal<string | null>(null);

    readonly form = this.fb.group({
        title: ['', [Validators.required, Validators.maxLength(200)]],
        objectives: this.fb.array<FormControl<string>>([]),
        materials: this.fb.array<FormControl<string>>([]),
        activities: this.fb.array<ActivityForm>([]),
        assessment: [''],
        changeSummary: ['', Validators.maxLength(1000)],
    });

    constructor() {
        if (this.planID) {
            this.load(this.planID);
        } else {
            this.setContent({ objectives: [''], materials: [''], activities: [], assessment: '' });
            this.addActivity();
        }
    }

    /** Saving always creates the next version. */
    get nextVersion(): number {
        return (this.plan()?.currentVersion ?? 0) + 1;
    }

    get objectives() {
        return this.form.controls.objectives;
    }
    get materials() {
        return this.form.controls.materials;
    }
    get activities() {
        return this.form.controls.activities;
    }

    addItem(list: FormArray<FormControl<string>>): void {
        list.push(this.fb.control(''));
    }

    addActivity(): void {
        this.activities.push(
            this.fb.group({
                minutes: [10, [Validators.required, Validators.min(1)]],
                title: ['', Validators.required],
                detail: [''],
            })
        );
    }

    remove(list: FormArray, index: number): void {
        list.removeAt(index);
    }

    reload(): void {
        if (this.planID) this.load(this.planID);
    }

    save(): void {
        if (this.form.invalid) {
            this.form.markAllAsTouched();
            return;
        }
        const draft = this.toDraft();
        const plan = this.plan();
        this.saving.set(true);
        this.errorMessage.set(null);

        const savedPlanID: Observable<string> = plan
            ? this.lessons.revise(plan.id, plan.currentVersion, draft).pipe(map((v) => v.planID))
            : this.lessons
                  .create(this.schoolID ?? '', this.departmentID ?? '', draft)
                  .pipe(map((created) => created.id));
        savedPlanID.subscribe({
            next: (id) => {
                this.saving.set(false);
                this.router.navigate(['/lessons', id]);
            },
            error: (error: HttpErrorResponse) => {
                this.saving.set(false);
                if (error.status === 409) {
                    this.conflict.set(true);
                    return;
                }
                this.errorMessage.set(
                    jsonApiErrorMessage(error as JsonApiErrorResponse, 'Unable to save the lesson plan.')
                );
            },
        });
    }

    toDraft(): LessonDraft {
        const value = this.form.getRawValue();
        const filled = (items: string[]) => items.map((i) => i.trim()).filter(Boolean);
        return {
            title: value.title.trim(),
            changeSummary: value.changeSummary.trim(),
            content: {
                objectives: filled(value.objectives),
                materials: filled(value.materials),
                activities: value.activities
                    .filter((activity) => activity.title.trim())
                    .map((activity) => ({
                        minutes: Number(activity.minutes),
                        title: activity.title.trim(),
                        detail: activity.detail.trim(),
                    })),
                assessment: value.assessment.trim(),
            },
        };
    }

    private load(planID: string): void {
        this.loading.set(true);
        this.conflict.set(false);
        forkJoin([this.lessons.plan(planID), this.lessons.versions(planID)]).subscribe({
            next: ([plan, versions]) => {
                const current = versions.find((v) => v.version === plan.currentVersion);
                this.plan.set(plan);
                this.form.controls.title.setValue(current?.title ?? plan.title);
                this.form.controls.changeSummary.setValue('');
                this.setContent(toLessonContent(current?.content));
                this.loading.set(false);
            },
            error: (error: HttpErrorResponse) => {
                this.loading.set(false);
                this.errorMessage.set(
                    jsonApiErrorMessage(error as JsonApiErrorResponse, 'Unable to load the lesson plan.')
                );
            },
        });
    }

    private setContent(content: LessonContent): void {
        const strings = (items: string[]) => items.map((item) => this.fb.control(item));
        this.form.setControl('objectives', this.fb.array(strings(content.objectives)));
        this.form.setControl('materials', this.fb.array(strings(content.materials)));
        this.form.setControl(
            'activities',
            this.fb.array(
                content.activities.map((activity) =>
                    this.fb.group({
                        minutes: [activity.minutes, [Validators.required, Validators.min(1)]],
                        title: [activity.title, Validators.required],
                        detail: [activity.detail],
                    })
                )
            )
        );
        this.form.controls.assessment.setValue(content.assessment);
    }
}
