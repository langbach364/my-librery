#!/bin/bash

commandTool="commandTool"
fileName="nameSocket.sock"

$commandTool && {
    echo "Công cụ đã xử lý xong"
    echo "SUCCESS" >$fileName
    rm -f $fileName
    exit 0
} || {
    echo "Công cụ xử lý thất bại"
    echo "FAIL" >$fileName
    rm -f $fileName
    exit 1
}