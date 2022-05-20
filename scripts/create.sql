drop table if exists users cascade;
create table users
(
    id         serial primary key,
    is_admin   boolean        not null,
    first_name varchar        not null,
    last_name  varchar        not null,
    username   varchar unique not null,
    email      varchar unique not null,
    password   varchar        not null,
    activated  bool           not null
);

drop table if exists tests cascade;
create table tests
(
    id            serial primary key,
    last_modified timestamp with time zone                             not null,
    final         bool                                                 not null,
    name          varchar                                              not null,
    public        boolean                                              not null,
    user_id       int default 0 references users on delete set default not null,
    task_id       int references tasks on delete cascade               not null,
    language      varchar                                              not null,
    code          varchar                                              not null
);

drop table if exists tasks cascade;
create table tasks
(
    id           serial primary key,
    author_id    int default 0 references users on delete set default not null,
    approver_id  int default 0 references users on delete set default not null,
    title        varchar                                              not null,
    difficulty   varchar                                              not null,
    is_published boolean                                              not null,
    added_on     timestamp with time zone                             not null,
    text         varchar                                              not null
);

drop table if exists user_solutions cascade;
create table user_solutions
(
    id            serial primary key,
    user_id       int default 0 references users on delete set default not null,
    task_id       int references tasks on delete cascade               not null,
    last_modified timestamp with time zone                             not null,
    language      varchar                                              not null,
    name          varchar                                              not null,
    public        boolean                                              not null,
    code          varchar                                              not null
);

drop table if exists user_solutions_tests cascade;
create table user_solutions_tests
(
    user_solution_id int references user_solutions on delete cascade not null,
    test_id          int references tests on delete cascade          not null,
    user_id          int references users on delete cascade          not null
);

drop table if exists user_solutions_results cascade;
create table user_solutions_results
(
    user_solution_id int references user_solutions on delete cascade not null,
    test_id          int references tests on delete cascade          not null,
    user_id          int references users on delete cascade          not null,
    exit_code        int                                             not null,
    output           varchar                                         not null,
    compilation_time float4                                          not null,
    real_time        float4                                          not null,
    kernel_time      float4                                          not null,
    user_time        float4                                          not null,
    max_ram_usage    float4                                          not null,
    binary_size      float4                                          not null
);

drop table if exists last_opened cascade;
create table last_opened
(
    user_id                         int references users on delete cascade          not null,
    task_id                         int references tasks on delete cascade          not null,
    user_solution_id_for_language_1 int references user_solutions on delete cascade not null,
    user_solution_id_for_language_2 int references user_solutions on delete cascade not null,
    language_1                      varchar                                         not null,
    language_2                      varchar                                         not null
);

drop table if exists tokens_for_password_reset cascade;
create table tokens_for_password_reset
(
    user_id int references users on delete cascade not null,
    token   varchar                                not null
);

drop table if exists tokens_for_registration cascade;
create table tokens_for_registration
(
    user_id int references users on delete cascade not null,
    token   varchar                                not null
);