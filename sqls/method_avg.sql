USE perfrt;
SELECT repo.name, mea.id AS measurement, c.committer_date, commit_hash, test.name AS test_name, r.id AS run, f.name AS class_name, m.name AS method_name, m.created_at AS method_started_at, m.ended_at AS method_ended_at, m.caller_id, m.own_duration, m.cumulative_duration, 
AVG(res.cpu_percent), STD(res.cpu_percent)
FROM commits AS c
INNER JOIN files AS f ON f.commit_id=c.id
INNER JOIN methods AS m ON m.file_id=f.id
INNER JOIN runs AS r ON m.run_id = r.id
INNER JOIN measurements AS mea On r.measurement_id=mea.id
INNER JOIN repositories AS repo ON mea.repository_id = repo.id
INNER JOIN testcases AS test ON test.id=r.test_case_id
INNER JOIN jvms AS jvm ON jvm.run_id = r.id
INNER JOIN resources AS res ON res.run_id = r.id
WHERE m.finished=true
AND res.timestamp >= (SELECT MIN(created_at) FROM methods AS met WHERE met.run_id=res.run_id GROUP BY met.run_id)
AND res.timestamp <= (SELECT MAX(ended_at) FROM methods AS met WHERE met.run_id=res.run_id GROUP BY met.run_id)     
GROUP BY repo.id, c.commit_hash, f.name, m.id, r.id, res.run_id     
ORDER BY repo.name, mea.id, c.committer_date, test.name, f.name, m.name, r.id;