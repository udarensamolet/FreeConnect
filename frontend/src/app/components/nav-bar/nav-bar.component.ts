import { Component, OnInit, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router, RouterModule } from '@angular/router';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatButtonModule } from '@angular/material/button';
import { Subscription } from 'rxjs';
import { AuthService } from '../../services/auth.service';
import { AlertService } from '../../services/alert.service';

@Component({
  selector: 'app-navbar',
  standalone: true,
  imports: [
    CommonModule,
    RouterModule,
    MatToolbarModule,
    MatButtonModule
  ],
  templateUrl: './nav-bar.component.html',
  styleUrls: ['./nav-bar.component.css']
})
export class NavbarComponent implements OnInit, OnDestroy {
  isLoggedIn = false;
  userRole: string | null = null;
  private authSubscription!: Subscription;

  constructor(
    private authService: AuthService,
    private router: Router,
    private alertService: AlertService
  ) {}

  ngOnInit(): void {
    this.authSubscription = this.authService.loggedIn$.subscribe(
      (loggedIn) => {
        this.isLoggedIn = loggedIn;
        if (this.isLoggedIn) {
          const currentUser = this.authService.getCurrentUser();
          this.userRole = currentUser ? currentUser.role : null;
        } else {
          this.userRole = null;
        }
      }
    );
  }

  ngOnDestroy(): void {
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