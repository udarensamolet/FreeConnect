import { Component, OnInit, OnDestroy } from '@angular/core';
import { Router, RouterModule } from '@angular/router';
import { CommonModule } from '@angular/common';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatButtonModule } from '@angular/material/button';
import { Subscription } from 'rxjs';
import { AuthService } from '../../services/auth.service';
import { AlertService } from '../../services/alert.service';

@Component({
  selector: 'app-navbar',
  templateUrl: './nav-bar.component.html',
  styleUrls: ['./nav-bar.component.css'],
  standalone: true,
  imports: [CommonModule, RouterModule, MatToolbarModule, MatButtonModule]
})
export class NavbarComponent implements OnInit, OnDestroy {
  isLoggedIn: boolean = false;
  private authSubscription!: Subscription;

  constructor(private authService: AuthService, private router: Router, private alertService: AlertService) {}

  ngOnInit(): void {
    // Subscribe to the loggedIn$ observable to reactively update the state.
    this.authSubscription = this.authService.loggedIn$.subscribe(
      (loggedIn) => (this.isLoggedIn = loggedIn)
    );
  }

  ngOnDestroy(): void {
    // Clean up the subscription.
    if (this.authSubscription) {
      this.authSubscription.unsubscribe();
    }
  }

  onLogout(): void {
    this.authService.logout().subscribe({
      next: () => {
        this.router.navigate(['/login']);
      },
      error: (error) => {
        this.alertService.alert(error.toString());
      }
    });
  }
}
