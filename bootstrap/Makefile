
all: ../data.go

pat.go: ../patterns.go
	echo '// THIS IS A GENERATED FILE' > pat.go
	echo '// DO NOT EDIT' >> pat.go
	echo >> pat.go
	perl -p -e 's/package textcat/package main/; s/GetPatterns/getPatterns/; s/MaxPatterns/maxPatterns/' ../patterns.go >> pat.go

bootstrap: pat.go bootstrap.go
	go build .

../data.go: bootstrap ../samples*/*.txt
	./bootstrap ../samples*/*.txt
	gofmt -w ../data.go
