#!/bin/bash
docker build -t clsung/csigo_build build/
docker run -d --name csigo_build clsung/csigo_build
docker cp csigo_build:/go/bin/csigosrv run/
docker stop csigo_build
docker rm csigo_build
docker rmi clsung/csigo_build
docker build -t clsung/csigosrv run/
