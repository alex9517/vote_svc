#!/bin/bash
#  Created : 2024-Apr-22
# Modified : 2024-May-14
# To start a GUI version of JMeter, use '$HOME/bin/jmeter.sh'

JMETER=/home/alex/Applications/apache-jmeter-5.6.3/bin/jmeter

# $JMETER -n -t ./VoteSvcGetLoad-20.jmx -l ./log-20.jtl -e -o ./votes-get-20/

# $JMETER -n -t ./VoteSvcGetLoad-100.jmx -l ./log-100.jtl -e -o ./votes-get-100/

# $JMETER -n -t ./VoteSvcGetLoad-1000.jmx -l ./log-1000.jtl -e -o ./votes-get-1000/

# $JMETER -n -t ./VoteSvcGetLoad-10000.jmx -l ./log-10000.jtl -e -o ./votes-get-10000/

# $JMETER -n -t ./VoteSvcGetResultsLoad-20.jmx -l ./log-20.jtl -e -o ./votes-get-res-20/

# $JMETER -n -t ./VoteSvcGetResultsLoad-100.jmx -l ./log-100.jtl -e -o ./votes-get-res-100/

# $JMETER -n -t ./VoteSvcGetResultsLoad-1000.jmx -l ./log-1000.jtl -e -o ./votes-get-res-1000/

# $JMETER -n -t ./VoteSvcGetResultsLoad-10000.jmx -l ./log-10000.jtl -e -o ./votes-get-res-10000/

# $JMETER -n -t ./VoteSvcUpdateLoad-100.jmx -l ./log-100.jtl -e -o ./votes-put-res-100/

$JMETER -n -t ./VoteSvcUpdateLoad-1000.jmx -l ./log-1000.jtl -e -o ./votes-put-res-1000/
