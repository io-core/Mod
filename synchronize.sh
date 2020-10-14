#!/bin/bash
git add *.Mod
git add *.Pkg
git add *.Tool
git add README.md
git commit -m 'sync local to master'
git push origin main
