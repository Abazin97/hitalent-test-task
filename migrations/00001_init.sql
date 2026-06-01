-- +goose Up

CREATE TABLE departments (
                             id BIGSERIAL PRIMARY KEY,
                             name VARCHAR(200) NOT NULL,
                             parent_id BIGINT NULL REFERENCES departments(id) ON DELETE CASCADE,
                             created_at TIMESTAMP NOT NULL DEFAULT NOW(),

                             CONSTRAINT chk_department_name
                                 CHECK (length(trim(name)) > 0)
);

CREATE TABLE employees (
                           id BIGSERIAL PRIMARY KEY,
                           department_id BIGINT NOT NULL REFERENCES departments(id) ON DELETE CASCADE,
                           full_name VARCHAR(200) NOT NULL,
                           position VARCHAR(200) NOT NULL,
                           hired_at DATE NULL,
                           created_at TIMESTAMP NOT NULL DEFAULT NOW(),

                           CONSTRAINT chk_employee_name
                               CHECK (length(trim(full_name)) > 0),

                           CONSTRAINT chk_employee_position
                               CHECK (length(trim(position)) > 0)
);

-- +goose Down

DROP TABLE employees;
DROP TABLE departments;