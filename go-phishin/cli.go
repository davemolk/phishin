package phishin

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
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
    Era string
    EraList []string `json:"years"`
}

func prettyPrintEra(w io.Writer, e EraOutput) error {
    _, err := fmt.Fprintf(w, "Era %s:\n%s\n", e.Era, strings.Join(e.EraList, ", "))
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