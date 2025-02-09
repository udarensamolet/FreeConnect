import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../environments/environment';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class ProposalService {
  private apiUrl = environment.apiBaseUrl; // e.g. "http://localhost:8080/api"

  constructor(private http: HttpClient) {}

  getAllProposals(): Observable<any> {
    return this.http.get(`${this.apiUrl}/proposals`);
  }

  createProposal(data: any): Observable<any> {
    return this.http.post(`${this.apiUrl}/proposals`, data);
  }

  getProposalById(proposalId: number): Observable<any> {
    return this.http.get(`${this.apiUrl}/proposals/${proposalId}`);
  }

  updateProposal(proposalId: number, data: any): Observable<any> {
    return this.http.put(`${this.apiUrl}/proposals/${proposalId}`, data);
  }


  deleteProposal(proposalId: number): Observable<any> {
    return this.http.delete(`${this.apiUrl}/proposals/${proposalId}`);
  }
}