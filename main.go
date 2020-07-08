package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/ghodss/yaml"
	"github.com/mitchellh/go-homedir"
)

var (
	hostsFile  = "/etc/hosts"
	configFile = filepath.Join(getHomedir(), ".config", "focus", "focus.yaml")
	backupFile = filepath.Join(getHomedir(), ".config", "focus", "hosts.bak")
)

var (
	profile  string
	duration time.Duration
)

type Config struct {
	Profiles map[string][]string `yaml:"profiles"`
}

func init() {
	flag.StringVar(&profile, "profile", "default", "The profile to use")
	flag.DurationVar(&duration, "timer", 0, "Stop blocking after a period of time")
}

func main() {
	flag.Parse()
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	// Set up the channel that we write to if an interrupt is sent
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Set up the channel that we write to once the timer ends
	timeout := make(chan time.Time, 1)
	if duration > 0 {
		go func() {
			t := <-time.After(duration)
			timeout <- t
		}()
	}

	cfg, err := readConfig(configFile)
	if err != nil {
		return err
	}
	websites, ok := cfg.Profiles[profile]
	if !ok {
		return fmt.Errorf("Profile %s not defined in config file", profile)
	}

	if err := backupHosts(); err != nil {
		return err
	}

	if err := appendEntriesToHosts(websitesToHostsEntries(websites)); err != nil {
		return err
	}

	log.Println("Blocking websites. Press ctrl+c to stop")
	if duration > 0 {
		log.Printf("%s remaining", duration)
	}

	select {
	case <-sigs:
		// break
	case <-timeout:
		log.Println("Time's up")
		err := beeep.Notify("Time's up", "Unblocking websites", "")
		if err != nil {
			return err
		}
	}

	if err := restoreHosts(); err != nil {
		return err
	}

	if err := deleteBackup(); err != nil {
		return err
	}

	return nil
}

func appendEntriesToHosts(entries string) error {
	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile(hostsFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintf(f, "\n%s\n", entries)

	log.Printf("Added entries to %s", hostsFile)

	return nil
}

func websitesToHostsEntries(websites []string) string {
	var entries []string
	for _, website := range websites {
		entry := fmt.Sprintf("127.0.0.1 %s www.%s", website, website)
		entries = append(entries, entry)
	}
	return strings.Join(entries, "\n")
}

func readConfig(file string) (*Config, error) {
	contents, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	if err := yaml.Unmarshal(contents, config); err != nil {
		return nil, err
	}
	return config, nil
}

func backupHosts() error {
	contents, err := ioutil.ReadFile(hostsFile)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(backupFile, contents, 0644); err != nil {
		return err
	}
	log.Printf("Backed up %s to %s", hostsFile, backupFile)
	return nil
}

func restoreHosts() error {
	contents, err := ioutil.ReadFile(backupFile)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(hostsFile, contents, 0644); err != nil {
		return err
	}
	log.Printf("Restored %s from %s", hostsFile, backupFile)
	return nil
}

func getHomedir() string {
	dir, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func deleteBackup() error {
	if err := os.Remove(backupFile); err != nil {
		return err
	}
	log.Printf("Removed old backup file %s", backupFile)
	return nil
}
