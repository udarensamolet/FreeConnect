import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute, Router } from '@angular/router';

import { MatCardModule } from '@angular/material/card';
import { MatTableModule } from '@angular/material/table';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';

import { ProjectService } from '../../services/project.service';
import { TaskService } from '../../services/task.service';

@Component({
  selector: 'app-project-details',
  standalone: true,
  imports: [
    CommonModule,
    MatCardModule,
    MatTableModule,
    MatButtonModule,
    MatIconModule
  ],
  templateUrl: './project-details.component.html',
  styleUrls: ['./project-details.component.css']
})
export class ProjectDetailsComponent implements OnInit {
  projectId!: number;
  project: any;
  tasks: any[] = [];

  displayedColumns: string[] = ['title', 'deadline', 'status', 'budget', 'actions'];

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private projectService: ProjectService,
    private taskService: TaskService
  ) {}

  ngOnInit(): void {
    this.projectId = Number(this.route.snapshot.paramMap.get('id'));
    this.loadProject();
    this.loadTasks();
  }

  loadProject(): void {
    this.projectService.getProjectById(this.projectId).subscribe({
      next: (res: any) => {
        this.project = res.project;
      },
      error: (err) => {
        console.error('Failed to fetch project details:', err);
      }
    });
  }

  loadTasks(): void {
    this.taskService.getTasksByProject(this.projectId).subscribe({
      next: (res: any) => {
        this.tasks = res.tasks || [];
      },
      error: (err) => {
        console.error('Failed to fetch tasks:', err);
      }
    });
  }

  onAddTask(): void {
    this.router.navigate(['/projects/', this.projectId, 'create-task']);
  }

  onEditTask(taskId: number): void {
    this.router.navigate([
      '/projects',
      this.projectId,
      'tasks',
      taskId,
      'edit'
    ]);
  }

  onTaskDetails(taskId: number): void {
    this.router.navigate(['/projects', this.projectId, 'tasks', taskId]);
  }
  
  onDeleteTask(taskId: number): void {
    console.log('Delete Task:', taskId);
    //this.loadTasks();
  }
}