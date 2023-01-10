# API Reference

### <a name="namespaces"></a> Namespaces

#### List registered namespaces

```
GET "v1/namespace"
```

```json
{
    "list": [
        {
            "id": "Infrastructure"
        }
    ]
}
```

#### Register new namespace

```
POST "v1/namespace"
Content-Type: "application/json"

{
  "id": "Infra"
}
```

#### Read registered namespace

```
GET "v1/namespace/{namespaceID}"
```

```json
{
	"id": "Infrastructure"
}
```

> When there's no such id, renders "204 No Content" response

#### Deregister namespace

```
DELETE "v1/namespace/{namespaceID}"
```

### <a name="tasks"></a> Tasks

#### List registered task IDs

```
GET "v1/task?namespace_id={namespaceID}"
```

```json
{
	"list": [
		{
			"id": "Backup_Test_Shards11-21",
			"namespace_id": "Infrastructure",
			"timeout": 360,
			"retries_left": 4,
			"updated_at": 1673362943,
			"status": "idle"
		},
		{
			"id": "Backup_Test_Shards1-10",
			"namespace_id": "Infrastructure",
			"timeout": 360,
			"retries_left": 3,
			"updated_at": 1673362984,
			"status": "failed"
		},
		{
			"id": "Backup_Test_Shards22-32",
			"namespace_id": "Infrastructure",
			"timeout": 360,
			"retries_left": 4,
			"updated_at": 1673363000,
			"status": "done"
		}
	]
}
```

#### Register new task

```
POST "v1/task"
Content-Type: "application/json"

{
  "id": "Backup_Test_Shards1-10",
  "namespace_id": "Infrastructure",
  "timeout": 360,
  "retries_left": 5,
  "status": "idle"
}
```

| Parameter    | Description          | Type    | Required |
|--------------|----------------------| ------- | -------- |
| id           | Task ID              | String  | True     |
| namespace_id | Namespace ID         | String  | True     |
| timeout      | Timeout (in seconds) | Integer | False    |
| retries_left | Retries left         | Integer | False    |
| status       | Status               | String  | False    |

#### Read registered task

```
GET "v1/task/{taskID}"
```

```json
{
	"id": "Backup_Test_Shards1-10",
	"namespace_id": "Infrastructure",
	"timeout": 360,
	"retries_left": 3,
	"updated_at": 1673362984,
	"status": "failed"
}
```

> When there's no such id, renders "204 No Content" response

#### Deregister task

```
DELETE "v1/task/{taskID}"
```

#### Capture tasks

```
GET "v1/task/capture?namespace_id={namespaceID}&limit={maxNumberOfTasks}"
```

```json
{
	"list": [
		{
			"id": "Backup_Test_Shards22-32",
			"namespace_id": "Infrastructure",
			"timeout": 360,
			"retries_left": 4,
			"updated_at": 1673362954,
			"status": "in_progress"
		},
		{
			"id": "Backup_Test_Shards1-10",
			"namespace_id": "Infrastructure",
			"timeout": 360,
			"retries_left": 3,
			"updated_at": 1673362954,
			"status": "in_progress"
		}
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
