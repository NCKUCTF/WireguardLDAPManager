package genkey

import (
    "fmt"
    "flag"
    "os"

    "WireguardLDAPManager/models/privatekey"
    "WireguardLDAPManager/models/wireguard"
)

var f *flag.FlagSet

func Usage() {
    fmt.Fprintf(f.Output(), `  %s
  or
  %s <key name>
      Generate a new wireguard key for this user.

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
    keyname := ""
    if len(subargs) == 1 {
        keyname = subargs[0]
    } else if len(subargs) > 1 {
        fmt.Fprintln(os.Stderr, "Bad args...")
        Usage()
        return
    }
    if keyname == "" || privatekey.ContainName(keyname) {
        if keyname != "" {
            fmt.Fprintln(os.Stderr, "Key name exist!")
        }
        keyname = privatekey.ReadNewName()
    }
    privatekey.Generate(keyname)
    wireguard.Reconfig()
    fmt.Println("Key generate success! Please use \"showconfig\" to get your wireguard config.")
}
