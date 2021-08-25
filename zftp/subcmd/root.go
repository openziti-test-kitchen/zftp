package subcmd

import (

	"github.com/spf13/cobra"
	"log"
	"zssh/zftplib"
)

const ExpectedServiceName = "zftp"

var flags = &zftplib.FtpFlags{}
var rootCmd = &cobra.Command{
	Use: "zftp",
	Long: "Zitified ftp",
	Short: "Zitified ftp",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init(){
	flags.InitFlags(rootCmd,ExpectedServiceName)
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

