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

export const APPLICATION_TYPES = [
    'KUCCPS_PLACEMENT',
    'SELF_SPONSORED_UNDERGRAD',
    'DIPLOMA',
    'MASTERS',
    'PHD',
    'TVET',
    'BRIDGING',
    'CERTIFICATE',
] as const;

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

export const SIS_SYNC_JOB_STATUSES = [
    'QUEUED',
    'RUNNING',
    'SUCCEEDED',
    'FAILED',
    'RETRY_READY',
] as const;

export const SIS_SYNC_EVENT_STATUSES = [
    'QUEUED',
    'PROCESSING',
    'SUCCEEDED',
    'FAILED',
    'RETRY_READY',
] as const;

export const SIS_SYNC_DIRECTIONS = ['INBOUND', 'OUTBOUND'] as const;

export const CUSTOM_FIELD_OWNERS = ['CONSTITUENT', 'APPLICATION'] as const;

export const CUSTOM_FIELD_DATA_TYPES = [
    'TEXT',
    'TEXTAREA',
    'NUMBER',
    'DATE',
    'SELECT',
    'BOOLEAN',
] as const;

export const IMPORT_BATCH_STATUSES = [
    'PREVIEWED',
    'VALIDATION_FAILED',
    'QUEUED',
    'PROCESSING',
    'COMPLETED',
    'FAILED',
] as const;

export const IMPORT_SOURCES = ['MANUAL_UPLOAD', 'SIS_EXPORT'] as const;

export const IMPORT_FILE_TYPES = ['CSV', 'XLSX'] as const;

export const IMPORT_TARGETS = ['CONSTITUENTS', 'APPLICATIONS'] as const;

export const SIS_SYNC_EVENT_TYPES = [
    'BATCH_TERMS_PULL',
    'BATCH_PROGRAMS_PULL',
    'BATCH_PERSON_MATCHES_PULL',
    'BATCH_ENROLLMENT_PULL',
    'APPLICATION_SUBMISSION',
    'APPLICATION_DECISION',
    'DOCUMENT_STATUS',
    'ENROLLMENT_INTENT',
    'KUCCPS_PLACEMENT_PULL',
    'KUCCPS_PLACEMENT_CONFIRM',
    'KNEC_RESULT_VERIFICATION',
    'IPRS_IDENTITY_VERIFICATION',
    'MPESA_STK_PUSH',
    'MPESA_C2B_CALLBACK',
    'MPESA_TRANSACTION_QUERY',
    'SMS_OUTBOUND',
    'SMS_DELIVERY_REPORT',
    'WHATSAPP_MESSAGE_SEND',
    'WHATSAPP_WEBHOOK_INBOUND',
] as const;

export const KENYA_ADAPTERS = [
    'KUCCPS',
    'KNEC',
    'IPRS',
    'MPESA_DARAJA',
    'CELCOM_AFRICA_SMS',
    'WHATSAPP_CLOUD',
] as const;

export const SIS_SYNC_OPERATIONS = [
    'PLACEMENT_PULL',
    'PLACEMENT_CONFIRM',
    'RESULT_VERIFICATION',
    'IDENTITY_VERIFICATION',
    'STK_PUSH',
    'C2B_CALLBACK',
    'TRANSACTION_QUERY',
    'SMS_OUTBOUND',
    'DELIVERY_REPORT_PULL',
    'WA_MESSAGE_SEND',
    'WA_WEBHOOK_INBOUND',
] as const;

export const NOTIFICATION_CHANNELS = ['SMS', 'WHATSAPP', 'EMAIL'] as const;

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
export type SisSyncJobStatus = (typeof SIS_SYNC_JOB_STATUSES)[number];
export type SisSyncEventStatus = (typeof SIS_SYNC_EVENT_STATUSES)[number];
export type SisSyncDirection = (typeof SIS_SYNC_DIRECTIONS)[number];
export type SisSyncEventType = (typeof SIS_SYNC_EVENT_TYPES)[number];
export type KenyaAdapter = (typeof KENYA_ADAPTERS)[number];
export type SisSyncOperation = (typeof SIS_SYNC_OPERATIONS)[number];
export type NotificationChannel = (typeof NOTIFICATION_CHANNELS)[number];
export type CustomFieldOwner = (typeof CUSTOM_FIELD_OWNERS)[number];
export type CustomFieldDataType = (typeof CUSTOM_FIELD_DATA_TYPES)[number];
export type ImportBatchStatus = (typeof IMPORT_BATCH_STATUSES)[number];
export type ImportSource = (typeof IMPORT_SOURCES)[number];
export type ImportFileType = (typeof IMPORT_FILE_TYPES)[number];
export type ImportTarget = (typeof IMPORT_TARGETS)[number];

export interface CustomFieldDefinition {
    id: string;
    owner: CustomFieldOwner;
    fieldKey: string;
    label: string;
    description?: string;
    dataType: CustomFieldDataType;
    required: boolean;
    options: string[];
    validation?: string;
    searchable: boolean;
    reportable: boolean;
    importable: boolean;
    exportable: boolean;
    displayOrder: number;
    active: boolean;
    dateCreated: string;
    dateUpdated: string;
}

export interface CustomFieldDefinitionRequest {
    owner: CustomFieldOwner;
    fieldKey: string;
    label: string;
    description?: string | null;
    dataType: CustomFieldDataType;
    required: boolean;
    options: string[];
    validation?: string | null;
    searchable: boolean;
    reportable: boolean;
    importable: boolean;
    exportable: boolean;
    displayOrder: number;
    active: boolean;
}

export interface CustomFieldDefinitionQuery {
    page?: number;
    rows?: number;
    orderBy?: string;
    owner?: CustomFieldOwner;
    active?: boolean;
}

export interface CustomFieldValue {
    id: string;
    definitionID: string;
    owner: CustomFieldOwner;
    ownerID: string;
    value: string;
    dateCreated: string;
    dateUpdated: string;
}

export interface CustomFieldValueRequest {
    definitionID: string;
    owner: CustomFieldOwner;
    ownerID: string;
    value: string;
}

export interface SisSyncJob {
    id: string;
    name: string;
    adapter: KenyaAdapter;
    operation: SisSyncOperation;
    status: SisSyncJobStatus;
    direction: SisSyncDirection;
    startedAt?: string;
    completedAt?: string;
    recordsPulled: number;
    recordsPushed: number;
    eventsRequeued: number;
    attempts: number;
    maxAttempts: number;
    nextRetryAt?: string;
    externalRef?: string;
    externalReceiptID?: string;
    errorCode?: string;
    errorDetail?: string;
    lastErrorAt?: string;
    failureReason?: string;
    retryable: boolean;
    owner: string;
    dateCreated: string;
    dateUpdated: string;
}

export interface SisSyncEvent {
    id: string;
    jobID?: string;
    adapter: KenyaAdapter;
    operation: SisSyncOperation;
    eventType: SisSyncEventType;
    status: SisSyncEventStatus;
    direction: SisSyncDirection;
    resourceType: string;
    resourceID: string;
    externalRef?: string;
    externalReceiptID?: string;
    payloadHash: string;
    attempts: number;
    maxAttempts: number;
    nextRetryAt?: string;
    errorCode?: string;
    errorDetail?: string;
    lastErrorAt?: string;
    failureReason?: string;
    auditMessage: string;
    dateCreated: string;
    dateUpdated: string;
}

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

export interface Program {
    id: string;
    externalSISID: string;
    name: string;
    code: string;
    description?: string;
    degreeLevel?: string;
    active: boolean;
    syncedAt?: string;
    dateCreated: string;
    dateUpdated: string;
}

export interface ProgramQuery {
    page?: number;
    rows?: number;
    orderBy?: string;
    active?: boolean;
}

export interface AcademicTerm {
    id: string;
    externalSISID: string;
    name: string;
    code: string;
    termType?: string;
    startDate: string;
    endDate: string;
    applicationStartDate?: string;
    applicationDeadline?: string;
    active: boolean;
    syncedAt?: string;
    dateCreated: string;
    dateUpdated: string;
}

export interface AcademicTermQuery {
    page?: number;
    rows?: number;
    orderBy?: string;
    active?: boolean;
}

export interface NotificationPreferences {
    smsOptIn: boolean;
    whatsAppOptIn: boolean;
    emailOptIn: boolean;
    priority: NotificationChannel[];
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

export interface KUCCPSPlacement {
    placementID: string;
    admissionNumber?: string;
    institutionCode: string;
    programmeCode: string;
    programmeName: string;
    placementYear: number;
    clusterCode?: string;
    clusterPoints?: number;
    weightedPointsNote?: string;
}

export interface ApplicationKCSESubject {
    subjectCode: string;
    grade: string;
    points: number;
}

export interface ApplicationKCSEResult {
    indexNumber: string;
    examYear: number;
    subjects: ApplicationKCSESubject[];
    meanGrade: string;
    meanPoints: number;
}

export interface Application {
    id: string;
    constituentID: string;
    programID: string;
    academicTermID: string;
    applicationType: ApplicationType;
    status: ApplicationStatus;
    kuccpsPlacement?: KUCCPSPlacement;
    kcseResult?: ApplicationKCSEResult;
    assignedReviewerID?: string;
    submittedAt?: string;
    dateCreated: string;
    dateUpdated: string;
}

export interface ApplicationRequest {
    constituentID: string;
    programID: string;
    academicTermID: string;
    applicationType: ApplicationType;
    kuccpsPlacement?: KUCCPSPlacement | null;
    kcseResult?: ApplicationKCSEResult | null;
    assignedReviewerID?: string | null;
}

export interface ApplicationQuery {
    page?: number;
    rows?: number;
    orderBy?: string;
    constituent_id?: string;
    program_id?: string;
    academic_term_id?: string;
    application_type?: ApplicationType;
    status?: ApplicationStatus;
}

export interface ApplicationTransition {
    id: string;
    applicationID: string;
    fromStatus: ApplicationStatus;
    toStatus: ApplicationStatus;
    actorID: string;
    reason?: string;
    note?: string;
    metadata?: Record<string, unknown>;
    dateCreated: string;
}

export interface ApplicationTransitionRequest {
    toStatus: Extract<ApplicationStatus, 'SUBMITTED' | 'WITHDRAWN'>;
    reason?: string | null;
    note?: string | null;
    metadata?: Record<string, unknown> | null;
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

export interface ImportBatch {
    id: string;
    source: ImportSource;
    fileType: ImportFileType;
    target: ImportTarget;
    status: ImportBatchStatus;
    fileName: string;
    storageKey?: string;
    uploadedByID: string;
    totalRows: number;
    validRows: number;
    invalidRows: number;
    duplicateRows: number;
    fieldMapping: Record<string, string>;
    invalidReportKey?: string;
    validationSummary?: string;
    committedAt?: string;
    dateCreated: string;
    dateUpdated: string;
}

export interface ImportBatchRequest {
    source: ImportSource;
    fileType: ImportFileType;
    target: ImportTarget;
    status: ImportBatchStatus;
    fileName: string;
    storageKey?: string | null;
    totalRows: number;
    validRows: number;
    invalidRows: number;
    duplicateRows: number;
    fieldMapping: Record<string, string>;
    invalidReportKey?: string | null;
    validationSummary?: string | null;
}

export interface ImportBatchQuery {
    page?: number;
    rows?: number;
    orderBy?: string;
    source?: ImportSource;
    file_type?: ImportFileType;
    target?: ImportTarget;
    status?: ImportBatchStatus;
    uploaded_by_id?: string;
}

export interface ImportInvalidRow {
    id: string;
    batchID: string;
    rowNumber: number;
    fieldName?: string;
    rawData: Record<string, string>;
    errorCode: string;
    errorDetail: string;
    dateCreated: string;
}

export interface ImportInvalidRowRequest {
    rowNumber: number;
    fieldName?: string | null;
    rawData: Record<string, string>;
    errorCode: string;
    errorDetail: string;
}

export interface ImportInvalidRowsRequest {
    rows: ImportInvalidRowRequest[];
}

export interface ImportInvalidRowQuery {
    page?: number;
    rows?: number;
    orderBy?: string;
    import_batch_id?: string;
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
