# Bresson Script (BRS)

**Bresson** est un langage de programmation interprété conçu pour la simplicité et la lisibilité. Inspiré par une approche minimaliste, il privilégie la clarté du code et la facilité d'apprentissage tout en offrant les fonctionnalités essentielles pour le développement d'applications.

---

# Liste de fonctions

---

### - bfile() 

```
bfile("nom_fichier")
```
Recherche le fichier depuis la racine où est exécutée la commande

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

# Conditions & Boucles

---
### - bif()

```
bif(nom == "Jean")
    bprint("Bonjour Jean!")
|
```
Exécute le bloc si la condition est vraie

---

### - belse

```
bif(age >= "18")
    bprint("Majeur")
belse
    bprint("Mineur")
|
```
Exécute le bloc si la condition bif est fausse

---

### - bwhile()

```
compteur = "1"
bwhile(compteur <= 5)
    bprint("Compteur: " + compteur)
    bcalc(compteur + 1) -> compteur
|
```
Répète le bloc tant que la condition est vraie

---

### - " | "

Ferme un bloc (if, else, while)

---

## VARIABLES SPÉCIALES

- `_bread` : Contient le contenu du dernier fichier lu
- `_brand` : Contient le dernier nombre aléatoire généré
- `_btime` : Contient la dernière date/heure formatée
- `_bard` : Contient le dernier argument de ligne de commande récupéré

---


