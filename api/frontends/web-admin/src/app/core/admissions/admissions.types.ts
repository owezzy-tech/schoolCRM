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

export const DOCUMENT_STATUSES = [
    'UPLOADED',
    'PENDING_REVIEW',
    'ACCEPTED',
    'REJECTED',
    'WAIVED',
    'EXPIRED',
    'SYNCED_TO_SIS',
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

export const LEAD_ASSIGNMENT_CRITERION_FIELDS = [
    'territory',
    'program_interest',
    'application_type',
    'academic_term',
    'lead_source',
    'event_attendance',
    'alpha_range',
] as const;

export const LEAD_ASSIGNMENT_CRITERION_OPERATORS = [
    'EQ',
    'IN',
    'BETWEEN',
] as const;

export const EVENT_ATTENDANCE_STATUSES = [
    'ANY',
    'REGISTERED',
    'ATTENDED',
    'NO_SHOW',
] as const;

export const APPLICATION_FEE_STATUSES = [
    'PENDING',
    'PAID',
    'FAILED',
    'WAIVED',
    'REFUNDED',
    'NOT_REQUIRED',
] as const;

export const APPLICATION_FEE_PROVIDERS = [
    'stripe',
    'square',
    'manual',
    'not_required',
] as const;

export type LeadScoreBand = (typeof LEAD_SCORE_BANDS)[number];
export type LeadScoreCriterionField =
    (typeof LEAD_SCORE_CRITERION_FIELDS)[number];
export type LeadScoreCriterionOperator =
    (typeof LEAD_SCORE_CRITERION_OPERATORS)[number];
export type AdmissionsRole = (typeof ADMISSIONS_ROLES)[number];
export type ApplicationStatus = (typeof APPLICATION_STATUSES)[number];
export type ApplicationType = (typeof APPLICATION_TYPES)[number];
export type DocumentStatus = (typeof DOCUMENT_STATUSES)[number];
export type LeadAssignmentCriterionField =
    (typeof LEAD_ASSIGNMENT_CRITERION_FIELDS)[number];
export type LeadAssignmentCriterionOperator =
    (typeof LEAD_ASSIGNMENT_CRITERION_OPERATORS)[number];
export type EventAttendanceStatus = (typeof EVENT_ATTENDANCE_STATUSES)[number];
export type ApplicationFeeStatus = (typeof APPLICATION_FEE_STATUSES)[number];
export type ApplicationFeeProvider = (typeof APPLICATION_FEE_PROVIDERS)[number];

export interface ApplicationFeeAuditEvent {
    actor: string;
    action: string;
    reason: string;
    timestamp: string;
}

export interface ApplicationFee {
    id: string;
    applicationID: string;
    amountCents: number;
    currency: string;
    status: ApplicationFeeStatus;
    provider: ApplicationFeeProvider;
    transactionID?: string;
    dueAt?: string;
    paidAt?: string;
    waiverGrantedBy?: string;
    waiverReason?: string;
    waiverGrantedAt?: string;
    auditTrail: ApplicationFeeAuditEvent[];
}

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

export interface LeadAssignmentCriterion {
    field: LeadAssignmentCriterionField;
    operator: LeadAssignmentCriterionOperator;
    values: string[];
}

export interface AssignmentAuditEvent {
    actor: string;
    action: string;
    reason: string;
    timestamp: string;
}

export interface LeadAssignmentRule {
    id: string;
    name: string;
    description: string;
    criteria: LeadAssignmentCriterion[];
    assignee: string;
    assigneeRole: Extract<AdmissionsRole, 'RECRUITER' | 'APPLICATION_REVIEWER'>;
    active: boolean;
    priority: number;
    lastMatchedCount: number;
    auditTrail: AssignmentAuditEvent[];
}

export interface AssignmentCandidate {
    id: string;
    applicant: string;
    program: string;
    term: string;
    source: string;
    territory: string;
    eventAttendance: EventAttendanceStatus;
    alphaKey: string;
    currentAssignee: string;
    recommendedAssignee: string;
    matchedRule: string;
    assignmentMode: 'AUTO' | 'MANUAL_OVERRIDE';
    overrideReason?: string;
    overrideActor?: string;
    overrideTimestamp?: string;
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

export interface ApplicationFormField {
    fieldName: string;
    fieldType: string;
    required: boolean;
    displayOrder: number;
    validation?: string;
}

export interface ApplicationChecklistTemplateItem {
    itemKey: string;
    documentName: string;
    description?: string;
    required: boolean;
    displayOrder: number;
}

export interface ApplicationFormTemplate {
    id: string;
    programID: string;
    academicTermID: string;
    applicationType: ApplicationType;
    name: string;
    description?: string;
    version: number;
    requiredFields: ApplicationFormField[];
    checklistItems: ApplicationChecklistTemplateItem[];
    active: boolean;
    priority: number;
    dateCreated: string;
    dateUpdated: string;
}

export interface ApplicationFormTemplateRequest {
    programID: string;
    academicTermID: string;
    applicationType: ApplicationType;
    name: string;
    description?: string | null;
    requiredFields: ApplicationFormField[];
    checklistItems: ApplicationChecklistTemplateItem[];
    active: boolean;
    priority: number;
}

export interface ApplicationFormTemplateQuery {
    page?: number;
    rows?: number;
    orderBy?: string;
    application_type?: ApplicationType;
    active?: boolean;
}

export interface ChecklistItem {
    id: string;
    applicationID: string;
    itemKey: string;
    documentName: string;
    description?: string;
    required: boolean;
    status: DocumentStatus;
    displayOrder: number;
    dateCreated: string;
    dateUpdated: string;
}

export interface ChecklistItemRequest {
    itemKey: string;
    documentName: string;
    description?: string | null;
    required: boolean;
    displayOrder: number;
}

export interface ChecklistItemQuery {
    page?: number;
    rows?: number;
    orderBy?: string;
    status?: DocumentStatus;
    required?: boolean;
}

export interface AdmissionsDocument {
    id: string;
    applicationID: string;
    checklistItemID: string;
    fileName: string;
    contentType: string;
    sizeBytes: number;
    storageKey: string;
    status: DocumentStatus;
    reviewerID?: string;
    reviewerNotes?: string;
    uploadedByID: string;
    uploadedAt: string;
    reviewedAt?: string;
    dateCreated: string;
    dateUpdated: string;
}

export interface AdmissionsDocumentRequest {
    checklistItemID: string;
    fileName: string;
    contentType: string;
    sizeBytes: number;
    storageKey: string;
}

export interface AdmissionsDocumentVerificationRequest {
    status: Extract<DocumentStatus, 'ACCEPTED' | 'REJECTED' | 'WAIVED'>;
    reviewerNotes?: string | null;
}

export interface AdmissionsDocumentQuery {
    page?: number;
    rows?: number;
    orderBy?: string;
    checklist_item_id?: string;
    status?: DocumentStatus;
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
