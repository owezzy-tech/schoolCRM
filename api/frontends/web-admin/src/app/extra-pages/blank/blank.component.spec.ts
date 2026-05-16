import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { BlankComponent } from './blank.component';
import { testProviders } from '../../testing/test-providers';
describe('BlankComponent', () => {
  let component: BlankComponent;
  let fixture: ComponentFixture<BlankComponent>;
  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      imports: [BlankComponent],
      providers: testProviders
    }).compileComponents();
  }));
  beforeEach(() => {
    fixture = TestBed.createComponent(BlankComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });
  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
