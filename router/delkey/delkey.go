package delkey

import (
    "log"
    "fmt"
    "flag"

    "WireguardLDAPManager/utils/config"
    "WireguardLDAPManager/models/privatekey"
    "WireguardLDAPManager/models/wireguard"
)

var f *flag.FlagSet

func Usage() {
    fmt.Fprintf(f.Output(), `  %s
  or
  %s <key name 1> <key name 2> ... <key name n>
      Delete wireguard keys for this user.

Options:
  -h    Print help message
`, f.Name(), f.Name())
    f.PrintDefaults()
}

func Setup(name string) {
    f = flag.NewFlagSet(name, flag.ExitOnError)
    f.Usage = Usage
}

func Run(args []string) {
    f.Parse(args)
    subargs := f.Args()
    
    run(subargs)
}

func run(subargs []string) {
    keynames := subargs
    if len(keynames) == 0 {
        for asking := true; asking; asking = config.Ask("Continue input next key name?"){
            keynames = append(keynames, privatekey.ReadName())
        }
    }
    for _, name := range keynames {
        if !privatekey.ContainName(name) {
            log.Fatalf("Key \"%s\" name exist!", name)
        }
    }
    privatekey.Delete(keynames)
    wireguard.Reconfig()
    fmt.Println("Key delete success!")
}
