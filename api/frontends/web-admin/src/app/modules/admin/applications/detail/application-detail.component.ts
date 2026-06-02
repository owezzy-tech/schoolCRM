import { ChangeDetectionStrategy, Component, inject } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatTabsModule } from '@angular/material/tabs';
import { ActivatedRoute, RouterLink } from '@angular/router';
import {
    ApplicationFee,
    ApplicationFeeStatus,
    ApplicationKCSEResult,
    KUCCPSPlacement,
    NotificationChannel,
    NotificationPreferences,
} from 'app/core/admissions/admissions.types';
import { ApplicationDocumentsComponent } from '../components/documents/application-documents.component';

interface TimelineEvent {
    icon: string;
    title: string;
    date: string;
}

interface Note {
    authorInitials: string;
    authorName: string;
    timestamp: string;
    body: string;
}

interface Insight {
    title: string;
    description: string;
}

interface ContactPreferenceRow {
    channel: NotificationChannel;
    label: string;
    destination: string;
    optIn: boolean;
    provider: string;
    consentNote: string;
}

@Component({
    selector: 'app-application-detail',
    standalone: true,
    imports: [
        RouterLink,
        MatButtonModule,
        MatIconModule,
        MatTabsModule,
        ApplicationDocumentsComponent,
    ],
    changeDetection: ChangeDetectionStrategy.OnPush,
    templateUrl: './application-detail.component.html',
})
export class ApplicationDetailComponent {
    private readonly route = inject(ActivatedRoute);
    readonly applicationId = this.route.snapshot.paramMap.get('id');

    readonly applicantName = 'Achieng Otieno';
    readonly status = 'In review';
    readonly applicationType = 'KUCCPS PLACEMENT';
    readonly kuccpsPlacement: KUCCPSPlacement = {
        placementID: 'KUCCPS-2026-003024',
        admissionNumber: 'SCM/2026/003024',
        institutionCode: 'SCM001',
        programmeCode: 'BSC-CS',
        programmeName: 'B.Sc. Computer Science',
        placementYear: 2026,
        clusterCode: 'CS01',
        clusterPoints: 42.318,
        weightedPointsNote:
            'Public KUCCPS approximation; exact points require private KNEC performance indices.',
    };
    readonly kcseResult: ApplicationKCSEResult = {
        indexNumber: '204003024',
        examYear: 2025,
        meanGrade: 'A-',
        meanPoints: 74,
        subjects: [
            { subjectCode: '101', grade: 'A-', points: 11 },
            { subjectCode: '102', grade: 'B+', points: 10 },
            { subjectCode: '121', grade: 'A', points: 12 },
            { subjectCode: '232', grade: 'A-', points: 11 },
        ],
    };
    readonly fee: ApplicationFee = {
        id: 'fee-app-3024',
        applicationID: 'APP-3024',
        amountCents: 250000,
        currency: 'KES',
        status: 'PAID',
        provider: 'manual',
        transactionID: 'mpesa-qe3024',
        paidAt: 'May 30, 2026 10:12 AM',
        auditTrail: [
            {
                actor: 'System',
                action: 'Payment recorded',
                reason: 'Stripe provider webhook marked application fee paid.',
                timestamp: 'May 30, 2026 10:12 AM',
            },
            {
                actor: 'Avery',
                action: 'Receipt reviewed',
                reason: 'Confirmed payment before moving application forward.',
                timestamp: 'May 30, 2026 10:45 AM',
            },
        ],
    };

    readonly summaryCards = [
        { label: 'Programme', value: 'B.Sc. Computer Science' },
        { label: 'Application type', value: this.applicationType },
        { label: 'Term', value: 'September 2026' },
        { label: 'Submitted', value: '2 hours ago' },
        { label: 'Application Fee', value: 'Paid · KES 2,500.00' },
    ];

    readonly applicantInfo = [
        { label: 'Name', value: 'Achieng Otieno' },
        { label: 'Email', value: 'achieng.otieno@example.ac.ke' },
        { label: 'Phone', value: '+254 712 345 678' },
        { label: 'DOB', value: 'May 14, 2007' },
        { label: 'National ID', value: '42817391' },
        { label: 'UPI', value: 'KEMIS7Q4A2' },
        { label: 'Address', value: 'P.O. Box 1024-40100, Kisumu' },
        {
            label: 'Essay excerpt',
            value: '"My interest in software engineering began through a county robotics club..."',
        },
    ];

    readonly notificationPreferences: NotificationPreferences = {
        smsOptIn: true,
        whatsAppOptIn: false,
        emailOptIn: true,
        priority: ['SMS', 'WHATSAPP', 'EMAIL'],
    };

    readonly contactPreferenceRows: ContactPreferenceRow[] = [
        {
            channel: 'SMS',
            label: 'SMS',
            destination: '+254 712 345 678',
            optIn: this.notificationPreferences.smsOptIn,
            provider: 'Celcom Africa SMS',
            consentNote: 'Default urgent-alert channel for Kenya applicants.',
        },
        {
            channel: 'WHATSAPP',
            label: 'WhatsApp',
            destination: '+254 712 345 678',
            optIn: this.notificationPreferences.whatsAppOptIn,
            provider: 'WhatsApp Cloud',
            consentNote: 'Off until the applicant explicitly opts in.',
        },
        {
            channel: 'EMAIL',
            label: 'Email',
            destination: 'achieng.otieno@example.ac.ke',
            optIn: this.notificationPreferences.emailOptIn,
            provider: 'SchoolCRM mailer',
            consentNote: 'Kept on for long-form admissions notices.',
        },
    ];

    readonly notes: Note[] = [
        {
            authorInitials: 'AM',
            authorName: 'Avery',
            timestamp: '1 hour ago',
            body: 'Strong KCSE profile and KUCCPS placement matches the selected programme catalog record.',
        },
        {
            authorInitials: 'PP',
            authorName: 'Priya',
            timestamp: 'Yesterday',
            body: 'KCSE result slip verified; waiting for finance receipt reconciliation.',
        },
        {
            authorInitials: 'SM',
            authorName: 'System',
            timestamp: '2 days ago',
            body: 'Application submitted successfully.',
        },
    ];

    readonly timelineEvents: TimelineEvent[] = [
        {
            icon: 'heroicons_outline:check-circle',
            title: 'Review started',
            date: 'Today, 9:00 AM',
        },
        {
            icon: 'heroicons_outline:document-text',
            title: 'All documents received',
            date: 'Yesterday, 2:30 PM',
        },
        {
            icon: 'heroicons_outline:arrow-up-tray',
            title: 'Transcript uploaded',
            date: 'Yesterday, 10:15 AM',
        },
        {
            icon: 'heroicons_outline:envelope',
            title: 'Recommendation 2 received',
            date: 'May 28, 4:00 PM',
        },
        {
            icon: 'heroicons_outline:envelope',
            title: 'Recommendation 1 received',
            date: 'May 27, 11:20 AM',
        },
        {
            icon: 'heroicons_outline:paper-airplane',
            title: 'Application submitted',
            date: 'May 26, 8:45 PM',
        },
    ];

    readonly insights: Insight[] = [
        {
            title: 'Recommended decision',
            description:
                'Strong admit based on KCSE mean grade and cluster estimate.',
        },
        { title: 'Risk flags', description: 'None identified.' },
        {
            title: 'Similar profiles',
            description:
                'Matches prior KUCCPS placement admits for B.Sc. Computer Science.',
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
        return new Intl.NumberFormat('en-KE', {
            style: 'currency',
            currency: fee.currency,
        }).format(fee.amountCents / 100);
    }

    preferenceIcon(channel: NotificationChannel): string {
        const icons: Record<NotificationChannel, string> = {
            SMS: 'heroicons_outline:chat-bubble-left',
            WHATSAPP: 'heroicons_outline:device-phone-mobile',
            EMAIL: 'heroicons_outline:envelope',
        };

        return icons[channel];
    }
}
