-- namespace
create table namespace
(
    id text not null
);
alter table namespace add constraint namespace_unique_id unique (id);

-- task
create table task
(
    id           text        not null,
    namespace_id text        not null,
    timeout      integer     not null
);
alter table task add constraint task_unique_id unique (id);

-- task state
create type task_status as enum ('timeout', 'failed', 'idle', 'in_progress', 'done');
create table task_state
(
	task_id      text        not null,
	retries_left smallint    not null,
	updated_at   bigint      not null,
	status       task_status not null
);
alter table task_state add constraint task_unique_id unique (task_id);
