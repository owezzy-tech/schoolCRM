import { DOCUMENT, DatePipe, NgTemplateOutlet } from '@angular/common';
import { HttpErrorResponse } from '@angular/common/http';
import {
    ChangeDetectionStrategy,
    Component,
    DestroyRef,
    Injector,
    afterNextRender,
    computed,
    inject,
    signal,
} from '@angular/core';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { MatButtonModule } from '@angular/material/button';
import { Title } from '@angular/platform-browser';
import { ActivatedRoute, RouterLink } from '@angular/router';
import {
    JsonApiErrorResponse,
    jsonApiErrorMessage,
} from 'app/core/api/json-api';
import { LessonsService } from 'app/core/lessons/lessons.service';
import {
    LESSON_STATUS_LABELS,
    LessonPlan,
    LessonVersion,
} from 'app/core/lessons/lessons.types';
import {
    Subject,
    catchError,
    forkJoin,
    map,
    merge,
    of,
    switchMap,
    tap,
} from 'rxjs';

@Component({
    selector: 'app-lesson-export',
    standalone: true,
    imports: [DatePipe, NgTemplateOutlet, MatButtonModule, RouterLink],
    changeDetection: ChangeDetectionStrategy.OnPush,
    templateUrl: './lesson-export.component.html',
    styleUrl: './lesson-export.component.scss',
})
export class LessonExportComponent {
    private readonly lessons = inject(LessonsService);
    private readonly route = inject(ActivatedRoute);
    private readonly document = inject(DOCUMENT);
    private readonly injector = inject(Injector);
    private readonly destroyRef = inject(DestroyRef);
    private readonly title = inject(Title);
    private readonly printRequests = new Subject<void>();

    readonly plan = signal<LessonPlan | null>(null);
    readonly version = signal<LessonVersion | null>(null);
    readonly loading = signal(true);
    readonly errorMessage = signal<string | null>(null);
    readonly statusLabels = LESSON_STATUS_LABELS;
    get routePlanID(): string {
        return this.route.snapshot.paramMap.get('planId') ?? '';
    }
    readonly publication = computed(() => {
        const version = this.version();
        return this.plan()?.publishedVersion === version?.version
            ? 'Live publication'
            : version?.publishedAt
              ? 'Previously published'
              : 'Not published';
    });

    constructor() {
        const previousTitle = this.title.getTitle();
        this.destroyRef.onDestroy(() => this.title.setTitle(previousTitle));
        merge(
            this.route.paramMap.pipe(map(() => false)),
            this.printRequests.pipe(map(() => true))
        )
            .pipe(
                tap(() => {
                    this.loading.set(true);
                    this.version.set(null);
                    this.plan.set(null);
                    this.errorMessage.set(null);
                }),
                switchMap((print) => {
                    const planID =
                        this.route.snapshot.paramMap.get('planId') ?? '';
                    const number = Number(
                        this.route.snapshot.paramMap.get('version')
                    );
                    return forkJoin({
                        plan: this.lessons.plan(planID),
                        versions: this.lessons.versions(planID),
                    }).pipe(
                        map(({ plan, versions }) => ({
                            plan,
                            version: versions.find((v) => v.version === number),
                            print,
                        })),
                        catchError((error: HttpErrorResponse) => {
                            this.errorMessage.set(
                                jsonApiErrorMessage(
                                    error as JsonApiErrorResponse,
                                    'This lesson version is unavailable or you no longer have access.'
                                )
                            );
                            return of(null);
                        })
                    );
                }),
                takeUntilDestroyed()
            )
            .subscribe((result) => {
                this.loading.set(false);
                if (!result) return;
                if (!result.version) {
                    this.errorMessage.set(
                        'This version does not exist or is not visible to you.'
                    );
                    return;
                }
                this.plan.set(result.plan);
                this.version.set(result.version);
                this.title.setTitle(
                    `${result.version.title} - version ${result.version.version}`
                );
                if (result.print) {
                    afterNextRender(
                        () => {
                            if (
                                !this.loading() &&
                                this.version()?.id === result.version?.id
                            ) {
                                this.document.defaultView?.print();
                            }
                        },
                        { injector: this.injector }
                    );
                }
            });
    }

    print(): void {
        // Refresh authorised reads before opening the print dialog, including after revocation.
        this.printRequests.next();
    }

    entries(value: unknown): [string, unknown][] | null {
        return value !== null &&
            typeof value === 'object' &&
            !Array.isArray(value)
            ? Object.entries(value)
            : null;
    }

    items(value: unknown): unknown[] | null {
        return Array.isArray(value) ? value : null;
    }

    label(key: string): string {
        const text = key
            .replace(/([a-z])([A-Z])/g, '$1 $2')
            .replace(/[_-]/g, ' ');
        return text.charAt(0).toUpperCase() + text.slice(1);
    }
}
