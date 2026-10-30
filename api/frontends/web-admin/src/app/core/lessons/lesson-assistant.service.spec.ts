import { provideHttpClient } from '@angular/common/http';
import { TestBed } from '@angular/core/testing';
import { TokenStorageService } from 'app/core/auth/token-storage.service';
import { firstValueFrom, toArray } from 'rxjs';
import { afterEach, describe, expect, it, vi } from 'vitest';
import {
    LessonAssistantService,
    readAssistantStream,
} from './lesson-assistant.service';
import { assistantEvent, GenerationRequest } from './lesson-assistant.types';

const request: GenerationRequest = {
    request_id: 'request-1',
    school_id: 'school-1',
    department_id: 'dept-1',
    framework: 'kenya-cbc',
    stage: 'grade-1',
    subject: 'english',
    revision: '2024',
    topic: 'School vocabulary',
    duration_minutes: 30,
};
const frame = (name: string, fields = {}) =>
    `event: ${name}\r\ndata: ${JSON.stringify({ protocolVersion: 1, requestID: request.request_id, ...fields })}\r\n\r\n`;
function body(text: string, everyByte = false) {
    const bytes = new TextEncoder().encode(text);
    return new ReadableStream<Uint8Array>({
        start(controller) {
            if (everyByte)
                for (const byte of bytes)
                    controller.enqueue(new Uint8Array([byte]));
            else controller.enqueue(bytes);
            controller.close();
        },
    });
}
afterEach(() => vi.unstubAllGlobals());
describe('lesson assistant transport', () => {
    it('decodes fragmented UTF-8 and CRLF with heartbeat comments', async () => {
        const events = [];
        const text =
            ': heartbeat\r\n\r\n' +
            frame('progress', {
                stage: 'generate',
                detail: 'Drafting Kiswahili — shuleni',
            }) +
            frame('completed');
        for await (const event of readAssistantStream(
            body(text, true),
            request.request_id
        ))
            events.push(event);
        expect(events).toEqual([
            {
                type: 'progress',
                stage: 'generate',
                detail: 'Drafting Kiswahili — shuleni',
            },
            { type: 'completed' },
        ]);
    });
    it('rejects another request and incompatible protocol instead of displaying them', () => {
        expect(() =>
            assistantEvent(
                'completed',
                { protocolVersion: 1, requestID: 'someone-else' },
                request.request_id
            )
        ).toThrow('incompatible request');
        expect(() =>
            assistantEvent(
                'completed',
                { protocolVersion: 2, requestID: request.request_id },
                request.request_id
            )
        ).toThrow('incompatible request');
        expect(() =>
            assistantEvent(
                'structured-result',
                {
                    protocolVersion: 1,
                    requestID: request.request_id,
                    result: { planID: 'p', title: 'x', version: 0 },
                },
                request.request_id
            )
        ).toThrow('invalid event');
    });
    it('reconnects with the same input and a fresh bearer token', async () => {
        let token = 'first-token';
        TestBed.configureTestingModule({
            providers: [
                provideHttpClient(),
                {
                    provide: TokenStorageService,
                    useValue: { get: () => token },
                },
            ],
        });
        const fetchMock = vi.fn(() =>
            Promise.resolve(
                new Response(
                    body(
                        frame('structured-result', {
                            result: { planID: 'p', title: 'Saved', version: 1 },
                        }) + frame('completed')
                    ),
                    {
                        headers: { 'Content-Type': 'text/event-stream' },
                    }
                )
            )
        );
        vi.stubGlobal('fetch', fetchMock);
        const service = TestBed.inject(LessonAssistantService);
        await firstValueFrom(service.generate(request).pipe(toArray()));
        token = 'fresh-token';
        await firstValueFrom(service.generate(request).pipe(toArray()));
        const first = fetchMock.mock.calls[0] as unknown as [
            string,
            RequestInit,
        ];
        const second = fetchMock.mock.calls[1] as unknown as [
            string,
            RequestInit,
        ];
        expect(first[1].body).toBe(second[1].body);
        expect(second[1].headers).toMatchObject({
            Authorization: 'Bearer fresh-token',
        });
    });
    it('rejects terminal success without a saved receipt', async () => {
        TestBed.configureTestingModule({
            providers: [
                provideHttpClient(),
                {
                    provide: TokenStorageService,
                    useValue: { get: () => 'token' },
                },
            ],
        });
        vi.stubGlobal(
            'fetch',
            vi.fn(() =>
                Promise.resolve(
                    new Response(body(frame('completed')), {
                        headers: { 'Content-Type': 'text/event-stream' },
                    })
                )
            )
        );
        await expect(
            firstValueFrom(
                TestBed.inject(LessonAssistantService)
                    .generate(request)
                    .pipe(toArray())
            )
        ).rejects.toThrow('without a saved draft receipt');
    });
    it('treats EOF without terminal success as interrupted, even after a draft receipt', async () => {
        TestBed.configureTestingModule({
            providers: [
                provideHttpClient(),
                {
                    provide: TokenStorageService,
                    useValue: { get: () => 'token' },
                },
            ],
        });
        vi.stubGlobal(
            'fetch',
            vi.fn(() =>
                Promise.resolve(
                    new Response(
                        body(
                            frame('structured-result', {
                                result: {
                                    planID: 'saved',
                                    title: 'Saved',
                                    version: 1,
                                },
                            })
                        ),
                        { headers: { 'Content-Type': 'text/event-stream' } }
                    )
                )
            )
        );
        await expect(
            firstValueFrom(
                TestBed.inject(LessonAssistantService)
                    .generate(request)
                    .pipe(toArray())
            )
        ).rejects.toThrow('Connection interrupted');
    });
    it('preserves JSON:API permission failures and cancels when unsubscribed', async () => {
        TestBed.configureTestingModule({
            providers: [
                provideHttpClient(),
                {
                    provide: TokenStorageService,
                    useValue: { get: () => 'token' },
                },
            ],
        });
        vi.stubGlobal(
            'fetch',
            vi.fn(() =>
                Promise.resolve(
                    new Response(
                        JSON.stringify({
                            errors: [{ detail: 'Teaching access revoked' }],
                        }),
                        { status: 403 }
                    )
                )
            )
        );
        await expect(
            firstValueFrom(
                TestBed.inject(LessonAssistantService).generate(request)
            )
        ).rejects.toMatchObject({
            status: 403,
            error: { errors: [{ detail: 'Teaching access revoked' }] },
        });
        let signal: AbortSignal | undefined;
        vi.stubGlobal(
            'fetch',
            vi.fn((_url: string, options: RequestInit) => {
                signal = options.signal as AbortSignal;
                return new Promise<Response>(() => {});
            })
        );
        const subscription = TestBed.inject(LessonAssistantService)
            .generate(request)
            .subscribe();
        subscription.unsubscribe();
        expect(signal?.aborted).toBe(true);
    });
    it('keeps non-JSON proxy failures as status-only HTTP errors', async () => {
        TestBed.configureTestingModule({
            providers: [
                provideHttpClient(),
                {
                    provide: TokenStorageService,
                    useValue: { get: () => 'token' },
                },
            ],
        });
        vi.stubGlobal(
            'fetch',
            vi.fn(() =>
                Promise.resolve(
                    new Response('<html>Bad gateway at internal-host</html>', {
                        status: 502,
                    })
                )
            )
        );
        await expect(
            firstValueFrom(
                TestBed.inject(LessonAssistantService).generate(request)
            )
        ).rejects.toMatchObject({ status: 502, error: null });
    });
});
