import { FuseCardComponent } from '@fuse/components/card';
import type { Meta, StoryObj } from '@storybook/angular';
import { moduleMetadata } from '@storybook/angular';

type CardArgs = Pick<FuseCardComponent, 'expanded'> & {
    title: string;
    curriculum: string;
    stage: string;
    description: string;
};

const meta: Meta<CardArgs> = {
    title: 'Fuse/Card',
    component: FuseCardComponent,
    decorators: [moduleMetadata({ imports: [FuseCardComponent] })],
    parameters: {
        docs: {
            description: {
                component:
                    'Layout examples using the existing Fuse card. These fixtures do not implement lesson-plan storage or publishing.',
            },
        },
    },
    args: {
        expanded: false,
        title: 'Lesson-plan draft',
        curriculum: 'Kenya CBC/CBE',
        stage: 'Grade 4',
        description:
            'Sample content for reviewing card spacing, typography and curriculum metadata.',
    },
    render: (args) => ({
        props: args,
        template: `
            <fuse-card [expanded]="expanded" class="w-full max-w-lg flex-col">
                <div class="p-6">
                    <p class="text-secondary text-sm">{{ curriculum }} · {{ stage }}</p>
                    <h2 class="mt-2 text-xl font-semibold">{{ title }}</h2>
                    <p class="mt-3">{{ description }}</p>
                </div>
                <div fuseCardExpansion class="border-t p-6">
                    <h3 class="font-semibold">Review sequence</h3>
                    <p class="mt-2">Teacher submission, HOD review, dean approval, then teacher publication.</p>
                </div>
            </fuse-card>
        `,
    }),
};

export default meta;
type Story = StoryObj<CardArgs>;

export const KenyaCurriculum: Story = {};

export const CambridgeCurriculum: Story = {
    args: {
        curriculum: 'International Cambridge',
        stage: 'Year 6',
        title: 'Inquiry-based lesson draft',
    },
};

export const Expanded: Story = {
    args: { expanded: true },
};
