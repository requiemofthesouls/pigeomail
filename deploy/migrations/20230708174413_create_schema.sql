-- +goose Up 
CREATE SCHEMA IF NOT EXISTS pigeomail;
-- +goose Down 
DROP SCHEMA IF EXISTS pigeomail;
