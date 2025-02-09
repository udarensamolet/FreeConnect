import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../environments/environment';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class ProjectService {
  private apiUrl = environment.apiBaseUrl;

  constructor(private http: HttpClient) {}

  getAllProjects(): Observable<any> {
    return this.http.get(`${this.apiUrl}/projects`);
  }

  getProject(projectId: number): Observable<any> {
    return this.http.get(`${this.apiUrl}/projects/${projectId}`);
  }

  createProject(projectData: any): Observable<any> {
    return this.http.post(`${this.apiUrl}/projects`, projectData);
  }

  updateProject(projectId: number, projectData: any): Observable<any> {
    return this.http.put(`${this.apiUrl}/projects/${projectId}`, projectData);
  }

  deleteProject(projectId: number): Observable<any> {
    return this.http.delete(`${this.apiUrl}/projects/${projectId}`);
  }
}