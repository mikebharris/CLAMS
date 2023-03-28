set timezone = 'UTC';

create schema if not exists "hacktionlab_workshops";
set schema 'hacktionlab_workshops';

create table hacktionlab_workshops.workshops
(
    id    serial primary key,
    title text unique not null
);

create table hacktionlab_workshops.people
(
    id       serial primary key,
    forename varchar(50)         null,
    surname  varchar(50)         null,
    email    varchar(100) unique not null
);

create table hacktionlab_workshops.roles
(
    id        serial primary key,
    role_name varchar(20) null
);

create table hacktionlab_workshops.workshop_signups
(
    id           serial primary key,
    people_id    int       not null,
    workshop_id  int       not null,
    role_id      int       not null,
    signed_up_on timestamp not null,

    foreign key (people_id)
        references people (id),
    foreign key (workshop_id)
        references workshops (id),
    foreign key (role_id)
        references roles (id)
);

create sequence trigger_notifications_id_seq as bigint start 1 increment 1 cache 1;
create table hacktionlab_workshops.trigger_notifications
(
    id      bigint       not null default nextval('trigger_notifications_id_seq'),
    message varchar(256) not null,
    primary key (id)
);

create or replace function fn_trigger() returns trigger as
$fn_trigger$
begin
    insert into hacktionlab_workshops.trigger_notifications (message)
    values ('{' ||
            '"WorkshopSignupId": ' || coalesce(new.id, old.id) || ',' ||
            '"WorkshopId": ' || coalesce(new.workshop_id, old.workshop_id) || ',' ||
            '"PeopleId": ' || coalesce(new.people_id, old.people_id) || ',' ||
            '"RoleId": ' || coalesce(new.role_id, old.role_id) ||
            '}');
    return null;
end;
$fn_trigger$ language plpgsql;

create trigger signups_trigger
    after insert or update
    on hacktionlab_workshops.workshop_signups
    for each row
execute procedure fn_trigger();