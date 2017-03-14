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

const (
	// __inline as defined in audiotransport/<predefine>:4
	__inline = 0
	// __inline__ as defined in audiotransport/<predefine>:5
	__inline__ = 0
	// __const as defined in audiotransport/<predefine>:8
	__const = 0
	// __stdc_hosted__ as defined in audiotransport/<predefine>:23
	__stdc_hosted__ = 1
	// __stdc_version__ as defined in audiotransport/<predefine>:24
	__stdc_version__ = int64(199901)
	// __stdc__ as defined in audiotransport/<predefine>:25
	__stdc__ = 1
	// __gnuc__ as defined in audiotransport/<predefine>:26
	__gnuc__ = 4
	// __flt_min__ as defined in audiotransport/<predefine>:29
	__flt_min__ = 0
	// __dbl_min__ as defined in audiotransport/<predefine>:30
	__dbl_min__ = 0
	// __ldbl_min__ as defined in audiotransport/<predefine>:31
	__ldbl_min__ = 0
	// __x86_64__ as defined in audiotransport/<predefine>:35
	__x86_64__ = 1
)
