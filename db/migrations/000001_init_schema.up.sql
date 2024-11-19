-- Active: 1720732902098@@127.0.0.1@3306
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    pb_id TEXT UNIQUE NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    admin INTEGER NOT NULL CHECK (admin IN (0, 1))
);

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

CREATE TABLE exams (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at DATETIME,
    date_of_passing DATETIME NOT NULL,
    exam_number INTEGER NOT NULL,
    professor_id INTEGER NOT NULL,
    template_id INTEGER NOT NULL,
    FOREIGN KEY (professor_id) REFERENCES users (id) FOREIGN KEY (template_id) REFERENCES templates (id)
);

CREATE TABLE exams_exercises (
    exam_id INTEGER NOT NULL,
    exercise_id INTEGER NOT NULL,
    FOREIGN KEY (exam_id) REFERENCES exams (id),
    FOREIGN KEY (exercise_id) REFERENCES exercises (id),
    PRIMARY KEY (exam_id, exercise_id)
);

CREATE TABLE templates (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    uploaded_by INTEGER NOT NULL,
    pb_file_id TEXT UNIQUE NOT NULL,
    template_name TEXT NOT NULL,
    FOREIGN KEY (uploaded_by) REFERENCES users (id)
);
