#!/bin/bash

docker run -w `pwd` -v `pwd`:`pwd` --net=host -it --rm mysql/mysql-server sh -c 'exec mysql -h localhost -umsvcdir -pmsvcdir microservicesdirtest'
