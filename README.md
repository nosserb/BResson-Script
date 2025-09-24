# Bresson Script (BRS)

Bresson Script (abrégé **BRS**) est un langage de script léger conçu pour faciliter l’écriture rapide de scripts en ligne de commande. Inspiré par un projet personnel, il combine la simplicité d’utilisation avec la puissance de Go pour exécuter des scripts, manipuler des fichiers, générer des nombres aléatoires et gérer le temps.

---

## Installation

1. Clone le dépôt GitHub :

```
git clone https://github.com/nosserb/BResson-Script.git
cd Bresson-Script
```
2. sur linux :

### Execute les commandes :

## exemple
```
bprint("Bienvenue dans le jeu de devinette !")
binput("Comment tu t'appelles ?") -> 'nom'

bprint("Salut " + 'nom' + "! Je pense à un nombre entre 1 et 10...")

bcalc(rand(1,10)) -> 'nombreMystere'

binput("Devine le nombre :") -> 'guess'

bif('guess' == 'nombreMystere')
    bprint("Bravo ! Tu as trouvé ! 🎉")
belse
    bprint("Dommage ! Le nombre était "  'nombreMystere')
|
```

# liste de fonction

### - bfile() 

```
bfile("nom_fichier")
```
vas rechercher le fichier depuis la racine ou est executer la commande

---

### - bprint()

```
bprint("Bonjour le monde!")
bprint("Valeur: " + maVariable)
```
Affiche du texte à l'écran avec support des couleurs

---

### - binput() -> variable

```
binput("Entrez votre nom: ") -> nom
bprint("Bonjour " + nom)
```
Demande une saisie utilisateur et stocke le résultat dans une variable

---

### - variable = valeur

```
nom = "Jean"
age = "25"
```
Assigne une valeur à une variable

---

### - bcalc() -> variable

```
bcalc(5 + 3) -> resultat
bcalc(10 * 2) -> produit
```
Effectue des calculs mathématiques (addition, multiplication)

---

### - bread()

```
bread("config.txt")
contenu = _bread
```
Lit le contenu d'un fichier et le stocke dans `_bread` (variable)

---

### - bwrite()

```
bwrite("output.txt", "Mon contenu")
```
Écrit du contenu dans un fichier

---

### - brename()

```
brename("old.txt", "new.txt")
```
Renomme un fichier

---

### - bif()

```
bif(nom == "Jean")
    bprint("Bonjour Jean!")
|
```
Condition bif - exécute le bloc si la condition est vraie

---

### - belse

```
bif(age >= "18")
    bprint("Majeur")
belse
    bprint("Mineur")
|
```
Bloc belse - s'exécute si la condition bif est fausse

---

### - 