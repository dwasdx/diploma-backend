package api

// ContextKey is a value for use with context.WithValue
type ContextKey struct {
	name string
}

// NewContextKey returns new context key
func NewContextKey(name string) *ContextKey {
	return &ContextKey{name: name}
}
func (k *ContextKey) String() string {
	return k.name
}
