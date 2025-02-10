// app.config.ts
import { ApplicationConfig } from '@angular/core';
import { provideHttpClient, withFetch, withInterceptors } from '@angular/common/http';
import { provideRouter } from '@angular/router';
import { routes } from './app.routes';
import { provideClientHydration, withEventReplay } from '@angular/platform-browser';
import { provideAnimationsAsync } from '@angular/platform-browser/animations/async';

// Import our functional interceptor
import { tokenInterceptorFactory } from './interceptors/token.interceptor';

export const appConfig: ApplicationConfig = {
  providers: [
    provideRouter(routes),
    provideClientHydration(withEventReplay()),
    provideAnimationsAsync(),
    provideHttpClient(
      // If you use fetch-based HttpClient:
      withFetch(),
      // Attach our token interceptor
      withInterceptors([tokenInterceptorFactory])
    )
  ]
};