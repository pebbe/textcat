#!/usr/bin/env perl

$state = 0;
for (<>) {
    if ($state == 0) {
	if (/"(.*).utf8": {/) {
	    print "$1\n";
	    $state = 1;
	}
    } elsif ($state == 1) {
	/"(.*)"/;
	$t = $1;
	$t =~ s/\\x([0-9a-f]{2})/chr(hex($1))/ge;
	print "$t\n";
	if (/".*":\s*\d*}/) {
	    print "******\n";
	    $state = 0;
	}
    }
}
