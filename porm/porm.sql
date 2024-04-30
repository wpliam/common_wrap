create database if not exists `orm_test`;

drop table if exists `t_test_data`;
create table if not exists `t_test_data`
(
    `id`          bigint        not null auto_increment comment '',
    `name`        varchar(100)  not null default '',
    `status`      int           not null default 0,
    `enable`      boolean       not null default true,
    `like`        int           not null default 0,
    `content`     text,
    `score`       decimal(3, 2) not null default 0,
    `create_time` datetime               default current_timestamp(),
    `update_time` datetime               default current_timestamp() on update current_timestamp(),
    primary key (`id`) USING BTREE
) engine = InnoDB
  CHARACTER SET = utf8mb4