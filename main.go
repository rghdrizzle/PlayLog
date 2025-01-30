package main


import (

  "github.com/rivo/tview"
  "fmt"
  "log"
  "os"
  "encoding/json"
  "strconv"
)

type Game struct{
  Name string `json:"name"`
  Status string `json:"status"`
  Score int `json:"score"`
}

var (
  GameLog = []Game{}
  gameFile = "gamelist.json"
)


// Load the gameList from the json file
func LoadGameLogFile(){
  if _, err := os.Stat(gameFile); err==nil{
    data , err := os.ReadFile(gameFile)
    if err!=nil{
      log.Fatal("Error reading json file:",err)
    }
    json.Unmarshal(data,&GameLog)
  }
}

// Save the gameFile
func SaveGameLog(){
  data , err := json.MarshalIndent(GameLog,""," ") // This is just making json pretty
  if err!=nil{
    log.Fatal(err)
  }
  os.WriteFile(gameFile,data,0644)
}

// Delete game func
func DeletGameLog(index int){
  if index <0 || index >=len(GameLog){
    fmt.Printf("Invalid index")
    return
  }
  GameLog = append(GameLog[:index],GameLog[index+1:]...)
  SaveGameLog()
}

func main(){
  app := tview.NewApplication()
  LoadGameLogFile()
  gameList := tview.NewTextView().SetDynamicColors(true).SetWordWrap(true)

  gameList.SetBorder(true).SetTitle("PlayLog!")

  refreshGameList := func(){
    gameList.Clear()
    if len(GameLog)==0{
      fmt.Fprintln(gameList,"No items in Play log")
    }else{
      for i, game := range GameLog{
        fmt.Fprintf(gameList,"[%d] %s (Status: %s) (Score: %d/10)",i+1,game.Name, game.Status, game.Score)
      }
    }
    
  }

  gameNameInput := tview.NewInputField().SetLabel("Name")
  gameStatusInput := tview.NewDropDown().SetLabel("Statis").AddOption("To Play",func(){}).AddOption("In Progress", func(){}).AddOption("Finished", func(){})
  gameScoreInput := tview.NewInputField().SetLabel("Score").SetAcceptanceFunc(tview.InputFieldInteger)
  gameIDInput := tview.NewInputField().SetLabel("Game ID to delete:")


  form := tview.NewForm().AddFormItem(gameNameInput).AddFormItem(gameStatusInput).AddFormItem(gameScoreInput).AddFormItem(gameIDInput).AddButton("Add",func(){
    name := gameNameInput.GetText()
    _,status := gameStatusInput.GetCurrentOption()
    score := gameScoreInput.GetText()
    if name !="" && score!= ""{
      scoreInt , err := strconv.Atoi(score)
      if err!=nil{
        fmt.Printf("Invalid score")
        return
      }
      fmt.Printf(status)
      GameLog = append(GameLog, Game{Name:name,Status:status,Score:scoreInt})
      SaveGameLog()
      refreshGameList()
      gameNameInput.SetText("")
      gameScoreInput.SetText("")
    }
  }).AddButton("Delete",func(){
      idStr := gameIDInput.GetText()
      if idStr ==""{
        fmt.Fprintln(gameList, "Please enter an game ID to delete.")
				return
      }
      id , err := strconv.Atoi(idStr)
      if err!=nil || id<1 || id>len(GameLog){
        fmt.Fprintln(gameList, "Invalid game ID")
        return
      }
      DeletGameLog(id-1)
      refreshGameList()
      gameIDInput.SetText("")
  }).AddButton("Exit",func(){
    app.Stop()
  })

  form.SetBorder(true).SetTitle("Add a Game").SetTitleAlign(tview.AlignLeft)

  flex := tview.NewFlex().AddItem(gameList,0,1,true).AddItem(form,0,2,false)

  refreshGameList()
  tview.NewImage()

  app.EnableMouse(true)

  if err :=app.SetRoot(flex,true).Run(); err!=nil{
    log.Fatal(err)
  }
}
