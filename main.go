// cmd subsurface implements simple transformations on top of
// subsurface xml files.
//
// Note that currently only a subset of the subsurface file
// format is supported and only one transformation is implemented
// (timeshift dives). But if you need to manipulate subsurface xml
// files this is as good a place as any to start.
//
// See : http://subsurface.hohndel.org/
package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
	"time"
)

type DiveLog struct {
	XMLName  xml.Name `xml:"divelog"`
	Program  string   `xml:"program,attr"`
	Version  string   `xml:"version,attr"`
	Settings string   `xml:"settings,omitempty"`
	Dives    []*Dive  `xml:"dives>dive"`
}
type Dive struct {
	Number       int           `xml:"number,attr,omitempty"`
	Date         string        `xml:"date,attr,omitempty"`
	TimeOfDay    string        `xml:"time,attr,omitempty"` // offset from midnight on Date
	RawDuration  string        `xml:"duration,attr,omitempty"`
	Rating       int           `xml:"rating,attr,omitempty"`
	Visibility   int           `xml:"visibility,attr,omitempty"`
	Location     *Location     `xml:"location,omitempty"`
	ComputerDive *ComputerDive `xml:"divecomputer,omitempty"`
	Notes        string        `xml:"notes"`
	Suit         string        `xml:"suit"`
	Buddy        string        `xml:"buddy"`
}

type Location struct {
	GPS  string `xml:"gps,attr,omitempty"`
	Name string `xml:",chardata"`
}

type ComputerDive struct {
	Model     string `xml:"model,attr"`
	DeviceID  string `xml:"deviceid,attr"`
	DiveID    string `xml:"diveid,attr"`
	Date      string `xml:"date,attr,omitempty"`
	TimeOfDay string `xml:"time,attr,omitempty"` // offset from midnight on Date
	Note      string `xml:"notes,omitempty"`
	DepthStat *struct {
		Max  string `xml:"max,attr"`
		Mean string `xml:"mean,attr"`
	} `xml:"depth,omitempty"`
	Temperature *struct {
		Water string `xml:"water,attr"`
	} `xml:"temperature,omitempty"`
	Samples []*struct {
		Time  string `xml:"time,attr,omitempty"`
		Depth string `xml:"depth,attr,omitempty"`
		Temp  string `xml:"temp,attr,omitempty"`
	} `xml:"sample,omitempty"`
}

const (
	ssDate = "2006-01-02"
	ssTime = "15:04:05"
)

func (d *Dive) Time() (time.Time, error) {
	date, err := time.Parse(ssDate, d.Date)
	if err != nil {
		return date, err
	}

	tim, err := time.Parse(ssTime, d.TimeOfDay)
	if err != nil {
		return date, err
	}

	y, m, day := date.Date()
	return tim.AddDate(y, int(m), day), nil
}

func (d *Dive) Duration() time.Duration {
	dur := strings.Fields(d.RawDuration)[0]
	dur = strings.Replace(dur, ":", "m", 2) + "s"
	dd, _ := time.ParseDuration(dur)
	return dd
}

func (d *Dive) SetTime(t time.Time) {
	d.Date = t.Format(ssDate)
	d.TimeOfDay = t.Format(ssTime)
}

func abort(msg string) {
	fmt.Fprint(os.Stderr, msg)
	os.Exit(1)
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [opts] [infile] [outfile]\n", path.Base(os.Args[0]))
	flag.PrintDefaults()
}

func main() {
	var (
		shiftDur = flag.Duration("d", 0, "move each dive by this duration")
	)

	flag.Parse()

	var (
		f   *os.File
		out *os.File
	)
	args := flag.Args()
	if len(args) < 1 {
		f = os.Stdin
	} else {
		var err error
		f, err = os.Open(args[0])
		if err != nil {
			abort(err.Error())
		}
		defer f.Close()
	}

	if len(args) < 2 {
		out = os.Stdout
	} else {
		var err error
		out, err = os.Create(args[1])
		if err != nil {
			abort(err.Error())
		}
		defer out.Close()
	}

	dec := xml.NewDecoder(f)
	s := &DiveLog{}
	err := dec.Decode(s)
	if err != nil {
		abort(err.Error())
	}

	for _, dive := range s.Dives {
		t, err := dive.Time()
		if err != nil {
			continue // skip
		}

		dive.SetTime(t.Add(*shiftDur))
	}

	enc := xml.NewEncoder(out)
	enc.Indent("", "   ")
	err = enc.Encode(s)
	if err != nil {
		abort(err.Error())
	}

}
