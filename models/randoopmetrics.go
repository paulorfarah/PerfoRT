package models

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"database/sql"

	"github.com/wcharczuk/go-chart/v2"
	"github.com/wcharczuk/go-chart/v2/drawing"
	"gorm.io/gorm"
)

type RandoopMetrics struct {
	Model
	CommitID uint
	Commit   Commit
	// ChangeID   uint
	// Change     Change
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
	NMEDiff    int
	EMEDiff    int
	AETNDiff   float64
	AETEDiff   float64
	AMUDiff    float64
	NMEPerc    float64
	EMEPerc    float64
	AETNPerc   float64
	AETEPerc   float64
	AMUPerc    float64
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

	rows, err := db.Query("SELECT rm.nme_before, rm.nme_after, rm.eme_before, rm.eme_after, rm.aetn_before, rm.aetn_after, rm.aete_before, rm.aete_after, rm.amu_before, rm.amu_after FROM randoopmetrics as rm INNER JOIN changes as c ON rm.change_id=c.id INNER JOIN commits as com ON c.commit_id = com.id ORDER BY committer_date;")
	defer rows.Close()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	var nmeBefore, nmeAfter, emeBefore, emeAfter []float64
	// var nmeAfter []float64
	var xValues []float64
	var count float64
	for rows.Next() {
		//read values
		var nmeBs, nmeAs, emeBs, emeAs, aetnBs, aetnAs, aeteBs, aeteAs, amuBs, amuAs string
		err = rows.Scan(&nmeBs, &nmeAs, &emeBs, &emeAs, &aetnBs, &aetnAs, &aeteBs, &aeteAs, &amuBs, &amuAs)
		if err != nil {
			fmt.Println(err.Error())
		}

		// x axis
		xValues = append(xValues, count)
		count++

		// nme
		nmeB, err := strconv.ParseFloat(nmeBs, 64)
		if err != nil {
			nmeB = 0
		}
		nmeBefore = append(nmeBefore, nmeB)
		nmeA, err := strconv.ParseFloat(nmeAs, 64)
		if err != nil {
			nmeA = 0
		}
		nmeAfter = append(nmeAfter, nmeA)

		// eme
		emeB, err := strconv.ParseFloat(emeBs, 64)
		if err != nil {
			emeB = 0
		}
		emeBefore = append(nmeBefore, emeB)
		emeA, err := strconv.ParseFloat(emeAs, 64)
		if err != nil {
			emeA = 0
		}
		emeAfter = append(nmeAfter, emeA)

	}

	PlotRandoopResults("nme", xValues, nmeBefore, nmeAfter)
	PlotRandoopResults("eme", xValues, emeBefore, emeAfter)
	// plot([]float64{1.0, 2.0, 3.0, 4.0}, []float64{0.0581, 0.0581, 0.0581, 0.0581})
}

func PlotRandoopResults(filename string, xValues, yValuesBef, yValuesAft []float64) {
	graph := chart.Chart{
		YAxis: chart.YAxis{
			Range: &chart.ContinuousRange{
				Min: 0.0,
				Max: 1,
			},
		},
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

	f, _ := os.Create(filename + ".png")
	defer f.Close()
	graph.Render(chart.PNG, f)

}

// func plot(xValues, yValuesBef []float64) {
// 	graph := chart.Chart{
// 		YAxis: chart.YAxis{
// 			Range: &chart.ContinuousRange{
// 				Min: 0.0,
// 				Max: 4,
// 			},
// 		},
// 		Series: []chart.Series{
// 			chart.ContinuousSeries{
// 				XValues: xValues,    //[]float64{1.0, 2.0, 3.0, 4.0},
// 				YValues: yValuesBef, //[]float64{1.0, 2.0, 3.0, 4.0},
// 			},
// 		},
// 	}

// 	// buffer := bytes.NewBuffer([]byte{})
// 	// err := graph.Render(chart.PNG, buffer)
// 	// if err != nil {
// 	// 	fmt.Println(err.Error())
// 	// }

// 	f, _ := os.Create("output2.png")
// 	defer f.Close()
// 	graph.Render(chart.PNG, f)
// }
