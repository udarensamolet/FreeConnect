import { Component, OnInit } from '@angular/core';
import { Router, RouterModule } from '@angular/router';
import { CommonModule } from '@angular/common';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatButtonModule } from '@angular/material/button';
import { AuthService } from '../../services/auth.service';
import { AlertService } from '../../services/alert.service';

@Component({
  selector: 'app-navbar',
  templateUrl: './navbar.component.html',
  styleUrls: ['./navbar.component.css'],
  standalone: true,
  imports: [CommonModule, RouterModule, MatToolbarModule, MatButtonModule]
})
export class NavbarComponent implements OnInit {
  isLoggedIn: boolean = false;

  constructor(private authService: AuthService, private router: Router, private alertService: AlertService) { }

  ngOnInit(): void {
    this.isLoggedIn = !!this.authService.getToken();
  }

  onLogout(): void {
    this.authService.logout().subscribe({
      next: () => {
        this.authService.removeToken();
        this.isLoggedIn = false;
        this.router.navigate(['/login']);
      },
      error: (error) => {
        this.alertService.alert(error.toString());
      }
    });
  }
}