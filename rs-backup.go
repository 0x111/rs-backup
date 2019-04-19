package main

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os/exec"
)

func main() {
	var err error
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

	fmt.Println(path)

	// cmd

	cmd := exec.Command("sleep", "5")
	log.Println("Running command and waiting for it to finish...")
	err = cmd.Run()
	log.Printf("Command finished with error: %v", err)
}
