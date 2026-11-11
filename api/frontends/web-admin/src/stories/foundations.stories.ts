import type { Meta, StoryObj } from '@storybook/angular';

/**
 * Visual reference for the Elimu tokens in `src/styles/design-tokens.css`.
 * Every swatch and sample reads the live CSS custom property, so these
 * stories change with the light/dark toolbar and never drift from the tokens.
 */
const meta: Meta = {
    title: 'Design System/Foundations',
    parameters: {
        docs: {
            description: {
                component:
                    'Tokens defined once in `src/styles/design-tokens.css` (light on `:root`, dark on `.dark`). See `design-spec.md` for the rules behind them.',
            },
        },
    },
};

export default meta;
type Story = StoryObj;

const page = `font-family: var(--elimu-font-sans); color: var(--elimu-fg);`;
const eyebrow = `font-family: var(--elimu-font-mono); font-size: var(--elimu-text-eyebrow);
    letter-spacing: var(--elimu-tracking-eyebrow); text-transform: uppercase; color: var(--elimu-muted);`;
const card = `background: var(--elimu-surface); border: 1px solid var(--elimu-border);
    border-radius: var(--elimu-radius); box-shadow: var(--elimu-shadow); padding: var(--elimu-space-5);`;

const colours = [
    ['bg', 'Page ground'],
    ['surface', 'Cards, topbar, sheets'],
    ['surface-2', 'Wells, tab rails'],
    ['fg', 'Primary text'],
    ['fg-soft', 'Derived: softened text'],
    ['muted', 'Secondary text, captions'],
    ['border', 'Hairlines, dividers'],
    ['accent', 'Brand: active marker, single primary action'],
    ['accent-soft', 'Derived: soft button, selection'],
    ['success', 'Status only'],
    ['warning', 'Status only'],
    ['danger', 'Status only'],
];

export const Colour: Story = {
    render: () => ({
        template: `
        <div style="${page} display: grid; gap: var(--elimu-space-3);
                    grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));">
            ${colours
                .map(
                    ([name, role]) => `
                <div style="${card} padding: 0; overflow: hidden;">
                    <div style="height: 72px; background: var(--elimu-${name});
                                border-bottom: 1px solid var(--elimu-border);"></div>
                    <div style="padding: var(--elimu-space-3) var(--elimu-space-4);">
                        <div style="font-family: var(--elimu-font-mono); font-size: 13px;">--elimu-${name}</div>
                        <div style="color: var(--elimu-muted); font-size: 13px;">${role}</div>
                    </div>
                </div>`
                )
                .join('')}
        </div>`,
    }),
};

export const Typography: Story = {
    render: () => ({
        template: `
        <div style="${page} display: grid; gap: var(--elimu-space-6); max-width: 72ch;">
            <div>
                <div style="${eyebrow}">Display · page title</div>
                <h1 style="font-family: var(--elimu-font-display); font-feature-settings: var(--elimu-font-display-features); font-size: var(--elimu-text-title);
                           font-weight: 600; letter-spacing: -0.01em; margin: var(--elimu-space-2) 0 0;">
                    Grade 4 Science lesson plans
                </h1>
            </div>
            <div>
                <div style="${eyebrow}">Body · 15px / 1.6</div>
                <p style="font-size: var(--elimu-text-body); line-height: 1.6; margin: var(--elimu-space-2) 0 0;">
                    Teachers draft lessons against the curriculum, submit them for review, and publish
                    approved versions to their class. Line length stays between 65 and 75 characters.
                </p>
            </div>
            <div>
                <div style="${eyebrow}">Mono · tabular numerals</div>
                <div style="font-family: var(--elimu-font-mono); font-variant-numeric: tabular-nums;
                            font-size: var(--elimu-text-data); margin-top: var(--elimu-space-2);">
                    1,204 learners · 87.5% attendance · KES 45,000.00
                </div>
            </div>
        </div>`,
    }),
};

const spaces = [1, 2, 3, 4, 5, 6, 8];

export const SpacingAndShape: Story = {
    name: 'Spacing and shape',
    render: () => ({
        template: `
        <div style="${page} display: grid; gap: var(--elimu-space-6);">
            <div style="display: grid; gap: var(--elimu-space-2);">
                ${spaces
                    .map(
                        (n) => `
                    <div style="display: flex; align-items: center; gap: var(--elimu-space-4);">
                        <code style="width: 140px; font-size: 13px;">--elimu-space-${n}</code>
                        <div style="height: 12px; width: var(--elimu-space-${n});
                                    background: var(--elimu-accent); border-radius: 2px;"></div>
                    </div>`
                    )
                    .join('')}
            </div>
            <div style="display: flex; gap: var(--elimu-space-4); flex-wrap: wrap;">
                <div style="${card} width: 200px;"><code>--elimu-radius</code></div>
                <div style="${card} width: 200px; border-radius: var(--elimu-radius-lg);">
                    <code>--elimu-radius-lg</code>
                </div>
            </div>
        </div>`,
    }),
};

export const Motion: Story = {
    render: () => ({
        styles: [
            `.press {
                font: 500 14px var(--elimu-font-sans);
                padding: var(--elimu-space-2) var(--elimu-space-4);
                border-radius: var(--elimu-radius);
                border: 0;
                background: var(--elimu-accent);
                color: var(--elimu-on-accent);
                cursor: pointer;
                transition: transform var(--elimu-duration-press) var(--elimu-ease-out);
            }
            .press:active { transform: scale(0.97); }
            .press:focus-visible { outline: var(--elimu-focus-ring); outline-offset: var(--elimu-focus-offset); }`,
        ],
        template: `
        <div style="${page} display: grid; gap: var(--elimu-space-4); max-width: 60ch;">
            <button class="press" type="button" style="justify-self: start;">Press me</button>
            <p style="color: var(--elimu-muted); margin: 0; line-height: 1.6;">
                Press feedback: <code>scale(0.97)</code> over <code>--elimu-duration-press</code>
                (160ms) with <code>--elimu-ease-out</code>. Popovers use 180ms and view changes 240ms.
                Only <code>transform</code> and <code>opacity</code> animate. With reduced motion
                every duration token becomes 0ms.
            </p>
        </div>`,
    }),
};
