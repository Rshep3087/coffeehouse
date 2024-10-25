-- add a column to store the image url of the recipe

ALTER TABLE public.recipes
ADD COLUMN image_url TEXT NULL;