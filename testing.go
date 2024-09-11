package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver

	g "github.com/gosnmp/gosnmp"
)

const connStr = "user=postgres password=123 dbname=birsens sslmode=disable"

func sendinttosql(val int, oid string, targetIp string) {

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	now := time.Now()

	// Format the date and time part

	insertQuery := `INSERT INTO public."sensorDatas" (value, timestamp,sensoroid,deviceip) VALUES ($1, $2, $3, $4)`
	// Execute a query
	_, err = db.Exec(insertQuery, val, now, oid, targetIp)
	if err != nil {
		log.Fatal(err)
	}
}
func main() {

	g.Default.Target = "192.168.9.113"
	fmt.Println("start")
	for {
		fmt.Println("this works every 0.5 minute!")

		err := g.Default.Connect()
		if err != nil {
			log.Fatalf("Connect() err: %v", err)
		}
		defer g.Default.Conn.Close()

		oids := []string{"1.3.6.1.2.1.1.1.0", "1.3.6.1.4.1.53389.4.1.0",
			"1.3.6.1.4.1.53389.4.2.0", "1.3.6.1.4.1.53389.4.3.0", "1.3.6.1.4.1.53389.4.4.0"}
		result, err2 := g.Default.Get(oids)
		if err2 != nil {
			log.Fatalf("Get() err: %v", err2)
		}

		for i, variable := range result.Variables {
			fmt.Printf("%d: oid: %s ", i, variable.Name)

			switch variable.Type {
			case g.OctetString:
				fmt.Printf("string: %s\n", string(variable.Value.([]byte)))
			default:

				intver := g.ToBigInt(variable.Value)

				fmt.Printf("number: %d\n", intver)
				intnopoint := *intver
				intval := intnopoint.Int64()
				integer := int(intval)
				oid := variable.Name
				if integer > 0 && integer < 999 {
					sendinttosql(integer, oid, g.Default.Target)
				}

			}
		}
		time.Sleep(time.Minute / 2)
	}
}
