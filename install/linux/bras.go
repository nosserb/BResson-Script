package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: bras <script.brs> [args...]")
		os.Exit(1)
	}

	scriptPath := os.Args[1]
	args := os.Args[2:]

	// Construire la commande pour exécuter bresson depuis le PATH
	cmdArgs := append([]string{scriptPath}, args...)
	cmd := exec.Command("bresson", cmdArgs...)

	// Rediriger stdin, stdout, stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Erreur d'exécution : %v\n", err)
		os.Exit(1)
	}
}
