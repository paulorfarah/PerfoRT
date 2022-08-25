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
SET sql_mode = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION';

my.cnf:
sql-mode="ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION"

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


FAQ
resources is full: 
du -hx --max-depth=1 .
=======
Setting Files
1) .env


2) .releases
Contains the list of hashes of commits of the target system to be measured by perfrt. The .releases file should have the name .releases plus character "_" (underline) and the package name configured in the .env file.
For example, for the target system apache commons-bcel, the package name is "org.apache.bcel.", so the releases file should be named as ".releases_org.apache.bcel.".

3) .tcignore
There are some testcases that runs without ending. To deal with this situation add the testcase in a .tcignore file. A .tcignore file is a list of testcases to be ignored during the performance measurement of the target system. The .tcignore file should have the name .tcignore plus character "_" (underline) and the package name configured in the .env file.
For example, for the target system apache commons-bcel, the package name is "org.apache.bcel.", so the testcase ignore list file should be named as ".tcignore_org.apache.bcel.".

If the testcase to be ignored is not listed in the ignore file, perfrt will stablish a timeout, also can be configured in the testcase_timeout of the .env file. However, ignore them is better because do not take extra time and use extra resources neither.

SELECT c.committer_date, commit_hash, r.id, f.name AS classname, m.name AS methodName, jvm.object_allocation_in_new_tlab_tlab_size   FROM perfrt.commits AS c INNER JOIN perfrt.files AS f ON f.commit_id=c.id INNER JOIN perfrt.methods AS m ON m.file_id=f.id INNER JOIN perfrt.runs AS r ON m.run_id = r.id INNER JOIN perfrt.jvms jvm ON jvm.run_id = r.id ORDER BY c.committer_date, f.name; 



INTO OUTFILE '/var/lib/mysql-files/jvm2.csv'
          FIELDS ENCLOSED BY '"'
          TERMINATED BY ';'
          ESCAPED BY '"'
          LINES TERMINATED BY '\r\n';

          

mysql perfrt -u root -p  < jvms.sql > openfire.tsv



SELECT repo.name, mea.id AS measurement, c.committer_date, commit_hash, r.id AS run, f.name AS classname, m.name AS methodName, m.own_duration, 
AVG(res.cpu_percent), STD(res.cpu_percent)
FROM commits AS c
INNER JOIN files AS f ON f.commit_id=c.id
INNER JOIN methods AS m ON m.file_id=f.id
INNER JOIN runs AS r ON m.run_id = r.id
INNER JOIN measurements AS mea On r.measurement_id=mea.id
INNER JOIN repositories AS repo ON mea.repository_id = repo.id
INNER JOIN jvms AS jvm ON jvm.run_id = r.id
INNER JOIN resources AS res ON res.run_id = r.id
WHERE res.timestamp >= (SELECT MIN(created_at) FROM methods AS met WHERE met.run_id=res.run_id GROUP BY met.run_id)
AND res.timestamp <= (SELECT MAX(ended_at) FROM methods AS met WHERE met.run_id=res.run_id GROUP BY met.run_id)     
GROUP BY repo.id, c.commit_hash, f.name, m.id, r.id, res.run_id     
ORDER BY repo.name, mea.id, c.committer_date, f.name, m.name, r.id;


-- INTO OUTFILE '/var/lib/mysql-files/jvm2.csv'
-- FIELDS ENCLOSED BY '"'
-- TERMINATED BY ';'
-- ESCAPED BY '"'
-- LINES TERMINATED BY '\r\n';


# read testcases:
mvn clean verify > maven_build.out
echo $(( $(cat maven_build.out | grep "Tests run" | grep -v "Time elapsed" | cut -d , -f 1 | cut -d " " -f 4 | tr "\n" "+") 0))
