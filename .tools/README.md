# Project-local Gentle AI

This directory vendors the official Linux AMD64 **Gentle AI v1.46.0** executable for this workspace only.

- **Binary:** `./.tools/gentle-ai`
- **Release:** `https://github.com/Gentleman-Programming/gentle-ai/releases/tag/v1.46.0`
- **Upstream tag/source commit:** `v1.46.0` / `b22a7eb8730e0e255c7a6d142aedfc606cbb020e`
- **Integrity:** `gentle-ai_1.46.0_linux_amd64.tar.gz` was verified against the official release `checksums.txt` with `sha256sum --check --strict` before extracting the binary.

Use this executable for repository-local `sdd-status` and `sdd-continue` checks; do not invoke its `install`, `sync`, or `upgrade` commands because those manage global user configuration.

## Capability boundary

The official v1.46.0 CLI exposes `sdd-status` and `sdd-continue`, but does **not** expose `sdd-verify-validate`, `sdd-attempt`, or `review`. The project must not manufacture verify-admission or review-receipt artifacts in their absence. Consequently, this version upgrade alone cannot unblock Issue #15's formal verification/archive/closure gates.
