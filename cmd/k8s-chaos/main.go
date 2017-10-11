package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"

	"github.com/jonstacks/k8s-chaos/pkg/config"
	"github.com/jonstacks/k8s-chaos/pkg/pod"
)

func handleError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	client, err := config.GetClient()
	handleError(err)

	namespace := flag.String("namespace", "", "The kubernetes namespace to target. (Required)")
	regex := flag.String("regex", ".", "The regex to filter pods in a namespace")
	period := flag.Int("period", 60, "The amount of time in seconds between deleting")
	maxKill := flag.Int("max-kill", 1, "The maximum number of pods to kill")

	flag.Parse()

	if *namespace == "" {
		flag.PrintDefaults()
		os.Exit(2)
	}

	r, err := regexp.Compile(*regex)
	if err != nil {
		fmt.Printf("Could not compile regex '%s'", *regex)
	}

	opts := pod.DefaultKillerOpts()
	opts.Namespace = *namespace
	opts.Period = time.Duration(*period) * time.Second
	opts.Regex = r
	opts.MaxKill = *maxKill

	pk := pod.NewKiller(client, opts)
	pk.Start()

	select {
	case <-sigc:
		fmt.Println("Stopping pod killer...")
		pk.Stop()
	}
}
