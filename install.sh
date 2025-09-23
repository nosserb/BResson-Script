#!/bin/bash

echo "Compilation de bras..."
go build -o bras bras.go

echo "Déplacement dans /usr/local/bin (sudo requis)..."
sudo mv bras /usr/local/bin/

echo "Ajout des permissions..."
sudo chmod +x /usr/local/bin/bras

echo "Installation terminée !"
echo "Tu peux maintenant lancer : bras tonfichier.brs"

echo "Pour l’utiliser :"
echo "chmod +x install.sh"
echo "./install.sh"