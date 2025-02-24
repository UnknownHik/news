CREATE TABLE IF NOT EXISTS "News" (
    "Id" BIGSERIAL PRIMARY KEY,
    "Title" TEXT NOT NULL,
    "Content" TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS "NewsCategories" (
    "NewsId" BIGINT NOT NULL,
    "CategoryId" BIGINT NOT NULL,
    PRIMARY KEY ("NewsId", "CategoryId"),
    FOREIGN KEY ("NewsId") REFERENCES "News"("Id") ON DELETE CASCADE
);

INSERT INTO "News" ("Title", "Content") VALUES
                                    ('news1', 'content1'),
                                    ('news2', 'content2'),
                                    ('news3', 'content3'),
                                    ('news4', 'content4'),
                                    ('news5', 'content5'),
                                    ('news6', 'content6'),
                                    ('news7', 'content7'),
                                    ('news8', 'content8'),
                                    ('news9', 'content9'),
                                    ('news10', 'content10'),
                                    ('news11', 'content11'),
                                    ('news12', 'content12'),
                                    ('news13', 'content13'),
                                    ('news14', 'content14'),
                                    ('news15', 'content15');

INSERT INTO "NewsCategories" ("NewsId", "CategoryId") VALUES
                                      (1, 3),
                                      (1, 6),
                                      (1, 7),
                                      (2, 1),
                                      (2, 2),
                                      (3, 5),
                                      (4, 1),
                                      (5, 2),
                                      (6, 1),
                                      (6, 4),
                                      (6, 7),
                                      (7, 5),
                                      (7, 6),
                                      (8, 2),
                                      (9, 1),
                                      (10, 7),
                                      (11, 4),
                                      (12, 2),
                                      (13, 5),
                                      (14, 5),
                                      (15, 7),
                                      (15, 8);