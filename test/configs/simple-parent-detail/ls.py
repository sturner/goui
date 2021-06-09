#!/usr/bin/env python
import os, json, sys
args = sys.argv[1:]
ls_dir = "."
if len(args) > 0:
    ls_dir = args[0]
print json.dumps(os.listdir(ls_dir))



