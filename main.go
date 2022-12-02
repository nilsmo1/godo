package main
import "fmt"
import "encoding/json"
import "os"
import gc "github.com/gbin/goncurses"

type Lists struct {
    Lists []Items `json:"lists"`
}

type Items struct {
    Title   string `json:"title"`
    Items []Item   `json:"items"`
}

type Item struct {
    Title       string `json:"title"`
    Description string `json:"description"`
    Status      bool   `json:"status"`
}

func print_item(scr *gc.Window, item Item, row int, idx int) {
    var done_char rune = ' '
    if row == idx { scr.AttrSet(gc.A_STANDOUT) }
    if item.Status { done_char = 'x' }
    scr.Printf("[%c] %d. Title: %s\n", done_char, idx+1, item.Title)
    scr.AttrSet(gc.A_NORMAL)
    scr.Printf("\tDescription:\n\t\t %s\n", item.Description)
}

func print_todos(scr *gc.Window, items Items) Items {
    var row int = 0
    scr.Clear()
    scr.Printf("TODO LIST: %s\n", items.Title)
    for idx, item := range items.Items { print_item(scr, item, row, idx) }
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
        for idx, item := range items.Items { print_item(scr, item, row, idx) }
        scr.Refresh()
    }
    return items
}

func print_list(scr *gc.Window, list Items, row int, idx int) {
    if row == idx { scr.AttrSet(gc.A_STANDOUT) }
    scr.Printf("Title: %s\n", list.Title)
    scr.AttrSet(gc.A_NORMAL)
    
}

func get_list_idx(scr *gc.Window, lists Lists) int { 
    var row int = 0
    for idx, list := range lists.Lists { print_list(scr, list, row, idx) }
    for {
        k := gc.KeyString(scr.GetChar())
        if k == "q"             { break }
        if k == "s" || k == "B" { if row+1 < len(lists.Lists) { row++ } }
        if k == "w" || k == "A" { if row   > 0 { row-- } }
        if k == "enter" { 
            return row
        }
        scr.Clear()
        for idx, list := range lists.Lists { print_list(scr, list, row, idx) }
        scr.Refresh()
    }
    
    return -1
}


func main() {
    json_content, read_err := os.ReadFile("test.json")
    if read_err != nil { fmt.Println(read_err); return }
    
    var lists Lists
    unmarshal_err := json.Unmarshal(json_content, &lists)
    if unmarshal_err != nil { fmt.Println(unmarshal_err); return }

    scr, init_err := gc.Init();
    if init_err != nil { fmt.Print(init_err); return }
    defer gc.End()

    gc.Echo(false)
    gc.Cursor(0)

    var list_idx int
    list_idx = get_list_idx(scr, lists)
    if list_idx == -1 { gc.End(); os.Exit(0) }

    var items Items
    items = lists.Lists[list_idx]
    lists.Lists[list_idx] = print_todos(scr, items)
    b, _ := json.MarshalIndent(lists, "", "\t")
    os.WriteFile("test.json", b, 0644)
}
