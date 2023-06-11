package main

import (
	_ "embed"
	"flag"
	"log"
	"time"

	"github.com/signintech/gopdf"
	"github.com/spf13/cobra"
)

//go:embed Inter.ttf
var interFont []byte

var (
	id    string
	title string

	logo string
	from string
	to   string
	date string
	due  string

	items      []string
	quantities []int
	rates      []float64

	tax      float64
	discount float64
	currency string

	note   string
	output string
)

func init() {
	generateCmd.Flags().StringVar(&id, "id", time.Now().Format("20060102"), "ID")
	generateCmd.Flags().StringVar(&title, "title", "INVOICE", "Title")

	generateCmd.Flags().Float64SliceVarP(&rates, "rate", "r", []float64{25}, "Rates")
	generateCmd.Flags().IntSliceVarP(&quantities, "quantity", "q", []int{2}, "Quantities")
	generateCmd.Flags().StringSliceVarP(&items, "item", "i", []string{"Paper Cranes"}, "Items")

	generateCmd.Flags().StringVarP(&logo, "logo", "l", "", "Company logo")
	generateCmd.Flags().StringVarP(&from, "from", "f", "Project Folded, Inc.", "Issuing company")
	generateCmd.Flags().StringVarP(&to, "to", "t", "Untitled Corporation, Inc.", "Recipient company")
	generateCmd.Flags().StringVar(&date, "date", time.Now().Format("Jan 02, 2006"), "Date")
	generateCmd.Flags().StringVar(&due, "due", time.Now().AddDate(0, 0, 14).Format("Jan 02, 2006"), "Payment due date")

	generateCmd.Flags().Float64Var(&tax, "tax", 0.13, "Tax")
	generateCmd.Flags().Float64VarP(&discount, "discount", "d", 0.0, "Discount")
	generateCmd.Flags().StringVarP(&currency, "currency", "c", "USD", "Currency")

	generateCmd.Flags().StringVarP(&note, "note", "n", "For debugging purposes.", "Note")
	generateCmd.Flags().StringVarP(&output, "output", "o", "invoice.pdf", "Output file (.pdf)")

	flag.Parse()
}

var rootCmd = &cobra.Command{
	Use:   "invoice",
	Short: "Invoice generates invoices from the command line.",
	Long:  `Invoice generates invoices from the command line.`,
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate an invoice",
	Long:  `Generate an invoice`,
	RunE: func(cmd *cobra.Command, args []string) error {
		pdf := gopdf.GoPdf{}
		pdf.Start(gopdf.Config{
			PageSize: *gopdf.PageSizeA4,
		})
		pdf.SetMargins(40, 40, 40, 40)
		pdf.AddPage()
		err := pdf.AddTTFFontData("Inter", interFont)
		if err != nil {
			return err
		}

		writeLogo(&pdf, logo, from)
		writeTitle(&pdf, title, id)
		writeBillTo(&pdf, to)
		writeHeaderRow(&pdf)
		subtotal := 0.0
		for i := range items {
			q := 1
			if len(quantities) > i {
				q = quantities[i]
			}

			r := 0.0
			if len(rates) > i {
				r = rates[i]
			}

			writeRow(&pdf, items[i], q, r)
			subtotal += float64(q) * r
		}
		writeNotes(&pdf, note)
		writeTotals(&pdf, subtotal, subtotal*tax, subtotal*discount)
		writeFooter(&pdf, id)
		return pdf.WritePdf(output)
	},
}

func main() {
	rootCmd.AddCommand(generateCmd)
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}