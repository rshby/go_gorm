-- CREATE TABLE SAMPE
CREATE TABLE sample (
                        id VARCHAR(100) NOT NULL UNIQUE PRIMARY KEY ,
                        name VARCHAR(100)
)ENGINE = InnoDB;


CREATE TABLE users (
    id VARCHAR(255) NOT NULL UNIQUE PRIMARY KEY ,
    password VARCHAR(255) NOT NULL ,
    name VARCHAR(255) NOT NULL ,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
)ENGINE = InnoDB;

-- add firstname middlename lastname in table users
ALTER TABLE users RENAME COLUMN name TO first_name;
ALTER TABLE users ADD COLUMN middle_name VARCHAR(255) NULL AFTER first_name;
ALTER TABLE users ADD COLUMN last_name VARCHAR(255) NULL AFTER middle_name;