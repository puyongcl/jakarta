#!/bin/sh

upgrade_test() {
  localDir=$1
  imgName=$(echo "${localDir}" | sed 's/\//\-/g')
  echo "${imgName}"
  sudo docker compose -f docker-compose-test.yml pull "${imgName}"
  sudo docker compose -f docker-compose-test.yml stop "${imgName}"
  sudo docker compose -f docker-compose-test.yml rm -f "${imgName}"
  sudo docker compose -f docker-compose-test.yml up -d "${imgName}"
}