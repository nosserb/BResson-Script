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
        p = strings.ReplaceAll(p, "\"", "") // enlève tous les guillemets
        p = replaceVars(p)
        result += p
    }
    return result
}

func findFile(name string) string {
    var result string
    filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return nil
        }
        if !info.IsDir() && strings.EqualFold(info.Name(), name) {
            result = path
            return fmt.Errorf("trouvé") // stoppe le walk dès qu'on trouve
        }
        return nil
    })
    return result
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
        name = strings.Trim(name, "\"")
        name = strings.TrimSuffix(name, ")")
        path := findFile(name)
        if path == "" {
            fmt.Println("Fichier introuvable :", name)
        } else {
            fmt.Println("Fichier trouvé :", path)
        }
        return
    }

    if strings.Contains(line, "=") {
        parts := strings.Split(line, "=")
        key := strings.TrimSpace(parts[0])
        value := strings.TrimSpace(parts[1])
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
