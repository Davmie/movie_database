drop table if exists actors cascade;
create table public.actors(
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    first_name VARCHAR(35) NOT NULL,
    last_name VARCHAR(35) NOT NULL,
    gender CHAR(1) NOT NULL,
    birthday DATE NOT NULL
);

drop table if exists movies cascade;
create table public.movies(
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    title VARCHAR(150) NOT NULL,
    description VARCHAR(1000) NOT NULL,
    release_date DATE NOT NULL,
    rating INT NOT NULL
);

drop table if exists movies_actors cascade;
create table public.movies_actors(
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    movie_id INT NOT NULL,
    foreign key (movie_id) references public.movies(id),
    actor_id INT NOT NULL,
    foreign key (actor_id) references public.actors(id)
);

drop table if exists users cascade;
create table public.users(
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    login VARCHAR(256) NOT NULL UNIQUE,
    password VARCHAR(128) NOT NULL,
    user_role VARCHAR(20) NOT NULL
);