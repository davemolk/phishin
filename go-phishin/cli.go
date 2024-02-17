package phishin

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
	"time"
)

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

func printJSONEras(w io.Writer, e ErasOutput) error {
    b, err := json.MarshalIndent(&e, "", "  ")
    if err != nil {
        return err
    }
    fmt.Fprintln(w, string(b))
    return nil
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

func printJSONEra(w io.Writer, e EraOutput) error {
    b, err := json.MarshalIndent(&e, "", "  ")
    if err != nil {
        return err
    }
    fmt.Fprintln(w, string(b))
    return nil
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

func printJSONYears(w io.Writer, years YearsOutput) error {
    b, err := json.MarshalIndent(&years, "", "  ")
    if err != nil {
        return err
    }
    fmt.Fprintln(w, string(b))
    return nil
}

type YearResponse struct {
    Data         []ConcertInfo `json:"data"`
}

type YearOutput struct {
    ConcertInfo         []ConcertInfo `json:"year"`
}

type ConcertInfo struct {
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

func printJSONYear(w io.Writer, years YearOutput) error {
    b, err := json.MarshalIndent(&years, "", "  ")
    if err != nil {
        return err
    }
    fmt.Fprintln(w, string(b))
    return nil
}

func prettyPrintYear(tw *tabwriter.Writer, year YearOutput, verbose bool) error {
    if verbose {
        fmt.Fprintln(tw, "Date:\tVenue:\tLocation:\tDuration:\tSoundboard:\tRemastered:")
        for _, y := range year.ConcertInfo {
            s := soundTreatment(y.Sbd)
            r := soundTreatment(y.Remastered)
            fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%s\n", y.Date, y.VenueName, y.Venue.Location, convertMillisecondToConcertDuration(int64(y.Duration)), s, r)
        }
        return tw.Flush()
    }
    fmt.Fprintln(tw, "Date:\tVenue:\tLocation:")
    for _, y := range year.ConcertInfo {
        fmt.Fprintf(tw, "%s\t%s\t%s\n", y.Date, y.VenueName, y.Venue.Location)
    }
    return tw.Flush()
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
    Alias       any       `json:"alias"`
    Original    bool      `json:"original"`
    Artist      any       `json:"artist"`
    Lyrics      string    `json:"lyrics"`
    TracksCount int       `json:"tracks_count"`
    UpdatedAt   time.Time `json:"updated_at"`
    Tracks      []Track `json:"tracks"`
}

type SongsResponse struct {
	Data         []Song `json:"data"`
}

type SongResponse struct {
    Data         Song `json:"data"`
}

type ConcertInfoBasic struct {
    ID         int       `json:"id"`
    Date       string    `json:"date"`
    Duration   int       `json:"duration"`
    Incomplete bool      `json:"incomplete"`
    Sbd        bool      `json:"sbd"`
    Remastered bool      `json:"remastered"`
    TourID     int       `json:"tour_id"`
    VenueID    int       `json:"venue_id"`
    LikesCount int       `json:"likes_count"`
    TaperNotes string    `json:"taper_notes"`
    UpdatedAt  time.Time `json:"updated_at"`
    VenueName  string    `json:"venue_name"`
    Location   string    `json:"location"`
}

type Tour struct {
    ID         int    `json:"id"`
    Name       string `json:"name"`
    ShowsCount int    `json:"shows_count"`
    Slug       string `json:"slug"`
    StartsOn   string `json:"starts_on"`
    EndsOn     string `json:"ends_on"`
    Shows      []ConcertInfoBasic `json:"shows"`
    UpdatedAt time.Time `json:"updated_at"`
}

type ToursResponse struct {
	Data []Tour `json:"data"`
}

type TourResponse struct {
    Data Tour `json:"data"`
}

type VenuesResponse struct {
	Data         []Venue `json:"data"`
}

type VenueResponse struct {
    Data         Venue `json:"data"`
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
	Data         []ConcertInfo `json:"data"`
}

type ShowResponse struct {
    
}

type ShowOnDateResponse struct {
    Data ConcertInfo `json:"data"`
}

type ShowsOnDayOfYear struct {
	Data []ConcertInfo `json:"data"`
}

type RandomShowResponse struct {
    Data ConcertInfo `json:"data"`
}

type TracksResponse struct {
    Data         []Track `json:"data"`
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
}

type TagsResponse struct{
    Data         []TagListItem `json:"data"`
}





// type Venue struct {
//     ID         int           `json:"id"`
//     Slug       string        `json:"slug"`
//     Name       string        `json:"name"`
//     OtherNames []string `json:"other_names"`
//     Latitude   float64       `json:"latitude"`
//     Longitude  float64       `json:"longitude"`
//     ShowsCount int           `json:"shows_count"`
//     Location   string        `json:"location"`
//     UpdatedAt  time.Time     `json:"updated_at"`
// }


// type ConcertTag struct {
//     ID       int         `json:"id"`
//     Name     string      `json:"name"`
//     Priority int         `json:"priority"`
//     Group    string      `json:"group"`
//     Color    string      `json:"color"`
//     Notes    string `json:"notes"`
// }