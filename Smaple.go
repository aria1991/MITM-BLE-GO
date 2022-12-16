package main

import (
	"fmt"
	"log"

	"github.com/paypal/gatt"
	"github.com/paypal/gatt/examples/option"
)

func onStateChanged(d gatt.Device, s gatt.State) {
	fmt.Println("State:", s)
	switch s {
	case gatt.StatePoweredOn:
		fmt.Println("Scanning for devices...")
		d.Scan([]gatt.UUID{}, false)
		return
	default:
		d.StopScanning()
	}
}

func onPeripheralDiscovered(p gatt.Peripheral, a *gatt.Advertisement, rssi int) {
	fmt.Printf("Peripheral discovered: %s (RSSI %d)\n", p.ID(), rssi)
	fmt.Println("  Name:", a.LocalName)
	fmt.Println("  Services:", p.Services())

	// Check if the device has the "MITM" flag set in the "Security Manager"
	// service, indicating that it supports Man-in-the-Middle (MITM) protection
	for _, service := range p.Services() {
		if service.UUID().Equal(gatt.MustParseUUID("1803")) {
			data, err := p.ReadCharacteristic(service.Characteristics()[0])
			if err != nil {
				log.Println("Failed to read characteristic:", err)
				continue
			}
			if data[0]&0x04 > 0 {
				fmt.Println("  Device supports MITM protection")
			} else {
				fmt.Println("  Device does NOT support MITM protection")
			}
			break
		}
	}
}

func main() {
	d, err := gatt.NewDevice(option.DefaultClientOptions...)
	if err != nil {
		log.Fatalf("Failed to open device, err: %s\n", err)
		return
	}

	// Register handlers
	d.Handle(
		gatt.PeripheralDiscovered(onPeripheralDiscovered),
		gatt.StateChanged(onStateChanged),
	)

	// Start the device
	d.Init(onStateChanged)
	select {}
}
