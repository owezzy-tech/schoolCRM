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
-- Description: Create admissions inquiries table
CREATE TABLE admissions_inquiries (
    inquiry_id           UUID      NOT NULL,
    constituent_id       UUID      NOT NULL,
    first_name           TEXT      NOT NULL,
    last_name            TEXT      NOT NULL,
    date_of_birth        TIMESTAMP NOT NULL,
    primary_email        TEXT      NOT NULL,
    primary_phone        TEXT      NOT NULL,
    program_of_interest  UUID      NULL,
    term_of_interest     UUID      NULL,
    source               TEXT      NOT NULL,
    utm_source           TEXT      NULL,
    utm_medium           TEXT      NULL,
    utm_campaign         TEXT      NULL,
    message              TEXT      NULL,
    status               TEXT      NOT NULL,
    date_created         TIMESTAMP NOT NULL,
    date_updated         TIMESTAMP NOT NULL,

    PRIMARY KEY (inquiry_id),
    FOREIGN KEY (constituent_id) REFERENCES admissions_constituents(constituent_id),
    FOREIGN KEY (program_of_interest) REFERENCES admissions_programs(program_id),
    FOREIGN KEY (term_of_interest) REFERENCES admissions_academic_terms(academic_term_id),
    CONSTRAINT admissions_inquiries_status CHECK (status IN ('NEW', 'CONTACTED', 'CONVERTED', 'CLOSED')),
    CONSTRAINT admissions_inquiries_source_not_empty CHECK (trim(source) <> '')
);

CREATE INDEX idx_admissions_inquiries_constituent ON admissions_inquiries (constituent_id);
CREATE INDEX idx_admissions_inquiries_email ON admissions_inquiries (primary_email);
CREATE INDEX idx_admissions_inquiries_source ON admissions_inquiries (source);
CREATE INDEX idx_admissions_inquiries_status ON admissions_inquiries (status);
CREATE INDEX idx_admissions_inquiries_created ON admissions_inquiries (date_created);

-- Version: 1.09
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

-- Version: 1.10
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

-- Version: 1.11
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

-- Version: 1.12
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

-- Version: 1.13
-- Description: Create admissions applicant profiles
CREATE TABLE admissions_applicant_profiles (
    applicant_profile_id  UUID      NOT NULL,
    user_id               UUID      NOT NULL,
    constituent_id        UUID      NOT NULL,
    is_active             BOOLEAN   NOT NULL,
    date_created          TIMESTAMP NOT NULL,
    date_updated          TIMESTAMP NOT NULL,

    PRIMARY KEY (applicant_profile_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (constituent_id) REFERENCES admissions_constituents(constituent_id) ON DELETE CASCADE,
    CONSTRAINT admissions_applicant_profiles_unique_user UNIQUE (user_id),
    CONSTRAINT admissions_applicant_profiles_unique_constituent UNIQUE (constituent_id)
);

CREATE INDEX idx_admissions_applicant_profiles_user ON admissions_applicant_profiles (user_id);
CREATE INDEX idx_admissions_applicant_profiles_constituent ON admissions_applicant_profiles (constituent_id);
CREATE INDEX idx_admissions_applicant_profiles_active ON admissions_applicant_profiles (is_active);

-- Version: 1.14
-- Description: Create admissions application form templates
CREATE TABLE admissions_application_form_templates (
    form_template_id  UUID      NOT NULL,
    program_id        UUID      NOT NULL,
    academic_term_id  UUID      NOT NULL,
    application_type  TEXT      NOT NULL,
    name              TEXT      NOT NULL,
    description       TEXT      NULL,
    version           INT       NOT NULL,
    required_fields   JSONB     NOT NULL,
    checklist_items   JSONB     NOT NULL,
    is_active         BOOLEAN   NOT NULL,
    priority          INT       NOT NULL,
    date_created      TIMESTAMP NOT NULL,
    date_updated      TIMESTAMP NOT NULL,

    PRIMARY KEY (form_template_id),
    FOREIGN KEY (program_id) REFERENCES admissions_programs(program_id),
    FOREIGN KEY (academic_term_id) REFERENCES admissions_academic_terms(academic_term_id),
    CONSTRAINT admissions_form_templates_type CHECK (application_type IN ('FRESHMAN', 'TRANSFER', 'GRADUATE')),
    CONSTRAINT admissions_form_templates_name_not_empty CHECK (trim(name) <> ''),
    CONSTRAINT admissions_form_templates_version_positive CHECK (version >= 1),
    CONSTRAINT admissions_form_templates_required_fields_array CHECK (jsonb_typeof(required_fields) = 'array'),
    CONSTRAINT admissions_form_templates_checklist_items_array CHECK (jsonb_typeof(checklist_items) = 'array'),
    CONSTRAINT admissions_form_templates_priority_non_negative CHECK (priority >= 0)
);

CREATE INDEX idx_admissions_form_templates_program ON admissions_application_form_templates (program_id);
CREATE INDEX idx_admissions_form_templates_term ON admissions_application_form_templates (academic_term_id);
CREATE INDEX idx_admissions_form_templates_type ON admissions_application_form_templates (application_type);
CREATE INDEX idx_admissions_form_templates_active ON admissions_application_form_templates (is_active);
CREATE INDEX idx_admissions_form_templates_priority ON admissions_application_form_templates (priority);

-- Version: 1.15
-- Description: Create admissions checklist and document intake tables
CREATE TABLE admissions_checklist_items (
    checklist_item_id  UUID      NOT NULL,
    application_id     UUID      NOT NULL,
    item_key           TEXT      NOT NULL,
    document_name      TEXT      NOT NULL,
    description        TEXT      NULL,
    is_required        BOOLEAN   NOT NULL,
    status             TEXT      NOT NULL,
    display_order      INT       NOT NULL,
    date_created       TIMESTAMP NOT NULL,
    date_updated       TIMESTAMP NOT NULL,

    PRIMARY KEY (checklist_item_id),
    FOREIGN KEY (application_id) REFERENCES admissions_applications(application_id) ON DELETE CASCADE,
    CONSTRAINT admissions_checklist_items_item_key_not_empty CHECK (trim(item_key) <> ''),
    CONSTRAINT admissions_checklist_items_document_name_not_empty CHECK (trim(document_name) <> ''),
    CONSTRAINT admissions_checklist_items_order_non_negative CHECK (display_order >= 0),
    CONSTRAINT admissions_checklist_items_status CHECK (status IN ('UPLOADED', 'PENDING_REVIEW', 'ACCEPTED', 'REJECTED', 'WAIVED', 'EXPIRED', 'SYNCED_TO_SIS')),
    CONSTRAINT admissions_checklist_items_unique_key UNIQUE (application_id, item_key)
);

CREATE INDEX idx_admissions_checklist_items_application ON admissions_checklist_items (application_id);
CREATE INDEX idx_admissions_checklist_items_status ON admissions_checklist_items (status);
CREATE INDEX idx_admissions_checklist_items_required ON admissions_checklist_items (is_required);
CREATE INDEX idx_admissions_checklist_items_order ON admissions_checklist_items (display_order);

CREATE TABLE admissions_documents (
    document_id        UUID      NOT NULL,
    application_id     UUID      NOT NULL,
    checklist_item_id  UUID      NOT NULL,
    file_name          TEXT      NOT NULL,
    content_type       TEXT      NOT NULL,
    size_bytes         BIGINT    NOT NULL,
    storage_key        TEXT      NOT NULL,
    status             TEXT      NOT NULL,
    reviewer_id        UUID      NULL,
    reviewer_notes     TEXT      NULL,
    uploaded_by_id     UUID      NOT NULL,
    uploaded_at        TIMESTAMP NOT NULL,
    reviewed_at        TIMESTAMP NULL,
    date_created       TIMESTAMP NOT NULL,
    date_updated       TIMESTAMP NOT NULL,

    PRIMARY KEY (document_id),
    FOREIGN KEY (application_id) REFERENCES admissions_applications(application_id) ON DELETE CASCADE,
    FOREIGN KEY (checklist_item_id) REFERENCES admissions_checklist_items(checklist_item_id) ON DELETE CASCADE,
    FOREIGN KEY (uploaded_by_id) REFERENCES users(user_id),
    FOREIGN KEY (reviewer_id) REFERENCES users(user_id),
    CONSTRAINT admissions_documents_file_name_not_empty CHECK (trim(file_name) <> ''),
    CONSTRAINT admissions_documents_content_type_not_empty CHECK (trim(content_type) <> ''),
    CONSTRAINT admissions_documents_size_positive CHECK (size_bytes > 0),
    CONSTRAINT admissions_documents_storage_key_not_empty CHECK (trim(storage_key) <> ''),
    CONSTRAINT admissions_documents_status CHECK (status IN ('UPLOADED', 'PENDING_REVIEW', 'ACCEPTED', 'REJECTED', 'WAIVED', 'EXPIRED', 'SYNCED_TO_SIS'))
);

CREATE INDEX idx_admissions_documents_application ON admissions_documents (application_id);
CREATE INDEX idx_admissions_documents_checklist_item ON admissions_documents (checklist_item_id);
CREATE INDEX idx_admissions_documents_status ON admissions_documents (status);
CREATE INDEX idx_admissions_documents_uploaded_by ON admissions_documents (uploaded_by_id);
CREATE INDEX idx_admissions_documents_reviewer ON admissions_documents (reviewer_id);
CREATE INDEX idx_admissions_documents_uploaded_at ON admissions_documents (uploaded_at);

-- Version: 1.16
-- Description: Create admissions lead scoring rules and scores
CREATE TABLE admissions_lead_score_rules (
    lead_score_rule_id  UUID      NOT NULL,
    name                TEXT      NOT NULL,
    description         TEXT      NULL,
    criteria            JSONB     NOT NULL,
    points              INT       NOT NULL,
    is_active           BOOLEAN   NOT NULL,
    priority            INT       NOT NULL,
    date_created        TIMESTAMP NOT NULL,
    date_updated        TIMESTAMP NOT NULL,

    PRIMARY KEY (lead_score_rule_id),
    CONSTRAINT admissions_lead_score_rules_name_not_empty CHECK (trim(name) <> ''),
    CONSTRAINT admissions_lead_score_rules_criteria_not_empty CHECK (jsonb_typeof(criteria) = 'array' AND jsonb_array_length(criteria) > 0),
    CONSTRAINT admissions_lead_score_rules_points CHECK (points >= 0),
    CONSTRAINT admissions_lead_score_rules_priority CHECK (priority >= 0)
);

CREATE INDEX idx_admissions_lead_score_rules_active ON admissions_lead_score_rules (is_active);
CREATE INDEX idx_admissions_lead_score_rules_priority ON admissions_lead_score_rules (priority);

CREATE TABLE admissions_lead_scores (
    lead_score_id    UUID      NOT NULL,
    constituent_id   UUID      NOT NULL,
    total_score      INT       NOT NULL,
    score_band       TEXT      NOT NULL,
    breakdown        JSONB     NOT NULL,
    recalculated_at  TIMESTAMP NOT NULL,
    date_created     TIMESTAMP NOT NULL,
    date_updated     TIMESTAMP NOT NULL,

    PRIMARY KEY (lead_score_id),
    FOREIGN KEY (constituent_id) REFERENCES admissions_constituents(constituent_id) ON DELETE CASCADE,
    CONSTRAINT admissions_lead_scores_unique_constituent UNIQUE (constituent_id),
    CONSTRAINT admissions_lead_scores_total_score CHECK (total_score >= 0),
    CONSTRAINT admissions_lead_scores_band CHECK (score_band IN ('COLD', 'WARM', 'HOT', 'READY_TO_APPLY')),
    CONSTRAINT admissions_lead_scores_breakdown_array CHECK (jsonb_typeof(breakdown) = 'array')
);

CREATE INDEX idx_admissions_lead_scores_constituent ON admissions_lead_scores (constituent_id);
CREATE INDEX idx_admissions_lead_scores_band ON admissions_lead_scores (score_band);
CREATE INDEX idx_admissions_lead_scores_total_score ON admissions_lead_scores (total_score DESC);

-- Version: 1.17
-- Description: Create admissions custom field registry and values
CREATE TABLE admissions_custom_field_definitions (
    custom_field_definition_id UUID      NOT NULL,
    owner                      TEXT      NOT NULL,
    field_key                  TEXT      NOT NULL,
    label                      TEXT      NOT NULL,
    description                TEXT      NULL,
    data_type                  TEXT      NOT NULL,
    is_required                BOOLEAN   NOT NULL,
    options                    TEXT[]    NOT NULL,
    validation                 TEXT      NULL,
    is_searchable              BOOLEAN   NOT NULL,
    is_reportable              BOOLEAN   NOT NULL,
    is_importable              BOOLEAN   NOT NULL,
    is_exportable              BOOLEAN   NOT NULL,
    display_order              INT       NOT NULL,
    is_active                  BOOLEAN   NOT NULL,
    date_created               TIMESTAMP NOT NULL,
    date_updated               TIMESTAMP NOT NULL,

    PRIMARY KEY (custom_field_definition_id),
    CONSTRAINT admissions_custom_field_definitions_owner CHECK (owner IN ('CONSTITUENT', 'APPLICATION')),
    CONSTRAINT admissions_custom_field_definitions_key_not_empty CHECK (trim(field_key) <> ''),
    CONSTRAINT admissions_custom_field_definitions_label_not_empty CHECK (trim(label) <> ''),
    CONSTRAINT admissions_custom_field_definitions_data_type CHECK (data_type IN ('TEXT', 'TEXTAREA', 'NUMBER', 'DATE', 'SELECT', 'BOOLEAN')),
    CONSTRAINT admissions_custom_field_definitions_select_options CHECK (data_type <> 'SELECT' OR cardinality(options) > 0),
    CONSTRAINT admissions_custom_field_definitions_order_non_negative CHECK (display_order >= 0),
    CONSTRAINT admissions_custom_field_definitions_unique_key UNIQUE (owner, field_key)
);

CREATE INDEX idx_admissions_custom_field_definitions_owner ON admissions_custom_field_definitions (owner);
CREATE INDEX idx_admissions_custom_field_definitions_active ON admissions_custom_field_definitions (is_active);
CREATE INDEX idx_admissions_custom_field_definitions_display_order ON admissions_custom_field_definitions (display_order);

CREATE TABLE admissions_custom_field_values (
    custom_field_value_id      UUID      NOT NULL,
    custom_field_definition_id UUID      NOT NULL,
    owner                      TEXT      NOT NULL,
    owner_id                   UUID      NOT NULL,
    value                      TEXT      NOT NULL,
    date_created               TIMESTAMP NOT NULL,
    date_updated               TIMESTAMP NOT NULL,

    PRIMARY KEY (custom_field_value_id),
    FOREIGN KEY (custom_field_definition_id) REFERENCES admissions_custom_field_definitions(custom_field_definition_id) ON DELETE CASCADE,
    CONSTRAINT admissions_custom_field_values_owner CHECK (owner IN ('CONSTITUENT', 'APPLICATION')),
    CONSTRAINT admissions_custom_field_values_value_not_empty CHECK (trim(value) <> ''),
    CONSTRAINT admissions_custom_field_values_unique_owner UNIQUE (custom_field_definition_id, owner, owner_id)
);

CREATE INDEX idx_admissions_custom_field_values_definition ON admissions_custom_field_values (custom_field_definition_id);
CREATE INDEX idx_admissions_custom_field_values_owner ON admissions_custom_field_values (owner, owner_id);

-- Version: 1.18
-- Description: Create admissions import batch and invalid row report tables
CREATE TABLE admissions_import_batches (
    import_batch_id    UUID      NOT NULL,
    source             TEXT      NOT NULL,
    file_type          TEXT      NOT NULL,
    target             TEXT      NOT NULL,
    status             TEXT      NOT NULL,
    file_name          TEXT      NOT NULL,
    storage_key        TEXT      NULL,
    uploaded_by_id     UUID      NOT NULL,
    total_rows         INT       NOT NULL,
    valid_rows         INT       NOT NULL,
    invalid_rows       INT       NOT NULL,
    duplicate_rows     INT       NOT NULL,
    field_mapping      JSONB     NOT NULL,
    invalid_report_key TEXT      NULL,
    validation_summary TEXT      NULL,
    committed_at       TIMESTAMP NULL,
    date_created       TIMESTAMP NOT NULL,
    date_updated       TIMESTAMP NOT NULL,

    PRIMARY KEY (import_batch_id),
    FOREIGN KEY (uploaded_by_id) REFERENCES users(user_id),
    CONSTRAINT admissions_import_batches_source CHECK (source IN ('MANUAL_UPLOAD', 'SIS_EXPORT')),
    CONSTRAINT admissions_import_batches_file_type CHECK (file_type IN ('CSV', 'XLSX')),
    CONSTRAINT admissions_import_batches_target CHECK (target IN ('CONSTITUENTS', 'APPLICATIONS')),
    CONSTRAINT admissions_import_batches_status CHECK (status IN ('PREVIEWED', 'VALIDATION_FAILED', 'QUEUED', 'PROCESSING', 'COMPLETED', 'FAILED')),
    CONSTRAINT admissions_import_batches_file_name_not_empty CHECK (trim(file_name) <> ''),
    CONSTRAINT admissions_import_batches_rows_non_negative CHECK (total_rows >= 0 AND valid_rows >= 0 AND invalid_rows >= 0 AND duplicate_rows >= 0),
    CONSTRAINT admissions_import_batches_rows_consistent CHECK (valid_rows + invalid_rows <= total_rows),
    CONSTRAINT admissions_import_batches_field_mapping_object CHECK (jsonb_typeof(field_mapping) = 'object')
);

CREATE INDEX idx_admissions_import_batches_target ON admissions_import_batches (target);
CREATE INDEX idx_admissions_import_batches_status ON admissions_import_batches (status);
CREATE INDEX idx_admissions_import_batches_uploaded_by ON admissions_import_batches (uploaded_by_id);
CREATE INDEX idx_admissions_import_batches_created ON admissions_import_batches (date_created);

CREATE TABLE admissions_import_invalid_rows (
    import_invalid_row_id UUID      NOT NULL,
    import_batch_id       UUID      NOT NULL,
    row_number            INT       NOT NULL,
    field_name            TEXT      NULL,
    raw_data              JSONB     NOT NULL,
    error_code            TEXT      NOT NULL,
    error_detail          TEXT      NOT NULL,
    date_created          TIMESTAMP NOT NULL,

    PRIMARY KEY (import_invalid_row_id),
    FOREIGN KEY (import_batch_id) REFERENCES admissions_import_batches(import_batch_id) ON DELETE CASCADE,
    CONSTRAINT admissions_import_invalid_rows_row_number CHECK (row_number > 0),
    CONSTRAINT admissions_import_invalid_rows_raw_data_object CHECK (jsonb_typeof(raw_data) = 'object'),
    CONSTRAINT admissions_import_invalid_rows_error_code_not_empty CHECK (trim(error_code) <> ''),
    CONSTRAINT admissions_import_invalid_rows_error_detail_not_empty CHECK (trim(error_detail) <> '')
);

CREATE INDEX idx_admissions_import_invalid_rows_batch ON admissions_import_invalid_rows (import_batch_id);
CREATE INDEX idx_admissions_import_invalid_rows_row_number ON admissions_import_invalid_rows (row_number);

-- Version: 1.19
-- Description: Create Kenya counties reference table
CREATE TABLE counties (
    code         TEXT      NOT NULL,
    name         TEXT      NOT NULL,
    date_created TIMESTAMP NOT NULL,
    date_updated TIMESTAMP NOT NULL,

    PRIMARY KEY (code),
    CONSTRAINT counties_code_not_empty CHECK (trim(code) <> ''),
    CONSTRAINT counties_name_not_empty CHECK (trim(name) <> '')
);

CREATE INDEX idx_counties_name ON counties (name);

-- Version: 1.20
-- Description: Create Kenya sub-counties reference table
CREATE TABLE sub_counties (
    code         TEXT      NOT NULL,
    county_code  TEXT      NOT NULL,
    name         TEXT      NOT NULL,
    date_created TIMESTAMP NOT NULL,
    date_updated TIMESTAMP NOT NULL,

    PRIMARY KEY (code),
    FOREIGN KEY (county_code) REFERENCES counties(code),
    CONSTRAINT sub_counties_code_not_empty CHECK (trim(code) <> ''),
    CONSTRAINT sub_counties_name_not_empty CHECK (trim(name) <> '')
);

CREATE INDEX idx_sub_counties_county_code ON sub_counties (county_code);
CREATE INDEX idx_sub_counties_name ON sub_counties (name);

-- Version: 1.21
-- Description: Add Kenya constituent identity identifiers and manual backfill review table
ALTER TABLE admissions_constituents
    ADD COLUMN national_id TEXT NULL,
    ADD COLUMN national_id_verified_at TIMESTAMP NULL,
    ADD COLUMN national_id_verified_by_adapter TEXT NULL,
    ADD COLUMN upi TEXT NULL,
    ADD COLUMN upi_verified_at TIMESTAMP NULL,
    ADD COLUMN upi_verified_by_adapter TEXT NULL,
    ADD COLUMN kcse_index_number TEXT NULL,
    ADD COLUMN kcse_index_verified_at TIMESTAMP NULL,
    ADD COLUMN kcse_index_verified_by_adapter TEXT NULL;

CREATE UNIQUE INDEX idx_admissions_constituents_national_id ON admissions_constituents (national_id) WHERE national_id IS NOT NULL;
CREATE UNIQUE INDEX idx_admissions_constituents_upi ON admissions_constituents (upi) WHERE upi IS NOT NULL;
CREATE UNIQUE INDEX idx_admissions_constituents_kcse_index_number ON admissions_constituents (kcse_index_number) WHERE kcse_index_number IS NOT NULL;

CREATE TABLE admissions_identity_backfill_reviews (
    identity_backfill_review_id UUID      NOT NULL,
    constituent_id              UUID      NOT NULL,
    external_sis_id             TEXT      NOT NULL,
    review_reason               TEXT      NOT NULL,
    status                      TEXT      NOT NULL,
    date_created                TIMESTAMP NOT NULL,
    date_updated                TIMESTAMP NOT NULL,

    PRIMARY KEY (identity_backfill_review_id),
    FOREIGN KEY (constituent_id) REFERENCES admissions_constituents(constituent_id) ON DELETE CASCADE,
    CONSTRAINT admissions_identity_backfill_reviews_reason_not_empty CHECK (trim(review_reason) <> ''),
    CONSTRAINT admissions_identity_backfill_reviews_status CHECK (status IN ('PENDING', 'RESOLVED')),
    CONSTRAINT admissions_identity_backfill_reviews_unique_constituent UNIQUE (constituent_id, external_sis_id)
);

INSERT INTO admissions_identity_backfill_reviews
    (identity_backfill_review_id, constituent_id, external_sis_id, review_reason, status, date_created, date_updated)
SELECT
    gen_random_uuid(),
    constituent_id,
    external_sis_id,
    'external_sis_id requires manual classification before migration to national_id, upi, or kcse_index_number',
    'PENDING',
    NOW(),
    NOW()
FROM
    admissions_constituents
WHERE
    external_sis_id IS NOT NULL
ON CONFLICT DO NOTHING;

CREATE INDEX idx_admissions_identity_backfill_reviews_constituent ON admissions_identity_backfill_reviews (constituent_id);
CREATE INDEX idx_admissions_identity_backfill_reviews_status ON admissions_identity_backfill_reviews (status);

-- Version: 1.22
-- Description: Localize admissions applications for Kenya
ALTER TABLE admissions_applications
    DROP CONSTRAINT admissions_applications_type,
    ADD COLUMN kuccps_placement JSONB NULL,
    ADD COLUMN kcse_result JSONB NULL,
    ADD CONSTRAINT admissions_applications_type CHECK (application_type IN ('KUCCPS_PLACEMENT', 'SELF_SPONSORED_UNDERGRAD', 'DIPLOMA', 'MASTERS', 'PHD', 'TVET', 'BRIDGING', 'CERTIFICATE'));

ALTER TABLE admissions_application_form_templates
    DROP CONSTRAINT admissions_form_templates_type,
    ADD CONSTRAINT admissions_form_templates_type CHECK (application_type IN ('KUCCPS_PLACEMENT', 'SELF_SPONSORED_UNDERGRAD', 'DIPLOMA', 'MASTERS', 'PHD', 'TVET', 'BRIDGING', 'CERTIFICATE'));
