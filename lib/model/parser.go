package model

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/foolin/pagser"
	"github.com/joaquinrovira/upv-oos-reservations/lib/model/daytime"
	"github.com/joaquinrovira/upv-oos-reservations/lib/model/timerange"
)

type TableData struct {
	ColumnHeaders []string       `pagser:"thead tr th->eachText()"`
	Rows          []TableDataRow `pagser:"tbody tr"`
}

type TableDataRow struct {
	Columns []struct {
		RawText string `pagser:"->html()"`
		URL     string `pagser:"a->attr(href)"`
	} `pagser:"td"`
}

var idxToWeekday = []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday, time.Saturday, time.Sunday}

// Find the target table from an HTTP response
func FindTable(res *http.Response) (selection *goquery.Selection, err error) { // Initialize goquery object with response body
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return
	}
	res.Body.Close()

	// https://www.codetable.net/hex/e9 == é
	// https://www.codetable.net/hex/e1 == á
	tableHeader := []string{"Horario", "Lunes", "Martes", "Mi\xe9rcoles", "Jueves", "Viernes", "S\xe1bado", "Domingo"}

	// Find target table, search for matching header
	selection = doc.Find("table.upv_listacolumnas").FilterFunction(func(i int, s *goquery.Selection) bool {
		matches := true

		// Select the received table headers and validate against tableHeader
		s.Find("thead > tr > th").EachWithBreak(func(i int, s *goquery.Selection) bool {
			if i > len(tableHeader) || s.Text() != tableHeader[i] {
				matches = false // Tag table as not matching
				return false    // Break early
			}

			return true
		})
		return matches
	})

	return
}

func ParseHTMLTable(s *goquery.Selection) (data TableData, err error) {
	p := pagser.New()
	err = p.ParseSelection(&data, s)
	return
}

func MarshalTable(data *TableData) (reservations *ReservationsWeek, err error) {
	reservations = NewReservarionsWeek()

	// Fill reservations.SlotTimes[] with their corresponding values
	reservations.SlotTimes, err = parseTimeSlotsFromTable(data)
	if err != nil {
		return
	}

	// Initialize slots for each weekday (monday, tuesday, wednesday, etc.)
	// Mon-Sun is represented by a different column starting at index 1 (monday) and ending at column 7 (sunday)
	for idx, weekDay := range idxToWeekday {
		reservations.Slots[weekDay], _ = parseWeekdaySlotsFromTable(data, idx)
	}

	return
}

func parseTimeSlotsFromTable(data *TableData) (slotTimes []timerange.TimeRange, err error) {
	for row := 0; row < len(data.Rows); row++ {
		// For each row of the table, the first column contains the time slot description
		// Time slot text follows the format 'XX:YY-AA:BB <some text>'
		// XX:YY being the start SlotTime
		// AA:BB being the end SlotTime
		timeSlotText := data.Rows[row].Columns[0].RawText

		// Get slot start and end times
		result := HHMMTimeRegex.FindAll([]byte(timeSlotText), 2)

		if len(result) != 2 {
			err = fmt.Errorf("time slot test mismatch on table row %v, missing time range (invalid text received: '%v')", row, timeSlotText)
			return
		}

		// Parse XX:YY
		t0, err := daytime.Parse(string(result[0]))
		if err != nil {
			return nil, err
		}

		// Parse AA:BB
		t1, err := daytime.Parse(string(result[1]))
		if err != nil {
			return nil, err
		}

		// Store time slot
		slotTimeRange := timerange.TimeRange{Start: t0, End: t1}
		slotTimes = append(slotTimes, slotTimeRange)
	}

	return
}

func parseWeekdaySlotsFromTable(data *TableData, weekdayIndex int) (slots []ReservationSlot, err error) {
	slots = make([]ReservationSlot, len(data.Rows)) // Prebuild slice with correct length

	for row := 0; row < len(data.Rows); row++ {
		dataRow := &data.Rows[row]
		slots[row], err = parseWeekdaySlotFromTableRow(dataRow, weekdayIndex)
		if err != nil {
			return
		}
	}

	return
}

func parseWeekdaySlotFromTableRow(row *TableDataRow, weekdayIndex int) (slot ReservationSlot, err error) {
	columnIndex := weekdayIndex + 1

	rawSlot := row.Columns[columnIndex]
	rawText := rawSlot.RawText
	urlText := rawSlot.URL

	// Empty slots are possible
	if strings.TrimSpace(rawText) == "" {
		return
	}

	/**
	The rawText may be sorrounded by either <a> or <span> tags in formats
	2 and 4 respectively (see comments below for more info). If there is
	such a surrounding tag, we must remove it.
	**/
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(rawText))
	if sel := doc.Find("html body").Children().First(); sel.Children().Length() != 0 {
		rawText, _ = sel.Html()
	}

	/**
	If the slot is freely accessible, the rawText will be in the format
	<span [some attirbutes...]> [RAW_TEXT] </span>
	**/
	// removeSpan := doc.Find()

	/**
	From this point forward, rawText will have a known format.
	We can encounter one  of several different possibilities:

	1 - The slot has not been reserved and has availability: 	"[ID]<br/>Solo Socios<br/>[NUM] libres"
	2 - The slot has not been reserved and has no availability: "[ID]<br/>Solo Socios<br/>Completo"
	3 - The slot has been reserved: 							"[ID]<br/>Ya inscrito"
	4 - The slot is freely accessible:							"[ID]<br>Acceso libre<br>Solo Socios"

	NOTE:
		[ID] 	represents the slot's ID and follows the regex '[A-Z]{2}[0-9]{3}'
		[NUM]	represents the slot's availability and follows the regex '[0-9]+'

	**/
	splitRawText := strings.Split(rawText, "<br/>")
	if len(splitRawText) < 2 {
		err = fmt.Errorf("invalid format, slot raw text has unknown format '%s'", rawText)
		return
	}

	// Line 0 will always be the slot's ID
	name := splitRawText[0]

	// Make URL object from urlText (if possible)
	var url *url.URL
	if urlText != "" {
		url, _ = url.Parse("https://intranet.upv.es/pls/soalu/" + urlText) // We can safely ignore this error
	}

	// If there is a URL associated with the slot, we can expect format 2 to appear
	var availability int64 = 0 // Default availability of 0
	if urlText != "" {
		numString := string(NumberRegex.Find([]byte(splitRawText[2])))
		if numString == "" {
			err = fmt.Errorf("invalid format, number expected but none found in string '%s'", splitRawText[2])
			return
		}
		availability, _ = strconv.ParseInt(numString, 10, 16) // We can safely ignore this error
	}

	// If the rawText contains matches this regex, we can assume the user is registered in this slot
	regstered := AlreadyRegisteredRegex.Match([]byte(rawText))

	// Build slot and return
	slot = ReservationSlot{name, availability, regstered, url}
	return
}
