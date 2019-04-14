package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gkarthiks/ip-geo-locator/src/types"
	"github.com/gkarthiks/ip-geo-locator/src/utils"
	"github.com/gorilla/mux"
	_ "github.com/rogpeppe/go-charset/data"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.Info("Server using logrus for logging and mux for routing")
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/ip/{key}", GetGeoIPDetail)
	srv := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	logrus.Infof("Server started listening on " + srv.Addr)
	logrus.Fatal(srv.ListenAndServe())
}

// GetGeoIPDetail returns the geo details of the IP
func GetGeoIPDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if strings.Contains(vars["key"], ",") {
		var finalResponse types.FinalResponse
		var jData, xmlBytes []byte
		var err error
		slicedIPAddress := strings.Split(vars["key"], ",")
		for _, ipAddress := range slicedIPAddress {
			logrus.Info(strings.TrimSpace(ipAddress))
			xmlBytes, err = getXMLResponseFromGeoAPI(ipAddress)
			if err != nil {
				jData = utils.FormulateErrorMessage(w, err)
			} else {
				finalResponseBundle(xmlBytes, &finalResponse)
			}
		}
		sendResponseOrError(err, jData, w, finalResponse)
	} else {
		logrus.Info("Requesting data for single IP address")
		var jData, xmlBytes []byte
		var finalResponse types.FinalResponse
		var err error
		vars := mux.Vars(r)
		xmlBytes, err = getXMLResponseFromGeoAPI(strings.TrimSpace(vars["key"]))
		if err != nil {
			jData = utils.FormulateErrorMessage(w, err)
		} else {
			finalResponseBundle(xmlBytes, &finalResponse)
		}
		sendResponseOrError(err, jData, w, finalResponse)
	}
}

func getXMLResponseFromGeoAPI(ipAddress string) (xmlBytes []byte, err error) {
	if xmlBytes, err = utils.GetXML("http://api.geoiplookup.net/?query=" + ipAddress); err != nil {
		return nil, err
	}
	return xmlBytes, nil
}

// finalResponseBundle is responsible for preparing the bundle of final response calls to the utility functions
func finalResponseBundle(xmlBytes []byte, finalResponse *types.FinalResponse) {
	var ipXMLData types.IP
	ipXMLData = utils.ConvertResponseToStructType(xmlBytes)
	res := utils.PopulateJSONResponse(ipXMLData)
	finalResponse.Responses = append(finalResponse.Responses, res)
}

// sendResponseOrError sends back the successful response message or error message accordingly
func sendResponseOrError(err error, jData []byte, w http.ResponseWriter, finalResponse types.FinalResponse) {
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write(jData)
	} else {
		jData, err := json.Marshal(finalResponse)
		if err != nil {
			logrus.Error(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jData)
	}
}
