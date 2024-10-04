CREATE TABLE IF NOT EXISTS public.saved_recipes (
    user_id integer NOT NULL,
    recipe_id integer NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, recipe_id),
    FOREIGN KEY (user_id) REFERENCES public.users (id) ON DELETE CASCADE,
    FOREIGN KEY (recipe_id) REFERENCES public.recipes (id) ON DELETE CASCADE
);