import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { First2Component } from './first2.component';

import { testProviders } from '../../testing/test-providers';
describe('First2Component', () => {
  let component: First2Component;
  let fixture: ComponentFixture<First2Component>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      imports: [First2Component],
      providers: testProviders
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(First2Component);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
