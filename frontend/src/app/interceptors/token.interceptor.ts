// token.interceptor.ts
import { HttpInterceptorFn } from '@angular/common/http';
import { inject, PLATFORM_ID } from '@angular/core';
import { isPlatformBrowser } from '@angular/common';
import { AuthService } from '../services/auth.service';

export const tokenInterceptorFactory: HttpInterceptorFn = (request, next) => {
  const platformId = inject(PLATFORM_ID);
  const authService = inject(AuthService);

  // Only attach the token in the browser (avoid SSR conflicts)
  if (isPlatformBrowser(platformId)) {
    const token = authService.getToken();
    if (token) {
      // Clone the request to set the Authorization header
      const cloned = request.clone({
        setHeaders: {
          Authorization: `Bearer ${token}`
        }
      });
      return next(cloned);
    }
  }

  // If no token, or we're on the server, just forward the request as-is
  return next(request);
};