import { ChangeDetectionStrategy, Component } from '@angular/core';
import { UntypedFormBuilder, UntypedFormGroup, FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';

@Component({
    selector: 'app-portal-inquiry',
    standalone: true,
    imports: [
        FormsModule,
        ReactiveFormsModule,
        MatButtonModule,
        MatFormFieldModule,
        MatIconModule,
        MatInputModule
    ],
    template: `
        <div class="flex flex-col flex-auto py-12 px-6 sm:px-10 lg:px-16">
            <div class="mx-auto w-full max-w-6xl text-center">
                <h1 class="text-3xl font-bold tracking-tight text-default sm:text-4xl">Ask admissions</h1>
                <p class="mt-2 mb-8 text-lg text-secondary">Send a question and we'll get back to you within one business day.</p>
                
                <div class="mx-auto max-w-3xl rounded-2xl border bg-card p-8 shadow-sm text-left">
                    <form [formGroup]="form" class="flex flex-col gap-4">
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
                                <mat-label>Phone (optional)</mat-label>
                                <input matInput type="tel" formControlName="phone">
                            </mat-form-field>
                            
                            <mat-form-field appearance="outline" class="w-full">
                                <mat-label>Program of interest</mat-label>
                                <input matInput formControlName="program" placeholder="e.g. Computer Science BSc">
                            </mat-form-field>
                            
                            <mat-form-field appearance="outline" class="w-full">
                                <mat-label>Intake term</mat-label>
                                <input matInput formControlName="term" placeholder="Fall 2026">
                            </mat-form-field>
                        </div>
                        
                        <mat-form-field appearance="outline" class="w-full">
                            <mat-label>Your question</mat-label>
                            <textarea matInput rows="6" formControlName="question" placeholder="Tell us what you'd like to know..."></textarea>
                        </mat-form-field>
                        
                        <p class="text-sm text-secondary">By submitting, you agree to our privacy policy.</p>
                        
                        <div class="mt-4 flex justify-end">
                            <button mat-flat-button color="primary" type="button">
                                <mat-icon svgIcon="heroicons_outline:paper-airplane" class="icon-size-5 mr-2"></mat-icon>
                                Send inquiry
                            </button>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    `,
    changeDetection: ChangeDetectionStrategy.OnPush,
    host: { class: 'flex w-full flex-auto flex-col' },
})
export class PortalInquiryComponent {
    form: UntypedFormGroup;

    constructor(private _formBuilder: UntypedFormBuilder) {
        this.form = this._formBuilder.group({
            firstName: [''],
            lastName: [''],
            email: [''],
            phone: [''],
            program: [''],
            term: [''],
            question: ['']
        });
    }
}
