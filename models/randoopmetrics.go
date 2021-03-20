package models

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"database/sql"

	"github.com/jinzhu/gorm"
	"github.com/wcharczuk/go-chart/v2"
	"github.com/wcharczuk/go-chart/v2/drawing"
)

type RandoopMetrics struct {
	Model
	ChangeID   uint
	Change     Change
	NMEBefore  string
	EMEBefore  string
	AETNBefore string
	AETEBefore string
	AMUBefore  string
	NMEAfter   string
	EMEAfter   string
	AETNAfter  string
	AETEAfter  string
	AMUAfter   string
	NMEDiff    string
	EMEDiff    string
	AETNDiff   string
	AETEDiff   string
	AMUDiff    string
	NMEPerc    string
	EMEPerc    string
	AETNPerc   string
	AETEPerc   string
	AMUPerc    string
}

func (p *RandoopMetrics) TableName() string {
	return "randoopmetrics"
}

func CreateRandoopMetrics(db *gorm.DB, rm *RandoopMetrics) (uint, error) {
	err := db.Create(rm).Error
	if err != nil {
		return 0, err
	}
	fmt.Println("New randoop metrics added: ")
	return rm.ID, nil
}

type randooResult struct {
	Date      time.Time
	NMEBefore float64
	NMEAfter  float64
}

func GetRandoopMetrics() {
	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")
	dbPort := os.Getenv("db_port")
	db, err := sql.Open("mysql", username+":"+password+"@tcp("+dbHost+":"+dbPort+")/"+dbName)

	// if there is an error opening the connection, handle it
	if err != nil {
		panic(err.Error())
	}

	// defer the close till after the main function has finished
	// executing
	defer db.Close()

	rows, err := db.Query("SELECT rm.aetn_before, rm.aetn_after FROM randoopmetrics as rm INNER JOIN changes as c ON rm.change_id=c.id INNER JOIN commits as com ON c.commit_id = com.id ORDER BY committer_date;")
	defer rows.Close()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	var nmeBefore []float64
	var nmeAfter []float64
	var yValues []float64
	var count float64
	for rows.Next() {
		var i string
		var j string
		err = rows.Scan(&i, &j)
		if err != nil {
			fmt.Println(err.Error()) // proper error handling instead of panic in your app
		}
		fi, err := strconv.ParseFloat(i, 64)
		if err != nil {
			fi = float64(0.0)
		}
		nmeBefore = append(nmeBefore, fi)
		fj, err := strconv.ParseFloat(j, 64)
		if err != nil {
			fj = float64(0.0)
		}
		nmeAfter = append(nmeAfter, fj)
		yValues = append(yValues, count)
		count++

	}
	fmt.Println(nmeBefore)
	fmt.Println(nmeAfter)
	fmt.Println(yValues)
	// PlotRandoopResults(nmeBefore, nmeAfter, yValues)
	plot([]float64{1.0, 2.0, 3.0, 4.0}, []float64{0.0581, 0.0581, 0.0581, 0.0581})
}

func plot(yValuesBef, xValues []float64) {
	graph := chart.Chart{
		YAxis: chart.YAxis{
			Range: &chart.ContinuousRange{
				Min: 0.0,
				Max: 0.1,
			},
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				XValues: xValues,    //[]float64{1.0, 2.0, 3.0, 4.0},
				YValues: yValuesBef, //[]float64{1.0, 2.0, 3.0, 4.0},
			},
		},
	}

	// buffer := bytes.NewBuffer([]byte{})
	// err := graph.Render(chart.PNG, buffer)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	f, _ := os.Create("output2.png")
	defer f.Close()
	graph.Render(chart.PNG, f)
}

func PlotRandoopResults(yValuesBef, yValuesAft, xValues []float64) {
	graph := chart.Chart{
		XAxis: chart.XAxis{
			TickPosition: chart.TickPositionBetweenTicks,
			ValueFormatter: func(v interface{}) string {
				typed := v.(float64)
				typedDate := chart.TimeFromFloat64(typed)
				return fmt.Sprintf("%d-%d\n%d", typedDate.Month(), typedDate.Day(), typedDate.Year())
			},
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				Style: chart.Style{
					StrokeColor: drawing.ColorRed,               // will supercede defaults
					FillColor:   drawing.ColorRed.WithAlpha(64), // will supercede defaults
				},
				XValues: xValues,    //[]float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: yValuesBef, //[]float64{1.0, 2.0, 1.5, 4.0, 2.7},
			},
			chart.ContinuousSeries{
				YAxis:   chart.YAxisSecondary,
				XValues: xValues,    //[]float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: yValuesAft, //[]float64{50.0, 40.0, 30.0, 20.0, 10.0},
			},
		},
	}

	f, _ := os.Create("output.png")
	defer f.Close()
	graph.Render(chart.PNG, f)

}
