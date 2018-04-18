ULOGD-IMAGE="fcrisciani/data-server:ulogd"
DATA-SERVER="fcrisciani/data-server:fakedata"
UI-SERVER="fcrisciani/data-server:ui"

.PHONY: server ulogd

server:
	docker build -t ${DATA-SERVER} data-server

ulogd:
	docker build -t ${ULOGD-IMAGE} ulogd

push:
	docker push ${ULOGD-IMAGE}
	docker push ${DATA-SERVER}
	docker push ${UI-SERVER}

ui:
	docker build -t ${UI-SERVER} graph

all: server ulogd
