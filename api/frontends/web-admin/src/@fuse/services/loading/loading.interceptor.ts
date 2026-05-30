import {
    HttpContextToken,
    HttpEvent,
    HttpHandlerFn,
    HttpRequest,
} from '@angular/common/http';
import { inject } from '@angular/core';
import { FuseLoadingService } from '@fuse/services/loading/loading.service';
import { Observable, finalize, take } from 'rxjs';

export const FUSE_SKIP_LOADING = new HttpContextToken<boolean>(() => false);

export const fuseLoadingInterceptor = (
    req: HttpRequest<unknown>,
    next: HttpHandlerFn
): Observable<HttpEvent<unknown>> => {
    const fuseLoadingService = inject(FuseLoadingService);
    let handleRequestsAutomatically = false;

    if (req.context.get(FUSE_SKIP_LOADING)) {
        return next(req);
    }

    fuseLoadingService.auto$.pipe(take(1)).subscribe((value) => {
        handleRequestsAutomatically = value;
    });

    // If the Auto mode is turned off, do nothing
    if (!handleRequestsAutomatically) {
        return next(req);
    }

    // Set the loading status to true
    fuseLoadingService._setLoadingStatusFromInterceptor(true, req.url);

    return next(req).pipe(
        finalize(() => {
            // Set the status to false if there are any errors or the request is completed
            fuseLoadingService._setLoadingStatusFromInterceptor(false, req.url);
        })
    );
};
