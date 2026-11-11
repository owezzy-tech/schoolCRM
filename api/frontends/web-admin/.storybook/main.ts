import type { StorybookConfig } from '@storybook/angular';

const config: StorybookConfig = {
    stories: ['../src/stories/**/*.stories.ts'],
    addons: ['@storybook/addon-docs', '@storybook/addon-themes'],
    framework: '@storybook/angular',
    staticDirs: ['../public'],
    core: {
        disableTelemetry: true,
    },
};

export default config;
