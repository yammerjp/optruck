package actions

type Action interface {
	Run() error
}

var _ Action = (*MirrorConfig)(nil)
