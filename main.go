package main
import "fmt"
import "encoding/json"
import "os"
import "strings"
import gc "github.com/gbin/goncurses"

var CHARS string = ` abcdefghijklmnopqrstuvwxyz
ABCDEFGHIJKLMNOPQRSTUVWXYZ
/%@,+<!|{:(-$>#'}=^"&*\;._?)`

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

func get_input(scr *gc.Window, title string) string {
    var wh, ww, nwh, nww, posy, posx int
    wh, ww = scr.MaxYX()
    nwh, nww   = 4, ww/4
    posy, posx = (wh-nwh)/2, (ww-nww)/2 

    input_window, window_err := gc.NewWindow(nwh, nww, posy, posx)
    if window_err != nil { gc.End(); fmt.Println(window_err) }
    defer gc.End()
    var k, buffer string
    var diff int
    buffer = ""
    for {
        input_window.Clear()
        input_window.Println()

        input_window.AttrSet(gc.A_BOLD)
        input_window.Printf("  %s:\n", title)
        input_window.AttrSet(gc.A_NORMAL)
        diff = len(buffer) + 5 - nww
        if diff > 0 {
            input_window.Printf("  ..%s", buffer[diff+2:])
        } else { input_window.Printf("  %s", buffer) }

        input_window.AttrSet(gc.A_UNDERLINE | gc.A_BLINK)
        input_window.Println(" ")
        input_window.AttrSet(gc.A_NORMAL)

        input_window.Box(gc.ACS_VLINE, gc.ACS_HLINE)

        kk := input_window.GetChar()
        k = gc.KeyString(kk)
        if k  == "enter" { break }
        if kk == 27 { return "" }
        if kk == 127 {
            if c := len(buffer); c != 0 {
                buffer = buffer[:c-1]
                continue
            } 
        }
        if strings.Contains(CHARS, k) { buffer+=k }
    }
    gc.End()
    return buffer
}

func print_title(scr *gc.Window, title string) {
    scr.AttrSet(gc.A_BOLD)
    scr.Printf("TODO-LIST: %s\n", title)
    scr.AttrSet(gc.A_NORMAL)
}

func print_item(scr *gc.Window, item Item, row int, idx int) {
    var done_char rune = ' '
    if row == idx { scr.AttrSet(gc.A_STANDOUT) }
    if item.Status { done_char = 'x' }
    scr.Printf("[%c] %d. Task: %s\n", done_char, idx+1, item.Task)
    scr.AttrSet(gc.A_NORMAL)
    if item.Description != "" {
        scr.AttrSet(gc.A_BOLD)
        scr.Print("  Description:\n")
        scr.AttrSet(gc.A_NORMAL)
        scr.Printf("    %s\n", item.Description)
    }
}



func print_todos(scr *gc.Window, items *Items) {
    var row int = 0
    scr.Clear()
    print_title(scr, items.Title)
    var k string
    for {
        scr.Clear()
        print_title(scr, items.Title)
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
            task_new := get_input(scr, "Task")
            if task_new == "" { break }
            desc_new := get_input(scr, "Description")
            item := Item { task_new, desc_new, false }
            (*items).Items = append(items.Items, item)
            if row < 0 { row = 0 }
        case "e":
            task_edit := get_input(scr, "Task")
            if task_edit == "" { break }
            description_edit := get_input(scr, "Description")
            if task_edit != "" { (*items).Items[row].Task = task_edit }
            (*items).Items[row].Description = description_edit
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
            title := get_input(scr, "Title")
            if title == "" { break }
            task_new := get_input(scr, "Task")
            desc_new := get_input(scr, "Description")
            if task_new == "" { task_new = "New task"}
            item  := Item  { task_new, desc_new, false }
            items := Items { title, []Item{item} }
            (*lists).Lists = append(lists.Lists, items)
            if row < 0 { row = 0 }
        case "e":
            list_edit := get_input(scr, "List")
            if list_edit != "" { (*lists).Lists[row].Title = list_edit }
        }
        if ret { break }
    }
    return -1
}


func main() {
    home, home_err := os.UserHomeDir()
    if home_err != nil { fmt.Println(home_err); return }
    path := home + "/.config/godo-lists.json"
    json_content, read_err := os.ReadFile(path)
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
    os.WriteFile(path, b, 0644)
}
