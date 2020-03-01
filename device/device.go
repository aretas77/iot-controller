package device

type NodeDevice struct {
	Network string

	Stop chan struct{}
}
