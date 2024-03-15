package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

func printJSON(w io.Writer, data any) error {
	b, err := json.MarshalIndent(&data, "", "  ")
	if err != nil {
		return fmt.Errorf("unable to convert data to bytes: %w", err)
	}
	fmt.Fprintln(w, string(b))
	return nil
}

// trueAsYes is a type so we can print "yes" if a bool is true
type trueAsYes bool

func (s trueAsYes) String() string {
	if s {
		return "yes"
	}
	return ""
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

func prettyPrintEras(w io.Writer, e ErasOutput) error {
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

func prettyPrintEra(w io.Writer, e EraOutput) error {
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

func prettyPrintYears(tw *tabwriter.Writer, years YearsOutput) error {
	fmt.Fprintln(tw, "Years:\tShow Count:")
	for _, y := range years.Years {
		fmt.Fprintf(tw, "%s\t%d\n", y.Date, y.ShowCount)
	}
	return tw.Flush()
}

type YearResponse struct {
	Data []Show `json:"data"`
}

type YearOutput struct {
	Shows ShowsOutput `json:"shows"`
}

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

func prettyPrintSongs(tw *tabwriter.Writer, songs SongsOutput) error {
	fmt.Fprintln(tw, "Title:\tPhish Original:\tOriginal Artist:\tTracksCount")
	for _, s := range songs.Songs {
		orig := trueAsYes(s.Original)
		fmt.Fprintf(tw, "%s\t%s\t%s\t%d\n", s.Title, orig, s.Artist, s.TracksCount)
	}
	fmt.Fprintln(tw)
	if songs.TotalEntries != 0 {
		fmt.Fprintf(tw, "Total Entries: %d\tTotal Pages: %d\tResult Page: %d\n", songs.TotalEntries, songs.TotalPages, songs.CurrentPage)
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

func prettyPrintSong(tw *tabwriter.Writer, song SongOutput) error {
	fmt.Fprintln(tw, "Title:\tID:\tPhish Original:\tOriginal Artist:\tTracksCount")
	fmt.Fprintf(tw, "%s\t%d\t%v\t%s\t%d\n", song.Title, song.ID, song.Original, song.Artist, song.TracksCount)
	fmt.Fprintln(tw)
	fmt.Fprintln(tw, "Tracks")
	fmt.Fprintln(tw, "ID:\tDate:\tVenue:\tLocation:\tDuration:\tMp3")
	for _, t := range song.Tracks {
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

func prettyPrintTours(tw *tabwriter.Writer, tours ToursOutput) error {
	fmt.Fprintln(tw, "Name:\tStarts On:\tEnds On:\tShows Count:")
	for _, t := range tours.Tours {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%d\n", t.Name, t.StartsOn, t.EndsOn, t.ShowsCount)
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

func prettyPrintTour(tw *tabwriter.Writer, tour TourOutput) error {
	fmt.Fprintln(tw, "Name:\tStarts On:\tEnds On:\tShow Count:")
	fmt.Fprintf(tw, "%s\t%s\t%s\t%d\n", tour.Name, tour.StartsOn, tour.EndsOn, tour.ShowsCount)
	fmt.Fprintln(tw)
	fmt.Fprintln(tw, "ID:\tDate:\tVenue:\tLocation:\tDuration:\tSoundboard:\tRemastered:")
	for _, show := range tour.Shows {
		sbd := trueAsYes(show.Sbd)
		r := trueAsYes(show.Remastered)
		fmt.Fprintf(tw, "%d\t%s\t%s\t%s\t%s\t%s\t%s\n", show.ID, show.Date, show.VenueName, show.Venue.Location, show.Duration, sbd, r)
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

func prettyPrintVenues(tw *tabwriter.Writer, venues VenuesOutput) error {
	fmt.Fprintln(tw, "Venue:\tLocation:\tShow Count:")
	for _, v := range venues.Venues {
		fmt.Fprintf(tw, "%s\t%s\t%d\n", v.Name, v.Location, v.ShowsCount)
	}
	fmt.Fprintln(tw)
	if venues.CurrentPage != 0 {
		fmt.Fprintf(tw, "Total Entries: %d\tTotal Pages: %d\tResult Page: %d\n", venues.TotalEntries, venues.TotalPages, venues.CurrentPage)
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

func prettyPrintVenue(tw *tabwriter.Writer, venue VenueOutput) error {
	fmt.Fprintln(tw, "Venue:\tLocation:\tShow Count:")
	fmt.Fprintf(tw, "%s\t%s\t%d\n", venue.Name, venue.Location, venue.ShowsCount)
	fmt.Fprintln(tw)
	if len(venue.ShowDates) == 0 {
		return tw.Flush()
	}
	fmt.Fprintln(tw, "Show Dates")
	for _, d := range venue.ShowDates {
		fmt.Fprintf(tw, "%s\n", d)
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

func prettyPrintShows(tw *tabwriter.Writer, shows ShowsOutput, verbose bool) error {
	if verbose {
		fmt.Fprintln(tw, "ID:\tDate:\tVenue:\tLocation:\tDuration:\tSoundboard:\tRemastered:")
		for _, s := range shows.Shows {
			sbd := trueAsYes(s.Sbd)
			r := trueAsYes(s.Remastered)
			fmt.Fprintf(tw, "%d\t%s\t%s\t%s\t%s\t%s\t%s\n", s.ID, s.Date, s.VenueName, s.Venue.Location, s.Duration, sbd, r)
		}
		// the year details response prints a ShowsOutput but won't have any entries, for example
		if shows.TotalEntries != 0 {
			fmt.Fprintln(tw)
			fmt.Fprintf(tw, "Total Entries: %d\tTotal Pages: %d\tResult Page: %d\n", shows.TotalEntries, shows.TotalPages, shows.CurrentPage)
		}
		return tw.Flush()
	}
	fmt.Fprintln(tw, "Date:\tVenue:\tLocation:\tDuration:")
	for _, s := range shows.Shows {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", s.Date, s.VenueName, s.Venue.Location, s.Duration)
	}
	// the year details response prints a ShowsOutput but won't have any entries, for example
	if shows.TotalEntries != 0 {
		fmt.Fprintln(tw)
		fmt.Fprintf(tw, "Total Entries: %d\tTotal Pages: %d\tResult Page: %d\n", shows.TotalEntries, shows.TotalPages, shows.CurrentPage)
	}
	return tw.Flush()
}

type ShowResponse struct {
	Data Show `json:"data"`
}

func convertShowToOutput(show Show) ShowOutput {
	o := ShowOutput{
		ID:         show.ID,
		Date:       show.Date,
		Duration:   convertMillisecondToConcertDuration(int64(show.Duration)),
		Sbd:        show.Sbd,
		Remastered: show.Remastered,
		Tags:       show.Tags,
		VenueName:  show.VenueName,
	}
	o.Venue = convertVenueToOutput(show.Venue)
	tracks := convertTracksToOutput(show.Tracks)
	o.Tracks = tracks.Tracks
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
	ID         int           `json:"id"`
	Date       string        `json:"date"`
	Duration   string        `json:"duration"`
	Sbd        bool          `json:"sbd"`
	Remastered bool          `json:"remastered"`
	Tags       []Tag         `json:"tags"`
	Venue      VenueOutput   `json:"venue"`
	VenueName  string        `json:"venue_name"`
	Tracks     []TrackOutput `json:"tracks"`
}

func prettyPrintShow(tw *tabwriter.Writer, show ShowOutput, verbose bool) error {
	if verbose {
		fmt.Fprintln(tw, "ID:\tDate:\tVenue:\tLocation:\tDuration:\tSoundboard:\tRemastered:")
		sbd := trueAsYes(show.Sbd)
		r := trueAsYes(show.Remastered)
		fmt.Fprintf(tw, "%d\t%s\t%s\t%s\t%s\t%s\t%s\n", show.ID, show.Date, show.VenueName, show.Venue.Location, show.Duration, sbd, r)
		fmt.Fprintln(tw)
		if len(show.Tags) != 0 {
			fmt.Fprintln(tw, "Show Tags:")
			tagInfo := convertTagsToString(show.Tags)
			fmt.Fprintf(tw, "%s\n", tagInfo)
			fmt.Fprintln(tw)
		}
		// should always have tracks but worth a check
		if len(show.Tracks) == 0 {
			return tw.Flush()
		}
		fmt.Fprintf(tw, "%s\n", show.Tracks[0].SetName)
		for i, t := range show.Tracks {
			if i > 1 && t.SetName != show.Tracks[i-1].SetName {
				fmt.Fprintln(tw)
				fmt.Fprintf(tw, "%s\t\n", show.Tracks[i].SetName)
			}
			tagInfo := convertTagsToString(t.Tags)
			fmt.Fprintf(tw, "%s\t%s\t%s\n", t.Title, t.Duration, tagInfo)
		}
		fmt.Fprintln(tw)
		fmt.Fprintln(tw, "Mp3:")
		for _, t := range show.Tracks {
			fmt.Fprintf(tw, "%s\t%s\n", t.Title, t.Mp3)
		}
		return tw.Flush()
	}
	fmt.Fprintln(tw, "Date:\tVenue:\tLocation:")
	fmt.Fprintf(tw, "%s\t%s\t%s\n", show.Date, show.VenueName, show.Venue.Location)
	fmt.Fprintln(tw)
	// should always have tracks but worth a check
	if len(show.Tracks) == 0 {
		return tw.Flush()
	}
	fmt.Fprintf(tw, "%s\n", show.Tracks[0].SetName)
	for i, t := range show.Tracks {
		if i > 1 && t.SetName != show.Tracks[i-1].SetName {
			fmt.Fprintln(tw)
			fmt.Fprintf(tw, "%s\t\n", show.Tracks[i].SetName)
		}
		fmt.Fprintf(tw, "%s\t%s\t\n", t.Title, t.Duration)
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

func prettyPrintTracks(tw *tabwriter.Writer, tracks TracksOutput) error {
	fmt.Fprintln(tw, "ID:\tDate:\tVenue:\tLocation:\tTitle:\tMp3:")
	for _, t := range tracks.Tracks {
		fmt.Fprintf(tw, "%d\t%s\t%s\t%s\t%s\t%s\n", t.ID, t.ShowDate, t.VenueName, t.VenueLocation, t.Title, t.Mp3)
	}
	fmt.Fprintln(tw)
	if tracks.TotalEntries != 0 {
		fmt.Fprintf(tw, "Total Entries: %d\tTotal Pages: %d\tResult Page: %d\n", tracks.TotalEntries, tracks.TotalPages, tracks.CurrentPage)
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

func prettyPrintTrack(tw *tabwriter.Writer, track TrackOutput) error {
	fmt.Fprintln(tw, "ID:\tDate:\tVenue:\tLocation:\tTitle:\tDuration\tSet\tMp3")
	fmt.Fprintf(tw, "%d\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n", track.ID, track.ShowDate, track.VenueName, track.VenueLocation, track.Title, track.Duration, track.SetName, track.Mp3)
	fmt.Fprintln(tw)
	if len(track.Tags) != 0 {
		fmt.Fprintln(tw, "Tags")
		fmt.Fprintln(tw, "Name:\tGroup:\tNotes:")
		for _, t := range track.Tags {
			fmt.Fprintf(tw, "%s\t%s\t%s\n", t.Name, t.Group, t.Notes)
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

func prettyPrintTags(tw *tabwriter.Writer, tags TagsOutput) error {
	fmt.Fprintln(tw, "Name:\tDescription:\tGroup:")
	for _, t := range tags.Tags {
		fmt.Fprintf(tw, "%s\t%s\t%s\n", t.Name, t.Description, t.Group)
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

func prettyPrintTag(tw *tabwriter.Writer, tag TagListItemOutput) error {
	fmt.Fprintln(tw, "Name:\tDescription:\tGroup:")
	fmt.Fprintf(tw, "%s\t%s\t%s\n", tag.Name, tag.Description, tag.Group)
	fmt.Fprintln(tw)
	showIDs := make([]string, 0, len(tag.ShowIds))
	for _, i := range tag.ShowIds {
		showIDs = append(showIDs, strconv.Itoa(i))
	}
	trackIDs := make([]string, 0, len(tag.TrackIds))
	for _, i := range tag.TrackIds {
		trackIDs = append(trackIDs, strconv.Itoa(i))
	}
	fmt.Fprintf(tw, "Show IDs Where %s Appears\n", tag.Name)
	fmt.Fprintf(tw, "%s\n", strings.Join(showIDs, ", "))
	fmt.Fprintln(tw)
	fmt.Fprintf(tw, "Track IDs Where %s Appears\n", tag.Name)
	fmt.Fprintf(tw, "%s\n", strings.Join(trackIDs, ", "))

	return tw.Flush()
}

// todo: confirm showtags format
type SearchResponse struct {
	Data struct {
		ExactShow  Show          `json:"exact_show,omitempty"`
		OtherShows []Show        `json:"other_shows,omitempty"`
		ShowTags   []interface{} `json:"show_tags,omitempty"`
		Songs      []Song        `json:"songs,omitempty"`
		Tags       []TagListItem `json:"tags,omitempty"`
		Tours      []Tour `json:"tours,omitempty"`
		TrackTags  []TrackTag    `json:"track_tags,omitempty"`
		Tracks     []Track `json:"tracks,omitempty"`
		Venues     []Venue       `json:"venues,omitempty"`
	} `json:"data,omitempty"`
}

type TrackTagOutput struct {
	ID             int       `json:"id"`
	TrackID        int       `json:"track_id"`
	TagID          int       `json:"tag_id"`
	Notes          string    `json:"notes"`
	Transcript     string    `json:"transcript"`
}

func prettyPrintTrackTags(tw *tabwriter.Writer, tags []TrackTagOutput) error {
	fmt.Fprintln(tw, "ID:\tTrackID:\tTagID:\tNotes:\tTranscript:")
	for _, tag := range tags {
		fmt.Fprintf(tw, "%d\t%d\t%d\t%s\t%s\n", tag.ID, tag.TrackID, tag.TagID, tag.Notes, tag.Transcript)
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
				ID: t.ID,
				TrackID: t.TrackID,
				TagID: t.TagID,
				Notes: t.Notes,
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
		ExactShow *ShowOutput `json:"exact_show,omitempty"`
		OtherShows []ShowOutput        `json:"other_shows,omitempty"`
		ShowTags   []any `json:"show_tags,omitempty"`
		Songs      []SongOutput        `json:"songs,omitempty"`
		Tags       []TagListItemOutput `json:"tags,omitempty"`
		Tours      []TourOutput `json:"tours,omitempty"`
		TrackTags  []TrackTagOutput    `json:"track_tags,omitempty"`
		Tracks     []TrackOutput `json:"tracks,omitempty"`
		Venues     []VenueOutput       `json:"venues,omitempty"`
	} `json:"results"`
}

func prettyPrintSearch(tw *tabwriter.Writer, search SearchOutput) error {
	fmt.Fprintln(tw, "search results:")
	fmt.Fprintln(tw)
	if search.Results.ExactShow != nil {
		fmt.Fprintln(tw, "exact show results:")
		if err := prettyPrintShow(tw, *search.Results.ExactShow, true); err != nil {
			return err
		}
		fmt.Fprintln(tw)
	}
	if len(search.Results.OtherShows) != 0 {
		fmt.Fprintln(tw, "shows results:")
		if err := prettyPrintShows(tw, ShowsOutput{Shows: search.Results.OtherShows}, true); err != nil {
			return err
		}
		fmt.Fprintln(tw)
	}
	if len(search.Results.ShowTags) != 0 {
		fmt.Fprintln(tw, "show tag results:")
		//
		fmt.Fprintln(tw)
	}
	if len(search.Results.Songs) != 0 {
		fmt.Fprintln(tw, "song results:")
		if err := prettyPrintSongs(tw, SongsOutput{Songs: search.Results.Songs}); err != nil {
			return err
		}
		fmt.Fprintln(tw)
	}
	if len(search.Results.Tags) != 0 {
		fmt.Fprintln(tw, "tag results:")
		if err := prettyPrintTags(tw, TagsOutput{Tags: search.Results.Tags}); err != nil {
			return err
		}
		fmt.Fprintln(tw)
	}
	if len(search.Results.Tours) != 0 {
		fmt.Fprintln(tw, "tour results:")
		if err := prettyPrintTours(tw, ToursOutput{Tours: search.Results.Tours}); err != nil {
			return err
		}
		fmt.Fprintln(tw)
	}
	if len(search.Results.TrackTags) != 0 {
		fmt.Fprintln(tw, "track tag results:")
		if err := prettyPrintTrackTags(tw, search.Results.TrackTags); err != nil {
			return err
		}
		fmt.Fprintln(tw)
	}
	if len(search.Results.Tracks) != 0 {
		fmt.Fprintln(tw, "track results:")
		if err := prettyPrintTracks(tw, TracksOutput{Tracks: search.Results.Tracks}); err != nil {
			return err
		}
		fmt.Fprintln(tw)
	}
	if len(search.Results.Venues) != 0 {
		fmt.Fprintln(tw, "venue results:")
		if err := prettyPrintVenues(tw, VenuesOutput{Venues: search.Results.Venues}); err != nil {
			return err
		}
		fmt.Fprintln(tw)
	}
	return nil
}