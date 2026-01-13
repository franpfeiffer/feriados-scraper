package scraper

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/franpfeiffer/feriados-scraper/internal/models"
)

const feriadosURL = "https://www.bcra.gob.ar/consulta-feriados-bancarios/"

type FeriadoScraper struct {
	client *http.Client
}

func NewFeriadoScraper() *FeriadoScraper {
	return &FeriadoScraper{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *FeriadoScraper) GetFeriados() ([]models.Feriado, error) {
	resp, err := s.client.Get(feriadosURL)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("html parse error: %w", err)
	}

	var feriados []models.Feriado
	currentYear := time.Now().Year()

	doc.Find("table tr").Each(func(i int, row *goquery.Selection) {
		if i == 0 {
			return
		}

		cols := row.Find("td")
		if cols.Length() < 2 {
			return
		}

		dateText := strings.TrimSpace(cols.Eq(0).Text())
		description := strings.TrimSpace(cols.Eq(1).Text())

		date, err := parseDate(dateText, currentYear)
		if err != nil {
			return
		}

		feriados = append(feriados, models.Feriado{
			Date:        date,
			Description: description,
		})
	})

	return feriados, nil
}

func parseDate(text string, year int) (time.Time, error) {
	months := map[string]time.Month{
		"enero":      time.January,
		"febrero":    time.February,
		"marzo":      time.March,
		"abril":      time.April,
		"mayo":       time.May,
		"junio":      time.June,
		"julio":      time.July,
		"agosto":     time.August,
		"septiembre": time.September,
		"octubre":    time.October,
		"noviembre":  time.November,
		"diciembre":  time.December,
	}

	text = strings.ToLower(strings.TrimSpace(text))

	var day int
	var monthName string

	_, err := fmt.Sscanf(text, "%d de %s", &day, &monthName)
	if err != nil {
		_, err = fmt.Sscanf(text, "%d y %*d de %s", &day, &monthName)
		if err != nil {
			return time.Time{}, fmt.Errorf("could not parse date: %s", text)
		}
	}

	month, ok := months[monthName]
	if !ok {
		return time.Time{}, fmt.Errorf("unknown month: %s", monthName)
	}

	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC), nil
}

