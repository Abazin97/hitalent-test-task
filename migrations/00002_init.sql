-- +goose Up

CREATE UNIQUE INDEX ux_departments_parent_name
    ON departments (
                    COALESCE(parent_id, 0),
                    lower(trim(name))
        );

-- +goose Down

DROP INDEX ux_departments_parent_name;