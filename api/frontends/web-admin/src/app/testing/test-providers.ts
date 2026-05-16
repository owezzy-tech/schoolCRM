import { EnvironmentProviders, importProvidersFrom, Provider } from '@angular/core';
import { provideHttpClient } from '@angular/common/http';
import { provideRouter } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { FeatherModule } from 'angular-feather';
import { allIcons } from 'angular-feather/icons';
import { Chart, registerables } from 'chart.js';
import { DirectionService, RightSidebarService } from '@core';

Chart.register(...registerables);

export const testProviders: Array<Provider | EnvironmentProviders> = [
  provideHttpClient(),
  provideRouter([]),
  DirectionService,
  RightSidebarService,
  importProvidersFrom(TranslateModule.forRoot(), FeatherModule.pick(allIcons)),
];
