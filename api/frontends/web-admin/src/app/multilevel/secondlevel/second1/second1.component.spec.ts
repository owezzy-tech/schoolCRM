import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { Second1Component } from './second1.component';

import { testProviders } from '../../../testing/test-providers';
describe('Second1Component', () => {
  let component: Second1Component;
  let fixture: ComponentFixture<Second1Component>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      imports: [Second1Component],
      providers: testProviders
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(Second1Component);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
