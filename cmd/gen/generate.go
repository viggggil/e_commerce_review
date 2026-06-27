package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"review_service/internal/conf"
	"strings"

	"github.com/go-kratos/kratos/v3/config"
	"github.com/go-kratos/kratos/v3/config/file"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string
	// flagDBSource overrides the database source in config.
	flagDBSource string
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
	flag.StringVar(&flagDBSource, "db_source", "", "database source override, eg: user:pass@tcp(127.0.0.1:3306)/review_service?charset=utf8mb4&parseTime=True&loc=Local")
}

func connectDB(cfg *conf.Data_Database) (*gorm.DB, error) {
	if cfg == nil {
		return nil, errors.New("GEN:connectDB fail, need config")
	}
	switch strings.ToLower((cfg.GetDriver())) {
	case "mysql":
		db, err := gorm.Open(mysql.Open(cfg.GetSource()))
		if err != nil {
			return nil, fmt.Errorf("connect db fail: %w", err)
		}
		return db, nil
	}
	return nil, fmt.Errorf("GEN:connectDB unsupported db driver %q", cfg.GetDriver())
}

func main() {
	flag.Parse()

	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}
	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}
	if bc.Data == nil || bc.Data.Database == nil {
		panic(errors.New("GEN: missing data.database config"))
	}
	if source := strings.TrimSpace(os.Getenv("GEN_DATABASE_SOURCE")); source != "" {
		bc.Data.Database.Source = source
	}
	if source := strings.TrimSpace(flagDBSource); source != "" {
		bc.Data.Database.Source = source
	}

	g := gen.NewGenerator(gen.Config{
		OutPath:       "../../internal/data/query",
		Mode:          gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldNullable: true,
	})

	db, err := connectDB(bc.Data.Database)
	if err != nil {
		panic(fmt.Errorf("%w\n\nIf you are using local MySQL on Ubuntu, root may be configured with auth_socket and reject TCP password login. Create a dedicated user or pass a working DSN with -db_source or GEN_DATABASE_SOURCE", err))
	}
	g.UseDB(db)
	g.ApplyBasic(g.GenerateAllTable()...)
	g.Execute()
}
