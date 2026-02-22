package cli

import "fmt"

// Small helpers for the CLI UX can live here.
func Confirm(prompt string) bool {
    var ans string
    fmt.Printf("%s [y/N]: ", prompt)
    _, err := fmt.Scanln(&ans)
    if err != nil {
        return false
    }
    return ans == "y" || ans == "Y"
}
