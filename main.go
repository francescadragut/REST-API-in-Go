package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

const (
	host     = "localhost"
	port     = "5432"
	user     = "francesca"
	password = "1314"
	dbname   = "postgres"
)

func printMessage(message string) {
	fmt.Println("")
	fmt.Println(message)
	fmt.Println("")
}

func setupDB() *sql.DB {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)
	db, err := sql.Open("postgres", dbinfo)

	checkErr(err)

	return db
}

type Company struct {
	Name            string `json:"Name"`
	CUI             string `json:"CUI"`
	Year            string `json:"Year"`
	CAEN            string `json:"CAEN"`
	CAENDescription string `json:"CAEN Description"`
	Employees       string `json:"Employees"`
	Loss            string `json:"Net Loss"`
	Profit          string `json:"Net Profit"`
	Turnover        string `json:"Turnover"`
	Stocks          string `json:"Stocks"`
}

type JsonResponse struct {
	Type    string    `json:"type"`
	Data    []Company `json:"data"`
	Message string    `json:"message"`
}

func GetCompanies(w http.ResponseWriter, r *http.Request) {
	db := setupDB()

	printMessage("Getting companies...")

	rows, err := db.Query("SELECT * FROM romanian_companies_financial_data_serial")

	checkErr(err)

	var companies []Company

	for rows.Next() {
		var id int
		var name string
		var cui string
		var year string
		var caen string
		var caenDescription string
		var employees string
		var loss string
		var profit string
		var turnover string
		var stocks string

		err = rows.Scan(&id, &name, &cui, &year, &caen, &caenDescription, &employees, &loss, &profit, &turnover, &stocks)

		checkErr(err)

		companies = append(companies, Company{Name: name, CUI: cui, Year: year, CAEN: caen, CAENDescription: caenDescription, Employees: employees, Loss: loss, Profit: profit, Turnover: turnover, Stocks: stocks})
	}
	var response = JsonResponse{Type: "success", Data: companies}
	json.NewEncoder(w).Encode(response)
}

func GetCompanyByCUI(w http.ResponseWriter, r *http.Request) {
	db := setupDB()

	params := mux.Vars(r)
	cui := params["cui"]

	fmt.Println("Getting company with CUI: " + cui)

	rows, err := db.Query("SELECT * FROM romanian_companies_financial_data_serial WHERE cui = $1", cui)

	checkErr(err)

	var companies []Company

	var count = 0
	var response = JsonResponse{}
	for rows.Next() {
		var id int
		var name string
		var cuiCode string
		var year string
		var caen string
		var caenDescription string
		var employees string
		var loss string
		var profit string
		var turnover string
		var stocks string

		err = rows.Scan(&id, &name, &cuiCode, &year, &caen, &caenDescription, &employees, &loss, &profit, &turnover, &stocks)

		checkErr(err)

		if cuiCode == cui {
			companies = append(companies, Company{Name: name, CUI: cuiCode, Year: year, CAEN: caen, CAENDescription: caenDescription, Employees: employees, Loss: loss, Profit: profit, Turnover: turnover, Stocks: stocks})
			count += 1
		}
	}

	if count == 0 {
		response = JsonResponse{Type: "error", Message: "The company does not exist."}
		json.NewEncoder(w).Encode(response)
	} else {
		response = JsonResponse{Type: "success", Data: companies}
		json.NewEncoder(w).Encode(response)
	}
}

func GetCompanyByName(w http.ResponseWriter, r *http.Request) {
	db := setupDB()

	params := mux.Vars(r)
	name := params["name"]

	fmt.Println("Getting company with name: " + name)

	rows, err := db.Query("SELECT * FROM romanian_companies_financial_data_serial WHERE denumire = $1", name)

	checkErr(err)

	var companies []Company

	var count = 0
	var response = JsonResponse{}
	for rows.Next() {
		var id int
		var companyName string
		var cui string
		var year string
		var caen string
		var caenDescription string
		var employees string
		var loss string
		var profit string
		var turnover string
		var stocks string

		err = rows.Scan(&id, &companyName, &cui, &year, &caen, &caenDescription, &employees, &loss, &profit, &turnover, &stocks)

		checkErr(err)

		if companyName == name {
			companies = append(companies, Company{Name: companyName, CUI: cui, Year: year, CAEN: caen, CAENDescription: caenDescription, Employees: employees, Loss: loss, Profit: profit, Turnover: turnover, Stocks: stocks})
			count += 1
		}
	}

	if count == 0 {
		response = JsonResponse{Type: "error", Message: "The company does not exist."}
		json.NewEncoder(w).Encode(response)
	} else {
		response = JsonResponse{Type: "success", Data: companies}
		json.NewEncoder(w).Encode(response)
	}
}

func GetCompaniesByYear(w http.ResponseWriter, r *http.Request) {
	db := setupDB()

	params := mux.Vars(r)
	year := params["year"]

	fmt.Println("Getting companies opened in year: " + year)

	rows, err := db.Query("SELECT * FROM romanian_companies_financial_data_serial WHERE an = $1", year)

	checkErr(err)

	var companies []Company

	var count = 0
	var response = JsonResponse{}
	for rows.Next() {
		var id int
		var name string
		var cui string
		var companyYear string
		var caen string
		var caenDescription string
		var employees string
		var loss string
		var profit string
		var turnover string
		var stocks string

		err = rows.Scan(&id, &name, &cui, &companyYear, &caen, &caenDescription, &employees, &loss, &profit, &turnover, &stocks)

		checkErr(err)

		if companyYear == year {
			companies = append(companies, Company{Name: name, CUI: cui, Year: companyYear, CAEN: caen, CAENDescription: caenDescription, Employees: employees, Loss: loss, Profit: profit, Turnover: turnover, Stocks: stocks})
			count += 1
		}
	}

	if count == 0 {
		response = JsonResponse{Type: "error", Message: "No companies opened in year " + year + "."}
		json.NewEncoder(w).Encode(response)
	} else {
		response = JsonResponse{Type: "success", Data: companies}
		json.NewEncoder(w).Encode(response)
	}
}

func CreateCompany(w http.ResponseWriter, r *http.Request) {
	companyName := r.FormValue("name")
	companyCUI := r.FormValue("cui")
	companyYear := r.FormValue("year")
	companyCAEN := r.FormValue("caen")
	companyCAENDescription := r.FormValue("caenDescription")
	companyEmployees := r.FormValue("employees")
	companyLoss := r.FormValue("loss")
	companyProfit := r.FormValue("profit")
	companyTurnover := r.FormValue("turnover")
	companyStocks := r.FormValue("stocks")

	var response = JsonResponse{}

	if companyName == "" || companyYear == "" || companyCUI == "" || companyCAEN == "" || companyCAENDescription == "" || companyEmployees == "" || companyLoss == "" || companyProfit == "" || companyTurnover == "" || companyStocks == "" {
		response = JsonResponse{Type: "error", Message: "You are missing a parameter of the company."}
	} else {
		db := setupDB()

		printMessage("Inserting company into DB")

		fmt.Println("Inserting new company with name: " + companyName + " and year: " + companyYear)

		var lastInsertID int
		err := db.QueryRow("INSERT INTO romanian_companies_financial_data_serial(denumire, cui, an, caen, den_caen, numar_mediu_de_salariati, pierdere_neta, profit_net, cifra_de_afaceri_neta, stocuri) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) returning id;", companyName, companyCUI, companyYear, companyCAEN, companyCAENDescription, companyEmployees, companyLoss, companyProfit, companyTurnover, companyStocks).Scan(&lastInsertID)

		checkErr(err)

		response = JsonResponse{Type: "success", Message: "The company has been inserted successfully!"}
	}
	json.NewEncoder(w).Encode(response)
}

func DeleteCompanyByName(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	name := params["name"]

	var response = JsonResponse{}

	if name == "" {
		response = JsonResponse{Type: "error", Message: "You are missing name parameter."}
	} else {
		db := setupDB()

		printMessage("Deleting company " + name + " from DB")

		_, err := db.Exec("DELETE FROM romanian_companies_financial_data_serial where denumire = $1", name)

		checkErr(err)

		response = JsonResponse{Type: "success", Message: "The company has been deleted successfully!"}
	}

	json.NewEncoder(w).Encode(response)
}

func DeleteCompanyByCUI(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	cui := params["cui"]

	var response = JsonResponse{}

	if cui == "" {
		response = JsonResponse{Type: "error", Message: "You are missing name parameter."}
	} else {
		db := setupDB()

		printMessage("Deleting company with CUI " + cui + " from DB")

		_, err := db.Exec("DELETE FROM romanian_companies_financial_data_serial where cui = $1", cui)

		checkErr(err)

		response = JsonResponse{Type: "success", Message: "The company has been deleted successfully!"}
	}

	json.NewEncoder(w).Encode(response)
}

func PutByCUI(w http.ResponseWriter, r *http.Request) {

	companyName := r.FormValue("name")
	companyCUI := r.FormValue("cui")
	companyYear := r.FormValue("year")
	companyCAENCode := r.FormValue("caen_code")
	companyCAENDescription := r.FormValue("caen_description")
	companyEmployees := r.FormValue("employees")
	companyLoss := r.FormValue("loss")
	companyProfit := r.FormValue("profit")
	companyTurnover := r.FormValue("turnover")
	companyStocks := r.FormValue("stocks")

	var response = JsonResponse{}

	db := setupDB()

	if companyCUI != "" {
		_, check := db.Exec("SELECT * FROM romanian_companies_financial_data_serial WHERE cui = $1", companyCUI)
		checkErr(check)

		if check == nil {
			_, create := db.Exec("INSERT INTO romanian_companies_financial_data_serial(denumire, cui, an, caen, den_caen, numar_mediu_de_salariati, pierdere_neta, profit_net, cifra_de_afaceri_neta, stocuri) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)", companyName, companyCUI, companyYear, companyCAENCode, companyCAENDescription, companyEmployees, companyLoss, companyProfit, companyTurnover, companyStocks)
			checkErr(create)
		}

		if companyName != "" {
			_, err := db.Exec("UPDATE romanian_companies_financial_data_serial SET denumire = $1 WHERE cui = $2;", companyName, companyCUI)
			checkErr(err)
		}
		if companyYear != "" {
			_, err := db.Exec("UPDATE romanian_companies_financial_data_serial SET an = $1 WHERE cui = $2;", companyYear, companyCUI)
			checkErr(err)
		}
		if companyCAENCode != "" {
			_, err := db.Exec("UPDATE romanian_companies_financial_data_serial SET caen = $1 WHERE cui = $2;", companyCAENCode, companyCUI)
			checkErr(err)
		}
		if companyCAENDescription != "" {
			_, err := db.Exec("UPDATE romanian_companies_financial_data_serial SET den_caen = $1 WHERE cui = $2;", companyCAENDescription, companyCUI)
			checkErr(err)
		}
		if companyEmployees != "" {
			_, err := db.Exec("UPDATE romanian_companies_financial_data_serial SET numar_mediu_de_salariati = $1 WHERE cui = $2;", companyEmployees, companyCUI)
			checkErr(err)
		}
		if companyLoss != "" {
			_, err := db.Exec("UPDATE romanian_companies_financial_data_serial SET pierdere_neta = $1 WHERE cui = $2;", companyLoss, companyCUI)
			checkErr(err)
		}
		if companyProfit != "" {
			_, err := db.Exec("UPDATE romanian_companies_financial_data_serial SET profit_net = $1 WHERE cui = $2;", companyProfit, companyCUI)
			checkErr(err)
		}
		if companyTurnover != "" {
			_, err := db.Exec("UPDATE romanian_companies_financial_data_serial SET cifra_de_afaceri_neta = $1 WHERE cui = $2;", companyTurnover, companyCUI)
			checkErr(err)
		}
		if companyStocks != "" {
			_, err := db.Exec("UPDATE romanian_companies_financial_data_serial SET stocuri = $1 WHERE cui = $2;", companyStocks, companyCUI)
			checkErr(err)
		}

	}

	response = JsonResponse{Type: "success", Message: "The movie has been updated successfully!"}
	json.NewEncoder(w).Encode(response)
}

func PatchByCUI(w http.ResponseWriter, r *http.Request) {
	db := setupDB()

	params := mux.Vars(r)
	cui := params["cui"]

	var response = JsonResponse{}

	if cui == "" {
		response = JsonResponse{Type: "error", Message: "You are missing cui parameter."}
		json.NewEncoder(w).Encode(response)
	} else {
		rows, check := db.Exec("SELECT * FROM romanian_companies_financial_data_serial WHERE cui = $1", cui)
		checkErr(check)
		rowsAff, _ := rows.RowsAffected()

		if rowsAff == 0 {
			response = JsonResponse{Type: "error", Message: "The company with the given cui does not exist."}
		} else {
			companyName := r.FormValue("name")
			companyYear := r.FormValue("year")
			companyCAENCode := r.FormValue("caen_code")
			companyCAENDescription := r.FormValue("caen_description")
			companyEmployees := r.FormValue("employees")
			companyLoss := r.FormValue("loss")
			companyProfit := r.FormValue("profit")
			companyTurnover := r.FormValue("turnover")
			companyStocks := r.FormValue("stocks")

			if companyName != "" {
				_, err := db.Exec("UPDATE romanian_companies_financial_data_serial SET denumire = $1 WHERE cui = $2;", companyName, cui)
				checkErr(err)
			}
			if companyYear != "" {
				_, err := db.Exec("UPDATE romanian_companies_financial_data_serial SET an = $1 WHERE cui = $2;", companyYear, cui)
				checkErr(err)
			}
			if companyCAENCode != "" {
				_, err := db.Exec("UPDATE romanian_companies_financial_data_serial SET caen = $1 WHERE cui = $2;", companyCAENCode, cui)
				checkErr(err)
			}
			if companyCAENDescription != "" {
				_, err := db.Exec("UPDATE romanian_companies_financial_data_serial SET den_caen = $1 WHERE cui = $2;", companyCAENDescription, cui)
				checkErr(err)
			}
			if companyEmployees != "" {
				_, err := db.Exec("UPDATE romanian_companies_financial_data_serial SET numar_mediu_de_salariati = $1 WHERE cui = $2;", companyEmployees, cui)
				checkErr(err)
			}
			if companyLoss != "" {
				_, err := db.Exec("UPDATE romanian_companies_financial_data_serial SET pierdere_neta = $1 WHERE cui = $2;", companyLoss, cui)
				checkErr(err)
			}
			if companyProfit != "" {
				_, err := db.Exec("UPDATE romanian_companies_financial_data_serial SET profit_net = $1 WHERE cui = $2;", companyProfit, cui)
				checkErr(err)
			}
			if companyTurnover != "" {
				_, err := db.Exec("UPDATE romanian_companies_financial_data_serial SET cifra_de_afaceri_neta = $1 WHERE cui = $2;", companyTurnover, cui)
				checkErr(err)
			}
			if companyStocks != "" {
				_, err := db.Exec("UPDATE romanian_companies_financial_data_serial SET stocuri = $1 WHERE cui = $2;", companyStocks, cui)
				checkErr(err)
			}
		}

	}
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/companies/", GetCompanies).Methods("GET")
	router.HandleFunc("/companies/name/{name}", GetCompanyByName).Methods("GET")
	router.HandleFunc("/companies/cui/{cui}", GetCompanyByCUI).Methods("GET")
	router.HandleFunc("/companies/year/{year}", GetCompaniesByYear).Methods("GET")
	router.HandleFunc("/companies/", CreateCompany).Methods("POST")
	router.HandleFunc("/companies/name/{name}", DeleteCompanyByName).Methods("DELETE")
	router.HandleFunc("/companies/cui/{cui}", DeleteCompanyByCUI).Methods("DELETE")
	router.HandleFunc("/companies/", PutByCUI).Methods("PUT")
	router.HandleFunc("/companies/{cui}/", PatchByCUI).Methods("PATCH")

	fmt.Println("Server at 8080")
	log.Fatal(http.ListenAndServe("localhost:8080", router))
}
