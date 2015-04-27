#!/bin/bash

sudo service ghostrunner stop

if [ -z "$1" ] 
then 
	echo "A task runner api key is required"
	exit
fi

if [ -z "$2" ] 
then 
	echo "A task runner secret key is required"
	exit
fi

if [ -f /etc/ghostrunner.conf ]
then
	rm /etc/ghostrunner.conf
fi

cp ghostrunner.conf /etc
sudo chmod 777 /etc/ghostrunner.conf

sed -i "s/\"ApiKey\":\"\"/\"ApiKey\":\"$1\"/g" /etc/ghostrunner.conf
sed -i "s/\"ApiSecret\":\"\"/\"ApiSecret\":\"$2\"/g" /etc/ghostrunner.conf

if [ -d /etc/ghostrunner.proc ]
then
	rm -rf /etc/ghostrunner.proc
fi

mkdir /etc/ghostrunner.proc
sudo chmod 777 /etc/ghostrunner.proc

if [ -d /var/log/ghostrunner ]
then
	rm -rf /var/log/ghostrunner
fi

mkdir /var/log/ghostrunner
sudo chmod 777 /var/log/ghostrunner

if [ -f /usr/bin/ghostrunner ]
then
	rm /usr/bin/ghostrunner
fi

cp ghostrunner /usr/bin
sudo chmod 777 /usr/bin/ghostrunner

if [ -f /etc/init.d/ghostrunner ]
then
	rm /etc/init.d/ghostrunner
fi

sudo chmod 777 init.d/ghostrunner
cp init.d/ghostrunner /etc/init.d

sudo update-rc.d ghostrunner defaults
sudo service ghostrunner start