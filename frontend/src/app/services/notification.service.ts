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

  getNotifications(userId: number): Observable<any> {
    return this.http.get(`${this.apiUrl}/notifications?userId=${userId}`);
  }

  // SSE at GET /api/updates
  getSSEUpdates(): Observable<string> {
    const url = `${this.apiUrl}/updates`; 

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
}