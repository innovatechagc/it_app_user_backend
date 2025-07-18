#!/bin/bash

echo "🚀 Iniciando desarrollo local..."

# Verificar si Firebase CLI está instalado
if ! command -v firebase &> /dev/null; then
    echo "❌ Firebase CLI no está instalado"
    echo "Instala con: npm install -g firebase-tools"
    exit 1
fi

# Verificar si Go está instalado
if ! command -v go &> /dev/null; then
    echo "❌ Go no está instalado"
    exit 1
fi

echo "📦 Instalando dependencias..."
cd functions
go mod tidy
cd ..

echo "🔥 Iniciando Firebase Emulators..."
firebase emulators:start --only functions

echo "✅ Servidor corriendo en:"
echo "   Functions: http://localhost:5001"
echo "   UI: http://localhost:4000"