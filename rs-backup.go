package main

import (
	"encoding/base64"
	"fmt"
	"github.com/0x111/rs-backup/helpers"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/smtp"
	"os"
	"os/exec"
	"strconv"
	"strings"
	_ "text/template"
	"time"
)

func main() {
	var err error
	// log start time
	start := time.Now()
	log.Printf("Starting backup process %s", start)
	// load config
	viper.SetConfigName("rs-backup")       // name of config file (without extension)
	viper.AddConfigPath("/etc/rs-backup/") // path to look for the config file in
	//viper.AddConfigPath("$HOME/.appname")  // call multiple times to add many search paths
	viper.AddConfigPath(".") // optionally look for config in the working directory

	err = viper.ReadInConfig() // Find and read the config file
	if err != nil {            // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	// Look for rsync, if non-existent don't allow running the app
	path, err := exec.LookPath("rsync")
	// If rsync not found, exit the execution of the program
	if err != nil {
		log.Fatal("We could not find rsync in your path!")
	}

	// cmd
	var args []string
	// Dynamically append arguments to the args
	if viper.GetBool("show_progress") == true {
		args = append(args, "--progress")
	}

	if viper.GetBool("force_ipv4") == true {
		args = append(args, "--ipv4")
	}

	if viper.GetBool("archive") == true {
		args = append(args, "--archive")
	}

	if viper.GetBool("verbose") == true {
		args = append(args, "--verbose")
	}

	if viper.GetBool("compress") == true {
		args = append(args, "--compress")
	}

	remoteShellCommand := viper.GetString("remote_shell_command")

	if remoteShellCommand != "" && len(remoteShellCommand) > 0 {
		args = append(args, "--rsh")
		sshCmd := fmt.Sprintf("%s", remoteShellCommand)
		args = append(args, sshCmd)
	}

	// create a temporary file
	file, err := ioutil.TempFile("", "rs-backup")
	logFileName := file.Name() + ".log"
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(logFileName)

	if viper.GetBool("log") == true {
		args = append(args, "--log-file="+logFileName)
	}

	localDirectoryPath := viper.GetString("local_directory_path")

	if localDirectoryPath != "" && len(localDirectoryPath) > 0 {
		args = append(args, localDirectoryPath)
	}

	remoteDirectoryPath := viper.GetString("remote_directory_path")

	if remoteDirectoryPath != "" && len(remoteDirectoryPath) > 0 {
		args = append(args, remoteDirectoryPath)
	}

	log.Println("Running command and waiting for it to finish...")
	_, err = exec.Command(path, args...).Output()

	if err != nil {
		log.Fatal(fmt.Printf("Command finished with error: %v", err))
	}

	sendTo := viper.GetStringSlice("mail.to")
	// Content-Type: text/html; charset="UTF-8";
	fileContent := helpers.ReadFileContent(logFileName)
	mailBody := `From: %s
To: %s
Subject: Backup run at %s
Content-Type: multipart/mixed; boundary=_rssbckkgthbscrpt14467_
Content-Transfer-Encoding: 7bit
--_rssbckkgthbscrpt14467_
MIME-Version: 1.0
Content-Type: text/html; charset="UTF-8";
Content-Transfer-Encoding: 7bit

Backup started at: %s<br />
Backup ended at: %s <br />
Backup took: %s <br />
Find the contents of the rsync log in the attached log file.
<br />

Mailed by <a href="https://github.com/0x111/rs-backup">0x111/rs-backup</a>.<br />
--_rssbckkgthbscrpt14467_
Content-Type: application/octet-stream; name="rsync.log"
Content-Disposition: attachment; filename="rsync.log"
Content-Transfer-Encoding: base64

%s
--_rssbckkgthbscrpt14467_--
`

	// variables to make ExamplePlainAuth compile, without adding
	// unnecessary noise there.
	var (
		from       = viper.GetString("mail.from")
		recipients = sendTo
	)

	end := time.Now().Format(time.RFC822)
	mailBody = fmt.Sprintf(mailBody, from, strings.Join(sendTo, ","), start.Format(time.RFC822), start.Format(time.RFC822), end, time.Since(start), base64.StdEncoding.EncodeToString(fileContent))
	msg := []byte(mailBody)

	// hostname is used by PlainAuth to validate the TLS certificate.
	hostname := viper.GetString("smtp.host")
	port := strconv.Itoa(viper.GetInt("smtp.port"))
	auth := smtp.PlainAuth("", viper.GetString("smtp.user"), viper.GetString("smtp.password"), hostname)

	err = smtp.SendMail(hostname+":"+port, auth, from, recipients, msg)
	if err != nil {
		log.Fatal(err)
	}
}
