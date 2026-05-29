import {
    Component,
    DestroyRef,
    inject,
    OnInit,
    ViewChild,
    ViewEncapsulation,
} from '@angular/core';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import {
    FormsModule,
    NgForm,
    ReactiveFormsModule,
    UntypedFormBuilder,
    UntypedFormGroup,
    Validators,
} from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { fuseAnimations } from '@fuse/animations';
import { FuseAlertComponent, FuseAlertType } from '@fuse/components/alert';
import { Actions, ofType } from '@ngrx/effects';
import { Store } from '@ngrx/store';
import { AuthApiActions, AuthPageActions } from 'app/core/auth/+state';

@Component({
    selector: 'auth-sign-in',
    templateUrl: './sign-in.component.html',
    encapsulation: ViewEncapsulation.None,
    animations: fuseAnimations,
    imports: [
        RouterLink,
        FuseAlertComponent,
        FormsModule,
        ReactiveFormsModule,
        MatFormFieldModule,
        MatInputModule,
        MatButtonModule,
        MatIconModule,
        MatCheckboxModule,
        MatProgressSpinnerModule,
    ],
})
export class AuthSignInComponent implements OnInit {
    @ViewChild('signInNgForm') signInNgForm: NgForm;

    alert: { type: FuseAlertType; message: string } = {
        type: 'success',
        message: '',
    };
    signInForm: UntypedFormGroup;
    showAlert: boolean = false;

    private readonly store = inject(Store);
    private readonly actions$ = inject(Actions);
    private readonly activatedRoute = inject(ActivatedRoute);
    private readonly formBuilder = inject(UntypedFormBuilder);
    private readonly router = inject(Router);
    private readonly destroyRef = inject(DestroyRef);

    ngOnInit(): void {
        this.signInForm = this.formBuilder.group({
            email: [
                'superadmin@example.com',
                [Validators.required, Validators.email],
            ],
            password: ['gophers', Validators.required],
            rememberMe: [''],
        });

        this.actions$
            .pipe(
                ofType(AuthApiActions.signInSuccess),
                takeUntilDestroyed(this.destroyRef)
            )
            .subscribe(() => {
                const redirectURL =
                    this.activatedRoute.snapshot.queryParamMap.get(
                        'redirectURL'
                    ) || '/signed-in-redirect';
                this.router.navigateByUrl(redirectURL);
            });

        this.actions$
            .pipe(
                ofType(AuthApiActions.signInFailure),
                takeUntilDestroyed(this.destroyRef)
            )
            .subscribe(({ error }) => {
                this.signInForm.enable();
                this.signInNgForm.resetForm();
                this.alert = { type: 'error', message: error };
                this.showAlert = true;
            });
    }

    signIn(): void {
        if (this.signInForm.invalid) {
            return;
        }
        this.signInForm.disable();
        this.showAlert = false;

        const { email, password } = this.signInForm.value;
        this.store.dispatch(AuthPageActions.signIn({ email, password }));
    }
}
