CREATE TABLE public.cinema (
	id serial4 NOT NULL,
	"name" varchar(255) NOT NULL,
	price int4 NOT NULL,
	CONSTRAINT cinema_pkey PRIMARY KEY (id)
);