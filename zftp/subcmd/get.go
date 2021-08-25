package subcmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"zssh/zftplib"
)

var getCmd = &cobra.Command{
	Use: "get <identityName>@<userName>:[RemotePath] [LocalPath]",
	Long: "use to download a singular file from ftp Server",
	Short: "use to download a singular file from ftp Server",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		var remoteFilePath string
		var localFilePath string
		var err error

		localFilePath = args[1]
		remoteFilePath = args[0]

		username, targetIdentity := flags.GetUserAndIdentity(remoteFilePath)
		remoteFilePath = zftplib.ParseFilePath(remoteFilePath)

		ftpConn := zftplib.EstablishClient(*flags,username,targetIdentity)
		defer ftpConn.Close()

		err = ftpConn.Login("Anonymous", "Anonymous")
		if err != nil {
			logrus.Fatalf("failure logining in to ftp Server")
		}

		lf, err := os.OpenFile(localFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
		if err != nil {
			 logrus.Errorf("error opening local file [%s] (%w)", localFilePath, err)
		}
		defer func() { _ = lf.Close() }()


		err = ftpConn.Download(remoteFilePath,lf)
		if err != nil {
			logrus.Fatalf("failed to download file: %s", remoteFilePath)
		}

		if err = ftpConn.Quit(); err != nil  {
			logrus.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}