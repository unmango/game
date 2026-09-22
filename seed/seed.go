// Package seed derives sub-seeds from a root seed and a path.
//
// Every node in a game is named by a slash separated path. The sub-seed for a
// path is a pure function of the root seed and the path, so a sub-game rooted
// at any path is its own deterministic world. See ADR 0002.
package seed

import (
	"crypto/sha256"
	"encoding/binary"
	"math/rand/v2"
	"strings"
)

// Derive returns the sub-seed for path under root.
//
// The value is the first eight bytes, big endian, of
// SHA-256(root as big endian int64 || path). The derivation is part of the
// framework's contract with every game built on it and never changes.
func Derive(root int64, path string) uint64 {
	h := sha256.New()
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], uint64(root))
	h.Write(buf[:])
	h.Write([]byte(path))
	return binary.BigEndian.Uint64(h.Sum(nil)[:8])
}

// Rand returns a deterministic generator for a sub-seed.
// Two calls with the same sub-seed yield the same sequence.
func Rand(sub uint64) *rand.Rand {
	return rand.New(rand.NewPCG(sub, 0))
}

// Join builds a path from segments, ignoring empty ones.
// Join("character", "strength") is "character/strength".
func Join(segments ...string) string {
	parts := segments[:0:0]
	for _, s := range segments {
		if s != "" {
			parts = append(parts, s)
		}
	}
	return strings.Join(parts, "/")
}
