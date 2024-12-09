#!/bin/bash

START_SLOT=4700014
END_SLOT=10579133

LEFT_SLOT_1=$((START_SLOT + 1))
LEFT_SLOT_2=$((START_SLOT + 2))

RIGHT_SLOT_1=$((END_SLOT - 1))
RIGHT_SLOT_2=$((END_SLOT - 2))

MIDDLE_SLOT=$(( (START_SLOT + END_SLOT) / 2 ))

RANDOM_SLOT=$((MIDDLE_SLOT + RANDOM % 1000 - 500))
BASE_URL1="http://localhost:8000/blockreward"
BASE_URL2="http://localhost:8000/syncduties"

TEST_SLOTS=($LEFT_SLOT_1 $LEFT_SLOT_2 $RIGHT_SLOT_1 $RIGHT_SLOT_2 $MIDDLE_SLOT $RANDOM_SLOT)

for SLOT in "${TEST_SLOTS[@]}"; do
  echo "Testing slot $SLOT..."
  RESPONSE1=$(curl -s -X 'GET' "${BASE_URL1}/${SLOT}" -H 'accept: application/json')
  RESPONSE2=$(curl -s -X 'GET' "${BASE_URL2}/${SLOT}" -H 'accept: application/json')

  if [[ $? -eq 0 ]]; then
    echo "Slot $SLOT: $RESPONSE1"
    echo "Slot $SLOT: $RESPONSE2"
  else
    echo "Slot $SLOT: Error occurred while fetching data."
  fi
done
