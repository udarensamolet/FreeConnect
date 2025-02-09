import { Inject, Injectable, PLATFORM_ID } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../environments/environment';
import { Observable, BehaviorSubject } from 'rxjs';
import { tap } from 'rxjs/operators';
import { isPlatformBrowser } from '@angular/common';

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  private apiUrl = environment.apiBaseUrl;

  private loggedInSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(this.hasToken());
  public loggedIn$ = this.loggedInSubject.asObservable();

  constructor(@Inject(PLATFORM_ID) private platformId: Object, private http: HttpClient) {}

  private hasToken(): boolean {
    return isPlatformBrowser(this.platformId) && !!localStorage.getItem('accessToken');
  }

  // 1) GET CURRENT USER
  public getCurrentUser(): any {
    const userStr = localStorage.getItem('currentUser');
    return userStr ? JSON.parse(userStr) : null;
  }
  

  // 2) LOGIN
  login(credentials: { email: string; password: string }): Observable<any> {
    return this.http.post(`${this.apiUrl}/login`, credentials).pipe(
      tap((response: any) => {
        if (response && response.token) {
          this.setToken(response.token);
        }
        if (response && response.user) {
          localStorage.setItem('currentUser', JSON.stringify(response.user));
        }
      })
    );
  }

  // 3) REGISTER 
  register(data: any): Observable<any> {
    return this.http.post(`${this.apiUrl}/register`, data);
  }

  // 4) LOGOUT
  logout(): Observable<any> {
    return this.http.post(`${this.apiUrl}/logout`, {}).pipe(
      tap(() => {
        this.removeToken();
        this.removeUser();
      })
    );
  }

  // 5) SET TOKEN
  setToken(token: string) {
    localStorage.setItem('accessToken', token);
    this.loggedInSubject.next(true);
  }

  // 6) GET TOKEN
  getToken(): string | null {
    if (isPlatformBrowser(this.platformId)) {
      return localStorage.getItem('accessToken');
    }
    return null;
  }

  // 7) REMOVE TOKEN
  removeToken() {
    localStorage.removeItem('accessToken');
    this.loggedInSubject.next(false);
  }

  // 8) REMOVE USER
  removeUser() {
    localStorage.removeItem('currentUser');
  }
}