import { HttpErrorResponse } from '@angular/common/http';
import {
    ChangeDetectionStrategy,
    Component,
    computed,
    DestroyRef,
    inject,
    signal,
} from '@angular/core';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { jsonApiErrorMessage } from 'app/core/api/json-api';
import { AuthService } from 'app/core/auth/auth.service';
import { LessonAssistantService } from 'app/core/lessons/lesson-assistant.service';
import {
    AssistantEvent,
    AssistantStreamError,
    GenerationReceipt,
    GenerationRequest,
    LessonThread,
} from 'app/core/lessons/lesson-assistant.types';
import { LessonScope, loadLessonScopes } from 'app/core/lessons/lesson-scopes';
import { LessonsService } from 'app/core/lessons/lessons.service';
import { LessonCitation, LessonVersion } from 'app/core/lessons/lessons.types';
import { SchoolAccessService } from 'app/core/school-access/school-access.service';
import { Subscription } from 'rxjs';
import { LessonContentComponent } from '../../shared/lesson-content.component';
import { LessonEvidenceComponent } from '../../shared/lesson-evidence.component';

@Component({
    selector: 'app-lesson-assistant',
    imports: [
        FormsModule,
        MatButtonModule,
        MatFormFieldModule,
        MatInputModule,
        MatSelectModule,
        RouterLink,
        LessonContentComponent,
        LessonEvidenceComponent,
    ],
    changeDetection: ChangeDetectionStrategy.OnPush,
    templateUrl: './lesson-assistant.component.html',
})
export class LessonAssistantComponent {
    private readonly assistant = inject(LessonAssistantService);
    private readonly lessons = inject(LessonsService);
    private readonly access = inject(SchoolAccessService);
    private readonly auth = inject(AuthService);
    private readonly route = inject(ActivatedRoute);
    private readonly router = inject(Router);
    private readonly destroy = inject(DestroyRef);
    private stream?: Subscription;
    private history?: Subscription;
    private requestRead?: Subscription;
    private versionRead?: Subscription;
    readonly scopes = signal<LessonScope[]>([]);
    readonly scopeKey = signal('');
    readonly scope = computed(
        () => this.scopes().find((s) => s.key === this.scopeKey()) ?? null
    );
    readonly loading = signal(true);
    readonly threads = signal<LessonThread[]>([]);
    readonly request = signal<GenerationRequest | null>(null);
    readonly receipt = signal<GenerationReceipt | null>(null);
    readonly version = signal<LessonVersion | null>(null);
    readonly citations = signal<LessonCitation[]>([]);
    readonly evidence = computed(() => ({ citations: this.citations() }));
    readonly busy = signal(false);
    readonly completed = signal(false);
    readonly progress = signal('');
    readonly errorMessage = signal<string | null>(null);
    readonly resultError = signal<string | null>(null);
    readonly approvalPrompt = signal(false);
    readonly framework = signal('kenya-cbc');
    readonly stage = signal('grade-1');
    readonly subject = signal('english');
    readonly revision = signal('2024');
    readonly topic = signal('');
    readonly duration = signal(30);
    readonly valid = computed(
        () =>
            !!this.scope() &&
            [
                this.framework(),
                this.stage(),
                this.subject(),
                this.revision(),
                this.topic(),
            ].every((v) => v.trim().length > 0) &&
            Number.isInteger(this.duration()) &&
            this.duration() >= 1 &&
            this.duration() <= 180
    );

    constructor() {
        this.loadScopes();
    }
    loadScopes(): void {
        this.loading.set(true);
        this.errorMessage.set(null);
        loadLessonScopes(this.access)
            .pipe(takeUntilDestroyed(this.destroy))
            .subscribe({
                next: (scopes) => {
                    this.scopes.set(
                        scopes.filter((s) => s.capabilities.includes('teach'))
                    );
                    const params = this.route.snapshot.queryParamMap;
                    const preferred = `${params.get('school')}:${params.get('department')}`;
                    this.scopeKey.set(
                        this.scopes().find((s) => s.key === preferred)?.key ??
                            this.scopes()[0]?.key ??
                            ''
                    );
                    this.loading.set(false);
                    this.refreshHistory();
                    if (params.get('thread'))
                        this.openRequest(params.get('thread')!);
                },
                error: (error) => {
                    this.loading.set(false);
                    this.fail(
                        error,
                        'Unable to load your lesson responsibilities.'
                    );
                },
            });
    }
    selectScope(key: string): void {
        if (this.busy() || !this.scopes().some((s) => s.key === key)) return;
        this.reset();
        this.scopeKey.set(key);
        this.threads.set([]);
        this.refreshHistory();
        this.navigateRequest(null);
    }
    refreshHistory(): void {
        const scope = this.scope();
        if (!scope) return;
        this.history?.unsubscribe();
        this.history = this.assistant
            .threads(scope.schoolID, scope.departmentID)
            .pipe(takeUntilDestroyed(this.destroy))
            .subscribe({
                next: (threads) => this.threads.set(threads),
                error: (error) =>
                    this.fail(error, 'Unable to load private request history.'),
            });
    }
    start(): void {
        const scope = this.scope();
        if (!scope || !this.valid() || this.busy() || this.request()) return;
        const request: GenerationRequest = {
            request_id: crypto.randomUUID(),
            school_id: scope.schoolID,
            department_id: scope.departmentID,
            framework: this.framework().trim(),
            stage: this.stage().trim(),
            subject: this.subject().trim(),
            revision: this.revision().trim(),
            topic: this.topic().trim(),
            duration_minutes: this.duration(),
        };
        this.requestRead?.unsubscribe(); // A late history read must not replace it.
        this.request.set(request);
        this.navigateRequest(request.request_id);
        this.reconnect();
    }
    openRequest(id: string): void {
        if (this.busy()) return;
        this.reset();
        this.requestRead = this.assistant
            .thread(id)
            .pipe(takeUntilDestroyed(this.destroy))
            .subscribe({
                next: (thread) => {
                    const key = `${thread.request.school_id}:${thread.request.department_id}`;
                    if (!this.scopes().some((s) => s.key === key)) {
                        this.errorMessage.set(
                            'You no longer have teaching access to this request.'
                        );
                        return;
                    }
                    this.scopeKey.set(key);
                    this.request.set(thread.request);
                    this.completed.set(thread.status === 'completed');
                    this.navigateRequest(id);
                    this.refreshHistory();
                    if (thread.result) {
                        this.receipt.set(thread.result);
                        this.approvalPrompt.set(true);
                        this.loadResult();
                    }
                },
                error: (error) =>
                    this.fail(error, 'This private request is unavailable.'),
            });
    }
    reconnect(): void {
        const request = this.request();
        if (!request || this.busy() || this.completed()) return;
        this.stream?.unsubscribe();
        this.errorMessage.set(null);
        this.busy.set(true);
        this.progress.set('Connecting to your saved request…');
        this.stream = this.assistant
            .generate(request)
            .pipe(takeUntilDestroyed(this.destroy))
            .subscribe({
                next: (event) => this.receive(event),
                error: (error) => {
                    this.busy.set(false);
                    if (
                        error instanceof HttpErrorResponse &&
                        error.status === 401
                    )
                        return this.signInAgain();
                    this.fail(
                        error,
                        'Connection interrupted. Reconnect to recover this same request.'
                    );
                    this.refreshHistory();
                },
                complete: () => {
                    this.busy.set(false);
                    this.refreshHistory();
                },
            });
    }
    newRequest(): void {
        if (this.busy()) return;
        this.reset();
        this.navigateRequest(null);
    }
    loadResult(): void {
        const result = this.receipt();
        if (!result) return;
        this.resultError.set(null);
        this.versionRead?.unsubscribe();
        this.versionRead = this.lessons
            .versions(result.planID)
            .pipe(takeUntilDestroyed(this.destroy))
            .subscribe({
                next: (versions) => {
                    this.version.set(
                        versions.find((v) => v.version === result.version) ??
                            null
                    );
                    if (!this.version())
                        this.resultError.set(
                            'The saved version is no longer visible to you.'
                        );
                },
                error: (error) =>
                    this.resultError.set(
                        jsonApiErrorMessage(
                            error,
                            'Unable to load the saved draft. Retry reading it.'
                        )
                    ),
            });
    }
    private receive(event: AssistantEvent): void {
        switch (event.type) {
            case 'progress':
                this.progress.set(event.detail);
                break;
            case 'citation':
                this.citations.set(event.citations);
                break;
            case 'structured-result':
                this.receipt.set(event.result);
                this.loadResult();
                break;
            case 'approval-request':
                this.approvalPrompt.set(
                    event.planID === this.receipt()?.planID &&
                        event.version === this.receipt()?.version
                );
                break;
            case 'completed':
                this.completed.set(true);
                this.progress.set('Draft saved.');
                break;
            case 'terminal-error':
                if (event.status === 401) return this.signInAgain();
                this.errorMessage.set(event.detail);
                this.progress.set('Request stopped.');
                break;
        }
    }
    private reset(): void {
        this.stop();
        this.request.set(null);
        this.receipt.set(null);
        this.version.set(null);
        this.citations.set([]);
        this.completed.set(false);
        this.approvalPrompt.set(false);
        this.progress.set('');
        this.errorMessage.set(null);
        this.resultError.set(null);
        this.busy.set(false);
    }
    private stop(): void {
        this.stream?.unsubscribe();
        this.requestRead?.unsubscribe();
        this.versionRead?.unsubscribe();
    }
    private navigateRequest(id: string | null): void {
        const scope = this.scope();
        this.router.navigate([], {
            relativeTo: this.route,
            replaceUrl: true,
            queryParams: {
                school: scope?.schoolID,
                department: scope?.departmentID,
                thread: id,
            },
        });
    }
    /** Go rejected the bearer, so a retry cannot succeed; sign in and return here. */
    private signInAgain(): void {
        this.stop();
        this.busy.set(false);
        this.auth.signOut();
        this.router.navigate(['/sign-in'], {
            queryParams: { redirectURL: this.router.url },
        });
    }
    private fail(error: unknown, fallback: string): void {
        this.errorMessage.set(
            error instanceof HttpErrorResponse
                ? jsonApiErrorMessage(error, fallback)
                : error instanceof AssistantStreamError
                  ? error.message
                  : fallback
        );
    }
}
