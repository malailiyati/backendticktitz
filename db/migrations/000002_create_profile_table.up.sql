CREATE TABLE public.profile (
	id serial4 NOT NULL,
	firstname varchar(100) NOT NULL,
	lastname varchar(100) NOT NULL,
	phone varchar(20) NOT NULL,
	users_id int4 NULL,
	created_at timestamp NULL,
	updated_at timestamp NULL,
	profile_picture text NULL,
	CONSTRAINT profile_pkey PRIMARY KEY (id)
);


-- public.profile foreign keys

ALTER TABLE public.profile ADD CONSTRAINT profile_users_id_fkey FOREIGN KEY (users_id) REFERENCES public.users(id);