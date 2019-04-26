CREATE TABLE daylevels(
 id serial PRIMARY KEY,
 focus integer NOT NULL,
 fischio_orecchie integer NOT NULL,
 power_energy integer NOT NULL,
 dormito integer NOT NULL,
 PR  integer NOT NULL,
 ansia  integer NOT NULL,
 arrabiato integer NOT NULL,
 irritato integer NOT NULL,
 depresso  integer NOT NULL,
 cinque_tibetani BOOLEAN NOT NULL,
 meditazione BOOLEAN NOT NULL,
 createdOn TIMESTAMP default current_timestamp
);

CREATE TABLE meditations(
 id serial PRIMARY KEY,
 meditation text NOT NULL,
 timesused integer NOT NULL,
 createdOn TIMESTAMP default current_timestamp
);

CREATE TABLE actionsmed(
 id serial PRIMARY KEY,
 action text NOT NULL,
 idmed integer NOT NULL,
 idusr integer NOT NULL,
 createdOn TIMESTAMP default current_timestamp
);

CREATE TABLE users (
id serial PRIMARY KEY,
first_name text NOT NULL,
last_name text NOT NULL,
user_name VARCHAR(50) UNIQUE,
user_pwd text NOT NULL,
idrole integer NOT NULL
);

CREATE TABLE roles (
id serial PRIMARY KEY,
role text NOT NULL
);

ALTER TABLE actionsmed
    ADD CONSTRAINT fk_actions_meditation FOREIGN KEY (idmed) REFERENCES meditations (id),
    ADD CONSTRAINT fk_actions_users FOREIGN KEY (idusr) REFERENCES users (id);


ALTER TABLE users
    ADD CONSTRAINT fk_users_roles FOREIGN KEY (idrole) REFERENCES roles (id);

ALTER TABLE meditations
  ADD COLUMN pref_month VARCHAR(15),
  ADD COLUMN pref_day integer;

ALTER TABLE users
  ADD COLUMN user_pwd text;

ALTER TABLE users
    ALTER COLUMN user_name TYPE VARCHAR(50),
    ADD UNIQUE (user_name);


INSERT INTO roles (role)
    VALUES ('user');


##UPDATE roles SET role = 'user';


GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO ml;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO ml;


ALTER TABLE daylevels
    ADD COLUMN uid integer,
    ADD CONSTRAINT fk_daylevels_users FOREIGN KEY (uid) REFERENCES users (id);


UPDATE daylevels SET uid = 1;


ALTER TABLE daylevels
    ALTER COLUMN uid TYPE integer;
    
ALTER TABLE daylevels
    ALTER COLUMN uid SET NOT NULL;