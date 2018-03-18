package main

import (
	"fmt"
	"github.com/fredrik-jansson-se/confd/pkg/confd_ll"
	"syscall"
)

const (
	MAXH = 3
	MAXC = 2
)

type child struct {
	dn        int64
	childattr string
	inuse     bool
}

type rfhead struct {
	dn        int64
	sector_id string
	children  [MAXC]child
	inuse     bool
}

var rfheads [MAXH]rfhead

func main() {

	confd_ll.Confd_init("probe", confd_ll.CONFD_DEBUG)

	sock, err := confd_ll.Cdb_connect("127.0.0.1", 4565, confd_ll.CDB_DATA_SOCKET)
	if err != nil {
		panic(err)
	}

	subsock, err := confd_ll.Cdb_connect("127.0.0.1", 4565, confd_ll.CDB_SUBSCRIPTION_SOCKET)
	if err != nil {
		panic(err)
	}

	err = confd_ll.Confd_load_schemas("127.0.0.1", 4565)
	if err != nil {
		panic(err)
	}

	spoint, err := confd_ll.Cdb_subscribe(subsock, 3, 0, "/root/NodeB/RFHead")
	if err != nil {
		panic(err)
	}

	fmt.Printf("spoint = %d\n", spoint)

	err = confd_ll.Cdb_subscribe_done(subsock)
	if err != nil {
		panic(err)
	}

	err = read_db(sock)
	if err != nil {
		panic(err)
	}

	dump_db()

	rfds := &syscall.FdSet{}
	for {
		FD_ZERO(rfds)
		FD_SET(rfds, subsock)

		syscall.Select(subsock+1, rfds, nil, nil, nil)
		if err != nil {
			panic(err)
		}

		if FD_ISSET(rfds, subsock) {
			sub_points := make([]int, 1)
			reslen, err := confd_ll.Cdb_read_subscription_socket(subsock,
				&sub_points[0])

			if err != nil {
				panic(err)
			}

			fmt.Printf("reslen = %d\n", reslen)
			if reslen <= 0 {
				continue
			}

			err = confd_ll.Cdb_start_session(sock, confd_ll.CDB_RUNNING)
			if err != nil {
				panic(err)
			}

			confd_ll.Cdb_end_session(sock)

			//     cdb_diff_iterate(subsock, sub_points[0], iter,
			//                      ITER_WANT_PREV, (void*)&sock);
			//     cdb_end_session(sock);

			//     /* Here is an alternative approach to checking a subtree */
			//     /* the function below will invoke cdb_diff_iterate */
			//     /* and check if any changes have beem made in the tagpath */
			//     /* described by tags[] */
			//     /* This still only applies to the subscription point which */
			//     /* is being used */

			//     struct xml_tag tags[] = {{root_root, root__ns},
			//                              {root_NodeB, root__ns},
			//                              {root_RFHead, root__ns},
			//                              {root_Child, root__ns}};
			//     int tagslen = sizeof(tags)/sizeof(tags[0]);
			//     /* /root/NodeB/RFHead/Child */
			//     int retv = cdb_diff_match(subsock, sub_points[0],
			//                               tags, tagslen);
			//     fprintf(stderr, "Diff match: %s\n", retv ? "yes" : "no");
			// }

			err = confd_ll.Cdb_sync_subscription_socket(subsock, confd_ll.CDB_DONE_PRIORITY)
			if err != nil {
				panic(err)
			}
			dump_db()

		}

	}

	return
}

func read_db(cdbsock int) error {
	err := confd_ll.Cdb_start_session(cdbsock, confd_ll.CDB_RUNNING)
	if err != nil {
		return err
	}

	n, err := confd_ll.Cdb_num_instances(cdbsock, "/root/NodeB/RFHead")
	if err != nil {
		return err
	}

	for i := uint(0); i < n; i++ {
		var key confd_ll.Confd_value_t
		err = confd_ll.Cdb_get(cdbsock, &key, fmt.Sprintf("/root/NodeB/RFHead[%d]/dn", i))
		if err != nil {
			return err
		}
		read_head(cdbsock, &key)
	}

	return confd_ll.Cdb_end_session(cdbsock)
}

func dump_db() {
	for i := 0; i < MAXH; i++ {
		if !rfheads[i].inuse {
			continue
		}
		fmt.Printf("HEAD %d <%s>\n", rfheads[i].dn, rfheads[i].sector_id)
		for j := 0; j < MAXC; j++ {
			if !rfheads[i].children[j].inuse {
				continue
			}
			fmt.Printf("   Child %d  <<%s>>\n",
				rfheads[i].children[j].dn,
				rfheads[i].children[j].childattr)
		}
	}
}

func read_head(cdbsock int, headkey *confd_ll.Confd_value_t) {
	pos := -1

	for i := 0; i < MAXH; i++ {
		if confd_ll.CONFD_GET_INT64(headkey) == rfheads[i].dn {
			pos = i
			break
		}
	}

	if pos == -1 { // pick first
		for i := 0; i < MAXH; i++ {
			if !rfheads[i].inuse {
				pos = i
				break
			}
		}
	}

	fmt.Printf("Picking %d\n", pos)

	hp := &rfheads[pos]

	err := confd_ll.Cdb_cd(cdbsock, fmt.Sprintf("/root/NodeB/RFHead{%d}", confd_ll.CONFD_GET_INT64(headkey)))
	if err != nil {
		panic(err)
	}

	hp.dn = confd_ll.CONFD_GET_INT64(headkey)
	hp.inuse = true
	hp.sector_id, err = confd_ll.Cdb_get_str(cdbsock, "SECTORID_ID", 255)
	if err != nil {
		panic(err)
	}

	n, err := confd_ll.Cdb_num_instances(cdbsock, "Child")
	for i := 0; i < MAXC; i++ {
		hp.children[i].inuse = false
	}

	for i := uint(0); i < n; i++ {
		hp.children[i].dn, err = confd_ll.Cdb_get_int64(cdbsock, fmt.Sprintf("Child[%d]/cdn", i))
		if err != nil {
			panic(err)
		}
		hp.children[i].childattr, err = confd_ll.Cdb_get_str(cdbsock, fmt.Sprintf("Child[%d]/childAttr", i), 255)
		if err != nil {
			panic(err)
		}
		hp.children[i].inuse = true
	}
}

func FD_SET(p *syscall.FdSet, i int) {
	p.Bits[i/64] |= 1 << uint(i) % 64
}

func FD_ISSET(p *syscall.FdSet, i int) bool {
	return (p.Bits[i/64] & (1 << uint(i) % 64)) != 0
}

func FD_ZERO(p *syscall.FdSet) {
	for i := range p.Bits {
		p.Bits[i] = 0
	}
}
