all: .docker

.docker: root.fxs confd.conf
	docker build -t confd_iter_c .
	touch .docker

root.fxs: root.yang
	confdc -c -o root.fxs -- root.yang

clean:
	docker rmi confd_iter_c || true
	rm .docker || true
	rm root.fxs || true
