package models

import "encoding/xml"

type Rk7NotifyEvent struct {
	XMLName   xml.Name `xml:"a"`
	RestCode  string   `xml:"RestCode,attr"`
	DateTime  string   `xml:"DateTime,attr"`
	Situation string   `xml:"Situation,attr"`
	SeqNumber string   `xml:"seqnumber,attr"`
	GUID      string   `xml:"guid,attr"`
	Name      string   `xml:"name,attr"`
	ShiftNum  string   `xml:"ShiftNum,attr"`
	ShiftDate string   `xml:"ShiftDate,attr"`

	// Вложенные элементы
	Station *Station `xml:"Station"`
	Server  *Server  `xml:"Server"`
	Waiter  *Waiter  `xml:"Waiter"`
	Order   *Order   `xml:"Order"`
	Item    *Item    `xml:"Item"`
}

// Station описывает элемент <Station>
type Station struct {
	ID      string `xml:"id,attr"`
	Code    string `xml:"code,attr"`
	Name    string `xml:"name,attr"`
	GUID    string `xml:"guid,attr"`
	NetName string `xml:"NetName,attr"`
}

// Server описывает элемент <Server>
type Server struct {
	ID      string `xml:"id,attr"`
	Code    string `xml:"code,attr"`
	Name    string `xml:"name,attr"`
	GUID    string `xml:"guid,attr"`
	NetName string `xml:"NetName,attr"`
}

// Waiter описывает <Waiter>
type Waiter struct {
	ID   string `xml:"id,attr"`
	Code string `xml:"code,attr"`
	Name string `xml:"name,attr"`
	Role *Role  `xml:"Role"`
}

// Role описывает <Role>
type Role struct {
	ID   string `xml:"id,attr"`
	Code string `xml:"code,attr"`
	Name string `xml:"name,attr"`
}

// Order описывает <Order>
type Order struct {
	Visit      string `xml:"visit,attr"`
	OrderIdent string `xml:"orderIdent,attr"`
	GUID       string `xml:"guid,attr"`
	OrderName  string `xml:"orderName,attr"`
	OrderSum   string `xml:"orderSum,attr"`
	OpenTime   string `xml:"openTime,attr"`
	Locked     string `xml:"locked,attr"`
}

// Item описывает <Item>
type Item struct {
	ClassName string `xml:"ClassName,attr"`
}
