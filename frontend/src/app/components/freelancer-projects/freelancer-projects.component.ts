import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';

// Angular Material modules
import { MatTableModule } from '@angular/material/table';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatIconModule } from '@angular/material/icon';

import { AuthService } from '../../services/auth.service';
import { ProjectService } from '../../services/project.service';

@Component({
  selector: 'app-freelancer-projects',
  standalone: true,
  imports: [
    CommonModule,
    MatTableModule,
    MatButtonModule,
    MatCardModule,
    MatIconModule
  ],
  templateUrl: './freelancer-projects.component.html',
  styleUrls: ['./freelancer-projects.component.css']
})
export class FreelancerProjectsComponent implements OnInit {
  displayedColumns: string[] = ['title', 'budget', 'status', 'creationDate', 'actions'];
  dataSource: any[] = []; // This holds all the projects for the logged-in freelancer

  constructor(
    private authService: AuthService,
    private projectService: ProjectService,
    private router: Router
  ) {}

  ngOnInit(): void {
    const currentUser = this.authService.getCurrentUser();
    if (!currentUser || currentUser.role !== 'freelancer') {
      // A fallback check in case the route guard wasn't used
      this.router.navigate(['/login']);
      return;
    }
    this.loadMyProjects(currentUser.user_id);
  }

  loadMyProjects(freelancerId: number) {
    this.projectService.getAllProjects().subscribe({
      next: (res: any) => {
        const allProjects = res.projects || [];
        // Filter by freelancer_id
        this.dataSource = allProjects.filter(
          (p: any) => p.freelancer_id === freelancerId
        );
      },
      error: (err) => {
        console.error('Failed to fetch projects', err);
      }
    });
  }

  onUpdate(projectId: number) {
    // Navigate to an "Update Project" page or do something else
    // Possibly blocked if only the client can update, but placeholder here
    this.router.navigate(['/projects', projectId, 'edit']);
  }

  onDetails(projectId: number) {
    // Show project details
    this.router.navigate(['/projects', projectId]);
  }

  onViewOffers(projectId: number) {
    // Navigate to proposals or some "offers" page
    this.router.navigate([`/projects/${projectId}/proposals`]);
  }
}
