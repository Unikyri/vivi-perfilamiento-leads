package ids

import (
	"sync"
	"testing"

	"github.com/oklog/ulid/v2"
)

func TestULID_ConcurrentOpaque(t *testing.T) {
	g := NuevoGeneradorID()
	const count = 128
	ids := make(chan string, count)
	var group sync.WaitGroup
	for i := 0; i < count; i++ {
		group.Add(1)
		go func() { defer group.Done(); ids <- g.Nuevo() }()
	}
	group.Wait()
	seen := make(map[string]struct{}, count)
	for i := 0; i < count; i++ {
		value := <-ids
		if len(value) != 26 || len(value) > 40 {
			t.Fatalf("ID length = %d, want 26 and <=40: %q", len(value), value)
		}
		if _, err := ulid.Parse(value); err != nil {
			t.Fatalf("ID is not parseable ULID: %v", err)
		}
		if _, ok := seen[value]; ok {
			t.Fatalf("duplicate ID %q", value)
		}
		seen[value] = struct{}{}
	}
}
