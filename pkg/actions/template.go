package actions

import (
	"fmt"
	"log/slog"

	"github.com/yammerjp/optruck/pkg/op"
	"github.com/yammerjp/optruck/pkg/output"
)

type TemplateConfig struct {
	Logger    *slog.Logger
	Target    op.Target
	Dest      output.Dest
	Overwrite bool
}

func (config TemplateConfig) Run() error {
	config.Logger.Info("Starting template action for item", "item", config.Target.ItemName)

	secretReference, err := config.Target.BuildClient().GetItem()
	if err != nil {
		return fmt.Errorf("failed to get secret reference: %w", err)
	}
	config.Logger.Info("Fetched secret reference", "secretReference", secretReference)

	config.Logger.Info("Starting template action for item", "item", config.Target.ItemName)
	err = config.Dest.Write(secretReference, config.Overwrite)
	if err != nil {
		return fmt.Errorf("failed to write secret reference: %w", err)
	}
	config.Logger.Info("Template written to %s successfully", "path", config.Dest.GetPath())

	return nil
}
