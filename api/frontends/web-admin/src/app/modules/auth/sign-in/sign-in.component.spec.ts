import { ComponentFixture, TestBed } from '@angular/core/testing';
import { Actions } from '@ngrx/effects';
import { Store } from '@ngrx/store';
import { provideRouter } from '@angular/router';
import { Subject } from 'rxjs';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { AuthPageActions } from 'app/core/auth/+state';
import { AuthSignInComponent } from './sign-in.component';

describe('AuthSignInComponent deployment credentials', () => {
    let fixture: ComponentFixture<AuthSignInComponent>;
    const dispatch = vi.fn();

    beforeEach(async () => {
        dispatch.mockClear();
        await TestBed.configureTestingModule({
            imports: [AuthSignInComponent],
            providers: [
                provideRouter([]),
                { provide: Store, useValue: { dispatch } },
                { provide: Actions, useValue: new Subject() },
            ],
        }).overrideComponent(AuthSignInComponent, {
            set: { template: '', imports: [] },
        }).compileComponents();
        fixture = TestBed.createComponent(AuthSignInComponent);
        fixture.autoDetectChanges();
        await fixture.whenStable();
    });

    it('starts with empty credentials and refuses a blank submission', () => {
        const component = fixture.componentInstance;
        expect(component.signInForm.get('email').value).toBe('');
        expect(component.signInForm.get('password').value).toBe('');
        component.signIn();
        expect(dispatch).not.toHaveBeenCalled();
    });

    it('dispatches only the credentials supplied by the user', () => {
        const component = fixture.componentInstance;
        const credentials = {
            email: 'person@schoolcrm.invalid',
            password: 'UserSuppliedFixtureOnly123!',
        };
        component.signInForm.patchValue(credentials);
        component.signIn();
        expect(dispatch).toHaveBeenCalledExactlyOnceWith(
            AuthPageActions.signIn(credentials)
        );
    });
});
