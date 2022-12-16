package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/paypal/gatt"
	"github.com/paypal/gatt/examples/option"
	"github.com/zmb3/ble"
	"github.com/zmb3/ble/linux"
)

// Check if device supports secure connections
func checkSecureConnections(p gatt.Peripheral) {
	for _, service := range p.Services() {
		if service.UUID().Equal(gatt.MustParseUUID("1803")) {
			data, err := p.ReadCharacteristic(service.Characteristics()[0])
			if err != nil {
				log.Println("Failed to read characteristic:", err)
				return
			}
			if data[0]&0x02 > 0 {
				fmt.Println("  Device supports secure connections")
			} else {
				fmt.Println("  Device does NOT support secure connections")
			}
			break
		}
	}
}

// Check if device is using a weak encryption key
func checkEncryptionKeySize(p gatt.Peripheral) {
	for _, service := range p.Services() {
		if service.UUID().Equal(gatt.MustParseUUID("1803")) {
			data, err := p.ReadCharacteristic(service.Characteristics()[1])
			if err != nil {
				log.Println("Failed to read characteristic:", err)
				return
			}
			keySize := data[0]
			if keySize < 16 {
				fmt.Println("  Device is using a weak encryption key (", keySize, " bytes)")
			} else {
				fmt.Println("  Device is using a strong encryption key (", keySize, " bytes)")
			}
			break
		}
	}
}

// Check if device is using a known-vulnerable Bluetooth protocol version
func checkProtocolVersion(p gatt.Peripheral) {
	for _, service := range p.Services() {
		if service.UUID().Equal(gatt.MustParseUUID("1803")) {
			data, err := p.ReadCharacteristic(service.Characteristics()[2])
			if err != nil {
				log.Println("Failed to read characteristic:", err)
				return
			}
			hciVersion := data[0]
			lmpVersion := data[1]
			fmt.Println("  HCI version:", hciVersion)
			fmt.Println("  LMP version:", lmpVersion)
	
