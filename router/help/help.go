package help

import (
    "fmt"
    "flag"
    "WireguardLDAPManager/router/reconfig"
    "WireguardLDAPManager/router/showconfig"
    "WireguardLDAPManager/router/genkey"
    "WireguardLDAPManager/router/delkey"
    "WireguardLDAPManager/router/clearkey"
)

var f *flag.FlagSet

var Usage func()

func Setup(name string) {
    f = flag.NewFlagSet(name, flag.ExitOnError)
    //flag
    f.Usage = Usage
}

func Run(args []string) {
    f.Parse(args)
    subargs := f.Args()
    if len(subargs) == 0 {
        subargs = append(subargs, "")
    }
    switch subargs[0] {
    case "":
        Usage()
    case "reconfig":
        reconfig.Usage()
    case "showconfig":
        showconfig.Usage()
    case "genkey":
        genkey.Usage()
    case "delkey":
        delkey.Usage()
    case "clearkey":
        clearkey.Usage()
    default:
        fmt.Fprintf(f.Output(), "Unknown command: '%s' (try without commands for a list of commands)\n", subargs[0])
    }
}
