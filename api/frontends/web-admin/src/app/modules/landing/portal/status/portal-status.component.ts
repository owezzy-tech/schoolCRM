import { ChangeDetectionStrategy, Component } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { NgClass, TitleCasePipe } from '@angular/common';

@Component({
    selector: 'app-portal-status',
    standalone: true,
    imports: [
        MatButtonModule,
        MatIconModule,
        MatProgressBarModule,
        NgClass,
        TitleCasePipe
    ],
    template: `
        <div class="flex flex-col flex-auto py-12 px-6 sm:px-10 lg:px-16">
            <div class="mx-auto w-full max-w-6xl">
                <!-- Header section -->
                <div class="flex flex-col items-start gap-4 sm:flex-row sm:items-center sm:justify-between">
                    <div class="flex flex-col">
                        <h1 class="text-3xl font-bold tracking-tight text-default sm:text-4xl">Application status</h1>
                        <h2 class="mt-1 text-lg text-secondary">APP-3018 &middot; Data Science MSc &middot; Fall 2026</h2>
                    </div>
                    <div class="inline-flex items-center rounded-full bg-yellow-100 px-3 py-1 text-sm font-medium text-yellow-900">
                        Under Review
                    </div>
                </div>

                <!-- Overall progress card -->
                <div class="mt-6 flex flex-col rounded-2xl border bg-card p-6 shadow-sm">
                    <h3 class="mb-4 text-lg font-semibold text-default">Overall progress</h3>
                    <mat-progress-bar mode="determinate" [value]="62" class="mb-2"></mat-progress-bar>
                    <p class="text-sm font-medium text-secondary">62% complete &middot; 3 of 5 stages done.</p>
                </div>

                <!-- 2-column grid -->
                <div class="mt-6 grid grid-cols-1 gap-6 md:grid-cols-2">
                    <!-- Timeline Card -->
                    <div class="flex flex-col rounded-2xl border bg-card p-6 shadow-sm">
                        <h3 class="mb-6 text-lg font-semibold text-default">Timeline</h3>
                        <div class="flex flex-col gap-6 relative">
                            <!-- Vertical dashed line -->
                            <div class="absolute bottom-4 left-[11px] top-4 border-l-2 border-dashed border-gray-300"></div>
                            
                            @for (item of timeline; track item.label) {
                                <div class="flex items-start gap-4 relative z-10">
                                    <div class="flex h-6 w-6 shrink-0 items-center justify-center rounded-full mt-0.5"
                                         [ngClass]="item.done ? 'bg-primary' : 'border-2 border-gray-300 bg-card'">
                                    </div>
                                    <div class="flex flex-col">
                                        <span class="font-medium text-default leading-tight">{{item.label}}</span>
                                        <span class="text-sm text-secondary">{{item.date}}</span>
                                    </div>
                                </div>
                            }
                        </div>
                    </div>

                    <!-- Documents Card -->
                    <div class="flex flex-col rounded-2xl border bg-card p-6 shadow-sm">
                        <h3 class="mb-2 text-lg font-semibold text-default">Documents</h3>
                        <div class="flex flex-col">
                            @for (doc of documents; track doc.label) {
                                <div class="flex items-center justify-between border-b py-3 last:border-b-0">
                                    <div class="flex items-center gap-3">
                                        <mat-icon svgIcon="heroicons_outline:document-text" class="text-secondary icon-size-5"></mat-icon>
                                        <span class="font-medium text-default">{{doc.label}}</span>
                                    </div>
                                    <div class="rounded-full px-2.5 py-0.5 text-xs font-medium"
                                         [ngClass]="{
                                            'bg-green-100 text-green-800': doc.status === 'verified',
                                            'bg-blue-100 text-blue-800': doc.status === 'received',
                                            'bg-red-100 text-red-800': doc.status === 'missing'
                                         }">
                                        {{doc.status | titlecase}}
                                    </div>
                                </div>
                            }
                        </div>
                    </div>
                </div>

                <!-- Counselor card -->
                <div class="mt-6 flex items-center gap-4 rounded-2xl bg-primary-50 p-6 dark:bg-primary-900/20">
                    <div class="flex h-12 w-12 shrink-0 items-center justify-center rounded-full bg-primary/10">
                        <mat-icon svgIcon="heroicons_outline:envelope" class="text-primary icon-size-6"></mat-icon>
                    </div>
                    <div class="flex flex-col">
                        <h3 class="font-semibold text-default">Reach your counselor</h3>
                        <p class="text-secondary mt-0.5">Maya Schultz &middot; maya&#64;northbrook.edu &middot; Replies within 1 business day.</p>
                    </div>
                </div>
            </div>
        </div>
    `,
    changeDetection: ChangeDetectionStrategy.OnPush,
    host: { class: 'flex w-full flex-auto flex-col' },
})
export class PortalStatusComponent {
    readonly timeline = [
        { label: 'Submitted', date: 'May 12, 2026', done: true },
        { label: 'Documents received', date: 'May 14, 2026', done: true },
        { label: 'Under review', date: 'May 18, 2026', done: true },
        { label: 'Interview', date: 'Scheduled for Jun 2', done: false },
        { label: 'Decision', date: 'Expected late June', done: false },
    ];

    readonly documents: ReadonlyArray<{ label: string; status: 'verified' | 'received' | 'missing' }> = [
        { label: 'Personal statement', status: 'verified' },
        { label: 'Transcripts', status: 'verified' },
        { label: 'Recommendation letter 1', status: 'received' },
        { label: 'Recommendation letter 2', status: 'missing' },
        { label: 'Test scores', status: 'verified' },
    ];
}
