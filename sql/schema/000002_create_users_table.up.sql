CREATE TABLE IF NOT EXISTS public.users (
    id integer NOT NULL GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    created_at timestamptz NOT NULL DEFAULT now(),
    name text NOT NULL,
    email citext NOT NULL UNIQUE,
    password_hash bytea NOT NULL,
    active boolean NOT NULL DEFAULT true,
    version integer NOT NULL DEFAULT 1
);