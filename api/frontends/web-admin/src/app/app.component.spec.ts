import { Component } from '@angular/core';
import { TestBed } from '@angular/core/testing';
import { provideRouter, RouterOutlet } from '@angular/router';
import { beforeEach, describe, expect, it } from 'vitest';

import { AppComponent } from './app.component';

@Component({
    selector: 'router-outlet',
    template: '',
})
class RouterOutletStubComponent {}

describe('AppComponent', () => {
    beforeEach(async () => {
        await TestBed.configureTestingModule({
            imports: [AppComponent],
            providers: [provideRouter([])],
        })
            .overrideComponent(AppComponent, {
                remove: {
                    imports: [RouterOutlet],
                },
                add: {
                    imports: [RouterOutletStubComponent],
                },
            })
            .compileComponents();
    });

    it('creates the app component', () => {
        const fixture = TestBed.createComponent(AppComponent);

        expect(fixture.componentInstance).toBeTruthy();
    });
});
