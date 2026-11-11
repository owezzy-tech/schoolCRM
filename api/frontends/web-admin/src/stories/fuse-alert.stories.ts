import { FuseAlertComponent } from '@fuse/components/alert';
import type { Meta, StoryObj } from '@storybook/angular';
import { moduleMetadata } from '@storybook/angular';
import { expect, fn, userEvent, waitFor, within } from 'storybook/test';

type AlertArgs = Pick<
    FuseAlertComponent,
    'appearance' | 'type' | 'dismissible' | 'dismissed' | 'showIcon'
> & {
    title: string;
    message: string;
    dismissedChanged: ReturnType<typeof fn>;
};

const meta: Meta<AlertArgs> = {
    title: 'Fuse/Alert',
    component: FuseAlertComponent,
    decorators: [moduleMetadata({ imports: [FuseAlertComponent] })],
    argTypes: {
        appearance: {
            control: 'select',
            options: ['soft', 'border', 'fill', 'outline'],
        },
        type: {
            control: 'select',
            options: [
                'primary',
                'accent',
                'warn',
                'basic',
                'info',
                'success',
                'warning',
                'error',
            ],
        },
        dismissedChanged: { control: false },
    },
    args: {
        appearance: 'soft',
        type: 'info',
        dismissible: false,
        dismissed: false,
        showIcon: true,
        title: 'Awaiting review',
        message: 'The Head of Department will review this lesson-plan version.',
        dismissedChanged: fn(),
    },
    render: (args) => ({
        props: args,
        template: `
            <fuse-alert
                style="max-width: 42rem"
                [appearance]="appearance"
                [type]="type"
                [dismissible]="dismissible"
                [dismissed]="dismissed"
                [showIcon]="showIcon"
                (dismissedChanged)="dismissedChanged($event)"
            >
                <span fuseAlertTitle>{{ title }}</span>
                {{ message }}
            </fuse-alert>
        `,
    }),
};

export default meta;
type Story = StoryObj<AlertArgs>;

export const AwaitingReview: Story = {};

export const Published: Story = {
    args: {
        type: 'success',
        title: 'Lesson plan published',
        message: 'This version has completed HOD review and dean approval.',
    },
};

export const ChangesRequested: Story = {
    args: {
        type: 'warning',
        title: 'Changes requested',
        message:
            'Update the assessment activities before resubmitting for review.',
    },
};

export const PermissionDenied: Story = {
    args: {
        type: 'error',
        title: 'Publication unavailable',
        message:
            'Dean approval is required before this version can be published.',
    },
};

export const Dismissible: Story = {
    args: {
        dismissible: true,
        type: 'success',
        title: 'Draft saved',
        message: 'Your draft is saved and can be edited before submission.',
    },
    play: async ({ canvasElement, args }) => {
        const canvas = within(canvasElement);
        await userEvent.click(canvas.getByRole('button'));
        await expect(args.dismissedChanged).toHaveBeenCalledWith(true);
        await waitFor(async () => {
            await expect(
                canvas.queryByText(args.message)
            ).not.toBeInTheDocument();
        });
    },
};
