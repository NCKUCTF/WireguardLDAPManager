package router

import (
    "fmt"
    "flag"
    "WireguardLDAPManager/router/help"
    "WireguardLDAPManager/router/reconfig"
    "WireguardLDAPManager/router/showconfig"
    "WireguardLDAPManager/router/genkey"
    "WireguardLDAPManager/router/delkey"
    "WireguardLDAPManager/router/clearkey"
)

var f *flag.FlagSet

func Usage() {
    fmt.Fprintf(f.Output(), `Usage: %s COMMAND

Options:
  -h    Print help message
`, f.Name())
    f.PrintDefaults()

    fmt.Fprintf(f.Output(), `
A list of commands is shown below. To get detailed usage and help for a
command, run:
  %s help COMMAND

Here is the list of commands available with a short syntax reminder. Use the
'help' command above to get full usage details.

  help
  reconfig
  showconfig <key name>
  genkey <key name>
  delkey <key name>
  clearkey
`, f.Name())
}

func Setup(name string) {
    f = flag.NewFlagSet(name, flag.ExitOnError)
    //flag
    f.Usage = Usage
    help.Usage = Usage
    help.Setup("help")
    reconfig.Setup("reconfig")
    showconfig.Setup("showconfig")
    genkey.Setup("genkey")
    delkey.Setup("delkey")
    clearkey.Setup("clearkey")
}

func Run(args []string) {
    f.Parse(args)
    subargs := f.Args()
    if len(subargs) == 0 {
        subargs = append(subargs, "")
    }
    switch subargs[0] {
    case "help":
        help.Run(subargs[1:])
    case "reconfig":
        reconfig.Run(subargs[1:])
    case "showconfig":
        showconfig.Run(subargs[1:])
    case "genkey":
        genkey.Run(subargs[1:])
    case "delkey":
        delkey.Run(subargs[1:])
    case "clearkey":
        clearkey.Run(subargs[1:])
    default:
        run(subargs)
    }
}

func run(subargs []string) {
    Usage()
}
