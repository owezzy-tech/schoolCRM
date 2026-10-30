import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import {
    JsonApiCollectionDocument,
    JsonApiDocument,
    unwrapJsonApiCollection,
    unwrapJsonApiResource,
} from 'app/core/api/json-api';
import { TokenStorageService } from 'app/core/auth/token-storage.service';
import { map, Observable } from 'rxjs';
import {
    AssistantEvent,
    assistantEvent,
    AssistantStreamError,
    GenerationRequest,
    LessonThread,
} from './lesson-assistant.types';

@Injectable({ providedIn: 'root' })
export class LessonAssistantService {
    private readonly http = inject(HttpClient);
    private readonly tokens = inject(TokenStorageService);

    threads(
        schoolID: string,
        departmentID: string
    ): Observable<LessonThread[]> {
        return this.http
            .get<
                JsonApiCollectionDocument<LessonThread>
            >('/v1/rag/lessons/threads', { params: { school_id: schoolID, department_id: departmentID } })
            .pipe(map((document) => unwrapJsonApiCollection(document).items));
    }
    thread(requestID: string): Observable<LessonThread> {
        return this.http
            .get<
                JsonApiDocument<LessonThread>
            >(`/v1/rag/lessons/threads/${encodeURIComponent(requestID)}`)
            .pipe(map(unwrapJsonApiResource));
    }
    /** Native fetch supports streamed POST with bearer headers; unsubscribe cancels the request. */
    generate(request: GenerationRequest): Observable<AssistantEvent> {
        return new Observable((subscriber) => {
            const abort = new AbortController();
            const token = this.tokens.get(); // Always read fresh on reconnect.
            const consume = async () => {
                if (!token)
                    throw new HttpErrorResponse({
                        status: 401,
                        error: {
                            errors: [
                                {
                                    detail: 'Sign in to continue your lesson request.',
                                },
                            ],
                        },
                    });
                const response = await fetch(
                    '/v1/rag/lessons/generate/stream',
                    {
                        method: 'POST',
                        signal: abort.signal,
                        credentials: 'same-origin',
                        headers: {
                            'Content-Type': 'application/json',
                            Accept: 'text/event-stream',
                            Authorization: `Bearer ${token}`,
                        },
                        body: JSON.stringify(request),
                    }
                );
                if (!response.ok)
                    throw new HttpErrorResponse({
                        status: response.status,
                        // Proxies may answer with HTML; keep only JSON:API details.
                        error: await response.json().catch(() => null),
                    });
                if (
                    !response.headers
                        .get('content-type')
                        ?.startsWith('text/event-stream') ||
                    !response.body
                )
                    throw new AssistantStreamError(
                        'The assistant stream is unavailable.'
                    );
                let terminal = false;
                let savedResult = false;
                for await (const event of readAssistantStream(
                    response.body,
                    request.request_id
                )) {
                    if (event.type === 'structured-result') savedResult = true;
                    if (event.type === 'completed' && !savedResult)
                        throw new AssistantStreamError(
                            'The assistant completed without a saved draft receipt. Reconnect this request.'
                        );
                    subscriber.next(event);
                    if (
                        event.type === 'completed' ||
                        event.type === 'terminal-error'
                    ) {
                        terminal = true;
                        break;
                    }
                }
                if (!terminal)
                    throw new AssistantStreamError(
                        'Connection interrupted. Reconnect to recover this same request.'
                    );
                subscriber.complete();
            };
            consume().catch((error) => {
                if (!subscriber.closed) subscriber.error(error);
            });
            return () => abort.abort();
        });
    }
}

/** SSE framing tolerates UTF-8, CRLF and arbitrary network chunk boundaries. */
export async function* readAssistantStream(
    body: ReadableStream<Uint8Array>,
    requestID: string
): AsyncGenerator<AssistantEvent> {
    const reader = body.getReader();
    const decoder = new TextDecoder();
    let buffer = '';
    try {
        while (true) {
            const chunk = await reader.read();
            buffer += decoder.decode(chunk.value, { stream: !chunk.done });
            let boundary: RegExpExecArray | null;
            while ((boundary = /\r?\n\r?\n/.exec(buffer))) {
                if (boundary.index > 262144)
                    throw new AssistantStreamError(
                        'The assistant event is too large.'
                    );
                const frame = buffer.slice(0, boundary.index);
                buffer = buffer.slice(boundary.index + boundary[0].length);
                let name = 'message';
                const data: string[] = [];
                for (const line of frame.split(/\r?\n/)) {
                    if (line.startsWith('event:')) name = line.slice(6).trim();
                    if (line.startsWith('data:'))
                        data.push(line.slice(5).replace(/^ /, ''));
                }
                if (data.length)
                    yield assistantEvent(
                        name,
                        JSON.parse(data.join('\n')),
                        requestID
                    );
            }
            if (buffer.length > 262144)
                throw new AssistantStreamError(
                    'The assistant event is too large.'
                );
            if (chunk.done) break;
        }
    } finally {
        await reader.cancel().catch(() => undefined);
        reader.releaseLock();
    }
}
