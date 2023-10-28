package reconfig

import (
    "fmt"
    "flag"
    "WireguardLDAPManager/models/wireguard"
)

var f *flag.FlagSet

func Usage() {
    fmt.Fprintf(f.Output(), `  %s
      Reconfig wireguard server.

Options:
  -h
        Print help message
`, f.Name())
    f.PrintDefaults()
}

func Setup(name string) {
    f = flag.NewFlagSet(name, flag.ExitOnError)
    //flag
    f.Usage = Usage
}

func Run(args []string) {
    f.Parse(args)
    subargs := f.Args()
    
    run(subargs)
}

func run(subargs []string) {
    wireguard.Reconfig()
}
