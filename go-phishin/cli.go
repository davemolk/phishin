package phishin

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
        return err
    }
    fmt.Fprintln(w, string(b))
    return nil
}

type ErasResponse struct {
    Data struct {
        One []string `json:"1.0"`
        Two []string `json:"2.0"`
        Three []string `json:"3.0"`
        Four []string `json:"4.0"`
    } `json:"data"`
}

type ErasOutput struct {
    One []string `json:"1.0"`
    Two []string `json:"2.0"`
    Three []string `json:"3.0"`
    Four []string `json:"4.0"`
}

func prettyPrintEras(w io.Writer, e ErasOutput) error {
    _, err := fmt.Fprintf(w,
        "Eras\n1.0: %v\n2.0: %v\n3.0: %v\n4.0: %v\n", strings.Join(e.One, ", "), strings.Join(e.Two, ", "),strings.Join(e.Three, ", "), strings.Join(e.Four, ", "),
    )
    return err
}

type EraResponse struct {
    Era []string `json:"data"`
}

type EraOutput struct {
    EraName string `json:"era"`
    Years []string `json:"years"`
}

func prettyPrintEra(w io.Writer, e EraOutput) error {
    _, err := fmt.Fprintf(w, "Era %s:\n%s\n", e.EraName, strings.Join(e.Years, ", "))
    return err
}

type Year struct {
    Date string `json:"date"`
    ShowCount int `json:"show_count"`
}

type YearsResponse struct {
    Data []Year `json:"data"`
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
    Data         []Show `json:"data"`
}

type Show struct {
    ID         int    `json:"id"`
    Date       string `json:"date"`
    Duration   int    `json:"duration"`
    Incomplete bool   `json:"incomplete"`
    Sbd        bool   `json:"sbd"`
    Remastered bool   `json:"remastered"`
    Tags       []Tag `json:"tags"`
    TourID int `json:"tour_id"`
    Venue  Venue `json:"venue"`
    VenueName  string `json:"venue_name"`
    TaperNotes string `json:"taper_notes"`
    LikesCount int    `json:"likes_count"`
    VenueID int `json:"venue_id"`
    Location string `json:"location"`
    Tracks     []Track `json:"tracks"`
    UpdatedAt time.Time `json:"updated_at"`
}

type Track struct {
    ID                int         `json:"id"`
    ShowID            int         `json:"show_id"`
    ShowDate          string      `json:"show_date"`
    VenueName         string      `json:"venue_name"`
    VenueLocation     string      `json:"venue_location"`
    Title             string      `json:"title"`
    Position          int         `json:"position"`
    Duration          int         `json:"duration"`
    JamStartsAtSecond int `json:"jam_starts_at_second"`
    Set               string      `json:"set"`
    SetName           string      `json:"set_name"`
    LikesCount        int         `json:"likes_count"`
    Slug              string      `json:"slug"`
    Tags              []Tag `json:"tags"`
    Mp3           string    `json:"mp3"`
    WaveformImage string    `json:"waveform_image"`
    SongIds       []int     `json:"song_ids"`
    UpdatedAt     time.Time `json:"updated_at"`
}

type Tag struct {
    ID             int         `json:"id"`
    Name           string      `json:"name"`
    Priority       int         `json:"priority"`
    Group          string      `json:"group"`
    Color          string      `json:"color"`
    Notes          string `json:"notes"`
    Transcript     string `json:"transcript"`
    StartsAtSecond int `json:"starts_at_second"`
    EndsAtSecond   int `json:"ends_at_second"`
}

type soundTreatment bool

func (s soundTreatment) String() string {
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

type Song struct {
    ID          int       `json:"id"`
    Slug        string    `json:"slug"`
    Title       string    `json:"title"`
    Alias       string       `json:"alias"`
    Original    bool      `json:"original"`
    Artist      string       `json:"artist"`
    Lyrics      string    `json:"lyrics"`
    TracksCount int       `json:"tracks_count"`
    UpdatedAt   time.Time `json:"updated_at"`
    Tracks      []Track `json:"tracks"`
}

type SongsResponse struct {
	Data         []Song `json:"data"`
}

type SongsOutput struct {
    Songs []Song `json:"songs"`
}

type SongResponse struct {
    Data         Song `json:"data"`
}

type SongOutput struct {
    Song `json:"song"`
}

func prettyPrintSongs(tw *tabwriter.Writer, songs SongsOutput) error {
    fmt.Fprintln(tw, "Title:\tAlias:\tOriginal:\tArtist:\tTracksCount")
    for _, s := range songs.Songs {
        // reuse bool stringer
        orig := soundTreatment(s.Original)
        fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%d\n", s.Title, s.Alias, orig, s.Artist, s.TracksCount)
    }
    return tw.Flush()
}

func prettyPrintSong(tw *tabwriter.Writer, song SongOutput) error {
    fmt.Fprintln(tw, "Title:\tAlias:\tOriginal:\tArtist:\tTracksCount")
    fmt.Fprintf(tw, "%s\t%s\t%v\t%s\t%d\n", song.Title, song.Alias, song.Original, song.Artist, song.TracksCount)
    fmt.Fprintln(tw)
    fmt.Println(tw, "Tracks")
    fmt.Fprintln(tw, "Date:\tVenue:\tLocation:\tMp3")
    for _, t := range song.Tracks {
        fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", t.ShowDate, t.VenueName, t.VenueLocation, t.Mp3)
    }
    return tw.Flush()
}

type Tour struct {
    ID         int    `json:"id"`
    Name       string `json:"name"`
    ShowsCount int    `json:"shows_count"`
    Slug       string `json:"slug"`
    StartsOn   string `json:"starts_on"`
    EndsOn     string `json:"ends_on"`
    Shows      []Show `json:"shows"`
    UpdatedAt time.Time `json:"updated_at"`
}

type ToursResponse struct {
	Data []Tour `json:"data"`
}

type ToursOutput struct {
	Tours []Tour `json:"tours"`
}

type TourResponse struct {
    Data Tour `json:"data"`
}

type TourOutput struct {
    Tour `json:"tour"`
}

func prettyPrintTours(tw *tabwriter.Writer, tours ToursOutput) error {
    fmt.Fprintln(tw, "Name:\tStarts On:\tEnds On:\tShow Count:")
    for _, t := range tours.Tours {
        fmt.Fprintf(tw, "%s\t%s\t%s\t%d\n", t.Name, t.StartsOn, t.EndsOn, t.ShowsCount)
    }
    return tw.Flush()
}

func prettyPrintTour(tw *tabwriter.Writer, tour TourOutput) error {
    fmt.Fprintln(tw, "Name:\tStarts On:\tEnds On:\tShow Count:")
    fmt.Fprintf(tw, "%s\t%s\t%s\t%d\n", tour.Name, tour.StartsOn, tour.EndsOn, tour.ShowsCount)
    fmt.Fprintln(tw)
    fmt.Fprintln(tw, "Date:\tVenue:\tLocation:\tDuration:\tSoundboard:\tRemastered:")
    for _, show := range tour.Shows {
        sbd := soundTreatment(show.Sbd)
        r := soundTreatment(show.Remastered)
        fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%s\n", show.Date, show.VenueName, show.Location, convertMillisecondToConcertDuration(int64(show.Duration)), sbd, r)
        fmt.Fprintln(tw)
    }
    return tw.Flush()
}

type VenuesResponse struct {
	Data         []Venue `json:"data"`
}

type VenuesOutput struct {
	Venues         []Venue `json:"venues"`
}

func prettyPrintVenues(tw *tabwriter.Writer, venues VenuesOutput) error {
    fmt.Fprintln(tw, "Venues:\tLocation:\tShow Count:")
    for _, v := range venues.Venues {
        fmt.Fprintf(tw, "%s\t%s\t%d\n", v.Name, v.Location, v.ShowsCount)
    }
    return tw.Flush()
}

type VenueResponse struct {
    Data         Venue `json:"data"`
}

type VenueOutput struct {
    Venue `json:"venue"`
}

func prettyPrintVenue(tw *tabwriter.Writer, venue VenueOutput) error {
    fmt.Fprintln(tw, "Venue:\tLocation:\tShow Count:")
    fmt.Fprintf(tw, "%s\t%s\t%d\n", venue.Name, venue.Location, venue.ShowsCount)
    fmt.Fprintln(tw, "\t")
    if len(venue.ShowDates) == 0 {
        return tw.Flush()
    }
    fmt.Fprintln(tw, "Show Dates")
    for _, d := range venue.ShowDates {
        fmt.Fprintf(tw, "%s\n", d)
    }
    return tw.Flush()
}

// used in venues list, details. subset used in years
type Venue struct {
    ID         int       `json:"id"`
    Slug       string    `json:"slug"`
    Name       string    `json:"name"`
    OtherNames []string     `json:"other_names"`
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

type ShowsResponse struct {
	Data         []Show `json:"data"`
}

type ShowsOutput struct {
    Shows []Show `json:"shows"`
}

type ShowResponse struct {
    Data         Show `json:"data"`
}

type ShowOutput struct {
    Show `json:"show"`
}

func prettyPrintShows(tw *tabwriter.Writer, shows ShowsOutput, verbose bool) error {
    if verbose {
        fmt.Fprintln(tw, "Date:\tVenue:\tLocation:\tDuration:\tSoundboard:\tRemastered:")
        for _, s := range shows.Shows {
            sbd := soundTreatment(s.Sbd)
            r := soundTreatment(s.Remastered)
            fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%s\n", s.Date, s.VenueName, s.Venue.Location, convertMillisecondToConcertDuration(int64(s.Duration)), sbd, r)
        }
        return tw.Flush()
    }
    fmt.Fprintln(tw, "Date:\tVenue:\tLocation:")
    for _, s := range shows.Shows {
        fmt.Fprintf(tw, "%s\t%s\t%s\n", s.Date, s.VenueName, s.Venue.Location)
    }
    return tw.Flush()
}

func prettyPrintShow(tw *tabwriter.Writer, show ShowOutput, verbose bool) error {
    if verbose {
        fmt.Fprintln(tw, "Date:\tVenue:\tLocation:\tDuration:\tSoundboard:\tRemastered:")
        sbd := soundTreatment(show.Sbd)
        r := soundTreatment(show.Remastered)
        fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%s\n", show.Date, show.VenueName, show.Venue.Location, convertMillisecondToConcertDuration(int64(show.Duration)), sbd, r)
        fmt.Fprintln(tw, "\t")
        // should always have tracks but check
        if len(show.Tracks) == 0 {
            return tw.Flush()
        }
        fmt.Fprintf(tw, "%s\n", show.Tracks[0].SetName)
        for i, t := range show.Tracks {
            if i > 1 && t.SetName != show.Tracks[i-1].SetName {
                fmt.Fprintln(tw, "\t")
                fmt.Fprintf(tw, "%s\t\n", show.Tracks[i].SetName)
            }
            var tags []string
            for _, t := range t.Tags {
                if t.Notes != "" {
                    // sometimes inserted mid-text
                    notes := strings.ReplaceAll(t.Notes, "\n", "")
                    notes = strings.ReplaceAll(notes, "\r", "")
                    tags = append(tags, fmt.Sprintf("%s: %s", t.Name, notes))
                } else {
                    tags = append(tags, t.Name)
                }
            }
            tagInfo := strings.Join(tags, ", ")
            fmt.Fprintf(tw, "%s\t%s\t%s\n", t.Title, convertMillisecondToConcertDuration(int64(t.Duration)), tagInfo)
        }
        fmt.Fprintln(tw, "\t")
        fmt.Fprintln(tw, "want to listen?",)
        for _, t := range show.Tracks {
            fmt.Fprintf(tw, "%s\t%s\n", t.Title, t.Mp3)
        }
        return tw.Flush()
    }
    fmt.Fprintln(tw, "Date:\tVenue:\tLocation:")
    fmt.Fprintf(tw, "%s\t%s\t%s\n", show.Date, show.VenueName, show.Venue.Location)
    fmt.Fprintln(tw, "\t")
    // should always have tracks but check
    if len(show.Tracks) == 0 {
        return tw.Flush()
    }
    fmt.Fprintf(tw, "%s\n", show.Tracks[0].SetName)
    for i, t := range show.Tracks {
        if i > 1 && t.SetName != show.Tracks[i-1].SetName {
            fmt.Fprintln(tw, "\t")
            fmt.Fprintf(tw, "%s\t\n", show.Tracks[i].SetName)
        }
        fmt.Fprintf(tw, "%s\t%s\t\n", t.Title, convertMillisecondToConcertDuration(int64(t.Duration)))
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
    Data         []Track `json:"data"`
}

type TracksOutput struct {
    Tracks []Track `json:"tracks"`
}

type TrackResponse struct {
    Data Track `json:"data"`
}

type TrackOutput struct {
    Track `json:"track"`
}

func prettyPrintTracks(tw *tabwriter.Writer, tracks TracksOutput) error {
    fmt.Fprintln(tw, "Date:\tVenue:\tLocation:\tTitle:\tMp3")
    for _, t := range tracks.Tracks {
        fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\n", t.ShowDate, t.VenueName, t.VenueLocation, t.Title, t.Mp3)
    }
    return tw.Flush()
}

func prettyPrintTrack(tw *tabwriter.Writer, track TrackOutput) error {
    fmt.Fprintln(tw, "Date:\tVenue:\tLocation:\tTitle:\tDuration\tSet\tMp3")
    fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n", track.ShowDate, track.VenueName, track.VenueLocation, track.Title, convertMillisecondToConcertDuration(int64(track.Duration)), track.SetName, track.Mp3)
    fmt.Fprintln(tw)
    if len(track.Track.Tags) != 0 {
        fmt.Println(tw, "Tags")
        fmt.Fprintln(tw, "Name:\tGroup:\tNotes:\tTranscript")
        for _, t := range track.Tags {
            fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", t.Name, t.Group, t.Notes, t.Transcript)
        }
    }
    return tw.Flush()
}

type TagListItem struct {
    ID          int       `json:"id"`
    Name        string    `json:"name"`
    Slug        string    `json:"slug"`
    Group       string    `json:"group"`
    Color       string    `json:"color"`
    Priority    int       `json:"priority"`
    Description string    `json:"description"`
    UpdatedAt   time.Time `json:"updated_at"`
    ShowIds     []int `json:"show_ids"`
    TrackIds    []int         `json:"track_ids"`
}

type TagsResponse struct{
    Data         []TagListItem `json:"data"`
}

type TagsOutput struct {
    Tags []TagListItem `json:"tags"`
}

type TagResponse struct {
    Data  TagListItem `json:"data"`
}

type TagOutput struct {
    TagListItem `json:"tag"`
}

func prettyPrintTags(tw *tabwriter.Writer, tags TagsOutput) error {
    fmt.Fprintln(tw, "Name:\tDescription:\tGroup:")
    for _, t := range tags.Tags {
        fmt.Fprintf(tw, "%s\t%s\t%s\n", t.Name, t.Description, t.Group)
    }
    return tw.Flush()
}

func prettyPrintTag(tw *tabwriter.Writer, tag TagOutput) error {
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

// wasn't able to get results for other fields
type SearchResponse struct {
    Data         struct {
        ExactShow  Show  `json:"exact_show"`
        OtherShows []Show `json:"other_shows"`
        ShowTags   []interface{} `json:"show_tags"`
        Songs      []Song `json:"songs"`
        Tags      []TagListItem `json:"tags"`
        Tours     []interface{} `json:"tours"`
        TrackTags []TrackTag `json:"track_tags"`
        Tracks []interface{} `json:"tracks"`
        Venues []Venue `json:"venues"`
    } `json:"data"`
}

type TrackTag struct {
    ID             int         `json:"id"`
    TrackID        int         `json:"track_id"`
    TagID          int         `json:"tag_id"`
    CreatedAt      time.Time   `json:"created_at"`
    Notes          string      `json:"notes"`
    StartsAtSecond int         `json:"starts_at_second"`
    EndsAtSecond   int `json:"ends_at_second"`
    Transcript     string      `json:"transcript"`
}
