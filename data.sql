create table jobs(job_id integer PRIMARY KEY, job_title varchar(35), min_salary integer, max_salary integer);
insert into jobs (job_id, job_title, min_salary, max_salary) values ( 1, 'pm', '100000', '150000');
insert into jobs (job_id, job_title, min_salary, max_salary) values ( 2, 'developer', '80000', '120000');
insert into jobs (job_id, job_title, min_salary, max_salary) values ( 3, 'tester', '60000', '100000');

create table locations (location_id integer PRIMARY KEY, street_address varchar(25), postal_code varchar(12), city varchar(30), state_province varchar(12), country_id integer);
insert into locations (location_id, street_address, postal_code, city, state_province, country_id) values ( 1, 'Pushkin str', '140105', 'Moscow', 'Moscow reg', 1);
insert into locations (location_id, street_address, postal_code, city, state_province, country_id) values ( 2, 'Kalatushkin str', '200100', 'Dubai', 'Dubai reg', 2);
insert into locations (location_id, street_address, postal_code, city, state_province, country_id) values ( 3, 'Wall str', '140105', 'New York', 'New York reg', 3);

create table departments(department_id integer PRIMARY KEY, department_name varchar(20), manager_id integer, location_id integer REFERENCES locations (location_id));
insert into departments (department_id, department_name, manager_id, location_id) values ( 1, 'IT department 001', 1, 1);
insert into departments (department_id, department_name, manager_id, location_id) values ( 2, 'IT department 002', 2, 2);
insert into departments (department_id, department_name, manager_id, location_id) values ( 3, 'IT department 003', 3, 3);

create table employees(employee_id integer PRIMARY KEY, first_name varchar(20), last_name varchar(25), email varchar(25), phone_number varchar(20), hire_date date, job_id integer REFERENCES jobs (job_id), salary integer, commission_pct integer, manager_id integer, department_id integer REFERENCES departments (department_id));
insert into employees(employee_id, first_name, last_name, email, phone_number, hire_date, job_id, salary, commission_pct, manager_id, department_id) values ( 1, 'Jon', 'Snow', 'jon@mail.ru', 88005555535, '2017-08-20', 1, 150000, 13, 1, 1);
insert into employees(employee_id, first_name, last_name, email, phone_number, hire_date, job_id, salary, commission_pct, manager_id, department_id) values ( 2, 'Tyrion', 'Lannister', 'tyrion@mail.ru', 88005555535, '2017-07-20', 2, 120000, 13, 1, 2);
insert into employees(employee_id, first_name, last_name, email, phone_number, hire_date, job_id, salary, commission_pct, manager_id, department_id) values ( 3, 'Daenerys', 'Targaryen', 'daenerys@mail.ru', 88005555535, '2017-06-20', 3, 100000, 13, 1, 3);
