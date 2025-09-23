package main

import (
    "os"
    "os/exec"
    "fmt"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: bras fichier.brs")
        return
    }
    file := os.Args[1]
    cmd := exec.Command("go", "run", "/Users/admin/BResson-Script/main.go", file)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    err := cmd.Run()
    if err != nil {
        fmt.Println("Erreur :", err)
    }
}
