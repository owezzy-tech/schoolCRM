import { ChangeDetectionStrategy, Component, OnInit, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterLink, ActivatedRoute, Router } from '@angular/router';
import { ReactiveFormsModule, UntypedFormBuilder, UntypedFormGroup, Validators } from '@angular/forms';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { MatRadioModule } from '@angular/material/radio';
import { MatDatepickerModule } from '@angular/material/datepicker';
import { MatNativeDateModule } from '@angular/material/core';
import { MOCK_EVENTS } from '../data/events.mock';

@Component({
    selector: 'app-event-form',
    standalone: true,
    imports: [
        CommonModule,
        RouterLink,
        ReactiveFormsModule,
        MatIconModule,
        MatButtonModule,
        MatFormFieldModule,
        MatInputModule,
        MatSelectModule,
        MatRadioModule,
        MatDatepickerModule,
        MatNativeDateModule
    ],
    changeDetection: ChangeDetectionStrategy.OnPush,
    templateUrl: './event-form.component.html'
})
export class EventFormComponent implements OnInit {
    private readonly fb = inject(UntypedFormBuilder);
    private readonly route = inject(ActivatedRoute);
    private readonly router = inject(Router);

    readonly isEditMode = signal(false);
    eventForm!: UntypedFormGroup;

    ngOnInit() {
        this.eventForm = this.fb.group({
            title: ['', Validators.required],
            type: ['open-day', Validators.required],
            status: ['draft', Validators.required],
            description: [''],
            startDate: [new Date(), Validators.required],
            startTime: ['09:00', Validators.required],
            endDate: [new Date(), Validators.required],
            endTime: ['10:00', Validators.required],
            isVirtual: [false, Validators.required],
            location: ['', Validators.required],
            capacity: [100, [Validators.required, Validators.min(0)]],
        });

        const id = this.route.snapshot.paramMap.get('id');
        if (id) {
            this.isEditMode.set(true);
            const evt = MOCK_EVENTS.find(e => e.id === id);
            if (evt) {
                const start = new Date(evt.start);
                const end = new Date(evt.end);
                
                this.eventForm.patchValue({
                    title: evt.title,
                    type: evt.type,
                    status: evt.status,
                    description: evt.description,
                    startDate: start,
                    startTime: start.toTimeString().slice(0, 5),
                    endDate: end,
                    endTime: end.toTimeString().slice(0, 5),
                    isVirtual: evt.isVirtual,
                    location: evt.location,
                    capacity: evt.capacity,
                });
            }
        }
    }

    save() {
        if (this.eventForm.invalid) {
            return;
        }
        
        // Mock save action
        console.log('Saved:', this.eventForm.value);
        this.router.navigate(['/events']);
    }
}
