bprint("=== Test du système de logging ===")

# Variables et calculs
nom = "Jean"
age = "25"
bcalc(age + 5) -> nouvel_age

bprint("Nom: " + nom)
bprint("Nouvel âge: " + nouvel_age)

# Test avec fichier temporaire
bwrite("test_log.txt", "Contenu de test")
bread("test_log.txt")

# Test avec timer
btimer("demo")
bsleep(1)
bendtimer("demo")

# Récupérer tous les logs
blogget
logs = _blog

# Convertir la liste de logs en chaîne de caractères pour écrire dans le fichier
logs_str = join(logs, "\n")  # Chaque log sur une ligne

bprint("=== LOGS COMPLETS ===")
bprint(logs_str)

# Sauvegarder tous les logs dans un fichier
bwrite("logs_complets.txt", logs_str)
bprint("Logs sauvegardés dans logs_complets.txt")
