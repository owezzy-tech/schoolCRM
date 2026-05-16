import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { Page404Component } from './page404.component';
import { testProviders } from '../../testing/test-providers';
describe('Page404Component', () => {
  let component: Page404Component;
  let fixture: ComponentFixture<Page404Component>;
  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      imports: [Page404Component],
      providers: testProviders
    }).compileComponents();
  }));
  beforeEach(() => {
    fixture = TestBed.createComponent(Page404Component);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });
  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
