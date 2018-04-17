ULOGD-IMAGE="fcrisciani/data-server:ulogd"
DATA-SERVER="fcrisciani/data-server:fakedata"

.PHONY: server ulogd

server:
	docker build -t ${DATA-SERVER} data-server

ulogd:
	docker build -t ${ULOGD-IMAGE} ulogd

push:
	docker push ${ULOGD-IMAGE}
	docker push ${DATA-SERVER}

all: server ulogd
