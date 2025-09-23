package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "strings"
)

var variables = make(map[string]string)

func executeLine(line string) {
    line = strings.TrimSpace(line)
    if line == "" || strings.HasPrefix(line, "#") {
        return 
    }

    if strings.HasPrefix(line, "bprint") {
        content := strings.TrimPrefix(line, "bprint(")
        content = strings.TrimSuffix(content, ")")
        content = strings.Trim(content, "\"")
        fmt.Println(content)
        return
    }

    if strings.Contains(line, "=") {
        parts := strings.Split(line, "=")
        key := strings.TrimSpace(parts[0])
        value := strings.TrimSpace(parts[1])
        variables[key] = value
        return
    }

    fmt.Println("Instruction non reconnue:", line)
}

func main() {
    if len(os.Args) < 2 {
        log.Fatal("Usage: go run main.go fichier.brs")
    }

    file := os.Args[1]
    f, err := os.Open(file)
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()

    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        executeLine(scanner.Text())
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
}
