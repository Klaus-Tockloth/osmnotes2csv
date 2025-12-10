#!/bin/bash

# Script bricht ab, wenn ein Befehl fehlschlägt
set -e

BUILD_DIR="build"

# Prüfen, ob das Build-Verzeichnis existiert
if [ ! -d "$BUILD_DIR" ]; then
    echo "Fehler: Verzeichnis '$BUILD_DIR' nicht gefunden."
    exit 1
fi

echo "Starte Komprimierung der Artefakte in '$BUILD_DIR'..."

# Durchlaufe alle Unterverzeichnisse in build/
for dir_path in "$BUILD_DIR"/*; do
    if [ -d "$dir_path" ]; then
        # Hole den reinen Ordnernamen (z.B. darwin-amd64)
        dir_name=$(basename "$dir_path")

        # In das Verzeichnis wechseln, damit das Zip keine Pfade enthält
        pushd "$dir_path" > /dev/null

        # Suche nach der Binary-Datei (ignoriert bereits existierende .zip Dateien)
        # Wir nehmen an, die Datei beginnt mit "osmnotes2csv"
        bin_file=$(find . -maxdepth 1 -type f -name "osmnotes2csv*" ! -name "*.zip" | head -n 1 | sed 's|^\./||')

        if [ -n "$bin_file" ]; then
            # Zielname: Ordnername-Dateiname.zip
            # Beispiel: darwin-amd64-osmnotes2csv.zip oder windows-amd64-osmnotes2csv.exe.zip
            zip_name="${dir_name}-${bin_file}.zip"

            # Zip erstellen (-q für quiet output)
            zip -q "$zip_name" "$bin_file"
            
            echo "Erstellt: $dir_path/$zip_name"
        else
            echo "Warnung: Keine ausführbare Datei in '$dir_name' gefunden."
        fi

        # Zurück ins Ursprungsverzeichnis
        popd > /dev/null
    fi
done

echo "Fertig."
