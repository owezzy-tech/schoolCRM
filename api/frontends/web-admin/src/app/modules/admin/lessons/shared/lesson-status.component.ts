import { ChangeDetectionStrategy, Component, computed, input } from '@angular/core';
import {
    LESSON_STATUS_LABELS,
    LESSON_STATUS_TONES,
    LessonStatus,
} from 'app/core/lessons/lessons.types';

const TONE_CLASSES = {
    neutral: 'bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-300',
    info: 'bg-blue-50 text-blue-800 dark:bg-blue-950 dark:text-blue-200',
    warning: 'bg-amber-50 text-amber-800 dark:bg-amber-950 dark:text-amber-200',
    success: 'bg-green-50 text-green-800 dark:bg-green-950 dark:text-green-200',
};

/** Status pill: dot plus text label, never colour alone. */
@Component({
    selector: 'app-lesson-status',
    standalone: true,
    changeDetection: ChangeDetectionStrategy.OnPush,
    template: `
        <span
            class="inline-flex items-center gap-1.5 rounded-full px-2.5 py-0.5 text-sm font-medium"
            [class]="toneClass()"
        >
            <span class="h-1.5 w-1.5 rounded-full bg-current" aria-hidden="true"></span>
            {{ label() }}
        </span>
    `,
})
export class LessonStatusComponent {
    readonly status = input.required<LessonStatus>();
    readonly label = computed(() => LESSON_STATUS_LABELS[this.status()]);
    readonly toneClass = computed(
        () => TONE_CLASSES[LESSON_STATUS_TONES[this.status()]]
    );
}
