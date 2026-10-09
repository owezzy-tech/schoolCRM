import {
    ChangeDetectionStrategy,
    Component,
    computed,
    input,
} from '@angular/core';
import { toLessonContent } from 'app/core/lessons/lessons.types';

import { LessonEvidenceComponent } from './lesson-evidence.component';

/** Read-only rendering of one version's lesson content. */
@Component({
    selector: 'app-lesson-content',
    standalone: true,
    imports: [LessonEvidenceComponent],
    changeDetection: ChangeDetectionStrategy.OnPush,
    template: `
        @let lesson = content();
        <div class="flex flex-col gap-6">
            <section>
                <h3
                    class="text-secondary mb-2 text-sm font-semibold uppercase tracking-wider"
                >
                    Prerequisites
                </h3>
                <ul class="list-disc space-y-1 pl-5">
                    @for (item of lesson.prerequisites ?? []; track $index) {
                        <li>{{ item }}</li>
                    } @empty {
                        <li>None recorded.</li>
                    }
                </ul>
            </section>
            <section>
                <h3
                    class="text-secondary mb-2 text-sm font-semibold uppercase tracking-wider"
                >
                    Objectives
                </h3>
                @if (lesson.objectives.length) {
                    <ul class="list-disc space-y-1 pl-5">
                        @for (objective of lesson.objectives; track $index) {
                            <li>{{ objective }}</li>
                        }
                    </ul>
                } @else {
                    <p class="text-secondary">None recorded.</p>
                }
            </section>
            <section>
                <h3
                    class="text-secondary mb-2 text-sm font-semibold uppercase tracking-wider"
                >
                    Materials
                </h3>
                @if (lesson.materials.length) {
                    <ul class="list-disc space-y-1 pl-5">
                        @for (material of lesson.materials; track $index) {
                            <li>{{ material }}</li>
                        }
                    </ul>
                } @else {
                    <p class="text-secondary">None recorded.</p>
                }
            </section>
            <section>
                <h3
                    class="text-secondary mb-2 text-sm font-semibold uppercase tracking-wider"
                >
                    Timed activities
                    @if (totalMinutes()) {
                        <span class="font-mono normal-case tabular-nums"
                            >· {{ totalMinutes() }} min</span
                        >
                    }
                </h3>
                @if (lesson.activities.length) {
                    <ol class="divide-y rounded-lg border">
                        @for (activity of lesson.activities; track $index) {
                            <li class="flex gap-4 p-3">
                                <span
                                    class="w-16 shrink-0 font-mono tabular-nums"
                                    >{{ activity.minutes }} min</span
                                >
                                <div>
                                    <div class="font-medium">
                                        {{ activity.title }}
                                    </div>
                                    <div class="text-secondary">
                                        {{ activity.detail }}
                                    </div>
                                </div>
                            </li>
                        }
                    </ol>
                } @else {
                    <p class="text-secondary">None recorded.</p>
                }
            </section>
            <section>
                <h3
                    class="text-secondary mb-2 text-sm font-semibold uppercase tracking-wider"
                >
                    Assessment
                </h3>
                <p class="whitespace-pre-line">
                    {{ lesson.assessment || 'None recorded.' }}
                </p>
            </section>
            <section>
                <h3
                    class="text-secondary mb-2 text-sm font-semibold uppercase tracking-wider"
                >
                    Differentiation
                </h3>
                <p class="whitespace-pre-line">
                    {{ lesson.differentiation || 'None recorded.' }}
                </p>
            </section>
            <app-lesson-evidence [value]="value()" />
        </div>
    `,
})
export class LessonContentComponent {
    readonly value = input.required<unknown>();
    readonly content = computed(() => toLessonContent(this.value()));
    readonly totalMinutes = computed(() =>
        this.content().activities.reduce((sum, a) => sum + a.minutes, 0)
    );
}
