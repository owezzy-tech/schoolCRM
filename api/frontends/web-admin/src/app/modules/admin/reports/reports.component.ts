import { ChangeDetectionStrategy, Component } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';

interface ReportStat {
    label: string;
    value: string;
}

interface ReportPanel {
    id: string;
    title: string;
    stats: ReportStat[];
}

@Component({
    selector: 'app-reports',
    standalone: true,
    imports: [MatButtonModule, MatIconModule],
    changeDetection: ChangeDetectionStrategy.OnPush,
    templateUrl: './reports.component.html',
})
export class ReportsComponent {
    readonly panels: readonly ReportPanel[] = [
        {
            id: 'funnel',
            title: 'Funnel overview',
            stats: [
                { label: 'Total Inquiries', value: '4,250' },
                { label: 'Started Apps', value: '1,840' },
                { label: 'Conversion', value: '43.2%' },
            ],
        },
        {
            id: 'source',
            title: 'Source attribution',
            stats: [
                { label: 'Organic Search', value: '38%' },
                { label: 'Direct Traffic', value: '24%' },
                { label: 'Paid Social', value: '18%' },
            ],
        },
        {
            id: 'reviewer',
            title: 'Reviewer load',
            stats: [
                { label: 'Avg Time/App', value: '14m' },
                { label: 'Pending Queue', value: '342' },
                { label: 'Completed (MTD)', value: '1,205' },
            ],
        },
        {
            id: 'yield',
            title: 'Yield by program',
            stats: [
                { label: 'Computer Science', value: '68%' },
                { label: 'Business', value: '54%' },
                { label: 'Engineering', value: '61%' },
            ],
        },
    ];
}
