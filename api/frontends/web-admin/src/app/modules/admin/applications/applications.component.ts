import { ChangeDetectionStrategy, Component } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { RouterLink } from '@angular/router';
import {
    ApplicationFee,
    ApplicationFeeStatus,
    ApplicationKCSEResult,
    ApplicationType,
    KUCCPSPlacement,
} from 'app/core/admissions/admissions.types';

interface ApplicationRow {
    id: string;
    applicant: string;
    program: string;
    programmeCode: string;
    term: string;
    applicationType: ApplicationType;
    status: string;
    progress: number;
    reviewer: string;
    score: number | null;
    kuccpsPlacement?: KUCCPSPlacement;
    kcseResult?: ApplicationKCSEResult;
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
    readonly filters = [
        'Status',
        'Programme',
        'Application type',
        'KCSE mean grade',
        'Reviewer',
        'Fee status',
    ];

    readonly applications: readonly ApplicationRow[] = [
        {
            id: 'APP-3024',
            applicant: 'Achieng Otieno',
            program: 'B.Sc. Computer Science',
            programmeCode: 'BSC-CS',
            term: 'September 2026',
            applicationType: 'KUCCPS_PLACEMENT',
            status: 'In review',
            progress: 80,
            reviewer: 'Avery',
            score: 86,
            kuccpsPlacement: {
                placementID: 'KUCCPS-2026-003024',
                admissionNumber: 'SCM/2026/003024',
                institutionCode: 'SCM001',
                programmeCode: 'BSC-CS',
                programmeName: 'B.Sc. Computer Science',
                placementYear: 2026,
                clusterCode: 'CS01',
                clusterPoints: 42.318,
                weightedPointsNote:
                    'Public KUCCPS approximation; exact points require KNEC PI data.',
            },
            kcseResult: this.kcse('204003024', 2025, 'A-', 74),
            fee: this.fee('APP-3024', 'PAID', 'manual', 250000, {
                transactionID: 'mpesa-qe3024',
                paidAt: 'May 30, 2026',
            }),
        },
        {
            id: 'APP-3023',
            applicant: 'Brian Wekesa',
            program: 'Diploma in ICT',
            programmeCode: 'DIP-ICT',
            term: 'January 2027',
            applicationType: 'DIPLOMA',
            status: 'Documents pending',
            progress: 45,
            reviewer: 'Priya',
            score: null,
            kcseResult: this.kcse('204003023', 2024, 'C+', 46),
            fee: this.fee('APP-3023', 'PENDING', 'manual', 150000, {
                dueAt: 'Jun 10, 2026',
            }),
        },
        {
            id: 'APP-3022',
            applicant: 'Njeri Mwangi',
            program: 'Bachelor of Commerce',
            programmeCode: 'BCOM',
            term: 'September 2026',
            applicationType: 'SELF_SPONSORED_UNDERGRAD',
            status: 'Submitted',
            progress: 60,
            reviewer: 'Diego',
            score: 71,
            kcseResult: this.kcse('204003022', 2025, 'B', 59),
            fee: this.fee('APP-3022', 'FAILED', 'manual', 250000),
        },
        {
            id: 'APP-3021',
            applicant: 'Hassan Ali',
            program: 'Certificate in Data Support',
            programmeCode: 'CERT-DS',
            term: 'May 2027',
            applicationType: 'CERTIFICATE',
            status: 'Interview',
            progress: 90,
            reviewer: 'Avery',
            score: 92,
            kcseResult: this.kcse('204003021', 2023, 'C', 39),
            fee: this.fee('APP-3021', 'WAIVED', 'manual', 100000, {
                waiverGrantedBy: 'Admissions Finance',
                waiverReason:
                    'Need-based waiver approved for certificate intake.',
                waiverGrantedAt: 'May 29, 2026',
            }),
        },
        {
            id: 'APP-3020',
            applicant: 'Grace Kiplagat',
            program: 'M.A. Education',
            programmeCode: 'MA-EDU',
            term: 'September 2026',
            applicationType: 'MASTERS',
            status: 'Decision pending',
            progress: 95,
            reviewer: 'Diego',
            score: 78,
            fee: this.fee('APP-3020', 'REFUNDED', 'manual', 300000, {
                transactionID: 'mpesa-rf3020',
            }),
        },
        {
            id: 'APP-3019',
            applicant: 'Peter Mutua',
            program: 'TVET Mechatronics',
            programmeCode: 'TVET-MECH',
            term: 'January 2027',
            applicationType: 'TVET',
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

    formatApplicationType(type: ApplicationType): string {
        return type.replaceAll('_', ' ');
    }

    formatAmount(fee: ApplicationFee): string {
        if (fee.status === 'NOT_REQUIRED') {
            return 'No fee';
        }

        return new Intl.NumberFormat('en-KE', {
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
            currency: 'KES',
            status,
            provider,
            auditTrail: [],
            ...overrides,
        };
    }

    private kcse(
        indexNumber: string,
        examYear: number,
        meanGrade: string,
        meanPoints: number
    ): ApplicationKCSEResult {
        return {
            indexNumber,
            examYear,
            meanGrade,
            meanPoints,
            subjects: [
                { subjectCode: '101', grade: meanGrade, points: 11 },
                { subjectCode: '102', grade: 'B+', points: 10 },
                { subjectCode: '121', grade: 'A-', points: 11 },
            ],
        };
    }
}
