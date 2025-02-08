import { Injectable, NgZone } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../environments/environment';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class NotificationService {
  private apiUrl = environment.apiBaseUrl;

  constructor(private http: HttpClient, private zone: NgZone) {}

  createNotification(data: any): Observable<any> {
    return this.http.post(`${this.apiUrl}/notifications`, data);
  }

  getNotificationById(notifId: number): Observable<any> {
    return this.http.get(`${this.apiUrl}/notifications/${notifId}`);
  }

  getNotificationsByUser(userId: number): Observable<any> {
    return this.http.get(`${this.apiUrl}/notifications/user/${userId}`);
  }

  updateNotification(notifId: number, data: any): Observable<any> {
    return this.http.put(`${this.apiUrl}/notifications/${notifId}`, data);
  }

  deleteNotification(notifId: number): Observable<any> {
    return this.http.delete(`${this.apiUrl}/notifications/${notifId}`);
  }

  getSSEUpdates(): Observable<string> {
    const url = `${this.apiUrl.replace('/api', '')}/updates`;

    return new Observable(observer => {
      const eventSource = new EventSource(url);

      eventSource.onmessage = (event) => {
        this.zone.run(() => {
          observer.next(event.data);
        });
      };

      eventSource.onerror = (error) => {
        this.zone.run(() => {
          observer.error(error);
        });
      };

      return () => {
        eventSource.close();
      };
    });
  }

  broadcastMessage(message: string): Observable<any> {
    return this.http.post(`${this.apiUrl}/broadcast`, { message });
  }
}