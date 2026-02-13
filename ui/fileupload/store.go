package fileupload

import "sync"

// FileMeta holds metadata about an uploaded file.
type FileMeta struct {
	ID       string
	Name     string
	Size     int64
	MimeType string
}

// Store is a thread-safe in-memory store for uploaded file metadata,
// keyed by "sessionID:componentID".
type Store struct {
	mu    sync.Mutex
	files map[string][]FileMeta
}

// NewStore creates a new empty file metadata store.
func NewStore() *Store {
	return &Store{files: make(map[string][]FileMeta)}
}

// Add appends a file to the store under the given key.
func (s *Store) Add(key string, meta FileMeta) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.files[key] = append(s.files[key], meta)
}

// Remove deletes a file by ID from the store.
func (s *Store) Remove(key, fileID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	files := s.files[key]
	for i, f := range files {
		if f.ID == fileID {
			s.files[key] = append(files[:i], files[i+1:]...)
			return
		}
	}
}

// List returns all files stored under the given key.
func (s *Store) List(key string) []FileMeta {
	s.mu.Lock()
	defer s.mu.Unlock()
	dst := make([]FileMeta, len(s.files[key]))
	copy(dst, s.files[key])
	return dst
}

// Clear removes all files under the given key.
func (s *Store) Clear(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.files, key)
}
