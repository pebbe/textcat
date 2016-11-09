
all: Taalrader.class test

test:
	java Taalrader

Taalrader.class: Taalrader.java textcat.jar
	javac Taalrader.java

textcat.jar: nl/xs4all/pebbe/textcat/Textcat.java nl/xs4all/pebbe/textcat/data
	rm -f nl/xs4all/pebbe/textcat/*.class
	javac nl/xs4all/pebbe/textcat/Textcat.java
	jar vcf textcat.jar nl/xs4all/pebbe/textcat/*.class nl/xs4all/pebbe/textcat/data

nl/xs4all/pebbe/textcat/data: ../data.go java2tab.pl
	./java2tab.pl ../data.go > nl/xs4all/pebbe/textcat/data
