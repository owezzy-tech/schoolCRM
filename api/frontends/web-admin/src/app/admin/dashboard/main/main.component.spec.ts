import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { MainComponent } from './main.component';
import { testProviders } from '../../../testing/test-providers';
describe('MainComponent', () => {
  let component: MainComponent;
  let fixture: ComponentFixture<MainComponent>;
  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      imports: [MainComponent],
      providers: testProviders
    }).compileComponents();
  }));
  beforeEach(() => {
    fixture = TestBed.createComponent(MainComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });
  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
