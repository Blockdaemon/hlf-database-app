#!/usr/bin/env python3
import os
import sys
import jinja2

sys.stdout.write(jinja2.Template(sys.stdin.read()).render(env=os.environ) + "\n")
