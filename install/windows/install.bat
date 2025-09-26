@echo off
setlocal enabledelayedexpansion

:: Bresson Script Language Installer for Windows
:: This script installs the Bresson language interpreter system-wide

echo.
echo ╔══════════════════════════════════════════════════════════════╗
echo ║                    Bresson Script Installer                  ║
echo ║                        Version 1.0.0                        ║
echo ╚══════════════════════════════════════════════════════════════╝
echo.

:: Check if Go is installed
echo [ÉTAPE] Vérification des prérequis...
go version >nul 2>&1
if errorlevel 1 (
    echo [ERREUR] Go n'est pas installé. Veuillez installer Go depuis https://golang.org/dl/
    pause
    exit /b 1
)

:: Check if source files exist
if not exist "bras.go" (
    echo [ERREUR] Fichier bras.go manquant
    pause
    exit /b 1
)

if not exist "..\main.go" (
    echo [ERREUR] Fichier main.go manquant
    pause
    exit /b 1
)

echo [SUCCÈS] Tous les prérequis sont satisfaits

:: Build binaries
echo.
echo [ÉTAPE] Compilation des binaires...
echo   - Compilation de bras.exe...
go build -o bras.exe bras.go
if errorlevel 1 (
    echo [ERREUR] Échec de la compilation de bras.go
    pause
    exit /b 1
)

echo   - Compilation de bresson.exe...
go build -o bresson.exe ..\main.go
if errorlevel 1 (
    echo [ERREUR] Échec de la compilation de main.go
    pause
    exit /b 1
)

echo [SUCCÈS] Compilation terminée

:: Install system-wide
echo.
echo [ÉTAPE] Installation système...
set "INSTALL_DIR=%ProgramFiles%\Bresson"

:: Create installation directory
if not exist "%INSTALL_DIR%" (
    mkdir "%INSTALL_DIR%"
)

:: Copy binaries
echo   - Installation dans %INSTALL_DIR%...
copy bras.exe "%INSTALL_DIR%\" >nul
copy bresson.exe "%INSTALL_DIR%\" >nul

:: Add to PATH
echo   - Ajout au PATH système...
for /f "tokens=2*" %%A in ('reg query "HKLM\SYSTEM\CurrentControlSet\Control\Session Manager\Environment" /v PATH 2^>nul') do set "CURRENT_PATH=%%B"

:: Check if already in PATH
echo !CURRENT_PATH! | findstr /i "%INSTALL_DIR%" >nul
if errorlevel 1 (
    reg add "HKLM\SYSTEM\CurrentControlSet\Control\Session Manager\Environment" /v PATH /t REG_EXPAND_SZ /d "!CURRENT_PATH!;%INSTALL_DIR%" /f >nul
    echo [SUCCÈS] Ajouté au PATH système
) else (
    echo [INFO] Déjà dans le PATH système
)

:: Create test script
echo.
echo [ÉTAPE] Test de l'installation...
echo # Test Bresson > test.brs
echo bprint("Hello from Bresson!") >> test.brs

:: Test installation
"%INSTALL_DIR%\bras.exe" test.brs >nul 2>&1
if errorlevel 1 (
    echo [ERREUR] Test d'installation échoué
) else (
    echo [SUCCÈS] Test d'installation réussi
)

:: Cleanup
del bras.exe bresson.exe test.brs >nul 2>&1

echo.
echo ╔══════════════════════════════════════════════════════════════╗
echo ║                     Installation terminée!                  ║
echo ╚══════════════════════════════════════════════════════════════╝
echo.
echo Redémarrez votre invite de commande pour utiliser la commande 'bras'
echo.
echo Utilisation:
echo   bras script.brs
echo   bras script.brs arg1 arg2
echo.
pause
