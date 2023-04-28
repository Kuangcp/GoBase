# http://192.168.16.151:12345/dolphinscheduler/projects/测试cdh/task-instance/list-paging?pageSize=10&pageNo=1&searchVal=&processInstanceId=&host=&stateType=&startDate=&endDate=&executorName=&processInstanceName=&processDefinitionId=&sort=&order= 

# curl -H "Content-Type:application/json" -X POST --data '{"id":17,"text":"http://192.168.16.151:12345/dolphinscheduler/projects/测试cdh/task-instance/list-paging?pageSize=10&pageNo=1&searchVal=&processInstanceId=&host=&stateType=&startDate=&endDate=&executorName=&processInstanceName=&processDefinitionId=&sort=&order=","document":{"title":"阿森松岛所445","number":223}}' http://127.0.0.1:8080/api/index?database=rrrr


curl -H "Content-Type:application/json" -X POST --data '{"query":"测试cdh","page":1,"limit":10,"order":"desc"}' http://127.0.0.1:8080/api/query?database=rrrr

