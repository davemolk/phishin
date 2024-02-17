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
	url := c.FormatURL("eras")
	got, err := c.getEras(ctx, url)
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
	url := c.FormatURL("eras")
	err := c.getAndPrintEras(ctx, url)
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
	url := c.FormatURL("eras")
	err := c.getAndPrintEras(ctx, url)
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
	url := c.FormatURL("eras")
	got, err := c.getEra(ctx, url)
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
	url := c.FormatURL("eras")
	err := c.getAndPrintEra(ctx, url)
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
	url := c.FormatURL("eras")
	err := c.getAndPrintEra(ctx, url)
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
	url := c.FormatURL("years")
	got, err := c.getYears(ctx, url)
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
	url := c.FormatURL("years")
	err := c.getAndPrintYears(ctx, url)
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
	want := ShowsOutput{
		Shows: []Show{
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
	url := c.FormatURL("years")
	got, err := c.getYear(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Shows) != len(want.Shows) {
		t.Errorf("got %d shows want %d shows", len(got.Shows), len(want.Shows))
	}	
	for i := range got.Shows {
		if got.Shows[i].Sbd != want.Shows[i].Sbd {
			t.Errorf("got %v want %v", got.Shows[i].Sbd, want.Shows[i].Sbd)
		}
		if got.Shows[i].Venue.Location != want.Shows[i].Venue.Location {
			t.Errorf("got %q want %q", got.Shows[i].Venue.Location, want.Shows[i].Venue.Location)
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
		url := c.FormatURL("years")
		err := c.getAndPrintYear(ctx, url)
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
		url := c.FormatURL("years")
		err := c.getAndPrintYear(ctx, url)
		if err != nil {
			t.Fatal(err)
		}
		got := buf.String()
		if got != want {
			t.Errorf("got \n%s want \n%s", got, want)
		}
	})
}

func TestGetShows(t *testing.T) {
	t.Parallel()
	ts := httptest.NewTLSServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "../testdata/shows.json")
		}))
	defer ts.Close()
	c := NewClient("dummy", os.Stdout)
	c.BaseURL = ts.URL
	c.HTTPClient = ts.Client()
	want := ShowsOutput{
		Shows: []Show{
			{
				Date: "1990-04-05",
				Sbd: true,
				VenueName: "J.J. McCabe's",
				Venue: Venue{
					Location: "Boulder, CO",
				},
			},
		},
	}
	ctx := context.Background()
	url := c.FormatURL("shows")
	got, err := c.getYear(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Shows) != len(want.Shows) {
		t.Errorf("got %d shows want %d shows", len(got.Shows), len(want.Shows))
	}	
	for i := range got.Shows {
		if got.Shows[i].Sbd != want.Shows[i].Sbd {
			t.Errorf("got %v want %v", got.Shows[i].Sbd, want.Shows[i].Sbd)
		}
		if got.Shows[i].Venue.Location != want.Shows[i].Venue.Location {
			t.Errorf("got %q want %q", got.Shows[i].Venue.Location, want.Shows[i].Venue.Location)
		}
	}
}

func TestGetAndPrintShowsText(t *testing.T) {
	t.Parallel()
	ts := httptest.NewTLSServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "../testdata/shows.json")
		}))
	defer ts.Close()
	t.Run("non-verbose", func(t *testing.T) {
		buf := &bytes.Buffer{}
		c := NewClient("dummy", buf)
		c.Output = buf
		c.BaseURL = ts.URL
		c.HTTPClient = ts.Client()
		want := `Date:       Venue:         Location:
1990-04-05  J.J. McCabe's  Boulder, CO
`
		ctx := context.Background()
		url := c.FormatURL("shows")
		err := c.getAndPrintYear(ctx, url)
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
		want := `Date:       Venue:         Location:    Duration:  Soundboard:  Remastered:
1990-04-05  J.J. McCabe's  Boulder, CO  2h 27m     yes          
`
		ctx := context.Background()
		url := c.FormatURL("shows")
		err := c.getAndPrintYear(ctx, url)
		if err != nil {
			t.Fatal(err)
		}
		got := buf.String()
		if got != want {
			t.Errorf("got \n%s want \n%s", got, want)
		}
	})
}