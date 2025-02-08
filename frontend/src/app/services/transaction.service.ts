import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../environments/environment';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class TransactionService {
  private apiUrl = environment.apiBaseUrl;  

  constructor(private http: HttpClient) {}

  createTransaction(data: any): Observable<any> {
    return this.http.post(`${this.apiUrl}/transactions`, data);
  }

  getTransactionById(transactionId: number): Observable<any> {
    return this.http.get(`${this.apiUrl}/transactions/${transactionId}`);
  }

  getTransactionsByProject(projectId: number): Observable<any> {
    return this.http.get(`${this.apiUrl}/projects/${projectId}/transactions`);
  }

  updateTransaction(transactionId: number, data: any): Observable<any> {
    return this.http.put(`${this.apiUrl}/transactions/${transactionId}`, data);
  }

  deleteTransaction(transactionId: number): Observable<any> {
    return this.http.delete(`${this.apiUrl}/transactions/${transactionId}`);
  }
}