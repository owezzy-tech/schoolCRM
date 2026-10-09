import {
    LessonCitation,
    lessonCitations,
    lessonContentRecord,
} from './lessons.types';

export interface GenerationRequest {
    request_id: string;
    school_id: string;
    department_id: string;
    framework: string;
    stage: string;
    subject: string;
    revision: string;
    topic: string;
    duration_minutes: number;
}
export interface GenerationReceipt {
    planID: string;
    version: number;
    title: string;
}
export interface LessonThread {
    id: string;
    request: GenerationRequest;
    status: 'pending' | 'completed';
    result: GenerationReceipt | null;
}
export type AssistantEvent =
    | {
          type: 'progress';
          stage: 'retrieve' | 'generate' | 'persist';
          detail: string;
      }
    | { type: 'citation'; citations: LessonCitation[] }
    | { type: 'structured-result'; result: GenerationReceipt }
    | {
          type: 'approval-request';
          planID: string;
          version: number;
          action: 'submit-for-review';
      }
    | { type: 'completed' }
    | { type: 'terminal-error'; status: number; detail: string };

/** A client-safe stream failure; other runtime errors are never displayed. */
export class AssistantStreamError extends Error {}

/** Reject incompatible or misrouted events at the transport boundary. */
export function assistantEvent(
    name: string,
    value: unknown,
    requestID: string
): AssistantEvent {
    const data = lessonContentRecord(value);
    if (data['protocolVersion'] !== 1 || data['requestID'] !== requestID) {
        throw new AssistantStreamError(
            'The assistant response belongs to an incompatible request.'
        );
    }
    switch (name) {
        case 'progress':
            if (
                (data['stage'] === 'retrieve' ||
                    data['stage'] === 'generate' ||
                    data['stage'] === 'persist') &&
                typeof data['detail'] === 'string'
            ) {
                return {
                    type: name,
                    stage: data['stage'],
                    detail: data['detail'],
                };
            }
            break;
        case 'citation':
            if (Array.isArray(data['citations']))
                return { type: name, citations: lessonCitations(data) };
            break;
        case 'structured-result': {
            const result = lessonContentRecord(data['result']);
            if (
                typeof result['planID'] === 'string' &&
                typeof result['title'] === 'string' &&
                typeof result['version'] === 'number' &&
                Number.isInteger(result['version']) &&
                result['version'] > 0
            ) {
                return {
                    type: name,
                    result: {
                        planID: result['planID'],
                        title: result['title'],
                        version: result['version'],
                    },
                };
            }
            break;
        }
        case 'approval-request':
            if (
                typeof data['planID'] === 'string' &&
                typeof data['version'] === 'number' &&
                Number.isInteger(data['version']) &&
                data['version'] > 0 &&
                data['action'] === 'submit-for-review'
            ) {
                return {
                    type: name,
                    planID: data['planID'],
                    version: data['version'],
                    action: data['action'],
                };
            }
            break;
        case 'terminal-error':
            if (
                typeof data['status'] === 'number' &&
                typeof data['detail'] === 'string'
            )
                return {
                    type: name,
                    status: data['status'],
                    detail: data['detail'],
                };
            break;
        case 'completed':
            return { type: name };
    }
    throw new AssistantStreamError('The assistant returned an invalid event.');
}
