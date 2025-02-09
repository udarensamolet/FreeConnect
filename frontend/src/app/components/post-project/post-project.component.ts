// post-project.component.ts
import { Component } from '@angular/core';
import { ProjectService } from '../../services/project.service';
import { AlertService } from '../../services/alert.service';

@Component({
  selector: 'app-post-project',
  templateUrl: './post-project.component.html'
})
export class PostProjectComponent {
  title = '';
  description = '';
  budget = 0;
  duration = 0;

  constructor(private projectService: ProjectService, private alertService: AlertService) {}

  onSubmit() {
    const projectData = {
      title: this.title,
      description: this.description,
      budget: this.budget,
      duration: this.duration
    };
    this.projectService.createProject(projectData).subscribe({
      next: (res) => {
        alert('Project created successfully!');
      },
      error: (err) => {
        console.error(err);
        alert('Failed to create project');
      }
    });
  }
}
