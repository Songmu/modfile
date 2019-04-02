#!/bin/sh

perl -MFile::Basename=basename -i -pe \
    's{"((?:cmd/go/)?internal/.+?)"}{q{"github.com/Songmu/modfile/internal/}.basename($1).q{"}}eg' \
    $(find . -name '*.go')
