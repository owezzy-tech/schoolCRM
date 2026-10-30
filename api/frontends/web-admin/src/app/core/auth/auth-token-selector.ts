import { AuthUtils } from './auth.utils';

export interface AuthorizationTokenContext {
    requestUrl: string;
    currentRouteUrl: string;
    staffToken: string;
    portalToken: string | null;
}

export function authorizationTokenForRequest({
    requestUrl,
    currentRouteUrl,
    staffToken,
    portalToken,
}: AuthorizationTokenContext): string | null {
    if (requestUrl.startsWith('/v1/auth/applicant-portal/')) {
        return null;
    }

    if (
        currentRouteUrl.startsWith('/portal') &&
        requestUrl.startsWith('/v1/admissions/') &&
        portalToken
    ) {
        return portalToken;
    }

    if (staffToken && !AuthUtils.isTokenExpired(staffToken)) {
        return staffToken;
    }

    return null;
}
