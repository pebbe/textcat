#!/usr/bin/env python3

# usage: python3 ./mkData.py /path/to/libexttextcat/*.lm | gofmt > data.go

import os.path
import re
import sys

reMore = re.compile('[ \t\n\r].*')

sys.stdout.buffer.write(b'''// THIS IS A GENERATED FILE
// DO NOT EDIT

package textcat

var data = map[string]map[string]int{
''')

for filename in sys.argv[1:]:
    lineno = 0
    with open(filename, "rt", encoding='utf-8') as fp:
        fname = os.path.basename(filename).replace('.lm', '')
        sys.stdout.buffer.write('"{}.utf8": {{\n'.format(fname).encode('utf8'))
        for line in fp:
            line = reMore.sub('', line).replace('"', '\\"')
            sys.stdout.buffer.write('"{}": {},\n'.format(line, lineno).encode('utf8'))
            lineno += 1
        sys.stdout.buffer.write(b'},\n')

sys.stdout.buffer.write(b'}\n')
