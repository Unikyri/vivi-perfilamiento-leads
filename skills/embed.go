// Package skillassets exposes the repository's prompt skills as an embedded filesystem.
package skillassets

import "embed"

// FS contains all Issue #24 SKILL.md files. Keeping the embed directive at the
// skills root avoids the forbidden .. patterns that Go rejects in go:embed.
//
//go:embed */SKILL.md
var FS embed.FS
