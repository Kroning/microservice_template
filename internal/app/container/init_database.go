package container

import (
	"github.com/Kroning/example_service/pkg/storage/postgresql"
)

// initDatabase initializes the database connection and runs migrations if enabled.
func (c *Container) initDatabase() {
	cfg := c.App.Cfg.DB

	c.DB = postgresql.New(c.App.Ctx, postgresql.Config{
		Master: postgresql.ReplicaConfig{
			Host:        cfg.Master.Host,
			Port:        cfg.Master.Port,
			User:        cfg.Master.User,
			Password:    cfg.Master.Password.Value(),
			Database:    cfg.Master.Database,
			MaxOpen:     cfg.Master.MaxOpen,
			MaxIdle:     cfg.Master.MaxIdle,
			MaxLifetime: cfg.Master.MaxLifetime,
			MaxIdleTime: cfg.Master.MaxIdleTime,
		},
		Metrics:        cfg.Metrics,
		Migrations:     cfg.Migrations,
		MigrationsPath: cfg.MigrationsPath,
	})
}
