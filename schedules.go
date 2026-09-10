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

var SCHEDULE = Database{
	Name: "TimeBlocks",
	ColumnNames: []string{
		"ScheduleID",
		"AccountID",
		"ApprovalStatus",
		"Notes",
	},
	FileURL: DATABASE_URL + "schedules.db",
	CreationQuery: `
	CREATE TABLE IF NOT EXISTS Schedules(
	ScheduleID INTEGER PRIMARY KEY AUTOINCREMENT,
	AccountID INT NOT NULL,
	ApprovalStatus STRING DEFAULT "NOT APPROVED" CHECK (ApprovalStatus IN ("UNSUBMITTED", "PENDING", "APPROVED", "REVISION REQUIRED")),
	Notes STRING DEFAULT "",
	FOREIGN KEY (AccountID) REFERENCES Accounts(accID)
	);
	`,
}

/*
Struct blueprint
Name: Schedule
Values: ScheduleID, AccountID, ApprovalStatus
Explanation: This is the struct used to make an item from the Schedule database.
*/
type Schedule struct {
	ScheduleID     int
	AccountID      int
	ApprovalStatus string
	Notes          string
}

/*
Function Blueprint
Name: getScheduleFromAcc
Return: Schedule
Parameters: http.ResponseWriter, *http.Request
Purpose: It needs to go into the Schedules table and find the table that is connected to the account.
*/

func getScheduleFromAcc(w http.ResponseWriter, r *http.Request) Schedule {
	loggedIn := checkLoggedIn(w, r)
	schedule := Schedule{
		ScheduleID: -1,
	}

	if loggedIn != "LOGGEDIN" {
		fmt.Println("The User is not logged in, cannot access the schedule.")
		return schedule
	}

	db, err := sql.Open("sqlite", SCHEDULE.FileURL)
	if err != nil {
		fmt.Println("There was an error while opening the Schedules database")
		fmt.Println(err)
		return schedule
	}

	_, err = db.Exec(SCHEDULE.CreationQuery)

	if err != nil {
		fmt.Println("Could not create the Schedule table.")
		fmt.Println(err)
		return schedule
	}

	query := `
	SELECT * FROM Schedules WHERE AccountID = ?
	`
	result, err := db.Query(query, accountInfo(w, r).ID)

	if err != nil {
		fmt.Println("Error in the query for a schedule")
		fmt.Println(err)
		return schedule
	}

	for result.Next() {
		err := result.Scan(&schedule.ScheduleID, &schedule.AccountID, &schedule.ApprovalStatus, &schedule.Notes)
		if err != nil {
			fmt.Println("There was an error in the scanning: ", err)
			return schedule
		}
	}
	err = result.Err()
	if err != nil {
		fmt.Println(err)
		return schedule
	}

	if schedule.ScheduleID == -1 {
		createSchedule(w, r)
		fmt.Println("Creating a new schedule.")
		return getScheduleFromAcc(w, r)
	}

	fmt.Println("Successfully retreived the schedule.", schedule)
	return schedule
}

/*
getScheduleID
takes w and r as parameters and then doesn't return anything to the caller, but rather to the page.
*/

func getScheduleID(w http.ResponseWriter, r *http.Request) {

	sched := getScheduleFromAcc(w, r)
	fmt.Println("This is the scheduleID: ", sched.ScheduleID)

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(map[string]string{
		"ScheduleID": strconv.Itoa(sched.ScheduleID),
	})
	if err != nil {
		fmt.Println("There was an error retrieving the data.", err)
	}
	fmt.Println("Successfully retrieved the data from", sched)
}

/*
Name: getScheduleFromSchedID
Parameters: int id
Purpose: This will return a schedule based on the schedule's ID number
*/
func getScheduleFromSchedID(id int) Schedule {
	schedule := Schedule{
		ScheduleID: -1,
	}

	db, err := sql.Open("sqlite", SCHEDULE.FileURL)
	if err != nil {
		fmt.Println("There was an error while opening the Schedules database")
		fmt.Println(err)
		return schedule
	}

	_, err = db.Exec(SCHEDULE.CreationQuery)

	if err != nil {
		fmt.Println("Could not create the Schedule table.")
		fmt.Println(err)
		return schedule
	}

	query := `
	SELECT * FROM Schedules WHERE ScheduleID = ?
	`
	result, err := db.Query(query, id)

	if err != nil {
		fmt.Println("Error in the query for a schedule")
		fmt.Println(err)
	}

	for result.Next() {
		err := result.Scan(&schedule.ScheduleID, &schedule.AccountID, &schedule.ApprovalStatus, &schedule.Notes)
		if err != nil {
			fmt.Println(err)
			return schedule
		}
	}
	err = result.Err()
	if err != nil {
		fmt.Println(err)
		return schedule
	}

	fmt.Println("Successfully retreived the schedule.", schedule)
	return schedule
}

/*
Name: createSchedule
Return: none
Parameters: w http.ResponseWriter, r *http.Request
Purpose: to create a new schedule that automatically connects to the user.
It will not return anything because there will also be a createScheduleHandler that will manage
the html and the responses. The goal of this function is to just create it correctly in
the database.
*/

func createSchedule(w http.ResponseWriter, r *http.Request) {
	loggedIn := checkLoggedIn(w, r)
	if loggedIn != "LOGGEDIN" {
		fmt.Println("User is not logged in, cannot create a new schedule.")
		return
	}

	db, err := sql.Open("sqlite", SCHEDULE.FileURL)

	if err != nil {
		fmt.Println("Error in opening the table", err)
		return
	}
	_, err = db.Exec(SCHEDULE.CreationQuery)
	if err != nil {
		fmt.Println("Could not run creation correctly for the Schedule database", err)
		return
	}

	account := accountInfo(w, r)

	newSchedule := Schedule{
		AccountID:      account.ID,
		ApprovalStatus: "UNSUBMITTED",
	}

	query := `
	INSERT INTO Schedules (AccountID, ApprovalStatus) VALUES (?, ?);
	`

	result, err := db.Exec(query, newSchedule.AccountID, newSchedule.ApprovalStatus)
	if err != nil {
		fmt.Println("There was an error in creating the new schedule.", err)
		return
	}

	id, err := result.LastInsertId()

	if err != nil {
		fmt.Println("There was an error retrieving the new id", err)
		return
	}

	newSchedule.ScheduleID = int(id)
	defer db.Close()
	return
}

/*
Name: updateScheduleNote
Parameters:w, r and r should carry with it a form with the new information and the
id of the schedule.
Return:none
Purpose: This handler updates a note that a schedule has in order for an admin to make commentary
and allow the employee then adress the issue.
*/
func updateScheduleNote(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Updating the note")
	err := r.ParseForm()
	if err != nil {
		fmt.Println("There was an error with the form", err)
		accountPage(w, r)
		return
	}
	id, err := strconv.Atoi(r.FormValue("ScheduleID"))

	if err != nil {
		fmt.Println("There is an error with the value sent", err)
		accountPage(w, r)
		return
	}
	fmt.Println("Successfully found the id:", id)
	sched := getScheduleFromSchedID(id)

	note := r.FormValue("Note")

	if !checkAdmin(w, r) && accountInfo(w, r).ID != sched.AccountID {
		// w.Header().Set("HX-Redirect", "/account/")
		// w.WriteHeader(http.StatusUnauthorized)
		accountPage(w, r)
		fmt.Println("The user is authorized to see schedule")
		return
	}

	db, err := sql.Open("sqlite", SCHEDULE.FileURL)
	if err != nil {
		fmt.Println("There was an error opening the schedule database, could not update approval.", err)
		return
	}
	query := `
	UPDATE Schedules SET Notes = ? WHERE ScheduleID = ?;
	`

	res, err := db.Exec(query, note, sched.ScheduleID)
	if err != nil {
		fmt.Println("There was an error updating the schedule approval.", err)
		return
	}
	db.Close()

	fmt.Println("Successfully changed the Note.", res)

	w.Header().Set("HX-Redirect", "/admin/schedules/approvals/")
	w.WriteHeader(http.StatusOK)
}

/*
Name: getApprovalsList
Purpose: This is to get a list of all of the schedules in the database that need approval
*/
func getApprovalsList(w http.ResponseWriter, r *http.Request) []Schedule {
	var approvals []Schedule
	if !checkAdmin(w, r) {
		fmt.Println("User is not authorized to see approvals")
		return approvals
	}

	db, err := sql.Open("sqlite", SCHEDULE.FileURL)
	if err != nil {
		fmt.Println("In the first one")
		fmt.Println(err)
		return approvals
	}

	query := `
	SELECT * FROM Schedules;
	`

	result, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
		return approvals
	}
	defer result.Close()

	for result.Next() {
		var schedule Schedule
		err := result.Scan(&schedule.ScheduleID, &schedule.AccountID, &schedule.ApprovalStatus, &schedule.Notes)
		if err != nil {
			fmt.Println(err)
			return approvals
		}

		approvals = append(approvals, schedule)
	}

	err = result.Err()

	if err != nil {
		fmt.Println(err)
		return approvals
	}
	fmt.Print("These are all of the schedules that need approval", approvals)
	return approvals
}

/*
Name getScheduleStatus
Params: w, r
Return: none
Purpose: to find and send to the site the status of the user's schedule.
*/
func getApprovalStatus(w http.ResponseWriter, r *http.Request) {
	sched := getScheduleFromAcc(w, r)

	files := []string{
		HTML_FILES["APPROVALSTATUS"],
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		fmt.Println("There was an error parsing the approval files", err)
		return
	}
	ts.Execute(w, sched)
}

func getNote(w http.ResponseWriter, r *http.Request) {
	sched := getScheduleFromAcc(w, r)

	files := []string{
		HTML_FILES["NOTE"],
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		fmt.Println("There was an error parsing the approval files", err)
		return
	}
	fmt.Println("Successfully retrieved the note:", sched.Notes)
	ts.Execute(w, sched)
}

/*
Name: approveSchedule
Parameters: writer and request. THE REQUEST IS EXPECTED TO HAVE A VALUE OF "ScheduleID".
Return:none
Purpose: This function is to handle when an admin selects approve on the admin page and
then go in and update a schedule to show that it is approved.
*/
func approveSchedule(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	id, err := strconv.Atoi(r.FormValue("ScheduleID"))
	if err != nil {
		fmt.Println("Invalid scheduleId sent.", id, err)
	}
	schedule := getScheduleFromSchedID(id)

	updateApproval(schedule, "APPROVED", "/admin/schedules/approvals/", w, r)

}

func revisionSchedule(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	id, err := strconv.Atoi(r.FormValue("ScheduleID"))
	if err != nil {
		fmt.Println("Invalid scheduleId sent.", id, err)
	}
	schedule := getScheduleFromSchedID(id)

	updateApproval(schedule, "REVISION REQUIRED", "/admin/schedules/approvals/", w, r)

}

/*
Name: sendScheduleForReview
Parameters: w, r

*/

func sendScheduleForReview(w http.ResponseWriter, r *http.Request) {

	schedule := getScheduleFromAcc(w, r)

	updateApproval(schedule, "PENDING", "/account/", w, r)
}

func unsubmittedSchedule(w http.ResponseWriter, r *http.Request) {

	schedule := getScheduleFromAcc(w, r)

	updateApproval(schedule, "UNSUBMITTED", "/account/", w, r)
	fmt.Println(schedule)
}

/*
Name: updateApproval
Parameters: Schedule, newStatus string
Return:None
Purpose: This function takes a schedule and updates its approval status to the one
given in the parameter in the database.
*/
func updateApproval(schedule Schedule, newStatus string, redirect string, w http.ResponseWriter, r *http.Request) {
	if !checkAdmin(w, r) && newStatus != "PENDING" {
		fmt.Println("User not authorized to make changes.")
		return
	}
	db, err := sql.Open("sqlite", SCHEDULE.FileURL)
	if err != nil {
		fmt.Println("There was an error opening the schedule database, could not update approval.", err)
		return
	}
	query := `
	UPDATE Schedules SET ApprovalStatus = ? WHERE ScheduleID = ?
	`

	res, err := db.Exec(query, newStatus, schedule.ScheduleID)
	if err != nil {
		fmt.Println("There was an error updating the schedule approval.", err)
		return
	}
	db.Close()

	fmt.Println("Successfully changed the approval.", res)

	w.Header().Set("HX-Redirect", redirect)
	w.WriteHeader(http.StatusOK)
}

/*
Name: getWeeklyTotal
Parameters: w, r where r has to have a header with the ScheduleID
Return: returns to the web page the number of hours
*/
