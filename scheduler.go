package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	_ "modernc.org/sqlite"
)

var TIMEBLOCKS = Database{
	Name: "TimeBlocks",
	ColumnNames: []string{
		"blockID",
		"ScheduleID",
		"dayOfWeek",
		"startTime",
		"endTime",
		"duration",
	},
	FileURL: DATABASE_URL + "timeblocks.db",
	CreationQuery: `
	CREATE TABLE IF NOT EXISTS TimeBlocks(
	blockID INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	ScheduleID INTEGER NOT NULL,
	dayOfWeek STRING NOT NULL CHECK (dayOfWeek IN ('Mon', 'Tue', 'Wed', 'Thur', 'Fri')),
	startTime TIME(0) NOT NULL,
	endTime TIME(0) NOT NULL,
	left INT NOT NULL DEFAULT 0,
	Duration string NOT NULL DEFAULT "0:30",
	FOREIGN KEY (ScheduleID) REFERENCES Schedules(ScheduleID)
	);
	`,
}

type TimeBlock struct {
	BlockID    int
	ScheduleID int
	DayOfWeek  string
	StartTime  string
	EndTime    string
	Left       int
	Duration   string
}

func newTimeBlock(w http.ResponseWriter, r *http.Request) {
	loggedIn := checkLoggedIn(w, r)
	if loggedIn != "LOGGEDIN" {
		fmt.Println("Could not make a new Time Block, user is not logged in.")
		return
	}
	tb := createNewTimeBlock(w, r)

	db, err := sql.Open("sqlite", TIMEBLOCKS.FileURL)
	if err != nil {
		fmt.Println("In the first one")
		fmt.Println(err)
		return
	}

	_, tableCreationErr := db.Exec(TIMEBLOCKS.CreationQuery)
	if tableCreationErr != nil {
		fmt.Println(tableCreationErr)
		return
	}
	fmt.Printf("About to create a new timeblock, %d", string(TIMEBLOCKS.FileURL))

	result, err := db.Exec(
		"INSERT INTO TimeBlocks (ScheduleID, dayOfWeek, startTime, endTime, left, duration) VALUES (?, ?, ?, ?, ?, ?)",
		tb.ScheduleID,
		tb.DayOfWeek,
		tb.StartTime,
		tb.EndTime,
		tb.Left,
		tb.Duration,
	)

	if err != nil {
		fmt.Printf("In the second one22, %d", TIMEBLOCKS.FileURL)
		fmt.Println(err)
		return
	}

	id, err := result.LastInsertId()
	fmt.Println(tb)

	if err != nil {
		fmt.Println("In the third one")
		fmt.Println(err)
		return
	}

	tb.BlockID = int(id)
	fmt.Println(tb)
	defer db.Close()

	// w.Write([]byte(`<div style="background:yellow;">TEST BLOCK</div>`))

	files :=
		[]string{
			HTML_FILES["TIMEBLOCK"],
		}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Executing %+v\n", tb)
	tbCapsule := []TimeBlock{
		tb,
	}
	err = ts.Execute(w, tbCapsule)
	if err != nil {
		fmt.Println(err)

	}
}

func createNewTimeBlock(w http.ResponseWriter, r *http.Request) TimeBlock {
	timeBlock := TimeBlock{}
	loggedIn := checkLoggedIn(w, r)
	if loggedIn != "LOGGEDIN" {
		return timeBlock
	}
	schedule := getScheduleFromAcc(w, r)
	if schedule.ScheduleID == -1 {
		createSchedule(w, r)
		schedule = getScheduleFromAcc(w, r)
	}
	timeBlock = TimeBlock{
		ScheduleID: schedule.ScheduleID,
		DayOfWeek:  "Mon",
		StartTime:  "8:00",
		EndTime:    "8:30",
		Left:       0,
		Duration:   "0:30",
	}
	return timeBlock
}

func updateTimeBlock(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update Called")

	newTimeBlocks := []TimeBlock{}

	err := json.NewDecoder(r.Body).Decode(&newTimeBlocks)

	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	db, err := sql.Open("sqlite", TIMEBLOCKS.FileURL)
	if err != nil {
		fmt.Println(err)
	}
	db.SetMaxOpenConns(1)
	_, err = db.Exec("PRAGMA busy_timeout = 5000")
	if err != nil {
		fmt.Print(err)
	}

	for _, newTimeBlock := range newTimeBlocks {
		fmt.Println("BlockID: ", newTimeBlock)

		query := `
	UPDATE TimeBlocks SET dayOfWeek = ?, StartTime = ?, EndTime = ?, Left = ?, Duration = ? WHERE blockID = ?
	`

		result, err := db.Exec(query, newTimeBlock.DayOfWeek, newTimeBlock.StartTime, newTimeBlock.EndTime, newTimeBlock.Left, newTimeBlock.Duration, newTimeBlock.BlockID)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("BlockID: ", newTimeBlock, result)

		result = result
	}
	defer db.Close()

	fmt.Println("Successfully Changed")
	return
}

func loadTimeBlocks(w http.ResponseWriter, r *http.Request) []TimeBlock {
	fmt.Println("Asked for all time blocks")

	err := r.ParseForm()
	if err != nil {
		fmt.Println("There was an error with the form", err)
		w.Header().Set("HX-Redirect", "/account/")
		w.WriteHeader(http.StatusBadRequest)
		return []TimeBlock{}
	}
	id, err := strconv.Atoi(r.FormValue("ScheduleID"))
	fmt.Println("This is the form value: ", r.FormValue("ScheduleID"))

	if err != nil {
		fmt.Println("There is an error with the value sent", err, id)
		w.Header().Set("HX-Redirect", "/")
		w.WriteHeader(http.StatusBadRequest)
		return []TimeBlock{}
	}

	sched := getScheduleFromSchedID(id)

	var timeBlocks []TimeBlock

	loggedIn := checkLoggedIn(w, r)
	if loggedIn != "LOGGEDIN" {
		return timeBlocks
	}

	db, err := sql.Open("sqlite", TIMEBLOCKS.FileURL)
	if err != nil {
		fmt.Println("In the first one")
		fmt.Println(err)
		return timeBlocks
	}

	query := `
	SELECT * FROM TimeBlocks WHERE scheduleID = ?;
	`

	result, err := db.Query(query, sched.ScheduleID)
	if err != nil {
		fmt.Println(err)
		return timeBlocks
	}
	defer result.Close()
	defer db.Close()

	for result.Next() {
		var timeBlock TimeBlock
		err := result.Scan(&timeBlock.BlockID, &timeBlock.ScheduleID, &timeBlock.DayOfWeek, &timeBlock.StartTime, &timeBlock.EndTime, &timeBlock.Left, &timeBlock.Duration)
		if err != nil {
			fmt.Println(err)
			return timeBlocks
		}

		timeBlocks = append(timeBlocks, timeBlock)
	}

	err = result.Err()

	if err != nil {
		fmt.Println(err)
		return timeBlocks
	}
	fmt.Print(timeBlocks)
	return timeBlocks
}

func loadTimeBlocksHandler(w http.ResponseWriter, r *http.Request) {
	timeBlocks := loadTimeBlocks(w, r)
	fmt.Println("Event handler called")

	files :=
		[]string{
			HTML_FILES["SAVEDTIMEBLOCK"],
		}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	fmt.Println("End of the Block")
	err = ts.Execute(w, timeBlocks)
	if err != nil {
		fmt.Println(err)
	}
}

func loadTimeBlocksEditHandler(w http.ResponseWriter, r *http.Request) {
	timeBlocks := loadTimeBlocks(w, r)
	fmt.Println("Event handler called")

	files :=
		[]string{
			HTML_FILES["TIMEBLOCK"],
		}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	fmt.Print("End of the edit handler")
	ts.Execute(w, timeBlocks)
}
