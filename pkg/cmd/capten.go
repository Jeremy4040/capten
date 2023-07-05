package cmd

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type CLIFormatter struct {
}

func (f *CLIFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var levelColor *color.Color
	switch entry.Level {
	case logrus.InfoLevel:
		levelColor = color.New(color.FgGreen)
	case logrus.WarnLevel:
		levelColor = color.New(color.FgYellow)
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = color.New(color.FgRed, color.Bold)
	default:
		levelColor = color.New()
	}
	message := fmt.Sprintf("[%s] %s\n", levelColor.Sprint(strings.ToUpper(entry.Level.String())), entry.Message)
	return []byte(message), nil
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "capten",
	Short: "",
	Long:  `command line tool for building cluster`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	logrus.SetFormatter(&CLIFormatter{})
	cobra.CheckErr(rootCmd.Execute())
}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "for creation of resources or cluster",
	Long:  ``,
}

var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "destroy created cluster",
	Long:  ``,
}

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "sets up cluster for usage",
	Long:  ``,
}

func readAndValidClusterFlags(cmd *cobra.Command) (clusterType string, cloudType string, err error) {
	clusterType, _ = cmd.Flags().GetString("type")
	if len(clusterType) == 0 {
		clusterType = "k3s"
	}
	if clusterType != "k3s" {
		err = fmt.Errorf("cluster type '%s' is not supported, supported types: k3s", clusterType)
		return
	}

	cloudType, _ = cmd.Flags().GetString("cloud")
	if len(cloudType) == 0 {
		cloudType = "aws"
	}
	if cloudType != "aws" {
		err = fmt.Errorf("cloud service '%s' is not supported, supported cloud serivces: aws", cloudType)
		return
	}
	return
}

func init() {
	clusterCreateSubCmd.PersistentFlags().String("cloud", "", "cloud service (default: aws)")
	clusterDestroySubCmd.PersistentFlags().String("cloud", "", "cloud service (default: aws)")
	clusterCreateSubCmd.PersistentFlags().String("type", "", "type of cluster (default: k3s)")
	clusterDestroySubCmd.PersistentFlags().String("type", "", "type of cluster (default: k3s)")

	createCmd.AddCommand(clusterCreateSubCmd)
	destroyCmd.AddCommand(clusterDestroySubCmd)
	setupCmd.AddCommand(appsCmd)
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(destroyCmd)
	rootCmd.AddCommand(setupCmd)
}
