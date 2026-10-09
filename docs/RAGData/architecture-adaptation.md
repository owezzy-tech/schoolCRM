# EduPlan features in SchoolCRM

The user selected the existing SchoolCRM stack, LangGraph orchestration and LlamaIndex retrieval.
The EduPlan blueprint is a feature and architecture reference, not a requirement to
migrate to React, LangChainGo, a new HTTP router or a separate repository.

| Concern | SchoolCRM direction | Blueprint reference |
| --- | --- | --- |
| Educator and admin workflows | Extend the existing Angular application | Page 1, frontend and core platform vision |
| Operational APIs and lesson persistence | Extend existing Go application/domain/storage layers and JSON:API response contract | Pages 1–2, API and schema |
| Curriculum ingestion, retrieval and AI generation | Use LlamaIndex within Python RAG adapters | Pages 1–2, discovery and generation |
| Agent workflows | Use LangGraph in Python with durable checkpoints and approvals | User decision |
| Lesson history | Store complete immutable version snapshots with a change summary and an atomic current-version pointer | Page 2, version history |
| Publishing | Teacher drafts; HOD reviews; dean approves every version; teacher publishes | Confirmed user decisions |
| School and department access | Reuse existing identity; enforce scoped authorisation before retrieval and state changes | Page 3, security |
| PDF export | Export an authorised, selected lesson version | Page 1, PDF export |
| Workflow and deployment tooling | Evaluate existing infrastructure against demonstrated durability and performance needs | Pages 1 and 3 |

## LlamaIndex boundary

The current Python query path embeds a question, calls `IVectorStore.search`, then calls
`ILLMProvider.answer`. Assembly in `api/services/RAG/infrastructure/dependencies.py`
currently wires no-op embeddings, an in-memory vector store and an echo LLM. These
are scaffold adapters, not working semantic retrieval or AI generation.

LlamaIndex belongs in Python adapters and infrastructure assembly. Preserve framework-free
domain types and application use cases. The existing embedding, vector-store and LLM
ports are the starting seams. Extend the responsible retrieval boundary where school,
subject, curriculum revision and source access filters are required; the present vector
store search signature has no such filters and is insufficient for private curriculum data.
Apply these filters during retrieval, before any model sees the content.

Use LlamaIndex for ingestion/indexing, retrieval and grounded generation. Qualification
calculations belong in deterministic domain rules with reviewed evidence. Retrieved
prose or model inference must not set a student's numeric points or invent entry rules.
Blueprint proposals must not enter the student admissions evidence collection. Keep
public admissions data, private school curriculum and product architecture distinct at
the retrieval boundary, even though this portable graph contains all three source files.

This decision selects LlamaIndex for subsequent implementation. No LlamaIndex package,
live vector database, model provider or new RAG endpoint is installed by this data build.
The reproducible graph itself remains a standard-library Python build.

## Repairs needed before implementing the blueprint schema

The sample tables have no school/department scope and the Casbin sample expresses
global roles only. The feature's school isolation must be enforced in both persistence
queries and server authorisation, including retrieval, generation, plan history and export.

The version table does not enforce append-only writes or unique versions per plan.
Its cascade deletion can remove history. Define retention/deletion rules, protect published
history, and make creation of a version plus the current-version pointer transactional.
The blueprint calls this delta storage but supplies full `snapshot_json`; use snapshots
until an explicit requirement establishes a need for diffs.

The free-form status column does not enforce workflow transitions. The teacher submits,
the Head of Department reviews, the dean approves every version including the initial
plan, and the teacher publishes that exact approved version.
Published content should reference an immutable version, not mutable draft content.

Vector dimensions must follow the selected embedding model. Do not assume the example's
1536 dimensions fit every provider. Record model identity and curriculum revision with
each index and plan's retrieval evidence. Validate generated lesson content against the
agreed schema before storing it as a draft.

The blueprint's uptime and latency figures are targets requiring measurements, not
properties guaranteed by SQLC, Go or a queue package. Temporal, Asynq, Casbin, SQLC and
Maroto remain candidate tools; select them only when existing repository components do
not meet the required behaviour.

The confirmed curriculum scope is Kenya CBC/CBE and 8-4-4 from Playgroup through
secondary, plus Cambridge Early Years through Year 9. Keep framework, stage, subject,
revision and access filters explicit. Higher Cambridge years are future expansion.
Google Calendar is the initial provider; Nextcloud is a separately assessed alternative.

Beads epic `schoolCRM-o58` owns delivery. Actual curriculum resources, lesson output
details, embedding/model provider and retention rules remain open. The current proposed
design and exact approval lifecycle are in [the diagram specification](../diagrams/school-ai-v1-design.md).
