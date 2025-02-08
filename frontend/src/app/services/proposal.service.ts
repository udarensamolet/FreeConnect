import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../environments/environment';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class ProposalService {
  private apiUrl = environment.apiBaseUrl;  

  constructor(private http: HttpClient) {}

  createProposal(data: any): Observable<any> {
    return this.http.post(`${this.apiUrl}/proposals`, data);
  }

  getProposalById(proposalId: number): Observable<any> {
    return this.http.get(`${this.apiUrl}/proposals/${proposalId}`);
  }

  getProposalsByProject(projectId: number): Observable<any> {
    return this.http.get(`${this.apiUrl}/projects/${projectId}/proposals`);
  }

  updateProposal(proposalId: number, data: any): Observable<any> {
    return this.http.put(`${this.apiUrl}/proposals/${proposalId}`, data);
  }

  deleteProposal(proposalId: number): Observable<any> {
    return this.http.delete(`${this.apiUrl}/proposals/${proposalId}`);
  }

  acceptProposal(proposalId: number): Observable<any> {
    return this.http.post(`${this.apiUrl}/proposals/${proposalId}/accept`, {});
  }
}