-- Version: 1.01
-- Description: Create table users
CREATE TABLE users (
	user_id       UUID        NOT NULL,
	name          TEXT        NOT NULL,
	email         TEXT UNIQUE NOT NULL,
	roles         TEXT[]      NOT NULL,
	password_hash TEXT        NOT NULL,
    department    TEXT        NULL,
    enabled       BOOLEAN     NOT NULL,
	date_created  TIMESTAMP   NOT NULL,
	date_updated  TIMESTAMP   NOT NULL,

	PRIMARY KEY (user_id)
);

-- Version: 1.02
-- Description: Create table products
CREATE TABLE products (
	product_id   UUID           NOT NULL,
    user_id      UUID           NOT NULL,
	name         TEXT           NOT NULL,
    cost         NUMERIC(10, 2) NOT NULL,
	quantity     INT            NOT NULL,
	date_created TIMESTAMP      NOT NULL,
	date_updated TIMESTAMP      NOT NULL,

	PRIMARY KEY (product_id),
	FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

-- Version: 1.03
-- Description: Add products view.
CREATE OR REPLACE VIEW view_products AS
SELECT
    p.product_id,
    p.user_id,
	p.name,
    p.cost,
	p.quantity,
    p.date_created,
    p.date_updated,
    u.name AS user_name
FROM
    products AS p
JOIN
    users AS u ON u.user_id = p.user_id

-- Version: 1.04
-- Description: Create table homes
CREATE TABLE homes (
    home_id       UUID       NOT NULL,
    type          TEXT       NOT NULL,
    user_id       UUID       NOT NULL,
    address_1     TEXT       NOT NULL,
    address_2     TEXT       NULL,
    zip_code      TEXT       NOT NULL,
    city          TEXT       NOT NULL,
    state         TEXT       NOT NULL,
    country       TEXT       NOT NULL,
    date_created  TIMESTAMP  NOT NULL,
    date_updated  TIMESTAMP  NOT NULL,

    PRIMARY KEY (home_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

-- Version: 1.05
-- Description: Create table audit
CREATE TABLE audit (
    id          UUID      NOT NULL,
    obj_id      UUID      NOT NULL,
    obj_domain  TEXT      NOT NULL,
    obj_name    TEXT      NOT NULL,
    actor_id    UUID      NOT NULL,
    action      TEXT      NOT NULL,
    data        JSONB     NULL,
    message     TEXT      NULL,
    timestamp   TIMESTAMP NOT NULL,

    PRIMARY KEY (id)
);

-- Version: 1.06
-- Description: Create admissions reference tables
CREATE TABLE admissions_programs (
    program_id      UUID      NOT NULL,
    external_sis_id TEXT      UNIQUE NOT NULL,
    name            TEXT      NOT NULL,
    code            TEXT      NOT NULL,
    description     TEXT      NULL,
    degree_level    TEXT      NULL,
    is_active       BOOLEAN   NOT NULL,
    synced_at       TIMESTAMP NULL,
    date_created    TIMESTAMP NOT NULL,
    date_updated    TIMESTAMP NOT NULL,

    PRIMARY KEY (program_id)
);

CREATE INDEX idx_admissions_programs_active ON admissions_programs (is_active);
CREATE INDEX idx_admissions_programs_code ON admissions_programs (code);

CREATE TABLE admissions_academic_terms (
    academic_term_id        UUID      NOT NULL,
    external_sis_id         TEXT      UNIQUE NOT NULL,
    name                    TEXT      NOT NULL,
    code                    TEXT      NOT NULL,
    term_type               TEXT      NULL,
    start_date              TIMESTAMP NOT NULL,
    end_date                TIMESTAMP NOT NULL,
    application_start_date  TIMESTAMP NULL,
    application_deadline    TIMESTAMP NULL,
    is_active               BOOLEAN   NOT NULL,
    synced_at               TIMESTAMP NULL,
    date_created            TIMESTAMP NOT NULL,
    date_updated            TIMESTAMP NOT NULL,

    PRIMARY KEY (academic_term_id),
    CONSTRAINT admissions_academic_terms_date_range CHECK (start_date < end_date),
    CONSTRAINT admissions_academic_terms_application_window CHECK (
        application_start_date IS NULL
        OR application_deadline IS NULL
        OR application_deadline >= application_start_date
    )
);

CREATE INDEX idx_admissions_academic_terms_active ON admissions_academic_terms (is_active);
CREATE INDEX idx_admissions_academic_terms_code ON admissions_academic_terms (code);

-- Version: 1.07
-- Description: Create admissions constituents table
CREATE TABLE admissions_constituents (
    constituent_id   UUID      NOT NULL,
    first_name       TEXT      NOT NULL,
    last_name        TEXT      NOT NULL,
    preferred_name   TEXT      NULL,
    middle_name      TEXT      NULL,
    suffix           TEXT      NULL,
    date_of_birth    TIMESTAMP NOT NULL,
    primary_email    TEXT      NOT NULL,
    primary_phone    TEXT      NOT NULL,
    external_sis_id  TEXT      UNIQUE NULL,
    lifecycle_stage  TEXT      NOT NULL,
    duplicate_status TEXT      NOT NULL,
    duplicate_of_id  UUID      NULL,
    sis_synced_at    TIMESTAMP NULL,
    date_created     TIMESTAMP NOT NULL,
    date_updated     TIMESTAMP NOT NULL,

    PRIMARY KEY (constituent_id),
    FOREIGN KEY (duplicate_of_id) REFERENCES admissions_constituents(constituent_id),
    CONSTRAINT admissions_constituents_lifecycle_stage CHECK (lifecycle_stage IN ('PROSPECT', 'INQUIRY', 'APPLICANT', 'ADMITTED', 'ENROLLED', 'ALUMNI')),
    CONSTRAINT admissions_constituents_duplicate_status CHECK (duplicate_status IN ('ACTIVE', 'MERGED', 'DUPLICATE_OF')),
    CONSTRAINT admissions_constituents_duplicate_consistency CHECK (
        (duplicate_status = 'ACTIVE' AND duplicate_of_id IS NULL)
        OR (duplicate_status IN ('MERGED', 'DUPLICATE_OF') AND duplicate_of_id IS NOT NULL)
    )
);

CREATE INDEX idx_admissions_constituents_email ON admissions_constituents (primary_email);
CREATE INDEX idx_admissions_constituents_phone ON admissions_constituents (primary_phone);
CREATE INDEX idx_admissions_constituents_lifecycle ON admissions_constituents (lifecycle_stage);
CREATE INDEX idx_admissions_constituents_duplicate_status ON admissions_constituents (duplicate_status);

-- Version: 1.08
-- Description: Create admissions duplicate review queue
CREATE TABLE admissions_duplicate_reviews (
    duplicate_review_id       UUID      NOT NULL,
    source_constituent_id     UUID      NOT NULL,
    candidate_constituent_id  UUID      NOT NULL,
    match_type                TEXT      NOT NULL,
    match_score               INT       NOT NULL,
    match_reason              TEXT      NOT NULL,
    status                    TEXT      NOT NULL,
    resolved_by               UUID      NULL,
    resolved_at               TIMESTAMP NULL,
    resolution_note           TEXT      NULL,
    date_created              TIMESTAMP NOT NULL,
    date_updated              TIMESTAMP NOT NULL,

    PRIMARY KEY (duplicate_review_id),
    FOREIGN KEY (source_constituent_id) REFERENCES admissions_constituents(constituent_id),
    FOREIGN KEY (candidate_constituent_id) REFERENCES admissions_constituents(constituent_id),
    CONSTRAINT admissions_duplicate_reviews_distinct_pair CHECK (source_constituent_id <> candidate_constituent_id),
    CONSTRAINT admissions_duplicate_reviews_match_type CHECK (match_type IN ('EXACT', 'FUZZY')),
    CONSTRAINT admissions_duplicate_reviews_match_score CHECK (match_score >= 0 AND match_score <= 100),
    CONSTRAINT admissions_duplicate_reviews_status CHECK (status IN ('PENDING', 'LINKED', 'MERGED', 'REJECTED', 'DEFERRED')),
    CONSTRAINT admissions_duplicate_reviews_resolution_consistency CHECK (
        (status = 'PENDING' AND resolved_by IS NULL AND resolved_at IS NULL)
        OR (status IN ('LINKED', 'MERGED', 'REJECTED', 'DEFERRED') AND resolved_by IS NOT NULL AND resolved_at IS NOT NULL)
    ),
    CONSTRAINT admissions_duplicate_reviews_unique_pair UNIQUE (source_constituent_id, candidate_constituent_id)
);

CREATE INDEX idx_admissions_duplicate_reviews_source ON admissions_duplicate_reviews (source_constituent_id);
CREATE INDEX idx_admissions_duplicate_reviews_candidate ON admissions_duplicate_reviews (candidate_constituent_id);
CREATE INDEX idx_admissions_duplicate_reviews_status ON admissions_duplicate_reviews (status);
CREATE INDEX idx_admissions_duplicate_reviews_match_type ON admissions_duplicate_reviews (match_type);

-- Version: 1.09
-- Description: Create admissions applications table
CREATE TABLE admissions_applications (
    application_id        UUID      NOT NULL,
    constituent_id        UUID      NOT NULL,
    program_id            UUID      NOT NULL,
    academic_term_id      UUID      NOT NULL,
    application_type      TEXT      NOT NULL,
    status                TEXT      NOT NULL,
    assigned_reviewer_id  UUID      NULL,
    submitted_at          TIMESTAMP NULL,
    date_created          TIMESTAMP NOT NULL,
    date_updated          TIMESTAMP NOT NULL,

    PRIMARY KEY (application_id),
    FOREIGN KEY (constituent_id) REFERENCES admissions_constituents(constituent_id),
    FOREIGN KEY (program_id) REFERENCES admissions_programs(program_id),
    FOREIGN KEY (academic_term_id) REFERENCES admissions_academic_terms(academic_term_id),
    CONSTRAINT admissions_applications_type CHECK (application_type IN ('FRESHMAN', 'TRANSFER', 'GRADUATE')),
    CONSTRAINT admissions_applications_status CHECK (status IN ('DRAFT', 'SUBMITTED', 'AWAITING_DOCUMENTS', 'READY_FOR_REVIEW', 'IN_REVIEW', 'DECISION_PENDING', 'ADMITTED', 'DENIED', 'WAITLISTED', 'DEFERRED', 'WITHDRAWN', 'ENROLLED'))
);

CREATE UNIQUE INDEX idx_admissions_applications_active_tuple
    ON admissions_applications (constituent_id, academic_term_id, program_id)
    WHERE status NOT IN ('DENIED', 'WITHDRAWN', 'ENROLLED');
CREATE INDEX idx_admissions_applications_constituent ON admissions_applications (constituent_id);
CREATE INDEX idx_admissions_applications_program ON admissions_applications (program_id);
CREATE INDEX idx_admissions_applications_academic_term ON admissions_applications (academic_term_id);
CREATE INDEX idx_admissions_applications_status ON admissions_applications (status);

-- Version: 1.10
-- Description: Create admissions application transition history table
CREATE TABLE admissions_application_transitions (
    application_transition_id UUID      NOT NULL,
    application_id            UUID      NOT NULL,
    from_status               TEXT      NOT NULL,
    to_status                 TEXT      NOT NULL,
    actor_id                  UUID      NOT NULL,
    reason                    TEXT      NULL,
    note                      TEXT      NULL,
    metadata                  JSONB     NULL,
    date_created              TIMESTAMP NOT NULL,

    PRIMARY KEY (application_transition_id),
    FOREIGN KEY (application_id) REFERENCES admissions_applications(application_id),
    CONSTRAINT admissions_application_transitions_from_status CHECK (from_status IN ('DRAFT', 'SUBMITTED', 'AWAITING_DOCUMENTS', 'READY_FOR_REVIEW', 'IN_REVIEW', 'DECISION_PENDING', 'ADMITTED', 'DENIED', 'WAITLISTED', 'DEFERRED', 'WITHDRAWN', 'ENROLLED')),
    CONSTRAINT admissions_application_transitions_to_status CHECK (to_status IN ('DRAFT', 'SUBMITTED', 'AWAITING_DOCUMENTS', 'READY_FOR_REVIEW', 'IN_REVIEW', 'DECISION_PENDING', 'ADMITTED', 'DENIED', 'WAITLISTED', 'DEFERRED', 'WITHDRAWN', 'ENROLLED'))
);

CREATE INDEX idx_admissions_application_transitions_application ON admissions_application_transitions (application_id);
CREATE INDEX idx_admissions_application_transitions_actor ON admissions_application_transitions (actor_id);
CREATE INDEX idx_admissions_application_transitions_created ON admissions_application_transitions (date_created);

-- Version: 1.11
-- Description: Create admissions staff profiles
CREATE TABLE admissions_staff_profiles (
    staff_profile_id  UUID      NOT NULL,
    user_id           UUID      NOT NULL,
    admissions_roles  TEXT[]    NOT NULL,
    is_active         BOOLEAN   NOT NULL,
    date_created      TIMESTAMP NOT NULL,
    date_updated      TIMESTAMP NOT NULL,

    PRIMARY KEY (staff_profile_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    CONSTRAINT admissions_staff_profiles_unique_user UNIQUE (user_id),
    CONSTRAINT admissions_staff_profiles_roles_not_empty CHECK (cardinality(admissions_roles) > 0),
    CONSTRAINT admissions_staff_profiles_roles CHECK (admissions_roles <@ ARRAY[
        'ADMISSIONS_ADMIN',
        'RECRUITER',
        'APPLICATION_REVIEWER',
        'MARKETING_MANAGER',
        'EVENT_MANAGER',
        'REPORT_VIEWER',
        'APPLICANT'
    ]::TEXT[])
);

CREATE INDEX idx_admissions_staff_profiles_user ON admissions_staff_profiles (user_id);
CREATE INDEX idx_admissions_staff_profiles_roles ON admissions_staff_profiles USING GIN (admissions_roles);
CREATE INDEX idx_admissions_staff_profiles_active ON admissions_staff_profiles (is_active);
