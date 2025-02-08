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

  createTask(data: any): Observable<any> {
    return this.http.post(`${this.apiUrl}/tasks`, data);
  }

  getTaskById(taskId: number): Observable<any> {
    return this.http.get(`${this.apiUrl}/tasks/${taskId}`);
  }

  getTasksByProject(projectId: number): Observable<any> {
    return this.http.get(`${this.apiUrl}/projects/${projectId}/tasks`);
  }

  updateTask(taskId: number, data: any): Observable<any> {
    return this.http.put(`${this.apiUrl}/tasks/${taskId}`, data);
  }

  deleteTask(taskId: number): Observable<any> {
    return this.http.delete(`${this.apiUrl}/tasks/${taskId}`);
  }
}