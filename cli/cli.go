package cli

import (
	"encoding/json"
	"fmt"
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

func prettyPrintEras(e ErasOutput) {
    fmt.Printf(
        "Eras\n1.0: %v\n2.0: %v\n3.0: %v\n4.0: %v\n", strings.Join(e.One, ", "), strings.Join(e.Two, ", "),strings.Join(e.Three, ", "), strings.Join(e.Four, ", "),
    )
}

func printJSONEras(e ErasOutput) error {
    b, err := json.MarshalIndent(&e, "", "  ")
    if err != nil {
        return err
    }
    fmt.Println(string(b))
    return nil
}

type EraResponse struct {
    Era []string `json:"data"`
}

type EraOutput struct {
    Era string
    EraList []string `json:"era"`
}

func prettyPrintEra(e EraOutput) {
    fmt.Printf("Era %s:\n%s\n", e.Era, strings.Join(e.EraList, ", "))
}

func printJSONEra(e EraOutput) error {
    b, err := json.MarshalIndent(&e, "", "  ")
    if err != nil {
        return err
    }
    fmt.Println(string(b))
    return nil
}