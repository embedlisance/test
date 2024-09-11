package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	_ "github.com/lib/pq" // PostgreSQL driver
)

func sendinttosql(val int, oid string, targetIp string) {

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	insertQuery := `INSERT INTO public."DeviceConfigs" (sensortype,oid,deviceip) VALUES ($1, $2, $3)`
	// Execute a query
	_, err = db.Exec(insertQuery, val, oid, targetIp)
	if err != nil {
		log.Fatal(err)
	}
}

// Define a struct to hold form data
type FormData struct {
	DeviceIP   string
	DeviceOID  string
	SensorType string
}

const connStr = "user=postgres password=123 dbname=birsens sslmode=disable"

// HTML template for the form
const formHTML = `
<!DOCTYPE html>
<html>
<head>
    <title>Device Form</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: #f4f4f4;
            margin: 0;
            padding: 0;
            display: flex;
            justify-content: center;
            align-items: center;
            height: 100vh;
        }
        .container {
            background: white;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
            width: 100%;
            max-width: 600px;
        }
        h1 {
            text-align: center;
            color: #333;
        }
        form {
            display: flex;
            flex-direction: column;
        }
        label {
            margin: 10px 0 5px;
            font-weight: bold;
            color: #555;
        }
        input[type="text"] {
            padding: 10px;
            border: 1px solid #ddd;
            border-radius: 4px;
            margin-bottom: 20px;
        }
        input[type="submit"] {
            background-color: #4CAF50;
            color: white;
            border: none;
            padding: 10px 15px;
            border-radius: 4px;
            cursor: pointer;
            font-size: 16px;
        }
        input[type="submit"]:hover {
            background-color: #45a049;
        }
        .result {
            margin-top: 20px;
            padding: 15px;
            background: #e7f1ff;
            border: 1px solid #b3d9ff;
            border-radius: 4px;
        }
        .result h2 {
            margin-top: 0;
            color: #333;
        }
        .result p {
            margin: 5px 0;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Device Form</h1>
        <form method="post" action="/">
            <label for="device-ip">Device IP:</label>
            <input type="text" id="device-ip" name="device_ip" required>

            <label for="device-oid">Device OID:</label>
            <input type="text" id="device-oid" name="device_oid" required>

            <label for="sensor-type">Sensor Type:</label>
            <input type="text" id="sensor-type" name="sensor_type" required>

            <input type="submit" value="Submit">
        </form>
        {{if .}}
        <div class="result">
			
            <h2>Submitted Data:</h2>
            <p>Device IP: {{.DeviceIP}}</p>
            <p>Device OID: {{.DeviceOID}}</p>
            <p>Sensor Type: {{.SensorType}}</p>
        </div>
        {{end}}
    </div>
</body>
</html>
`

// Handler function to display the form and process submissions
func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Parse form data
		deviceIP := r.FormValue("device_ip")
		deviceOID := r.FormValue("device_oid")
		sensorType := r.FormValue("sensor_type")
		failed := "INSERT FAILED"

		num, err := strconv.Atoi(sensorType)
		if err != nil {
			fmt.Println("Error:", err)
			// Create FormData object to pass to template
			data := FormData{
				DeviceIP:   failed,
				DeviceOID:  failed,
				SensorType: failed,
			}

			// Render the template with form data
			tmpl, err := template.New("form").Parse(formHTML)
			if err != nil {
				http.Error(w, "Error parsing template", http.StatusInternalServerError)
				return
			}
			tmpl.Execute(w, data)
			return
		}
		fmt.Println("Converted number:", num)

		// Create FormData object to pass to template
		data := FormData{
			DeviceIP:   deviceIP,
			DeviceOID:  deviceOID,
			SensorType: sensorType,
		}

		// Render the template with form data
		tmpl, err := template.New("form").Parse(formHTML)
		if err != nil {
			http.Error(w, "Error parsing template", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, data)
		sendinttosql(num, deviceOID, deviceIP)
		return
	}

	// Render the empty form
	tmpl, err := template.New("form").Parse(formHTML)
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Starting server on :8080")
	http.ListenAndServe(":8080", nil)
}
