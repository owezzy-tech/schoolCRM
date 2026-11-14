# Agent Instructions

This project uses **bd** (beads) for issue tracking. Run `bd onboard` to get started.

## Quick Reference

```bash
bd ready              # Find available work
bd show <id>          # View issue details
bd update <id> --claim  # Claim work atomically
bd close <id>         # Complete work
bd dolt push          # Push beads data to remote
```

## Non-Interactive Shell Commands

**ALWAYS use non-interactive flags** with file operations to avoid hanging on confirmation prompts.

Shell commands like `cp`, `mv`, and `rm` may be aliased to include `-i` (interactive) mode on some systems, causing the agent to hang indefinitely waiting for y/n input.

**Use these forms instead:**
```bash
# Force overwrite without prompting
cp -f source dest           # NOT: cp source dest
mv -f source dest           # NOT: mv source dest
rm -f file                  # NOT: rm file

# For recursive operations
rm -rf directory            # NOT: rm -r directory
cp -rf source dest          # NOT: cp -r source dest
```

**Other commands that may prompt:**
- `scp` - use `-o BatchMode=yes` for non-interactive
- `ssh` - use `-o BatchMode=yes` to fail instead of prompting
- `apt-get` - use `-y` flag
- `brew` - use `HOMEBREW_NO_AUTO_UPDATE=1` env var

## API Response Contract

- All backend HTTP API handlers registered through `foundation/web` must return JSON:API v1.1 documents.
- Success responses use `Content-Type: application/vnd.api+json` and include top-level `jsonapi: { version: "1.1" }` plus `data`.
- Error responses use top-level `errors`; surface client-safe messages in `errors[0].detail` so frontends can display API failure messages.
- Do not return ad hoc `{ "error": ... }` or bare domain JSON from normal API handlers. Add new response behavior through the central `foundation/web.Respond` pipeline.
- Frontends consuming backend APIs must unwrap JSON:API `data.attributes` for single resources, `data` plus `meta` for collections, and `errors[0].detail` for failed requests.

## Architecture and delivery

- Before implementing School AI features, read `docs/school-ai-ard.md` for ownership, contracts, invariants, acceptance conditions and unresolved decisions. Update the ARD when an authorised decision changes those requirements.
- Use Git Flow: `main` holds releases; `develop` integrates work. Create one `feature/<bead-id>-<slug>` branch from current `origin/develop` for each bead, including documentation and ordinary bug fixes. Use `release/<version>` from `develop` for releases and `hotfix/<bead-id>-<slug>` from `main` for urgent production fixes. Release and hotfix changes return to both `main` and `develop` through PRs. Preserve unrelated local changes, using an isolated worktree when needed.
- For every bead, run the affected quality gates and the `thermo-nuclear-code-quality-review` skill against the full PR diff. The user's `/thermal-nuclear-review` means this skill. Resolve blocking findings before delivery; if the skill is unavailable, report the missing review and keep the bead open.
- Commit and push each validated bead, then open its own PR against `develop`, or `main` for a hotfix. Include the bead ID, acceptance evidence, review verdict and material risks in the PR. Record the PR URL in Beads and close the bead only after validation, review and PR creation succeed. PR creation is required; merging remains a separate action.
- Finish one bead's PR before starting the next. A parent epic stays open until its child features satisfy their acceptance conditions. Documentation delivery does not close implementation beads.

## Agent Checkpoints

Stop and ask the user before you do any of these. Do not decide alone, even when the change looks small or obviously right.

1. **Changing an existing test's expected result.** This includes status codes, error codes, messages and golden values. Adding a new test is fine. Do not edit what an existing test expects unless the user approved that behaviour change. List each changed expectation in the handoff.
2. **Changing Git hooks, CI configuration or the scope of any check.** This covers `.husky/`, `scripts/git-hooks/`, `.circleci/`, linter configuration, and any change that skips, narrows or weakens a check. Never use `--no-verify`. If a hook blocks you, stop and report it.
3. **Starting or stopping local services.** This covers Docker containers, databases and dev servers. Ask first, unless the task brief already approves that service.
4. **A test and the code disagree.** Do not change either one to make them agree. Report which one you think is wrong and why, then wait for a decision.

## Handoff Format

End every task, and every stop at a checkpoint, with these four sections:

- **Changed:** each file or behaviour that changed, and why. List every changed test expectation separately.
- **Checks run:** each check with its command, result and environment (for example: `CGO_ENABLED=0 go test ./...`, pass, local macOS with Docker `servicetest`).
- **Not run or unfinished:** checks you skipped, work left open, and the reason for each.
- **Questions for the reviewer:** decisions you made that need confirmation, and open risks.

## Frontend Design Rules

These apply to `api/frontends/web-admin`. The tokens in `src/styles/design-tokens.css` and the rules in `design-spec.md` are the source of truth.

1. **Use tokens, never raw values.** Read colours, spacing, radii, shadows and type from the `--elimu-*` custom properties or the Fuse theme classes. Do not write a hex colour, a new shadow or a one-off radius.
2. **Start from a story.** Before you build or change a component, find the story in `src/stories/` that shows the same state and follow its markup and wording.
3. **Add a story for every new state.** Use real product wording in fixture content, never placeholder text. Add a `play` test when the component has an interaction.
4. **Check both themes.** Look at the result in light and dark in Storybook, and run the stories' `play` tests. Report both in the handoff.
5. **No token fits? Stop and ask.** A missing token is a design decision. Do not invent a value that looks close.

## Code Review Rules

### Authorization

- Check every school and lesson command against the caller's active membership on the server. Hiding a button in the frontend or reading `/v1/me/school-memberships` grants nothing.
- Only the author may read an unsubmitted lesson draft. Other members see the newest submitted version.

### API contract

- Flag any handler that bypasses `foundation/web.Respond` or returns a body that is not JSON:API v1.1, as defined in "API Response Contract" above.

### Migrations

- Add schema changes as a new `-- Version:` block at the end of `business/sdk/migrate/sql/migrate.sql`. Flag any edit to a block that already exists on `develop`.

<!-- BEGIN BEADS INTEGRATION v:1 profile:minimal hash:ca08a54f -->
## Beads Issue Tracker

This project uses **bd (beads)** for issue tracking. Run `bd prime` to see full workflow context and commands.

### Quick Reference

```bash
bd ready              # Find available work
bd show <id>          # View issue details
bd update <id> --claim  # Claim work
bd close <id>         # Complete work
```

### Rules

- Use `bd` for ALL task tracking — do NOT use TodoWrite, TaskCreate, or markdown TODO lists
- Run `bd prime` for detailed command reference and session close protocol
- Use `bd remember` for persistent knowledge — do NOT use MEMORY.md files
- Beads is initialized **only at the repository root**. Always run `bd` commands from `/Users/owen_adirah/GolandProjects/schoolCRM`; do NOT run `bd init` or create nested `.beads/` directories in subprojects such as `api/frontends/web-admin`.
- Create a bead before feature implementation, claim it with `bd update <id> --claim`, and close it only after validation passes.
- Never use `bd edit`; it opens an interactive editor. Use inline flags such as `--description`, `--acceptance`, `--notes`, and `--design`.
- Use numeric priorities only: `0` critical, `1` high, `2` medium, `3` low, `4` backlog.

## Session Completion

**When ending a work session**, you MUST complete ALL steps below. Work is NOT complete until `git push` succeeds.

**MANDATORY WORKFLOW:**

1. **File issues for remaining work** - Create issues for anything that needs follow-up
2. **Run quality gates** (if code changed) - Tests, linters, builds
3. **Update issue status** - Close finished work, update in-progress items
4. **PUSH TO REMOTE** - This is MANDATORY:
   ```bash
   git pull --rebase
   bd dolt push
   git push
   git status  # MUST show "up to date with origin"
   ```
5. **Clean up** - Clear stashes, prune remote branches
6. **Verify** - All changes committed AND pushed
7. **Hand off** - Provide context for next session

**CRITICAL RULES:**
- Work is NOT complete until `git push` succeeds
- NEVER stop before pushing - that leaves work stranded locally
- NEVER say "ready to push when you are" - YOU must push
- If push fails, resolve and retry until it succeeds
<!-- END BEADS INTEGRATION -->

<!-- BEGIN BEADS CODEX SETUP: generated by bd setup codex -->
## Beads Issue Tracker

Use Beads (`bd`) for durable task tracking in repositories that include it. Use the `beads` skill at `.agents/skills/beads/SKILL.md` (project install) or `~/.agents/skills/beads/SKILL.md` (global install) for Beads workflow guidance, then use the `bd` CLI for issue operations.

### Quick Reference

```bash
bd ready                # Find available work
bd show <id>            # View issue details
bd update <id> --claim  # Claim work
bd close <id>           # Complete work
bd prime                # Refresh Beads context
```

### Rules

- Use `bd` for all task tracking; do not create markdown TODO lists.
- Run `bd prime` when Beads context is missing or stale. Codex 0.129.0+ can load Beads context automatically through native hooks; use `/hooks` to inspect or toggle them.
- Keep persistent project memory in Beads via `bd remember`; do not create ad hoc memory files.

**Architecture in one line:** issues live in a local Dolt DB; sync uses `refs/dolt/data` on your git remote; `.beads/issues.jsonl` is a passive export. See https://github.com/gastownhall/beads/blob/main/docs/SYNC_CONCEPTS.md for details and anti-patterns.
<!-- END BEADS CODEX SETUP -->
