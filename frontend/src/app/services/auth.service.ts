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

  // This subject holds the current user object (or null if no user).
  private userSubject: BehaviorSubject<any | null>;
  // Public observable for any interested subscriber (e.g., Navbar).
  public user$: Observable<any | null>;

  constructor(
    @Inject(PLATFORM_ID) private platformId: Object,
    private http: HttpClient
  ) {
    // On startup, read 'currentUser' from localStorage if in the browser.
    let initialUser: any | null = null;
    if (isPlatformBrowser(this.platformId)) {
      const userStr = localStorage.getItem('currentUser');
      initialUser = userStr ? JSON.parse(userStr) : null;
    }

    // Initialize the BehaviorSubject with the user we found (or null).
    this.userSubject = new BehaviorSubject<any | null>(initialUser);
    this.user$ = this.userSubject.asObservable();
  }

  // 1) GET CURRENT USER (sync)
  public getCurrentUser(): any | null {
    // The BehaviorSubject always has the current user in memory.
    return this.userSubject.value;
  }

  // 2) LOGIN
  login(credentials: { email: string; password: string }): Observable<any> {
    return this.http.post(`${this.apiUrl}/login`, credentials).pipe(
      tap((response: any) => {
        // If success, store token and user in localStorage (browser only).
        if (response && response.token) {
          this.setToken(response.token);
        }
        if (response && response.user && isPlatformBrowser(this.platformId)) {
          localStorage.setItem('currentUser', JSON.stringify(response.user));
          // Also push user into the BehaviorSubject
          this.userSubject.next(response.user);
        }
      })
    );
  }

  // 3) REGISTER
  register(data: any): Observable<any> {
    return this.http.post(`${this.apiUrl}/register`, data);
  }

  // 4) LOGOUT (no backend call)
  logout(): Observable<void> {
    return new Observable<void>((observer) => {
      this.removeToken();
      this.removeUser();
      observer.next();
      observer.complete();
    });
  }

  // 5) SET TOKEN
  setToken(token: string): void {
    if (isPlatformBrowser(this.platformId)) {
      localStorage.setItem('accessToken', token);
    }
  }

  // 6) GET TOKEN
  getToken(): string | null {
    if (isPlatformBrowser(this.platformId)) {
      return localStorage.getItem('accessToken');
    }
    return null;
  }

  // 7) REMOVE TOKEN
  removeToken(): void {
    if (isPlatformBrowser(this.platformId)) {
      localStorage.removeItem('accessToken');
    }
  }

  // 8) REMOVE USER
  removeUser(): void {
    if (isPlatformBrowser(this.platformId)) {
      localStorage.removeItem('currentUser');
    }
    // Also signal that there's now no user in memory.
    this.userSubject.next(null);
  }
}