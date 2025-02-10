import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../environments/environment';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class TaskService {
  private apiUrl = environment.apiBaseUrl;

  constructor(private http: HttpClient) {}

  createTask(taskData: any): Observable<any> {
    return this.http.post(`${this.apiUrl}/tasks`, taskData);
  }
  
  getTasksByProject(projectId: number): Observable<any> {
    return this.http.get(`${this.apiUrl}/projects/${projectId}/tasks`);
  }

  deleteTask(taskId: number): Observable<any> {
    return this.http.delete(`${this.apiUrl}/tasks/${taskId}`);
  }

  getTaskById(taskId: number): Observable<any> {
    return this.http.get(`${this.apiUrl}/tasks/${taskId}`);
  }
  
  updateTask(taskId: number, taskData: any): Observable<any> {
    return this.http.put(`${this.apiUrl}/tasks/${taskId}`, taskData);
  }
}