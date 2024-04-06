CREATE TABLE public.people (
    id SERIAL PRIMARY KEY,
    name VARCHAR(64) NOT NULL,
    surname VARCHAR(64) NOT NULL,
    patronymic VARCHAR(64) NOT NULL
);

CREATE TABLE public.cars (
    id SERIAL PRIMARY KEY,
    owner INTEGER REFERENCES public.people(id) NOT NULL,
    regnum VARCHAR(20) UNIQUE NOT NULL,
    mark VARCHAR(255) NOT NULL,
    model VARCHAR(255) NOT NULL,
    year INTEGER NOT NULL
);
