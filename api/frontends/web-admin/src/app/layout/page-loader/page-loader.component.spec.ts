import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { PageLoaderComponent } from './page-loader.component';
import { testProviders } from '../../testing/test-providers';
describe('PageLoaderComponent', () => {
  let component: PageLoaderComponent;
  let fixture: ComponentFixture<PageLoaderComponent>;
  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      imports: [PageLoaderComponent],
      providers: testProviders
    }).compileComponents();
  }));
  beforeEach(() => {
    fixture = TestBed.createComponent(PageLoaderComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });
  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
