package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var variables = make(map[string]string)
var blockStack []bool
var arrays = make(map[string][]string)
var functions = make(map[string][]string)
var timers = make(map[string]time.Time)
var currentColor = ""

var cmdArgs []string

func extractVarName(line string) string {
	parts := strings.Split(line, "->")
	if len(parts) < 2 {
		return ""
	}
	return strings.TrimSpace(parts[1])
}

func replaceVars(s string) string {
	for k, v := range variables {
		s = strings.ReplaceAll(s, k, v)
	}
	return s
}

func parseStringExpression(expr string) string {
	expr = strings.TrimSpace(expr)
	if !strings.Contains(expr, "+") {
		return replaceVars(strings.Trim(expr, "\""))
	}

	parts := strings.Split(expr, "+")
	result := ""
	for _, part := range parts {
		part = strings.TrimSpace(part)
		part = strings.Trim(part, "\"")
		part = replaceVars(part)
		result += part
	}
	return result
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

func evalExpression(expr string) string {
	expr = strings.TrimSpace(expr)
	if strings.HasPrefix(expr, "bfile") {
		inner := strings.TrimPrefix(expr, "bfile")
		inner = strings.Trim(inner, "() ")
		inner = replaceVars(inner)
		return findFile(inner)
	}
	return replaceVars(expr)
}

func shouldExecute() bool {
	if len(blockStack) == 0 {
		return true
	}
	for _, shouldExec := range blockStack {
		if !shouldExec {
			return false
		}
	}
	return true
}

func evaluateMath(expr string) (float64, error) {
	expr = strings.ReplaceAll(expr, " ", "")
	if strings.Contains(expr, "+") {
		parts := strings.Split(expr, "+")
		sum := 0.0
		for _, part := range parts {
			val, err := strconv.ParseFloat(part, 64)
			if err != nil {
				return 0, err
			}
			sum += val
		}
		return sum, nil
	}
	if strings.Contains(expr, "*") {
		parts := strings.Split(expr, "*")
		if len(parts) == 2 {
			a, err1 := strconv.ParseFloat(parts[0], 64)
			b, err2 := strconv.ParseFloat(parts[1], 64)
			if err1 == nil && err2 == nil {
				return a * b, nil
			}
		}
	}
	return strconv.ParseFloat(expr, 64)
}

func printWithColor(text string) {
	switch currentColor {
	case "rouge":
		fmt.Printf("\033[31m%s\033[0m\n", text)
	case "vert":
		fmt.Printf("\033[32m%s\033[0m\n", text)
	case "bleu":
		fmt.Printf("\033[34m%s\033[0m\n", text)
	case "jaune":
		fmt.Printf("\033[33m%s\033[0m\n", text)
	default:
		fmt.Println(text)
	}
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

	if strings.HasPrefix(line, "bwhile") {
		cond := strings.TrimPrefix(line, "bwhile")
		cond = strings.Trim(cond, "() ")

		shouldExecBlock := false
		if strings.Contains(cond, "<=") {
			parts := strings.Split(cond, "<=")
			if len(parts) == 2 {
				left := strings.TrimSpace(parts[0])
				right := strings.TrimSpace(parts[1])
				leftVal, err1 := strconv.Atoi(replaceVars(left))
				rightVal, err2 := strconv.Atoi(replaceVars(right))
				if err1 == nil && err2 == nil {
					shouldExecBlock = leftVal <= rightVal
				}
			}
		}
		blockStack = append(blockStack, shouldExecBlock)
		return
	}

	if strings.HasPrefix(line, "bif") {
		cond := strings.TrimPrefix(line, "bif")
		cond = strings.Trim(cond, "() ")

		shouldExecBlock := false
		if strings.Contains(cond, "==") {
			parts := strings.Split(cond, "==")
			if len(parts) == 2 {
				left := strings.TrimSpace(parts[0])
				right := strings.TrimSpace(parts[1])
				left = strings.Trim(left, "\"")
				right = strings.Trim(right, "\"")
				shouldExecBlock = replaceVars(left) == replaceVars(right)
			}
		}
		blockStack = append(blockStack, shouldExecBlock)
		return
	}

	if strings.HasPrefix(line, "belse") {
		if len(blockStack) > 0 {
			blockStack[len(blockStack)-1] = !blockStack[len(blockStack)-1]
		}
		return
	}

	if !shouldExecute() {
		return
	}

	if strings.HasPrefix(line, "bard") {
		s := strings.TrimPrefix(line, "bard")
		s = strings.Trim(s, "() ")
		index, err := strconv.Atoi(s)
		if err == nil && index > 0 && index < len(cmdArgs) {
			variables["_bard"] = cmdArgs[index]
			fmt.Printf("Argument %d: %s\n", index, cmdArgs[index])
		}
		return
	}

	if strings.HasPrefix(line, "binput") {
		parts := strings.Split(line, "->")
		prompt := strings.TrimSpace(parts[0][6:])
		prompt = strings.Trim(prompt, `()"`)
		varName := ""
		if len(parts) > 1 {
			varName = strings.TrimSpace(parts[1])
		}

		fmt.Print(prompt + " ")
		var input string
		fmt.Scanln(&input)

		if varName != "" {
			variables[varName] = input
		}
	}

	if strings.HasPrefix(line, "bcalc") {
		expr := strings.TrimSpace(line[5:])
		result, _ := evaluateMath(expr)
		varName := extractVarName(line)
		if varName != "" {
			variables[varName] = fmt.Sprintf("%v", result)
		}
		return
	}

	if strings.HasPrefix(line, "bread") {
		filename := strings.TrimPrefix(line, "bread")
		filename = strings.Trim(filename, "()")
		filename = strings.Trim(filename, "\"")
		filename = replaceVars(filename)

		content, err := ioutil.ReadFile(filename)
		if err == nil {
			variables["_bread"] = string(content)
			fmt.Printf("Lu fichier: %s\n", filename)
		} else {
			fmt.Printf("Erreur lecture: %s\n", err)
		}
		return
	}

	if strings.HasPrefix(line, "bwrite") {
		s := strings.TrimPrefix(line, "bwrite")
		s = strings.Trim(s, "() ")

		parts := []string{}
		inQuotes := false
		current := ""
		for _, char := range s {
			if char == '"' {
				if inQuotes {
					parts = append(parts, current)
					current = ""
					inQuotes = false
				} else {
					inQuotes = true
				}
			} else if inQuotes {
				current += string(char)
			}
		}

		if len(parts) == 2 {
			filename := replaceVars(parts[0])
			content := replaceVars(parts[1])
			err := ioutil.WriteFile(filename, []byte(content), 0644)
			if err == nil {
				fmt.Printf("Écrit fichier: %s\n", filename)
			} else {
				fmt.Printf("Erreur écriture: %s\n", err)
			}
		}
		return
	}

	if strings.HasPrefix(line, "bsleep") {
		s := strings.TrimPrefix(line, "bsleep")
		s = strings.Trim(s, "() ")
		seconds, err := strconv.Atoi(replaceVars(s))
		if err == nil {
			fmt.Printf("Attente %d secondes...\n", seconds)
			time.Sleep(time.Duration(seconds) * time.Second)
		}
		return
	}

	if strings.HasPrefix(line, "btimer") {
		name := strings.TrimPrefix(line, "btimer")
		name = strings.Trim(name, "()")
		name = strings.Trim(name, "\"")
		name = replaceVars(name)
		timers[name] = time.Now()
		fmt.Printf("Timer '%s' démarré\n", name)
		return
	}

	if strings.HasPrefix(line, "bendtimer") {
		name := strings.TrimPrefix(line, "bendtimer")
		name = strings.Trim(name, "()")
		name = strings.Trim(name, "\"")
		name = replaceVars(name)
		if startTime, exists := timers[name]; exists {
			duration := time.Since(startTime)
			fmt.Printf("Timer '%s': %v\n", name, duration)
			delete(timers, name)
		}
		return
	}

	if strings.HasPrefix(line, "bcolor") {
		color := strings.TrimPrefix(line, "bcolor")
		color = strings.Trim(color, "()")
		color = strings.Trim(color, "\"")
		currentColor = replaceVars(color)
		fmt.Printf("Couleur changée: %s\n", currentColor)
		return
	}

	if strings.HasPrefix(line, "bprint") {
		s := strings.TrimPrefix(line, "bprint")
		s = strings.TrimSpace(s)
		s = strings.Trim(s, `()"`)
		s = parseStringExpression(s)

		for k, v := range variables {
			s = strings.ReplaceAll(s, "'"+k+"'", v)
		}

		printWithColor(s)
		return
	}

	if strings.HasPrefix(line, "bfile") {
		name := strings.TrimPrefix(line, "bfile")
		name = strings.Trim(name, "()\\ ")
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
		name = strings.TrimSpace(name)
		name = strings.Trim(name, "()")
		name = strings.Trim(name, "\"")
		name = replaceVars(name)
		fmt.Println("Exécution :", name)
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
				fmt.Printf("Random généré: %d\n", val)
			}
		}
		return
	}

	if strings.HasPrefix(line, "btime") {
		s := strings.TrimPrefix(line, "btime")
		s = strings.Trim(s, "()\\ ")
		now := time.Now()
		parts := strings.Fields(s)
		res := ""
		for i, p := range parts {
			switch p {
			case "j":
				res += fmt.Sprintf("%02d", now.Day())
			case "m":
				res += fmt.Sprintf("%02d", now.Month())
			case "a":
				res += fmt.Sprintf("%d", now.Year())
			case "h":
				res += fmt.Sprintf("%02d", now.Hour())
			case "min":
				res += fmt.Sprintf("%02d", now.Minute())
			case "s":
				res += fmt.Sprintf("%02d", now.Second())
			}
			if i < len(parts)-1 {
				res += " "
			}
		}
		variables["_btime"] = res
		fmt.Println("Temps:", res)
		return
	}

	if strings.HasPrefix(line, "brename") {
		s := strings.TrimPrefix(line, "brename")
		s = strings.Trim(s, "() ")
		parts := []string{}
		inQuotes := false
		current := ""
		for _, char := range s {
			if char == '"' {
				if inQuotes {
					parts = append(parts, current)
					current = ""
					inQuotes = false
				} else {
					inQuotes = true
				}
			} else if inQuotes {
				current += string(char)
			}
		}

		if len(parts) == 2 {
			oldName := replaceVars(parts[0])
			newName := replaceVars(parts[1])
			err := os.Rename(oldName, newName)
			if err != nil {
				fmt.Println("Erreur renommage:", err)
			} else {
				fmt.Println("Renommé:", oldName, "->", newName)
			}
		} else {
			fmt.Println("Usage: brename(\"ancien.txt\" \"nouveau.txt\")")
		}
		return
	}

	if strings.Contains(line, "=") && !strings.Contains(line, "==") {
		parts := strings.SplitN(line, "=", 2)
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if value == "_binput" || value == "_bcalc" || value == "_bread" || value == "_bard" {
			if val, exists := variables[value]; exists {
				variables[key] = val
				fmt.Printf("Variable %s = %s\n", key, variables[key])
			}
		} else {
			value = strings.Trim(value, "\"")
			variables[key] = replaceVars(value)
			fmt.Printf("Variable %s = %s\n", key, variables[key])
		}
		return
	}

}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go fichier.brs [args...]")
		return
	}

	cmdArgs = os.Args

	file := os.Args[1]
	f, err := os.Open(file)
	if err != nil {
		fmt.Println("Erreur ouverture :", err)
		return
	}
	defer f.Close()

	fmt.Println("=== Exécution du script Bresson ===")
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		runLine(scanner.Text())
	}
	fmt.Println("=== Fin d'exécution ===")
}
