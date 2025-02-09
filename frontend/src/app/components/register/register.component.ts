import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';

import { MatCardModule } from '@angular/material/card';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatSelectModule } from '@angular/material/select';

import { AuthService } from '../../services/auth.service';

@Component({
  selector: 'app-register',
  templateUrl: './register.component.html',
  styleUrls: ['./register.component.css',
  ],
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

  constructor(private authService: AuthService) {}

  onSubmit() {
    if (!this.formData.name) {
      alert('Name cannot be empty');
      return;
    }
    if (!this.formData.email) {
      alert('Email cannot be empty');
      return;
    }
    if (!this.formData.password) {
      alert('Password cannot be empty');
      return;
    }
    if (this.formData.password !== this.formData.confirmPassword) {
      alert('Passwords do not match!');
      return;
    }

    this.authService.register(this.formData).subscribe({
      next: (response) => {
        alert('Registration successful!');
        console.log(response);
      },
      error: (error) => {
        alert('Registration failed!');
        console.error(error);
      }
    });
  }
}