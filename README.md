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
Lit le contenu d'un fichier et le stocke dans `_bread`

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

### - bwhile()

```
compteur = "1"
bwhile(compteur <= 5)
    bprint("Compteur: " + compteur)
    bcalc(compteur + 1) -> compteur
|
```
Boucle while - répète tant que la condition est vraie

---

### - " | "

Ferme un bloc (if, else, while)

---

### - bsleep()

```
bsleep(3)
bprint("3 secondes écoulées")
```
Met en pause l'exécution pendant X secondes

---

### - brand()

```
brand(1, 10)
nombre = _brand
```
Génère un nombre aléatoire entre min et max, stocké dans `_brand`

--- 

### - btime()

```
btime(j m a h min)
bprint("Date: " + _btime)
```
Récupère la date/heure actuelle, stockée dans `_btime`
**Formats :** `j` (jour), `m` (mois), `a` (année), `h` (heure), `min` (minute), `s` (seconde)

---

### - bcolor()

```
bcolor("rouge")
bprint("Texte en rouge")
bcolor("vert")
bprint("Texte en vert")
```
Change la couleur d'affichage pour les prochains `bprint`
**Couleurs :** `rouge`, `vert`, `bleu`, `jaune`

---

### - btimer()

```
btimer("monTimer")
```
Démarre un timer avec un nom

---

### - bendtimer()

```
btimer("test")
bsleep(2)
bendtimer("test")
```
Arrête un timer et affiche la durée écoulée

---

### - bard()

```
bard(1)
premier_arg = _bard
```
Récupère un argument de ligne de commande, stocké dans `_bard`

---

### - bstartfile()

```
bstartfile("script.exe")
```
Indique l'exécution d'un fichier (affichage uniquement)

---

### - ` # `

```
# Ceci est un commentaire
bprint("Hello World")
```
Ligne de commentaire (ignorée à l'exécution)

---

## VARIABLE SPECIAL

- `_bread` : Contient le contenu du dernier fichier lu
- `_brand` : Contient le dernier nombre aléatoire généré
- `_btime` : Contient la dernière date/heure formatée
- `_bard` : Contient le dernier argument de ligne de commande récupéré

---


