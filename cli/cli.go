package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

type GenericResponse struct {
	TotalEntries int `json:"total_entries"`
	TotalPages   int `json:"total_pages"`
	Page         int `json:"page"`
	Data         any `json:"data"`
}

func printJSON(w io.Writer, data any) error {
	b, err := json.MarshalIndent(&data, "", "  ")
	if err != nil {
		return fmt.Errorf("unable to convert data to bytes: %w", err)
	}
	fmt.Fprintln(w, string(b))
	return nil
}

type PrettyPrinter interface {
	PrettyPrint(io.Writer, bool) error
}

func PrintResults(w io.Writer, pp PrettyPrinter, json, verbose bool) error {
	if json {
		return printJSON(w, pp)
	}
	return pp.PrettyPrint(w, verbose)
}

type trueAsYes bool

func (s trueAsYes) String() string {
	if s {
		return "yes"
	}
	return "no"
}

func convertMillisecondToConcertDuration(ms int64) string {
	var msInSecond int64 = 1000
	var nsInSecond int64 = 1000000
	t := time.Unix(ms/msInSecond, (ms%msInSecond)*nsInSecond).UTC()
	// early shows, like 1983.12.02, are under an hour
	if t.Hour() != 0 {
		return fmt.Sprintf("%dh %dm", t.Hour(), t.Minute())
	}
	return fmt.Sprintf("%dm %ds", t.Minute(), t.Second())
}

////////////////////////
/* Convenience Types */
//////////////////////

// Show is a convenience struct to hold the show data in the API response.
type Show struct {
	ID         int       `json:"id"`
	Date       string    `json:"date"`
	Duration   int       `json:"duration"`
	Incomplete bool      `json:"incomplete"`
	Sbd        bool      `json:"sbd"`
	Remastered bool      `json:"remastered"`
	Tags       []Tag     `json:"tags"`
	TourID     int       `json:"tour_id"`
	Venue      Venue     `json:"venue"`
	VenueName  string    `json:"venue_name"`
	TaperNotes string    `json:"taper_notes"`
	LikesCount int       `json:"likes_count"`
	VenueID    int       `json:"venue_id"`
	Location   string    `json:"location"`
	Tracks     []Track   `json:"tracks"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Song is a convenience struct to hold the song data in the API response.
type Song struct {
	ID          int       `json:"id"`
	Slug        string    `json:"slug"`
	Title       string    `json:"title"`
	Alias       string    `json:"alias"`
	Original    bool      `json:"original"`
	Artist      string    `json:"artist"`
	Lyrics      string    `json:"lyrics"`
	TracksCount int       `json:"tracks_count"`
	UpdatedAt   time.Time `json:"updated_at"`
	Tracks      []Track   `json:"tracks"`
}

// Tag is a convenience struct to hold the tag data in the API response.
type Tag struct {
	Name  string `json:"name"`
	Group string `json:"group"`
	Notes string `json:"notes"`
}

// TagListItem is a convenience struct to hold the tag data in the API response
// for the /tags endpoint.
type TagListItem struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Group       string    `json:"group"`
	Color       string    `json:"color"`
	Priority    int       `json:"priority"`
	Description string    `json:"description"`
	UpdatedAt   time.Time `json:"updated_at"`
	ShowIds     []int     `json:"show_ids"`
	TrackIds    []int     `json:"track_ids"`
}

// Tour is a convenience struct to hold the tour data in the API response.
type Tour struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	ShowsCount int       `json:"shows_count"`
	Slug       string    `json:"slug"`
	StartsOn   string    `json:"starts_on"`
	EndsOn     string    `json:"ends_on"`
	Shows      []Show    `json:"shows"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Track is a convenience struct to hold the track data in the API response.
type Track struct {
	ID                int       `json:"id"`
	ShowID            int       `json:"show_id"`
	ShowDate          string    `json:"show_date"`
	VenueName         string    `json:"venue_name"`
	VenueLocation     string    `json:"venue_location"`
	Title             string    `json:"title"`
	Position          int       `json:"position"`
	Duration          int       `json:"duration"`
	JamStartsAtSecond int       `json:"jam_starts_at_second"`
	Set               string    `json:"set"`
	SetName           string    `json:"set_name"`
	LikesCount        int       `json:"likes_count"`
	Slug              string    `json:"slug"`
	Tags              []Tag     `json:"tags"`
	Mp3               string    `json:"mp3"`
	WaveformImage     string    `json:"waveform_image"`
	SongIds           []int     `json:"song_ids"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// TrackTag is a convenience struct to hold the tag data found in Tracks data.
type TrackTag struct {
	ID             int       `json:"id"`
	TrackID        int       `json:"track_id"`
	TagID          int       `json:"tag_id"`
	CreatedAt      time.Time `json:"created_at"`
	Notes          string    `json:"notes"`
	StartsAtSecond int       `json:"starts_at_second"`
	EndsAtSecond   int       `json:"ends_at_second"`
	Transcript     string    `json:"transcript"`
}

// Venue is a convenience struct to hold the venue data in the API response.
type Venue struct {
	ID         int       `json:"id"`
	Slug       string    `json:"slug"`
	Name       string    `json:"name"`
	OtherNames []string  `json:"other_names"`
	Latitude   float64   `json:"latitude"`
	Longitude  float64   `json:"longitude"`
	Location   string    `json:"location"`
	City       string    `json:"city"`
	State      string    `json:"state"`
	Country    string    `json:"country"`
	ShowsCount int       `json:"shows_count"`
	ShowDates  []string  `json:"show_dates"`
	ShowIds    []int     `json:"show_ids"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Year is a convenience struct to hold the year data in the API response.
type Year struct {
	Date      string `json:"date"`
	ShowCount int    `json:"show_count"`
}

///////////
/* Eras */
/////////

type ErasResponse struct {
	Data struct {
		One   []string `json:"1.0"`
		Two   []string `json:"2.0"`
		Three []string `json:"3.0"`
		Four  []string `json:"4.0"`
	} `json:"data"`
}

type ErasOutput struct {
	One   []string `json:"1.0"`
	Two   []string `json:"2.0"`
	Three []string `json:"3.0"`
	Four  []string `json:"4.0"`
}

func (e ErasOutput) PrettyPrint(w io.Writer, verbose bool) error {
	_, err := fmt.Fprintf(w,
		"Eras\n1.0: %v\n2.0: %v\n3.0: %v\n4.0: %v\n", strings.Join(e.One, ", "), strings.Join(e.Two, ", "), strings.Join(e.Three, ", "), strings.Join(e.Four, ", "),
	)
	return err
}

type EraResponse struct {
	Era []string `json:"data"`
}

type EraOutput struct {
	EraName string   `json:"era"`
	Years   []string `json:"years"`
}

func (e EraOutput) PrettyPrint(w io.Writer, verbose bool) error {
	_, err := fmt.Fprintf(w, "Era %s:\n%s\n", e.EraName, strings.Join(e.Years, ", "))
	return err
}

////////////
/* Years */
//////////

type YearsResponse struct {
	TotalEntries int    `json:"total_entries"`
	TotalPages   int    `json:"total_pages"`
	Page         int    `json:"page"`
	Data         []Year `json:"data"`
}

type YearsOutput struct {
	Years []Year `json:"years"`
}

func (y YearsOutput) PrettyPrint(w io.Writer, verbose bool) error {
	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', tabwriter.DiscardEmptyColumns)
	fmt.Fprintln(tw, "Years:\tShow Count:")
	for _, year := range y.Years {
		fmt.Fprintf(tw, "%s\t%d\n", year.Date, year.ShowCount)
	}
	return tw.Flush()
}

type YearResponse struct {
	Data []Show `json:"data"`
}

type YearOutput struct {
	Shows ShowsOutput `json:"shows"`
}

// func (y YearOutput) PrettyPrint(w io.Writer, verbose bool) error {

// }

type SongsResponse struct {
	TotalEntries int    `json:"total_entries"`
	TotalPages   int    `json:"total_pages"`
	Page         int    `json:"page"`
	Data         []Song `json:"data"`
}

type SongsOutput struct {
	TotalEntries int          `json:"total_entries"`
	TotalPages   int          `json:"total_pages"`
	CurrentPage  int          `json:"current_page"`
	Songs        []SongOutput `json:"songs"`
}

func (s SongsOutput) PrettyPrint(w io.Writer, verbose bool) error {
	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', tabwriter.DiscardEmptyColumns)
	fmt.Fprintln(tw, "Title:\tOriginal Artist:\tTracksCount:")
	for _, song := range s.Songs {
		artist := "Phish"
		if !song.Original {
			artist = song.Artist
		}
		fmt.Fprintf(tw, "%s\t%s\t%d\n", song.Title, artist, song.TracksCount)
	}
	fmt.Fprintln(tw)
	if s.TotalEntries != 0 {
		fmt.Fprintf(tw, "Total Entries: %d\tTotal Pages: %d\tResult Page: %d\n", s.TotalEntries, s.TotalPages, s.CurrentPage)
	}
	return tw.Flush()
}

type SongResponse struct {
	Data Song `json:"data"`
}

func convertSongToOutput(song Song) SongOutput {
	o := SongOutput{
		ID:          song.ID,
		Title:       song.Title,
		Original:    song.Original,
		Artist:      song.Artist,
		TracksCount: song.TracksCount,
	}
	tracks := convertTracksToOutput(song.Tracks)
	o.Tracks = tracks.Tracks
	return o
}

type SongOutput struct {
	ID          int           `json:"id"`
	Title       string        `json:"title"`
	Original    bool          `json:"original"`
	Artist      string        `json:"artist"`
	TracksCount int           `json:"tracks_count"`
	Tracks      []TrackOutput `json:"tracks"`
}

func (s SongOutput) PrettyPrint(w io.Writer, verbose bool) error {
	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', tabwriter.DiscardEmptyColumns)
	fmt.Fprintln(tw, "Title:\tID:\tOriginal Artist:\tTracksCount:")
	artist := "Phish"
	if !s.Original {
		artist = s.Artist
	}
	fmt.Fprintf(tw, "%s\t%d\t%s\t%d\n", s.Title, s.ID, artist, s.TracksCount)
	fmt.Fprintln(tw)
	fmt.Fprintln(tw, "Tracks")
	fmt.Fprintln(tw, "ID:\tDate:\tVenue:\tLocation:\tDuration:\tMp3")
	for _, t := range s.Tracks {
		fmt.Fprintf(tw, "%d\t%s\t%s\t%s\t%s\t%s\n", t.ID, t.ShowDate, t.VenueName, t.VenueLocation, t.Duration, t.Mp3)
	}
	return tw.Flush()
}

type ToursResponse struct {
	Data []Tour `json:"data"`
}

func convertToursToOutput(tours []Tour) ToursOutput {
	tt := make([]TourOutput, 0, len(tours))
	for _, t := range tours {
		tour := TourOutput{
			Name:       t.Name,
			ShowsCount: t.ShowsCount,
			StartsOn:   t.StartsOn,
			EndsOn:     t.EndsOn,
		}
		shows := convertShowsToOutput(t.Shows)
		tour.Shows = shows.Shows
		tt = append(tt, tour)
	}
	return ToursOutput{Tours: tt}
}

type ToursOutput struct {
	Tours []TourOutput `json:"tours"`
}

func (t ToursOutput) PrettyPrint(w io.Writer, verbose bool) error {
	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', tabwriter.DiscardEmptyColumns)
	fmt.Fprintln(tw, "Name:\tStarts On:\tEnds On:\tShows Count:")
	for _, tour := range t.Tours {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%d\n", tour.Name, tour.StartsOn, tour.EndsOn, tour.ShowsCount)
	}
	return tw.Flush()
}

type TourResponse struct {
	Data Tour `json:"data"`
}

type TourOutput struct {
	Name       string       `json:"name"`
	ShowsCount int          `json:"shows_count"`
	StartsOn   string       `json:"starts_on"`
	EndsOn     string       `json:"ends_on"`
	Shows      []ShowOutput `json:"shows"`
}

func (t TourOutput) PrettyPrint(w io.Writer, verbose bool) error {
	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', tabwriter.DiscardEmptyColumns)
	fmt.Fprintln(tw, "Name:\tStarts On:\tEnds On:\tShow Count:")
	fmt.Fprintf(tw, "%s\t%s\t%s\t%d\n", t.Name, t.StartsOn, t.EndsOn, t.ShowsCount)
	fmt.Fprintln(tw)
	fmt.Fprintln(tw, "ID:\tDate:\tVenue:\tLocation:\tDuration:\tSoundboard:\tRemastered:")
	for _, show := range t.Shows {
		sbd := trueAsYes(show.Sbd)
		r := trueAsYes(show.Remastered)
		fmt.Fprintf(tw, "%d\t%s\t%s\t%s\t%s\t%s\t%s\n", show.ID, show.Date, show.VenueName, show.VenueLocation, show.Duration, sbd, r)
	}
	return tw.Flush()
}

type VenuesResponse struct {
	TotalEntries int     `json:"total_entries"`
	TotalPages   int     `json:"total_pages"`
	Page         int     `json:"page"`
	Data         []Venue `json:"data"`
}

type VenuesOutput struct {
	TotalEntries int           `json:"total_entries"`
	TotalPages   int           `json:"total_pages"`
	CurrentPage  int           `json:"current_page"`
	Venues       []VenueOutput `json:"venues"`
}

func (v VenuesOutput) PrettyPrint(w io.Writer, verbose bool) error {
	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', tabwriter.DiscardEmptyColumns)
	fmt.Fprintln(tw, "Venue:\tLocation:\tShow Count:")
	for _, venue := range v.Venues {
		fmt.Fprintf(tw, "%s\t%s\t%d\n", venue.Name, venue.Location, venue.ShowsCount)
	}
	fmt.Fprintln(tw)
	if v.CurrentPage != 0 {
		fmt.Fprintf(tw, "Total Entries: %d\tTotal Pages: %d\tResult Page: %d\n", v.TotalEntries, v.TotalPages, v.CurrentPage)
	}
	return tw.Flush()
}

type VenueResponse struct {
	Data Venue `json:"data"`
}

func convertVenueToOutput(venue Venue) VenueOutput {
	return VenueOutput{
		Name:       venue.Name,
		Location:   venue.Location,
		ShowsCount: venue.ShowsCount,
		ShowDates:  venue.ShowDates,
	}
}

type VenueOutput struct {
	Name       string   `json:"name"`
	Location   string   `json:"location"`
	ShowsCount int      `json:"shows_count"`
	ShowDates  []string `json:"show_dates"`
}

func (v VenueOutput) PrettyPrint(w io.Writer, verbose bool) error {
	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', tabwriter.DiscardEmptyColumns)
	fmt.Fprintln(tw, "Venue:\tLocation:\tShow Count:")
	fmt.Fprintf(tw, "%s\t%s\t%d\n", v.Name, v.Location, v.ShowsCount)
	fmt.Fprintln(tw)
	if len(v.ShowDates) == 0 {
		return tw.Flush()
	}
	fmt.Fprintln(tw, "Show Dates")
	for _, d := range v.ShowDates {
		fmt.Fprintln(tw, d)
	}
	return tw.Flush()
}

type ShowsResponse struct {
	TotalEntries int    `json:"total_entries"`
	TotalPages   int    `json:"total_pages"`
	Page         int    `json:"page"`
	Data         []Show `json:"data"`
}

func convertShowsToOutput(data []Show) ShowsOutput {
	shows := make([]ShowOutput, 0, len(data))
	for _, s := range data {
		show := convertShowToOutput(s)
		shows = append(shows, show)
	}
	return ShowsOutput{
		Shows: shows,
	}
}

type ShowsOutput struct {
	TotalEntries int          `json:"total_entries"`
	TotalPages   int          `json:"total_pages"`
	CurrentPage  int          `json:"current_page"`
	Shows        []ShowOutput `json:"shows"`
}

func (s ShowsOutput) PrettyPrint(w io.Writer, verbose bool) error {
	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', tabwriter.DiscardEmptyColumns)
	if verbose {
		fmt.Fprintln(tw, "ID:\tDate:\tVenue:\tLocation:\tDuration:\tSoundboard:\tRemastered:")
		for _, show := range s.Shows {
			sbd := trueAsYes(show.Sbd)
			r := trueAsYes(show.Remastered)
			fmt.Fprintf(tw, "%d\t%s\t%s\t%s\t%s\t%s\t%s\n", show.ID, show.Date, show.VenueName, show.VenueLocation, show.Duration, sbd, r)
		}
		// the year details response prints a ShowsOutput but won't have any entries, for example
		if s.TotalEntries != 0 {
			fmt.Fprintln(tw)
			fmt.Fprintf(tw, "Total Entries: %d\tTotal Pages: %d\tResult Page: %d\n", s.TotalEntries, s.TotalPages, s.CurrentPage)
		}
		return tw.Flush()
	}
	fmt.Fprintln(tw, "Date:\tVenue:\tLocation:\tDuration:")
	for _, show := range s.Shows {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", show.Date, show.VenueName, show.VenueLocation, show.Duration)
	}
	// the year details response prints a ShowsOutput but won't have any entries, for example
	if s.TotalEntries != 0 {
		fmt.Fprintln(tw)
		fmt.Fprintf(tw, "Total Entries: %d\tTotal Pages: %d\tResult Page: %d\n", s.TotalEntries, s.TotalPages, s.CurrentPage)
	}
	return tw.Flush()
}

type ShowResponse struct {
	Data Show `json:"data"`
}

func convertShowToOutput(show Show) ShowOutput {
	o := ShowOutput{
		ID:            show.ID,
		Date:          show.Date,
		Duration:      convertMillisecondToConcertDuration(int64(show.Duration)),
		Sbd:           show.Sbd,
		Remastered:    show.Remastered,
		Tags:          show.Tags,
		VenueName:     show.VenueName,
		VenueLocation: show.Location,
	}
	o.Venue = convertVenueToOutput(show.Venue)
	tracks := convertTracksToOutput(show.Tracks)
	o.Tracks = tracks.Tracks
	// some callers have the Location field populated, while
	// others have that information in the embedded Venue struct,
	// so make sure we have that information in one place going forward.
	location := show.Location
	if location == "" {
		location = show.Venue.Location
	}
	o.VenueLocation = location
	return o
}

func convertTagsToString(tags []Tag) string {
	if len(tags) == 0 {
		return ""
	}
	var tt []string
	for _, t := range tags {
		if t.Notes != "" {
			// sometimes inserted mid-text
			notes := strings.ReplaceAll(t.Notes, "\n", "")
			notes = strings.ReplaceAll(notes, "\r", "")
			tt = append(tt, fmt.Sprintf("%s: %s", t.Name, notes))
		} else {
			tt = append(tt, t.Name)
		}
	}
	return strings.Join(tt, ", ")
}

type ShowOutput struct {
	ID            int           `json:"id"`
	Date          string        `json:"date"`
	Duration      string        `json:"duration"`
	Sbd           bool          `json:"sbd"`
	Remastered    bool          `json:"remastered"`
	Tags          []Tag         `json:"tags"`
	Venue         VenueOutput   `json:"venue"`
	VenueName     string        `json:"venue_name"`
	VenueLocation string        `json:"location"`
	Tracks        []TrackOutput `json:"tracks"`
}

func (s ShowOutput) PrettyPrint(w io.Writer, verbose bool) error {
	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', tabwriter.DiscardEmptyColumns)
	if verbose {
		fmt.Fprintln(tw, "ID:\tDate:\tVenue:\tLocation:\tDuration:\tSoundboard:\tRemastered:")
		sbd := trueAsYes(s.Sbd)
		r := trueAsYes(s.Remastered)
		fmt.Fprintf(tw, "%d\t%s\t%s\t%s\t%s\t%s\t%s\n", s.ID, s.Date, s.VenueName, s.VenueLocation, s.Duration, sbd, r)
		fmt.Fprintln(tw)
		if len(s.Tags) != 0 {
			fmt.Fprintln(tw, "Show Tags:")
			tagInfo := convertTagsToString(s.Tags)
			fmt.Fprintln(tw, tagInfo)
			fmt.Fprintln(tw)
		}
		// should always have tracks but worth a check
		if len(s.Tracks) == 0 {
			return tw.Flush()
		}
		longestTitleLen := 0
		for _, t := range s.Tracks {
			if len(t.Title) > longestTitleLen {
				longestTitleLen = len(t.Title)
			}
		}
		fmt.Fprintln(tw, s.Tracks[0].SetName)
		for i, t := range s.Tracks {
			if i > 1 && t.SetName != s.Tracks[i-1].SetName {
				fmt.Fprintln(tw)
				fmt.Fprintln(tw, s.Tracks[i].SetName)
			}
			// we want the title - duration distance the same
			// across sets, so make all titles the same length
			toAdd := longestTitleLen - len(t.Title)
			title := t.Title + strings.Repeat(" ", toAdd)
			fmt.Fprintf(tw, "%s\t%s\n", title, t.Duration)
		}
		fmt.Fprintln(tw)
		fmt.Fprintln(tw, "Track Info:")
		for _, t := range s.Tracks {
			fmt.Fprintln(tw, t.Title)
			fmt.Fprintln(tw, t.Mp3)
			tagInfo := convertTagsToString(t.Tags)
			if tagInfo != "" {
				fmt.Fprintln(tw, tagInfo)
			}
			fmt.Fprintln(tw)
		}
		return tw.Flush()
	}
	fmt.Fprintln(tw, "Date:\tVenue:\tLocation:")
	fmt.Fprintf(tw, "%s\t%s\t%s\n", s.Date, s.VenueName, s.VenueLocation)
	fmt.Fprintln(tw)
	// should always have tracks but worth a check
	if len(s.Tracks) == 0 {
		return tw.Flush()
	}
	longestTitleLen := 0
	for _, t := range s.Tracks {
		if len(t.Title) > longestTitleLen {
			longestTitleLen = len(t.Title)
		}
	}
	fmt.Fprintln(tw, s.Tracks[0].SetName)
	for i, t := range s.Tracks {
		if i > 1 && t.SetName != s.Tracks[i-1].SetName {
			fmt.Fprintln(tw)
			fmt.Fprintln(tw, s.Tracks[i].SetName)
		}
		toAdd := longestTitleLen - len(t.Title)
		title := t.Title + strings.Repeat(" ", toAdd)
		fmt.Fprintf(tw, "%s\t%s\n", title, t.Duration)
	}
	return tw.Flush()
}

type ShowOnDateResponse struct {
	Data Show `json:"data"`
}

type ShowsOnDayOfYear struct {
	Data []Show `json:"data"`
}

type RandomShowResponse struct {
	Data Show `json:"data"`
}

type TracksResponse struct {
	TotalEntries int     `json:"total_entries"`
	TotalPages   int     `json:"total_pages"`
	Page         int     `json:"page"`
	Data         []Track `json:"data"`
}

func convertTracksToOutput(data []Track) TracksOutput {
	tracks := make([]TrackOutput, 0, len(data))
	for _, t := range data {
		tracks = append(tracks, convertTrackToOutput(t))
	}
	return TracksOutput{
		Tracks: tracks,
	}
}

type TracksOutput struct {
	TotalEntries int           `json:"total_entries"`
	TotalPages   int           `json:"total_pages"`
	CurrentPage  int           `json:"current_page"`
	Tracks       []TrackOutput `json:"tracks"`
}

func (t TracksOutput) PrettyPrint(w io.Writer, verbose bool) error {
	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', tabwriter.DiscardEmptyColumns)
	fmt.Fprintln(tw, "ID:\tDate:\tVenue:\tLocation:\tTitle:\tMp3:")
	for _, track := range t.Tracks {
		fmt.Fprintf(tw, "%d\t%s\t%s\t%s\t%s\t%s\n", track.ID, track.ShowDate, track.VenueName, track.VenueLocation, track.Title, track.Mp3)
	}
	fmt.Fprintln(tw)
	if t.TotalEntries != 0 {
		fmt.Fprintf(tw, "Total Entries: %d\tTotal Pages: %d\tResult Page: %d\n", t.TotalEntries, t.TotalPages, t.CurrentPage)
	}
	return tw.Flush()
}

type TrackResponse struct {
	Data Track `json:"data"`
}

func convertTrackToOutput(track Track) TrackOutput {
	return TrackOutput{
		ID:            track.ID,
		ShowDate:      track.ShowDate,
		VenueName:     track.VenueName,
		VenueLocation: track.VenueLocation,
		Title:         track.Title,
		Duration:      convertMillisecondToConcertDuration(int64(track.Duration)),
		SetName:       track.SetName,
		Tags:          track.Tags,
		Mp3:           track.Mp3,
	}
}

type TrackOutput struct {
	ID            int    `json:"id"`
	ShowDate      string `json:"show_date"`
	VenueName     string `json:"venue_name"`
	VenueLocation string `json:"venue_location"`
	Title         string `json:"title"`
	Duration      string `json:"duration"`
	SetName       string `json:"set_name"`
	Tags          []Tag  `json:"tags"`
	Mp3           string `json:"mp3"`
}

func (t TrackOutput) PrettyPrint(w io.Writer, verbose bool) error {
	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', tabwriter.DiscardEmptyColumns)
	fmt.Fprintln(tw, "ID:\tDate:\tVenue:\tLocation:\tTitle:\tDuration\tSet\tMp3")
	fmt.Fprintf(tw, "%d\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n", t.ID, t.ShowDate, t.VenueName, t.VenueLocation, t.Title, t.Duration, t.SetName, t.Mp3)
	fmt.Fprintln(tw)
	if len(t.Tags) != 0 {
		fmt.Fprintln(tw, "Tags")
		fmt.Fprintln(tw, "Name:\tGroup:\tNotes:")
		for _, tag := range t.Tags {
			fmt.Fprintf(tw, "%s\t%s\t%s\n", tag.Name, tag.Group, tag.Notes)
		}
	}
	return tw.Flush()
}

type TagsResponse struct {
	Data []TagListItem `json:"data"`
}

type TagsOutput struct {
	Tags []TagListItemOutput `json:"tags"`
}

func (t TagsOutput) PrettyPrint(w io.Writer, verbose bool) error {
	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', tabwriter.DiscardEmptyColumns)
	fmt.Fprintln(tw, "Name:\tDescription:\tGroup:")
	for _, tag := range t.Tags {
		fmt.Fprintf(tw, "%s\t%s\t%s\n", tag.Name, tag.Description, tag.Group)
	}
	return tw.Flush()
}

type TagResponse struct {
	Data TagListItem `json:"data"`
}

func convertTagListItemToOutput(t TagListItem) TagListItemOutput {
	o := TagListItemOutput{
		Name:        t.Name,
		Group:       t.Group,
		Description: t.Description,
		ShowIds:     t.ShowIds,
		TrackIds:    t.TrackIds,
	}
	return o
}

type TagListItemOutput struct {
	Name        string `json:"name"`
	Group       string `json:"group"`
	Description string `json:"description"`
	ShowIds     []int  `json:"show_ids"`
	TrackIds    []int  `json:"track_ids"`
}

func (t TagListItemOutput) PrettyPrint(w io.Writer, verbose bool) error {
	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', tabwriter.DiscardEmptyColumns)
	fmt.Fprintln(tw, "Name:\tDescription:\tGroup:")
	fmt.Fprintf(tw, "%s\t%s\t%s\n", t.Name, t.Description, t.Group)
	fmt.Fprintln(tw)
	showIDs := make([]string, 0, len(t.ShowIds))
	for _, i := range t.ShowIds {
		showIDs = append(showIDs, strconv.Itoa(i))
	}
	trackIDs := make([]string, 0, len(t.TrackIds))
	for _, i := range t.TrackIds {
		trackIDs = append(trackIDs, strconv.Itoa(i))
	}
	fmt.Fprintf(tw, "Show IDs Where %s Appears\n", t.Name)
	fmt.Fprintf(tw, "%s\n", strings.Join(showIDs, ", "))
	fmt.Fprintln(tw)
	fmt.Fprintf(tw, "Track IDs Where %s Appears\n", t.Name)
	fmt.Fprintf(tw, "%s\n", strings.Join(trackIDs, ", "))

	return tw.Flush()
}

// todo: confirm showtags format
type SearchResponse struct {
	Data struct {
		ExactShow  Show          `json:"exact_show,omitempty"`
		OtherShows []Show        `json:"other_shows,omitempty"`
		ShowTags   []interface{} `json:"show_tags"`
		Songs      []Song        `json:"songs,omitempty"`
		Tags       []TagListItem `json:"tags,omitempty"`
		Tours      []Tour        `json:"tours,omitempty"`
		TrackTags  []TrackTag    `json:"track_tags,omitempty"`
		Tracks     []Track       `json:"tracks,omitempty"`
		Venues     []Venue       `json:"venues,omitempty"`
	} `json:"data"`
}

type TrackTagOutput struct {
	ID         int    `json:"id"`
	TrackID    int    `json:"track_id"`
	TagID      int    `json:"tag_id"`
	Notes      string `json:"notes"`
	Transcript string `json:"transcript"`
}

type TrackTagsOutput struct {
	Tags []TrackTagOutput
}

func (t TrackTagsOutput) PrettyPrint(w io.Writer, verbose bool) error {
	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', tabwriter.DiscardEmptyColumns)
	for _, tag := range t.Tags {
		fmt.Fprintln(tw, "ID:\tTrackID:\tTagID:")
		fmt.Fprintf(tw, "%d\t%d\t%d\n", tag.ID, tag.TrackID, tag.TagID)
		if tag.Notes != "" {
			fmt.Fprintln(tw)
			fmt.Fprintln(tw, "Notes:")
			notes := strings.ReplaceAll(tag.Notes, "&gt;", ">")
			fmt.Fprintln(tw, notes)
		}
		if tag.Transcript != "" {
			if tag.Notes != "" {
				fmt.Fprintln(tw)
			}
			fmt.Fprintln(tw, "Transcript:")
			fmt.Fprintln(tw, tag.Transcript)
		}
		fmt.Fprintln(tw)
	}
	return tw.Flush()
}

func convertSearchToSearchOutput(s SearchResponse) SearchOutput {
	o := SearchOutput{}
	if s.Data.ExactShow.ID != 0 {
		show := convertShowToOutput(s.Data.ExactShow)
		o.Results.ExactShow = &show
	}
	if len(s.Data.OtherShows) != 0 {
		shows := convertShowsToOutput(s.Data.OtherShows)
		o.Results.OtherShows = shows.Shows
	}

	// todo show tags

	if len(s.Data.Songs) != 0 {
		songs := make([]SongOutput, 0, len(s.Data.Songs))
		for _, song := range s.Data.Songs {
			songs = append(songs, convertSongToOutput(song))
		}
		o.Results.Songs = songs
	}
	if len(s.Data.Tags) != 0 {
		tags := make([]TagListItemOutput, 0, len(s.Data.Tags))
		for _, tag := range s.Data.Tags {
			tags = append(tags, convertTagListItemToOutput(tag))
		}
		o.Results.Tags = tags
	}
	if len(s.Data.Tours) != 0 {
		tours := convertToursToOutput(s.Data.Tours)
		o.Results.Tours = tours.Tours
	}
	if len(s.Data.TrackTags) != 0 {
		tags := make([]TrackTagOutput, 0, len(s.Data.TrackTags))
		for _, t := range s.Data.TrackTags {
			tags = append(tags, TrackTagOutput{
				ID:         t.ID,
				TrackID:    t.TrackID,
				TagID:      t.TagID,
				Notes:      t.Notes,
				Transcript: t.Transcript,
			})
		}
		o.Results.TrackTags = tags
	}
	if len(s.Data.Tracks) != 0 {
		tracks := convertTracksToOutput(s.Data.Tracks)
		o.Results.Tracks = tracks.Tracks
	}
	if len(s.Data.Venues) != 0 {
		venues := make([]VenueOutput, 0, len(s.Data.Venues))
		for _, venue := range s.Data.Venues {
			venues = append(venues, convertVenueToOutput(venue))
		}
		o.Results.Venues = venues
	}
	return o
}

type SearchOutput struct {
	Results struct {
		ExactShow  *ShowOutput         `json:"exact_show,omitempty"`
		OtherShows []ShowOutput        `json:"other_shows,omitempty"`
		ShowTags   []any               `json:"show_tags,omitempty"`
		Songs      []SongOutput        `json:"songs,omitempty"`
		Tags       []TagListItemOutput `json:"tags,omitempty"`
		Tours      []TourOutput        `json:"tours,omitempty"`
		TrackTags  []TrackTagOutput    `json:"track_tags,omitempty"`
		Tracks     []TrackOutput       `json:"tracks,omitempty"`
		Venues     []VenueOutput       `json:"venues,omitempty"`
	} `json:"results"`
}

func (s SearchOutput) PrettyPrint(w io.Writer, verbose bool) error {
	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', tabwriter.DiscardEmptyColumns)
	var results bool
	if s.Results.ExactShow != nil {
		results = true
		fmt.Fprintln(tw, "*** EXACT SHOW RESULTS ***")
		if err := s.Results.ExactShow.PrettyPrint(w, true); err != nil {
			return err
		}
		fmt.Fprintln(tw)
	}
	if len(s.Results.OtherShows) != 0 {
		results = true
		fmt.Fprintln(tw, "*** SHOW RESULTS ***")
		so := ShowsOutput{Shows: s.Results.OtherShows}
		if err := so.PrettyPrint(w, true); err != nil {
			return err
		}
		fmt.Fprintln(tw)
	}
	if len(s.Results.ShowTags) != 0 {
		results = true
		fmt.Fprintln(tw, "*** SHOW TAG RESULTS ***")
		// todo
		fmt.Fprintln(tw)
	}
	if len(s.Results.Songs) != 0 {
		results = true
		fmt.Fprintln(tw, "*** SONG RESULTS ***")
		so := SongsOutput{Songs: s.Results.Songs}
		if err := so.PrettyPrint(w, false); err != nil {
			return err
		}
		fmt.Fprintln(tw)
	}
	if len(s.Results.Tags) != 0 {
		results = true
		fmt.Fprintln(tw, "*** TAG RESULTS ***")
		to := TagsOutput{Tags: s.Results.Tags}
		if err := to.PrettyPrint(w, false); err != nil {
			return err
		}
		fmt.Fprintln(tw)
	}
	if len(s.Results.Tours) != 0 {
		results = true
		fmt.Fprintln(tw, "*** TOUR RESULTS ***")
		to := ToursOutput{Tours: s.Results.Tours}
		if err := to.PrettyPrint(w, false); err != nil {
			return err
		}
		fmt.Fprintln(tw)
	}
	if len(s.Results.TrackTags) != 0 {
		results = true
		fmt.Fprintln(tw, "*** TRACK TAG RESULTS ***")
		to := TrackTagsOutput{Tags: s.Results.TrackTags}
		if err := to.PrettyPrint(w, false); err != nil {
			return err
		}
		fmt.Fprintln(tw)
	}
	if len(s.Results.Tracks) != 0 {
		results = true
		fmt.Fprintln(tw, "*** TRACK RESULTS ***")
		to := TracksOutput{Tracks: s.Results.Tracks}
		if err := to.PrettyPrint(w, false); err != nil {
			return err
		}
		fmt.Fprintln(tw)
	}
	if len(s.Results.Venues) != 0 {
		results = true
		fmt.Fprintln(tw, "*** VENUE RESULTS ***")
		vo := VenuesOutput{Venues: s.Results.Venues}
		if err := vo.PrettyPrint(w, false); err != nil {
			return err
		}
		fmt.Fprintln(tw)
	}
	if !results {
		fmt.Fprint(os.Stderr, searchTips)
	}
	return nil
}

type WriteCounter struct {
	ContentLength int64
	TotalWritten  int64
	Name          string
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.TotalWritten += int64(n)
	wc.PrintProgress()
	return n, nil
}

func (wc *WriteCounter) PrintProgress() {
	fmt.Printf("\r%s", strings.Repeat(" ", 70))
	fmt.Printf("\rdownloaded %s of %s", humanizeBytes(wc.TotalWritten), wc.Name)
}

func humanizeBytes(b int64) string {
	base := 1024.0
	sizes := []string{" B", " KiB", " MiB", " GiB"}

	if b < 10 {
		return fmt.Sprintf("%d B", b)
	}
	e := math.Floor(logn(float64(b), base))
	suffix := sizes[int(e)]
	val := math.Floor(float64(b)/math.Pow(base, e)*10+0.5) / 10
	f := "%.0f %s"
	if val < 10 {
		f = "%.1f %s"
	}
	return fmt.Sprintf(f, val, suffix)
}

func logn(n, b float64) float64 {
	return math.Log(n) / math.Log(b)
}
