import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';

import { MatCardModule } from '@angular/material/card';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';

import { AuthService } from '../../services/auth.service';
import { ProjectService } from '../../services/project.service';

@Component({
  selector: 'app-edit-project',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    MatCardModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule
  ],
  templateUrl: './edit-project.component.html',
  styleUrls: ['./edit-project.component.css']
})
export class EditProjectComponent implements OnInit {
  projectId!: number;
  formData = {
    title: '',
    description: '',
    budget: 0,
    duration: 0
  };

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private authService: AuthService,
    private projectService: ProjectService
  ) {}

  ngOnInit(): void {
    // 1) Get the :id from route
    this.projectId = Number(this.route.snapshot.paramMap.get('id'));

    // 2) Fetch the project from the API
    this.loadProject();
  }

  loadProject(): void {
    this.projectService.getProjectById(this.projectId).subscribe({
      next: (res: any) => {
        const project = res.project;
        // Patch the formData with existing project info
        this.formData.title = project.title;
        this.formData.description = project.description;
        this.formData.budget = project.budget;
        this.formData.duration = project.duration;
      },
      error: (err) => {
        console.error('Failed to load project:', err);
      }
    });
  }

  onSubmit(): void {
    // Build payload with updated data
    const payload = {
      title: this.formData.title,
      description: this.formData.description,
      budget: this.formData.budget,
      duration: this.formData.duration
    };

    this.projectService.updateProject(this.projectId, payload).subscribe({
      next: () => {
        // After successful update, navigate somewhere (e.g. /my-projects)
        this.router.navigate(['/my-projects']);
      },
      error: (err) => {
        console.error('Update failed:', err);
      }
    });
  }
}