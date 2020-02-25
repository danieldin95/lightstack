package api

import (
	"github.com/danieldin95/lightstar/compute/libvirtc"
	"github.com/danieldin95/lightstar/http/schema"
	"github.com/danieldin95/lightstar/libstar"
	"github.com/danieldin95/lightstar/storage"
	"github.com/danieldin95/lightstar/storage/libvirts"
	"github.com/gorilla/mux"
	"net/http"
	"path"
	"strings"
)

type Disk struct {

}

func DiskConf2XML(conf *schema.DiskConf) (*libvirtc.DiskXML, error) {
	// create new disk firstly.
	size := libstar.ToBytes(conf.Size, conf.Unit)
	slot := libstar.H2D8(conf.Slot)
	vol, err := NewVolume(conf.Name, Slot2Disk(slot), size)
	if err != nil {
		return nil, err
	}
	xml := libvirtc.DiskXML{
		Type:   "file",
		Device: "disk",
		Driver: libvirtc.DiskDriverXML{
			Name: "qemu",
			Type: vol.Target.Format.Type,
		},
		Source: libvirtc.DiskSourceXML{
			File: vol.Target.Path,
		},
		Target: libvirtc.DiskTargetXML{
			Bus: conf.Bus,
			Dev: libvirtc.DISK.Slot2Dev(conf.Bus, slot),
		},
		Address: libvirtc.AddressXML{
			Type:   "pci",
			Domain: "0x0000",
			Bus:    libvirtc.DISK_BUS,
			Slot:   conf.Slot,
		},
	}
	return &xml, nil
}

func (disk Disk) Router(router *mux.Router) {
	router.HandleFunc("/api/instance/{id}/disk", disk.POST).Methods("POST")
	router.HandleFunc("/api/instance/{id}/disk/{dev}", disk.DELETE).Methods("DELETE")
}

func (disk Disk) POST(w http.ResponseWriter, r *http.Request) {
	conf := &schema.DiskConf{}
	if err := GetData(r, conf); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	uuid, _ := GetArg(r, "id")
	dom, err := libvirtc.LookupDomainByUUIDString(uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dom.Free()

	if conf.Name == "" {
		conf.Name, _ = dom.GetName()
	}
	xmlObj, err := DiskConf2XML(conf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	libstar.Debug("Disk.POST: %s", xmlObj.Encode())
	flags := libvirtc.DOMAIN_DEVICE_MODIFY_SAVE
	if err := dom.AttachDeviceFlags(xmlObj.Encode(), flags); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ResponseMsg(w, 0, "success")
}

func (disk Disk) DELETE(w http.ResponseWriter, r *http.Request) {
	uuid, _ := GetArg(r, "id")
	dom, err := libvirtc.LookupDomainByUUIDString(uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dom.Free()

	dev, _ := GetArg(r, "dev")
	xml := libvirtc.NewDomainXMLFromDom(dom, true)
	if xml == nil {
		http.Error(w, "Cannot get domain's descXML", http.StatusInternalServerError)
		return
	}

	if xml.Devices.Disks != nil {
		for _, disk := range xml.Devices.Disks {
			if disk.Target.Dev != dev {
				continue
			}
			// found deivice
			flags := libvirtc.DOMAIN_DEVICE_MODIFY_SAVE
			if err := dom.DetachDeviceFlags(disk.Encode(), flags); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			file := disk.Source.File
			if strings.HasPrefix(file, storage.Location) &&
				(strings.HasSuffix(file, ".img") || strings.HasSuffix(file, ".qcow2")) {
				dir := path.Dir(file)
				volume := path.Base(file)
				pool := path.Base(dir)
				libvirts.RemoveVolume(libvirts.ToDomainPool(pool), volume)
			}
		}
	}
	ResponseMsg(w, 0, "success")
}
