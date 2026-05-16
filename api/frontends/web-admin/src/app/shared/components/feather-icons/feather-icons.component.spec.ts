import { ComponentFixture, TestBed } from '@angular/core/testing';

import { FeatherIconsComponent } from './feather-icons.component';

import { testProviders } from '../../../testing/test-providers';
describe('FeatherIconsComponent', () => {
  let component: FeatherIconsComponent;
  let fixture: ComponentFixture<FeatherIconsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [FeatherIconsComponent],
      providers: testProviders
    })
      .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(FeatherIconsComponent);
    component = fixture.componentInstance;
    component.icon = 'home';
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
