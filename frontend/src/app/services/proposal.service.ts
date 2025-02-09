import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../environments/environment';

@Injectable({ providedIn: 'root' })
export class ProposalService {
  private apiUrl = environment.apiBaseUrl;

  constructor(private http: HttpClient) {}

  createProposal(proposalData: any) {
    return this.http.post(`${this.apiUrl}/proposals`, proposalData);
  }

  getProposalsByProject(projectId: number) {
    return this.http.get(`${this.apiUrl}/projects/${projectId}/proposals`);
  }
}