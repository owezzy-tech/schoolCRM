export const LEAD_SCORE_BANDS = [
    'COLD',
    'WARM',
    'HOT',
    'READY_TO_APPLY',
] as const;

export const LEAD_SCORE_CRITERION_FIELDS = [
    'lifecycle_stage',
    'application_type',
    'application_status',
    'program_id',
    'academic_term_id',
] as const;

export const LEAD_SCORE_CRITERION_OPERATORS = ['EQ', 'IN'] as const;

export const LIFECYCLE_STAGES = [
    'PROSPECT',
    'INQUIRY',
    'APPLICANT',
    'ADMITTED',
    'ENROLLED',
    'ALUMNI',
] as const;

export const APPLICATION_TYPES = ['FRESHMAN', 'TRANSFER', 'GRADUATE'] as const;

export const APPLICATION_STATUSES = [
    'DRAFT',
    'SUBMITTED',
    'AWAITING_DOCUMENTS',
    'READY_FOR_REVIEW',
    'IN_REVIEW',
    'DECISION_PENDING',
    'ADMITTED',
    'DENIED',
    'WAITLISTED',
    'DEFERRED',
    'WITHDRAWN',
    'ENROLLED',
] as const;

export const ADMISSIONS_ROLES = [
    'ADMISSIONS_ADMIN',
    'RECRUITER',
    'APPLICATION_REVIEWER',
    'MARKETING_MANAGER',
    'EVENT_MANAGER',
    'REPORT_VIEWER',
    'APPLICANT',
] as const;

export type LeadScoreBand = (typeof LEAD_SCORE_BANDS)[number];
export type LeadScoreCriterionField =
    (typeof LEAD_SCORE_CRITERION_FIELDS)[number];
export type LeadScoreCriterionOperator =
    (typeof LEAD_SCORE_CRITERION_OPERATORS)[number];
export type AdmissionsRole = (typeof ADMISSIONS_ROLES)[number];

export interface LeadScoreCriterion {
    field: LeadScoreCriterionField;
    operator: LeadScoreCriterionOperator;
    values: string[];
}

export interface LeadScoreRule {
    id: string;
    name: string;
    description?: string;
    criteria: LeadScoreCriterion[];
    points: number;
    active: boolean;
    priority: number;
    dateCreated: string;
    dateUpdated: string;
}

export interface LeadScoreRuleRequest {
    name: string;
    description?: string | null;
    criteria: LeadScoreCriterion[];
    points: number;
    active: boolean;
    priority: number;
}

export interface LeadScoreRuleResult {
    ruleID: string;
    name: string;
    points: number;
    matched: boolean;
    reason: string;
}

export interface LeadScore {
    id: string;
    constituentID: string;
    totalScore: number;
    band: LeadScoreBand;
    breakdown: LeadScoreRuleResult[];
    recalculatedAt: string;
    dateCreated: string;
    dateUpdated: string;
}

export interface LeadScoreQuery {
    page?: number;
    rows?: number;
    orderBy?: string;
    band?: LeadScoreBand;
    min_score?: number;
}

export interface LeadScoreRuleQuery {
    page?: number;
    rows?: number;
    orderBy?: string;
    active?: boolean;
}

export interface Inquiry {
    id: string;
    constituentID: string;
    firstName: string;
    lastName: string;
    dateOfBirth: string;
    primaryEmail: string;
    primaryPhone: string;
    programOfInterest?: string;
    termOfInterest?: string;
    source: string;
    utmSource?: string;
    utmMedium?: string;
    utmCampaign?: string;
    message?: string;
    status: string;
    dateCreated: string;
    dateUpdated: string;
}

export interface InquiryRequest {
    firstName: string;
    lastName: string;
    dateOfBirth: string;
    primaryEmail: string;
    primaryPhone: string;
    programOfInterest?: string | null;
    termOfInterest?: string | null;
    source: string;
    utmSource?: string | null;
    utmMedium?: string | null;
    utmCampaign?: string | null;
    message?: string | null;
}
