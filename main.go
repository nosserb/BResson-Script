package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var variables = make(map[string]string)
var blockStack []string

func replaceVars(s string) string {
	for k, v := range variables {
		s = strings.ReplaceAll(s, k, v)
	}
	return s
}

func findFile(name string) string {
	var result string
	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() && info.Name() == name {
			result = path
			return filepath.SkipDir
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

	if line == "|" {
		if len(blockStack) > 0 {
			blockStack = blockStack[:len(blockStack)-1]
		}
		return
	}

	if strings.HasPrefix(line, "bprint") {
		s := strings.TrimPrefix(line, "bprint")
		s = strings.Trim(s, "() ")
		s = replaceVars(s)
		fmt.Println(s)
		return
	}

	if strings.HasPrefix(line, "bfile") {
		name := strings.TrimPrefix(line, "bfile")
		name = strings.Trim(name, "()\" ")
		path := findFile(replaceVars(name))
		if path == "" {
			fmt.Println("Fichier introuvable :", name)
		} else {
			fmt.Println("Fichier trouvé :", path)
		}
		return
	}

	if strings.HasPrefix(line, "bstartfile") {
		name := strings.TrimPrefix(line, "bstartfile")
		name = strings.Trim(name, "()\" ")
		name = replaceVars(name)
		fmt.Println("Exécution :", name)
		return
	}

	if strings.HasPrefix(line, "bif") {
		cond := strings.TrimPrefix(line, "bif")
		cond = strings.Trim(cond, "() ")
		parts := strings.Split(cond, "==")
		if len(parts) == 2 {
			left := strings.TrimSpace(parts[0])
			right := strings.TrimSpace(parts[1])
			if replaceVars(left) == replaceVars(right) {
				variables["_bif_true"] = "1"
			} else {
				variables["_bif_true"] = "0"
			}
		}
		blockStack = append(blockStack, "if")
		return
	}

	if strings.HasPrefix(line, "belse") {
		if variables["_bif_true"] == "0" {
			variables["_bif_true"] = "1"
		} else {
			variables["_bif_true"] = "0"
		}
		blockStack = append(blockStack, "else")
		return
	}

	if strings.HasPrefix(line, "brand") {
		s := strings.TrimPrefix(line, "brand")
		s = strings.Trim(s, "() ")
		parts := strings.Fields(s)
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
		s := strings.TrimPrefix(line, "btime")
		s = strings.Trim(s, "()\" ")
		now := time.Now()
		parts := strings.Fields(s)
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
		s := strings.TrimPrefix(line, "brename")
		s = strings.Trim(s, "() ")
		parts := strings.Split(s, "\" \"")
		if len(parts) == 2 {
			oldName := replaceVars(strings.Trim(parts[0], "\""))
			newName := replaceVars(strings.Trim(parts[1], "\""))
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
		parts := strings.SplitN(line, "=", 2)
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
