package main
import "fmt"
import "encoding/json"
import "os"
import gc "github.com/gbin/goncurses"

type Items struct {
    Title   string `json:"title"`
    Items []Item   `json:"items"`
}

type Item struct {
    Title       string `json:"title"`
    Description string `json:"description"`
    Status      bool   `json:"status"`
}

func print_todos(scr *gc.Window, items Items, row int) Items {
    var done_char rune = ' '
    scr.Printf("TODO LIST: %s\n", items.Title)

    for idx, item := range items.Items {
        if row == idx { scr.AttrSet(gc.A_STANDOUT) }
        if item.Status { done_char = 'x' }
        scr.Printf("[%c] %d. Title: %s\n", done_char, idx+1, item.Title)
        scr.AttrSet(gc.A_NORMAL)
        scr.Printf("\tDescription: %s\n", item.Description)
    }
    for {
        k := gc.KeyString(scr.GetChar())
        if k == "q"             { break }
        if k == "s" || k == "B" { if row+1 < len(items.Items) { row++ } }
        if k == "w" || k == "A" { if row   > 0 { row-- } }
        if k == "enter" { 
            items.Items[row].Status = !items.Items[row].Status
        }
        scr.Clear()
        scr.Printf("TODO LIST: %s\n", items.Title)
        for idx, item := range items.Items {
            done_char = ' '
            if row == idx { scr.AttrSet(gc.A_STANDOUT) }
            if item.Status { done_char = 'x' }
            scr.Printf("[%c] %d. Title: %s\n",done_char, idx+1, item.Title)
            scr.AttrSet(gc.A_NORMAL)
            done_char = ' '
            scr.Printf("\tDescription: %s\n", item.Description)
        }
        scr.Refresh()
    }
    return items
}


func main() {
    json_content, read_err := os.ReadFile("test.json")
    if read_err != nil { fmt.Println(read_err); return }

    var items Items
    unmarshal_err := json.Unmarshal(json_content, &items)
    if unmarshal_err != nil { fmt.Println(unmarshal_err); return }

    scr, init_err := gc.Init();
    if init_err != nil { fmt.Print(init_err); return }
    defer gc.End()

    gc.Echo(false)
    gc.Cursor(0)

    var row int = 0
    items = print_todos(scr, items, row)
    
    b, _ := json.MarshalIndent(items, "", "\t")
    os.WriteFile("test.json", b, 0644)

}
