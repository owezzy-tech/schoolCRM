import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { Dashboard2Component } from './dashboard2.component';
import { testProviders } from '../../../testing/test-providers';
describe('Dashboard2Component', () => {
  let component: Dashboard2Component;
  let fixture: ComponentFixture<Dashboard2Component>;
  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      imports: [Dashboard2Component],
      providers: testProviders
    }).compileComponents();
  }));
  beforeEach(() => {
    fixture = TestBed.createComponent(Dashboard2Component);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });
  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
