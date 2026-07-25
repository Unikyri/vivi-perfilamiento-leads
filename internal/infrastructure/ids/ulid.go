package ids

import (
	"crypto/rand"
	"io"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
)

// GeneradorID creates opaque, monotonic ULIDs safe for concurrent callers.
type GeneradorID struct {
	mu      sync.Mutex
	entropy io.Reader
	now     func() time.Time
}

func NuevoGeneradorID() *GeneradorID {
	return &GeneradorID{entropy: ulid.Monotonic(rand.Reader, 0), now: time.Now}
}

var _ interface{ Nuevo() string } = (*GeneradorID)(nil)

func (g *GeneradorID) Nuevo() string {
	g.mu.Lock()
	defer g.mu.Unlock()
	return ulid.MustNew(ulid.Timestamp(g.now().UTC()), g.entropy).String()
}
