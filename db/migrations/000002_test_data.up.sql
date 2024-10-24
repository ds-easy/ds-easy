INSERT INTO
    users (
        pb_id,
        first_name,
        last_name,
        email,
        admin
    )
VALUES (
        'bji9r0uujemtcii',
        'admin',
        'admin',
        'admin@dseasy.com',
        true
    );

INSERT INTO
    lessons (lesson_name, year, subject)
values (
        'integration',
        'terminale',
        'maths'
    );

INSERT INTO
    lessons (lesson_name, year, subject)
values (
        'les reactions acido basiques',
        'terminale',
        'physics'
    );

INSERT INTO
    lessons (lesson_name, year, subject)
values (
        'derivation',
        'seconde',
        'maths'
    );

INSERT INTO
    lessons (lesson_name, year, subject)
values (
        'algebre vectoriel',
        'premiere',
        'maths'
    );

INSERT INTO
    exercises (
        exercise_name,
        exercise_path,
        lesson_id,
        uploaded_by
    )
values (
        'exercise 1',
        'exercises/maths/derivation/exercise1.tex',
        3,
        1
    );

INSERT INTO
    exercises (
        exercise_name,
        exercise_path,
        lesson_id,
        uploaded_by
    )
values (
        'exercise 2',
        'exercises/maths/derivation/exercise2.tex',
        3,
        1
    );

INSERT INTO
    exercises (
        exercise_name,
        exercise_path,
        lesson_id,
        uploaded_by
    )
values (
        'exercise 3',
        'exercises/maths/derivation/exercise3.tex',
        3,
        1
    );

INSERT INTO
    exercises (
        exercise_name,
        exercise_path,
        lesson_id,
        uploaded_by
    )
values (
        'exercise 4',
        'exercises/maths/derivation/exercise4.tex',
        3,
        1
    );

INSERT INTO
    exercises (
        exercise_name,
        exercise_path,
        lesson_id,
        uploaded_by
    )
values (
        'exercise 5',
        'exercises/maths/derivation/exercise5.tex',
        3,
        1
    );

INSERT INTO
    exercises (
        exercise_name,
        exercise_path,
        lesson_id,
        uploaded_by
    )
values (
        'exercise 6',
        'exercises/maths/derivation/exercise6.tex',
        3,
        1
    );
