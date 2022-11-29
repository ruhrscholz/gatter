BEGIN;
    DROP TABLE IF EXISTS favorites;
    DROP TABLE IF EXISTS statuses;
    DROP INDEX IF EXISTS ind_statuses_public_id;
    DROP TABLE IF EXISTS sessions;
    DROP TABLE IF EXISTS users;
    DROP TABLE IF EXISTS accounts;
    DROP TABLE IF EXISTS apps;
COMMIT;