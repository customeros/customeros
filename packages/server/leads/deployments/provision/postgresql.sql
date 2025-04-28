-- auto-generated definition
create table if not exists app_keys
(
  id     bigserial
    primary key,
  app_id varchar(255) not null,
  key    varchar(255) not null,
  active boolean      not null
);

ALTER TABLE app_keys
    owner TO postgres;

