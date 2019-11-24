#!/usr/bin/env bash

e=0

go get ./...

for d in $(go list ./... | grep -v vendor); do
    echo "running tests in $d"

    go test -v -race -cover $d > report.txt

    if [ $? -ne 0 ]
    then
      e=$?
    fi

    cat report.txt

    if [ -d /var/tests ]
    then
       cat report.txt | go-junit-report > /var/tests/${d//\//_}.xml
       cp report.txt /var/tests/${d//\//_}.txt
    fi

    rm report.txt
done

exit $e