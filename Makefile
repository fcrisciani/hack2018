ULOGD-IMAGE="fcrisciani/argus:ulogd"
DATA-SERVER="fcrisciani/argus:data-server"
UI-SERVER="fcrisciani/argus:ui"

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
