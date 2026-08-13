-- +goose Up 
CREATE TABLE images (
    id UUID PRIMARY KEY, 
    objectKey TEXT NOT NULL UNIQUE, 
    ext TEXT NOT NULL, 
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL, 
    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down 
DROP TABLE images; 