package main

import (
    "crypto/hmac"
    "crypto/rand"
    "crypto/sha256"
    "encoding/base64"
    "encoding/hex"
    "io"
    "os"
    "time"
    "path/filepath"
    "errors"


    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"
    "fyne.io/fyne/v2/dialog"
    "fyne.io/fyne/v2/layout"

)


type SigningKey struct {
    Key       string    `json:"key"`
    CreatedAt time.Time `json:"created_at"`
}

func generateRandomKey(length int) (string, error) {
    bytes := make([]byte, length)
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}

func saveKeyToFile(key string, outputPath string) error {
    return os.WriteFile(outputPath, []byte(key), 0600)
}

func createSignature(filePath string, key string) (string, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return "", err
    }
    defer file.Close()

    h := hmac.New(sha256.New, []byte(key))
    if _, err := io.Copy(h, file); err != nil {
        return "", err
    }

    return hex.EncodeToString(h.Sum(nil)), nil
}

func validateSignature(keyPath, filePath, sigPath string) (bool, error) {
    // Calculate current signature
    currentSig, err := createSignature(keyPath, filePath)
    if err != nil {
        return false, err
    }

    // Read stored signature
    storedSig, err := os.ReadFile(sigPath)
    if err != nil {
        return false, err
    }

    return currentSig == string(storedSig), nil
}

func main() {

    myApp := app.New()
    window := myApp.NewWindow("Sodexo - Firma de documentos")

    // Create tabs for different screens
    tabs := container.NewAppTabs(
        container.NewTabItem("Firma documento", createSigningTab(window)),
        container.NewTabItem("Valida firma", createValidationTab(window)),
        container.NewTabItem("Genera clave", createGenerateKeyTab(window)),
    )

    window.SetContent(tabs)
    window.Resize(fyne.NewSize(800, 600))
    window.ShowAndRun()

}

func createSigningTab(window fyne.Window) fyne.CanvasObject {
    // Get default application directory
    execPath, _ := os.Executable()
    defaultKeyPath := filepath.Join(filepath.Dir(execPath), "signing_key.txt")

    // Create entry widgets
    keyEntry := widget.NewEntry()
    keyEntry.SetText(defaultKeyPath)
    fileEntry := widget.NewEntry()

    // Create browse buttons
    browseKeyBtn := widget.NewButton("Buscar", func() {
        fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
            if err == nil && reader != nil {
                keyEntry.SetText(reader.URI().Path())
            }
        }, window)
        fd.Show()
    })

    browseFileBtn := widget.NewButton("Buscar", func() {
        fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
            if err == nil && reader != nil {
                fileEntry.SetText(reader.URI().Path())
            }
        }, window)
        fd.Show()
    })

    // Create sign button with custom styling
    signBtn := widget.NewButton("Firmar", func() {
        if keyEntry.Text == "" || fileEntry.Text == "" {
            dialog.ShowError(errors.New("debe completar las rutas de ambos archivos"), window)
            return
        }

        signature, err := createSignature(keyEntry.Text, fileEntry.Text)
        if err != nil {
            dialog.ShowError(err, window)
            return
        }

        // Save signature to file
        sigPath := fileEntry.Text + ".sig"
        err = os.WriteFile(sigPath, []byte(signature), 0644)
        if err != nil {
            dialog.ShowError(err, window)
            return
        }

        dialog.ShowInformation("Success", "Firma fue grabada en: "+sigPath, window)
    })
    
    // Style the sign button
    signBtn.Importance = widget.HighImportance // Makes the button green
    
    // Create a container for the sign button that centers it and sets its width to 50%
    signBtnContainer := container.NewHBox(
        layout.NewSpacer(),
        container.NewGridWrap(
            fyne.NewSize(300, 40), // 300 pixels is approximately 50% of the default window width
            signBtn,
        ),
        layout.NewSpacer(),
    )

    // Layout
    keyBox := container.NewBorder(nil, nil, nil, browseKeyBtn, keyEntry)
    fileBox := container.NewBorder(nil, nil, nil, browseFileBtn, fileEntry)

    // Add spacing between fileBox and signBtn
    spacer := widget.NewSeparator()

    return container.NewVBox(
        widget.NewLabel("Clave:"),
        keyBox,
        widget.NewLabel("Archivo:"),
        fileBox,
        layout.NewSpacer(), // Add vertical space
        spacer,             // Add visual separator
        layout.NewSpacer(), // Add more vertical space
        signBtnContainer,   // Centered sign button
        layout.NewSpacer(), // Add bottom padding
    )
}

 
func createValidationTab(window fyne.Window) fyne.CanvasObject {
    // Get default application directory
    execPath, _ := os.Executable()
    defaultKeyPath := filepath.Join(filepath.Dir(execPath), "signing_key.txt")

    // Create entry widgets
    keyEntry := widget.NewEntry()
    keyEntry.SetText(defaultKeyPath)
    fileEntry := widget.NewEntry()
    sigEntry := widget.NewEntry()

    // Create browse buttons
    browseKeyBtn := widget.NewButton("Buscar", func() {
        fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
            if err == nil && reader != nil {
                keyEntry.SetText(reader.URI().Path())
            }
        }, window)
        fd.Show()
    })

    browseFileBtn := widget.NewButton("Buscar", func() {
        fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
            if err == nil && reader != nil {
                filePath := reader.URI().Path()
                fileEntry.SetText(filePath)
                sigEntry.SetText(filePath + ".sig")
            }
        }, window)
        fd.Show()
    })

    browseSigBtn := widget.NewButton("Buscar", func() {
        fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
            if err == nil && reader != nil {
                sigEntry.SetText(reader.URI().Path())
            }
        }, window)
        fd.Show()
    })

    // Create validate button
    validateBtn := widget.NewButton("Validar", func() {
        if keyEntry.Text == "" || fileEntry.Text == "" || sigEntry.Text == "" {
            dialog.ShowError(errors.New("debe seleccionar los archivos de: clave, archivo y firma"), window)
            return
        }

        isValid, err := validateSignature(keyEntry.Text, fileEntry.Text, sigEntry.Text)
        if err != nil {
            dialog.ShowError(err, window)
            return
        }

        if isValid {
            dialog.ShowInformation("Success", "La firma es valida", window)
        } else {
            dialog.ShowInformation("Invalid", "La firma no es valida", window)
        }
    })

    // Layout
    keyBox := container.NewBorder(nil, nil, nil, browseKeyBtn, keyEntry)
    fileBox := container.NewBorder(nil, nil, nil, browseFileBtn, fileEntry)
    sigBox := container.NewBorder(nil, nil, nil, browseSigBtn, sigEntry)

    return container.NewVBox(
        widget.NewLabel("Clave:"),
        keyBox,
        widget.NewLabel("Archivo:"),
        fileBox,
        widget.NewLabel("Firma:"),
        sigBox,
        validateBtn,
    )
}


func createGenerateKeyTab(window fyne.Window) fyne.CanvasObject {
    generateBtn := widget.NewButton("Genera clave", func() {
        // Generate random key
        key, err := generateRandomKey(1024)
        if err != nil {
            dialog.ShowError(err, window)
            return
        }

        // Save key to file
        execPath, _ := os.Executable()
        keyPath := filepath.Join(filepath.Dir(execPath), "signing_key.txt")

        err = saveKeyToFile(key, keyPath)
        if err != nil {
            dialog.ShowError(err, window)
            return
        }

        dialog.ShowInformation("Success", "El archivo de clave se ha grabado exitosamente: "+ keyPath, window)
    })

    return container.NewVBox(
        widget.NewLabel("Click en el botón para generar una nueva clave:"),
        generateBtn,
    )
}