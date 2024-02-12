package phishin

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestFormatURL(t *testing.T) {
	dummy := "dummy"
	endpoint := "songs"
	t.Run("just endpoint", func(t *testing.T) {
		t.Parallel()
		c := NewClient(dummy)
		got := c.FormatURL(endpoint)
		want := fmt.Sprintf("https://phish.in/api/v1/%s", endpoint)
		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	})
	t.Run("add query", func(t *testing.T) {
		t.Parallel()
		c := NewClient(dummy)
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
		c := NewClient(dummy)
		c.Parameters = []string{"per_page=3", "page=1", "sort_attr=date", "sort_dir=desc"}
		got := c.FormatURL(endpoint)
		want := fmt.Sprintf("https://phish.in/api/v1/%s?per_page=3&page=1&sort_attr=date&sort_dir=desc", endpoint)
		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	})
	t.Run("can't have query and parameters", func(t *testing.T) {
		t.Parallel()
		c := NewClient(dummy)
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
	c := NewClient("dummy")
	c.BaseURL = ts.URL
	c.HTTPClient = ts.Client()
	want := ErasOutput{
		One: []string{"1983-1987","1988","1989","1990","1991","1992","1993","1994","1995","1996","1997","1998","1999","2000"},
		Two: []string{"2002","2003","2004"},
		Three: []string{"2009","2010","2011","2012","2013","2014","2015","2016","2017","2018","2019","2020"},
		Four: []string{"2021","2022","2023"},
	}
	got, err := c.getEras()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}	
}

func TestGetAndPrintErasText(t *testing.T) {
	buf := &bytes.Buffer{}
	t.Parallel()
	ts := httptest.NewTLSServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "../testdata/eras.json")
		}))
	defer ts.Close()
	c := NewClient("dummy")
	c.Output = buf
	c.BaseURL = ts.URL
	c.HTTPClient = ts.Client()
	want := `Eras
1.0: 1983-1987, 1988, 1989, 1990, 1991, 1992, 1993, 1994, 1995, 1996, 1997, 1998, 1999, 2000
2.0: 2002, 2003, 2004
3.0: 2009, 2010, 2011, 2012, 2013, 2014, 2015, 2016, 2017, 2018, 2019, 2020
4.0: 2021, 2022, 2023
`
	err := c.getAndPrintEras()
	if err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
}

func TestGetAndPrintErasJSON(t *testing.T) {
	buf := &bytes.Buffer{}
	t.Parallel()
	ts := httptest.NewTLSServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "../testdata/simple_eras.json")
		}))
	defer ts.Close()
	c := NewClient("dummy")
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
	err := c.getAndPrintEras()
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
	c := NewClient("dummy")
	c.BaseURL = ts.URL
	c.HTTPClient = ts.Client()
	c.Query = "3.0"
	want := EraOutput{
		Era: "3.0",
		EraList: []string{"2009","2010","2011","2012","2013","2014","2015","2016","2017","2018","2019","2020"},
	}
	got, err := c.getEra()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}	
}

func TestGetAndPrintEraText(t *testing.T) {
	buf := &bytes.Buffer{}
	t.Parallel()
	ts := httptest.NewTLSServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "../testdata/era.json")
		}))
	defer ts.Close()
	c := NewClient("dummy")
	c.Output = buf
	c.BaseURL = ts.URL
	c.Query = "3.0"
	c.HTTPClient = ts.Client()
	want := `Era 3.0:
2009, 2010, 2011, 2012, 2013, 2014, 2015, 2016, 2017, 2018, 2019, 2020
`
	err := c.getAndPrintEra()
	if err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
}

func TestGetAndPrintEraJSON(t *testing.T) {
	buf := &bytes.Buffer{}
	t.Parallel()
	ts := httptest.NewTLSServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "../testdata/era.json")
		}))
	defer ts.Close()
	c := NewClient("dummy")
	c.Output = buf
	c.BaseURL = ts.URL
	c.HTTPClient = ts.Client()
	c.Query = "3.0"
	c.PrintJSON = true
	want := `{
  "Era": "3.0",
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
	err := c.getAndPrintEra()
	if err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
}