# API Reference

1. [User API](#user-api)  
1.1 [Tags](#user-api-tags)  
1.2 [Pools](#user-api-pools)  
1.3 [Tasks](#user-api-tasks)  
1.4 [Workers](#user-api-workers)  

2. [Worker API](#worker-api)  
2.1 [Tasks](#worker-api-tasks)

## <a name="user-api"></a> User API

### <a name="user-api-tags"></a> Tags

#### List registered tag IDs

```
GET "v1/user/tag"
```

#### Register new tag

```
POST "v1/user/tag"
Content-Type: "application/json"

{
  "id": "Backup_Download"
}
```

| Parameter | Description | Type   | Required |
| --------- | ----------- | ------ | -------- |
| id        | Tag ID      | String | True     |

#### Read registered tag

```
GET "v1/user/tag/{tagID}"
```

#### Delete registered tag

```
DELETE "v1/user/tag/{tagID}"
```

### <a name="user-api-pools"></a> Pools

#### List registered pool IDs

```
GET "v1/user/pool"
```

#### Update or register new pool

##### Register new pool

```
POST "v1/user/pool"
Content-Type: "application/json"

{
  "id": "Backup",
  "description": "Download and test latest backups",
  "tag_ids": ["Backup_Download"]
}
```

##### Update registered pool

```
PUT "v1/user/pool/{poolID}"
Content-Type: "application/json"

{
  "id": "Backup",
  "description": "Download and test latest backups",
  "tag_ids": ["Backup_Download", "Backup_Test"]
}
```

| Parameter   | Description | Type         | Required |
| ----------- | ----------- | ------------ | -------- |
| id          | Pool ID     | String       | True     |
| description | Description | String       | False    |
| tag_ids     | Tag IDs     | String array | True     |

#### Read registered pool

```
GET "v1/user/pool/{poolID}"
```

#### Delete registered pool

```
DELETE "v1/user/pool/{poolID}"
```

### <a name="user-api-tasks"></a> Tasks

#### List registered task IDs

```
GET "v1/user/task"
```

#### Register new task

```
POST "v1/user/task"
Content-Type: "application/json"

{
  "id": "Backup_Download_Daily",
  "pool_id": "Backup",
  "timeout": 360,
  "retries_left": 5,
  "status": "idle"
}
```

| Parameter    | Description          | Type    | Required |
| ------------ | -------------------- | ------- | -------- |
| id           | Task ID              | String  | True     |
| pool_id      | Pool ID              | String  | True     |
| timeout      | Timeout (in seconds) | Integer | False    |
| retries_left | Retries left         | Integer | False    |
| status       | Status               | String  | False    |

#### Read registered task

```
GET "v1/user/task/{taskID}"
```

#### Delete registered task

```
DELETE "v1/user/task/{taskID}"
```

### <a name="user-api-workers"></a> Workers

#### List registered worker IDs

```
GET "v1/user/worker"
```

#### Register new worker

```
POST "v1/user/worker"
Content-Type: "application/json"

{
  "id": "BackupWorker1",
  "tag_ids": ["Backup_Download"]
}
```

| Parameter    | Description | Type          | Required |
| ------------ | ----------- | ------------- | -------- |
| id           | Worker ID   | String        | True     |
| tag_ids      | Tag IDs     | String array  | True     |

#### Read registered worker

```
GET "v1/user/worker/{workerID}"
```

#### Delete registered worker

```
DELETE "v1/user/worker/{workerID}"
```

## <a name="worker-api"></a> Worker API

### <a name="worker-api-tasks"></a> Tasks

#### Capture tasks

```
GET "v1/worker/task/capture?limit={taskLimit}"
X-Ninhydrin-Worker-Token: {workerToken}
```

#### Change task status

```
POST "v1/worker/task/{taskID}/status"
X-Ninhydrin-Worker-Token: {workerToken}
Content-Type: "application/json"

{
  "status": "done"
}
```

| Parameter | Description | Type    | Required |
| --------- | ----------- | ------- | -------- |
| status    | Task status | String  | True     |
