# Archive Report: Issue #27 — Deployable Demo

## Closure

- **Change:** `issue-27-deploy-demo`
- **Issue:** #27 — `[B6] Deploy a Heroku, seed del demo y guion de entrega`
- **Archived:** 2026-07-26
- **Branch:** `chore/archive-issue-27`
- **Base:** `origin/main` (`852ba01`)
- **Workspace:** `/tmp/vivi-issue-27-archive`
- **Artifact store:** hybrid (Engram primary + OpenSpec repository)
- **Final archive:** `openspec/changes/archive/2026-07-26-issue-27-deploy-demo/`
- **Status:** archived and closed
- **Tasks:** 14/14 complete; 0 unchecked implementation or remediation tasks
- **Product scope:** no product code, configuration, deploy state, secrets, provider, index, or commit was changed by archive

## Source-of-Truth Spec Synchronization

All four change-domain specs were complete new main specifications; no corresponding main spec existed. Each delta was copied byte-for-byte before the change folder was moved.

| Domain | Action | Source | Destination | SHA-256 |
|---|---|---|---|---|
| `demo-control` | Created | `openspec/changes/issue-27-deploy-demo/specs/demo-control/spec.md` | `openspec/specs/demo-control/spec.md` | `afe8487e12c80c4c03af7fca2ff4e731ad0d064e54916afe4457eaaf19475b00` |
| `despliegue-heroku` | Created | `openspec/changes/issue-27-deploy-demo/specs/despliegue-heroku/spec.md` | `openspec/specs/despliegue-heroku/spec.md` | `b6aea524cb33588182150933ccb45c3dba84924b58f5373c6051909ed9eed292` |
| `estaticos-spa` | Created | `openspec/changes/issue-27-deploy-demo/specs/estaticos-spa/spec.md` | `openspec/specs/estaticos-spa/spec.md` | `7d7f6071d5814cc3d020496f7c826aad99e382ffc22eb7665d08dd32b68de628` |
| `seed-demo` | Created | `openspec/changes/issue-27-deploy-demo/specs/seed-demo/spec.md` | `openspec/specs/seed-demo/spec.md` | `39f9f75dd3d13794b7d3a644f6bf161a38910f07f0fcd11171a106dfdf636889` |

The active change directory was then moved intact to the dated archive. No delta section required destructive modification of an existing main spec.

## Verification and Task Evidence

- Native status: `nextRecommended: archive`, `archive: ready`, `blockedReasons: []`.
- Native task progress: 14 total, 14 completed, 0 pending, `allComplete: true`.
- Persisted tasks artifact: every task `1.1`–`5.2` and remediation task `R1`/`R2` is checked `[x]`.
- Verification envelope: `gentle-ai.verify-result/v1`, `pass_with_warnings`, 0 blockers, 0 critical findings, 8/8 requirements, 12/12 scenarios.
- Verification evidence revision: `sha256:2b263a3bdf6f2177775dbbdb13594dda5239f68ee90faaa96b1552b85e78b0eb`.
- Verification commands: `go test ./... -count=1` exit 0; `go build ./...` exit 0; recorded output hashes are retained in `verify-report.md`.
- Runtime authority: attempt ordinal 1 passed, `complete: true`, `decision_required: false`, `next_action: complete`, with no active attempt.
- Five verification warnings remain non-blocking and are preserved in the archived report: no live PostgreSQL run, documentation-governed acceptance guidance, operator-run public `/salud`, pre-existing npm audit findings, and unavailable Wiki Doc 12 source.

## Native Review Authority and Binding

The exact native command used before archive was:

```text
/tmp/gentle-ai-2.2.0-rc.1-inspect/gentle-ai_2.2.0-rc.1_linux_amd64 review validate --gate post-apply --cwd /tmp/vivi-issue-27-archive
```

Result: **allow** (`allowed: true`, `action: continue`), because the authoritative transaction, current repository target, and content-bound artifacts match.

| Field | Value |
|---|---|
| Lineage | `review-3225019e6219019b` |
| Generation | `1` |
| Terminal state | `approved` |
| Binding revision | `sha256:54c70eefcc12b7aac61b848f4c7ec8b513a0ab7a26f22fc934ab26e64f59d3ea` |
| Authority/store revision | `sha256:0f7d09300358696106bd49a02459b17623fdb85f11d578d85188d77abda350af` |
| Receipt hash | `sha256:9897ae65b1c430bd4ab0e9230a8ddee16748a8bcbfb9b1dcfd7256e909ad5665` |
| Base tree | `0319db36af47893a41d60a7f0f317200d4347b28` |
| Candidate tree | `3950efbbd67f87d6a0c29fed73695d0636b40f4d` |
| Paths digest | `sha256:b3ce6147e149fd060ad1315f1ce0efbf6aa1c98da8b275acd05fb45e68e2edde` |
| Policy hash | `sha256:34fb63d7f29f8613cd4431382b1057398a4816f8a4c20fc34677fffc80a184f6` |
| Ledger hash | `sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| Fix delta hash | `sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` (empty) |
| Evidence hash | `sha256:9e541403149f70ab1d2a124558888307133010e229c0489cf1acc0d904b126a5` |
| Base relationship | valid |
| Review risk | high; four selected lenses; zero findings |

The immutable authority files are Git-common-dir state, not project artifacts:

- `/home/daikyri/Workspace/Hackathon-Colsubsidio/vivi-perfilamiento-leads/.git/gentle-ai/review-transactions/v2/review-3225019e6219019b/review-state.json`
- `/home/daikyri/Workspace/Hackathon-Colsubsidio/vivi-perfilamiento-leads/.git/gentle-ai/review-transactions/v2/review-3225019e6219019b/review-receipt.json`
- `/home/daikyri/Workspace/Hackathon-Colsubsidio/vivi-perfilamiento-leads/.git/gentle-ai/review-transactions/v2/review-3225019e6219019b/finalize-attempt-journal.json`
- `/home/daikyri/Workspace/Hackathon-Colsubsidio/vivi-perfilamiento-leads/.git/gentle-ai/review-transactions/v2/review-3225019e6219019b/final-evidence/verification.txt`

The persisted `verify-report.md` retains its own historical verification narrative and earlier review metadata; the native post-apply authority above is the current archive gate and was independently validated immediately before archive.

## Engram Traceability

Full Engram observations read before archive:

| Artifact | Observation |
|---|---:|
| project SDD initialization | `#421` |
| state (pre-archive) | `#617` |
| proposal | `#619` |
| spec | `#620` |
| design | `#621` |
| tasks | `#622` |
| apply-progress | `#624` |
| verify-report | `#625` |

No exact Engram review topics were present for `review/transaction`, `review/ledger`, `review/receipt`, or `review/gate-context`; native Git-common-dir authority is therefore the authoritative review evidence. The Engram archive report is persisted under `sdd/issue-27-deploy-demo/archive-report` after this file is written.

## Archived Contents and Verification Checklist

Archive path: `openspec/changes/archive/2026-07-26-issue-27-deploy-demo/`

- `proposal.md` ✅
- `specs/demo-control/spec.md` ✅
- `specs/despliegue-heroku/spec.md` ✅
- `specs/estaticos-spa/spec.md` ✅
- `specs/seed-demo/spec.md` ✅
- `design.md` ✅
- `tasks.md` ✅ (14/14 checked)
- `apply-progress.md` ✅
- `verify-report.md` ✅
- `state.yaml` ✅ (closed archive state)
- `archive-report.md` ✅
- Active directory `openspec/changes/issue-27-deploy-demo` absent ✅
- Main specs present and byte-identical to their deltas ✅
- No product/config/deploy/secrets/provider changes ✅
- No commit or push performed ✅

**Status: archived and closed.**
