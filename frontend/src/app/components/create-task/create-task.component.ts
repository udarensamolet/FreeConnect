import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';

import { MatCardModule } from '@angular/material/card';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';

import { TaskService } from '../../services/task.service';

@Component({
  selector: 'app-create-task',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    MatCardModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule
  ],
  templateUrl: './create-task.component.html',
  styleUrls: ['./create-task.component.css']
})
export class CreateTaskComponent implements OnInit {
  projectId!: number;

  formData = {
    title: '',
    description: '',
    deadline: '',
    budget: 0
  };

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private taskService: TaskService
  ) {}

  ngOnInit(): void {
    // Retrieve :id from the URL
    this.projectId = Number(this.route.snapshot.paramMap.get('id'));
  }

  onSubmit(): void {
    // Convert a YYYY-MM-DD string into an ISO string
    const isoDeadline = new Date(this.formData.deadline).toISOString();

    const payload = {
      title: this.formData.title,
      description: this.formData.description,
      deadline: isoDeadline,
      budget: this.formData.budget,
      project_id: this.projectId
    };

    this.taskService.createTask(payload).subscribe({
      next: () => {
        // Navigate back to the project details after creating the task
        this.router.navigate(['/projects', this.projectId]);
      },
      error: (err) => {
        console.error('Failed to create task:', err);
      }
    });
  }
}