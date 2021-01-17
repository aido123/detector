package main

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/aido123/detector/pkg/containerservice/detector"
	"github.com/aido123/detector/pkg/containerservice/detector/detectorapi"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	log "github.com/sirupsen/logrus"
)

type DetectionLine struct {
	Status      string
	Priority    int
	Description string
}

var (
	counters      map[string]prometheus.Gauge
	ids           []string
	tenant        string
	subscription  string
	resourceGroup string
	cluster       string
	pollDelay     int
	apiTimeout    int
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	ids = strings.Split(os.Getenv("DETECTOR_IDS"), ",")
	subscription = os.Getenv("AZURE_SUBSCRIPTION_ID")
	resourceGroup = os.Getenv("RESOURCE_GROUP")
	cluster = os.Getenv("CLUSTER")
	pollDelay, _ = strconv.Atoi(os.Getenv("POLL_DELAY"))
	apiTimeout, _ = strconv.Atoi(os.Getenv("API_TIMEOUT"))

	//Create Prometheus gauge counter map for each detector id
	counters = make(map[string]prometheus.Gauge)
	for _, id := range ids {
		counters[id] = promauto.NewGauge(prometheus.GaugeOpts{
			Name: "detector_" + strings.Replace(id, "-", "_", -1),
			Help: "Detector metric " + "detector_" + strings.Replace(id, "-", "_", -1),
		})
	}
}

func main() {

	//TODO: panic if vars not set
	process()
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}

func process() {

	detectorClient := detector.NewDetectorClient(subscription)

	go func() {
		for {
			authorizer, err := auth.NewAuthorizerFromEnvironment()
			if err != nil {
				log.Error("Authorization Error: " + err.Error())
				time.Sleep(time.Duration(pollDelay) * time.Second)
			}
			detectorClient.Authorizer = authorizer

			for _, id := range ids {
				detectionLines := getDetectionLines(id, detectorClient)
				metricDetectionLines(id, detectionLines)
			}
			time.Sleep(time.Duration(pollDelay) * time.Second)
		}

	}()
}

func getDetectionLines(id string, detectorClient detectorapi.DetectorClientAPI) []DetectionLine {

	//TODO consider moving start and end time format to function so we can unit test output
	now := time.Now()
	startTime := now.Add(time.Duration(-(pollDelay/60)-1) * time.Minute)
	endTime := now.Add(time.Duration(-1) * time.Minute)
	startTimeFormat := strings.Replace(startTime.Format("2006-01-02 15:04"), " ", "%20", -1)
	endTimeFormat := strings.Replace(endTime.Format("2006-01-02 15:04"), " ", "%20", -1)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(apiTimeout))
	defer cancel()

	out, err := detectorClient.Get(ctx, resourceGroup, cluster, id, startTimeFormat, endTimeFormat)

	detectionLines := []DetectionLine{}

	if err != nil {
		log.Error("Error calling detector api with id " + id + ". " + err.Error())
		return detectionLines
	} else {

		//TODO check if this is actually valid out.Properties.Dataset[0].Table.Rows[0][0]
		for _, logindividual := range out.Properties.Dataset[0].Table.Rows {
			logConcat := strings.Join(logindividual, " ")
			if logindividual[0] == "Success" {
				log.WithFields(log.Fields{
					"detectid": id,
				}).Debug(logConcat)
				detectionLines = append(detectionLines, DetectionLine{logindividual[0], 1, logConcat})
			} else if logindividual[0] == "Info" {
				log.WithFields(log.Fields{
					"detectid": id,
				}).Info(logConcat)
				detectionLines = append(detectionLines, DetectionLine{logindividual[0], 2, logConcat})
			} else if logindividual[0] == "Warning" {
				log.WithFields(log.Fields{
					"detectid": id,
				}).Warn(logConcat)
				detectionLines = append(detectionLines, DetectionLine{logindividual[0], 3, logConcat})
			} else if logindividual[0] == "Critical" {
				log.WithFields(log.Fields{
					"detectid": id,
				}).Error(logConcat)
				detectionLines = append(detectionLines, DetectionLine{logindividual[0], 4, logConcat})
			} else {
				log.WithFields(log.Fields{
					"detectid": id,
				}).Debug(logConcat)
				detectionLines = append(detectionLines, DetectionLine{"Other", 0, logConcat})
			}
		}

	}
	return detectionLines
}

func metricDetectionLines(id string, detectionLines []DetectionLine) {
	criticality := 0
	for _, dl := range detectionLines {
		if dl.Priority > criticality {
			criticality = dl.Priority
		}
	}
	counters[id].Add(float64(criticality))
}
