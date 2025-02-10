import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { MatCardModule } from '@angular/material/card';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';

import { Router } from '@angular/router';
import { AuthService } from '../../services/auth.service';
import { ProjectService } from '../../services/project.service';

@Component({
  selector: 'app-post-project',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    MatCardModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule
  ],
  templateUrl: './post-project.component.html',
  styleUrls: ['./post-project.component.css']
})
export class PostProjectComponent implements OnInit {
  formData = {
    title: '',
    description: '',
    budget: 0,
    duration: 0
  };

  user: any;  // We'll store the current user (client) here

  constructor(
    private authService: AuthService,
    private projectService: ProjectService,
    private router: Router
  ) {}

  ngOnInit(): void {
    // Get current user from AuthService
    this.user = this.authService.getCurrentUser();
    // If user is not client, they shouldn't be here, but we rely on the guard anyway
  }

  onSubmit(): void {
    // Build the payload
    const payload = {
      ...this.formData,
      client_id: this.user.user_id
      // status defaults to "open" in your backend
      // creation_date is handled by your backend's default
    };

    this.projectService.createProject(payload).subscribe({
      next: (response) => {
        // After success, maybe navigate to "my-projects"
        this.router.navigate(['/my-projects']);
      },
      error: (err) => {
        console.error('Failed to create project:', err);
      }
    });
  }
}