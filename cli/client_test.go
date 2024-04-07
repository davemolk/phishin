package cli

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

var (
	updateGolden = flag.Bool("update", false, "update golden files")
)

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

func getGoldenValue(t *testing.T, goldenFile, actual string, update bool) string {
	t.Helper()
	path := "../testdata/" + goldenFile
	f, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = f.Close() }()
	if update {
		if _, err := f.WriteString(actual); err != nil {
			t.Fatal(err)
		}
		return actual
	}

	b, err := io.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

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

func TestFromArgs(t *testing.T) {
	t.Parallel()
	c := NewClient("dummy", io.Discard)
	t.Run("error for unrecognized command", func(t *testing.T) {
		if err := c.fromArgs([]string{"phish"}); err == nil {
			t.Error("wanted error, got nil")
		}
	})
	t.Run("show-on-date errors with no query", func(t *testing.T) {
		err := c.fromArgs([]string{"show-on-date", "-s", ""})
		if err == nil {
			t.Error("wanted error, got nil")
		}
	})
	t.Run("show-on-date does not error with query", func(t *testing.T) {
		if err := c.fromArgs([]string{"show-on-date", "-s", "1994-10-31"}); err != nil {
			t.Errorf("wanted nil, got %v", err)
		}
	})
	t.Run("shows-on-day-of-year errors with no query", func(t *testing.T) {
		err := c.fromArgs([]string{"shows-on-day-of-year", "-s", ""})
		if err == nil {
			t.Error("wanted error, got nil")
		}
	})
	t.Run("shows-on-day-of-year does not error with query", func(t *testing.T) {
		if err := c.fromArgs([]string{"shows-on-day-of-year", "-s", "10-31"}); err != nil {
			t.Errorf("wanted nil, got %v", err)
		}
	})
	t.Run("random-show doesn't take a query param", func(t *testing.T) {
		if err := c.fromArgs([]string{"random-show", "-s", "10-31"}); err != nil {
			t.Errorf("wanted nil, got %v", err)
		}
		if c.Query != "" {
			t.Errorf("got %q wanted nothing", c.Query)
		}
	})
	t.Run("search errors with no query", func(t *testing.T) {
		err := c.fromArgs([]string{"search", "-s", ""})
		if err == nil {
			t.Error("wanted error, got nil")
		}
	})
	t.Run("search does not error with query", func(t *testing.T) {
		if err := c.fromArgs([]string{"search", "-s", "costume"}); err != nil {
			t.Errorf("wanted nil, got %v", err)
		}
	})
	t.Run("eras, tours, and tags don't take pagination or sort params", func(t *testing.T) {
		t.Run("eras no pagination", func(t *testing.T) {
			if err := c.fromArgs([]string{"eras", "-pp", "15"}); err != nil {
				t.Errorf("wanted nil, got %v", err)
			}
			if len(c.Parameters) != 0 {
				t.Errorf("got %d wanted 0", len(c.Parameters))
			}
		})
		t.Run("tours no pagination", func(t *testing.T) {
			if err := c.fromArgs([]string{"tours", "-pp", "15"}); err != nil {
				t.Errorf("wanted nil, got %v", err)
			}
			if len(c.Parameters) != 0 {
				t.Errorf("got %d wanted 0", len(c.Parameters))
			}
		})
		t.Run("tags no pagination", func(t *testing.T) {
			if err := c.fromArgs([]string{"tags", "-pp", "15"}); err != nil {
				t.Errorf("wanted nil, got %v", err)
			}
			if len(c.Parameters) != 0 {
				t.Errorf("got %d wanted 0", len(c.Parameters))
			}
		})
		t.Run("eras no sort", func(t *testing.T) {
			if err := c.fromArgs([]string{"eras", "-dir", "asc"}); err != nil {
				t.Errorf("wanted nil, got %v", err)
			}
			if len(c.Parameters) != 0 {
				t.Errorf("got %d wanted 0", len(c.Parameters))
			}
		})
		t.Run("tours no sort", func(t *testing.T) {
			if err := c.fromArgs([]string{"tours", "-dir", "asc"}); err != nil {
				t.Errorf("wanted nil, got %v", err)
			}
			if len(c.Parameters) != 0 {
				t.Errorf("got %d wanted 0", len(c.Parameters))
			}
		})
		t.Run("tags no sort", func(t *testing.T) {
			if err := c.fromArgs([]string{"tags", "-dir", "asc"}); err != nil {
				t.Errorf("wanted nil, got %v", err)
			}
			if len(c.Parameters) != 0 {
				t.Errorf("got %d wanted 0", len(c.Parameters))
			}
		})
	})
	t.Run("include_show_counts=true added to years", func(t *testing.T) {
		if err := c.fromArgs([]string{"years"}); err != nil {
			t.Errorf("wanted nil, got %v", err)
		}
		want := "include_show_counts=true"
		if len(c.Parameters) != 1 {
			t.Errorf("got param length %d, wanted 1", len(c.Parameters))
		}
		if c.Parameters[0] != "include_show_counts=true" {
			t.Errorf("got %q wanted %q", c.Parameters[0], want)
		}
		// reset parameters
		c.Parameters = nil
	})
	t.Run("songs does not support tag flag", func(t *testing.T) {
		if err := c.fromArgs([]string{"songs", "-tag", "sbd"}); err != nil {
			t.Errorf("wanted nil, got %v", err)
		}
		if len(c.Parameters) != 0 {
			t.Errorf("got %d wanted 0", len(c.Parameters))
		}
	})
	t.Run("venues does not support tag flag", func(t *testing.T) {
		if err := c.fromArgs([]string{"venues", "-tag", "sbd"}); err != nil {
			t.Errorf("wanted nil, got %v", err)
		}
		if len(c.Parameters) != 0 {
			t.Errorf("got %d wanted 0", len(c.Parameters))
		}
	})
	t.Run("perPage of 20 will not be added to params list", func(t *testing.T) {
		if err := c.fromArgs([]string{"venues", "-pp", "20"}); err != nil {
			t.Errorf("wanted nil, got %v", err)
		}
		if len(c.Parameters) != 0 {
			t.Errorf("got %d wanted 0", len(c.Parameters))
		}
	})
	t.Run("perPage of < 1 will not be added to params list", func(t *testing.T) {
		if err := c.fromArgs([]string{"venues", "-pp", "0"}); err != nil {
			t.Errorf("wanted nil, got %v", err)
		}
		if len(c.Parameters) != 0 {
			t.Errorf("got %d wanted 0", len(c.Parameters))
		}
	})
	t.Run("perPage of > 1 and !20 will  be added to params list", func(t *testing.T) {
		if err := c.fromArgs([]string{"venues", "-pp", "10"}); err != nil {
			t.Errorf("wanted nil, got %v", err)
		}
		if len(c.Parameters) != 1 {
			t.Errorf("got %d wanted 1", len(c.Parameters))
		}
		want := "per_page=10"
		if c.Parameters[0] != want {
			t.Errorf("got %q want %q", c.Parameters[0], want)
		}
		// reset
		c.Parameters = nil
	})
	t.Run("page < 2 will not be set", func(t *testing.T) {
		if err := c.fromArgs([]string{"venues", "-p", "0"}); err != nil {
			t.Errorf("wanted nil, got %v", err)
		}
		if len(c.Parameters) != 0 {
			t.Errorf("got %d wanted 0", len(c.Parameters))
		}
	})
	t.Run("page > 1 are set", func(t *testing.T) {
		if err := c.fromArgs([]string{"venues", "-p", "10"}); err != nil {
			t.Errorf("wanted nil, got %v", err)
		}
		if len(c.Parameters) != 1 {
			t.Errorf("got %d wanted 1", len(c.Parameters))
		}
		want := "page=10"
		if c.Parameters[0] != want {
			t.Errorf("got %q want %q", c.Parameters[0], want)
		}
		// reset
		c.Parameters = nil
	})
	t.Run("sort directions other than asc and desc are ignored", func(t *testing.T) {
		if err := c.fromArgs([]string{"venues", "-dir", "phish"}); err != nil {
			t.Errorf("wanted nil, got %v", err)
		}
		if len(c.Parameters) != 0 {
			t.Errorf("got %d wanted 0", len(c.Parameters))
		}
		t.Run("accepts asc", func(t *testing.T) {
			if err := c.fromArgs([]string{"venues", "-dir", "asc"}); err != nil {
				t.Errorf("wanted nil, got %v", err)
			}
			if len(c.Parameters) != 1 {
				t.Errorf("got %d wanted 1", len(c.Parameters))
			}
			want := "sort_dir=asc"
			if c.Parameters[0] != want {
				t.Errorf("got %q want %q", c.Parameters[0], want)
			}
			// reset
			c.Parameters = nil
		})
		t.Run("accepts desc", func(t *testing.T) {
			if err := c.fromArgs([]string{"venues", "-sort-dir", "desc"}); err != nil {
				t.Errorf("wanted nil, got %v", err)
			}
			if len(c.Parameters) != 1 {
				t.Errorf("got %d wanted 1", len(c.Parameters))
			}
			want := "sort_dir=desc"
			if c.Parameters[0] != want {
				t.Errorf("got %q want %q", c.Parameters[0], want)
			}
			// reset
			c.Parameters = nil
		})
	})
	t.Run("sort attr not added to params if blank", func(t *testing.T) {
		if err := c.fromArgs([]string{"venues", "-sort-attr", ""}); err != nil {
			t.Errorf("wanted nil, got %v", err)
		}
		if len(c.Parameters) != 0 {
			t.Errorf("got %d wanted 0", len(c.Parameters))
		}
		t.Run("attr otherwise not validated, just added", func(t *testing.T) {
			if err := c.fromArgs([]string{"venues", "-a", "phish"}); err != nil {
				t.Errorf("wanted nil, got %v", err)
			}
			want := "sort_attr=phish"
			if c.Parameters[0] != want {
				t.Errorf("got %q want %q", c.Parameters[0], want)
			}
			// reset
			c.Parameters = nil
		})
	})
	t.Run("empty tag won't be added to params", func(t *testing.T) {
		if err := c.fromArgs([]string{"shows", "-tag", ""}); err != nil {
			t.Errorf("wanted nil, got %v", err)
		}
		if len(c.Parameters) != 0 {
			t.Errorf("got %d wanted 0", len(c.Parameters))
		}
		t.Run("non-empty tag will be added", func(t *testing.T) {
			if err := c.fromArgs([]string{"shows", "-tag", "sbd"}); err != nil {
				t.Errorf("wanted nil, got %v", err)
			}
			want := "tag=sbd"
			if c.Parameters[0] != want {
				t.Errorf("got %q want %q", c.Parameters[0], want)
			}
			// reset
			c.Parameters = nil
		})
	})
}

func TestClientRun(t *testing.T) {
	c := NewClient("dummy", nil)
	tt := []struct {
		name      string
		serveFile string
		path      string
		golden    string
		json      bool
		verbose   bool
		query     string
		raw       bool
	}{
		{
			name:      "eras",
			serveFile: "../testdata/eras.json",
			path:      "eras",
			golden:    "eras.golden",
			json:      false,
			verbose:   false,
			query:     "",
			raw:       false,
		},
		{
			name:      "eras json",
			serveFile: "../testdata/simple_eras.json",
			path:      "eras",
			golden:    "eras.json.golden",
			json:      true,
			verbose:   false,
			query:     "",
			raw:       false,
		},
		{
			name:      "era",
			serveFile: "../testdata/era.json",
			path:      "eras",
			golden:    "era.golden",
			json:      false,
			verbose:   false,
			query:     "3.0",
			raw:       false,
		},
		{
			name:      "era",
			serveFile: "../testdata/era.json",
			path:      "eras",
			golden:    "era.json.golden",
			json:      true,
			verbose:   false,
			query:     "3.0",
			raw:       false,
		},
		{
			name:      "years",
			serveFile: "../testdata/years.json",
			path:      "years",
			golden:    "years.golden",
			json:      false,
			verbose:   false,
			query:     "",
			raw:       false,
		},
		{
			name:      "year non-verbose",
			serveFile: "../testdata/year.json",
			path:      "years",
			golden:    "year.nonverbose.golden",
			json:      false,
			verbose:   false,
			query:     "1994",
			raw:       false,
		},
		{
			name:      "year verbose",
			serveFile: "../testdata/year.json",
			path:      "years",
			golden:    "year.verbose.golden",
			json:      false,
			verbose:   true,
			query:     "1994",
			raw:       false,
		},
		{
			name:      "shows non-verbose",
			serveFile: "../testdata/shows.json",
			path:      "shows",
			golden:    "shows.nonverbose.golden",
			json:      false,
			verbose:   false,
			query:     "",
			raw:       false,
		},
		{
			name:      "shows verbose",
			serveFile: "../testdata/shows.json",
			path:      "shows",
			golden:    "shows.verbose.golden",
			json:      false,
			verbose:   true,
			query:     "",
			raw:       false,
		},
		{
			name:      "show non-verbose",
			serveFile: "../testdata/show.json",
			path:      "shows",
			golden:    "show.nonverbose.golden",
			json:      false,
			verbose:   false,
			query:     "1990-04-05",
			raw:       false,
		},
		{
			name:      "show verbose",
			serveFile: "../testdata/show.json",
			path:      "shows",
			golden:    "show.verbose.golden",
			json:      false,
			verbose:   true,
			query:     "1990-04-05",
			raw:       false,
		},
		{
			name:      "venues",
			serveFile: "../testdata/venues.json",
			path:      "venues",
			golden:    "venues.golden",
			json:      false,
			verbose:   false,
			query:     "",
			raw:       false,
		},
		{
			name:      "venue",
			serveFile: "../testdata/venue.json",
			path:      "venues",
			golden:    "venue.golden",
			json:      false,
			verbose:   false,
			query:     "the-academy",
			raw:       false,
		},
		{
			name:      "venue raw",
			serveFile: "../testdata/nectars.json",
			path:      "venues",
			golden:    "nectars.golden",
			json:      false,
			verbose:   false,
			query:     "nectar-s",
			raw:       true,
		},
		{
			name:      "tags",
			serveFile: "../testdata/tags.json",
			path:      "tags",
			golden:    "tags.golden",
			json:      false,
			verbose:   false,
			query:     "",
			raw:       false,
		},
		{
			name:      "tag",
			serveFile: "../testdata/tag.json",
			path:      "tags",
			golden:    "tag.golden",
			json:      false,
			verbose:   false,
			query:     "jamcharts",
			raw:       false,
		},
		{
			name:      "tours",
			serveFile: "../testdata/tours.json",
			path:      "tours",
			golden:    "tours.golden",
			json:      false,
			verbose:   false,
			query:     "",
			raw:       false,
		},
		{
			name:      "tour",
			serveFile: "../testdata/tour.json",
			path:      "tours",
			golden:    "tour.golden",
			json:      false,
			verbose:   false,
			query:     "1985-tour",
			raw:       false,
		},
		{
			name:      "songs",
			serveFile: "../testdata/songs.json",
			path:      "songs",
			golden:    "songs.golden",
			json:      false,
			verbose:   false,
			query:     "",
			raw:       false,
		},
		{
			name:      "song",
			serveFile: "../testdata/song.json",
			path:      "songs",
			golden:    "song.golden",
			json:      false,
			verbose:   false,
			query:     "david-bowie",
			raw:       false,
		},
		{
			name:      "tracks",
			serveFile: "../testdata/tracks.json",
			path:      "tracks",
			golden:    "tracks.golden",
			json:      false,
			verbose:   false,
			query:     "",
			raw:       false,
		},
		{
			name:      "track",
			serveFile: "../testdata/track.json",
			path:      "tracks",
			golden:    "track.golden",
			json:      false,
			verbose:   false,
			query:     "stash",
			raw:       false,
		},
		{
			name:      "search",
			serveFile: "../testdata/boulder_search.json",
			path:      "search",
			golden:    "search.golden",
			json:      false,
			verbose:   false,
			query:     "boulder",
			raw:       false,
		},
	}
	for _, tc := range tt {
		ctx := context.Background()
		buf := &bytes.Buffer{}
		c.Output = buf
		c.PrintJSON = tc.json
		c.Verbose = tc.verbose
		c.Query = tc.query
		c.RawOutput = tc.raw
		ts := httptest.NewTLSServer(http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				http.ServeFile(w, r, tc.serveFile)
			}))
		defer ts.Close()
		c.BaseURL = ts.URL
		c.HTTPClient = ts.Client()
		err := c.run(ctx, tc.path)
		if err != nil {
			t.Fatal(err)
		}
		got := buf.String()
		want := getGoldenValue(t, tc.golden, got, *updateGolden)
		if got != want {
			t.Errorf("got\n%s want\n%s", got, want)
		}
	}
}

func TestGetEras(t *testing.T) {
	t.Parallel()
	ts := httptest.NewTLSServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "../testdata/eras.json")
		}))
	defer ts.Close()
	buf := &bytes.Buffer{}
	c := NewClient("dummy", buf)
	c.BaseURL = ts.URL
	c.HTTPClient = ts.Client()
	want := ErasOutput{
		One:   []string{"1983-1987", "1988", "1989", "1990", "1991", "1992", "1993", "1994", "1995", "1996", "1997", "1998", "1999", "2000"},
		Two:   []string{"2002", "2003", "2004"},
		Three: []string{"2009", "2010", "2011", "2012", "2013", "2014", "2015", "2016", "2017", "2018", "2019", "2020"},
		Four:  []string{"2021", "2022", "2023"},
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
		Years:   []string{"2009", "2010", "2011", "2012", "2013", "2014", "2015", "2016", "2017", "2018", "2019", "2020"},
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
				Date:      "1983-1987",
				ShowCount: 34,
			},
			{
				Date:      "1988",
				ShowCount: 44,
			},
			{
				Date:      "1989",
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
		Shows: []ShowOutput{
			{
				ID:         135,
				Date:       "1994-04-04",
				Duration:   "2h 40m",
				Sbd:        true,
				Remastered: false,
				Tags: []Tag{
					{
						Name: "SBD", Group: "Audio",
					},
				},
				VenueName: "The Flynn Theatre",
				Venue: VenueOutput{
					Name:       "The Flynn Theatre",
					Location:   "Burlington, VT",
					ShowsCount: 4,
				},
				VenueLocation: "Burlington, VT",
				Tracks: []TrackOutput{
					{
						ID:            2553,
						ShowDate:      "1994-04-04",
						VenueName:     "The Flynn Theatre",
						VenueLocation: "Burlington, VT",
						Title:         "Divided Sky",
						Duration:      "13m 31s",
						SetName:       "Set 1",
						Tags:          []Tag{},
						Mp3:           "https://phish.in/audio/000/002/553/2553.mp3",
					},
					{
						ID:            2554,
						ShowDate:      "1994-04-04",
						VenueName:     "The Flynn Theatre",
						VenueLocation: "Burlington, VT",
						Title:         "Sample in a Jar",
						Duration:      "4m 59s",
						SetName:       "Set 1",
						Tags:          []Tag{},
						Mp3:           "https://phish.in/audio/000/002/554/2554.mp3",
					},
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
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got \n%v want \n%v", got, want)
	}
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
		TotalEntries: 1759,
		TotalPages:   88,
		CurrentPage:  1,
		Shows: []ShowOutput{
			{
				ID:         696,
				Date:       "1990-04-05",
				Duration:   "2h 27m",
				Sbd:        true,
				Remastered: false,
				Tags: []Tag{
					{
						Name:  "SBD",
						Group: "Audio",
					},
				},
				VenueName:     "J.J. McCabe's",
				VenueLocation: "Boulder, CO",
				Venue: VenueOutput{
					Name:       "J.J. McCabe's",
					Location:   "Boulder, CO",
					ShowsCount: 1,
				},
				Tracks: []TrackOutput{
					{
						ID:            14073,
						ShowDate:      "1990-04-05",
						VenueName:     "J.J. McCabe's",
						VenueLocation: "Boulder, CO",
						Title:         "Possum",
						Duration:      "6m 48s",
						SetName:       "Set 1",
						Tags: []Tag{
							{
								Name:  "SBD",
								Group: "Audio",
							},
						},
						Mp3: "https://phish.in/audio/000/014/073/14073.mp3",
					},
					{
						ID:            14074,
						ShowDate:      "1990-04-05",
						VenueName:     "J.J. McCabe's",
						VenueLocation: "Boulder, CO",
						Title:         "Ya Mar",
						Duration:      "7m 7s",
						SetName:       "Set 1",
						Tags: []Tag{
							{
								Name:  "SBD",
								Group: "Audio",
							},
							{
								Name:  "Tease",
								Group: "Song Content",
								Notes: "Theme from Bonanza by Ray Evans and\n Jay Livingston",
							},
						},
						Mp3: "https://phish.in/audio/000/014/074/14074.mp3",
					},
				},
			},
		},
	}
	ctx := context.Background()
	url := c.FormatURL("shows")
	got, err := c.getShows(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestGetShow(t *testing.T) {
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
	c := NewClient("dummy", os.Stdout)
	c.BaseURL = ts.URL
	c.HTTPClient = ts.Client()
	// grab a subset to spot-check values
	want := ShowOutput{
		ID:         696,
		Date:       "1990-04-05",
		Duration:   "2h 27m",
		Sbd:        true,
		Remastered: false,
		Tags: []Tag{
			{
				Name:  "SBD",
				Group: "Audio",
			},
		},
		VenueName: "J.J. McCabe's",
		Venue: VenueOutput{
			Name:       "J.J. McCabe's",
			Location:   "Boulder, CO",
			ShowsCount: 1,
		},
		Tracks: []TrackOutput{
			{
				ID:            14073,
				ShowDate:      "1990-04-05",
				VenueName:     "J.J. McCabe's",
				VenueLocation: "Boulder, CO",
				Title:         "Possum",
				Duration:      "6m 48s",
				SetName:       "Set 1",
				Tags: []Tag{
					{
						Name:  "SBD",
						Group: "Audio",
					},
				},
				Mp3: "https://phish.in/audio/000/014/073/14073.mp3",
			},
			{
				ID:            14074,
				ShowDate:      "1990-04-05",
				VenueName:     "J.J. McCabe's",
				VenueLocation: "Boulder, CO",
				Title:         "Ya Mar",
				Duration:      "7m 7s",
				SetName:       "Set 1",
				Tags: []Tag{
					{
						Name:  "SBD",
						Group: "Audio",
					},
					{
						Name:  "Tease",
						Group: "Song Content",
						Notes: "Theme from Bonanza by Ray Evans and\n Jay Livingston",
					},
				},
				Mp3: "https://phish.in/audio/000/014/074/14074.mp3",
			},
		},
	}
	ctx := context.Background()
	c.Query = query
	url := c.FormatURL(path)
	got, err := c.getShow(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	// spot-check a few values instead of writing up the whole
	// struct -- we'll confirm we have what we want in the
	// print show text test.
	if got.Date != want.Date {
		t.Errorf("got %v want %v", got.Date, want.Date)
	}
	if got.Sbd != want.Sbd {
		t.Errorf("got %v want %v", got.Sbd, want.Sbd)
	}
	if got.Venue.Location != want.Venue.Location {
		t.Errorf("got %q want %q", got.Venue.Location, want.Venue.Location)
	}
	for i := 0; i < 2; i++ {
		if !reflect.DeepEqual(got.Tracks[i], want.Tracks[i]) {
			t.Errorf("got \n%v want \n%v", got.Tracks[i], want.Tracks[i])
		}
	}
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
		TotalEntries: 666,
		TotalPages:   34,
		CurrentPage:  1,
		Venues: []VenueOutput{
			{
				Name:       "The Base Lodge, Johnson State College",
				Location:   "Johnson, VT",
				ShowsCount: 2,
				ShowDates:  []string{"1988-03-11", "1989-04-14"},
			},
			{
				Name:       "The Academy",
				Location:   "New York, NY",
				ShowsCount: 1,
				ShowDates:  []string{"1991-07-15"},
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
		Name:       "The Academy",
		Location:   "New York, NY",
		ShowsCount: 1,
		ShowDates:  []string{"1991-07-15"},
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
				Name:        "Costume",
				Description: "Musical costume sequence",
				Group:       "Set Content",
			},
			{
				Name:        "Audience",
				Description: "Contribution from audience during performance",
				Group:       "Song Content",
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
		Name:        "Jamcharts",
		Description: "Phish.net Jam Charts selections (phish.net/jamcharts)",
		Group:       "Curated Selections",
		ShowIds:     []int{3},
		TrackIds:    []int{1, 2},
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
		Tours: []TourOutput{
			{
				Name:       "1983 Tour",
				StartsOn:   "1983-12-02",
				EndsOn:     "1983-12-02",
				ShowsCount: 1,
				Shows: []ShowOutput{
					{
						ID:            1324,
						Date:          "1983-12-02",
						Duration:      "17m 11s",
						Sbd:           true,
						Remastered:    false,
						Venue:         VenueOutput{},
						Tags:          nil,
						VenueName:     "Harris-Millis Cafeteria, University of Vermont",
						Tracks:        []TrackOutput{},
						VenueLocation: "Burlington, VT",
					},
				},
			},
			{
				Name:       "1984 Tour",
				StartsOn:   "1984-11-03",
				EndsOn:     "1984-12-01",
				ShowsCount: 2,
				Shows: []ShowOutput{
					{
						ID:            1334,
						Date:          "1984-11-03",
						Duration:      "1h 10m",
						Sbd:           false,
						Remastered:    false,
						Venue:         VenueOutput{},
						Tags:          nil,
						VenueName:     "Slade Hall, University of Vermont",
						VenueLocation: "Burlington, VT",
						Tracks:        []TrackOutput{},
					},
					{
						ID:            2,
						Date:          "1984-12-01",
						Duration:      "1h 35m",
						Sbd:           true,
						Remastered:    false,
						Venue:         VenueOutput{},
						Tags:          nil,
						VenueName:     "Nectar's",
						VenueLocation: "Burlington, VT",
						Tracks:        []TrackOutput{},
					},
				},
			},
		},
	}
	ctx := context.Background()
	url := c.FormatURL("tours")
	got, err := c.getTours(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got.Tours, want.Tours) {
		t.Errorf("got \n%v \nwant \n%v", got.Tours, want.Tours)
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
		Name:       "1985 Tour",
		StartsOn:   "1985-03-04",
		EndsOn:     "1985-11-23",
		ShowsCount: 6,
		Shows: []ShowOutput{
			{
				ID:            3,
				Date:          "1985-03-04",
				Duration:      "40m 14s",
				Sbd:           true,
				Remastered:    false,
				Venue:         VenueOutput{},
				Tags:          nil,
				VenueName:     "Hunt's",
				VenueLocation: "Burlington, VT",
				Tracks:        []TrackOutput{},
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
		TotalEntries: 942,
		TotalPages:   48,
		CurrentPage:  1,
		Songs: []SongOutput{
			{
				ID:          84,
				Title:       "Billy Breathes",
				Original:    true,
				Artist:      "",
				TracksCount: 64,
				Tracks:      []TrackOutput{},
			},
			{
				Title:    "Arc",
				Original: false,
				Artist:   "Arctic Monkeys",
				Tracks:   []TrackOutput{},
			},
		},
	}
	ctx := context.Background()
	url := c.FormatURL("songs")
	got, err := c.getSongs(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got \n%v \nwant\n%v", got, want)
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
	want := SongOutput{
		ID:          979,
		Title:       "David Bowie",
		Original:    true,
		Artist:      "",
		TracksCount: 447,
		Tracks: []TrackOutput{
			{
				ID:            115,
				ShowDate:      "1986-10-31",
				VenueName:     "Sculpture Room, Goddard College",
				VenueLocation: "Plainfield, VT",
				Title:         "David Bowie",
				Duration:      "10m 19s",
				SetName:       "Set 2",
				Tags: []Tag{
					{
						Name:  "SBD",
						Group: "Audio",
					},
					{
						Name:  "Jamcharts",
						Group: "Curated Selections",
						Notes: "Earliest known live version. Jam is played at a slowed tempo initially, but picks up speed and intensity as it develops.",
					},
				},
				Mp3: "https://phish.in/audio/000/000/115/115.mp3",
			},
		},
	}
	ctx := context.Background()
	c.Query = query
	url := c.FormatURL(path)
	got, err := c.getSong(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got \n%v \nwant\n%v", got, want)
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
		TotalEntries: 35069,
		TotalPages:   1754,
		CurrentPage:  1,
		Tracks: []TrackOutput{
			{
				ID:            4270,
				Title:         "Maze",
				ShowDate:      "1994-10-07",
				VenueName:     "Stabler Arena, Lehigh University",
				VenueLocation: "Bethlehem, PA",
				Duration:      "11m 13s",
				SetName:       "Set 2",
				Tags:          []Tag{},
				Mp3:           "https://phish.in/audio/000/004/270/4270.mp3",
			},
			{
				ID:            6693,
				Title:         "Stash",
				ShowDate:      "1993-04-09",
				VenueName:     "State Theatre",
				VenueLocation: "Minneapolis, MN",
				Duration:      "11m 15s",
				SetName:       "Set 1",
				Tags: []Tag{
					{
						Name:  "SBD",
						Group: "Audio",
					},
					{
						Name:  "Jamcharts",
						Group: "Curated Selections",
						Notes: "Several minutes of growly, percussive, dissonant, and atypical jamming.",
					},
				},
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
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got \n%v \nwant\n%v", got, want)
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
	want := TrackOutput{
		ID:            6693,
		Title:         "Stash",
		ShowDate:      "1993-04-09",
		VenueName:     "State Theatre",
		VenueLocation: "Minneapolis, MN",
		Duration:      "11m 15s",
		SetName:       "Set 1",
		Tags: []Tag{
			{
				Name:  "SBD",
				Group: "Audio",
			},
			{
				Name:  "Jamcharts",
				Group: "Curated Selections",
				Notes: "Several minutes of growly, percussive, dissonant, and atypical jamming.",
			},
		},
		Mp3: "https://phish.in/audio/000/006/693/6693.mp3",
	}
	ctx := context.Background()
	c.Query = query
	url := c.FormatURL(path)
	got, err := c.getTrack(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got \n%v \nwant\n%v", got, want)
	}
}

func TestGetSearch(t *testing.T) {
	t.Parallel()
	query := "boulder"
	path := "search"
	ts := httptest.NewTLSServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != fmt.Sprintf("/%s/%s", path, query) {
				t.Fatalf("wrong url: %s", r.URL.Path)
			}
			http.ServeFile(w, r, "../testdata/boulder_search.json")
		}))
	defer ts.Close()
	c := NewClient("dummy", os.Stdout)
	c.BaseURL = ts.URL
	c.HTTPClient = ts.Client()
	want := SearchOutput{
		Results: struct {
			ExactShow  *ShowOutput         "json:\"exact_show,omitempty\""
			OtherShows []ShowOutput        "json:\"other_shows,omitempty\""
			ShowTags   []any               "json:\"show_tags,omitempty\""
			Songs      []SongOutput        "json:\"songs,omitempty\""
			Tags       []TagListItemOutput "json:\"tags,omitempty\""
			Tours      []TourOutput        "json:\"tours,omitempty\""
			TrackTags  []TrackTagOutput    "json:\"track_tags,omitempty\""
			Tracks     []TrackOutput       "json:\"tracks,omitempty\""
			Venues     []VenueOutput       "json:\"venues,omitempty\""
		}{
			TrackTags: []TrackTagOutput{
				{
					ID:         54793,
					TrackID:    10882,
					TagID:      16,
					Notes:      "Colonel Forbin braves thousands of falling rocks and boulders, their collective force transforming the mountainside into the face of the Great and Knowledgeable Icculus, who, in an act of benevolence, calls upon the Famous Mockingbird to retrieve the Helping Friendly Book and save the people of Gamehendge.",
					Transcript: "TREY: Okay, outside right now it's snowing, there's clouds in the sky. Come with us now, lifting up slowly, off the ground. Picture outside this building, thousands, you know, hundreds of miles of clouds over us right now, snow coming down over us cold. Slowly we're lifting up... You can see up above; looking down, looking down you see the building and the streets going by and bodies of water.",
				},
			},
			Venues: []VenueOutput{
				{
					Name:       "Balch Fieldhouse, University of Colorado",
					Location:   "Boulder, CO",
					ShowsCount: 2,
				},
			},
		},
	}
	ctx := context.Background()
	c.Query = query
	url := c.FormatURL(path)
	got, err := c.getSearch(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got \n%v \nwant\n%v", got, want)
	}
}
