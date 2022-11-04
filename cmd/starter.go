package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/api"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/common"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/config"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/db/postgresdb"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/logger"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/route"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/server"
	"github.com/spf13/cobra"
)

// Instance - Represents an instance of the server
var gatewayServer *server.Instance

var rootCmd = &cobra.Command{
	Use:   "EdgeCraft Web Server",
	Short: "A gateway server for communication with EdgeCraft Web Server",
	Long: `EdgeCraft Web Server is a gateway server for managing kubernetes's resource.
server run command
  'edgecraft'
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {

		// Interrupt handler
		c := make(chan os.Signal)
		go func() {
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			logger.Infof("Received %s signal", <-c)
			gatewayServer.Shutdown()
		}()

		// Start server
		gatewayServer.Start()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// init - called on package loading
func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	// k8s config flag
	rootCmd.Flags().String("kubeconfig", "", "The path to the kubeconfig used to connect to the kubernetes API server and the Kubelets")
}

func initConfig() {
	// create default logger
	err := logger.New()
	if err != nil {
		logger.Fatalf("Could not instantiate log %ss", err.Error())
	}

	// load config file
	conf, err := config.Load()
	if err != nil {
		logger.Fatalf("Could not load configuration: %s", err.Error())
		os.Exit(0)
	}

	// Make K8s Client for multiple clusters using File/ConfigMap
	kubeconfig, err := rootCmd.Flags().GetString("kubeconfig")
	if err != nil {
		logger.Warnf("invalid kubeconfig parameter", err.Error())
	}
	logger.Infof("specified kubeconfig parameter is '%s'", kubeconfig)

	config.SetupClusters(kubeconfig)

	// Load Message Files (i18n)
	common.LoadMessages(conf.API.LangPath, conf.API.Langs)

	// create server instance & initialize the server & workers
	gatewayServer = server.NewInstance(conf)
	gatewayServer.Init()

	// Establish database connection
	DB, err := postgresdb.NewConnection(conf.DB)
	if err != nil {
		logger.WithError(err).Fatal("Could not open database connection")
		return
	}
	// DB 정보 저장
	gatewayServer.DB = DB

	// create API instance & Initialize API
	api, err := api.New(conf.API, DB, &gatewayServer.Worker)
	if err != nil {
		logger.WithError(err).Fatal("Could not create api instance and initialize")
		return
	}

	// httpserver's route setting
	route.SetRoutes(api, gatewayServer)

	logger.Info("Intialize done.")
}
