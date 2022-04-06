package repositories

type ErrNotFound struct {
}

func (s ErrNotFound) Error() string {
	return "object not found in repository"
}
