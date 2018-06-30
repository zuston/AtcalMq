# 迁移脚本,将日志文件迁移到 hdfs 中
TIME=`date -d "yesterday" +%Y-%m-%d`
echo $TIME
createCommand="/usr/hadoop/bin/hdfs dfs -mkdir /mqBackup/"${TIME}
$createCommand
putCommand="/usr/hadoop/bin/hdfs dfs -put /opt/aneBackup/*"${TIME}" /mqBackup/"${TIME}
$putCommand
