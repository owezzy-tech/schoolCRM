import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { LockedComponent } from './locked.component';
import { testProviders } from '../../testing/test-providers';
describe('LockedComponent', () => {
  let component: LockedComponent;
  let fixture: ComponentFixture<LockedComponent>;
  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      imports: [LockedComponent],
      providers: testProviders
    }).compileComponents();
  }));
  beforeEach(() => {
    fixture = TestBed.createComponent(LockedComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });
  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
