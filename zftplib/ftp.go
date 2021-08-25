/*
	Copyright NetFoundry, Inc.

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

	https://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

package zftplib

import (
	"fmt"
	"github.com/gonutz/ftp-client/ftp"
	"github.com/openziti/sdk-golang/ziti"
	"github.com/openziti/sdk-golang/ziti/config"
	"github.com/pkg/errors"
	"github.com/pkg/sftp"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)


func SendFile(client *sftp.Client, localPath string, remotePath string) error {
	localFile, err := ioutil.ReadFile(localPath)

	if err != nil {
		return errors.Wrapf(err, "unable to read local file %v", localFile)
	}

	rmtFile, err := client.OpenFile(remotePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC)

	if err != nil {
		return errors.Wrapf(err, "unable to open remote file %v", remotePath)
	}
	defer rmtFile.Close()

	_, err = rmtFile.Write(localFile)
	if err != nil {
		return err
	}

	return nil
}

func RetrieveRemoteFiles(client *sftp.Client, localPath string, remotePath string) error {

	rf, err := client.Open(remotePath)
	if err != nil {
		return fmt.Errorf("error opening remote file [%s] (%w)", remotePath, err)
	}
	defer func() { _ = rf.Close() }()

	lf, err := os.OpenFile(localPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error opening local file [%s] (%w)", localPath, err)
	}
	defer func() { _ = lf.Close() }()

	_, err = io.Copy(lf, rf)
	if err != nil {
		return fmt.Errorf("error copying remote file to local [%s] (%w)", remotePath, err)
	}
	logrus.Infof("%s => %s", remotePath, localPath)

	return nil
}

func EstablishClient(f FtpFlags, userName string, targetIdentity string) *ftp.Connection {
	ctx := ziti.NewContextWithConfig(getConfig(f.ZConfig))
	_, ok := ctx.GetService(f.ServiceName)
	if !ok {
		logrus.Fatalf("service not found: %s", f.ServiceName)
	}

	dialOptions := &ziti.DialOptions{
		ConnectTimeout: 0,
		Identity:       targetIdentity,
		AppData:        nil,
	}
	svc, err := ctx.DialWithOptions(f.ServiceName, dialOptions)
	if err != nil {
		logrus.Fatalf("error when dialing service name %s. %v", f.ServiceName, err)
	}

	ftpConn, err := ftp.ConnectOn(svc)
	if err != nil {
		log.Fatal(err)
	}

	return ftpConn
}

func (f *FtpFlags) DebugLog(msg string, args ...interface{}) {
	if f.Debug {
		logrus.Infof(msg, args...)
	}
}

func getConfig(cfgFile string) (zitiCfg *config.Config) {
	zitiCfg, err := config.NewFromFile(cfgFile)
	if err != nil {
		log.Fatalf("failed to load ziti configuration file: %v", err)
	}
	return zitiCfg
}

// AppendBaseName tags file name on back of remotePath if the path is blank or a directory/*
func AppendBaseName(c *sftp.Client, remotePath string, localPath string, debug bool) string {
	localPath = filepath.Base(localPath)
	if remotePath == "" {
		remotePath = filepath.Base(localPath)
	} else {
		info, err := c.Lstat(remotePath)
		if err == nil && info.IsDir() {
			remotePath = filepath.Join(remotePath, localPath)
		} else if debug {
			logrus.Infof("Remote File/Directory: %s doesn't exist [%v]", remotePath, err)
		}
	}
	return remotePath
}

func After(value string, a string) string {
	// Get substring after a string.
	pos := strings.LastIndex(value, a)
	if pos == -1 {
		return ""
	}
	adjustedPos := pos + len(a)
	if adjustedPos >= len(value) {
		return ""
	}
	return value[adjustedPos:len(value)]
}
