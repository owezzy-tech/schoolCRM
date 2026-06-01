import { ChangeDetectionStrategy, Component } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { RouterLink } from '@angular/router';
import {
    ApplicationFee,
    ApplicationFeeStatus,
} from 'app/core/admissions/admissions.types';

interface Application {
    id: string;
    applicant: string;
    program: string;
    term: string;
    status: string;
    progress: number;
    reviewer: string;
    score: number | null;
    fee: ApplicationFee;
}

@Component({
    selector: 'app-applications',
    standalone: true,
    imports: [RouterLink, MatButtonModule, MatIconModule],
    changeDetection: ChangeDetectionStrategy.OnPush,
    templateUrl: './applications.component.html',
})
export class ApplicationsComponent {
    readonly filters = ['Status', 'Program', 'Term', 'Reviewer', 'Fee status'];

    readonly applications: readonly Application[] = [
        {
            id: 'APP-3024',
            applicant: 'Sofia Martinez',
            program: 'B.Sc. CS',
            term: 'Fall 2026',
            status: 'In review',
            progress: 80,
            reviewer: 'Avery',
            score: 86,
            fee: this.fee('APP-3024', 'PAID', 'stripe', 15000, {
                transactionID: 'pi_3024_fall26',
                paidAt: 'May 30, 2026',
            }),
        },
        {
            id: 'APP-3023',
            applicant: 'James Okoro',
            program: 'M.A. Edu',
            term: 'Spring 2027',
            status: 'Documents pending',
            progress: 45,
            reviewer: 'Priya',
            score: null,
            fee: this.fee('APP-3023', 'PENDING', 'stripe', 15000, {
                dueAt: 'Jun 10, 2026',
            }),
        },
        {
            id: 'APP-3022',
            applicant: 'Priya Patel',
            program: 'B.A. Econ',
            term: 'Fall 2026',
            status: 'Submitted',
            progress: 60,
            reviewer: 'Diego',
            score: 71,
            fee: this.fee('APP-3022', 'FAILED', 'stripe', 15000),
        },
        {
            id: 'APP-3021',
            applicant: 'Liam Chen',
            program: 'B.Sc. Bio',
            term: 'Fall 2026',
            status: 'Interview',
            progress: 90,
            reviewer: 'Avery',
            score: 92,
            fee: this.fee('APP-3021', 'WAIVED', 'manual', 15000, {
                waiverGrantedBy: 'Maya Schultz',
                waiverReason: 'Need-based waiver approved by admissions admin.',
                waiverGrantedAt: 'May 29, 2026',
            }),
        },
        {
            id: 'APP-3020',
            applicant: 'Aisha Bello',
            program: 'M.B.A.',
            term: 'Fall 2026',
            status: 'Decision pending',
            progress: 95,
            reviewer: 'Diego',
            score: 78,
            fee: this.fee('APP-3020', 'REFUNDED', 'stripe', 15000, {
                transactionID: 're_3020_refund',
            }),
        },
        {
            id: 'APP-3019',
            applicant: 'Noah Williams',
            program: 'B.Sc. CS',
            term: 'Fall 2026',
            status: 'Admitted',
            progress: 100,
            reviewer: 'Avery',
            score: 95,
            fee: this.fee('APP-3019', 'NOT_REQUIRED', 'not_required', 0),
        },
    ];

    feeStatusClass(status: ApplicationFeeStatus): string {
        switch (status) {
            case 'PAID':
                return 'bg-green-100 text-green-700';
            case 'PENDING':
                return 'bg-amber-100 text-amber-700';
            case 'FAILED':
                return 'bg-red-100 text-red-700';
            case 'WAIVED':
                return 'bg-purple-100 text-purple-700';
            case 'REFUNDED':
                return 'bg-blue-100 text-blue-700';
            case 'NOT_REQUIRED':
                return 'bg-slate-100 text-secondary';
        }
    }

    formatFeeStatus(status: ApplicationFeeStatus): string {
        return status.replaceAll('_', ' ');
    }

    formatAmount(fee: ApplicationFee): string {
        if (fee.status === 'NOT_REQUIRED') {
            return 'No fee';
        }

        return new Intl.NumberFormat('en-US', {
            style: 'currency',
            currency: fee.currency,
        }).format(fee.amountCents / 100);
    }

    private fee(
        applicationID: string,
        status: ApplicationFeeStatus,
        provider: ApplicationFee['provider'],
        amountCents: number,
        overrides: Partial<ApplicationFee> = {}
    ): ApplicationFee {
        return {
            id: `fee-${applicationID.toLowerCase()}`,
            applicationID,
            amountCents,
            currency: 'USD',
            status,
            provider,
            auditTrail: [],
            ...overrides,
        };
    }
}
