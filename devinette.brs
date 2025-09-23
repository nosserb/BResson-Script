player = "Guy"
brand(1 10)
secret = _brand

bprint("Salut " + player + "! Devine le nombre entre 1 et 10.")

bif(player == "Guy")
    guess = "5"
    bif(guess == secret)
        bprint("Bravo, tu as trouvé du premier coup!")
        bstartfile("celebration.brs")
    belse
        bprint("Raté! Essaie encore.")
        guess = "7"
        bif(guess == secret)
            bprint("Bien joué au deuxième essai!")
        belse
            guess = "3"
            bif(guess == secret)
                bprint("Enfin, tu as trouvé!")
            belse
                bprint("Désolé, tu as perdu. Le nombre était " + secret)
            |
        |
    |
|

btime("j m a h min s")
logtime = _btime

brename("log.txt" "log_" + player + ".txt")
bprint("Le jeu a été joué le " + logtime)
