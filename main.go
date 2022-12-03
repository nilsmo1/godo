package main
import "fmt"
import "encoding/json"
import "os"
import gc "github.com/gbin/goncurses"

type Lists struct { Lists []Items `json:"lists"` }

type Items struct {
    Title   string `json:"title"`
    Items []Item   `json:"items"`
}

type Item struct {
    Task        string `json:"title"`
    Description string `json:"description"`
    Status      bool   `json:"status"`
}

func print_item(scr *gc.Window, item Item, row int, idx int) {
    var done_char rune = ' '
    if row == idx { scr.AttrSet(gc.A_STANDOUT) }
    if item.Status { done_char = 'x' }
    scr.Printf("[%c] %d. Task: %s\n", done_char, idx+1, item.Task)
    scr.AttrSet(gc.A_NORMAL)
    scr.Printf("\tDescription:\n\t\t %s\n", item.Description)
}

func print_todos(scr *gc.Window, items *Items) {
    var row int = 0
    scr.Clear()
    scr.Printf("TODO LIST: %s\n", items.Title)
    var k string
    for {
        scr.Clear()
        scr.Printf("TODO LIST: %s\n", items.Title)
        for idx, item := range items.Items { print_item(scr, item, row, idx) }
        scr.Refresh()
        k = gc.KeyString(scr.GetChar())
        switch k {
        case "q": return
        case "s", "B": if row+1 < len(items.Items) { row++ }
        case "w", "A": if row   > 0                { row-- }
        case "enter": (*items).Items[row].Status = !items.Items[row].Status
        case "d": 
            if len(items.Items) == 0 { break }
            (*items).Items = append(items.Items[:row], items.Items[row+1:]...)
            if row >= len(items.Items) { row = len(items.Items)-1 }
        case "n": 
            item := Item { "New task", "new task description", false }
            (*items).Items = append(items.Items, item)
            if row < 0 { row = 0 }
        }
    }
}

func print_list(scr *gc.Window, list Items, row int, idx int) {
    if row == idx { scr.AttrSet(gc.A_STANDOUT) }
    scr.Printf("Title: %s\n", list.Title)
    scr.AttrSet(gc.A_NORMAL)
    
}

func get_list_idx(scr *gc.Window, lists *Lists) int { 
    var row int = 0
    var k string
    var ret bool
    for {
        scr.Clear()
        for idx, list := range lists.Lists { print_list(scr, list, row, idx) }
        scr.Refresh()
        k = gc.KeyString(scr.GetChar())
        switch k { 
        case "q": ret = true
        case "s", "B": if row+1 < len(lists.Lists) { row++ } 
        case "w", "A": if row   > 0                { row-- } 
        case "enter": return row 
        case "d": 
            if len(lists.Lists) == 0 { break }
            (*lists).Lists = append(lists.Lists[:row], lists.Lists[row+1:]...)
            if row >= len(lists.Lists) { row = len(lists.Lists)-1 }
        case "n": 
            item  := Item  { "New task", "new task description", false }
            items := Items { "New todo-list", []Item{item} }
            (*lists).Lists = append(lists.Lists, items)
            if row < 0 { row = 0 }
        }
        if ret { break }
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
    var items Items
    for {
        scr.Clear()
        list_idx = get_list_idx(scr, &lists)
        if list_idx == -1 { break }

        items = lists.Lists[list_idx]
        print_todos(scr, &items)
        lists.Lists[list_idx] = items
    }
    b, _ := json.MarshalIndent(lists, "", "\t")
    os.WriteFile("test.json", b, 0644)
}
