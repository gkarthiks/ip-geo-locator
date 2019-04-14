package utils

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gkarthiks/ip-geo-locator/src/types"
	"github.com/rogpeppe/go-charset/charset"
	"github.com/sirupsen/logrus"
)

// GetXML reterives the XML data from geo ip api
func GetXML(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, fmt.Errorf("GET error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("Status error: %v", resp.StatusCode)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("Read body: %v", err)
	}
	return data, nil
}

// FormulateErrorMessage formulates the JSON error message
func FormulateErrorMessage(w http.ResponseWriter, err error) (js []byte) {
	logrus.Error(err)
	w.Header().Set("Content-Type", "application/json")
	errString := err.Error()
	var errorResponses types.ErrResponse
	errorResponses.Status = "Error"
	errorResponses.ErrorMessage = errString
	js, _ = json.Marshal(errorResponses)
	return js
}

// PopulateJSONResponse populates the JSON response structure
func PopulateJSONResponse(ipRequestedData types.IP) (res types.Response) {
	res.City = ipRequestedData.Results.City
	res.Countrycode = ipRequestedData.Results.Countrycode
	res.Countryname = ipRequestedData.Results.Countryname
	res.IP = ipRequestedData.Results.IP
	res.Isp = ipRequestedData.Results.ISP
	res.Host = ipRequestedData.Results.Host
	res.Latitude = ipRequestedData.Results.Latitude
	res.Longitude = ipRequestedData.Results.Longitude
	return res
}

// ConvertResponseToStructType converts the iso-8559-1 encoded xml into defined struct
func ConvertResponseToStructType(xmlBytes []byte) (ipStructData types.IP) {
	reader := bytes.NewReader(xmlBytes)
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReader
	err := decoder.Decode(&ipStructData)
	if err != nil {
		logrus.Error("decoder error:", err)
	}
	return ipStructData
}
