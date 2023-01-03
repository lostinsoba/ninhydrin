# API Reference

### <a name="tasks"></a> Tasks

#### List registered task IDs

```
GET "v1/task"
```

```
{
    "list": [
        "Backup_Test_Shards1-10",
        "Backup_Test_Shards11-21",
        "Backup_Test_Shards22-32"
    ]
}
```

#### Register new task

```
POST "v1/task"
Content-Type: "application/json"

{
  "id": "Backup_Test_Shards1-10",
  "timeout": 360,
  "retries_left": 5,
  "status": "idle"
}
```

| Parameter    | Description          | Type    | Required |
| ------------ | -------------------- | ------- | -------- |
| id           | Task ID              | String  | True     |
| timeout      | Timeout (in seconds) | Integer | False    |
| retries_left | Retries left         | Integer | False    |
| status       | Status               | String  | False    |

#### Read registered task

```
GET "v1/task/{taskID}"
```

```
{
    "id": "Backup_Test_Shards1-10",
    "timeout": 360,
    "retries_left": 4,
    "updated_at": 1672758864,
    "status": "in_progress"
}
```

#### Delete registered task

```
DELETE "v1/task/{taskID}"
```

#### Capture tasks

```
GET "v1/task/capture?limit=n"
```

```
{
    "list": [
        "Backup_Test_Shards11-21",
        "Backup_Test_Shards22-32"
    ]
}
```

#### Release tasks

```
PUT "v1/task/release"
Content-Type: "application/json"

{
  "status": "done",
  "task_ids": ["Backup_Test_Shards1-10"]
}
```

| Parameter | Description      | Type            | Required |
| --------- | ---------------- | --------------- | -------- |
| status    | Status           | String          | True     |
| task_ids  | List of task IDs | List of strings | True     |
