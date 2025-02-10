import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatCardModule } from '@angular/material/card';
import { AuthService } from '../../services/auth.service';
import { ProjectService } from '../../services/project.service';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [
    CommonModule,
    MatCardModule
  ],
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.css']
})
export class HomeComponent implements OnInit {
  isLoggedIn = false;
  user: any = null;  // { user_id, role, name, etc. }
  projects: any[] = [];  // relevant projects for client/freelancer
  welcomeMessage = '';    // "Welcome, John Doe!"
  noProjectsMessage = 'You don\'t have any created/assigned projects at this moment.';

  constructor(
    private authService: AuthService,
    private projectService: ProjectService
  ) {}

  ngOnInit(): void {
    this.user = this.authService.getCurrentUser();
    this.isLoggedIn = !!this.user;

    if (this.isLoggedIn) {
      // Example: "Welcome, John Doe!"
      this.welcomeMessage = `Welcome, ${this.user.name || 'User'}!`;

      // If user is client/freelancer, load projects
      if (this.user.role === 'client' || this.user.role === 'freelancer') {
        this.loadProjects();
      }
      // If user is admin, no projects to show for now
    }
  }

  private loadProjects(): void {
    this.projectService.getAllProjects().subscribe({
      next: (response: any) => {
        // Suppose "response.projects" is the list
        const allProjects = response.projects || [];

        if (this.user.role === 'client') {
          // Filter where client_id == user_id
          this.projects = allProjects.filter(
            (p: any) => p.client_id === this.user.user_id
          );
        } else if (this.user.role === 'freelancer') {
          // Filter where freelancer_id == user_id
          this.projects = allProjects.filter(
            (p: any) => p.freelancer_id === this.user.user_id
          );
        }
      },
      error: (err: any) => {
        console.error('Failed to load projects', err);
      }
    });
  }
}
