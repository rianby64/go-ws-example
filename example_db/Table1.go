package example

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/lib/pq"
	grid "github.com/rianby64/arca-grid"
	arca "github.com/rianby64/arca-ws-jsonrpc"
)

// BindTable1WithPg whatever
func BindTable1WithPg(
	s *arca.JSONRPCExtensionWS,
	connStr string,
	db *sql.DB,
	dbName string,
	mirrors *[]*sql.DB,
) *grid.Grid {

	type Table1 struct {
		ID   int64
		Num1 float64
		Num2 float64
	}

	g := grid.Grid{}

	var queryHandler grid.RequestHandler = func(
		requestParams *interface{},
		context *interface{},
		notify grid.NotifyCallback,
	) (interface{}, error) {
		rows, err := db.Query(`
		SELECT
			"ID",
			"Num1",
			"Num2"
		FROM "Table1"
		ORDER BY "ID"
		`)
		if err != nil {
			log.Fatal(err)
		}

		var results []Table1

		var iID interface{}
		var iNum1 interface{}
		var iNum2 interface{}

		for rows.Next() {
			err := rows.Scan(
				&iID,
				&iNum1,
				&iNum2,
			)
			if err != nil {
				log.Fatal(err)
			}

			var ID int64
			var Num1 float64
			var Num2 float64

			if iID != nil {
				ID = iID.(int64)
			}
			if iNum1 != nil {
				Num1 = iNum1.(float64)
			}
			if iNum2 != nil {
				Num2 = iNum2.(float64)
			}

			results = append(results, Table1{
				ID:   ID,
				Num1: Num1,
				Num2: Num2,
			})
		}

		rows.Close()
		return results, nil
	}

	var updateHandler grid.RequestHandler = func(
		requestParams *interface{},
		context *interface{},
		notify grid.NotifyCallback,
	) (interface{}, error) {
		params := (*requestParams).(map[string]interface{})
		setters := []string{}
		for key, value := range params {
			if key == "ID" {
				continue
			}
			if key == "Num1" || key == "Num2" {
				Value := value.(float64)
				setters = append(setters, fmt.Sprintf(`"%v"=%v`, key, Value))
			}
		}
		strSetters := strings.Join(setters, ",")
		ID := params["ID"].(float64)
		query := fmt.Sprintf(`
		UPDATE "Table1"
			SET %v
			WHERE "ID"='%v';
		`, strSetters, ID)
		_, err := db.Exec(query)
		if err != nil {
			log.Println(err)
		}
		for _, mirror := range *mirrors {
			_, err := mirror.Exec(query)
			if err != nil {
				log.Println(err)
			}
		}
		return nil, nil
	}

	var insertHandler grid.RequestHandler = func(
		requestParams *interface{},
		context *interface{},
		notify grid.NotifyCallback,
	) (interface{}, error) {
		params := (*requestParams).(map[string]interface{})
		fields := []string{}
		values := []string{}
		for key, value := range params {
			if key == "ID" {
				continue
			}
			fields = append(fields, fmt.Sprintf(`"%v"`, key))
			if key == "Num1" || key == "Num2" {
				Value := value.(float64)
				values = append(values, fmt.Sprintf(`%v`, Value))
			}
		}
		strValues := strings.Join(values, ",")
		strFields := strings.Join(fields, ",")
		query := fmt.Sprintf(`
		INSERT INTO "Table1"(%v)
			VALUES(%v);
		`, strFields, strValues)
		_, err := db.Exec(query)
		if err != nil {
			log.Println(err)
		}
		for _, mirror := range *mirrors {
			_, err := mirror.Exec(query)
			if err != nil {
				log.Println(err)
			}
		}
		return nil, nil
	}

	var deleteHandler grid.RequestHandler = func(
		requestParams *interface{},
		context *interface{},
		notify grid.NotifyCallback,
	) (interface{}, error) {
		params := (*requestParams).(map[string]interface{})
		ID := params["ID"].(float64)

		query := fmt.Sprintf(`
		DELETE FROM "Table1"
			WHERE "ID"='%v';
		`, ID)
		_, err := db.Exec(query)
		if err != nil {
			log.Println(err)
		}
		for _, mirror := range *mirrors {
			_, err := mirror.Exec(query)
			if err != nil {
				log.Println(err)
			}
		}
		return nil, nil
	}

	methods := grid.QUID{
		Query:  &queryHandler,
		Update: &updateHandler,
		Insert: &insertHandler,
		Delete: &deleteHandler,
	}

	BindArcaWithGrid(connStr, s, &g, &methods, "Table1")

	/*
		Here I have to handle the way to redirect the notifications...
		Still I'm exploring
	*/

	type pgNotifyJSONRPC struct {
		Method string
		Source string
		Db     string
		Result interface{}
	}

	reportProblem := func(_ pq.ListenerEventType, err error) {
		if err != nil {
			log.Fatalln(err)
		}
	}

	minReconn := 10 * time.Second
	maxReconn := time.Minute
	listener := pq.NewListener(connStr, minReconn, maxReconn, reportProblem)
	err := listener.Listen("jsonrpc")
	if err != nil {
		panic(err)
	}

	go (func() {
		for {
			msg, ok := <-listener.Notify
			if !ok {
				return
			}
			var notification pgNotifyJSONRPC
			payload := []byte(msg.Extra)

			err := json.Unmarshal(payload, &notification)
			if err != nil {
				log.Println(err, ":: Notification ERROR")
			}

			var context interface{} = map[string]string{
				"Source": notification.Source,
				"Db":     notification.Db,
			}
			var response arca.JSONRPCresponse

			response.Method = notification.Method
			response.Context = context
			response.Result = notification.Result

			log.Println(dbName, response, ":: Notification")

			s.Broadcast(&response)
		}
	})()
	return &g
}
