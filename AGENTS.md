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
