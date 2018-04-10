package rubbernecker

// PersistanceEngine interface should ensure any backing service will follow the
// same set of rules.
type PersistanceEngine interface {
	Get(key string) (interface{}, error)
	Put(key string, value interface{}) error
}
