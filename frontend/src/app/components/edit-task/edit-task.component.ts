import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';

import { MatCardModule } from '@angular/material/card';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatSelectModule } from '@angular/material/select';

import { TaskService } from '../../services/task.service';

@Component({
  selector: 'app-edit-task',
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
  templateUrl: './edit-task.component.html',
  styleUrls: ['./edit-task.component.css']
})
export class EditTaskComponent implements OnInit {
  projectId!: number;
  taskId!: number;

  // We'll store form data here
  formData = {
    title: '',
    description: '',
    deadline: '',
    budget: 0,
    status: 'open'
  };

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private taskService: TaskService
  ) {}

  ngOnInit(): void {
    this.projectId = Number(this.route.snapshot.paramMap.get('projectId'));
    this.taskId = Number(this.route.snapshot.paramMap.get('taskId'));

    // Load existing task data
    this.loadTask();
  }

  loadTask(): void {
    this.taskService.getTaskById(this.taskId).subscribe({
      next: (res: any) => {
        const task = res.task;
        // Patch formData
        this.formData.title = task.title;
        this.formData.description = task.description;
        // Convert the existing deadline to YYYY-MM-DD if needed
        // For example, if it's an ISO string:
        const dateObj = new Date(task.deadline);
        this.formData.deadline = dateObj.toISOString().split('T')[0];

        this.formData.budget = task.budget;
        this.formData.status = task.status || 'open';
      },
      error: (err) => {
        console.error('Failed to load task:', err);
      }
    });
  }

  onSubmit(): void {
    // Convert the deadline back to an ISO string if needed
    const isoDeadline = new Date(this.formData.deadline).toISOString();

    const payload = {
      title: this.formData.title,
      description: this.formData.description,
      deadline: isoDeadline,
      budget: this.formData.budget,
      status: this.formData.status
    };

    this.taskService.updateTask(this.taskId, payload).subscribe({
      next: () => {
        // Once updated, navigate back to project details
        this.router.navigate(['/projects', this.projectId]);
      },
      error: (err) => {
        console.error('Failed to update task:', err);
      }
    });
  }
}