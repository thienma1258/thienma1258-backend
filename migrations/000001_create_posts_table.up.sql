CREATE TABLE posts (
 id serial NOT NULL  PRIMARY KEY,
 user_id varchar(255) DEFAULT NULL,
 title varchar(255) NOT NULL,
 slug varchar(255) NOT NULL UNIQUE,
 views INT NOT NULL DEFAULT '0',
 image varchar(255) NOT NULL,
 body text NOT NULL ,
 published SMALLINT NOT NULL,
 created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
 updated_at timestamp NOT NULL
);