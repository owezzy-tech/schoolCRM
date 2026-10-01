import { provideHttpClient } from '@angular/common/http';
import { inject, provideEnvironmentInitializer } from '@angular/core';
import { provideAnimations } from '@angular/platform-browser/animations';
import { withThemeByClassName } from '@storybook/addon-themes';
import type { Preview } from '@storybook/angular';
import {
    applicationConfig,
    componentWrapperDecorator,
} from '@storybook/angular';
import { IconsService } from 'app/core/icons/icons.service';

const preview: Preview = {
    tags: ['autodocs'],
    parameters: {
        layout: 'fullscreen',
        controls: {
            expanded: true,
        },
    },
    decorators: [
        applicationConfig({
            providers: [
                provideAnimations(),
                provideHttpClient(),
                provideEnvironmentInitializer(() => {
                    inject(IconsService);
                }),
            ],
        }),
        withThemeByClassName({
            themes: {
                Light: 'theme-default light',
                Dark: 'theme-default dark',
            },
            defaultTheme: 'Light',
        }),
        componentWrapperDecorator(
            (story) =>
                `<main class="bg-default text-default" style="min-height: 100vh; padding: 2rem;">${story}</main>`
        ),
    ],
};

export default preview;
