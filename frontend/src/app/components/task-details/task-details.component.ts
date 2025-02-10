import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute, Router } from '@angular/router';

import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { TaskService } from '../../services/task.service';

@Component({
  selector: 'app-task-details',
  standalone: true,
  imports: [
    CommonModule,
    MatCardModule,
    MatButtonModule
  ],
  templateUrl: './task-details.component.html',
  styleUrls: ['./task-details.component.css']
})
export class TaskDetailsComponent implements OnInit {
  projectId!: number;
  taskId!: number;
  task: any;

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private taskService: TaskService
  ) {}

  ngOnInit(): void {
    // read :projectId and :taskId from URL
    this.projectId = Number(this.route.snapshot.paramMap.get('projectId'));
    this.taskId = Number(this.route.snapshot.paramMap.get('taskId'));

    this.loadTask();
  }

  loadTask(): void {
    this.taskService.getTaskById(this.taskId).subscribe({
      next: (res: any) => {
        this.task = res.task;
      },
      error: (err) => {
        console.error('Failed to load task:', err);
      }
    });
  }

  onBackToProject(): void {
    // Navigate back to the project details page
    this.router.navigate(['/projects', this.projectId]);
  }
}
