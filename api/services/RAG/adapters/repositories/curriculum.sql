-- Python-owned retrieval data. This does not write Go-owned business tables.
CREATE EXTENSION IF NOT EXISTS vector;
CREATE TABLE IF NOT EXISTS rag_curriculum_sources (
    id UUID PRIMARY KEY,
    school_id UUID NOT NULL,
    department_id UUID NOT NULL,
    framework TEXT NOT NULL,
    stage TEXT NOT NULL,
    subject TEXT NOT NULL,
    revision TEXT NOT NULL,
    sha256 TEXT NOT NULL,
    metadata JSONB NOT NULL,
    original BYTEA NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending','approved','withdrawn')),
    reviewed_by TEXT,
    review_note TEXT,
    model_identity TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CHECK (status = 'pending' OR (reviewed_by IS NOT NULL AND model_identity IS NOT NULL))
);
CREATE UNIQUE INDEX IF NOT EXISTS rag_curriculum_original_scope ON rag_curriculum_sources
    (school_id, department_id, framework, stage, subject, revision, sha256,
     (metadata->>'first_page'), (metadata->>'last_page'));
CREATE TABLE IF NOT EXISTS rag_curriculum_chunks (
    source_id UUID NOT NULL REFERENCES rag_curriculum_sources(id),
    ordinal INTEGER NOT NULL CHECK (ordinal >= 0),
    page INTEGER NOT NULL CHECK (page > 0),
    content TEXT NOT NULL,
    embedding vector(1024),
    PRIMARY KEY (source_id, ordinal)
);
CREATE TABLE IF NOT EXISTS rag_curriculum_audit (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    source_id UUID NOT NULL REFERENCES rag_curriculum_sources(id),
    actor TEXT NOT NULL,
    action TEXT NOT NULL,
    note TEXT NOT NULL,
    occurred_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS rag_curriculum_scope ON rag_curriculum_sources
    (school_id, department_id, framework, stage, subject, revision, status, model_identity);
