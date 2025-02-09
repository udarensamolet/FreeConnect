import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../environments/environment';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class InvoiceService {
  private apiUrl = environment.apiBaseUrl;  

  constructor(private http: HttpClient) {}

  createInvoice(data: any): Observable<any> {
    return this.http.post(`${this.apiUrl}/invoices`, data);
  }

  getInvoiceById(invoiceId: number): Observable<any> {
    return this.http.get(`${this.apiUrl}/invoices/${invoiceId}`);
  }

  getInvoicesByProject(projectId: number): Observable<any> {
    return this.http.get(`${this.apiUrl}/projects/${projectId}/invoices`);
  }

  updateInvoice(invoiceId: number, data: any): Observable<any> {
    return this.http.put(`${this.apiUrl}/invoices/${invoiceId}`, data);
  }

  deleteInvoice(invoiceId: number): Observable<any> {
    return this.http.delete(`${this.apiUrl}/invoices/${invoiceId}`);
  }
}