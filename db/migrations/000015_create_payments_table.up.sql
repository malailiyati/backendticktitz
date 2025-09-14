CREATE TABLE public.payments (
	id serial4 NOT NULL,
	"method" varchar(70) NOT NULL,
	CONSTRAINT payments_pkey PRIMARY KEY (id)
);