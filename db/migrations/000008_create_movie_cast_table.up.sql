CREATE TABLE public.movie_cast (
	movie_id int4 NULL,
	cast_id int4 NULL
);


-- public.movie_cast foreign keys

ALTER TABLE public.movie_cast ADD CONSTRAINT movie_cast_cast_id_fkey FOREIGN KEY (cast_id) REFERENCES public.casts(id);
ALTER TABLE public.movie_cast ADD CONSTRAINT movie_cast_movie_id_fkey FOREIGN KEY (movie_id) REFERENCES public.movies(id);