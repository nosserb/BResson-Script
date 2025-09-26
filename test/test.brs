# Démonstration du système de logging
blog

bprint("=== Test du système de logging ===")

# Variables et calculs
nom = "Jean"
age = "25"
bcalc(age + 5) -> nouvel_age

bprint("Nom: " + nom)
bprint("Nouvel âge: " + nouvel_age)

# Test avec fichier
bwrite("test_log.txt", "Contenu de test")
bread("test_log.txt")

# Test avec timer
btimer("demo")
bsleep(1)
bendtimer("demo")

# Récupérer tous les logs
blogget
logs = _blog

bprint("=== LOGS COMPLETS ===")
bprint(logs)

# Sauvegarder les logs dans un fichier
bwrite("logs_complets.txt", logs)
bprint("Logs sauvegardés dans logs_complets.txt")
