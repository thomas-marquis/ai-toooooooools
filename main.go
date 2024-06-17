package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/spf13/cobra"
)

func fetchHTMLContent(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func generatePDF(url, outputFileName string) error {
	// Créer un contexte
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Créer un autre contexte avec un timeout pour l'opération
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Buffer pour le PDF
	var buf []byte

	// Exécuter les tâches
	if err := chromedp.Run(ctx, chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			buf, _, err = page.PrintToPDF().WithPrintBackground(true).Do(ctx)
			return err
		}),
	}); err != nil {
		return err
	}

	// Écrire le PDF dans un fichier
	if err := os.WriteFile(outputFileName, buf, 0644); err != nil {
		return err
	}

	return nil
}

func main() {
	var outputFileName string
	var rootCmd = &cobra.Command{
		Use:   "html2pdf",
		Short: "Convert HTML to PDF",
		Long:  `A tool to convert HTML pages to PDF files using chromedp.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				log.Fatalf("Usage: %s <url> <output_file>", cmd.Use)
			}

			url := args[0]
			if len(args) > 1 {
				outputFileName = args[1]
			} else {
				outputFileName = "output.pdf"
			}

			// Générer le PDF
			err := generatePDF(url, outputFileName)
			if err != nil {
				log.Fatalf("Failed to create PDF: %v", err)
			}

			fmt.Printf("PDF saved as %s\n", outputFileName)
		},
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
