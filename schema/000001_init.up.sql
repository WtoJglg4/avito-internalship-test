CREATE TABLE features (
    id serial not null unique,
    name varchar(64)
);

CREATE TABLE tags (
    id serial not null unique,
    name varchar(64)
);


CREATE TABLE banners (
    id serial not null unique,
    feature_id int references features(id) on delete set null,
    content  json not null, 
    is_active bool not null default true
);

CREATE TABLE banner_tags (
    banner_id int references banners(id) on delete cascade,
    tag_id int references tags(id) on delete cascade
);

CREATE INDEX features_id_idx ON features(id);

CREATE INDEX tags_id_idx ON tags(id);

CREATE INDEX banners_id_idx ON banners(id);