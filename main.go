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

// Variables globales pour stocker l'état du programme
var mesVariables = make(map[string]string)    // Stocke toutes les variables créées
var pileDeBlocs []bool                        // Pile pour gérer les if/while imbriqués
var mesTableaux = make(map[string][]string)   // Stocke les tableaux
var mesFonctions = make(map[string][]string)  // Stocke les fonctions définies
var mesTimers = make(map[string]time.Time)    // Stocke les timers
var couleurActuelle = ""                      // Couleur actuelle pour l'affichage
var argumentsCommande []string                // Arguments passés au programme
var loggingActive = false                     // Active/désactive le logging
var bufferLogs []string                       // Stocke tous les logs

// FONCTIONS DE LOGGING

// Ajoute un message au système de logs
func ajouterLog(message string) {
	if loggingActive {
		horodatage := time.Now().Format("15:04:05")
		entreeLog := fmt.Sprintf("[%s] %s", horodatage, message)
		bufferLogs = append(bufferLogs, entreeLog)
		fmt.Println(entreeLog)
	}
}

// Affiche du texte et l'ajoute aux logs si activé
func afficherTexte(texte string) {
	fmt.Println(texte)
	if loggingActive {
		ajouterLog(texte)
	}
}

// Retourne tous les logs sous forme de chaîne
func obtenirLogsCommeTexte() string {
	return strings.Join(bufferLogs, "\n")
}

// FONCTIONS UTILITAIRES 

// Extrait le nom de variable après "->"
func extraireNomVariable(ligne string) string {
	parties := strings.Split(ligne, "->")
	if len(parties) < 2 {
		return ""
	}
	return strings.TrimSpace(parties[1])
}

// Remplace toutes les variables dans une chaîne par leurs valeurs
func remplacerVariables(chaine string) string {
	for nom, valeur := range mesVariables {
		chaine = strings.ReplaceAll(chaine, nom, valeur)
	}
	return chaine
}

// Analyse une expression de texte avec concaténation
func analyserExpressionTexte(expression string) string {
	expression = strings.TrimSpace(expression)
	
	// Si pas de concaténation, retourner directement
	if !strings.Contains(expression, "+") {
		return remplacerVariables(strings.Trim(expression, "\""))
	}

	// Joindre des textes avec +
	parties := strings.Split(expression, "+")
	resultat := ""
	for _, partie := range parties {
		partie = strings.TrimSpace(partie)
		partie = strings.Trim(partie, "\"")
		partie = remplacerVariables(partie)
		resultat += partie
	}
	return resultat
}

// Cherche un fichier dans le répertoire courant et ses sous-dossiers
func chercherFichier(nom string) string {
	var resultat string
	filepath.Walk(".", func(chemin string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() && info.Name() == nom {
			resultat = chemin
			return filepath.SkipDir
		}
		return nil
	})
	return resultat
}

// evaluerExpression - Évalue une expression (comme bfile())
func evaluerExpression(expression string) string {
	expression = strings.TrimSpace(expression)
	if strings.HasPrefix(expression, "bfile") {
		contenu := strings.TrimPrefix(expression, "bfile")
		contenu = strings.Trim(contenu, "() ")
		contenu = remplacerVariables(contenu)
		return chercherFichier(contenu)
	}
	return remplacerVariables(expression)
}

// FONCTIONS DE CONTRÔLE DE FLUX 

// Vérifie si on doit exécuter la ligne courante (gestion des blocs if/while)
func doitExecuter() bool {
	if len(pileDeBlocs) == 0 {
		return true
	}
	for _, doitExec := range pileDeBlocs {
		if !doitExec {
			return false
		}
	}
	return true
}

// FONCTIONS MATHÉMATIQUES

// Évalue une expression mathématique simple
func calculerMath(expression string) (float64, error) {
	expression = strings.ReplaceAll(expression, " ", "")
	
	// Addition
	if strings.Contains(expression, "+") {
		parties := strings.Split(expression, "+")
		somme := 0.0
		for _, partie := range parties {
			valeur, err := strconv.ParseFloat(partie, 64)
			if err != nil {
				return 0, err
			}
			somme += valeur
		}
		return somme, nil
	}
	
	// Multiplication
	if strings.Contains(expression, "*") {
		parties := strings.Split(expression, "*")
		if len(parties) == 2 {
			a, err1 := strconv.ParseFloat(parties[0], 64)
			b, err2 := strconv.ParseFloat(parties[1], 64)
			if err1 == nil && err2 == nil {
				return a * b, nil
			}
		}
	}
	
	// Nombre simple
	return strconv.ParseFloat(expression, 64)
}

// AFFICHAGE 

// Affiche du texte avec la couleur définie
func afficherAvecCouleur(texte string) {
	switch couleurActuelle {
	case "rouge":
		fmt.Printf("\033[31m%s\033[0m\n", texte)
	case "vert":
		fmt.Printf("\033[32m%s\033[0m\n", texte)
	case "bleu":
		fmt.Printf("\033[34m%s\033[0m\n", texte)
	case "jaune":
		fmt.Printf("\033[33m%s\033[0m\n", texte)
	default:
		fmt.Println(texte)
	}
}

// FONCTION PRINCIPALE D'EXÉCUTION

// Analyse et exécute une ligne de code Bresson
func executerLigne(ligne string) {
	ligne = strings.TrimSpace(ligne)
	
	// Ignorer les lignes vides et les commentaires
	if ligne == "" || strings.HasPrefix(ligne, "#") {
		return
	}

	// COMMANDES DE LOGGING
	
	// Active le système de logging
	if strings.HasPrefix(ligne, "blog") {
		loggingActive = true
		if bufferLogs == nil {
			bufferLogs = []string{}
		}
		// Ne pas réinitialiser _blog si déjà existant
		if _, existe := mesVariables["_blog"]; !existe {
			mesVariables["_blog"] = ""
		}
		ajouterLog("Système de logging activé")
		fmt.Println("=== LOGGING ACTIVÉ ===")
		return
	}

	// Récupère tous les logs dans la variable _blog
	if strings.HasPrefix(ligne, "blogget") {
		mesVariables["_blog"] = obtenirLogsCommeTexte()
		ajouterLog("Logs récupérés dans la variable _blog")
		return
	}

	// Efface tous les logs
	if strings.HasPrefix(ligne, "blogclear") {
		bufferLogs = []string{}
		mesVariables["_blog"] = ""
		ajouterLog("Tous les logs ont été effacés")
		return
	}

	// CONTRÔLE DE FLUX
	
	// Ferme un bloc (if, while, etc.)
	if ligne == "|" {
		if len(pileDeBlocs) > 0 {
			pileDeBlocs = pileDeBlocs[:len(pileDeBlocs)-1]
		}
		return
	}

	// Début d'une boucle while
	if strings.HasPrefix(ligne, "bwhile") {
		condition := strings.TrimPrefix(ligne, "bwhile")
		condition = strings.Trim(condition, "() ")

		doitExecuterBloc := false
		if strings.Contains(condition, "<=") {
			parties := strings.Split(condition, "<=")
			if len(parties) == 2 {
				gauche := strings.TrimSpace(parties[0])
				droite := strings.TrimSpace(parties[1])
				valeurGauche, err1 := strconv.Atoi(remplacerVariables(gauche))
				valeurDroite, err2 := strconv.Atoi(remplacerVariables(droite))
				if err1 == nil && err2 == nil {
					doitExecuterBloc = valeurGauche <= valeurDroite
				}
			}
		}
		pileDeBlocs = append(pileDeBlocs, doitExecuterBloc)
		ajouterLog(fmt.Sprintf("Boucle while: condition=%s, exécution=%v", condition, doitExecuterBloc))
		return
	}

	// Début d'une condition if
	if strings.HasPrefix(ligne, "bif") {
		condition := strings.TrimPrefix(ligne, "bif")
		condition = strings.Trim(condition, "() ")

		doitExecuterBloc := false
		if strings.Contains(condition, "==") {
			parties := strings.Split(condition, "==")
			if len(parties) == 2 {
				gauche := strings.TrimSpace(parties[0])
				droite := strings.TrimSpace(parties[1])
				gauche = strings.Trim(gauche, "\"")
				droite = strings.Trim(droite, "\"")
				doitExecuterBloc = remplacerVariables(gauche) == remplacerVariables(droite)
			}
		}
		pileDeBlocs = append(pileDeBlocs, doitExecuterBloc)
		ajouterLog(fmt.Sprintf("Condition if: %s, résultat=%v", condition, doitExecuterBloc))
		return
	}

	//  Bloc else (inverse la condition précédente)
	if strings.HasPrefix(ligne, "belse") {
		if len(pileDeBlocs) > 0 {
			pileDeBlocs[len(pileDeBlocs)-1] = !pileDeBlocs[len(pileDeBlocs)-1]
		}
		ajouterLog("Bloc else activé")
		return
	}

	// Si on est dans un bloc qui ne doit pas s'exécuter, ignorer
	if !doitExecuter() {
		return
	}

	// COMMANDES D'ENTRÉE/SORTIE

	// Récupère un argument de la ligne de commande
	if strings.HasPrefix(ligne, "bard") {
		parametre := strings.TrimPrefix(ligne, "bard")
		parametre = strings.Trim(parametre, "() ")
		index, err := strconv.Atoi(parametre)
		if err == nil && index > 0 && index < len(argumentsCommande) {
			mesVariables["_bard"] = argumentsCommande[index]
			fmt.Printf("Argument %d: %s\n", index, argumentsCommande[index])
			ajouterLog(fmt.Sprintf("Argument %d récupéré: %s", index, argumentsCommande[index]))
		}
		return
	}

	// Demande une saisie utilisateur
	if strings.HasPrefix(ligne, "binput") {
		parties := strings.Split(ligne, "->")
		invite := strings.TrimSpace(parties[0][6:]) // Enlever "binput"
		invite = strings.Trim(invite, `()"`)
		nomVariable := ""
		if len(parties) > 1 {
			nomVariable = strings.TrimSpace(parties[1])
		}

		fmt.Print(invite + " ")
		var saisie string
		fmt.Scanln(&saisie)

		if nomVariable != "" {
			mesVariables[nomVariable] = saisie
			ajouterLog(fmt.Sprintf("Saisie utilisateur: %s = %s", nomVariable, saisie))
		}
		return
	}

	// Affiche du texte
	if strings.HasPrefix(ligne, "bprint") {
		texte := strings.TrimPrefix(ligne, "bprint")
		texte = strings.TrimSpace(texte)
		texte = strings.Trim(texte, `()"`)
		texte = analyserExpressionTexte(texte)

		// Remplacer les variables entre quotes
		for nom, valeur := range mesVariables {
			texte = strings.ReplaceAll(texte, "'"+nom+"'", valeur)
		}

		afficherAvecCouleur(texte)
		ajouterLog(fmt.Sprintf("Affichage: %s", texte))
		return
	}

	// COMMANDES MATHÉMATIQUES

	// Effectue un calcul mathématique
	if strings.HasPrefix(ligne, "bcalc") {
		expression := strings.TrimSpace(ligne[5:])
		resultat, _ := calculerMath(expression)
		nomVariable := extraireNomVariable(ligne)
		if nomVariable != "" {
			mesVariables[nomVariable] = fmt.Sprintf("%v", resultat)
			ajouterLog(fmt.Sprintf("Calcul: %s = %v", nomVariable, resultat))
		}
		return
	}

	// Génère un nombre aléatoire
	if strings.HasPrefix(ligne, "brand") {
		parametre := strings.TrimPrefix(ligne, "brand")
		parametre = strings.Trim(parametre, "() ")
		parties := strings.Fields(parametre)
		if len(parties) == 2 {
			min, err1 := strconv.Atoi(parties[0])
			max, err2 := strconv.Atoi(parties[1])
			if err1 == nil && err2 == nil {
				rand.Seed(time.Now().UnixNano())
				valeur := rand.Intn(max-min+1) + min
				mesVariables["_brand"] = strconv.Itoa(valeur)
				fmt.Printf("Nombre aléatoire généré: %d\n", valeur)
				ajouterLog(fmt.Sprintf("Nombre aléatoire: %d (entre %d et %d)", valeur, min, max))
			}
		}
		return
	}

	// COMMANDES DE FICHIERS 

	// Lit le contenu d'un fichier
	if strings.HasPrefix(ligne, "bread") {
		nomFichier := strings.TrimPrefix(ligne, "bread")
		nomFichier = strings.Trim(nomFichier, "()")
		nomFichier = strings.Trim(nomFichier, "\"")
		nomFichier = remplacerVariables(nomFichier)

		contenu, err := ioutil.ReadFile(nomFichier)
		if err == nil {
			mesVariables["_bread"] = string(contenu)
			fmt.Printf("Fichier lu: %s\n", nomFichier)
			ajouterLog(fmt.Sprintf("Fichier lu: %s (%d caractères)", nomFichier, len(contenu)))
		} else {
			fmt.Printf("Erreur de lecture: %s\n", err)
			ajouterLog(fmt.Sprintf("Erreur lecture fichier %s: %s", nomFichier, err))
		}
		return
	}

	// Écrit du contenu dans un fichier
	if strings.HasPrefix(ligne, "bwrite") {
		parametre := strings.TrimPrefix(ligne, "bwrite")
		parametre = strings.Trim(parametre, "() ")

		// Lire le texte entre guillemets
		parties := []string{}
		entreGuillemets := false
		actuel := ""
		for _, caractere := range parametre {
			if caractere == '"' {
				if entreGuillemets {
					parties = append(parties, actuel)
					actuel = ""
					entreGuillemets = false
				} else {
					entreGuillemets = true
				}
			} else if entreGuillemets {
				actuel += string(caractere)
			}
		}

		if len(parties) == 2 {
			nomFichier := remplacerVariables(parties[0])
			contenu := remplacerVariables(parties[1])
			err := ioutil.WriteFile(nomFichier, []byte(contenu), 0644)
			if err == nil {
				fmt.Printf("Fichier écrit: %s\n", nomFichier)
				ajouterLog(fmt.Sprintf("Fichier écrit: %s (%d caractères)", nomFichier, len(contenu)))
			} else {
				fmt.Printf("Erreur d'écriture: %s\n", err)
				ajouterLog(fmt.Sprintf("Erreur écriture fichier %s: %s", nomFichier, err))
			}
		}
		return
	}

	// Cherche un fichier
	if strings.HasPrefix(ligne, "bfile") {
		nom := strings.TrimPrefix(ligne, "bfile")
		nom = strings.Trim(nom, "()\\ ")
		chemin := chercherFichier(remplacerVariables(nom))
		if chemin == "" {
			fmt.Println("Fichier introuvable :", nom)
			ajouterLog(fmt.Sprintf("Fichier introuvable: %s", nom))
		} else {
			fmt.Println("Fichier trouvé :", chemin)
			ajouterLog(fmt.Sprintf("Fichier trouvé: %s", chemin))
		}
		return
	}

	// Renomme un fichier
	if strings.HasPrefix(ligne, "brename") {
		parametre := strings.TrimPrefix(ligne, "brename")
		parametre = strings.Trim(parametre, "() ")
		
		// Lire le texte entre guillemets
		parties := []string{}
		entreGuillemets := false
		actuel := ""
		for _, caractere := range parametre {
			if caractere == '"' {
				if entreGuillemets {
					parties = append(parties, actuel)
					actuel = ""
					entreGuillemets = false
				} else {
					entreGuillemets = true
				}
			} else if entreGuillemets {
				actuel += string(caractere)
			}
		}

		if len(parties) == 2 {
			ancienNom := remplacerVariables(parties[0])
			nouveauNom := remplacerVariables(parties[1])
			err := os.Rename(ancienNom, nouveauNom)
			if err != nil {
				fmt.Println("Erreur de renommage:", err)
				ajouterLog(fmt.Sprintf("Erreur renommage %s -> %s: %s", ancienNom, nouveauNom, err))
			} else {
				fmt.Println("Renommé:", ancienNom, "->", nouveauNom)
				ajouterLog(fmt.Sprintf("Fichier renommé: %s -> %s", ancienNom, nouveauNom))
			}
		} else {
			fmt.Println("Usage: brename(\"ancien.txt\" \"nouveau.txt\")")
		}
		return
	}

	// COMMANDES DE TEMPS

	// Faire une pause
	if strings.HasPrefix(ligne, "bsleep") {
		parametre := strings.TrimPrefix(ligne, "bsleep")
		parametre = strings.Trim(parametre, "() ")
		secondes, err := strconv.Atoi(remplacerVariables(parametre))
		if err == nil {
			fmt.Printf("Attente de %d secondes...\n", secondes)
			ajouterLog(fmt.Sprintf("Pause: %d secondes", secondes))
			time.Sleep(time.Duration(secondes) * time.Second)
		}
		return
	}

	// Démarre un timer
	if strings.HasPrefix(ligne, "btimer") {
		nom := strings.TrimPrefix(ligne, "btimer")
		nom = strings.Trim(nom, "()")
		nom = strings.Trim(nom, "\"")
		nom = remplacerVariables(nom)
		mesTimers[nom] = time.Now()
		fmt.Printf("Timer '%s' démarré\n", nom)
		ajouterLog(fmt.Sprintf("Timer '%s' démarré", nom))
		return
	}

	// Arrête un timer et affiche la durée
	if strings.HasPrefix(ligne, "bendtimer") {
		nom := strings.TrimPrefix(ligne, "bendtimer")
		nom = strings.Trim(nom, "()")
		nom = strings.Trim(nom, "\"")
		nom = remplacerVariables(nom)
		if heureDebut, existe := mesTimers[nom]; existe {
			duree := time.Since(heureDebut)
			fmt.Printf("Timer '%s': %v\n", nom, duree)
			ajouterLog(fmt.Sprintf("Timer '%s' terminé: %v", nom, duree))
			delete(mesTimers, nom)
		}
		return
	}

	// Récupère la date/heure actuelle
	if strings.HasPrefix(ligne, "btime") {
		parametre := strings.TrimPrefix(ligne, "btime")
		parametre = strings.Trim(parametre, "()\\ ")
		maintenant := time.Now()
		parties := strings.Fields(parametre)
		resultat := ""
		
		for i, partie := range parties {
			switch partie {
			case "j":
				resultat += fmt.Sprintf("%02d", maintenant.Day())
			case "m":
				resultat += fmt.Sprintf("%02d", maintenant.Month())
			case "a":
				resultat += fmt.Sprintf("%d", maintenant.Year())
			case "h":
				resultat += fmt.Sprintf("%02d", maintenant.Hour())
			case "min":
				resultat += fmt.Sprintf("%02d", maintenant.Minute())
			case "s":
				resultat += fmt.Sprintf("%02d", maintenant.Second())
			}
			if i < len(parties)-1 {
				resultat += " "
			}
		}
		mesVariables["_btime"] = resultat
		fmt.Println("Date/Heure:", resultat)
		ajouterLog(fmt.Sprintf("Date/heure récupérée: %s", resultat))
		return
	}

	// COMMANDES D'AFFICHAGE

	// Change la couleur d'affichage
	if strings.HasPrefix(ligne, "bcolor") {
		couleur := strings.TrimPrefix(ligne, "bcolor")
		couleur = strings.Trim(couleur, "()")
		couleur = strings.Trim(couleur, "\"")
		couleurActuelle = remplacerVariables(couleur)
		fmt.Printf("Couleur changée: %s\n", couleurActuelle)
		ajouterLog(fmt.Sprintf("Couleur changée: %s", couleurActuelle))
		return
	}

	// Exécute un fichier (simulation)
	if strings.HasPrefix(ligne, "bstartfile") {
		nom := strings.TrimPrefix(ligne, "bstartfile")
		nom = strings.TrimSpace(nom)
		nom = strings.Trim(nom, "()")
		nom = strings.Trim(nom, "\"")
		nom = remplacerVariables(nom)
		fmt.Println("Exécution de:", nom)
		ajouterLog(fmt.Sprintf("Exécution: %s", nom))
		return
	}

	// GESTION DES VARIABLES

	// Assignation de variable ( = mais pas ==)
	if strings.Contains(ligne, "=") && !strings.Contains(ligne, "==") {
		parties := strings.SplitN(ligne, "=", 2)
		nom := strings.TrimSpace(parties[0])
		valeur := strings.TrimSpace(parties[1])

		// Variables spéciales (résultats de commandes)
		if valeur == "_binput" || valeur == "_bcalc" || valeur == "_bread" || valeur == "_bard" || valeur == "_blog" {
			if val, existe := mesVariables[valeur]; existe {
				mesVariables[nom] = val
				fmt.Printf("Variable %s = %s\n", nom, mesVariables[nom])
				ajouterLog(fmt.Sprintf("Variable assignée: %s = %s", nom, mesVariables[nom]))
			}
		} else {
			// Variable normale
			valeur = strings.Trim(valeur, "\"")
			mesVariables[nom] = remplacerVariables(valeur)
			fmt.Printf("Variable %s = %s\n", nom, mesVariables[nom])
			ajouterLog(fmt.Sprintf("Variable créée: %s = %s", nom, mesVariables[nom]))
		}
		return
	}
}

// FONCTION PRINCIPALE

// Point d'entrée du programme
func main() {
	// Vérifier qu'un fichier a été fourni
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go fichier.brs [arguments...]")
		return
	}

	// Sauvegarder les arguments de la ligne de commande
	argumentsCommande = os.Args

	// Ouvrir le fichier script
	fichier := os.Args[1]
	f, err := os.Open(fichier)
	if err != nil {
		fmt.Println("Erreur d'ouverture du fichier:", err)
		return
	}
	defer f.Close()

	// Exécuter le script ligne par ligne
	fmt.Println("=== Exécution du script Bresson ===")
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		executerLigne(scanner.Text())
	}
	fmt.Println("=== Fin d'exécution ===")
}
