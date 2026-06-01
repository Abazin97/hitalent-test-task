-- +goose Up

ALTER TABLE departments
    ADD CONSTRAINT chk_department_name
        CHECK (length(trim(name)) > 0);

ALTER TABLE employees
    ADD CONSTRAINT chk_employee_name
        CHECK (length(trim(full_name)) > 0);

ALTER TABLE employees
    ADD CONSTRAINT chk_employee_position
        CHECK (length(trim(position)) > 0);

-- +goose Down

ALTER TABLE employees
DROP CONSTRAINT chk_employee_position;

ALTER TABLE employees
DROP CONSTRAINT chk_employee_name;

ALTER TABLE departments
DROP CONSTRAINT chk_department_name;