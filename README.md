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
max_connections

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
