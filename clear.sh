# 清理脚本
TIME=`date -d "yesterday" +%Y-%m-%d`
echo $TIME
clearCommand="rm /opt/aneBackup/*.log."${TIME}
$clearCommand