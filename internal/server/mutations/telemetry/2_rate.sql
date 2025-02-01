create table if not exists "rates"
(
    queue_id     text not null,
    metric_name  text not null,
    metric_value real not null,
    timestamp    real not null,
    labels       text not null
);

create index if not exists queue_index
    on rates (queue_id);

create index if not exists metric_index
    on rates (metric_name);

create index if not exists timestamp_index
    on rates (timestamp);