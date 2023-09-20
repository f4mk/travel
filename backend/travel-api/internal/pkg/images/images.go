package images

type Status string

const (
	Pending Status = "pending"
	Loaded  Status = "loaded"
	Deleted Status = "deleted"
)
