-- +goose Up
-- +goose StatementBegin
CREATE TABLE if not exists users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    hashed_password VARCHAR(255) NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE TABLE if not exists tasks (
    id uuid PRIMARY KEY,
    user_id BIGINT NOT NULL,
    translator VARCHAR(255) NOT NULL,
    code text NOT NULL,
    result text NOT NULL,
    task_status VARCHAR(16) NOT NULL,
    stderr text,
    stdout text,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now(),

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE 
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE if exists tasks;
DROP TABLE if exists users;
-- +goose StatementEnd
