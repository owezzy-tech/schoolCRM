import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { SigninComponent } from './signin.component';
import { testProviders } from '../../testing/test-providers';
describe('SigninComponent', () => {
  let component: SigninComponent;
  let fixture: ComponentFixture<SigninComponent>;
  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      imports: [SigninComponent],
      providers: testProviders
    }).compileComponents();
  }));
  beforeEach(() => {
    fixture = TestBed.createComponent(SigninComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });
  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
