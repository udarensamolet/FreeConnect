import { Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../environments/environment';

@Injectable({
  providedIn: 'root'
})
export class ProjectService {
  private apiUrl = environment.apiBaseUrl;

  constructor(private http: HttpClient) { }

  createProject(projectData: any): Observable<any> {
    return this.http.post(`${this.apiUrl}/projects`, projectData);
  }

  getAllProjects(): Observable<any> {
    return this.http.get<any>(`${this.apiUrl}/projects`);
  }

  getProjects(search?: string, minBudget?: number, maxBudget?: number, status?: string): Observable<any> {
    let params = new HttpParams();
    if (search) {
      params = params.set('search', search);
    }
    if (minBudget !== undefined) {
      params = params.set('minBudget', minBudget.toString());
    }
    if (maxBudget !== undefined) {
      params = params.set('maxBudget', maxBudget.toString());
    }
    if (status) {
      params = params.set('status', status);
    }
  
    return this.http.get(`${this.apiUrl}/projects`, { params });
  }

  getProjectById(id: number): Observable<any> {
    return this.http.get(`${this.apiUrl}/projects/${id}`);
  }
}