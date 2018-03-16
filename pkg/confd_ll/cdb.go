package confd_ll

import _ "fmt"
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

int go_cdb_connect(const char* host, uint16_t port, int cdb_sock_type) {
	struct sockaddr_storage ss;
	size_t ss_size = lookup(host, port, &ss);
	if (!ss_size) {
		return -1;
	}

	struct sockaddr_in* sin = (struct sockaddr_in*) &ss;
	struct sockaddr_in6* sin6 = (struct sockaddr_in6*) &ss;
	const struct sockaddr* sa = (struct sockaddr*) &ss;

	int s = socket(ss.ss_family, SOCK_STREAM, IPPROTO_TCP);

	if (s < 0) {
		return -1;
	}

	if (CONFD_OK != cdb_connect(s, cdb_sock_type, sa, ss_size)) {
		close(s);
		return -1;
	}

	return s;
}

int _cdb_subscribe(int sock, int priority, int nspace, int *spoint, const char *fmt) {
	return cdb_subscribe(sock, priority, nspace, spoint, fmt);
}
*/
import "C"

const (
	CDB_DATA_SOCKET         = C.CDB_DATA_SOCKET
	CDB_SUBSCRIPTION_SOCKET = C.CDB_SUBSCRIPTION_SOCKET
)

func Cdb_connect(host string, port uint16, cdb_sock_type int) (int, error) {
	chost := C.CString(host)
	defer C.free(unsafe.Pointer(chost))
	res, err := C.go_cdb_connect(chost, C.uint16_t(port), C.int(cdb_sock_type))
	return int(res), err
}

func Cdb_subscribe(sock int, priority int, nspace int, path string) (int, error) {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	var spoint C.int

	res := C._cdb_subscribe(C.int(sock), C.int(priority), C.int(nspace), &spoint, cpath)

	if res != C.CONFD_OK {
		return -1, Confd_lasterr()
	}

	return int(spoint), nil
}

func Cdb_subscribe_done(sock int) error {
	if C.CONFD_OK != C.cdb_subscribe_done(C.int(sock)) {
		return Confd_lasterr()
	}
	return nil
}
