-- add a new column to the recipes table to store which user created the recipe. This will allow us to query for all recipes created by a specific user. We'll also add a foreign key constraint to ensure that the user_id references an existing user in the users table.

ALTER TABLE public.recipes
ADD COLUMN user_id integer NOT NULL DEFAULT 1,
ADD CONSTRAINT fk_user_id
FOREIGN KEY (user_id)
REFERENCES public.users (id);
