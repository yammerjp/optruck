package actions

type Action interface {
	Run() error
}

var _ Action = (*UploadConfig)(nil)
var _ Action = (*TemplateConfig)(nil)
var _ Action = (*MirrorConfig)(nil)
