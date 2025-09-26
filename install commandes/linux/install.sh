#!/bin/bash

# Bresson Script Language Installer for Linux/macOS
# This script installs the Bresson language interpreter system-wide

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
INSTALL_DIR="/usr/local/bin"
LIB_DIR="/usr/local/lib/bresson"
VERSION="1.0.0"

print_header() {
    echo -e "${BLUE}"
    echo "╔══════════════════════════════════════════════════════════════╗"
    echo "║                    Bresson Script Installer                  ║"
    echo "║                        Version $VERSION                         ║"
    echo "╚══════════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
}

print_step() {
    echo -e "${YELLOW}[ÉTAPE]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCÈS]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERREUR]${NC} $1"
}

check_requirements() {
    print_step "Vérification des prérequis..."
    
    # Check if Go is installed
    if ! command -v go &> /dev/null; then
        print_error "Go n'est pas installé. Veuillez installer Go depuis https://golang.org/dl/"
        exit 1
    fi
    
    # Check Go version
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    print_success "Go version $GO_VERSION détectée"
    
    # Check if we have the required files
    if [[ ! -f "bras.go" ]] || [[ ! -f "main.go" ]]; then
        print_error "Fichiers source manquants (bras.go ou main.go)"
        exit 1
    fi
    
    print_success "Tous les prérequis sont satisfaits"
}

build_binaries() {
    print_step "Compilation des binaires..."
    
    # Build the launcher
    echo "  - Compilation de bras (launcher)..."
    go build -o bras bras.go
    if [[ $? -ne 0 ]]; then
        print_error "Échec de la compilation de bras.go"
        exit 1
    fi
    
    # Build the interpreter
    echo "  - Compilation de l'interpréteur..."
    go build -o bresson main.go
    if [[ $? -ne 0 ]]; then
        print_error "Échec de la compilation de main.go"
        exit 1
    fi
    
    print_success "Compilation terminée"
}

install_system_wide() {
    print_step "Installation système..."
    
    # Check if we need sudo
    if [[ ! -w "$INSTALL_DIR" ]]; then
        echo "  - Privilèges administrateur requis pour l'installation dans $INSTALL_DIR"
        SUDO="sudo"
    else
        SUDO=""
    fi
    
    # Create lib directory if it doesn't exist
    $SUDO mkdir -p "$LIB_DIR"
    
    # Install the interpreter in lib directory
    echo "  - Installation de l'interpréteur dans $LIB_DIR..."
    $SUDO cp bresson "$LIB_DIR/"
    $SUDO chmod +x "$LIB_DIR/bresson"
    
    # Install the launcher in bin directory
    echo "  - Installation du launcher dans $INSTALL_DIR..."
    $SUDO cp bras "$INSTALL_DIR/"
    $SUDO chmod +x "$INSTALL_DIR/bras"
    
    print_success "Installation système terminée"
}

update_path() {
    print_step "Vérification du PATH..."
    
    # Check if /usr/local/bin is in PATH
    if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
        echo "  - $INSTALL_DIR n'est pas dans votre PATH"
        echo "  - Ajoutez cette ligne à votre ~/.bashrc ou ~/.zshrc :"
        echo -e "    ${YELLOW}export PATH=\"$INSTALL_DIR:\$PATH\"${NC}"
        echo "  - Puis rechargez votre shell avec: source ~/.bashrc"
    else
        print_success "$INSTALL_DIR est déjà dans votre PATH"
    fi
}

create_test_script() {
    print_step "Création d'un script de test..."
    
    cat > test.brs << 'EOF'
# Script de test pour Bresson
bprint("=== Test du langage Bresson ===")
bprint("Bonjour depuis Bresson Script!")

binput("Entrez votre nom: ") -> nom
bprint("Salut " + nom + "!")

bcalc(2 + 3) -> resultat
bprint("2 + 3 = " + resultat)

bprint("Test terminé avec succès!")
EOF
    
    print_success "Script de test créé: test.brs"
}

test_installation() {
    print_step "Test de l'installation..."
    
    # Test if bras command is available
    if command -v bras &> /dev/null; then
        print_success "Commande 'bras' disponible"
        
        # Test with the test script
        echo "  - Test avec le script de démonstration..."
        echo "Jean" | bras test.brs > /dev/null 2>&1
        if [[ $? -eq 0 ]]; then
            print_success "Test d'exécution réussi"
        else
            print_error "Test d'exécution échoué"
        fi
    else
        print_error "Commande 'bras' non trouvée. Vérifiez votre PATH."
    fi
}

cleanup() {
    print_step "Nettoyage..."
    rm -f bras bresson test.brs
    print_success "Fichiers temporaires supprimés"
}

show_usage() {
    echo -e "${GREEN}"
    echo "╔══════════════════════════════════════════════════════════════╗"
    echo "║                     Installation terminée!                  ║"
    echo "╚══════════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
    echo
    echo "Utilisation:"
    echo "  bras script.brs              # Exécuter un script"
    echo "  bras script.brs arg1 arg2    # Avec des arguments"
    echo
    echo "Exemple:"
    echo "  echo 'bprint(\"Hello World!\")' > hello.brs"
    echo "  bras hello.brs"
    echo
    echo -e "Documentation complète: ${BLUE}README.md${NC}"
}

# Main installation process
main() {
    print_header
    
    check_requirements
    build_binaries
    install_system_wide
    update_path
    create_test_script
    test_installation
    cleanup
    show_usage
    
    echo
    print_success "Installation de Bresson Script terminée!"
}

# Handle script arguments
case "${1:-}" in
    --help|-h)
        echo "Usage: $0 [--help|--uninstall]"
        echo "  --help      Affiche cette aide"
        echo "  --uninstall Désinstalle Bresson Script"
        exit 0
        ;;
    --uninstall)
        print_step "Désinstallation de Bresson Script..."
        sudo rm -f "$INSTALL_DIR/bras"
        sudo rm -rf "$LIB_DIR"
        print_success "Bresson Script désinstallé"
        exit 0
        ;;
    "")
        main
        ;;
    *)
        echo "Option inconnue: $1"
        echo "Utilisez --help pour l'aide"
        exit 1
        ;;
esac
