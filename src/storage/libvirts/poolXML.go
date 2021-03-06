package libvirts

import (
	"encoding/xml"
	"github.com/danieldin95/lightstar/src/libstar"
)

type PoolXML struct {
	XMLName    xml.Name      `xml:"pool" json:"-"`
	Type       string        `xml:"type,attr" json:"type"`
	Name       string        `xml:"name" json:"name"`
	UUID       string        `xml:"uuid" json:"uuid"`
	Source     SourceXML     `xml:"source" json:"source"`
	Capacity   CapacityXML   `xml:"capacity" json:"capacity"`
	Allocation AllocationXML `xml:"allocation" json:"allocation"`
	Available  AvailableXML  `xml:"available" json:"available"`
	Target     TargetXML     `xml:"target" json:"target"`
}

func (pol *PoolXML) Decode(xmlData string) error {
	if err := xml.Unmarshal([]byte(xmlData), pol); err != nil {
		libstar.Error("PoolXML.Decode %s", err)
		return err
	}
	return nil
}

func (pol *PoolXML) Encode() string {
	data, err := xml.Marshal(pol)
	if err != nil {
		libstar.Error("PoolXML.Encode %s", err)
		return ""
	}
	return string(data)
}

type AvailableXML struct {
	XMLName xml.Name `xml:"available" json:"-"`
	Unit    string   `xml:"unit,attr" json:"unit"`
	Value   string   `xml:",chardata" json:"value"`
}

type SourceXML struct {
	XMLName xml.Name  `xml:"source" json:"-"`
	Host    HostXML   `xml:"host" json:"host"`
	Dir     DirXML    `xml:"dir" json:"dir"`
	Format  FormatXML `xml:"format" json:"format"`
}

type HostXML struct {
	XMLName xml.Name `xml:"host" json:"-"`
	Name    string   `xml:"name,attr" json:"name"`
}

type DirXML struct {
	XMLName xml.Name `xml:"dir" json:"-"`
	Path    string   `xml:"path,attr" json:"path"`
}
