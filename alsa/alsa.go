// WARNING: This file has automatically been generated on Tue, 14 Mar 2017 01:22:35 EDT.
// By https://git.io/c-for-go. DO NOT EDIT.

package alsa

/*
#cgo LDFLAGS: -L. -lalsa -lasound -lpulse -lpulse-simple
#include <alsa/asoundlib.h>
#include <pulse/simple.h>
#include <pulse/error.h>
#include "../alsa-bindings/alsa.h"
#include "../alsa-bindings/pulse.h"
#include <stdlib.h>
#include "cgo_helpers.h"
*/
import "C"
import "unsafe"

// Init_playback function as declared in alsa-bindings/alsa.h:1
func Init_playback(Device string, Samplerate int32, Channels int32) int32 {
	cDevice, _ := unpackPCharString(Device)
	cSamplerate, _ := (C.int)(Samplerate), cgoAllocsUnknown
	cChannels, _ := (C.int)(Channels), cgoAllocsUnknown
	__ret := C.init_playback(cDevice, cSamplerate, cChannels)
	__v := (int32)(__ret)
	return __v
}

// Alsa_writei function as declared in alsa-bindings/alsa.h:2
func Alsa_writei(b *[]byte, Len int32) int32 {
	addr := &((*b)[0])
	cBytes := unsafe.Pointer(addr)
	cLen, _ := (C.int)(Len), cgoAllocsUnknown
	__ret := C.alsa_writei(cBytes, cLen)
	__v := (int32)(__ret)
	return __v
}

// Play_bytes function as declared in alsa-bindings/alsa.h:3
func Play_bytes(HandleIdx int32, b *[]byte, Len int32) int32 {
	cHandleIdx, _ := (C.int)(HandleIdx), cgoAllocsUnknown
	addr := &((*b)[0])
	cBytes := unsafe.Pointer(addr)
	cLen, _ := (C.int)(Len), cgoAllocsUnknown
	__ret := C.play_bytes(cHandleIdx, cBytes, cLen)
	__v := (int32)(__ret)
	return __v
}

// Close_playback function as declared in alsa-bindings/alsa.h:4
func Close_playback(HandleIdx int32) int32 {
	cHandleIdx, _ := (C.int)(HandleIdx), cgoAllocsUnknown
	__ret := C.close_playback(cHandleIdx)
	__v := (int32)(__ret)
	return __v
}

// Pa_init function as declared in alsa-bindings/pulse.h:1
func Pa_init(Name string, Device string, Samplerate int32, Channels int32, IsPlayback int32) int32 {
	cName, _ := unpackPCharString(Name)
	cDevice, _ := unpackPCharString(Device)
	cSamplerate, _ := (C.int)(Samplerate), cgoAllocsUnknown
	cChannels, _ := (C.int)(Channels), cgoAllocsUnknown
	cIsPlayback, _ := (C.int)(IsPlayback), cgoAllocsUnknown
	__ret := C.pa_init(cName, cDevice, cSamplerate, cChannels, cIsPlayback)
	__v := (int32)(__ret)
	return __v
}

// Pa_read function as declared in alsa-bindings/pulse.h:2
func Pa_handle_read(Idx int32, b *[]byte, Len int32) int32 {
	cIdx, _ := (C.int)(Idx), cgoAllocsUnknown
	addr := &((*b)[0])
	cBuf := unsafe.Pointer(addr)
	cLen, _ := (C.int)(Len), cgoAllocsUnknown
	__ret := C.pa_handle_read(cIdx, cBuf, cLen)
	__v := (int32)(__ret)
	return __v
}

// Pa_write function as declared in alsa-bindings/pulse.h:3
func Pa_handle_write(Idx int32, b *[]byte, Len int32) int32 {
	cIdx, _ := (C.int)(Idx), cgoAllocsUnknown
	addr := &((*b)[0])
	cBuf := unsafe.Pointer(addr)
	cLen, _ := (C.int)(Len), cgoAllocsUnknown
	__ret := C.pa_handle_write(cIdx, cBuf, cLen)
	__v := (int32)(__ret)
	return __v
}

// Pa_release function as declared in alsa-bindings/pulse.h:4
func Pa_release(Idx int32) int32 {
	cIdx, _ := (C.int)(Idx), cgoAllocsUnknown
	__ret := C.pa_release(cIdx)
	__v := (int32)(__ret)
	return __v
}
