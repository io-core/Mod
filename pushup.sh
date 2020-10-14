#!/bin/bash
git add *.Mod
git add *.Pkg
git add README.md
git commit -m 'sync local to upstream'
git push origin
