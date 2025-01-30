# firmadoc

Esta es una aplicación que permite la generación aleatorea de una clave de 1024 caracteres y firmar documentos con esta clave.

## Complilación

** Para generar un binario con el mismo sistema operativo

```bash
fyne package -executable firmadoc -icon icon.png
```

** Para generar un binario hacia otro sistema operativo

```bash
CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ fyne package -os windows -icon icon.png
```

# Guía de Compilación Cruzada con Go y Fyne

## Tabla de Contenidos

- [Compilación Cruzada de Windows desde macOS](#compilación-cruzada-de-windows-desde-macos)
- [Creación de Iconos](#creación-de-iconos)
- [Compilación para macOS](#compilación-para-macos)


## Compilación Cruzada de Windows desde macOS

### Requisitos Previos
```bash
brew install mingw-w64
```

### Comando de Compilación Básico
```bash
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go build -o app.exe
```

### Compilación con Fyne y Recursos
1. Crear archivo de recursos:
```bash
echo "IDI_ICON1 ICON \"icon.ico\"" > appicon.rc
```

2. Compilar recursos:
```bash
x86_64-w64-mingw32-windres appicon.rc -O coff -o appicon.syso
```

3. Compilación final:
```bash
CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 go build -ldflags -H=windowsgui -o SigningApp.exe
```

## Creación de Iconos

### SVG para Icono de Firma de Documentos
```svg
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512">
  <!-- Background -->
  <rect x="96" y="32" width="320" height="448" fill="#ffffff" stroke="#000000" stroke-width="16"/>
  
  <!-- Document Lines -->
  <line x1="144" y1="128" x2="368" y2="128" stroke="#666666" stroke-width="8"/>
  <line x1="144" y1="176" x2="368" y2="176" stroke="#666666" stroke-width="8"/>
  <line x1="144" y1="224" x2="368" y2="224" stroke="#666666" stroke-width="8"/>
  
  <!-- Signature Line -->
  <line x1="144" y1="320" x2="368" y2="320" stroke="#000000" stroke-width="4"/>
  
  <!-- Signature -->
  <path d="M144 320 C200 280, 240 360, 280 320 S320 280, 368 320" 
        fill="none" 
        stroke="#0066cc" 
        stroke-width="8"
        stroke-linecap="round"/>
        
  <!-- Decorative Corner -->
  <path d="M96 32 L136 32 L136 72" 
        fill="none" 
        stroke="#000000" 
        stroke-width="16"/>
</svg>
```

### Conversión de SVG a PNG
```bash
brew install librsvg
rsvg-convert -h 256 icon.svg > icon.png
```

## Compilación para macOS

### Usando Fyne Package
```bash
fyne package -os darwin -icon icon.png
```

### Usando Go Build
```bash
# Compilar el binario
go build -o MyApp

# Crear el .app con icono
fyne package -executable MyApp -icon icon.png
```

### Requisitos Previos para macOS
```bash
# Instalar herramientas de Xcode
xcode-select --install

# Instalar herramientas de Fyne
go install fyne.io/fyne/v2/cmd/fyne@latest
```

## Notas Adicionales
- Asegúrate de tener todas las dependencias instaladas antes de comenzar
- Los iconos deben estar en el formato correcto para cada sistema operativo
- Para macOS, el resultado será un archivo .app
- Para Windows, el resultado será un archivo .exe


