import { ChangeDetectionStrategy, Component } from '@angular/core';
import { UntypedFormBuilder, UntypedFormGroup, FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { NgClass } from '@angular/common';

@Component({
    selector: 'app-portal-apply',
    standalone: true,
    imports: [
        FormsModule,
        ReactiveFormsModule,
        MatButtonModule,
        MatFormFieldModule,
        MatIconModule,
        MatInputModule,
        MatProgressBarModule,
        NgClass
    ],
    template: `
        <div class="flex flex-col flex-auto py-12 px-6 sm:px-10 lg:px-16">
            <div class="mx-auto w-full max-w-6xl">
                <div class="mb-8 flex flex-col">
                    <span class="mb-2 text-sm font-semibold tracking-wider text-secondary uppercase">Step 1 of 6 - Account</span>
                    <h1 class="text-3xl font-bold tracking-tight text-default sm:text-4xl">Apply for Fall 2026</h1>
                </div>
                
                <div class="mb-10 w-full">
                    <mat-progress-bar mode="determinate" [value]="16"></mat-progress-bar>
                </div>
                
                <div class="grid grid-cols-1 gap-8 lg:grid-cols-[1fr_320px]">
                    <!-- Left Column - Form Card -->
                    <div class="rounded-2xl border bg-card p-8 shadow-sm">
                        <form [formGroup]="form" class="flex flex-col">
                            <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
                                <mat-form-field appearance="outline" class="w-full">
                                    <mat-label>First name</mat-label>
                                    <input matInput formControlName="firstName">
                                </mat-form-field>
                                
                                <mat-form-field appearance="outline" class="w-full">
                                    <mat-label>Last name</mat-label>
                                    <input matInput formControlName="lastName">
                                </mat-form-field>
                                
                                <mat-form-field appearance="outline" class="w-full">
                                    <mat-label>Email</mat-label>
                                    <input matInput type="email" formControlName="email">
                                </mat-form-field>
                                
                                <mat-form-field appearance="outline" class="w-full">
                                    <mat-label>Phone</mat-label>
                                    <input matInput type="tel" formControlName="phone">
                                </mat-form-field>
                                
                                <mat-form-field appearance="outline" class="w-full">
                                    <mat-label>Password</mat-label>
                                    <input matInput type="password" formControlName="password">
                                </mat-form-field>
                                
                                <mat-form-field appearance="outline" class="w-full">
                                    <mat-label>Confirm password</mat-label>
                                    <input matInput type="password" formControlName="confirmPassword">
                                </mat-form-field>
                            </div>
                            
                            <div class="mt-8 flex items-center justify-between">
                                <button mat-stroked-button type="button">
                                    <mat-icon svgIcon="heroicons_outline:arrow-left" class="icon-size-5"></mat-icon>
                                    <span class="ml-2">Back</span>
                                </button>
                                <button mat-flat-button color="primary" type="button">
                                    <span class="mr-2">Continue</span>
                                    <mat-icon svgIcon="heroicons_outline:arrow-right" class="icon-size-5"></mat-icon>
                                </button>
                            </div>
                        </form>
                    </div>
                    
                    <!-- Right Column - Sidebar -->
                    <div class="flex flex-col gap-6">
                        <!-- Progress Card -->
                        <div class="rounded-2xl border bg-card p-6 shadow-sm">
                            <h2 class="mb-6 text-sm font-bold tracking-wider text-secondary uppercase">Progress</h2>
                            <div class="flex flex-col gap-4">
                                @for (step of steps; track step.number) {
                                    <div class="flex items-center gap-4">
                                        <div class="flex h-8 w-8 items-center justify-center rounded-full text-sm font-bold"
                                             [ngClass]="step.active ? 'bg-primary text-on-primary' : 'border border-gray-300 text-secondary'">
                                            {{step.number}}
                                        </div>
                                        <span class="text-base" [ngClass]="step.active ? 'font-bold text-default' : 'font-medium text-secondary'">
                                            {{step.label}}
                                        </span>
                                    </div>
                                }
                            </div>
                        </div>
                        
                        <!-- Need Help Card -->
                        <div class="rounded-2xl bg-gray-50 p-6 dark:bg-gray-800">
                            <h2 class="mb-2 text-lg font-semibold text-default">Need help?</h2>
                            <p class="text-secondary leading-relaxed">
                                Email admissions&#64;northbrook.edu or chat with a counselor weekdays 9am-6pm ET.
                            </p>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    `,
    changeDetection: ChangeDetectionStrategy.OnPush,
    host: { class: 'flex w-full flex-auto flex-col' },
})
export class PortalApplyComponent {
    readonly steps = [
        { number: 1, label: 'Account', active: true },
        { number: 2, label: 'Program', active: false },
        { number: 3, label: 'Background', active: false },
        { number: 4, label: 'Essay', active: false },
        { number: 5, label: 'Documents', active: false },
        { number: 6, label: 'Review', active: false },
    ];
    
    form: UntypedFormGroup;

    constructor(private _formBuilder: UntypedFormBuilder) {
        this.form = this._formBuilder.group({
            firstName: [''],
            lastName: [''],
            email: [''],
            phone: [''],
            password: [''],
            confirmPassword: ['']
        });
    }
}
