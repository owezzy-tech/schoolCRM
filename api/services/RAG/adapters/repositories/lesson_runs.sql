CREATE TABLE IF NOT EXISTS rag_lesson_runs (
    actor_id UUID NOT NULL,
    request_id UUID NOT NULL,
    input_hash TEXT NOT NULL,
    request JSONB NOT NULL,
    result JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (actor_id, request_id)
);
