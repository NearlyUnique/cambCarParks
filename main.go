package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type (
	// CarPark data
	CarPark struct {
		Name   string
		Spaces int
	}
	// RawHTML .
	RawHTML struct {
		XMLName xml.Name `xml:"r"`
		Name    []string `xml:"h2>a"`
		Value   []string `xml:"p>strong"`
	}
)

func main() {
	t := time.Now()
	defer func() { fmt.Printf("%v", time.Since(t)) }()
	data, err := GetCarParkData("https://www.cambridge.gov.uk/jdi_parking_ajax/complete")
	if err != nil {
		if v, ok := err.(*url.Error); ok {
			fmt.Printf("URL Error:\n\tTemp:%t\n\tTimeout:%t\n\tOP:%s\n\tURL:%s\n\t%v\n", v.Temporary(), v.Timeout(), v.Op, v.URL, v.Err)
		} else {
			fmt.Printf("Something went wrong %v\n", err)
		}
		os.Exit(1)
	}
	for _, cp := range data {
		fmt.Println(cp)
	}
	fmt.Println()
}
func (c CarPark) String() string {
	return fmt.Sprintf("%s,%d", c.Name, c.Spaces)
}

// GetCarParkData get from internet
func GetCarParkData(url string) ([]CarPark, error) {
	var c = &http.Client{
		Timeout: time.Second * 10,
	}
	response, err := c.Get(url)
	if err != nil {
		return []CarPark{}, err
	}
	defer response.Body.Close()
	buf, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return []CarPark{}, err
	}
	return Parse(string(buf))
}

// Parse the raw html
func Parse(data string) ([]CarPark, error) {
	if len(data) == 0 {
		return []CarPark{}, fmt.Errorf("Missing data")
	}
	if data[0] != '<' {
		return []CarPark{}, fmt.Errorf("Does not look like XML")
	}
	raw := RawHTML{}
	err := xml.Unmarshal([]byte("<r>"+data+"</r>"), &raw)
	if err != nil {
		return nil, err
	}
	cps := []CarPark{}
	for i, name := range raw.Name {
		cp := CarPark{
			Name:   name,
			Spaces: spacesFromHTML(raw.Value[i]),
		}
		cps = append(cps, cp)
	}
	return cps, nil
}
func spacesFromHTML(html string) int {
	spaces := 0
	parts := strings.Split(html, " ")
	if len(parts) == 2 && len(parts[0]) > 0 {
		spaces, _ = strconv.Atoi(parts[0])
	}
	return spaces
}
