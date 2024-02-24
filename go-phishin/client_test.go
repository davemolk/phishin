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
	query := "1994"
	path := "years"
	ts := httptest.NewTLSServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != fmt.Sprintf("/%s/%s", path, query) {
				t.Fatalf("wrong url: %s", r.URL.Path)
			}
			http.ServeFile(w, r, "../testdata/year.json")
		}))
	defer ts.Close()
	c := NewClient("dummy", os.Stdout)
	c.BaseURL = ts.URL
	c.HTTPClient = ts.Client()
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
	c.Query = query
	url := c.FormatURL(path)
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
	query := "1994"
	path := "years"
	ts := httptest.NewTLSServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != fmt.Sprintf("/%s/%s", path, query) {
				t.Fatalf("wrong url: %s", r.URL.Path)
			}
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
		c.Query = query
		url := c.FormatURL(path)
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
		c.Query = query
		url := c.FormatURL(path)
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

// func TestGetShow(t *testing.T) {
// 	t.Parallel()
// 	query := "1990-04-05"
// 	path := "shows"
// 	ts := httptest.NewTLSServer(http.HandlerFunc(
// 		func(w http.ResponseWriter, r *http.Request) {
// 			if r.URL.Path != fmt.Sprintf("/%s/%s", path, query) {
// 				t.Fatalf("wrong url: %s", r.URL.Path)
// 			}
// 			http.ServeFile(w, r, "../testdata/show.json")
// 		}))
// 	defer ts.Close()
// 	c := NewClient("dummy", os.Stdout)
// 	c.BaseURL = ts.URL
// 	c.HTTPClient = ts.Client()
// 	// grab a subset to spot-check values
// 	want := ShowOutput{
// 		Show: Show{
// 			Date: "1990-04-05",
// 			Sbd: true,
// 			VenueName: "J.J. McCabe's",
// 			Venue: Venue{
// 				Location: "Boulder, CO",
// 			},
// 		},
// 	}
// 	ctx := context.Background()
// 	c.Query = query
// 	url := c.FormatURL(path)
// 	got, err := c.getShow(ctx, url)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if got.Date != want.Date {
// 		t.Errorf("got %v want %v", got.Date, want.Date)
// 	}
// 	if got.Sbd != want.Sbd {
// 		t.Errorf("got %v want %v", got.Sbd, want.Sbd)
// 	}
// 	if got.Venue.Location != want.Venue.Location {
// 		t.Errorf("got %q want %q", got.Venue.Location, want.Venue.Location)
// 	}
// }

func TestGetAndPrintShowText(t *testing.T) {
	t.Parallel()
	query := "1990-04-05"
	path := "shows"
	ts := httptest.NewTLSServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != fmt.Sprintf("/%s/%s", path, query) {
				t.Fatalf("wrong url: %s", r.URL.Path)
			}
			http.ServeFile(w, r, "../testdata/show.json")
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
            
Set 1
Possum                   6m 48s   
Ya Mar                   7m 7s    
David Bowie              11m 23s  
Carolina                 2m 1s    
The Oh Kee Pa Ceremony   1m 45s   
Suzy Greenberg           5m 19s   
You Enjoy Myself         12m 40s  
The Lizards              10m 12s  
Fire                     4m 20s   
                         
Set 2                    
Reba                     11m 39s  
Uncle Pen                5m 14s   
Jesus Just Left Chicago  8m 10s   
AC/DC Bag                6m 23s   
Donna Lee                3m 24s   
Tweezer                  10m 0s   
Fee                      5m 14s   
Cavern                   4m 59s   
Mike's Song              6m 23s   
I Am Hydrogen            2m 19s   
Weekapaug Groove         7m 35s   
If I Only Had a Brain    3m 10s   
Contact                  6m 21s   
                         
Encore                   
Golgi Apparatus          4m 41s  
`
		ctx := context.Background()
		c.Query = query
		url := c.FormatURL(path)
		err := c.getAndPrintShow(ctx, url)
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
            
Set 1
Possum                   6m 48s   SBD
Ya Mar                   7m 7s    SBD, Tease: Theme from Bonanza by Ray Evans and Jay Livingston
David Bowie              11m 23s  SBD, Tease: Theme from Bonanza by Ray Evans and Jay Livingston, Tease: Wipe Out by The Surfaris
Carolina                 2m 1s    SBD, A Cappella
The Oh Kee Pa Ceremony   1m 45s   SBD
Suzy Greenberg           5m 19s   SBD
You Enjoy Myself         12m 40s  SBD, Tease: Flash Light by Parliament
The Lizards              10m 12s  SBD
Fire                     4m 20s   SBD
                         
Set 2                    
Reba                     11m 39s  SBD
Uncle Pen                5m 14s   SBD
Jesus Just Left Chicago  8m 10s   SBD, Guest: Dan Mosebee on harmonica
AC/DC Bag                6m 23s   SBD
Donna Lee                3m 24s   SBD
Tweezer                  10m 0s   SBD, Tease: Dave's Energy Guide
Fee                      5m 14s   SBD
Cavern                   4m 59s   SBD, Alt Lyric: "...taking turns at *stabbing* her; the brothel wife then grabbed the knife and slashed me on the tongue; I turned the blade back on the bitch and dropped her in the dung...a cushion convector, a *penile collector*..."
Mike's Song              6m 23s   SBD
I Am Hydrogen            2m 19s   SBD
Weekapaug Groove         7m 35s   SBD
If I Only Had a Brain    3m 10s   SBD
Contact                  6m 21s   SBD
                         
Encore                   
Golgi Apparatus          4m 41s  SBD
                         
want to listen?
Possum                   https://phish.in/audio/000/014/073/14073.mp3
Ya Mar                   https://phish.in/audio/000/014/074/14074.mp3
David Bowie              https://phish.in/audio/000/014/075/14075.mp3
Carolina                 https://phish.in/audio/000/014/076/14076.mp3
The Oh Kee Pa Ceremony   https://phish.in/audio/000/014/077/14077.mp3
Suzy Greenberg           https://phish.in/audio/000/014/078/14078.mp3
You Enjoy Myself         https://phish.in/audio/000/014/079/14079.mp3
The Lizards              https://phish.in/audio/000/014/080/14080.mp3
Fire                     https://phish.in/audio/000/014/081/14081.mp3
Reba                     https://phish.in/audio/000/014/082/14082.mp3
Uncle Pen                https://phish.in/audio/000/014/083/14083.mp3
Jesus Just Left Chicago  https://phish.in/audio/000/014/084/14084.mp3
AC/DC Bag                https://phish.in/audio/000/014/085/14085.mp3
Donna Lee                https://phish.in/audio/000/014/086/14086.mp3
Tweezer                  https://phish.in/audio/000/014/087/14087.mp3
Fee                      https://phish.in/audio/000/014/088/14088.mp3
Cavern                   https://phish.in/audio/000/014/089/14089.mp3
Mike's Song              https://phish.in/audio/000/014/090/14090.mp3
I Am Hydrogen            https://phish.in/audio/000/014/091/14091.mp3
Weekapaug Groove         https://phish.in/audio/000/014/092/14092.mp3
If I Only Had a Brain    https://phish.in/audio/000/014/093/14093.mp3
Contact                  https://phish.in/audio/000/014/094/14094.mp3
Golgi Apparatus          https://phish.in/audio/000/014/095/14095.mp3
`
		ctx := context.Background()
		c.Query = query
		url := c.FormatURL(path)
		err := c.getAndPrintShow(ctx, url)
		if err != nil {
			t.Fatal(err)
		}
		got := buf.String()
		if got != want {
			t.Errorf("got \n%s want \n%s", got, want)
		}
	})
}

func TestGetVenues(t *testing.T) {
	t.Parallel()
	ts := httptest.NewTLSServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "../testdata/venues.json")
		}))
	defer ts.Close()
	c := NewClient("dummy", os.Stdout)
	c.BaseURL = ts.URL
	c.HTTPClient = ts.Client()
	want := VenuesOutput{
		Venues: []VenueOutput{
			{
				Name: "The Base Lodge, Johnson State College",
				Location: "Johnson, VT",
				ShowsCount: 2,
				ShowDates: []string{"1988-03-11","1989-04-14"},
			},
			{
				Name: "The Academy",
				Location: "New York, NY",
				ShowsCount: 1,
				ShowDates: []string{"1991-07-15"},
			},
		},
	}
	ctx := context.Background()
	url := c.FormatURL("venues")
	got, err := c.getVenues(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestGetAndPrintVenuesText(t *testing.T) {
	t.Parallel()
	ts := httptest.NewTLSServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "../testdata/venues.json")
		}))
	defer ts.Close()
	buf := &bytes.Buffer{}
	c := NewClient("dummy", buf)
	c.Output = buf
	c.BaseURL = ts.URL
	c.HTTPClient = ts.Client()
	want := `Venues:                                Location:     Show Count:
The Base Lodge, Johnson State College  Johnson, VT   2
The Academy                            New York, NY  1
`
	ctx := context.Background()
	url := c.FormatURL("venues")
	err := c.getAndPrintVenues(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	if got != want {
		t.Errorf("got \n%s want \n%s", got, want)
	}
	
}

func TestGetVenue(t *testing.T) {
	t.Parallel()
	query := "the-academy"
	path := "venues"
	ts := httptest.NewTLSServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != fmt.Sprintf("/%s/%s", path, query) {
				t.Fatalf("wrong url: %s", r.URL.Path)
			}
			http.ServeFile(w, r, "../testdata/venue.json")
		}))
	defer ts.Close()
	c := NewClient("dummy", os.Stdout)
	c.BaseURL = ts.URL
	c.HTTPClient = ts.Client()
	want := VenueOutput{
		Name: "The Academy",
		Location: "New York, NY",
		ShowsCount: 1,
		ShowDates: []string{"1991-07-15"},
	}
	ctx := context.Background()
	c.Query = query
	url := c.FormatURL(path)
	got, err := c.getVenue(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestGetAndPrintVenueText(t *testing.T) {
	t.Parallel()
	query := "the-academy"
	path := "venues"
	ts := httptest.NewTLSServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != fmt.Sprintf("/%s/%s", path, query) {
				t.Fatalf("wrong url: %s", r.URL.Path)
			}
			http.ServeFile(w, r, "../testdata/venue.json")
		}))
	defer ts.Close()
	buf := &bytes.Buffer{}
	c := NewClient("dummy", buf)
	c.Output = buf
	c.BaseURL = ts.URL
	c.HTTPClient = ts.Client()
	want := `Venue:       Location:     Show Count:
The Academy  New York, NY  1
             
Show Dates
1991-07-15
`
	ctx := context.Background()
	c.Query = query
	url := c.FormatURL(path)
	err := c.getAndPrintVenue(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	if got != want {
		t.Errorf("got \n%s want \n%s", got, want)
	}
}

func TestGetTags(t *testing.T) {
	t.Parallel()
	ts := httptest.NewTLSServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "../testdata/tags.json")
		}))
	defer ts.Close()
	c := NewClient("dummy", os.Stdout)
	c.BaseURL = ts.URL
	c.HTTPClient = ts.Client()
	want := TagsOutput{
		Tags: []TagListItemOutput{
			{
				Name: "Costume",
				Description: "Musical costume sequence",
				Group: "Set Content",
			},
			{
				Name: "Audience",
				Description: "Contribution from audience during performance",
				Group: "Song Content",
			},
		},
	}
	ctx := context.Background()
	url := c.FormatURL("tags")
	got, err := c.getTags(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestGetAndPrintTagsText(t *testing.T) {
	t.Parallel()
	ts := httptest.NewTLSServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "../testdata/tags.json")
		}))
	defer ts.Close()
	buf := &bytes.Buffer{}
	c := NewClient("dummy", buf)
	c.Output = buf
	c.BaseURL = ts.URL
	c.HTTPClient = ts.Client()
	want := `Name:     Description:                                   Group:
Costume   Musical costume sequence                       Set Content
Audience  Contribution from audience during performance  Song Content
`
	ctx := context.Background()
	url := c.FormatURL("tags")
	err := c.getAndPrintTags(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	if got != want {
		t.Errorf("got \n%s want \n%s", got, want)
	}
	
}

func TestGetTag(t *testing.T) {
	t.Parallel()
	query := "jamcharts"
	path := "tags"
	ts := httptest.NewTLSServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != fmt.Sprintf("/%s/%s", path, query) {
				t.Fatalf("wrong url: %s", r.URL.Path)
			}
			http.ServeFile(w, r, "../testdata/tag.json")
		}))
	defer ts.Close()
	c := NewClient("dummy", os.Stdout)
	c.BaseURL = ts.URL
	c.HTTPClient = ts.Client()
	want := TagListItemOutput{
		Name: "Jamcharts",
		Description: "Phish.net Jam Charts selections (phish.net/jamcharts)",
		Group: "Curated Selections",
		ShowIds: []int{3},
		TrackIds: []int{1, 2},
	}
	ctx := context.Background()
	c.Query = query
	url := c.FormatURL(path)
	got, err := c.getTag(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestGetAndPrintTagText(t *testing.T) {
	t.Parallel()
	query := "jamcharts"
	path := "tags"
	ts := httptest.NewTLSServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != fmt.Sprintf("/%s/%s", path, query) {
				t.Fatalf("wrong url: %s", r.URL.Path)
			}
			http.ServeFile(w, r, "../testdata/tag.json")
		}))
	defer ts.Close()
	buf := &bytes.Buffer{}
	c := NewClient("dummy", buf)
	c.Output = buf
	c.BaseURL = ts.URL
	c.HTTPClient = ts.Client()
	want := `Name:      Description:                                           Group:
Jamcharts  Phish.net Jam Charts selections (phish.net/jamcharts)  Curated Selections

Show IDs Where Jamcharts Appears
3

Track IDs Where Jamcharts Appears
1, 2
`
	ctx := context.Background()
	c.Query = query
	url := c.FormatURL(path)
	err := c.getAndPrintTag(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	if got != want {
		t.Errorf("got \n%s want \n%s", got, want)
	}
}

func TestGetTours(t *testing.T) {
	t.Parallel()
	ts := httptest.NewTLSServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "../testdata/tours.json")
		}))
	defer ts.Close()
	c := NewClient("dummy", os.Stdout)
	c.BaseURL = ts.URL
	c.HTTPClient = ts.Client()
	want := ToursOutput{
		Tours: []Tour{
			{
				Name: "1983 Tour",
				StartsOn: "1983-12-02",
				EndsOn: "1983-12-02",
				ShowsCount: 1,
			},
			{
				Name: "1984 Tour",
				StartsOn: "1984-11-03",
				EndsOn: "1984-12-01",
				ShowsCount: 2,
			},
		},
	}
	ctx := context.Background()
	url := c.FormatURL("tours")
	got, err := c.getTours(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	for i, v := range got.Tours {
		if v.Name != want.Tours[i].Name {
			t.Errorf("got %s want %s", v.Name, want.Tours[i].Name)
		}
		if v.StartsOn != want.Tours[i].StartsOn {
			t.Errorf("got %s want %s", v.StartsOn, want.Tours[i].StartsOn)
		}
		if v.EndsOn != want.Tours[i].EndsOn {
			t.Errorf("got %s want %s", v.EndsOn, want.Tours[i].EndsOn)
		}
		if v.ShowsCount != want.Tours[i].ShowsCount {
			t.Errorf("got %d want %d", v.ShowsCount, want.Tours[i].ShowsCount)
		}
	}
}

func TestGetAndPrintToursText(t *testing.T) {
	t.Parallel()
	ts := httptest.NewTLSServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "../testdata/tours.json")
		}))
	defer ts.Close()
	buf := &bytes.Buffer{}
	c := NewClient("dummy", buf)
	c.Output = buf
	c.BaseURL = ts.URL
	c.HTTPClient = ts.Client()
	want := `Name:      Starts On:  Ends On:    Show Count:
1983 Tour  1983-12-02  1983-12-02  1
1984 Tour  1984-11-03  1984-12-01  2
`
	ctx := context.Background()
	url := c.FormatURL("tours")
	err := c.getAndPrintTours(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	if got != want {
		t.Errorf("got \n%s want \n%s", got, want)
	}
	
}

func TestGetTour(t *testing.T) {
	t.Parallel()
	query := "1985-tour"
	path := "tours"
	ts := httptest.NewTLSServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != fmt.Sprintf("/%s/%s", path, query) {
				t.Fatalf("wrong url: %s", r.URL.Path)
			}
			http.ServeFile(w, r, "../testdata/tour.json")
		}))
	defer ts.Close()
	c := NewClient("dummy", os.Stdout)
	c.BaseURL = ts.URL
	c.HTTPClient = ts.Client()
	want := TourOutput{
		Name: "1985 Tour",
		StartsOn: "1985-03-04",
		EndsOn: "1985-11-23",
		ShowsCount: 6,
		Shows: []Show{
			{
				Date: "1985-03-04",
			},
		},
	}
	ctx := context.Background()
	c.Query = query
	url := c.FormatURL(path)
	got, err := c.getTour(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestGetAndPrintTourText(t *testing.T) {
	t.Parallel()
	query := "1985-tour"
	path := "tours"
	ts := httptest.NewTLSServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != fmt.Sprintf("/%s/%s", path, query) {
				t.Fatalf("wrong url: %s", r.URL.Path)
			}
			http.ServeFile(w, r, "../testdata/tour.json")
		}))
	defer ts.Close()
	buf := &bytes.Buffer{}
	c := NewClient("dummy", buf)
	c.Output = buf
	c.BaseURL = ts.URL
	c.HTTPClient = ts.Client()
	want := `Name:      Starts On:  Ends On:    Show Count:
1985 Tour  1985-03-04  1985-11-23  6

Date:       Venue:  Location:       Duration:  Soundboard:  Remastered:
1985-03-04  Hunt's  Burlington, VT  40m 14s    yes          
`
	ctx := context.Background()
	c.Query = query
	url := c.FormatURL(path)
	err := c.getAndPrintTour(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	if got != want {
		t.Errorf("got \n%s want \n%s", got, want)
	}
}

func TestGetSongs(t *testing.T) {
	t.Parallel()
	ts := httptest.NewTLSServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "../testdata/songs.json")
		}))
	defer ts.Close()
	c := NewClient("dummy", os.Stdout)
	c.BaseURL = ts.URL
	c.HTTPClient = ts.Client()
	want := SongsOutput{
		Songs: []Song{
			{
				Title: "Billy Breathes",
				Original: true,
				Artist: "",
			},
			{
				Title: "Arc",
				Original: false,
				Artist: "Arctic Monkeys",
			},
		},
	}
	ctx := context.Background()
	url := c.FormatURL("songs")
	got, err := c.getSongs(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	for i, v := range got.Songs {
		if v.Title != want.Songs[i].Title {
			t.Errorf("got %s want %s", v.Title, want.Songs[i].Title)
		}
		if v.Original != want.Songs[i].Original {
			t.Errorf("got %v want %v", v.Original, want.Songs[i].Original)
		}
		if v.Artist != want.Songs[i].Artist {
			t.Errorf("got %s want %s", v.Artist, want.Songs[i].Artist)
		}
	}
}

func TestGetAndPrintSongsText(t *testing.T) {
	t.Parallel()
	ts := httptest.NewTLSServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "../testdata/songs.json")
		}))
	defer ts.Close()
	buf := &bytes.Buffer{}
	c := NewClient("dummy", buf)
	c.Output = buf
	c.BaseURL = ts.URL
	c.HTTPClient = ts.Client()
	want := `Title:          Phish Original:  Original Artist:  TracksCount
Billy Breathes  yes                                64
Arc                              Arctic Monkeys    0
`
	ctx := context.Background()
	url := c.FormatURL("songs")
	err := c.getAndPrintSongs(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	if got != want {
		t.Errorf("got \n%s want \n%s", got, want)
	}
}

func TestGetSong(t *testing.T) {
	t.Parallel()
	query := "david-bowie"
	path := "songs"
	ts := httptest.NewTLSServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != fmt.Sprintf("/%s/%s", path, query) {
				t.Fatalf("wrong url: %s", r.URL.Path)
			}
			http.ServeFile(w, r, "../testdata/song.json")
		}))
	defer ts.Close()
	c := NewClient("dummy", os.Stdout)
	c.BaseURL = ts.URL
	c.HTTPClient = ts.Client()
	// grab a subset to spot-check values
	want := Song{
		Title: "David Bowie",
		Alias: "",
		Original: true,
		Artist: "",
		TracksCount: 447,
	}
	ctx := context.Background()
	c.Query = query
	url := c.FormatURL(path)
	got, err := c.getSong(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	if got.Title != want.Title {
		t.Errorf("got %s want %s", got.Title, want.Title)
	}
	if got.Alias != want.Alias {
		t.Errorf("got %s want %s", got.Alias, want.Alias)
	}
	if got.Original != want.Original {
		t.Errorf("got %v want %v", got.Original, want.Original)
	}
	if got.Artist != want.Artist {
		t.Errorf("got %s want %s", got.Artist, want.Artist)
	}
	if got.TracksCount != want.TracksCount {
		t.Errorf("got %d want %d", got.TracksCount, want.TracksCount)
	}
}

func TestGetAndPrintSongText(t *testing.T) {
	t.Parallel()
	query := "david-bowie"
	path := "songs"
	ts := httptest.NewTLSServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != fmt.Sprintf("/%s/%s", path, query) {
				t.Fatalf("wrong url: %s", r.URL.Path)
			}
			http.ServeFile(w, r, "../testdata/song.json")
		}))
	defer ts.Close()
	buf := &bytes.Buffer{}
	c := NewClient("dummy", buf)
	c.Output = buf
	c.BaseURL = ts.URL
	c.HTTPClient = ts.Client()
	want := `Title:       Phish Original:  Original Artist:  TracksCount
David Bowie  true                               447

Tracks
Date:       Venue:                           Location:             Mp3
1986-10-31  Sculpture Room, Goddard College  Plainfield, VT        https://phish.in/audio/000/000/115/115.mp3
1986-12-06  The Ranch                        South Burlington, VT  https://phish.in/audio/000/000/147/147.mp3
1986-12-06  The Ranch                        South Burlington, VT  https://phish.in/audio/000/000/145/145.mp3
`
	ctx := context.Background()
	c.Query = query
	url := c.FormatURL(path)
	err := c.getAndPrintSong(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	if got != want {
		t.Errorf("got \n%s want \n%s", got, want)
	}
}

func TestGetTracks(t *testing.T) {
	t.Parallel()
	ts := httptest.NewTLSServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "../testdata/tracks.json")
		}))
	defer ts.Close()
	c := NewClient("dummy", os.Stdout)
	c.BaseURL = ts.URL
	c.HTTPClient = ts.Client()
	want := TracksOutput{
		Tracks: []Track{
			{
				Title: "Maze",
				ShowDate: "1994-10-07",
				VenueName: "Stabler Arena, Lehigh University",
				VenueLocation: "Bethlehem, PA",
				Mp3: "https://phish.in/audio/000/004/270/4270.mp3",
			},
			{
				Title: "Stash",
				ShowDate: "1993-04-09",
				VenueName: "State Theatre",
				VenueLocation: "Minneapolis, MN",
				Mp3: "https://phish.in/audio/000/006/693/6693.mp3",
			},
		},
	}
	ctx := context.Background()
	url := c.FormatURL("tracks")
	got, err := c.getTracks(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	for i, v := range got.Tracks {
		if v.Title != want.Tracks[i].Title {
			t.Errorf("got %s want %s", v.Title, want.Tracks[i].Title)
		}
		if v.ShowDate != want.Tracks[i].ShowDate {
			t.Errorf("got %s want %s", v.ShowDate, want.Tracks[i].ShowDate)
		}
		if v.VenueLocation != want.Tracks[i].VenueLocation {
			t.Errorf("got %s want %s", v.VenueLocation, want.Tracks[i].VenueLocation)
		}
		if v.VenueName != want.Tracks[i].VenueName {
			t.Errorf("got %s want %s", v.VenueName, want.Tracks[i].VenueName)
		}
		if v.Mp3 != want.Tracks[i].Mp3 {
			t.Errorf("got %s want %s", v.Mp3, want.Tracks[i].Mp3)
		}
	}
}

func TestGetAndPrintTracksText(t *testing.T) {
	t.Parallel()
	ts := httptest.NewTLSServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "../testdata/tracks.json")
		}))
	defer ts.Close()
	buf := &bytes.Buffer{}
	c := NewClient("dummy", buf)
	c.Output = buf
	c.BaseURL = ts.URL
	c.HTTPClient = ts.Client()
	want := `Date:       Venue:                            Location:        Title:  Mp3
1994-10-07  Stabler Arena, Lehigh University  Bethlehem, PA    Maze    https://phish.in/audio/000/004/270/4270.mp3
1993-04-09  State Theatre                     Minneapolis, MN  Stash   https://phish.in/audio/000/006/693/6693.mp3
`
	ctx := context.Background()
	url := c.FormatURL("tracks")
	err := c.getAndPrintTracks(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	if got != want {
		t.Errorf("got \n%s want \n%s", got, want)
	}
}

func TestGetTrack(t *testing.T) {
	t.Parallel()
	query := "stash"
	path := "tracks"
	ts := httptest.NewTLSServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != fmt.Sprintf("/%s/%s", path, query) {
				t.Fatalf("wrong url: %s", r.URL.Path)
			}
			http.ServeFile(w, r, "../testdata/track.json")
		}))
	defer ts.Close()
	c := NewClient("dummy", os.Stdout)
	c.BaseURL = ts.URL
	c.HTTPClient = ts.Client()
	// grab a subset to spot-check values
	want := Track{
		Title: "Stash",
		ShowDate: "1993-04-09",
		VenueName: "State Theatre",
		VenueLocation: "Minneapolis, MN",
		Mp3: "https://phish.in/audio/000/006/693/6693.mp3",
	}
	ctx := context.Background()
	c.Query = query
	url := c.FormatURL(path)
	got, err := c.getTrack(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	if got.Title != want.Title {
		t.Errorf("got %s want %s", got.Title, want.Title)
	}
	if got.ShowDate != want.ShowDate {
		t.Errorf("got %s want %s", got.ShowDate, want.ShowDate)
	}
	if got.VenueName != want.VenueName {
		t.Errorf("got %s want %s", got.VenueName, want.VenueName)
	}
	if got.VenueLocation != want.VenueLocation {
		t.Errorf("got %s want %s", got.VenueLocation, want.VenueLocation)
	}
	if got.Mp3 != want.Mp3 {
		t.Errorf("got %s want %s", got.Mp3, want.Mp3)
	}
}

func TestGetAndPrintTrackText(t *testing.T) {
	t.Parallel()
	query := "stash"
	path := "tracks"
	ts := httptest.NewTLSServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != fmt.Sprintf("/%s/%s", path, query) {
				t.Fatalf("wrong url: %s", r.URL.Path)
			}
			http.ServeFile(w, r, "../testdata/track.json")
		}))
	defer ts.Close()
	buf := &bytes.Buffer{}
	c := NewClient("dummy", buf)
	c.Output = buf
	c.BaseURL = ts.URL
	c.HTTPClient = ts.Client()
	want := `Date:       Venue:         Location:        Title:  Duration  Set    Mp3
1993-04-09  State Theatre  Minneapolis, MN  Stash   11m 15s   Set 1  https://phish.in/audio/000/006/693/6693.mp3

Tags
Name:      Group:              Notes:
SBD        Audio               
Jamcharts  Curated Selections  Several minutes of growly, percussive, dissonant, and atypical jamming.
`
	ctx := context.Background()
	c.Query = query
	url := c.FormatURL(path)
	err := c.getAndPrintTrack(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	if got != want {
		t.Errorf("got \n%s want \n%s", got, want)
	}
}