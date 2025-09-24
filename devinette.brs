bprint("Bienvenue dans le jeu de devinette !")
binput("Comment tu t'appelles ?") -> 'nom'

bprint("Salut " + 'nom' + "! Je pense à un nombre entre 1 et 10...")

bcalc(rand(1,10)) -> 'nombreMystere'

binput("Devine le nombre :") -> 'guess'
bcolor("rouge")
bif('guess' == 'nombreMystere')
    bcolor("vert")
    bprint("Bravo ! Tu as trouvé ! 🎉")
belse
    bprint("Dommage ! Le nombre était "  'nombreMystere')
|
