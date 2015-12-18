#!/bin/bash
docker run -w `pwd` -v `pwd`:`pwd` --net=host -it --rm mysql/mysql-server sh -c 'exec mysql -h 0.0.0.0 -uroot -proot -e "GRANT ALL PRIVILEGES ON * . * TO msvcdir;"'
docker run -w `pwd` -v `pwd`:`pwd` --net=host -it --rm mysql/mysql-server sh -c 'exec mysql -h 0.0.0.0 -umsvcdir -pmsvcdir -e "create database microservicesdirtest;"'
docker run -w `pwd` -v `pwd`:`pwd` --net=host -it --rm mysql/mysql-server sh -c 'exec mysql -h 0.0.0.0 -umsvcdir -pmsvcdir microservicesdirtest < sql/create_schema.sql'
