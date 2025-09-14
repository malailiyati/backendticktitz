CREATE TABLE public.movies (
	id serial4 NOT NULL,
	title varchar(255) NOT NULL,
	director_id int4 NOT NULL,
	poster varchar(255) NOT NULL,
	background_poster varchar(255) NOT NULL,
	releasedate date NULL,
	duration interval NULL,
	synopsis text NULL,
	popularity int4 NULL,
	created_at timestamp NULL,
	updated_at timestamp NULL,
	deleted_at timestamp NULL,
	CONSTRAINT movies_pkey PRIMARY KEY (id)
);


-- public.movies foreign keys

ALTER TABLE public.movies ADD CONSTRAINT movies_director_id_fkey FOREIGN KEY (director_id) REFERENCES public.director(id);