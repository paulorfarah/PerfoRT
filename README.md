# PerfoRT: Performance Regression Tool

PerfoRT is a tool that automates the measurement of performance regression in Java projects and helps developers to mine performance regressions of software repositories. PerfoRT allows developers to automatically extract runtime performance information from Java projects, such as the number of
calls and time duration of versions, packages, classes, and methods.  it provides information related to testing code coverage metrics, process, and system utilization behavior, as well as to Java Virtual Machine (JVM) events. 


## Installation instructions


###### requirements:
- java
- golang
- maven
- mysql

###### mysql configurations:
https://www.digitalocean.com/community/tutorials/how-to-move-a-mysql-data-directory-to-a-new-location-on-ubuntu-16-04

ALTER USER 'root'@'localhost' IDENTIFIED WITH mysql_native_password by 'password';

sudo vim /etc/mysql/mysql.conf.d/mysqld.cnf

max_connections = 9999
sql-mode="ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION"
key_buffer_size=0
innodb_buffer_pool_size = 5G
innodb_stats_on_metadata = 0


show variables like 'max_connections';
show variables like 'sql_mode';
show variables like 'innodb_buffer_pool_size';
show variables like 'innodb_log_file_size';

###### Jacoco:
- $ wget https://search.maven.org/remotecontent?filepath=org/jacoco/jacoco/0.8.6/jacoco-0.8.6.zip
- $ unzip jacoco-0.8.6.zip /path/to/PerfoRT
  
###### environment variables:
~/.profile:
export PATH=$PATH:/mnt/sda4/go/bin:/mnt/sda4/apache-maven-3.8.6/bin
export MAVEN_OPTS="-javaagent:/<path>/<to>/PerfoRT/jacoco-0.8.6/jacocoagent.jar"
export JAVA_HOME=/usr/lib/jvm/java-1.11.0-openjdk-amd64

 ## Perfort settings:
  
1. .env
  This file contains settings to run PerfoRT.
  
2. .versions
  This folder stores files that contains lists of hashes of commits of the target system to be measured by PerfoRT. The file should have the name of the package name configured in the .env file.
For example, for the target system apache commons-bcel, the package name is "org.apache.bcel.", so the releases file should be named as ".releases_org.apache.bcel.".

3. .tcignore
There are some testcases that runs without ending. To deal with this situation add the testcase in a .tcignore file. A .tcignore file is a list of testcases to be ignored during the PerfoRTance measurement of the target system. The .tcignore file should have the name with the package name of the project. For example, for the target system apache commons-bcel, the package name is "org.apache.bcel.", so the testcase ignore list file should be named as "org.apache.bcel.".

If the testcase to be ignored is not listed in the ignore file, PerfoRT will stablish a timeout, also can be configured in the testcase_timeout of the .env file. However, ignore them is better because do not take extra time and use extra resources neither.

#### How to contribute?
  Just submit a PR.
  
#### How to cite?
  Please, use the following entry to cite PerfoRT:
```bibtex
@InProceedings{farahvergilio2023perfort,,
   title      = {PerfoRT: A Tool for Software Performance Regression},
   author     = {Paulo Roberto Farah and Silvia Regina Vergilio},
   booktitle    = {ACM/SPEC International Conference on Performance Engineering},
   year       = 2023,
   note       = {Available on https://github.com/paulorfarah/perfort}
 }
```

 
  

