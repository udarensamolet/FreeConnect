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

  getAllTransactions(): Observable<any> {
    return this.http.get(`${this.apiUrl}/transactions`);
  }

  createTransaction(data: any): Observable<any> {
    return this.http.post(`${this.apiUrl}/transactions`, data);
  }

  getTransactionById(transactionId: number): Observable<any> {
    return this.http.get(`${this.apiUrl}/transactions/${transactionId}`);
  }
}
