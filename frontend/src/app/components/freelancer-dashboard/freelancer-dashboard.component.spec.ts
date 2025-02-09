import { ComponentFixture, TestBed } from '@angular/core/testing';

import { FreelancerDashboardComponent } from './freelancer-dashboard.component';

describe('FreelancerDashboardComponent', () => {
  let component: FreelancerDashboardComponent;
  let fixture: ComponentFixture<FreelancerDashboardComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [FreelancerDashboardComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(FreelancerDashboardComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
