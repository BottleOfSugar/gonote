package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func saveToFile(filename string, content string) error {
	return ioutil.WriteFile(filename, []byte(content), 0644)
}

func openFile(filename string) (string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func main() {
	// Inicjalizacja aplikacji Fyne
	myApp := app.New()
	myWindow := myApp.NewWindow("Notatnik")

	// Tworzenie edytora
	textArea := widget.NewMultiLineEntry()
	textArea.SetPlaceHolder("Wpisz tutaj swoje notatki...")

	// Etykieta do wyświetlania otwartego folderu
	folderLabel := widget.NewLabel("Otwarte folder: Brak")

	// Przycisk "Zapisz"
	saveButton := widget.NewButton("Zapisz", func() {
		dialog.ShowFileSave(func(uc fyne.URIWriteCloser, err error) {
			if uc == nil || err != nil {
				return
			}
			err = saveToFile(uc.URI().Path(), textArea.Text)
			if err != nil {
				dialog.ShowError(fmt.Errorf("Błąd zapisu pliku: %v", err), myWindow)
				return
			}
			folderLabel.SetText("Otwarte folder: " + filepath.Dir(uc.URI().Path()))
			uc.Close()
		}, myWindow)
	})

	// Przycisk "Otwórz"
	openButton := widget.NewButton("Otwórz", func() {
		dialog.ShowFileOpen(func(uc fyne.URIReadCloser, err error) {
			if uc == nil || err != nil {
				return
			}
			content, err := openFile(uc.URI().Path())
			if err != nil {
				dialog.ShowError(fmt.Errorf("Błąd otwarcia pliku: %v", err), myWindow)
				return
			}
			textArea.SetText(content)
			folderLabel.SetText("Otwarte folder: " + filepath.Dir(uc.URI().Path()))
			uc.Close()
		}, myWindow)
	})

	// Przycisk "Wyczyść"
	clearButton := widget.NewButton("Wyczyść", func() {
		textArea.SetText("")
		folderLabel.SetText("Otwarte folder: Brak")
	})

	// Ułożenie elementów GUI w pionie
	content := container.NewVBox(
		folderLabel,
		textArea,
		container.NewHBox(saveButton, openButton, clearButton),
	)

	// Ustawienie zawartości okna
	myWindow.SetContent(content)

	// Ustawienia okna
	myWindow.Resize(fyne.NewSize(600, 400))

	// Wyświetlenie okna
	myWindow.ShowAndRun()
}
