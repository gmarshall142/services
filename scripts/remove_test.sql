-- users
delete from users where lastname = 'Victor';
delete from users where lastname = 'Martin Luther';

-- bikes
delete from bikes where name like 'certs%';

-- reset sequence value
select setval('users_id_seq', coalesce(max(id), 1)) from users;
select setval('bike_id_seq', coalesce(max(id), 1)) from bikes;
select setval('bikerims_id_seq', coalesce(max(id), 1)) from bikerims;
