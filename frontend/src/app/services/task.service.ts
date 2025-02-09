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

  createTask(projectId: number, data: any): Observable<any> {
    return this.http.post(`${this.apiUrl}/projects/${projectId}/tasks`, data);
  }

  getTasksByProject(projectId: number): Observable<any> {
    return this.http.get(`${this.apiUrl}/projects/${projectId}/tasks`);
  }

  getTaskById(projectId: number, taskId: number): Observable<any> {
    return this.http.get(`${this.apiUrl}/projects/${projectId}/tasks/${taskId}`);
  }

  updateTask(projectId: number, taskId: number, data: any): Observable<any> {
    return this.http.put(`${this.apiUrl}/projects/${projectId}/tasks/${taskId}`, data);
  }

  deleteTask(projectId: number, taskId: number): Observable<any> {
    return this.http.delete(`${this.apiUrl}/projects/${projectId}/tasks/${taskId}`);
  }
}