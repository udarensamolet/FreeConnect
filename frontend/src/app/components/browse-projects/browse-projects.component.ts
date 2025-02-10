import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';

// Angular Material modules
import { MatTableModule } from '@angular/material/table';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatIconModule } from '@angular/material/icon';

import { ProjectService } from '../../services/project.service';

@Component({
  selector: 'app-browse-projects',
  standalone: true,
  imports: [
    CommonModule,
    MatTableModule,
    MatButtonModule,
    MatCardModule,
    MatIconModule
  ],
  templateUrl: './browse-projects.component.html',
  styleUrls: ['./browse-projects.component.css']
})
export class BrowseUnassignedProjectsComponent implements OnInit {
  displayedColumns: string[] = ['title', 'budget', 'status', 'creationDate', 'actions'];
  dataSource: any[] = [];

  constructor(
    private projectService: ProjectService,
    private router: Router
  ) {}

  ngOnInit(): void {
    this.loadUnassignedProjects();
  }

  loadUnassignedProjects(): void {
    this.projectService.getAllProjects().subscribe({
      next: (res: any) => {
        const allProjects = res.projects || [];
        // Filter out only "unassigned" projects
        // If your backend sets `freelancer_id` to null, do:
        this.dataSource = allProjects.filter((p: any) => !p.freelancer_id);
        
        // Alternatively, you might check `p.status === 'open'`
        // or combine both conditions as you like:
        // .filter((p: any) => !p.freelancer_id && p.status === 'open');
      },
      error: (err) => {
        console.error('Failed to fetch unassigned projects', err);
      }
    });
  }

  onDetails(projectId: number) {
    // Navigate to project details
    this.router.navigate(['/projects', projectId]);
  }
}
