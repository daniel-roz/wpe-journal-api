CREATE TABLE IF NOT EXISTS journals (
    uuid UUID PRIMARY KEY,
    title TEXT NOT NULL,
    category TEXT NOT NULL,
    organization TEXT NOT NULL,
    date TEXT NOT NULL,
    article TEXT NOT NULL,
    images JSONB NOT NULL DEFAULT '[]'::jsonb,

    CONSTRAINT journals_category_check
    CHECK (
        category IN (
            'church',
            'classroom',
            'community',
            'court',
            'culture'
        )
    )
);