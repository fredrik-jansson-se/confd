package confd_ll

import _ "syscall"
import "fmt"
import "unsafe"

/*
#include <confd_lib.h>
#include <confd_cdb.h>
#include <arpa/inet.h>
#include <netinet/in.h>
#include <errno.h>
#include <string.h>
#include <unistd.h>
#include <stdlib.h>
#include "helper.h"

void go_confd_init(const char* name, enum confd_debug_level lvl) {
	confd_init(name, 0, lvl);
}

int go_confd_load_schemas(const char* host, uint16_t port) {
	struct sockaddr_storage ss;
	size_t ss_size = lookup(host, port, &ss);
	if (!ss_size) {
		return 1;
	}
	return confd_load_schemas((struct sockaddr*) &ss, ss_size);
}


*/
import "C"

const (
	CONFD_SILENT      = C.CONFD_SILENT
	CONFD_DEBUG       = C.CONFD_DEBUG
	CONFD_TRACE       = C.CONFD_TRACE
	CONFD_PROTO_TRACE = C.CONFD_PROTO_TRACE
)

func Confd_init(name string, lvl int) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	C.go_confd_init(cname, C.enum_confd_debug_level(lvl))
}

func Confd_load_schemas(host string, port uint16) error {
	chost := C.CString(host)
	defer C.free(unsafe.Pointer(chost))
	res, err := C.go_confd_load_schemas(chost, C.uint16_t(port))
	if res != C.CONFD_OK {
		// return fmt.Errorf("Confd_load_schemas failed with %d: %s", res, C.GoString(C.strerror(C.errno)))
		return err
	}
	return err
}

func Confd_lasterr() error {
	return fmt.Errorf(C.GoString(C.confd_lasterr()))
}
