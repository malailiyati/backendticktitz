CREATE TABLE public.schedule (
	id serial4 NOT NULL,
	movie_id int4 NULL,
	"date" date NULL,
	time_id int4 NULL,
	location_id int4 NULL,
	cinema_id int4 NULL,
	CONSTRAINT schedule_pkey PRIMARY KEY (id)
);


-- public.schedule foreign keys

ALTER TABLE public.schedule ADD CONSTRAINT schedule_cinema_id_fkey FOREIGN KEY (cinema_id) REFERENCES public.cinema(id);
ALTER TABLE public.schedule ADD CONSTRAINT schedule_location_id_fkey FOREIGN KEY (location_id) REFERENCES public."location"(id);
ALTER TABLE public.schedule ADD CONSTRAINT schedule_movie_id_fkey FOREIGN KEY (movie_id) REFERENCES public.movies(id);
ALTER TABLE public.schedule ADD CONSTRAINT schedule_time_id_fkey FOREIGN KEY (time_id) REFERENCES public."time"(id);