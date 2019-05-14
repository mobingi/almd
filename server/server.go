package server

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/mobingi/oceand/pkg/util"
	"github.com/mobingi/oceand/pkg/watch"
	"github.com/mobingi/oceand/server/options"
)

const configPathEnv = "CONFIG_PATH"

func NewOceandCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "oceand",
		Long: "it is for collect kubernetes cluster status",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
	return cmd
}

func run() error {
	log.Print("running")
	configPath := util.ReadEnvOrDie(configPathEnv)
	log.Print("config path:", configPath)
	o, err := options.NewOptionsFromFile(configPath)
	if err != nil {
		log.Print("read config err:", err)
		return err
	}

	if err := watch.Watch(o.ID, o.Token); err != nil {
		return err
	}

	select {}
}
