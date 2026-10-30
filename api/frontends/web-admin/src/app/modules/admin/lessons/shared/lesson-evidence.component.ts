import {
    ChangeDetectionStrategy,
    Component,
    computed,
    input,
} from '@angular/core';
import {
    lessonCitations,
    lessonContentRecord,
} from 'app/core/lessons/lessons.types';

@Component({
    selector: 'app-lesson-evidence',
    changeDetection: ChangeDetectionStrategy.OnPush,
    template: `
        <section aria-label="Curriculum evidence">
            <h3 class="mb-3 text-sm font-semibold uppercase tracking-wider">
                {{ heading() }}
            </h3>
            <ul class="space-y-3">
                @for (citation of citations(); track $index) {
                    <li class="rounded-lg border p-3">
                        <a
                            [href]="citation.sourceURL"
                            target="_blank"
                            rel="noopener noreferrer"
                            class="text-primary underline"
                            >{{ citation.title }}</a
                        >
                        <p>
                            PDF page {{ citation.page }} · revision
                            {{ citation.revision }}
                        </p>
                        @if (citation.passage) {
                            <blockquote
                                class="my-2 whitespace-pre-wrap border-l-2 pl-3"
                            >
                                {{ citation.passage }}
                            </blockquote>
                        }
                        <details class="mt-2 text-sm">
                            <summary class="cursor-pointer">
                                Evidence provenance
                            </summary>
                            <dl class="mt-2 break-all">
                                <dt class="font-medium">Source ID</dt>
                                <dd>{{ citation.sourceID }}</dd>
                                @if (citation.authority) {
                                    <dt class="font-medium">Authority</dt>
                                    <dd>{{ citation.authority }}</dd>
                                }
                                @if (citation.sourceSHA256) {
                                    <dt class="font-medium">Source SHA-256</dt>
                                    <dd>{{ citation.sourceSHA256 }}</dd>
                                }
                                @if (citation.indexRevision) {
                                    <dt class="font-medium">Index revision</dt>
                                    <dd>{{ citation.indexRevision }}</dd>
                                }
                                @if (citation.embeddingModel) {
                                    <dt class="font-medium">Embedding model</dt>
                                    <dd>{{ citation.embeddingModel }}</dd>
                                }
                            </dl>
                        </details>
                    </li>
                } @empty {
                    <li class="text-secondary">No references recorded.</li>
                }
            </ul>
            @if (origin(); as origin) {
                <p class="mt-4 text-sm">
                    Reused from plan {{ origin.planID }}, version
                    {{ origin.version }}. Approvals were not copied.
                </p>
            }
        </section>
    `,
})
export class LessonEvidenceComponent {
    readonly value = input.required<unknown>();
    readonly heading = input('Curriculum references');
    readonly citations = computed(() => lessonCitations(this.value()));
    readonly origin = computed(() => {
        const origin = lessonContentRecord(
            lessonContentRecord(this.value())['reusedFrom']
        );
        return typeof origin['planID'] === 'string' &&
            typeof origin['version'] === 'number'
            ? { planID: origin['planID'], version: origin['version'] }
            : null;
    });
}
