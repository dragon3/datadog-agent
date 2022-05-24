// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package model

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"regexp"
	"unsafe"
)

// containerIDPattern is the pattern of a container ID
var containerIDPattern = regexp.MustCompile(fmt.Sprintf(`([[:xdigit:]]{%v})`, sha256.Size*2))

// FindContainerID extracts the first sub string that matches the pattern of a container ID
func FindContainerID(s string) string {
	return containerIDPattern.FindString(s)
}

// SliceToArray copy src bytes to dst. Destination should have enough space
func SliceToArray(src []byte, dst unsafe.Pointer) {
	for i := range src {
		*(*byte)(unsafe.Pointer(uintptr(dst) + uintptr(i))) = src[i]
	}
}

// UnmarshalStringArray extract array of string for array of byte
func UnmarshalStringArray(data []byte) ([]string, error) {
	var result []string
	length := uint32(len(data))

	for i := uint32(0); i < length; {
		if i+4 >= length {
			return result, ErrStringArrayOverflow
		}
		// size of arg
		n := ByteOrder.Uint32(data[i : i+4])
		if n == 0 {
			return result, nil
		}
		i += 4

		if i+n > length {
			// truncated
			part := bytes.SplitN(data[i:length-1], []byte{0}, 2)[0]
			arg := make([]byte, len(part))
			copy(arg, part)

			return append(result, string(arg)), ErrStringArrayOverflow
		}

		part := bytes.SplitN(data[i:i+n], []byte{0}, 2)[0]
		arg := make([]byte, len(part))
		copy(arg, part)

		result = append(result, string(arg))

		i += n
	}

	return result, nil
}

// UnmarshalString unmarshal string
func UnmarshalString(data []byte, size int) (string, error) {
	if len(data) < size {
		return "", ErrNotEnoughData
	}

	part := bytes.SplitN(data[:size], []byte{0}, 2)[0]
	str := make([]byte, len(part))
	copy(str, part)

	return string(str), nil
}

// UnmarshalPrintableString unmarshal printable string
func UnmarshalPrintableString(data []byte, size int) (string, error) {
	if len(data) < size {
		return "", ErrNotEnoughData
	}

	str, err := UnmarshalString(data, size)
	if err != nil {
		return "", err
	}
	if !IsPrintable(str) {
		return "", ErrNonPrintable
	}

	return str, nil
}
