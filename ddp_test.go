package ddp

import (
	"fmt"
	"reflect"
	"testing"
)

func TestPacketDecoding(t *testing.T) {
	testCases := []struct {
		description string
		bytes       []byte
		expected    DDPPacket
	}{
		{
			description: "simple DDP message",
			bytes:       []byte{0x40, 0x0f, 0x92, 0x01, 0x45, 0x67, 0x89, 0x10, 0x10, 0x11},
			expected: DDPPacket{
				Header: DDPHeader{
					Flags: DDPHeaderFlags{
						Version:  1,
						Timecode: false,
						Storage:  false,
						Reply:    false,
						Query:    false,
						Push:     false,
					},
					SequenceNumber: 15,
					DataType: DDPDataType{
						Standard:     false,
						Type:         2,
						BitsPerPixel: 4,
					},
					SourceOrDestinationID: 1,
					DataOffset:            1164413200,
					DataLength:            4113,
				},
				Data: []byte{},
			},
		},
		{
			description: "DDP message with timecode",
			bytes:       []byte{0x50, 0x0f, 0x92, 0x01, 0x45, 0x67, 0x89, 0x10, 0x10, 0x11, 0x10, 0x11, 0x11, 0x10},
			expected: DDPPacket{
				Header: DDPHeader{
					Flags: DDPHeaderFlags{
						Version:  1,
						Timecode: true,
						Storage:  false,
						Reply:    false,
						Query:    false,
						Push:     false,
					},
					SequenceNumber: 15,
					DataType: DDPDataType{
						Standard:     false,
						Type:         2,
						BitsPerPixel: 4,
					},
					SourceOrDestinationID: 1,
					DataOffset:            1164413200,
					DataLength:            4113,
					Timecode: &DDPTimecode{
						Seconds:           4113,
						FractionalSeconds: 4368,
					},
				},
				Data: []byte{},
			},
		},
	}

	for _, testCase := range testCases {

		actual, err := Decode(testCase.bytes)

		if err != nil {
			t.Errorf("Test '%s' failed to decode chunk properly", testCase.description)
			fmt.Println(err)
		}

		if !reflect.DeepEqual(actual, testCase.expected) {
			t.Errorf("Test '%s' failed to decode chunk properly", testCase.description)
			fmt.Printf("expected: %+v\n", testCase.expected)
			fmt.Printf("actual: %+v\n", actual)
		}
	}
}

func TestPacketEncoding(t *testing.T) {
	testCases := []struct {
		description string
		expected    []byte
		packet      DDPPacket
	}{
		{
			description: "simple DDP message",
			expected:    []byte{0x40, 0x0f, 0x92, 0x01, 0x45, 0x67, 0x89, 0x10, 0x10, 0x11},
			packet: DDPPacket{
				Header: DDPHeader{
					Flags: DDPHeaderFlags{
						Version:  1,
						Timecode: false,
						Storage:  false,
						Reply:    false,
						Query:    false,
						Push:     false,
					},
					SequenceNumber: 15,
					DataType: DDPDataType{
						Standard:     false,
						Type:         2,
						BitsPerPixel: 4,
					},
					SourceOrDestinationID: 1,
					DataOffset:            1164413200,
					DataLength:            4113,
				},
				Data: []byte{},
			},
		},
		{
			description: "DDP message with timecode",
			expected:    []byte{0x50, 0x0f, 0x92, 0x01, 0x45, 0x67, 0x89, 0x10, 0x10, 0x11, 0x10, 0x11, 0x11, 0x10},
			packet: DDPPacket{
				Header: DDPHeader{
					Flags: DDPHeaderFlags{
						Version:  1,
						Timecode: true,
						Storage:  false,
						Reply:    false,
						Query:    false,
						Push:     false,
					},
					SequenceNumber: 15,
					DataType: DDPDataType{
						Standard:     false,
						Type:         2,
						BitsPerPixel: 4,
					},
					SourceOrDestinationID: 1,
					DataOffset:            1164413200,
					DataLength:            4113,
					Timecode: &DDPTimecode{
						Seconds:           4113,
						FractionalSeconds: 4368,
					},
				},
				Data: []byte{},
			},
		},
	}

	for _, testCase := range testCases {

		actual, err := Encode(testCase.packet)

		if err != nil {
			t.Errorf("Test '%s' failed to decode chunk properly", testCase.description)
			fmt.Println(err)
		}

		if !reflect.DeepEqual(actual, testCase.expected) {
			t.Errorf("Test '%s' failed to decode chunk properly", testCase.description)
			fmt.Printf("expected: %+v\n", testCase.expected)
			fmt.Printf("actual: %+v\n", actual)
		}
	}
}
