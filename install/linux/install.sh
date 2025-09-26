#!/bin/bash
set -e

# ==========================
# Bresson Script Installer
# ==========================

GO_VERSION="1.23.2"
GO_URL="https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz"
BIN_DIR="$HOME/bin"
VERSION="1.0.0"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}Bresson Script Installer v$VERSION${NC}"

# --------------------------
# Installer Go si absent
# --------------------------
if ! command -v go &> /dev/null; then
    echo -e "${YELLOW}[INFO] Go non trouvé, installation dans $BIN_DIR/go${NC}"
    mkdir -p "$BIN_DIR"
    wget -q "$GO_URL" -O /tmp/go.tar.gz
    tar -C "$BIN_DIR" -xzf /tmp/go.tar.gz
    rm /tmp/go.tar.gz
    export PATH="$BIN_DIR/go/bin:$PATH"
    if ! grep -q "bin/go/bin" ~/.bashrc; then
        echo 'export PATH=$HOME/bin/go/bin:$PATH' >> ~/.bashrc
    fi
    echo -e "${GREEN}[OK] Go installé${NC}"
else
    echo -e "${GREEN}[OK] Go détecté: $(go version)${NC}"
fi

# --------------------------
# Détecter dossier source
# --------------------------
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"   # remonte 2 niveaux pour atteindre la racine

# Vérifier que les fichiers existent
if [[ ! -f "$SCRIPT_DIR/bras.go" ]]; then
    echo -e "${RED}[ERREUR] bras.go introuvable dans $SCRIPT_DIR${NC}"
    exit 1
fi
if [[ ! -f "$REPO_ROOT/main.go" ]]; then
    echo -e "${RED}[ERREUR] main.go introuvable dans $REPO_ROOT${NC}"
    exit 1
fi

# --------------------------
# Compiler Bresson
# --------------------------
echo -e "${YELLOW}[INFO] Compilation des binaires...${NC}"

go build -o bras "$SCRIPT_DIR/bras.go"
go build -o bresson "$REPO_ROOT/main.go"

# Créer bin si pas existant
mkdir -p "$BIN_DIR"

# Copier les binaires
cp bras bresson "$BIN_DIR/"
chmod +x "$BIN_DIR/bras" "$BIN_DIR/bresson"

echo -e "${GREEN}[OK] Binaries installés dans $BIN_DIR${NC}"

# --------------------------
# Vérifier PATH
# --------------------------
if [[ ":$PATH:" != *":$BIN_DIR:"* ]]; then
    echo 'export PATH=$HOME/bin:$PATH' >> ~/.bashrc
    export PATH="$BIN_DIR:$PATH"
    echo -e "${YELLOW}[INFO] $BIN_DIR ajouté au PATH. Rechargez le shell ou faites source ~/.bashrc${NC}"
fi

# --------------------------
# Créer script de test
# --------------------------
cat > test.brs << 'EOF'
bprint("=== Test du langage Bresson ===")
bprint("Bonjour depuis Bresson Script!")

binput("Entrez votre nom: ") -> nom
bprint("Salut " + nom + "!")

bcalc(2 + 3) -> resultat
bprint("2 + 3 = " + resultat)

bprint("Test terminé avec succès!")
EOF

echo -e "${GREEN}[OK] Script test.brs créé${NC}"

# --------------------------
# Test rapide
# --------------------------
echo -e "${YELLOW}[INFO] Test rapide d'exécution...${NC}"
echo "Jean" | "$BIN_DIR/bras" test.brs

echo -e "${GREEN}[INSTALLATION TERMINÉE] Vous pouvez maintenant utiliser 'bras script.brs' depuis n'importe quel dossier${NC}"
