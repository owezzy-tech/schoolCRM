import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { SidebarComponent } from './sidebar.component';
import { testProviders } from '../../testing/test-providers';
describe('SidebarComponent', () => {
  let component: SidebarComponent;
  let fixture: ComponentFixture<SidebarComponent>;
  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      imports: [SidebarComponent],
      providers: testProviders
    }).compileComponents();
  }));
  beforeEach(() => {
    fixture = TestBed.createComponent(SidebarComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });
  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
