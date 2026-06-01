# Admissions CRM Accessibility Release Gates

This gate makes Section 508 and WCAG AA checks testable for admissions portal and staff review releases.

## Scope

Covered flows:

- Applicant portal entry, application draft, status, document upload, and event registration.
- Staff application review workspace, document verification, decision entry, and audit timeline.
- Admissions status badges, progress indicators, tables, forms, and async messages.

Out of scope for this gate:

- RAG user interfaces.
- Native mobile accessibility.
- Full third-party payment provider accessibility certification.

## WCAG AA basics checklist

Forms:

- [ ] Every input has a visible label.
- [ ] Help text and error text are associated with the input via `aria-describedby` where the component allows it.
- [ ] Errors use text, not color alone.
- [ ] Errors and success messages use `role="alert"` or `aria-live="polite"`.
- [ ] Required fields are identified in text or label copy.
- [ ] Final submission and decision actions explain disabled or blocked states.

Keyboard and focus:

- [ ] All actions are reachable with Tab and Shift+Tab.
- [ ] Focus order follows visual reading order.
- [ ] Visible focus rings are present on links, buttons, form controls, tabs, and upload controls.
- [ ] No keyboard trap exists in portal forms, review tabs, upload controls, or future dialogs.
- [ ] Disabled actions either leave the tab order or explain why they are unavailable.

Status, tables, and visual indicators:

- [ ] Status badges include readable text labels and are not color-only.
- [ ] Progress bars include an accessible label or adjacent text summary.
- [ ] Decorative icons are hidden from assistive technology.
- [ ] Meaningful icon-only controls include `aria-label`.
- [ ] Data tables use `scope="col"` for column headers and `scope="row"` where a row heading exists.
- [ ] Empty states describe the action needed to recover.

Contrast and layout:

- [ ] Body text and status labels meet 4.5:1 contrast in light and dark themes.
- [ ] Interactive component boundaries meet 3:1 contrast.
- [ ] Focus indicators meet 3:1 contrast against adjacent colors.
- [ ] Text remains readable at 200% zoom without horizontal scrolling except data tables.

## Keyboard acceptance scripts

### Applicant portal entry

1. Open `/portal`.
2. Tab to `Start application`, then `Check status`, then `Application email`, then `Get access`.
3. Submit an empty or invalid email state.
4. Confirm the error message is announced and visually rendered as text.
5. Shift+Tab back through the same controls without losing context.

### Applicant application draft

1. Open `/portal/apply`.
2. Tab through first name, last name, email, phone, password, confirm password, Back, and Continue.
3. Confirm the progress indicator has text describing current progress.
4. Confirm the fee status badge reads as text, not color only.

### Staff review workspace

1. Open `/applications/review`.
2. Tab through Applications back link, disabled submit-preview action, assigned application buttons, and the main review tabs.
3. Confirm locked records announce why they cannot be selected.
4. In Decision Entry, tab through decision, transition reason, committee note, Save draft, and disabled Submit decision.
5. Confirm the document table headers are announced as column headers.

### Document verification

1. Open an application document checklist.
2. Tab to the load form, application ID, and Load checklist button.
3. Tab through review section tabs and document action buttons.
4. Confirm status, required, upload date, and reviewer note values are not color-only.

## Release evidence

Each release must capture:

- Date and reviewer.
- Browser and assistive technology used.
- Manual keyboard script result for portal and staff review.
- Automated scan tool result, if available.
- Contrast check result for light and dark mode.
- Exceptions, linked issue IDs, and planned remediation.

## Evidence template

```text
Release:
Reviewer:
Date:
Browser:
Assistive technology:

Portal keyboard script: PASS / FAIL
Staff review keyboard script: PASS / FAIL
Document verification keyboard script: PASS / FAIL
Automated scan: PASS / FAIL / NOT RUN
Contrast check: PASS / FAIL

Exceptions:
- [issue-id] Description, severity, owner, due date
```
