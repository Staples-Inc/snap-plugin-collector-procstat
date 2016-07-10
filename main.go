package main

import (
	"os"

	"github.com/intelsdi-x/snap-plugin-collector-procstat/procstat"
	"github.com/intelsdi-x/snap/control/plugin"
)

func main() {
	plugin.Start(procstat.Meta(), procstat.New(), os.Args[1])
}
