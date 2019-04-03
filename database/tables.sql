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
 createdOn TIMESTAMP NOT NULL
);