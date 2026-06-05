import {
    ApplicationConfig,
    importProvidersFrom,
    inject,
    isDevMode,
    provideAppInitializer,
} from '@angular/core';
import { LuxonDateAdapter } from '@angular/material-luxon-adapter';
import { DateAdapter, MAT_DATE_FORMATS } from '@angular/material/core';
import { provideAnimations } from '@angular/platform-browser/animations';
import { provideRouter, withInMemoryScrolling } from '@angular/router';
import { provideFuse } from '@fuse';
import { TranslocoService, provideTransloco } from '@jsverse/transloco';
import { provideEffects } from '@ngrx/effects';
import { provideRouterStore } from '@ngrx/router-store';
import { provideStore } from '@ngrx/store';
import { provideStoreDevtools } from '@ngrx/store-devtools';
import { appRoutes } from 'app/app.routes';
import { AuthEffects } from 'app/core/auth/+state/auth.effects';
import { provideAuth } from 'app/core/auth/auth.provider';
import { provideIcons } from 'app/core/icons/icons.provider';
import { MockApiService } from 'app/mock-api';
import {
    NgxUiLoaderConfig,
    NgxUiLoaderHttpModule,
    NgxUiLoaderModule,
    PB_DIRECTION,
} from 'ngx-ui-loader';
import { firstValueFrom } from 'rxjs';
import { metaReducers, rootReducers } from './core/store/store.logger';
import { TranslocoHttpLoader } from './core/transloco/transloco.http-loader';

const ngxUiLoaderConfig: NgxUiLoaderConfig = {
    hasProgressBar: true,
    pbColor: '#4f46e5',
    pbDirection: PB_DIRECTION.leftToRight,
    pbThickness: 3,
    fgsSize: 0,
    overlayColor: 'rgba(0,0,0,0)',
};

export const appConfig: ApplicationConfig = {
    providers: [
        provideAnimations(),
        importProvidersFrom(
            NgxUiLoaderModule.forRoot(ngxUiLoaderConfig),
            NgxUiLoaderHttpModule.forRoot({
                showForeground: true,
            })
        ),
        provideRouter(
            appRoutes,
            withInMemoryScrolling({ scrollPositionRestoration: 'enabled' })
        ),

        provideStore(rootReducers, {
            metaReducers,
            runtimeChecks: {
                // strictStateImmutability and strictActionImmutability are enabled by default
                strictStateSerializability: true,
                strictActionSerializability: true,
                strictActionWithinNgZone: false,
                strictActionTypeUniqueness: true,
            },
        }),
        provideRouterStore(),
        provideEffects(AuthEffects),
        provideStoreDevtools({
            name: 'School CRM Store',
            maxAge: 25,
            logOnly: !isDevMode(),
        }),
        // Material Date Adapter
        {
            provide: DateAdapter,
            useClass: LuxonDateAdapter,
        },
        {
            provide: MAT_DATE_FORMATS,
            useValue: {
                parse: {
                    dateInput: 'D',
                },
                display: {
                    dateInput: 'DDD',
                    monthYearLabel: 'LLL yyyy',
                    dateA11yLabel: 'DD',
                    monthYearA11yLabel: 'LLLL yyyy',
                },
            },
        },

        // Transloco Config
        provideTransloco({
            config: {
                availableLangs: [
                    {
                        id: 'en',
                        label: 'English',
                    },
                    {
                        id: 'tr',
                        label: 'Turkish',
                    },
                ],
                defaultLang: 'en',
                fallbackLang: 'en',
                reRenderOnLangChange: true,
                prodMode: !isDevMode(),
            },
            loader: TranslocoHttpLoader,
        }),
        provideAppInitializer(() => {
            const translocoService = inject(TranslocoService);
            const defaultLang = translocoService.getDefaultLang();
            translocoService.setActiveLang(defaultLang);

            return firstValueFrom(translocoService.load(defaultLang));
        }),

        // Fuse
        provideAuth(),
        provideIcons(),
        provideFuse({
            mockApi: {
                delay: 0,
                service: MockApiService,
            },
            fuse: {
                layout: 'classy',
                scheme: 'light',
                screens: {
                    sm: '600px',
                    md: '960px',
                    lg: '1280px',
                    xl: '1440px',
                },
                theme: 'theme-default',
                themes: [
                    {
                        id: 'theme-default',
                        name: 'Default',
                    },
                    {
                        id: 'theme-brand',
                        name: 'Brand',
                    },
                    {
                        id: 'theme-teal',
                        name: 'Teal',
                    },
                    {
                        id: 'theme-rose',
                        name: 'Rose',
                    },
                    {
                        id: 'theme-purple',
                        name: 'Purple',
                    },
                    {
                        id: 'theme-amber',
                        name: 'Amber',
                    },
                ],
            },
        }),
    ],
};
