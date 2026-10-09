import {
    ChangeDetectionStrategy,
    Component,
    computed,
    effect,
    inject,
    input,
    signal,
} from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { RouterLink } from '@angular/router';
import { jsonApiErrorMessage } from 'app/core/api/json-api';
import { LessonsService } from 'app/core/lessons/lessons.service';
import { Subscription } from 'rxjs';

@Component({
    selector: 'app-lesson-reuse',
    imports: [
        FormsModule,
        MatButtonModule,
        MatFormFieldModule,
        MatInputModule,
        RouterLink,
    ],
    changeDetection: ChangeDetectionStrategy.OnPush,
    template: `
        @if (allowed()) {
            <section
                class="bg-card rounded-xl border p-4"
                aria-label="Reuse lesson version"
            >
                <h2 class="mb-2 font-medium">Reuse version {{ version() }}</h2>
                <p class="text-secondary mb-3 text-sm">
                    Create your own draft with this version's content and
                    evidence. It needs new HOD review and dean approval.
                </p>
                @if (createdID()) {
                    <a mat-flat-button [routerLink]="['/lessons', createdID()]"
                        >Open your new draft</a
                    >
                } @else {
                    <mat-form-field class="w-full">
                        <mat-label>New draft title (optional)</mat-label>
                        <input
                            matInput
                            maxlength="200"
                            [readonly]="requestID() !== null"
                            [ngModel]="title()"
                            (ngModelChange)="title.set($event)"
                        />
                    </mat-form-field>
                    <button
                        mat-stroked-button
                        type="button"
                        [disabled]="busy()"
                        (click)="reuse()"
                    >
                        {{
                            requestID()
                                ? 'Retry creating draft'
                                : 'Create my draft'
                        }}
                    </button>
                }
                @if (errorMessage()) {
                    <p class="mt-2 text-red-700" role="alert">
                        {{ errorMessage() }}
                    </p>
                }
            </section>
        }
    `,
})
export class LessonReuseComponent {
    private readonly lessons = inject(LessonsService);
    readonly planID = input.required<string>();
    readonly version = input.required<number>();
    readonly allowed = input(false);
    readonly title = signal('');
    readonly busy = signal(false);
    readonly requestID = signal<string | null>(null);
    readonly createdID = signal<string | null>(null);
    readonly errorMessage = signal<string | null>(null);
    private command?: Subscription;
    private readonly source = computed(
        () => `${this.planID()}:${this.version()}`
    );
    constructor() {
        effect((onCleanup) => {
            this.source();
            onCleanup(() => this.command?.unsubscribe());
            this.busy.set(false);
            this.requestID.set(null);
            this.createdID.set(null);
            this.title.set('');
            this.errorMessage.set(null);
        });
    }
    reuse(): void {
        if (!this.allowed() || this.busy() || this.createdID()) return;
        this.requestID.update((id) => id ?? crypto.randomUUID());
        this.busy.set(true);
        this.errorMessage.set(null);
        this.command = this.lessons
            .reuse(
                this.planID(),
                this.version(),
                this.requestID()!,
                this.title().trim()
            )
            .subscribe({
                next: (plan) => {
                    this.createdID.set(plan.id);
                    this.busy.set(false);
                },
                error: (error) => {
                    this.busy.set(false);
                    this.errorMessage.set(
                        jsonApiErrorMessage(
                            error,
                            'Unable to confirm the new draft. Retry this same command.'
                        )
                    );
                },
            });
    }
}
