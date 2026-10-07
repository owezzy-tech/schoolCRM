import { Capability } from 'app/core/school-access/school-access.types';

export type LessonStatus =
    | 'draft'
    | 'hod_review'
    | 'dean_approval'
    | 'changes_requested'
    | 'approved'
    | 'published';

export type LessonDecision = 'approve' | 'request_changes';

export const LESSON_STATUS_LABELS: Record<LessonStatus, string> = {
    draft: 'Draft',
    hod_review: 'In HOD review',
    dean_approval: 'In dean approval',
    changes_requested: 'Changes requested',
    approved: 'Approved',
    published: 'Published',
};

export type LessonStatusTone = 'neutral' | 'info' | 'warning' | 'success';

export const LESSON_STATUS_TONES: Record<LessonStatus, LessonStatusTone> = {
    draft: 'neutral',
    hod_review: 'info',
    dean_approval: 'info',
    changes_requested: 'warning',
    approved: 'success',
    published: 'success',
};

export interface LessonActivity {
    minutes: number;
    title: string;
    detail: string;
}

/**
 * Lesson content shape used by this UI until schoolCRM-o58.3 confirms the
 * curriculum schema; the API currently accepts any JSON object.
 */
export interface LessonContent {
    objectives: string[];
    materials: string[];
    activities: LessonActivity[];
    assessment: string;
    prerequisites?: string[];
    differentiation?: string;
}

export interface LessonCitation {
    sourceID: string;
    page: number;
    title: string;
    sourceURL: string;
    revision: string;
}

/** Keep the complete snapshot when editing. Read-only metadata must survive saves. */
export function lessonContentRecord(value: unknown): Record<string, unknown> {
    return value !== null && typeof value === 'object' && !Array.isArray(value)
        ? (value as Record<string, unknown>)
        : {};
}

export function lessonCitations(value: unknown): LessonCitation[] {
    const entries = lessonContentRecord(value)['citations'];
    if (!Array.isArray(entries)) return [];
    return entries.filter((entry): entry is LessonCitation => {
        const citation = lessonContentRecord(entry);
        return (
            typeof citation['sourceID'] === 'string' &&
            typeof citation['title'] === 'string' &&
            typeof citation['page'] === 'number' &&
            Number.isInteger(citation['page']) &&
            citation['page'] > 0 &&
            typeof citation['revision'] === 'string' &&
            typeof citation['sourceURL'] === 'string' &&
            /^https?:\/\//.test(citation['sourceURL'])
        );
    });
}

export interface LessonPlan {
    id: string;
    schoolID: string;
    departmentID: string;
    authorID: string;
    title: string;
    status: LessonStatus;
    /** Newest version this viewer may see; other authors' drafts are hidden. */
    currentVersion: number;
    publishedVersion: number | null;
    dateCreated: string;
    dateUpdated: string;
}

export interface LessonVersion {
    id: string;
    planID: string;
    version: number;
    title: string;
    content: unknown;
    changeSummary: string;
    authorID: string;
    status: LessonStatus;
    reviewerID: string | null;
    reviewedAt: string | null;
    approverID: string | null;
    approvedAt: string | null;
    publishedAt: string | null;
    reviewFeedback: string;
    approvalFeedback: string;
    dateCreated: string;
}

export interface LessonDraft {
    title: string;
    content: LessonContent;
    changeSummary: string;
}

export type LessonAction = 'edit' | 'submit' | 'review' | 'approve' | 'publish';

export interface LessonActions {
    actions: LessonAction[];
    /** The viewer holds the capability but authored or already reviewed this version. */
    separationOfDuties: boolean;
}

export const SEPARATION_OF_DUTIES_MESSAGE =
    'Author, HOD reviewer and dean approver must be three different people.';

/**
 * Mirrors the server's workflow rules to decide which controls to offer for
 * the plan's current version. The server re-checks every command.
 */
export function lessonActions(
    plan: LessonPlan,
    version: LessonVersion,
    viewerID: string,
    capabilities: Capability[]
): LessonActions {
    const can = (capability: Capability) => capabilities.includes(capability);
    const isAuthor = plan.authorID === viewerID;
    const actions: LessonAction[] = [];
    let separationOfDuties = false;

    if (isAuthor && can('teach')) {
        actions.push('edit');
        if (version.status === 'draft') actions.push('submit');
        if (version.status === 'approved') actions.push('publish');
    }
    if (can('review_lessons') && version.status === 'hod_review') {
        if (isAuthor) separationOfDuties = true;
        else actions.push('review');
    }
    if (can('approve_lessons') && version.status === 'dean_approval') {
        if (isAuthor || version.reviewerID === viewerID)
            separationOfDuties = true;
        else actions.push('approve');
    }

    return { actions, separationOfDuties };
}

/** Reads stored content defensively: older or foreign shapes render as empty sections. */
export function toLessonContent(value: unknown): LessonContent {
    const source = lessonContentRecord(value);
    const strings = (item: unknown): string[] =>
        Array.isArray(item)
            ? item.filter((entry): entry is string => typeof entry === 'string')
            : [];
    const activities = Array.isArray(source['activities'])
        ? source['activities'].flatMap((entry): LessonActivity[] => {
              if (entry === null || typeof entry !== 'object') return [];
              const activity = entry as Record<string, unknown>;
              return [
                  {
                      minutes: Number(activity['minutes']) || 0,
                      title: String(activity['title'] ?? ''),
                      detail: String(activity['detail'] ?? ''),
                  },
              ];
          })
        : [];

    return {
        ...('prerequisites' in source
            ? { prerequisites: strings(source['prerequisites']) }
            : {}),
        ...('differentiation' in source
            ? {
                  differentiation:
                      typeof source['differentiation'] === 'string'
                          ? source['differentiation']
                          : '',
              }
            : {}),
        objectives: strings(source['objectives']),
        materials: strings(source['materials']),
        activities,
        assessment:
            typeof source['assessment'] === 'string'
                ? source['assessment']
                : '',
    };
}
