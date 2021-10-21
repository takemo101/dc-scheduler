#### diagram ####

PLANTUML_JAR_URL = https://sourceforge.net/projects/plantuml/files/plantuml.jar/download
DIAGRAM_DIRECTORY = document
DIAGRAM_EXTENSION = .puml
DIAGRAM_SRC := $(wildcard ./$(DIAGRAM_DIRECTORY)/*/*$(DIAGRAM_EXTENSION))
DIAGRAM_FULL_SRC := $(basename $(DIAGRAM_SRC))
DIAGRAM_UML := $(addsuffix ${DIAGRAM_EXTENSION}, $(DIAGRAM_FULL_SRC))
DIAGRAM_PNG := $(addsuffix .png, $(DIAGRAM_FULL_SRC))
DIAGRAM_SVG := $(addsuffix .svg, $(DIAGRAM_FULL_SRC))

# download jar
download-jar:
	curl -sSfL $(PLANTUML_JAR_URL) -o plantuml.jar

# create png
png:
	java -jar plantuml.jar -tpng $(DIAGRAM_UML)

# create svg
svg:
	java -jar plantuml.jar -tsvg $(DIAGRAM_UML)

# clear
clear:
	make clear-resource
	rm -f plantuml.jar

# clear-resource
clear-resource:
	rm -f $(DIAGRAM_PNG) $(DIAGRAM_SVG)

# output png
png/%: $(DIAGRAM_DIRECTORY)/%${DIAGRAM_EXTENSION}
	java -jar plantuml.jar -tpng $^

# output svg
svg/%: $(DIAGRAM_DIRECTORY)/%${DIAGRAM_EXTENSION}
	java -jar plantuml.jar -tsvg $^


#### docker ####

BINARY_DIRECTORY = bin
LOLI_SSH_KEY = .ssh/id_rsa
LOLI_SSH_PORT = 35135
LOLI_SSH_USER = proud-iki-7985
LOLI_SSH_HOST = ssh-1.mc.lolipop.jp
LOLI_PROJECT_DIRECTORY = /var/app/dc-scheduler

# app resource setup
app-setup:
	cp config.example.yml config.yml
	cp config.testing.example.yml config.testing.yml
	npm run prod

# container build
docker-build:
	docker-compose build --no-cache mysql pma redis mailhog

# container start
docker-start:
	docker-compose up -d mysql pma redis mailhog

# container stop
docker-stop:
	docker-compose stop

# app test
app-testing:
	rm -f fiber.testing.sqlite
	go test -v ./test

# app build
app-build:
	go build -o ./go-app ./main.go
	go build -o ./go-app-cli ./cli/main.go

# brew install FiloSottile/musl-cross/musl-cross を実行してから
# app build for linux
app-linux-build:
	mkdir -p ./${BINARY_DIRECTORY}
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=1 CC=/usr/local/bin/x86_64-linux-musl-cc go build --ldflags '-linkmode external -extldflags "-static"' -a -v -o ./${BINARY_DIRECTORY}/go-linux-app ./main.go
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=1 CC=/usr/local/bin/x86_64-linux-musl-cc go build --ldflags '-linkmode external -extldflags "-static"' -a -v -o ./${BINARY_DIRECTORY}/go-linux-app-cli ./cli/main.go

# app deploy for lolipop
app-loli-deploy:
	make app-linux-build
	ssh -p ${LOLI_SSH_PORT} -i ${LOLI_SSH_KEY} ${LOLI_SSH_USER}@${LOLI_SSH_HOST} 'mkdir -p ${LOLI_PROJECT_DIRECTORY}/${BINARY_DIRECTORY}'
	scp -r -P ${LOLI_SSH_PORT} -i ${LOLI_SSH_KEY} ./${BINARY_DIRECTORY}/go-linux-app ${LOLI_SSH_USER}@${LOLI_SSH_HOST}:${LOLI_PROJECT_DIRECTORY}/${BINARY_DIRECTORY}
	scp -r -P ${LOLI_SSH_PORT} -i ${LOLI_SSH_KEY} ./${BINARY_DIRECTORY}/go-linux-app-cli ${LOLI_SSH_USER}@${LOLI_SSH_HOST}:${LOLI_PROJECT_DIRECTORY}/${BINARY_DIRECTORY}

app-loli-deply-copy:
	scp -r -P ${LOLI_SSH_PORT} -i ${LOLI_SSH_KEY} ./config ${LOLI_SSH_USER}@${LOLI_SSH_HOST}:${LOLI_PROJECT_DIRECTORY}
	scp -r -P ${LOLI_SSH_PORT} -i ${LOLI_SSH_KEY} ./resource ${LOLI_SSH_USER}@${LOLI_SSH_HOST}:${LOLI_PROJECT_DIRECTORY}
	scp -r -P ${LOLI_SSH_PORT} -i ${LOLI_SSH_KEY} ./static ${LOLI_SSH_USER}@${LOLI_SSH_HOST}:${LOLI_PROJECT_DIRECTORY}
	scp -r -P ${LOLI_SSH_PORT} -i ${LOLI_SSH_KEY} ./config.example.yml ${LOLI_SSH_USER}@${LOLI_SSH_HOST}:${LOLI_PROJECT_DIRECTORY}
	scp -r -P ${LOLI_SSH_PORT} -i ${LOLI_SSH_KEY} ./config.testing.example.yml ${LOLI_SSH_USER}@${LOLI_SSH_HOST}:${LOLI_PROJECT_DIRECTORY}
	scp -r -P ${LOLI_SSH_PORT} -i ${LOLI_SSH_KEY} ./Makefile ${LOLI_SSH_USER}@${LOLI_SSH_HOST}:${LOLI_PROJECT_DIRECTORY}
	ssh -p ${LOLI_SSH_PORT} -i ${LOLI_SSH_KEY} ${LOLI_SSH_USER}@${LOLI_SSH_HOST} 'mkdir -p ${LOLI_PROJECT_DIRECTORY}/storage/public'

loli-login:
	ssh -p 35135 -i .ssh/id_rsa proud-iki-7985@ssh-1.mc.lolipop.jp
