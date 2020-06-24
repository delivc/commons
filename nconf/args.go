package nconf

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

// RootArgs ..
type RootArgs struct {
	Prefix  string
	EnvFile string
}

// Setup ...
func (args *RootArgs) Setup(config interface{}, version string) (logrus.FieldLogger, error) {
	// first load the logger and BugSnag config
	rootConfig := &struct {
		Log     *LoggingConfig
		BugSnag *BugSnagConfig
	}{}
	if err := LoadFromEnv(args.Prefix, args.EnvFile, rootConfig); err != nil {
		return nil, errors.Wrap(err, "Failed to load the logging configuration")
	}

	log, err := ConfigureLogging(rootConfig.Log)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create the logger")
	}
	if version == "" {
		version = "unknown"
	}
	log = log.WithField("version", version)

	if err := SetupBugSnag(rootConfig.BugSnag, version); err != nil {
		return nil, errors.Wrap(err, "Failed to configure bugsnag")
	}

	if config != nil {
		// second load the config for this project
		if err := LoadFromEnv(args.Prefix, args.EnvFile, config); err != nil {
			return log, errors.Wrap(err, "Failed to load the config object")
		}
		log.Debug("Loaded configuration")
	}
	return log, nil
}

// MustSetup forces Setup
func (args *RootArgs) MustSetup(config interface{}, version string) logrus.FieldLogger {
	logger, err := args.Setup(config, version)
	if err != nil {
		if logger != nil {
			logger.WithError(err).Fatal("Failed to setup configuration")
		} else {
			panic(fmt.Sprintf("Failed to setup configuratio: %s", err.Error()))
		}
	}

	return logger
}

// ConfigFlag  reads the config flag
func (args *RootArgs) ConfigFlag() *pflag.Flag {
	return &pflag.Flag{
		Name:      "config",
		Shorthand: "c",
		Usage:     "A .env file to load configuration from",
		Value:     newStringValue("", &args.EnvFile),
	}
}

// PrefixFlag ...
func (args *RootArgs) PrefixFlag() *pflag.Flag {
	return &pflag.Flag{
		Name:      "prefix",
		Shorthand: "p",
		Usage:     "A prefix to search for when looking for env vars",
		Value:     newStringValue("", &args.Prefix),
	}
}

type stringValue string

func newStringValue(val string, p *string) *stringValue {
	*p = val
	return (*stringValue)(p)
}

func (s *stringValue) Set(val string) error {
	*s = stringValue(val)
	return nil
}
func (s *stringValue) Type() string   { return "string" }
func (s *stringValue) String() string { return string(*s) }
