import { ChangeDetectionStrategy, Component } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { MatTabsModule } from '@angular/material/tabs';
import { RouterLink } from '@angular/router';
import {
    ApplicationStatus,
    DocumentStatus,
} from 'app/core/admissions/admissions.types';

type ReviewDecision = Extract<
    ApplicationStatus,
    'ADMITTED' | 'DENIED' | 'WAITLISTED' | 'DEFERRED'
>;

interface AssignedApplication {
    id: string;
    applicant: string;
    program: string;
    term: string;
    status: Extract<
        ApplicationStatus,
        'READY_FOR_REVIEW' | 'IN_REVIEW' | 'DECISION_PENDING'
    >;
    statusLabel: string;
    reviewer: string;
    authorized: boolean;
    score: number;
    checklistComplete: number;
    checklistTotal: number;
    submittedAt: string;
    dueAt: string;
    risk: 'Low' | 'Medium' | 'High';
}

interface ReviewDocument {
    name: string;
    status: DocumentStatus;
    required: boolean;
    reviewerNotes: string;
}

interface ReviewTransition {
    fromStatus: ApplicationStatus;
    toStatus: ApplicationStatus;
    actor: string;
    reason: string;
    timestamp: string;
}

interface ReviewAuditEvent {
    timestamp: string;
    actor: string;
    action: string;
    entity: string;
}

interface DecisionForm {
    decision: ReviewDecision;
    reason: string;
    committeeNote: string;
}

const ASSIGNED_APPLICATIONS: AssignedApplication[] = [
    {
        id: 'APP-3024',
        applicant: 'Sofia Martinez',
        program: 'B.Sc. Computer Science',
        term: 'Fall 2026',
        status: 'IN_REVIEW',
        statusLabel: 'In review',
        reviewer: 'Avery',
        authorized: true,
        score: 86,
        checklistComplete: 5,
        checklistTotal: 5,
        submittedAt: 'May 30, 8:45 AM',
        dueAt: 'Jun 3, 5:00 PM',
        risk: 'Low',
    },
    {
        id: 'APP-3023',
        applicant: 'James Okoro',
        program: 'M.A. Education',
        term: 'Spring 2027',
        status: 'READY_FOR_REVIEW',
        statusLabel: 'Ready for review',
        reviewer: 'Avery',
        authorized: true,
        score: 74,
        checklistComplete: 4,
        checklistTotal: 5,
        submittedAt: 'May 29, 2:20 PM',
        dueAt: 'Jun 5, 12:00 PM',
        risk: 'Medium',
    },
    {
        id: 'APP-3020',
        applicant: 'Aisha Bello',
        program: 'M.B.A.',
        term: 'Fall 2026',
        status: 'DECISION_PENDING',
        statusLabel: 'Decision pending',
        reviewer: 'Diego',
        authorized: false,
        score: 78,
        checklistComplete: 5,
        checklistTotal: 5,
        submittedAt: 'May 27, 10:05 AM',
        dueAt: 'Jun 2, 4:00 PM',
        risk: 'High',
    },
];

const REVIEW_DOCUMENTS: ReviewDocument[] = [
    {
        name: 'Official transcript',
        status: 'ACCEPTED',
        required: true,
        reviewerNotes: 'GPA verified against submitted coursework.',
    },
    {
        name: 'Personal statement',
        status: 'ACCEPTED',
        required: true,
        reviewerNotes: 'Strong alignment with program outcomes.',
    },
    {
        name: 'Recommendation letter 1',
        status: 'ACCEPTED',
        required: true,
        reviewerNotes: 'Teacher recommendation complete.',
    },
    {
        name: 'Recommendation letter 2',
        status: 'PENDING_REVIEW',
        required: true,
        reviewerNotes: 'Needs reviewer confirmation before committee vote.',
    },
];

const REVIEW_TRANSITIONS: ReviewTransition[] = [
    {
        fromStatus: 'READY_FOR_REVIEW',
        toStatus: 'IN_REVIEW',
        actor: 'Avery',
        reason: 'Reviewer opened assigned queue item.',
        timestamp: 'May 30, 9:00 AM',
    },
    {
        fromStatus: 'AWAITING_DOCUMENTS',
        toStatus: 'READY_FOR_REVIEW',
        actor: 'System',
        reason: 'Required checklist reached review threshold.',
        timestamp: 'May 29, 4:10 PM',
    },
];

const REVIEW_AUDIT_EVENTS: ReviewAuditEvent[] = [
    {
        timestamp: 'May 30, 10:45 AM',
        actor: 'Avery',
        action: 'Update',
        entity: 'Application (APP-3024)',
    },
    {
        timestamp: 'May 30, 10:30 AM',
        actor: 'System',
        action: 'Create',
        entity: 'ApplicationTransition (APP-3024)',
    },
    {
        timestamp: 'May 29, 4:10 PM',
        actor: 'System',
        action: 'Update',
        entity: 'Checklist (APP-3024)',
    },
];

@Component({
    selector: 'app-staff-review-workspace',
    imports: [
        FormsModule,
        RouterLink,
        MatButtonModule,
        MatFormFieldModule,
        MatIconModule,
        MatInputModule,
        MatSelectModule,
        MatTabsModule,
    ],
    templateUrl: './staff-review-workspace.component.html',
    changeDetection: ChangeDetectionStrategy.OnPush,
})
export class StaffReviewWorkspaceComponent {
    readonly currentReviewer = 'Avery';
    readonly applications = ASSIGNED_APPLICATIONS;
    readonly documents = REVIEW_DOCUMENTS;
    readonly transitions = REVIEW_TRANSITIONS;
    readonly auditEvents = REVIEW_AUDIT_EVENTS;
    readonly decisions: ReviewDecision[] = [
        'ADMITTED',
        'DENIED',
        'WAITLISTED',
        'DEFERRED',
    ];

    selectedApplication = ASSIGNED_APPLICATIONS[0];
    decisionForm: DecisionForm = {
        decision: 'ADMITTED',
        reason: 'Strong academic profile and complete checklist.',
        committeeNote:
            'Recommend admit with standard enrollment checklist follow-up.',
    };

    selectApplication(application: AssignedApplication): void {
        if (!application.authorized) {
            return;
        }

        this.selectedApplication = application;
    }

    authorizedApplications(): readonly AssignedApplication[] {
        return this.applications.filter(
            (application) => application.authorized
        );
    }

    checklistPercent(application: AssignedApplication): number {
        if (application.checklistTotal === 0) {
            return 0;
        }

        return Math.round(
            (application.checklistComplete / application.checklistTotal) * 100
        );
    }

    statusClass(status: ApplicationStatus): string {
        switch (status) {
            case 'READY_FOR_REVIEW':
                return 'bg-blue-100 text-blue-700';
            case 'IN_REVIEW':
                return 'bg-amber-100 text-amber-700';
            case 'DECISION_PENDING':
                return 'bg-purple-100 text-purple-700';
            case 'ADMITTED':
            case 'ENROLLED':
                return 'bg-green-100 text-green-700';
            case 'DENIED':
            case 'WITHDRAWN':
                return 'bg-red-100 text-red-700';
            default:
                return 'bg-slate-100 text-secondary';
        }
    }

    documentStatusClass(status: DocumentStatus): string {
        switch (status) {
            case 'ACCEPTED':
                return 'bg-green-100 text-green-700';
            case 'REJECTED':
                return 'bg-red-100 text-red-700';
            case 'WAIVED':
                return 'bg-purple-100 text-purple-700';
            case 'PENDING_REVIEW':
            case 'UPLOADED':
                return 'bg-amber-100 text-amber-700';
            default:
                return 'bg-slate-100 text-secondary';
        }
    }

    riskClass(risk: AssignedApplication['risk']): string {
        switch (risk) {
            case 'Low':
                return 'text-green-700 bg-green-100';
            case 'Medium':
                return 'text-amber-700 bg-amber-100';
            case 'High':
                return 'text-red-700 bg-red-100';
        }
    }

    formatStatus(status: string): string {
        return status.replaceAll('_', ' ');
    }
}
