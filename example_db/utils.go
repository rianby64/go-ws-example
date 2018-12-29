package example

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/lib/pq"
	grid "github.com/rianby64/arca-grid"
	arca "github.com/rianby64/arca-ws-jsonrpc"
)

// BindArcaWithGrid whatever
func BindArcaWithGrid(
	connStr string,
	s *arca.JSONRPCExtensionWS,
	g *grid.Grid,
	methods *grid.QUID,
	source string,
) {

	var queryMethod arca.JSONRequestHandler = g.Query
	var updateMethod arca.JSONRequestHandler = g.Update
	var insertMethod arca.JSONRequestHandler = g.Insert
	var deleteMethod arca.JSONRequestHandler = g.Delete

	s.RegisterMethod(source, "read", &queryMethod)
	s.RegisterMethod(source, "update", &updateMethod)
	s.RegisterMethod(source, "insert", &insertMethod)
	s.RegisterMethod(source, "delete", &deleteMethod)

	g.Register(methods)
}

// ConnectNotifyWithArca whatever
func ConnectNotifyWithArca(
	connStr string,
	dbName string, dbNamePrimary string,
	s *arca.JSONRPCExtensionWS,
	dbs *map[string]*sql.DB,
) {

	type pgNotifyJSONRPC struct {
		Method  string
		Source  string
		Db      string
		Primary bool
		View    bool
		Result  interface{}
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
				log.Println("Disconnected", dbName)
				return
			}
			var notification pgNotifyJSONRPC
			payload := []byte(msg.Extra)

			err := json.Unmarshal(payload, &notification)
			if err != nil {
				log.Fatalln(payload, "Cant unmarshal it")
			}

			var context interface{} = map[string]interface{}{
				"Source":  notification.Source,
				"Db":      notification.Db,
				"Primary": notification.Primary,
				"View":    notification.View,
			}
			var response arca.JSONRPCresponse

			response.Method = notification.Method
			response.Context = context
			response.Result = notification.Result

			request := arca.JSONRPCrequest{}
			request.Method = notification.Method
			request.Context = map[string]interface{}{
				"Source": notification.Source,
				"Db":     notification.Db,
			}
			request.Params = notification.Result

			if notification.Primary && notification.Db == dbNamePrimary {
				for dbNameContext := range *dbs {
					if dbNameContext != dbNamePrimary {
						request.Context.(map[string]interface{})["Db"] = dbNameContext
						s.ProcessRequest(&request)
					}
				}
			}

			if notification.View {
				request.Context.(map[string]interface{})["Db"] = dbNamePrimary
				s.ProcessRequest(&request)
				continue
			}
			s.Broadcast(&response)

		}
	})()
}
