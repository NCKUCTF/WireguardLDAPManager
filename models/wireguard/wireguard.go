package wireguard

import (
    "fmt"
    "log"
    "io"
    "os"
    "os/exec"
    "bufio"
    "strings"
    "strconv"
    "path/filepath"
    "regexp"

    "github.com/google/uuid"

    "WireguardManager/utils/config"
)

func init() {
    _, err := exec.LookPath("wg")
    if err != nil {
        log.Fatalln("Please install wireguard.")
    }
}

func GetAllName() []string {
    out, _ := exec.Command("wg", "show", "interfaces").Output()
    return strings.Split(strings.TrimSpace(string(out)), " ")
}

func ReadName() string {
    fmt.Println("Please select your server name.")
    names := GetAllName()
    for i := 0; i < len(names); i++ {
        fmt.Printf("%d: %s\n", i+1, names[i])
    }
    indexstr := ""
    var index int64
    var err error
    for index, err = strconv.ParseInt(indexstr, 10, 64); err != nil || int(index) < 1 || int(index) > len(names); index, err = strconv.ParseInt(indexstr, 10, 64) {
        fmt.Printf("> ")
        fmt.Scanln(&indexstr)
    }
    return names[index - 1]
}

func GetConfig(name string) (conf string, data map[string]string) {
    f, err := os.Open(filepath.Join(config.WGpath, fmt.Sprintf("%s.conf", name)))
    if err != nil {
        log.Fatalln(err)
    }
    defer f.Close()

    scanner := bufio.NewScanner(f)
    conf = ""
    data = make(map[string]string)
    ininterface := false
    rmcomment := regexp.MustCompile(`^#\s*`)
    spliteq := regexp.MustCompile(`\s*=\s*`)
    matchend := regexp.MustCompile(`^# BEGIN .*$`)
    for scanner.Scan() {
        now := scanner.Text()
        if now == "" {
            continue
        } else if now == "[Interface]" {
            ininterface = true
        } else if !ininterface {
            nowvar := rmcomment.ReplaceAllString(now, "")
            if spliteq.MatchString(nowvar) {
                data[spliteq.Split(nowvar, 2)[0]] = spliteq.Split(nowvar, 2)[1]
            }
        } else if matchend.MatchString(now) {
            break
        } else {
            if spliteq.MatchString(now) {
                data[spliteq.Split(now, 2)[0]] = spliteq.Split(now, 2)[1]
            }
        }
        conf += fmt.Sprintf("\n%s", now)
    }

    if err := scanner.Err(); err != nil {
        log.Fatalln(err)
    }

    return
}

func SetConfig(name, conf string) {
    f, err := os.OpenFile(filepath.Join(config.WGpath, fmt.Sprintf("%s.conf", name)), os.O_RDWR, 0600)
    if err != nil {
        log.Fatalln(err)
    }
    _, err = f.WriteString(conf)
    if err != nil {
        log.Fatalln(err)
    }
    f.Close()
    out, _ := exec.Command("wg-quick", "strip", name).Output()

    tmpname := uuid.New().String()
    f, err = os.Create(fmt.Sprintf("/tmp/%s", tmpname))
    if err != nil {
        log.Fatalln(err)
    }
    _, err = f.Write(out)
    if err != nil {
        log.Fatalln(err)
    }
    f.Close()
    defer os.Remove(fmt.Sprintf("/tmp/%s", tmpname))
    exec.Command("wg", "syncconf", name, fmt.Sprintf("/tmp/%s", tmpname)).Run()
}

func Pubkey(privkey string) string {
    cmd := exec.Command("wg", "pubkey")
    stdin, err := cmd.StdinPipe()
    if err != nil {
        log.Fatalln(err)
    }
    stdout, err := cmd.StdoutPipe()
    if err != nil {
        log.Fatalln(err)
    }
    io.WriteString(stdin, privkey)
    stdin.Close()
    if err = cmd.Start(); err != nil {
        log.Fatalln(err)
    }
    out, _ := io.ReadAll(stdout)
    if err = cmd.Wait(); err != nil {
        log.Fatalln(err)
    }
    return strings.TrimSpace(string(out))
}
