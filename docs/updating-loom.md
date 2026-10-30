# Updating loom in this repo

How to move this repo's installed framework (agents, skills, rules and platform configs) to the
version in your `loom` binary, and keep it current. `scripts/update-framework.sh` does the work;
this page explains what it does and what to review.

**Today**: ai-assistant-dot-files **v3.1.0**, installed 2026-08-01 (commit `9e9daa2`), as committed
copies, with no install manifest. **Target**: v3.7.0, or whatever `loom install --dry-run` reports.

## Why a plain `loom install` is not enough here

`loom install` overwrites only paths its manifest (`.loom-manifest.json`) says it installed. The
v3.1.0 install predates the manifest. So a plain install treats all 233 existing framework files as
this repo's own and skips them; `--force` does not change that. It would add the 90 new files and
leave the rest at v3.1.0. (`loom update` would be the obvious command, but it is not implemented yet.)

The script fixes this, in four steps:

1. A dry run lists the platforms loom detects here, and every existing file it would skip as "not
   installed by loom".
2. It removes exactly those files. Each is at a path loom ships, so each is an old framework copy.
3. It installs one platform at a time, so detection does not depend on the files it just removed.
   The manifest merges across runs.
4. It lists the framework files loom no longer ships, for you to review.

Tested on a clean clone of this repo (2026-10-29, v3.7.0): 233 files taken over, 229 updated, 46
added, 0 still skipped, and 91 stale files reported.

## Before you start

1. **Build a current `loom`**, from ai-assistant-dot-files:
   ```bash
   cd ~/Projects/Rieken/ai-assistant-dot-files && git pull && go install ./cmd/loom
   ```
   The framework is embedded in the binary, so the binary's build is the version you install.
2. **Commit or stash this repo's working tree.** The script refuses to run on a dirty tree, so the
   update is one reviewable diff.
3. **Clean up leftovers**: `git worktree prune` (`.claude/worktrees/lucid-kare` is prunable), and
   delete `.cursor/*.bak.*` and `DOMAIN_DICTIONARY.md.bak.*` once you have checked them.

## Update

```bash
git switch -c chore/framework-update
scripts/update-framework.sh            # report only: version, platforms, files to take over
scripts/update-framework.sh --apply    # do it; changes are staged, not committed
```

## Review

**Platforms.** On this machine loom detects seven: claude-code, cursor, github-copilot,
openai-codex, roo-code, jetbrains and cline. Only five have their markers committed:
`.aiassistant`, `.junie` and `.clinerules` are untracked here. A clone, CI, or the factory's
worktrees see five. Commit those directories or drop the platforms; do not leave it to each machine.

**Stale files** (`stale-framework-files.txt`, 91 in the test run). Delete them, then delete the list:

- `.agents/` (all of it): v3.7.0 no longer writes it. If you use Antigravity, run
  `loom install --target . --copy --platform gemini` and check where it writes before deleting.
- `.claude/agents/test-driven-developer.md` and its `.cursor` copy: retired by ADR-009.
- `go-backend` rules in `.cursor/rules` and `.github/instructions`: renamed upstream.
- `.windsurfrules`: Windsurf is not detected here (no `.windsurf/`).
- `AGENTS.md`: openai-codex now writes `.openai.md`. But `AGENTS.md` is also how loom detects Codex
  when there is no `.codex/`. Delete it only if you do not use Codex, or add `.codex/` first.

```bash
xargs -a stale-framework-files.txt git rm -q -- && rm stale-framework-files.txt
```

**Files the update never touches**, because they are this repo's: `CLAUDE.md`,
`ARCHITECTURE_RULES.md`, `DOMAIN_DICTIONARY.md`, `PROJECT_OVERVIEW.md`, `features/`, `docs/`.
Check two by hand:

- `CLAUDE.md`'s "Feature Delivery Workflow" section (around line 56) may describe the retired
  markdown pipeline. Update that section from the framework template; keep Saturday's constraints.
- `DOMAIN_DICTIONARY.md` was overwritten by the v3.1.0 install (`9e9daa2` changed it and left a
  `.bak`). Make sure it holds Saturday's terms.

## Verify and commit

```bash
test -f .loom-manifest.json && echo "manifest written"
loom health
loom intake check                 # 13 orders valid (checked 2026-10-29)
git diff --cached --stat | tail -1
git commit -m "chore(framework): upgrade ai-assistant-dot-files to v3.7.0" \
  -m "First install with a manifest. Old framework copies were taken over and reinstalled, so
later updates can tell loom's files from this repo's."
```

Merge it before promoting any order in `intake/`. The factory runs the agents committed on `main`, so
orders promoted earlier would run on v3.1.0 agents.

## Every update after that

```bash
cd ~/Projects/Rieken/ai-assistant-dot-files && git pull && go install ./cmd/loom
cd ~/Projects/Rieken/saturday-monorepo && git switch -c chore/framework-update
scripts/update-framework.sh --apply
```

With a manifest in place, the script takes nothing over. It reports 0 to take over, loom updates
every file it installed and nobody edited, and the stale list catches files upstream retired.
`loom install` still reports skips, with a reason:

- `edited since loom installed it`: someone changed a framework file here. Move the change into
  `shared/` in ai-assistant-dot-files if it is general, or rerun with `--force` to discard it.
- `not installed by loom`: this repo's own file. Nothing to do.

## Factory files

Not written by the update; this repo owns them:

- `.claude/delivery-policy.yaml`: see `docs/factory-setup.md`.
- `.claude/factory/target.yaml`: once intake v2 lands.
- `.gitignore`: add `.loom/` (factory worktrees and index), as ai-assistant-dot-files does.
