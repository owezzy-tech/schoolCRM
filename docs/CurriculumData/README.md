# Curriculum source acquisition

Acquired on 3 October 2026 for `schoolCRM-o58.2.1`.

The user supplied the Cambridge IGCSE curriculum page and the CBC Education Kenya article. The article links to KICD; the Kenyan inventory follows the official KICD regular-learner pages rather than treating the article as curriculum evidence.

## Files and results

| Collection | Result | Inventory |
| --- | --- | --- |
| Cambridge IGCSE | 366 public syllabus and syllabus-update PDFs from 104 subject pages, 199,257,606 bytes, 10,061 pages | `cambridge-igcse.jsonl` |
| Kenya CBC/CBE | 220 published document links across pre-primary, lower primary and Grades 4–12; no downloadable PDFs | `kenya-cbc.jsonl` |

Raw downloads are stored locally at `/Users/owen_adirah/Documents/SchoolCRM/curriculum/2026-10-03`. Each Cambridge inventory record gives its relative filename, original and resolved URL, subject page, displayed revision label, SHA256 checksum, byte count and page count. `acquisition.json` records the source pages and validation totals. The binaries are not committed to Git.

All 366 downloaded files passed SHA256 verification and parsing with `pypdf`; every file contains pages. Parsing does not prove extraction quality or curriculum suitability.

The KICD pages embed Google Drive previews. Of 220 download attempts, 218 returned a non-PDF download-error page and two failed at the HTTP request. An inspected error page explicitly says the owner permits downloads only for the owner and editors. No HTML error responses were saved as PDFs. Follow-up `schoolCRM-o58.2.2` tracks authorised downloadable originals.

The Kenyan collection covers regular learners. Special-needs and teacher-education collections are outside this acquisition. Stage and subject labels come from publication pages and still require verification against the originals.

## Review before retrieval

Downloaded documents remain `unreviewed`, with revisions `unverified`. Cambridge collections include past, current and future examination years and update notices. A syllabus year is an examination year; acquisition date does not make a revision current. Before indexing, bind each document to its verified subject, examination years, revision, school access and approval state.

Cambridge IGCSE is upper secondary. Downloading these resources does not silently replace the previously confirmed Early Years through Year 9 runtime scope. Keep IGCSE distinct until its application scope is confirmed. These pages supply neither legacy Kenyan 8-4-4 resources nor complete Cambridge Early Years/Primary/Lower Secondary resources.

## Provider decisions

The user selected BGE embeddings through local Ollama in development and Cloudflare Workers in production, with DeepSeek for generation. BGE-M3 is the selected multilingual variant: `bge-m3` in Ollama and `@cf/baai/bge-m3` in Cloudflare Workers AI.

Model availability was checked against the [Ollama model page](https://ollama.com/library/bge-m3) and [Cloudflare model page](https://developers.cloudflare.com/workers-ai/models/bge-m3/). The local Ollama endpoint was unavailable during acquisition. No models were installed and no production account, credentials or deployment was changed.

Pin the actual local model digest and verify the embedding dimension at setup. Record provider, model, preprocessing and index revision. Use separate provider indexes and reindex on a provider change until measured compatibility permits reuse; matching vector dimensions alone is insufficient. The DeepSeek model ID and serving endpoint remain to be selected. This decision does not establish a policy for sending private student data to external providers.
