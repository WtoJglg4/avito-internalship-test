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
    content  jsonb not null, 
    is_active bool not null default true,
    version integer not null, 
    created_at timestamp not null, 
    updated_at timestamp not null, 
    tags_hash varchar(255) not null
);

CREATE TABLE banner_tags (
    banner_id int references banners(id) on delete cascade,
    tag_id int references tags(id) on delete cascade
);

CREATE TABLE banner_versions (
    id serial not null unique,
    banner_id int references banners(id) on delete cascade,
    content json not null,
    version integer not null, 
    updated_at timestamp not null,
    is_active bool not null default false
);

CREATE INDEX banner_versions_banner_id_idx ON banner_versions(banner_id);

CREATE INDEX features_id_idx ON features(id);

CREATE INDEX tags_id_idx ON tags(id);

CREATE INDEX banners_id_idx ON banners(id);

CREATE TABLE users
(
    id              serial not null unique,
    login           varchar(255) not null unique, 
    password_hash   varchar(255) not null,
    role            varchar(255) not null
);