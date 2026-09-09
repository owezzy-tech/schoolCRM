import {
    HttpErrorResponse,
    HttpEvent,
    HttpHandlerFn,
    HttpRequest,
} from '@angular/common/http';
import { inject } from '@angular/core';
import { Router } from '@angular/router';
import { AuthService } from 'app/core/auth/auth.service';
import { authorizationTokenForRequest } from 'app/core/auth/auth-token-selector';
import { AuthUtils } from 'app/core/auth/auth.utils';
import { PortalAuthService } from 'app/core/portal/portal-auth.service';
import { Observable, catchError, throwError } from 'rxjs';

/**
 * Intercept
 *
 * @param req
 * @param next
 */
export const authInterceptor = (
    req: HttpRequest<unknown>,
    next: HttpHandlerFn
): Observable<HttpEvent<unknown>> => {
    const authService = inject(AuthService);
    const portalAuthService = inject(PortalAuthService);
    const router = inject(Router);

    // Clone the request object
    let newReq = req.clone();

    // Request
    //
    // If the access token didn't expire, add the Authorization header.
    // We won't add the Authorization header if the access token expired.
    // This will force the server to return a "401 Unauthorized" response
    // for the protected API routes which our response interceptor will
    // catch and delete the access token from the local storage while logging
    // the user out from the app.
    const accessToken = authorizationTokenForRequest({
        requestUrl: req.url,
        currentRouteUrl: router.url,
        staffToken: authService.accessToken,
        portalToken: portalAuthService.accessToken(),
    });

    if (accessToken) {
        newReq = req.clone({
            headers: req.headers.set('Authorization', 'Bearer ' + accessToken),
        });
    }

    // Response
    return next(newReq).pipe(
        catchError((error) => {
            // Catch "401 Unauthorized" responses
            if (
                error instanceof HttpErrorResponse &&
                error.status === 401 &&
                shouldEndSession(newReq, authService)
            ) {
                // Sign out
                authService.signOut();

                // Reload the app
                location.reload();
            }

            return throwError(() => error);
        })
    );
};

function shouldEndSession(
    req: HttpRequest<unknown>,
    authService: AuthService
): boolean {
    if (!authService.accessToken || AuthUtils.isTokenExpired(authService.accessToken)) {
        return true;
    }

    return req.url.startsWith('/v1/auth/');
}
