#!/bin/sh

upgrade_pre() {
  localDir=$1
  imgName=$(echo "${localDir}" | sed 's/\//\-/g')
  echo "${imgName}"
  sudo docker compose -f docker-compose-pre.yml pull "${imgName}"
  sudo docker compose -f docker-compose-pre.yml stop "${imgName}"
  sudo docker compose -f docker-compose-pre.yml rm -f "${imgName}"
  sudo docker compose -f docker-compose-pre.yml up -d "${imgName}"

  return
}