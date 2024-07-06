#!/bin/bash
eval $(awk '{print "export " $1}' ./configs/dev.yaml | sed 's/: /=/g')
