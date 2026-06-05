import {
    provideHttpClient,
    withInterceptors,
    withInterceptorsFromDi,
} from '@angular/common/http';
import {
    EnvironmentProviders,
    Provider,
    inject,
    provideEnvironmentInitializer,
} from '@angular/core';
import { authInterceptor } from 'app/core/auth/auth.interceptor';
import { AuthService } from 'app/core/auth/auth.service';

export const provideAuth = (): Array<Provider | EnvironmentProviders> => {
    return [
        provideHttpClient(
            withInterceptors([authInterceptor]),
            withInterceptorsFromDi()
        ),
        provideEnvironmentInitializer(() => inject(AuthService)),
    ];
};
