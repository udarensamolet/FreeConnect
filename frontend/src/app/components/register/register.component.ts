import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';

import { MatCardModule } from '@angular/material/card';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatSelectModule } from '@angular/material/select';

import { Router } from '@angular/router';
import { AuthService } from '../../services/auth.service';
import { AlertService } from '../../services/alert.service';

@Component({
  selector: 'app-register',
  templateUrl: './register.component.html',
  styleUrls: ['./register.component.css'],
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    MatCardModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    MatSelectModule
  ],
})
export class RegisterComponent {
  formData = {
    name: '',
    email: '',
    role: '',
    password: '',
    confirmPassword: ''
  };

  constructor(
    private authService: AuthService,
    private router: Router,
    private alertService: AlertService
  ) {}

  onSubmit() {
    if (!this.formData.name) {
      ('Name cannot be empty');
      return;
    }
    if (!this.formData.email) {
      this.alertService.alert('Email cannot be empty');
      return;
    }
    if (!this.formData.password) {
      this.alertService.alert('Password cannot be empty');
      return;
    }
    if (this.formData.password !== this.formData.confirmPassword) {
      this.alertService.alert('Passwords do not match!');
      return;
    }

    this.authService.register(this.formData).subscribe({
      next: (response) => {
        this.router.navigate(['/login']);
      },
      error: (error) => {
        this.alertService.alert('Registration failed!');
      }
    });
  }
}