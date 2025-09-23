package main

import (
    "bufio"
    "fmt"
    "math/rand"
    "os"
    "strconv"
    "strings"
    "time"
)

var variables = make(map[string]string)

func replaceVars(s string) string {
    for k, v := range variables {
        s = strings.ReplaceAll(s, k, v)
    }
    return s
}

func runLine(line string) {
    line = strings.TrimSpace(line)
    if line == "" || strings.HasPrefix(line, "#") {
        return
    }

    if strings.HasPrefix(line, "bprint") {
        s := strings.TrimPrefix(line, "bprint(")
        s = strings.TrimSuffix(s, ")")
        s = strings.Trim(s, "\"")
        s = replaceVars(s)
        fmt.Println(s)
        return
    }

    if strings.HasPrefix(line, "bfile") {
        name := strings.TrimPrefix(line, "bfile(")
        name = strings.Trim(name, "\"")
        fmt.Println("Trouvé :", name)
        return
    }

    if strings.HasPrefix(line, "bstartfile") {
        name := strings.TrimPrefix(line, "bstartfile(")
        name = strings.Trim(name, "\"")
        fmt.Println("Exécution :", name)
        return
    }

    if strings.HasPrefix(line, "bif") {
        cond := strings.TrimPrefix(line, "bif(")
        cond = strings.TrimSuffix(cond, ")")
        cond = strings.TrimSpace(cond)
        parts := strings.Split(cond, "==")
        if len(parts) == 2 {
            left := strings.TrimSpace(parts[0])
            right := strings.TrimSpace(parts[1])
            if variables[left] == right {
                variables["_bif_true"] = "1"
            } else {
                variables["_bif_true"] = "0"
            }
        }
        return
    }

    if strings.HasPrefix(line, "belse") {
        if variables["_bif_true"] == "0" {
            variables["_bif_true"] = "1"
        } else {
            variables["_bif_true"] = "0"
        }
        return
    }

    if strings.HasPrefix(line, "brand") {
        s := strings.TrimPrefix(line, "brand(")
        s = strings.TrimSuffix(s, ")")
        s = strings.TrimSpace(s)
        parts := strings.Split(s, " ")
        if len(parts) == 2 {
            min, err1 := strconv.Atoi(parts[0])
            max, err2 := strconv.Atoi(parts[1])
            if err1 == nil && err2 == nil {
                rand.Seed(time.Now().UnixNano())
                val := rand.Intn(max-min+1) + min
                variables["_brand"] = strconv.Itoa(val)
                fmt.Println("Random :", val)
            }
        }
        return
    }

    if strings.HasPrefix(line, "btime") {
        s := strings.TrimPrefix(line, "btime(")
        s = strings.TrimSuffix(line, ")")
        s = strings.Trim(s, "\"")
        now := time.Now()
        parts := strings.Split(s, " ")
        res := ""
        for _, p := range parts {
            switch p {
            case "j":
                res += fmt.Sprintf("%02d ", now.Day())
            case "m":
                res += fmt.Sprintf("%02d ", now.Month())
            case "a":
                res += fmt.Sprintf("%d ", now.Year())
            case "h":
                res += fmt.Sprintf("%02d ", now.Hour())
            case "min":
                res += fmt.Sprintf("%02d ", now.Minute())
            case "s":
                res += fmt.Sprintf("%02d ", now.Second())
            }
        }
        res = strings.TrimSpace(res)
        variables["_btime"] = res
        fmt.Println(res)
        return
    }

    if strings.HasPrefix(line, "brename") {
        s := strings.TrimPrefix(line, "brename(")
        s = strings.TrimSuffix(s, ")")
        s = strings.Trim(s, "\"")
        parts := strings.Split(s, "\" \"")
        if len(parts) == 2 {
            oldName := strings.TrimSpace(parts[0])
            newName := strings.TrimSpace(parts[1])
            err := os.Rename(oldName, newName)
            if err != nil {
                fmt.Println("Erreur :", err)
            } else {
                fmt.Println("Renommé :", oldName, "->", newName)
            }
        } else {
            fmt.Println("Usage: brename(\"ancien.txt\" \"nouveau.txt\")")
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
