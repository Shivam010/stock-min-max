package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	commodities = map[string]string{
		"gold":     "Gold",
		"silver":   "Silver",
		"copper":   "Copper",
		"zinc":     "Zinc",
		"crudeoil": "Crudeoil",
		"cardamom": "Cardamom",
		"cotton":   "Cotton",
		"lead":     "Lead",
	}
	endPoint = "https://www.mcxindia.com/BackPage.aspx/GetGraphForScrip"

	// interval of size Slot min
	listOfSlots = []int{1, 5, 15, 30, 60}
)

func reader(com string) *bytes.Reader {
	return bytes.NewReader([]byte(`{"Commodity":"` + com + `"}`))
}

type mcxResponse struct {
	D struct {
		Data struct {
			Expiry    string  `json:"Expiry"`
			MaxDate   float64 `json:"MaxDate"`
			MinDate   float64 `json:"MinDate"`
			Commodity string  `json:"ScripName"`
			Values    []struct {
				Time  float64 `json:"x"`
				Price float64 `json:"y"`
			} `json:"IntradayGraphPlot"`
		}
	}
}

// Final Result object
type Result struct {
	// From represents whether the result is from Maxima
	// or Minima considering interval of size Slot min
	From string `json:"type"`
	Slot int    `json:"interval"`

	// Retrieved information - self explanatory
	Time    time.Time `json:"time"`
	Price   float64   `json:"price"`
	Message string    `json:"message"`
}

// Final high low values
type HighLowCap struct {
	Min         float64 `json:"min"`
	Max         float64 `json:"max"`
	BuyCap      float64 `json:"buy_cap"`
	SellCap     float64 `json:"sell_cap"`
	Description string  `json:"description"`
}

type Values struct {
	Minima *Result     `json:"minima"`
	Maxima *Result     `json:"maxima"`
	Result *HighLowCap `json:"result"`
}

type FinalResult struct {
	Summary   map[string]interface{} `json:"summary"`
	Average   *HighLowCap            `json:"average"`
	Intervals map[string]*Values     `json:"intervals"`
}

func epochToTime(_u float64) time.Time {
	u := int64(_u)
	// u is Timestamp in milliseconds
	return time.Unix(u/1000, u%1000*1000000).In(loc)
}

// Minima returns the nearest time(divided in `slot` seconds) at which
// a minimum price of the stock is obtained, along a message (if any)
func Minima(r *mcxResponse, slot int) *Result {
	if slot == 0 {
		return &Result{
			From:    "Minima",
			Slot:    slot,
			Time:    time.Time{},
			Price:   0,
			Message: "invalid slot value",
		}
	}
	rLen := len(r.D.Data.Values)
	if rLen == 0 {
		return &Result{
			From:    "Minima",
			Slot:    slot,
			Time:    time.Time{},
			Price:   0,
			Message: "Not enough data",
		}
	}
	i := rLen - 1
	val := r.D.Data.Values[i] // current time

	// at least `slot` amount of values should be obtained
	if i < slot {
		return &Result{
			From:    "Minima",
			Slot:    slot,
			Time:    epochToTime(val.Time),
			Price:   val.Price,
			Message: "Not enough data",
		}
	}
	i -= slot
	pre := r.D.Data.Values[i] // past time

	// Now, going in past

	// up -> past values are higher
	for val.Price <= pre.Price {
		// NOTE: This is redundant iterative loop as it will be checked in `maxima` also
		val = pre
		if i == 0 {
			// As first value is maxima so, the last value will be treated as minima, in this case
			pre = r.D.Data.Values[rLen-1]
			return &Result{
				From:    "Minima",
				Slot:    slot,
				Time:    epochToTime(pre.Time),
				Price:   pre.Price,
				Message: "No minima is obtained currently",
			}
		}
		i -= slot
		if i < 0 {
			i = 0
		}
		pre = r.D.Data.Values[i]
	}

	// down -> searching for nearest lowest value
	for val.Price >= pre.Price {
		val = pre
		if i == 0 {
			return &Result{
				From:    "Minima",
				Slot:    slot,
				Time:    epochToTime(pre.Time),
				Price:   pre.Price,
				Message: "Opened at lowest",
			}
		}
		i -= slot
		if i < 0 {
			i = 0
		}
		pre = r.D.Data.Values[i]
	}
	return &Result{
		From:    "Minima",
		Slot:    slot,
		Time:    epochToTime(val.Time),
		Price:   val.Price,
		Message: "",
	}
}

// Maxima returns the nearest time(divided in `slot` seconds) at which
// a maximum price of the stock is obtained, along a message (if any)
func Maxima(r *mcxResponse, slot int) *Result {
	if slot == 0 {
		return &Result{
			From:    "Maxima",
			Slot:    slot,
			Time:    time.Time{},
			Price:   0,
			Message: "invalid slot value",
		}
	}
	rLen := len(r.D.Data.Values)
	if rLen == 0 {
		return &Result{
			From:    "Maxima",
			Slot:    slot,
			Time:    time.Time{},
			Price:   0,
			Message: "Not enough data",
		}
	}
	i := rLen - 1
	val := r.D.Data.Values[i] // current time

	// at least `slot` amount of values should be obtained
	if i < slot {
		return &Result{
			From:    "Maxima",
			Slot:    slot,
			Time:    epochToTime(val.Time),
			Price:   val.Price,
			Message: "Not enough data",
		}
	}
	i -= slot
	pre := r.D.Data.Values[i] // past time

	// Now, going in past

	// up -> past values are lower
	for val.Price >= pre.Price {
		// NOTE: This is redundant iterative loop as it will be checked in `minima` also
		val = pre
		if i == 0 {
			// As first value is minima so, the last value will be treated as maxima, in this case
			pre = r.D.Data.Values[rLen-1]
			return &Result{
				From:    "Maxima",
				Slot:    slot,
				Time:    epochToTime(pre.Time),
				Price:   pre.Price,
				Message: "No maxima is obtained currently",
			}
		}
		i -= slot
		if i < 0 {
			i = 0
		}
		pre = r.D.Data.Values[i]
	}

	// down -> searching for nearest lowest value
	for val.Price <= pre.Price {
		val = pre
		if i == 0 {
			return &Result{
				From:    "Maxima",
				Slot:    slot,
				Time:    epochToTime(pre.Time),
				Price:   pre.Price,
				Message: "Opened at Highest",
			}
		}
		i -= slot
		if i < 0 {
			i = 0
		}
		pre = r.D.Data.Values[i]
	}
	return &Result{
		From:    "Maxima",
		Slot:    slot,
		Time:    epochToTime(val.Time),
		Price:   val.Price,
		Message: "",
	}
}

func extract(b []byte) *response {
	raw := &mcxResponse{}
	if err := json.Unmarshal(b, &raw); err != nil {
		// if some error obtained from MCX server
		mp := map[string]interface{}{}
		Ignore(json.Unmarshal(b, &mp))
		return &response{
			Error: "Sorry, something went wrong :(",
			Data:  mp,
		}
	}

	chMin := make(chan *Result, 8)
	chMax := make(chan *Result, 8)

	for _, slot := range listOfSlots {
		// `slot` Minute

		go func(s int) {
			chMin <- Minima(raw, s)
		}(slot)

		go func(s int) {
			chMax <- Maxima(raw, s)
		}(slot)
	}

	//wg.Wait()

	final := &FinalResult{
		Summary: map[string]interface{}{
			"expiry_date": raw.D.Data.Expiry,
			"to":          epochToTime(raw.D.Data.MaxDate),
			"from":        epochToTime(raw.D.Data.MinDate),
			"commodity":   raw.D.Data.Commodity,
		},
		Intervals: make(map[string]*Values, len(listOfSlots)),
		Average:   &HighLowCap{},
	}

	wg := &sync.WaitGroup{}
	minL := make(map[int]*Result, len(listOfSlots))
	maxL := make(map[int]*Result, len(listOfSlots))

	wg.Add(2)
	// list results from channel of Maxima'
	go func() {
		for i := 0; i < len(listOfSlots); i++ {
			c := <-chMax
			maxL[c.Slot] = c
		}
		close(chMax)
		wg.Done()
	}()

	// list results from channel of Minima'
	go func() {
		for i := 0; i < len(listOfSlots); i++ {
			c := <-chMin
			minL[c.Slot] = c
		}
		close(chMin)
		wg.Done()
	}()

	// wait till the process finishes
	wg.Wait()

	// Now, calculate Average results
	min := math.MaxFloat64
	max := math.SmallestNonzeroFloat64
	for _, slot := range listOfSlots {
		if minL[slot].Price < min {
			min = minL[slot].Price
		}
		if maxL[slot].Price > max {
			max = maxL[slot].Price
		}
		c := (maxL[slot].Price - minL[slot].Price) / 4
		// fill
		final.Intervals[fmt.Sprintf("%v minute", slot)] = &Values{
			Minima: minL[slot],
			Maxima: maxL[slot],
			Result: &HighLowCap{
				Min:         minL[slot].Price,
				Max:         maxL[slot].Price,
				BuyCap:      minL[slot].Price + c,
				SellCap:     maxL[slot].Price - c,
				Description: `Max: A ; Min: B ; SellCap: X ; BuyCap: Y`,
			},
		}
	}
	c := (max - min) / 4
	final.Average = &HighLowCap{
		Min:         min,
		Max:         max,
		BuyCap:      min + c,
		SellCap:     max - c,
		Description: `Max: A ; Min: B ; SellCap: X ; BuyCap: Y`,
	}
	return &response{Data: final}
}

func process(com string) *response {

	com = commodities[strings.ToLower(com)]
	if com == "" {
		return &response{Error: "invalid commodity"}
	}

	// extract data
	req, err := http.NewRequest(http.MethodPost, endPoint, reader(com))
	if err != nil {
		return &response{Error: err.Error()}
	}
	ctx, cancel := context.WithTimeout(req.Context(), 1*time.Minute)
	defer cancel()
	req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Referer", "https://www.mcxindia.com/home")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return &response{Error: err.Error()}
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return &response{Error: err.Error()}
	}

	return extract(b)
}
