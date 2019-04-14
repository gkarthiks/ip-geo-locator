package types

import "encoding/xml"

// IP the overall XML data received from geo api
type IP struct {
	XMLName xml.Name `xml:"ip"`

	Results Results `xml:"results>result"`
}

// Results the results tag data in XML from geo api
type Results struct {
	IP          string `xml:"ip" json:"ip"`
	Host        string `xml:"host" json:"host"`
	ISP         string `xml:"isp" json:"isp"`
	City        string `xml:"city" json:"city"`
	Countrycode string `xml:"countrycode" json:"countrycode"`
	Countryname string `xml:"countryname" json:"countryname"`
	Latitude    string `xml:"latitude" json:"latitude"`
	Longitude   string `xml:"longitude" json:"longitude"`
}

// FinalResponse is the final structure for json response as payload
type FinalResponse struct {
	Responses []Response
}

// Response is the data encapsulated for every IP data
type Response struct {
	IP          string `json:"ip"`
	Host        string `json:"host"`
	Isp         string `json:"isp"`
	City        string `json:"city"`
	Countrycode string `json:"countrycode"`
	Countryname string `json:"countryname"`
	Latitude    string `json:"latitude"`
	Longitude   string `json:"longitude"`
}

// ErrResponse used to send the error responses
type ErrResponse struct {
	Status       string `json:"status"`
	ErrorMessage string `json:"errorMessage"`
}
