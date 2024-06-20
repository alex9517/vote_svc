#!/bin/bash
#  Created : 2023-Apr-05
# Modified : 2024-Apr-05

echo "UNIT TESTS"
echo "=========="
echo "pkg/service"
echo "-----------"
echo "Service layer tests would not work without access to the database instance!"
read -p "Continue?(Y/n): " ANS
if [ "$ANS" == "N" ] || [ "$ANS" == "n" ];
then
  exit 0
fi
cd ./pkg/service
go test -v

echo ""
echo "pkg/endpoint"
echo "------------"
read -p "Continue?(Y/n): " ANS
if [ "$ANS" == "N" ] || [ "$ANS" == "n" ];
then
  exit 0
fi
cd ../endpoint
go test -v

echo ""
echo "pkg/http"
echo "--------"
read -p "Continue?(Y/n): " ANS
if [ "$ANS" == "N" ] || [ "$ANS" == "n" ];
then
  exit 0
fi
cd ../http
go test -v

# --- END ---
