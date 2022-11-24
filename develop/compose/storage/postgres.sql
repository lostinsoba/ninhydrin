-- tag
create table tag
(
    id text not null
);
alter table tag add constraint tag_unique_id unique (id);

-- pool
create table pool
(
    id          text   not null,
    description text   not null,
    tag_ids     text[] not null
);
alter table pool add constraint pool_unique_id unique (id);

-- task
create type task_status as enum ('idle', 'in_progress', 'timeout', 'failed', 'done');
create table task
(
    id           text        not null,
    pool_id      text        not null,
    timeout      interval    not null,
    retries_left smallint    not null,
    updated_at   timestamp   not null default (now() at time zone 'utc'),
    status       task_status not null
);
alter table task add constraint task_unique_id unique (id);

-- worker
create table worker
(
    id      text   not null,
    tag_ids text[] not null
);
alter table worker add constraint worker_unique_id unique (id);
