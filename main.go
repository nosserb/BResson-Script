package main

import (
    "bufio"
    "fmt"
    "os"
    "path/filepath"
    "strings"
)

var variables = make(map[string]string)

func replaceVars(s string) string {
    for k, v := range variables {
        s = strings.ReplaceAll(s, k, v)
    }
    return s
}

func evalString(s string) string {
    parts := strings.Split(s, "+")
    result := ""
    for _, p := range parts {
        p = strings.TrimSpace(p)
        if strings.HasPrefix(p, "\"") && strings.HasSuffix(p, "\"") {
            p = strings.Trim(p, "\"")
        } else {
            p = replaceVars(p)
        }
        result += p
    }
    return result
}

func findFile(name string) string {
    var foundPath string
    filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return nil
        }
        if !info.IsDir() && strings.EqualFold(info.Name(), name) {
            foundPath = path
            return filepath.SkipDir
        }
        return nil
    })
    return foundPath
}

func runLine(line string) {
    line = strings.TrimSpace(line)
    if line == "" || strings.HasPrefix(line, "#") {
        return
    }

    if strings.HasPrefix(line, "bprint") {
        s := strings.TrimPrefix(line, "bprint(")
        s = strings.TrimSuffix(s, ")")
        fmt.Println(evalString(s))
        return
    }

    if strings.HasPrefix(line, "bfile") {
        name := strings.TrimPrefix(line, "bfile(")
        name = strings.TrimSpace(name)
        name = strings.ReplaceAll(name, "\"", "")
        name = strings.TrimSuffix(name, ")")
        name = strings.TrimSpace(name)
        path := findFile(name)
        if path == "" {
            fmt.Println("Fichier introuvable :", name)
        } else {
            fmt.Println(path)
        }
        return
    }

    if strings.HasPrefix(line, "bstartfile") {
        arg := strings.TrimPrefix(line, "bstartfile(")
        arg = strings.TrimSuffix(arg, ")")
        arg = strings.TrimSpace(arg)

        if strings.HasPrefix(arg, "'") && strings.HasSuffix(arg, "'") {
            arg = strings.Trim(arg, "'")
            arg = strings.TrimSpace(arg)
            if strings.HasPrefix(arg, "bfile") {
                inner := strings.TrimPrefix(arg, "bfile(")
                inner = strings.TrimSuffix(inner, ")")
                inner = evalString(inner)
                arg = findFile(inner)
            } else {
                arg = evalString(arg)
            }
        }

        if arg == "" {
            fmt.Println("Fichier introuvable")
            return
        }

        f, err := os.Open(arg)
        if err != nil {
            fmt.Println("Erreur ouverture :", err)
            return
        }
        defer f.Close()
        scanner := bufio.NewScanner(f)
        for scanner.Scan() {
            runLine(scanner.Text())
        }
        return
    }

    if strings.Contains(line, "=") {
        parts := strings.Split(line, "=")
        key := strings.TrimSpace(parts[0])
        value := strings.TrimSpace(parts[1])
        value = strings.Trim(value, "\"")
        variables[key] = value
        return
    }

    fmt.Println("Instruction inconnue :", line)
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: go run main.go fichier.brs")
        return
    }

    file := os.Args[1]
    f, err := os.Open(file)
    if err != nil {
        fmt.Println("Erreur ouverture :", err)
        return
    }
    defer f.Close()

    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        runLine(scanner.Text())
    }
}
