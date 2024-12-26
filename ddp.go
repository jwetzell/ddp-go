package ddp

import (
	"encoding/binary"
	"errors"
)

type DDPDataType struct {
	Standard     bool
	Type         uint8
	BitsPerPixel uint8
}
type DDPHeaderFlags struct {
	Version  uint8
	Timecode bool
	Storage  bool
	Reply    bool
	Query    bool
	Push     bool
}

type DDPTimecode struct {
	Seconds           uint16
	FractionalSeconds uint16
}

type DDPHeader struct {
	Timecode              *DDPTimecode
	DataOffset            uint32
	DataLength            uint16
	Flags                 DDPHeaderFlags
	DataType              DDPDataType
	SequenceNumber        uint8
	SourceOrDestinationID uint8
}

type DDPPacket struct {
	Header DDPHeader
	Data   []byte
}

func Decode(bytes []byte) (DDPPacket, error) {
	if len(bytes) < 10 {
		return DDPPacket{}, errors.New("DDP packet must be at least 10 bytes")
	}

	flagsByte := bytes[0]

	flags := DDPHeaderFlags{
		Version:  (flagsByte & 0xc0) >> 6,
		Timecode: (flagsByte&0x10)>>4 == 1,
		Storage:  (flagsByte&0x08)>>3 == 1,
		Reply:    (flagsByte&0x04)>>2 == 1,
		Query:    (flagsByte&0x02)>>1 == 1,
		Push:     (flagsByte & 0x01) == 1,
	}

	sequenceNumber := bytes[1] & 0x0f

	dataTypeByte := bytes[2]

	dataType := DDPDataType{
		Standard:     dataTypeByte>>7 == 0,
		Type:         (dataTypeByte & 0x38) >> 3,
		BitsPerPixel: dataTypeByte & 0x07,
	}

	if dataType.BitsPerPixel > 1 {
		if dataType.BitsPerPixel == 2 {
			dataType.BitsPerPixel = 4
		} else if dataType.BitsPerPixel == 3 {
			dataType.BitsPerPixel = 8
		} else if dataType.BitsPerPixel == 4 {
			dataType.BitsPerPixel = 16
		} else if dataType.BitsPerPixel == 5 {
			dataType.BitsPerPixel = 24
		} else if dataType.BitsPerPixel == 6 {
			dataType.BitsPerPixel = 32
		} else {
			return DDPPacket{}, errors.New("unknown DDP size")
		}
	}

	sourceOrDestinationID := bytes[3]
	dataOffset := (uint32(bytes[4]) << 24) + (uint32(bytes[5]) << 16) + (uint32(bytes[6]) << 8) + uint32(bytes[7])
	dataLength := (uint16(bytes[8]) << 8) + uint16(bytes[9])

	header := DDPHeader{
		Flags:                 flags,
		SequenceNumber:        sequenceNumber,
		DataType:              dataType,
		SourceOrDestinationID: sourceOrDestinationID,
		DataOffset:            dataOffset,
		DataLength:            dataLength,
	}

	dataStart := 10
	if flags.Timecode {
		if len(bytes) < 14 {
			return DDPPacket{}, errors.New("DDP packet with timecode must be at least 14 bytes")
		}
		dataStart = 14
		header.Timecode = &DDPTimecode{
			Seconds:           uint16(bytes[10])<<8 + uint16(bytes[11]),
			FractionalSeconds: uint16(bytes[12])<<8 + uint16(bytes[13]),
		}
	}

	return DDPPacket{
		Header: header,
		Data:   bytes[dataStart:],
	}, nil
}

func Encode(packet DDPPacket) ([]byte, error) {
	flagByte := packet.Header.Flags.Version << 6

	if packet.Header.Flags.Timecode {
		flagByte += 16
	}

	if packet.Header.Flags.Storage {
		flagByte += 8
	}
	if packet.Header.Flags.Reply {
		flagByte += 4
	}
	if packet.Header.Flags.Query {
		flagByte += 2
	}
	if packet.Header.Flags.Push {
		flagByte += 1
	}

	var dataTypeByte uint8 = 0

	if !packet.Header.DataType.Standard {
		dataTypeByte += 128
	}

	dataTypeByte += packet.Header.DataType.Type << 3

	if packet.Header.DataType.BitsPerPixel == 1 {
		dataTypeByte += 1
	} else if packet.Header.DataType.BitsPerPixel == 4 {
		dataTypeByte += 2
	} else if packet.Header.DataType.BitsPerPixel == 8 {
		dataTypeByte += 3
	} else if packet.Header.DataType.BitsPerPixel == 16 {
		dataTypeByte += 4
	} else if packet.Header.DataType.BitsPerPixel == 24 {
		dataTypeByte += 5
	} else if packet.Header.DataType.BitsPerPixel == 32 {
		dataTypeByte += 6
	} else if packet.Header.DataType.BitsPerPixel == 0 {
		dataTypeByte += 0
	} else {
		return []byte{}, errors.New("DDP packet has unsupported bits per pixel value")
	}

	bytes := []byte{}

	bytes = append(bytes, flagByte)
	bytes = append(bytes, packet.Header.SequenceNumber)
	bytes = append(bytes, dataTypeByte)
	bytes = append(bytes, packet.Header.SourceOrDestinationID)

	bytes = binary.BigEndian.AppendUint32(bytes, packet.Header.DataOffset)
	bytes = binary.BigEndian.AppendUint16(bytes, packet.Header.DataLength)

	if packet.Header.Flags.Timecode {
		if packet.Header.Timecode != nil {
			bytes = binary.BigEndian.AppendUint16(bytes, packet.Header.Timecode.Seconds)
			bytes = binary.BigEndian.AppendUint16(bytes, packet.Header.Timecode.FractionalSeconds)
		} else {
			return []byte{}, errors.New("DDP packet had timecode flag set but no timecode value supplied")
		}
	}

	// const bytes = new Uint8Array([flagByte, sequenceNumberByte, dataTypeByte, packet.Header.sourceOrDestinationID])

	bytes = append(bytes, packet.Data...)
	return bytes, nil
}
