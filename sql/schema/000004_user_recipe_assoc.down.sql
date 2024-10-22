-- drop the user_id column from the recipes table

ALTER TABLE public.recipes
DROP COLUMN user_id;