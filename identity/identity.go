// Package identity holds the three facts that make a server one player: a
// root seed, a creation time, and an ed25519 keypair. See ADR 0001.
package identity

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

// FileName is the identity file inside a data directory.
const FileName = "identity.json"

// Identity is a server's persistent state.
type Identity struct {
	Seed       int64              `json:"seed"`
	CreatedAt  time.Time          `json:"created_at"`
	PrivateKey ed25519.PrivateKey `json:"private_key"`
}

// PublicKey returns the key that identifies this server to other servers.
func (id Identity) PublicKey() ed25519.PublicKey {
	return id.PrivateKey.Public().(ed25519.PublicKey)
}

// Generate creates a fresh identity with a random seed and keypair.
func Generate(now time.Time) (Identity, error) {
	var buf [8]byte
	if _, err := rand.Read(buf[:]); err != nil {
		return Identity{}, fmt.Errorf("identity: seed: %w", err)
	}
	_, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return Identity{}, fmt.Errorf("identity: keypair: %w", err)
	}
	return Identity{
		Seed:       int64(binary.BigEndian.Uint64(buf[:])),
		CreatedAt:  now.UTC().Truncate(time.Second),
		PrivateKey: priv,
	}, nil
}

// Load reads the identity from dir.
// It returns fs.ErrNotExist when no identity has been written.
func Load(dir string) (Identity, error) {
	data, err := os.ReadFile(filepath.Join(dir, FileName))
	if err != nil {
		return Identity{}, err
	}
	var id Identity
	if err := json.Unmarshal(data, &id); err != nil {
		return Identity{}, fmt.Errorf("identity: parse: %w", err)
	}
	if len(id.PrivateKey) != ed25519.PrivateKeySize {
		return Identity{}, errors.New("identity: private key has the wrong size")
	}
	return id, nil
}

// Save writes the identity to dir, creating it if needed.
// The file is readable only by the owner.
func Save(dir string, id Identity) error {
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(id, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, FileName), data, 0o600)
}

// LoadOrGenerate loads the identity from dir, generating and saving one on
// first start.
func LoadOrGenerate(dir string, now time.Time) (Identity, error) {
	id, err := Load(dir)
	if err == nil {
		return id, nil
	}
	if !errors.Is(err, fs.ErrNotExist) {
		return Identity{}, err
	}
	id, err = Generate(now)
	if err != nil {
		return Identity{}, err
	}
	return id, Save(dir, id)
}
