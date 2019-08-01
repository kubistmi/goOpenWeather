CREATE TABLE IF NOT EXISTS cities (
    id                  integer     primary key,
    city                varchar     not null,
    country             varchar     not null,
    lon                 numeric     not null,
    lat                 numeric     not null
);

TRUNCATE TABLE cities;

CREATE TABLE IF NOT EXISTS weather(
	city_id             integer     not null, --references cities(id),
	conditions          JSONB       not null,
	temperature         numeric     not null,
	pressure            smallint    not null,
	humidity            smallint    not null,
	temp_min            numeric     not null,
	temp_max            numeric     not null,
    visibility          smallint    not null,
	winddir             smallint    not null,
	windspeed           numeric     not null,
    clouds              smallint    not null,
	sunrise             integer     not null,
	sunset              integer     not null,
	timezone            integer     not null,
	extraction_time     integer     not null
);
