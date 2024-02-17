package phishin

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

func TestFormatURL(t *testing.T) {
	dummy := "dummy"
	endpoint := "songs"
	t.Run("just endpoint", func(t *testing.T) {
		t.Parallel()
		c := NewClient(dummy, os.Stdout)
		got := c.FormatURL(endpoint)
		want := fmt.Sprintf("https://phish.in/api/v1/%s", endpoint)
		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	})
	t.Run("add query", func(t *testing.T) {
		t.Parallel()
		c := NewClient(dummy, os.Stdout)
		query := "harry-hood"
		c.Query = query
		got := c.FormatURL(endpoint)
		want := fmt.Sprintf("https://phish.in/api/v1/%s/%s", endpoint, query)
		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	})
	t.Run("add parameters", func(t *testing.T) {
		t.Parallel()
		c := NewClient(dummy, os.Stdout)
		c.Parameters = []string{"per_page=3", "page=1", "sort_attr=date", "sort_dir=desc"}
		got := c.FormatURL(endpoint)
		want := fmt.Sprintf("https://phish.in/api/v1/%s?per_page=3&page=1&sort_attr=date&sort_dir=desc", endpoint)
		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	})
	t.Run("can't have query and parameters", func(t *testing.T) {
		t.Parallel()
		c := NewClient(dummy, os.Stdout)
		query := "harry-hood"
		c.Query = query
		c.Parameters = []string{"per_page=3", "page=1", "sort_attr=date", "sort_dir=desc"}
		got := c.FormatURL(endpoint)
		want := fmt.Sprintf("https://phish.in/api/v1/%s/%s", endpoint, query)
		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	})
}

func TestGetEras(t *testing.T) {
	t.Parallel()
	ts := httptest.NewTLSServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "../testdata/eras.json")
		}))
	defer ts.Close()
	c := NewClient("dummy", os.Stdout)
	c.BaseURL = ts.URL
	c.HTTPClient = ts.Client()
	want := ErasOutput{
		One: []string{"1983-1987","1988","1989","1990","1991","1992","1993","1994","1995","1996","1997","1998","1999","2000"},
		Two: []string{"2002","2003","2004"},
		Three: []string{"2009","2010","2011","2012","2013","2014","2015","2016","2017","2018","2019","2020"},
		Four: []string{"2021","2022","2023"},
	}
	ctx := context.Background()
	got, err := c.getEras(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}	
}

func TestGetAndPrintErasText(t *testing.T) {
	t.Parallel()
	buf := &bytes.Buffer{}
	ts := httptest.NewTLSServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "../testdata/eras.json")
		}))
	defer ts.Close()
	c := NewClient("dummy", buf)
	c.BaseURL = ts.URL
	c.HTTPClient = ts.Client()
	want := `Eras
1.0: 1983-1987, 1988, 1989, 1990, 1991, 1992, 1993, 1994, 1995, 1996, 1997, 1998, 1999, 2000
2.0: 2002, 2003, 2004
3.0: 2009, 2010, 2011, 2012, 2013, 2014, 2015, 2016, 2017, 2018, 2019, 2020
4.0: 2021, 2022, 2023
`
	ctx := context.Background()
	err := c.getAndPrintEras(ctx)
	if err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
}

func TestGetAndPrintErasJSON(t *testing.T) {
	t.Parallel()
	buf := &bytes.Buffer{}
	ts := httptest.NewTLSServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "../testdata/simple_eras.json")
		}))
	defer ts.Close()
	c := NewClient("dummy", buf)
	c.Output = buf
	c.BaseURL = ts.URL
	c.HTTPClient = ts.Client()
	c.PrintJSON = true
	want := `{
  "1.0": [
    "1992",
    "1993",
    "1994",
    "1995",
    "1996"
  ],
  "2.0": null,
  "3.0": null,
  "4.0": null
}
`
	ctx := context.Background()
	err := c.getAndPrintEras(ctx)
	if err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
}

func TestGetEra(t *testing.T) {
	t.Parallel()
	ts := httptest.NewTLSServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "../testdata/era.json")
		}))
	defer ts.Close()
	c := NewClient("dummy", os.Stdout)
	c.BaseURL = ts.URL
	c.HTTPClient = ts.Client()
	c.Query = "3.0"
	want := EraOutput{
		EraName: "3.0",
		Years: []string{"2009","2010","2011","2012","2013","2014","2015","2016","2017","2018","2019","2020"},
	}
	ctx := context.Background()
	got, err := c.getEra(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}	
}

func TestGetAndPrintEraText(t *testing.T) {
	t.Parallel()
	buf := &bytes.Buffer{}
	ts := httptest.NewTLSServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "../testdata/era.json")
		}))
	defer ts.Close()
	c := NewClient("dummy", buf)
	c.Output = buf
	c.BaseURL = ts.URL
	c.Query = "3.0"
	c.HTTPClient = ts.Client()
	want := `Era 3.0:
2009, 2010, 2011, 2012, 2013, 2014, 2015, 2016, 2017, 2018, 2019, 2020
`
	ctx := context.Background()
	err := c.getAndPrintEra(ctx)
	if err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
}

func TestGetAndPrintEraJSON(t *testing.T) {
	t.Parallel()
	buf := &bytes.Buffer{}
	ts := httptest.NewTLSServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "../testdata/era.json")
		}))
	defer ts.Close()
	c := NewClient("dummy", buf)
	c.Output = buf
	c.BaseURL = ts.URL
	c.HTTPClient = ts.Client()
	c.Query = "3.0"
	c.PrintJSON = true
	want := `{
  "era": "3.0",
  "years": [
    "2009",
    "2010",
    "2011",
    "2012",
    "2013",
    "2014",
    "2015",
    "2016",
    "2017",
    "2018",
    "2019",
    "2020"
  ]
}
`
	ctx := context.Background()
	err := c.getAndPrintEra(ctx)
	if err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
}

func TestGetYears(t *testing.T) {
	t.Parallel()
	ts := httptest.NewTLSServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "../testdata/years.json")
		}))
	defer ts.Close()
	c := NewClient("dummy", os.Stdout)
	c.BaseURL = ts.URL
	c.HTTPClient = ts.Client()
	want := YearsOutput{
		Years: []Year{
			{
				Date: "1983-1987",
				ShowCount: 34,
			},
			{
				Date: "1988",
				ShowCount: 44,
			},
			{
				Date: "1989",
				ShowCount: 64,
			},
		},
	}

	ctx := context.Background()
	got, err := c.getYears(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}	
}

func TestGetAndPrintYearsText(t *testing.T) {
	t.Parallel()
	buf := &bytes.Buffer{}
	ts := httptest.NewTLSServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "../testdata/years.json")
		}))
	defer ts.Close()
	c := NewClient("dummy", buf)
	c.Output = buf
	c.BaseURL = ts.URL
	c.HTTPClient = ts.Client()
	want := `Years:     Show Count:
1983-1987  34
1988       44
1989       64
`
	ctx := context.Background()	
	err := c.getAndPrintYears(ctx)
	if err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
}

func TestGetYear(t *testing.T) {
	t.Parallel()
	ts := httptest.NewTLSServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "../testdata/year.json")
		}))
	defer ts.Close()
	c := NewClient("dummy", os.Stdout)
	c.BaseURL = ts.URL
	c.HTTPClient = ts.Client()
	// simplified version so we can spot-check a few values
	want := YearOutput{
		ConcertInfo: []ConcertInfo{
			{
				Date: "1994-04-04",
				Sbd: true,
				VenueName: "The Flynn Theatre",
				Venue: Venue{
					Location: "Burlington, VT",
				},
			},
			{
				Date: "1994-04-05",
				Sbd: false,
				VenueName: "The Metropolis",
				Venue: Venue{
					Location: "Montréal, Québec, Canada",
				},
			},
			{
				Date: "1994-04-06",
				Sbd: false,
				VenueName: "Concert Hall",
				Venue: Venue{
					Location: "Toronto, Ontario, Canada",
				},
			},
		},
	}
	ctx := context.Background()
	got, err := c.getYear(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.ConcertInfo) != len(want.ConcertInfo) {
		t.Errorf("got %d ConcertInfo want %d ConcertInfo", len(got.ConcertInfo), len(want.ConcertInfo))
	}	
	for i := range got.ConcertInfo {
		if got.ConcertInfo[i].Sbd != want.ConcertInfo[i].Sbd {
			t.Errorf("got %v want %v", got.ConcertInfo[i].Sbd, want.ConcertInfo[i].Sbd)
		}
		if got.ConcertInfo[i].Venue.Location != want.ConcertInfo[i].Venue.Location {
			t.Errorf("got %q want %q", got.ConcertInfo[i].Venue.Location, want.ConcertInfo[i].Venue.Location)
		}
	}
}

func TestGetAndPrintYearText(t *testing.T) {
	t.Parallel()
	ts := httptest.NewTLSServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "../testdata/year.json")
		}))
	defer ts.Close()
	t.Run("non-verbose", func(t *testing.T) {
		buf := &bytes.Buffer{}
		c := NewClient("dummy", buf)
		c.Output = buf
		c.BaseURL = ts.URL
		c.HTTPClient = ts.Client()
		want := `Date:       Venue:             Location:
1994-04-04  The Flynn Theatre  Burlington, VT
1994-04-05  The Metropolis     Montréal, Québec, Canada
1994-04-06  Concert Hall       Toronto, Ontario, Canada
`
		ctx := context.Background()
		err := c.getAndPrintYear(ctx)
		if err != nil {
			t.Fatal(err)
		}
		got := buf.String()
		if got != want {
			t.Errorf("got \n%s want \n%s", got, want)
		}
	})
	t.Run("verbose", func(t *testing.T) {
		buf := &bytes.Buffer{}
		c := NewClient("dummy", buf)
		c.Output = buf
		c.BaseURL = ts.URL
		c.HTTPClient = ts.Client()
		c.Verbose = true
		want := `Date:       Venue:             Location:                 Duration:  Soundboard:  Remastered:
1994-04-04  The Flynn Theatre  Burlington, VT            2h 40m     yes          
1994-04-05  The Metropolis     Montréal, Québec, Canada  2h 22m                  
1994-04-06  Concert Hall       Toronto, Ontario, Canada  2h 19m                  
`
		ctx := context.Background()
		err := c.getAndPrintYear(ctx)
		if err != nil {
			t.Fatal(err)
		}
		got := buf.String()
		if got != want {
			t.Errorf("got \n%s want \n%s", got, want)
		}
	})
}