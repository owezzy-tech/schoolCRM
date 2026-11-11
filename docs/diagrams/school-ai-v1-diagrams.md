# School assistant v1 diagrams

Interactive standalone Archify HTML diagrams for the proposed SchoolCRM v1.
Open the HTML files in a browser. The viewer supports dark/light themes and export.
These diagrams describe the selected architecture and future behaviour; they do not
claim LangGraph, LlamaIndex or the new business workflows are already deployed.

| Diagram | HTML | Source | Final validation receipt |
| --- | --- | --- | --- |
| Runtime architecture | [Open](../../.archify/architecture-school-ai-v1-20261001-195546/school-ai-v1.html) | [JSON](../../.archify/architecture-school-ai-v1-20261001-195546/candidate.json) | [Receipt](../../.archify/architecture-school-ai-v1-20261001-195546/review-3/school-ai-v1.finalize-summary.json) |
| Curriculum and admissions evidence | [Open](../../.archify/dataflow-school-curriculum-20261001-195546/school-curriculum.html) | [JSON](../../.archify/dataflow-school-curriculum-20261001-195546/candidate.json) | [Receipt](../../.archify/dataflow-school-curriculum-20261001-195546/school-curriculum.finalize-summary.json) |
| Role-aware agent actions | [Open](../../.archify/workflow-school-agent-actions-20261001-195546/school-agent-actions.html) | [JSON](../../.archify/workflow-school-agent-actions-20261001-195546/candidate.json) | [Receipt](../../.archify/workflow-school-agent-actions-20261001-195546/school-agent-actions.finalize-summary.json) |
| Lesson version lifecycle | [Open](../../.archify/lifecycle-school-lesson-versions-20261001-195546/school-lesson-versions.html) | [JSON](../../.archify/lifecycle-school-lesson-versions-20261001-195546/candidate.json) | [Receipt](../../.archify/lifecycle-school-lesson-versions-20261001-195546/school-lesson-versions.finalize-summary.json) |

The diagrams pin repository evidence to commit
`bea454bdb4f420cac2074c91adbd89a6c0e9cc74` and the
[confirmed design specification](school-ai-v1-design.md).

All four passed `validate`, `deliver`, strict artifact `check`, and real Chrome
`browser-check` under the showcase profile, with no gate diagnostics. Candidate and
HTML SHA-256 values were independently matched against the final receipts.
Manual perceptual review was not performed. Route-readability recommendations remain
advisory in the receipts. The required architecture placement review was attempted;
it introduced another crossing, so the original accepted placement was restored and
finalized again. `review-2` retains that historical attempt; `review-3` is authoritative.

Confirmed scope:

- Kenya CBC/CBE and 8-4-4, Playgroup through secondary; Cambridge Early Years through
  Year 9. Keep framework, stage, subject and revision distinct. Curriculum files remain needed.
- Every lesson version: teacher draft, Head of Department review, dean approval,
  then teacher publication. Published versions remain immutable while revisions are drafts.
- Existing Go/Angular/Python stack, LangGraph orchestration and LlamaIndex retrieval.
- Google Calendar initially. [Google's current quota/pricing guidance](https://developers.google.com/workspace/calendar/api/guides/quota)
  says standard API use has no additional cost but quota excess charging is planned later
  in 2026. Workspace/account costs are separate.
- [Nextcloud Calendar](https://docs.nextcloud.com/server/latest/user_manual/en/groupware/calendar.html)
  is an open-source alternative requiring hosting and another integration. Its documented
  free/busy visibility is limited to users on the same instance, so universal feature parity
  with Google Calendar is not assumed.

Beads epic: `schoolCRM-o58`. Diagram task: `schoolCRM-rii`.
