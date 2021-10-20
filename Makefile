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
	mkdir ./${BINARY_DIRECTORY}
	go build -o ./${BINARY_DIRECTORY}/go-app ./main.go
	go build -o ./${BINARY_DIRECTORY}/go-app-cli ./cli/main.go


# heroku script
heroku-script
	make app-build
	cp ./.heroku/config.heroku.yml ./config.yml
