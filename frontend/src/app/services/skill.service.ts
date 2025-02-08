import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../environments/environment';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class SkillService {
  private apiUrl = environment.apiBaseUrl; 

  constructor(private http: HttpClient) {}

  createSkill(data: any): Observable<any> {
    return this.http.post(`${this.apiUrl}/skills`, data);
  }

  getAllSkills(): Observable<any> {
    return this.http.get(`${this.apiUrl}/skills`);
  }

  getSkillById(skillId: number): Observable<any> {
    return this.http.get(`${this.apiUrl}/skills/${skillId}`);
  }

  updateSkill(skillId: number, data: any): Observable<any> {
    return this.http.put(`${this.apiUrl}/skills/${skillId}`, data);
  }

  deleteSkill(skillId: number): Observable<any> {
    return this.http.delete(`${this.apiUrl}/skills/${skillId}`);
  }
}