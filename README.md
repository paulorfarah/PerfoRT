# perfrt  - performance regression testing tool

1) install maven and/or gradle (depend on the projects will be analyzed):
- $ wget https://dlcdn.apache.org/maven/maven-3/3.8.4/binaries/apache-maven-3.8.4-bin.tar.gz
- $ tar -xzvf apache-maven-3.8.4-bin.tar.gz
- $ export PATH=/path/to/apache-maven-3.8.4/bin:$PATH

installation

golang
mysql
maven
java

mysql configurations
https://www.digitalocean.com/community/tutorials/how-to-move-a-mysql-data-directory-to-a-new-location-on-ubuntu-16-04


show variables like 'sql_mode';
SET sql_mode = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION';

max_connections
set global max_connections = 999999;

2) download jacoco
- $ wget https://search.maven.org/remotecontent?filepath=org/jacoco/jacoco/0.8.6/jacoco-0.8.6.zip
- $ unzip jacoco-0.8.6.zip /path/to/perfrt

3) download and configure async-profiler: 
- $ wget https://github.com/jvm-profiling-tools/async-profiler/releases/download/v2.6/async-profiler-2.6-linux-x64.tar.gz
- $ tar -xzvf async-profiler-2.6-linux-x64.tar.gz 
- $ sudo apt install openjdk-11-dbg (or openjdk-8-dbg)
- $ sudo sysctl kernel.perf_event_paranoid=1
- $ sudo sysctl kernel.kptr_restrict=0
  
4) configure environment variable MAVEN_OPTS:
- export MAVEN_OPTS=-agentpath:path/to/async-profiler-2.5.1-linux-x64/build/libasyncProfiler.so=start,event=wall,file=profile.txt

Setting Files
1) .env


2) .releases
Contains the list of hashes of commits of the target system to be measured by perfrt. The .releases file should have the name .releases plus character "_" (underline) and the package name configured in the .env file.
For example, for the target system apache commons-bcel, the package name is "org.apache.bcel.", so the releases file should be named as ".releases_org.apache.bcel.".

3) .tcignore
There are some testcases that runs without ending. To deal with this situation add the testcase in a .tcignore file. A .tcignore file is a list of testcases to be ignored during the performance measurement of the target system. The .tcignore file should have the name .tcignore plus character "_" (underline) and the package name configured in the .env file.
For example, for the target system apache commons-bcel, the package name is "org.apache.bcel.", so the testcase ignore list file should be named as ".tcignore_org.apache.bcel.".

If the testcase to be ignored is not listed in the ignore file, perfrt will stablish a timeout, also can be configured in the testcase_timeout of the .env file. However, ignore them is better because do not take extra time and use extra resources neither.