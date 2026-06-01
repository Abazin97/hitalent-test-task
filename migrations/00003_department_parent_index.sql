-- +goose Up

CREATE UNIQUE INDEX ux_departments_parent_name
    ON departments (
                    parent_id,
                    lower(trim(name))
        );

-- +goose Down

DROP INDEX ux_departments_parent_name;