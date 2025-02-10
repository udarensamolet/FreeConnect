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

  createProject(projectData: any): Observable<any> {
    return this.http.post(`${this.apiUrl}/projects`, projectData);
  }

  updateProject(projectId: number, data: any): Observable<any> {
    return this.http.put(`${this.apiUrl}/projects/${projectId}`, data);
  }
  
  getProjectById(id: number): Observable<any> {
    return this.http.get(`${this.apiUrl}/projects/${id}`);
  }  
}