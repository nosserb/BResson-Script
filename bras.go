package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func findInterpreterDir(dir string) (string, bool) {
	candidates := []string{
		"bresson.exe", "bresson_interpreter.exe", "interpreter.exe",
		"bresson", "bresson_interpreter", "interpreter",
	}
	for _, c := range candidates {
		p := filepath.Join(dir, c)
		if fileExists(p) {
			return p, true
		}
	}
	mainGo := filepath.Join(dir, "main.go")
	if fileExists(mainGo) {
		return mainGo, false
	}
	parentMain := filepath.Join(dir, "..", "main.go")
	if fileExists(parentMain) {
		return parentMain, false
	}
	return "", false
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: bras <script.brs> [args...]")
		os.Exit(1)
	}

	execPath, err := os.Executable()
	if err != nil {
		fmt.Println("Erreur:", err)
		os.Exit(1)
	}
	dir := filepath.Dir(execPath)

	interp, isBinary := findInterpreterDir(dir)
	if interp == "" {
		cwd, _ := os.Getwd()
		interp = filepath.Join(cwd, "main.go")
		isBinary = false
		if !fileExists(interp) {
			fmt.Println("Impossible de trouver main.go ni binaire d'interpréteur dans le dossier de bras ni dans le dossier courant.")
			os.Exit(1)
		}
	}

	args := os.Args[1:]
	var cmd *exec.Cmd
	if isBinary {
		cmd = exec.Command(interp, args...)
	} else {
		cmdArgs := append([]string{"run", interp}, args...)
		cmd = exec.Command("go", cmdArgs...)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		fmt.Println("Erreur :", err)
		os.Exit(1)
	}
}
