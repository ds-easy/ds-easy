-- Active: 1720732902098@@127.0.0.1@3306
DROP TABLE IF EXISTS users;

CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    admin INTEGER NOT NULL CHECK (admin IN (0, 1))
);

DROP TABLE IF EXISTS lessons;

CREATE TABLE lessons (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    lesson_name TEXT NOT NULL,
    year TEXT NOT NULL CHECK (
        year IN (
            'premiere',
            'seconde',
            'terminale'
        )
    ),
    subject TEXT NOT NULL CHECK (
        subject IN ('maths', 'physics')
    )
);

DROP TABLE IF EXISTS exercises;

CREATE TABLE exercises (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    exercise_name TEXT UNIQUE NOT NULL,
    exercise_path TEXT UNIQUE NOT NULL,
    lesson_id INTEGER NOT NULL,
    uploaded_by INTEGER NOT NULL,
    FOREIGN KEY (lesson_id) REFERENCES lessons (id),
    FOREIGN KEY (uploaded_by) REFERENCES users (id)
);

DROP TABLE IF EXISTS exams;

CREATE TABLE exams (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    date_of_passing DATETIME,
    exam_number INTEGER,
    professor_id INTEGER,
    FOREIGN KEY (professor_id) REFERENCES users (id)
);

DROP TABLE IF EXISTS exams_exercises;

CREATE TABLE exams_exercises (
    exam_id INTEGER,
    exercise_id INTEGER,
    FOREIGN KEY (exam_id) REFERENCES exams (id),
    FOREIGN KEY (exercise_id) REFERENCES exercises (id),
    PRIMARY KEY (exam_id, exercise_id)
);