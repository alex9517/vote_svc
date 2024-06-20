#!/bin/bash
#  Created : 2023-Feb-01
# Modified : 2024-Apr-04

echo "UNIT TESTS"
echo "=========="
echo "pkg/service"
echo "-----------"
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

echo ""
echo "pkg/grpc"
echo "--------"
read -p "Continue?(Y/n): " ANS
if [ "$ANS" == "N" ] || [ "$ANS" == "n" ];
then
  exit 0
fi
cd ../grpc
go test -v

cd -

# --- END ---
