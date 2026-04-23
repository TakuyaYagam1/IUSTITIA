-- +goose Up
-- +goose StatementBegin

CREATE TABLE users (
    id         TEXT PRIMARY KEY,                    -- uuid
    username   TEXT UNIQUE NOT NULL,
    password   TEXT NOT NULL,                       -- bcrypt
    role       TEXT NOT NULL CHECK (role IN ('citizen', 'prosecutor', 'judge', 'registrar')),
    dome       TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE cases (
    id                     TEXT PRIMARY KEY,               -- uuid
    seq_num                INTEGER UNIQUE NOT NULL,
    defendant              TEXT NOT NULL,
    crime                  TEXT NOT NULL,
    status                 TEXT NOT NULL DEFAULT 'draft'
                           CHECK (status IN ('draft', 'open', 'assigned', 'hearing', 'closed')),
    verdict                TEXT
                           CHECK (verdict IN ('guilty', 'acquitted', 'dismissed') OR verdict IS NULL),
    classified_note        TEXT,
    author_id              TEXT REFERENCES users(id) ON DELETE SET NULL,
    assigned_prosecutor_id TEXT REFERENCES users(id) ON DELETE SET NULL,
    created_at             DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_cases_status              ON cases(status);
CREATE INDEX idx_cases_created_at          ON cases(created_at);
CREATE INDEX idx_cases_author              ON cases(author_id);
CREATE INDEX idx_cases_assigned_prosecutor ON cases(assigned_prosecutor_id);

CREATE TABLE complaints (
    id            TEXT PRIMARY KEY,                 -- uuid
    case_id       TEXT NOT NULL REFERENCES cases(id) ON DELETE CASCADE,
    author_id     TEXT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    text          TEXT NOT NULL,
    evidence_url  TEXT,
    evidence_data TEXT,
    created_at    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_complaints_case_id ON complaints(case_id);

CREATE TABLE archive (
    id              TEXT PRIMARY KEY,               -- uuid
    case_id         TEXT REFERENCES cases(id) ON DELETE SET NULL,
    defendant       TEXT NOT NULL,
    final_verdict   TEXT NOT NULL
                    CHECK (final_verdict IN ('guilty', 'acquitted', 'dismissed')),
    sentence        TEXT,
    classified_note TEXT,
    archived_at     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_archive_case_id      ON archive(case_id);
CREATE INDEX idx_archive_archived_at  ON archive(archived_at);

CREATE TABLE documents (
    id         TEXT PRIMARY KEY,                    -- uuid
    case_id    TEXT NOT NULL REFERENCES cases(id) ON DELETE CASCADE,
    author_id  TEXT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    content    TEXT NOT NULL,
    template   TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_documents_case_id ON documents(case_id);

CREATE TABLE case_opinions (
    id                  TEXT PRIMARY KEY,           -- uuid
    case_id             TEXT NOT NULL UNIQUE REFERENCES cases(id) ON DELETE CASCADE,
    prosecutor_id       TEXT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    preliminary_verdict TEXT NOT NULL
                        CHECK (preliminary_verdict IN ('guilty', 'acquitted', 'dismissed')),
    reasoning           TEXT NOT NULL,
    filed_at            DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_case_opinions_case       ON case_opinions(case_id);
CREATE INDEX idx_case_opinions_prosecutor ON case_opinions(prosecutor_id);

CREATE TABLE mtb_directives (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    directive_code  TEXT UNIQUE NOT NULL,
    secret_payload  TEXT NOT NULL,
    classification  TEXT NOT NULL DEFAULT 'classified'
                    CHECK (classification IN ('public', 'classified', 'top-secret')),
    issued_at       DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_mtb_directives_classification ON mtb_directives(classification);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_mtb_directives_classification;
DROP TABLE IF EXISTS mtb_directives;

DROP INDEX IF EXISTS idx_case_opinions_prosecutor;
DROP INDEX IF EXISTS idx_case_opinions_case;
DROP TABLE IF EXISTS case_opinions;

DROP INDEX IF EXISTS idx_documents_case_id;
DROP TABLE IF EXISTS documents;

DROP INDEX IF EXISTS idx_archive_archived_at;
DROP INDEX IF EXISTS idx_archive_case_id;
DROP TABLE IF EXISTS archive;

DROP INDEX IF EXISTS idx_complaints_case_id;
DROP TABLE IF EXISTS complaints;

DROP INDEX IF EXISTS idx_cases_assigned_prosecutor;
DROP INDEX IF EXISTS idx_cases_author;
DROP INDEX IF EXISTS idx_cases_created_at;
DROP INDEX IF EXISTS idx_cases_status;
DROP TABLE IF EXISTS cases;

DROP TABLE IF EXISTS users;

-- +goose StatementEnd
