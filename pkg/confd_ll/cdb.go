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

int _cdb_get(int sock, struct confd_value* v, const char* fmt) {
	return cdb_get(sock, v, fmt);
}

int _cdb_num_instances(int sock, const char *fmt) {
	return cdb_num_instances(sock, fmt);
}

int _cdb_cd(int sock, const char* fmt) {
	return cdb_cd(sock, fmt);
}

int _cdb_get_int64(int sock, int64_t* val, const char* fmt) {
	return cdb_get_int64(sock, val, fmt);
}

int _cdb_get_str(int sock, char* rval, int n, const char* fmt) {
	return cdb_get_str(sock, (char*) rval, n, fmt);
}
*/
import "C"

const (
	CDB_DATA_SOCKET         = C.CDB_DATA_SOCKET
	CDB_SUBSCRIPTION_SOCKET = C.CDB_SUBSCRIPTION_SOCKET

	CDB_RUNNING = C.CDB_RUNNING

	CDB_DONE_PRIORITY = C.CDB_DONE_PRIORITY
)

type Confd_value_t = C.struct_confd_value

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

func Cdb_start_session(sock int, db int) error {
	return resToError(C.cdb_start_session(C.int(sock), C.enum_cdb_db_type(db)))
}

func Cdb_end_session(sock int) error {
	return resToError(C.cdb_end_session(C.int(sock)))
}

func Cdb_set_namespace(sock int, ns int) error {
	return resToError(C.cdb_set_namespace(C.int(sock), C.int(ns)))
}

func Cdb_num_instances(sock int, path string) (uint, error) {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	res := C._cdb_num_instances(C.int(sock), cpath)

	if res < 0 {
		return 0, Confd_lasterr()
	}

	return uint(res), nil
}

func Cdb_get(sock int, v *Confd_value_t, path string) error {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	return resToError(C._cdb_get(C.int(sock), v, cpath))
}

func Cdb_cd(sock int, path string) error {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	return resToError(C._cdb_cd(C.int(sock), cpath))
}

func Cdb_get_str(sock int, path string, maxlen int) (string, error) {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	cval := make([]byte, maxlen+1)

	res := C._cdb_get_str(C.int(sock), (*C.char)(unsafe.Pointer(&cval[0])), C.int(maxlen), cpath)

	if res != C.CONFD_OK {
		return "", Confd_lasterr()
	}

	return string(cval), nil
}

func Cdb_get_int64(sock int, path string) (int64, error) {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	var val C.int64_t

	res := C._cdb_get_int64(C.int(sock), &val, cpath)

	if res != C.CONFD_OK {
		return 0, Confd_lasterr()
	}

	return int64(val), nil
}

func Cdb_read_subscription_socket(sock int, sub_points *int) (int, error) {
	var resultlen C.int

	res := C.cdb_read_subscription_socket(C.int(sock),
		(*C.int)(unsafe.Pointer(sub_points)),
		(*C.int)(unsafe.Pointer(&resultlen)))
	if res != C.CONFD_OK {
		return 0, Confd_lasterr()
	}

	return int(resultlen), nil
}

func Cdb_sync_subscription_socket(sock int, st int) error {
	return resToError(C.cdb_sync_subscription_socket(C.int(sock), C.enum_cdb_subscription_sync_type(st)))
}
