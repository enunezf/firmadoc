package main

import (
    //"crypto/hmac"
    "crypto/rand"
    //"crypto/sha256"
    "encoding/base64"
    //"encoding/hex"
    "fmt"
    //"io"
    "log"
    "os"
    "time"
    "path/filepath"


    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"
    "fyne.io/fyne/v2/dialog"

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


func main() {

    myApp := app.New()
    window := myApp.NewWindow("Sodexo - Firma de documentos")

	var outputPath = "signing_key.txt"

	key := generateRandomKey(1024)
	
	if err := saveKeyToFile(key, outputPath); err != nil {
		log.Fatal("Error saving key:", err)
	}
	
	fmt.Printf("Key generated and saved to %s\n", outputPath)
	fmt.Printf("Key: %s\n", key)

    // Create tabs for different screens
    tabs := container.NewAppTabs(
        container.NewTabItem("Generate Key", createGenerateKeyTab(window)),
    )

    window.SetContent(tabs)
    window.Resize(fyne.NewSize(800, 600))
    window.ShowAndRun()

}

func createGenerateKeyTab(window fyne.Window) fyne.CanvasObject {
    generateBtn := widget.NewButton("Genera clave", func() {
        // Generate random key
        if key, err := generateRandomKey(1024) {
            dialog.ShowError(err, window)
            return
        }

        // Save key to file
        execPath, _ := os.Executable()
        keyPath := filepath.Join(filepath.Dir(execPath), "signing_key.txt")
        if err := os.WriteFile(keyPath, key, 0644); err != nil {
            dialog.ShowError( err, window)
            return
        }

        dialog.ShowInformation("Success", "Key generated and saved to "+keyPath, window)
    })

    return container.NewVBox(
        widget.NewLabel("Click the button below to generate a new signing key:"),
        generateBtn,
    )
}