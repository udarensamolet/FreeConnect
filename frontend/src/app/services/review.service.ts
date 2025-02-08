import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../environments/environment';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class ReviewService {
  private apiUrl = environment.apiBaseUrl;

  constructor(private http: HttpClient) {}

  createReview(data: any): Observable<any> {
    return this.http.post(`${this.apiUrl}/reviews`, data);
  }

  getReviewById(reviewId: number): Observable<any> {
    return this.http.get(`${this.apiUrl}/reviews/${reviewId}`);
  }

  getReviewsByProject(projectId: number): Observable<any> {
    return this.http.get(`${this.apiUrl}/projects/${projectId}/reviews`);
  }

  updateReview(reviewId: number, data: any): Observable<any> {
    return this.http.put(`${this.apiUrl}/reviews/${reviewId}`, data);
  }

  deleteReview(reviewId: number): Observable<any> {
    return this.http.delete(`${this.apiUrl}/reviews/${reviewId}`);
  }
}