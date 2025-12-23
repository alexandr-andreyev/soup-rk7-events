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
	Station   *Station   `xml:"Station"`
	Server    *Server    `xml:"Server"`
	Waiter    *Waiter    `xml:"Waiter"`
	Order     *Order     `xml:"Order"`
	Item      *Item      `xml:"Item"`
	Session   *Session   `xml:"Session"`
	ChangeLog *ChangeLog `xml:"ChangeLog"`
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
	Table      *Table `xml:"Table"`
	Kdsstate   string `xml:"kdsstate,attr"`
}

type Table struct {
	ID   string `xml:"id,attr"`
	Code string `xml:"code,attr"`
	Name string `xml:"name,attr"`
	GUID string `xml:"guid,attr"`
}

// Item описывает <Item>
type Item struct {
	ClassName string `xml:"ClassName,attr"`
}

type Session struct {
	Uni          string `xml:"uni,attr"`
	LineGUID     string `xml:"line_guid,attr"`
	State        string `xml:"state,attr"`
	SessionID    string `xml:"sessionID,attr"`
	IsDraft      string `xml:"isDraft,attr"`
	RemindTime   string `xml:"remindTime,attr"`
	StartService string `xml:"startService,attr"`
	Printed      string `xml:"printed,attr"`
	CookMins     string `xml:"cookMins,attr"`

	Dishes []Dish `xml:"Dish"`
}

type Dish struct {
	ID       string `xml:"id,attr"`
	Code     string `xml:"code,attr"`
	Name     string `xml:"name,attr"`
	GUID     string `xml:"guid,attr"`
	Uni      string `xml:"uni,attr"`
	LineGUID string `xml:"line_guid,attr"`
	State    string `xml:"state,attr"`
	Price    string `xml:"price,attr"`
	Amount   string `xml:"amount,attr"`
	Quantity string `xml:"quantity,attr"`
	SrcQty   string `xml:"srcQuantity,attr"`
	New      string `xml:"new,attr"`

	Modis []Modi `xml:"Modi"`
}

type Modi struct {
	ID       string `xml:"id,attr"`
	Code     string `xml:"code,attr"`
	Name     string `xml:"name,attr"`
	GUID     string `xml:"guid,attr"`
	Uni      string `xml:"uni,attr"`
	LineGUID string `xml:"line_guid,attr"`
	State    string `xml:"state,attr"`
	Count    string `xml:"count,attr"`
}

type ChangeLog struct {
	Dishes []ChangeDish `xml:"Dish"`
	Modis  []ChangeModi `xml:"Modi"`
}

type ChangeDish struct {
	ID       string `xml:"id,attr"`
	Name     string `xml:"name,attr"`
	OldValue string `xml:"oldvalue,attr"`
	NewValue string `xml:"newvalue,attr"`
}

type ChangeModi struct {
	ID       string `xml:"id,attr"`
	Name     string `xml:"name,attr"`
	OldValue string `xml:"oldvalue,attr"`
	NewValue string `xml:"newvalue,attr"`
}
