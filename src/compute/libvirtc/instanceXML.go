package libvirtc

import (
	"encoding/xml"
	"github.com/danieldin95/lightstar/src/libstar"
)

var (
	PciDomain       = "0x00"
	PciDiskBus      = "0x01"
	PciInterfaceBus = "0x02"
	PciFunc         = "0x00"
)

type DomainXML struct {
	XMLName  xml.Name    `xml:"domain" json:"-"`
	Id       string      `xml:"id,attr" json:"id"`
	Type     string      `xml:"type,attr" json:"type"` // kvm
	Name     string      `xml:"name" json:"name"`
	UUID     string      `xml:"uuid" json:"uuid"`
	OS       OSXML       `xml:"os" json:"os"`
	Features FeaturesXML `xml:"features" json:"features"`
	CPU      CPUXML      `xml:"cpu" json:"cpu"`
	VCPU     VCPUXML     `xml:"vcpu" json:"vcpu"`
	Memory   MemXML      `xml:"memory" json:"memory"`
	CurMem   CurMemXML   `xml:"currentMemory" json:"currentMemory"`
	Devices  DevicesXML  `xml:"devices" json:"devices"`
}

type CPUXML struct {
	XMLName xml.Name `xml:"cpu" json:"-"`
	Mode    string   `xml:"mode,attr,omitempty" json:"mode"`   // host-model, host-passthrough
	Check   string   `xml:"check,attr,omitempty" json:"check"` // partial, full
	Match   string   `xml:"match,attr,omitempty" json:"match"` // exact
}

func NewDomainXMLFromDom(dom *Domain, secure bool) *DomainXML {
	if dom == nil {
		return nil
	}
	var err error
	var xmlData string
	xmlData, err = dom.GetXMLDesc(secure)
	if err != nil {
		return nil
	}
	instXml := &DomainXML{}
	if err := instXml.Decode(xmlData); err != nil {
		return nil
	}
	return instXml
}

func (domain *DomainXML) Decode(xmlData string) error {
	if err := xml.Unmarshal([]byte(xmlData), domain); err != nil {
		libstar.Error("DomainXML.Decode %s", err)
		return err
	}
	return nil
}

func (domain *DomainXML) Encode() string {
	data, err := xml.Marshal(domain)
	if err != nil {
		libstar.Error("DomainXML.Encode %s", err)
		return ""
	}
	return string(data)
}

func (domain *DomainXML) GraphicsAddr(format string) (string, string) {
	if len(domain.Devices.Graphics) == 0 {
		return "", ""
	}
	for _, g := range domain.Devices.Graphics {
		if g.Type == format {
			return g.Listen, g.Port
		}
	}
	return "", ""
}

type VCPUXML struct {
	XMLName   xml.Name `xml:"vcpu" json:"-"`
	Placement string   `xml:"placement,attr" json:"placement"` // static
	Value     string   `xml:",chardata" json:"Value"`
}

func (cpu *VCPUXML) Decode(xmlData string) error {
	if err := xml.Unmarshal([]byte(xmlData), cpu); err != nil {
		libstar.Error("VCPUXML.Decode %s", err)
		return err
	}
	return nil
}

type FeaturesXML struct {
	XMLName xml.Name `xml:"features" json:"-"`
	Acpi    *AcpiXML `xml:"acpi,omitempty" json:"acpi"`
	Apic    *ApicXML `xml:"apic,omitempty" json:"apic"`
	Pae     *PaeXML  `xml:"pae,omitempty" json:"pae"`
}

func (feat *FeaturesXML) Decode(xmlData string) error {
	if err := xml.Unmarshal([]byte(xmlData), feat); err != nil {
		libstar.Error("FeaturesXML.Decode %s", err)
		return err
	}
	return nil
}

type AcpiXML struct {
	XMLName xml.Name `xml:"acpi" json:"-"`
	Value   string   `xml:",chardata" json:"value"`
}

type ApicXML struct {
	XMLName xml.Name `xml:"apic" json:"-"`
	Value   string   `xml:",chardata" json:"value"`
}

type PaeXML struct {
	XMLName xml.Name `xml:"pae" json:"-"`
	Value   string   `xml:",chardata" json:"value"`
}

type MemXML struct {
	XMLName xml.Name `xml:"memory" json:"-"`
	Type    string   `xml:"unit,attr" json:"unit"`
	Value   string   `xml:",chardata" json:"value"`
}

func (mem *MemXML) Decode(xmlData string) error {
	if err := xml.Unmarshal([]byte(xmlData), mem); err != nil {
		libstar.Error("MemXML.Decode %s", err)
		return err
	}
	return nil
}

type CurMemXML struct {
	XMLName xml.Name `xml:"currentMemory" json:"-"`
	Type    string   `xml:"unit,attr" json:"unit"`
	Value   string   `xml:",chardata" json:"value"`
}

func (cmem *CurMemXML) Decode(xmlData string) error {
	if err := xml.Unmarshal([]byte(xmlData), cmem); err != nil {
		libstar.Error("CurMemXML.Decode %s", err)
		return err
	}
	return nil
}

type OSXML struct {
	XMLName  xml.Name      `xml:"os" json:"-"`
	Type     OSTypeXML     `xml:"type" json:"type"`
	Boot     []OSBootXML   `xml:"boot" json:"boot"` //<bootmenu enable='yes'/>
	BootMenu OSBootMenuXML `xml:"bootmenu" json:"bootmenu"`
}

func (osx *OSXML) Decode(xmlData string) error {
	if err := xml.Unmarshal([]byte(xmlData), osx); err != nil {
		libstar.Error("OSXML.Decode %s", err)
		return err
	}
	return nil
}

type OSTypeXML struct {
	XMLName xml.Name `xml:"type" json:"-"`
	Arch    string   `xml:"arch,attr" json:"arch"` // x86_64
	Machine string   `xml:"machine,attr" json:"machine"`
	Value   string   `xml:",chardata" json:"value"` // hvm
}

type OSBootXML struct {
	XMLName xml.Name `xml:"boot" json:"-"`
	Dev     string   `xml:"dev,attr" json:"dev"` // hd, cdrom, network
}

type OSBootMenuXML struct {
	XMLName xml.Name `xml:"bootmenu" json:"-"`
	Enable  string   `xml:"enable,attr" json:"enable"` // yes
}

type DevicesXML struct {
	XMLName     xml.Name           `xml:"devices" json:"-"`
	Graphics    []GraphicsXML      `xml:"graphics" json:"graphics"`
	Disks       []DiskXML          `xml:"disk" json:"disk"`
	Interfaces  []InterfaceXML     `xml:"interface" json:"interface"`
	Controllers []ControllerXML    `xml:"controller" json:"controller"`
	Inputs      []InputDeviceXML   `xml:"input" json:"input"`
	Sound       SoundDeviceXML     `xml:"sound" json:"sound"`
	Video       VideoDeviceXML     `xml:"video" json:"video"`
	Channels    []ChannelDeviceXML `xml:"channel" json:"channel"`
}

func (devices *DevicesXML) Decode(xmlData string) error {
	if err := xml.Unmarshal([]byte(xmlData), devices); err != nil {
		libstar.Error("DevicesXML.Decode %s", err)
		return err
	}
	return nil
}

// Bus 0: root
// Bus 1: disk, ide
// Bus 2: interface
// Bus 3: reverse
type ControllerXML struct {
	XMLName xml.Name    `xml:"controller" json:"-"`
	Type    string      `xml:"type,attr" json:"type"`
	Index   string      `xml:"index,attr" json:"port"`
	Model   string      `xml:"model,attr,omitempty" json:"model"` // pci-root, pci-bridge.
	Address *AddressXML `xml:"address,omitempty" json:"address"`
}

func (ctl *ControllerXML) Decode(xmlData string) error {
	if err := xml.Unmarshal([]byte(xmlData), ctl); err != nil {
		libstar.Error("ControllerXML.Decode %s", err)
		return err
	}
	return nil
}

type AddressXML struct {
	XMLName    xml.Name `xml:"address" json:"-"`
	Type       string   `xml:"type,attr,omitempty" json:"type"` // pci and drive.
	Domain     string   `xml:"domain,attr,omitempty" json:"domain"`
	Bus        string   `xml:"bus,attr,omitempty" json:"bus"`
	Slot       string   `xml:"slot,attr,omitempty" json:"slot"`
	Function   string   `xml:"function,attr,omitempty" json:"function"`
	Target     string   `xml:"target,attr,omitempty" json:"target"`
	Unit       string   `xml:"unit,attr,omitempty" json:"unit"`
	Controller string   `xml:"controller,attr,omitempty" json:"controller"`
}

type GraphicsXML struct {
	XMLName  xml.Name `xml:"graphics" json:"-"`
	Type     string   `xml:"type,attr" json:"type"` // vnc, spice
	Port     string   `xml:"port,attr" json:"port"`
	AutoPort string   `xml:"autoport,attr,omitempty" json:"autoport"`
	Listen   string   `xml:"listen,attr" json:"listen"`
	Password string   `xml:"passwd,attr,omitempty" json:"password"`
}

func (graphics *GraphicsXML) Decode(xmlData string) error {
	if err := xml.Unmarshal([]byte(xmlData), graphics); err != nil {
		libstar.Error("GraphicsXML.Decode %s", err)
		return err
	}
	return nil
}

func (graphics *GraphicsXML) Encode() string {
	data, err := xml.Marshal(graphics)
	if err != nil {
		libstar.Error("GraphicsXML.Encode %s", err)
		return ""
	}
	return string(data)
}

type DiskXML struct {
	XMLName xml.Name      `xml:"disk" json:"-"`
	Type    string        `xml:"type,attr" json:"type"`
	Device  string        `xml:"device,attr" json:"device"`
	Driver  DiskDriverXML `xml:"driver" json:"driver"`
	Source  DiskSourceXML `xml:"source" json:"source"`
	Target  DiskTargetXML `xml:"target" json:"target"`
	Address *AddressXML   `xml:"address,omitempty" json:"address"`
}

func (disk *DiskXML) Decode(xmlData string) error {
	if err := xml.Unmarshal([]byte(xmlData), disk); err != nil {
		libstar.Error("DiskXML.Decode %s", err)
		return err
	}
	return nil
}

func (disk *DiskXML) Encode() string {
	data, err := xml.Marshal(disk)
	if err != nil {
		libstar.Error("DiskXML.Encode %s", err)
		return ""
	}
	return string(data)
}

type DiskDriverXML struct {
	XMLName xml.Name `xml:"driver" json:"-"`
	Type    string   `xml:"type,attr" json:"type"`
	Name    string   `xml:"name,attr" json:"name"`
}

type DiskSourceXML struct {
	XMLName xml.Name `xml:"source" json:"-"`
	File    string   `xml:"file,attr,omitempty" json:"file"`
	Device  string   `xml:"dev,attr,omitempty" json:"device"`
}

type DiskTargetXML struct {
	XMLName xml.Name `xml:"target" json:"-"`
	Bus     string   `xml:"bus,attr,omitempty" json:"bus"`
	Dev     string   `xml:"dev,attr,omitempty" json:"dev"`
}

type InterfaceXML struct {
	XMLName     xml.Name             `xml:"interface" json:"-"`
	Type        string               `xml:"type,attr" json:"type"`
	Mac         InterfaceMacXML      `xml:"mac" json:"mac"`
	Source      InterfaceSourceXML   `xml:"source" json:"source"`
	Model       InterfaceModelXML    `xml:"model" json:"model"`
	Driver      *InterfaceDriverXML  `xml:"driver" json:"driver"`
	Target      InterfaceTargetXML   `xml:"target" json:"target"`
	VirtualPort *InterfaceVirPortXML `xml:"virtualport,omitempty" json:"virtualport"`
	Address     *AddressXML          `xml:"address,omitempty" json:"address"`
}

func (int *InterfaceXML) Decode(xmlData string) error {
	if err := xml.Unmarshal([]byte(xmlData), int); err != nil {
		libstar.Error("InterfaceXML.Decode %s", err)
		return err
	}
	return nil
}

func (int *InterfaceXML) Encode() string {
	data, err := xml.Marshal(int)
	if err != nil {
		libstar.Error("InterfaceXML.Encode %s", err)
		return ""
	}
	return string(data)
}

type InterfaceMacXML struct {
	XMLName xml.Name `xml:"mac" json:"-"`
	Address string   `xml:"address,attr,omitempty" json:"address"`
}

type InterfaceSourceXML struct {
	XMLName xml.Name    `xml:"source" json:"-"`
	Bridge  string      `xml:"bridge,attr,omitempty" json:"bridge"`
	Network string      `xml:"network,attr,omitempty" json:"network"`
	Address *AddressXML `xml:"address,omitempty" json:"address"`
}

type InterfaceModelXML struct {
	XMLName xml.Name `xml:"model" json:"-"`
	Type    string   `xml:"type,attr" json:"type"` //rtl8139, virtio, e1000
}

type InterfaceDriverXML struct {
	XMLName xml.Name `xml:"driver" json:"-"`
	Name    string   `xml:"name,attr" json:"name"` //vfio
}

type InterfaceTargetXML struct {
	XMLName xml.Name `xml:"target" json:"-"`
	Bus     string   `xml:"bus,attr,omitempty" json:"bus"`
	Dev     string   `xml:"dev,attr,omitempty" json:"dev"`
}

type InterfaceVirPortXML struct {
	XMLName xml.Name `xml:"virtualport" json:"-"`
	Type    string   `xml:"type,attr,omitempty" json:"type"` //openvswitch
}

type InputDeviceXML struct {
	XMLName xml.Name `xml:"input" json:"-"`
	Type    string   `xml:"type,attr" json:"type"`
	Bus     string   `xml:"bus,attr" json:"bus"`
}

type SoundDeviceXML struct {
	XMLName xml.Name `xml:"sound" json:"-"`
	Model   string   `xml:"model,attr" json:"model"` // ac97, ich6
}

type VideoDeviceXML struct {
	XMLName xml.Name      `xml:"video" json:"-"`
	Model   VideoModelXML `xml:"model" json:"model"`
}

type VideoModelXML struct {
	XMLName xml.Name `xml:"model" json:"-"`
	Type    string   `xml:"type,attr" json:"type"` // qxl, cirrus
}

type ChannelDeviceXML struct {
	XMLName xml.Name         `xml:"channel" json:"-"`
	Type    string           `xml:"type,attr" json:"type"`
	Source  ChannelSourceXML `xml:"source" json:"source"`
	Target  ChannelTargetXML `xml:"target" json:"target"`
}

type ChannelTargetXML struct {
	XMLName xml.Name `xml:"target" json:"-"`
	Type    string   `xml:"type,attr" json:"type"`
	Name    string   `xml:"name,attr" json:"name"`
}

type ChannelSourceXML struct {
	XMLName xml.Name `xml:"source" json:"-"`
	Channel string   `xml:"channel,attr" json:"channel"`
}
