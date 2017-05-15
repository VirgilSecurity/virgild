-- +migrate Up
CREATE TABLE applications (id varchar(24), card_id VARCHAR(64) NOT NULL,name VARCHAR(20) NOT NULL,bundle VARCHAR NOT NULL,description VARCHAR(1000),created_at varchar(28) NOT NULL,updated_at varchar(28) NOT NULL, PRIMARY KEY (id));
CREATE TABLE tokens (id varchar(24),name VARCHAR(20) NOT NULL,value VARCHAR NOT NULL,is_active BOOLEAN,application_id varchar(24) NOT NULL, created_at varchar(28) NOT NULL,updated_at varchar(28) NOT NULL, PRIMARY KEY (id));
CREATE UNIQUE INDEX applications_unique_cardid_idx ON applications (card_id);
CREATE UNIQUE INDEX tokens_unique_value_idx ON tokens (value);
CREATE INDEX tokens_applicationid_idx ON tokens (application_id);

-- +migrate Down
DROP TABLE applications;
DROP TABLE tokens;
DROP INDEX applications_unique_cardid_idx;
DROP INDEX tokens_unique_value_idx;
DROP INDEX tokens_applicationid_idx;
