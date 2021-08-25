package subcmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"zssh/zftplib"
)

var putCmd = &cobra.Command{
	Use: "put [LocalPath] <identityName>@<userName>:[RemotePath]",
	Long: "use to upload a singular file to ftp Server",
	Short: "use to upload a singular file to ftp Server",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		var remoteFilePath string
		var localFilePath string
		var err error

		localFilePath = args[0]
		remoteFilePath = args[1]

		localFilePath, err = filepath.Abs(localFilePath)
		if err != nil {
			logrus.Fatalf("cannot determine absolute local file path, unrecognized file name: %s", localFilePath)
		}
		if _, err := os.Stat(localFilePath); err != nil {
			logrus.Fatal(err)
		}
		flags.DebugLog("           local path: %s", localFilePath)

		data, err := os.Open(localFilePath)

		username, targetIdentity := flags.GetUserAndIdentity(remoteFilePath)
		remoteFilePath = zftplib.ParseFilePath(remoteFilePath)

		ftpConn := zftplib.EstablishClient(*flags,username,targetIdentity)
		defer ftpConn.Close()

		err = ftpConn.Login("Anonymous", "Anonymous")
		if err != nil {
			logrus.Fatalf("failure logining in to ftp Server")
		}

		ftpConn.Upload(data, remoteFilePath)

		if err := ftpConn.Quit(); err != nil  {
			logrus.Fatal(err)
		}


	},
}

func init() {
	rootCmd.AddCommand(putCmd)
}

