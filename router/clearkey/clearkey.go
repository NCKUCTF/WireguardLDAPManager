package clearkey

import (
    "fmt"
    "flag"

    "WireguardLDAPManager/utils/config"
    "WireguardLDAPManager/models/privatekey"
    "WireguardLDAPManager/models/wireguard"
)

var f *flag.FlagSet
var noask bool

func Usage() {
    fmt.Fprintf(f.Output(), `  %s
      Clear all wireguard keys for this user.

Options:
  -h    Print help message
`, f.Name())
    f.PrintDefaults()
}

func Setup(name string) {
    f = flag.NewFlagSet(name, flag.ExitOnError)
    f.BoolVar(&noask, "y", false, "Run anyway without ask")
    f.Usage = Usage
}

func Run(args []string) {
    f.Parse(args)
    subargs := f.Args()
    
    run(subargs)
}

func run(subargs []string) {
    if !noask && !config.Ask("Continue clear all keys?") {
        fmt.Println("Key clear cancel.")
        return
    }

    privatekey.Clear()
    wireguard.Reconfig()
    fmt.Println("Key clear success!")
}
