package confd_ll

/*
#include <arpa/inet.h>
#include <errno.h>
#include <string.h>
#include <confd_lib.h>

size_t lookup(const char* host, uint16_t port, struct sockaddr_storage* ss) {
	memset(ss, 0, sizeof(struct sockaddr_storage));

	struct sockaddr_in* sin = (struct sockaddr_in*) ss;
	struct sockaddr_in6* sin6 = (struct sockaddr_in6*) ss;

	if (1 == inet_pton(AF_INET, host, &sin->sin_addr)) {
		sin->sin_port=htons(port);
		sin->sin_family = AF_INET;
		return sizeof(struct sockaddr_in);
	}
	else if (1 == inet_pton(AF_INET6, host, &sin6->sin6_addr)) {
		sin6->sin6_port=htons(port);
		sin6->sin6_family = AF_INET;
		return sizeof(struct sockaddr_in6);
	}
	errno = EHOSTUNREACH;
	return 0;
}

int64_t _CONFD_GET_INT64(const confd_value_t* val) {
	return CONFD_GET_INT64(val);
}

*/
import "C"

func resToError(res C.int) error {
	if res != C.CONFD_OK {
		return Confd_lasterr()
	}
	return nil
}

func CONFD_GET_INT64(val *Confd_value_t) int64 {
	return int64(C._CONFD_GET_INT64(val))
}
