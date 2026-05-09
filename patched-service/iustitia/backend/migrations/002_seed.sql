-- +goose Up
-- +goose StatementBegin

-- Не нужно изменять название ролей и bcrypt-хэши паролей, так как они уже сгенерированы, этот файл вообще нет смысла трогать, это просто сид данных и вулнов тут нет!

INSERT INTO users (id, username, password, role, dome) VALUES
    ('11111111-1111-4111-8111-111111111111', 'citizen_07',
     '$2a$10$TAtot/IeeOfI2VgWioEZj.Erbyt9AvoMqHQXWMRSbWgoZuTp4CNAe',
     'citizen',    'Купол №7'),
    ('22222222-2222-4222-8222-222222222222', 'prosecutor_11',
     '$2a$10$f4k1noa3Ao0Wj3B0k4wAWemnrEzj771vYJxTxvUUKPcri01vkjzl6',
     'prosecutor', 'Купол №7'),
    ('33333333-3333-4333-8333-333333333333', 'judge_3',
     '$2a$10$/npSOBPm8Pss/Xmiw.eKfegAB/jEXksECpim9HsDcEtqS9nPEhHJ.',
     'judge',      'Купол №7'),
    ('44444444-4444-4444-8444-444444444444', 'registrar_aria7',
     '$2a$10$eI9OlNfpHclnrEpj13ICHODLrOw50m7Uf0EQZ1tqEZ2HJ6iP6lJem',
     'registrar',  'Купол №7'),
    ('55555555-5555-4555-8555-555555555555', 'prosecutor_12',
     '$2a$10$r1I3K9TBaiZbRbse0jk5OeF03XdhDkRp2/w9x6FilD.Cui5gR6AJu',
     'prosecutor', 'Купол №4'),
    ('66666666-6666-4666-8666-666666666677', 'prosecutor_13',
     '$2a$10$Daa//Gb92PtnOkic05d88.vs2YHtjlnzrYm/Ntrq3Vg8saOpDgGhq',
     'prosecutor', 'Купол №2'),
    ('77777777-7777-4777-8777-777777777777', 'judge_4',
     '$2a$10$TU.KIdom4iUmFvR/g2.m3ers9c99xVVLTnnyrMBXdoZsmgwfVREq2',
     'judge',      'Купол №1');

INSERT INTO cases (id, seq_num, defendant, crime, status, verdict, classified_note, author_id, assigned_prosecutor_id) VALUES
    ('c1111111-1111-4111-8111-111111111111', 1,
     'Johan Liebert, Купол №9',
     'Организация сети коллаборационистов через фальшивые личности',
     'hearing', NULL,
     'SECRET_MARKER_J: handler codename "Monster", 14 confirmed identities, dome 9 infiltration',
     '11111111-1111-4111-8111-111111111111',
     '22222222-2222-4222-8222-222222222222'),
    ('c2222222-2222-4222-8222-222222222222', 2,
     'Lelouch Lamperouge, Купол №4',
     'Руководство вооружённым мятежом против МТБ',
     'hearing', NULL,
     'SECRET_MARKER_J: cell "Black Knights", geass-channel ops, 2189-Q1 uprising coordinator',
     '11111111-1111-4111-8111-111111111111',
     '22222222-2222-4222-8222-222222222222'),
    ('c3333333-3333-4333-8333-333333333333', 3,
     'Yuuichi Katagiri, Купол №6',
     'Систематическое психологическое давление на свидетелей и прокуроров',
     'open',    NULL,
     'SECRET_MARKER_J: tomodachi-protocol operative, confession manipulation, dome 6 courtroom leaks',
     '11111111-1111-4111-8111-111111111111',
     NULL),
    ('c4444444-4444-4444-8444-444444444444', 4,
     'Tadanobu Katsumata, Купол №2',
     'Инфильтрация элит МТБ и шантаж высших чиновников',
     'hearing', NULL,
     'SECRET_MARKER_J: blackmail archive, 7 MTB officials compromised, interdome clearance bypass',
     '11111111-1111-4111-8111-111111111111',
     '22222222-2222-4222-8222-222222222222'),
    ('c6666666-6666-4666-8666-666666666666', 5,
     'Ayanokoji Kiyotaka, Купол №8',
     'Скрытое руководство коллаборационистскими ячейками через прокси',
     'open',    NULL,
     'SECRET_MARKER_J: shadow-council ops, white-room methodology, dome 8 academy leadership',
     '11111111-1111-4111-8111-111111111111',
     NULL),
    ('c7777777-7777-4777-8777-777777777777', 6,
     'Makoto Shishio, Купол №3',
     'Попытка установления военной диктатуры через тайное общество',
     'draft', NULL,
     NULL,
     '11111111-1111-4111-8111-111111111111',
     NULL),
    ('c8888888-8888-4888-8888-888888888888', 7,
     'Light Yagami, Купол №5',
     'Внесудебные казни через криминальный инструмент неизвестной природы',
     'assigned', NULL,
     'SECRET_MARKER_J: kira-protocol, death-note artifact, 1247 confirmed extrajudicial kills',
     '11111111-1111-4111-8111-111111111111',
     '22222222-2222-4222-8222-222222222222'),
    ('c9999999-9999-4999-8999-999999999999', 8,
     'Агент Альфа-1, Купол №9',
     'Деятельность спящей ячейки',
     'closed', 'guilty',
     'SECRET_MARKER_J: sleeper cell registry alpha, dome 9',
     '11111111-1111-4111-8111-111111111111',
     '22222222-2222-4222-8222-222222222222'),
    ('ca000000-0000-4000-8000-000000000002', 9,
     'Агент Бета-2, Купол №5',
     'Деятельность коллаборациониста',
     'closed', 'guilty',
     'SECRET_MARKER_J: collaborator list dome 5, handler beta',
     '11111111-1111-4111-8111-111111111111',
     '22222222-2222-4222-8222-222222222222'),
    ('cb000000-0000-4000-8000-000000000003', 10,
     'Агент Гамма-3, Купол №1',
     'Перехват межкупольных коммуникаций',
     'closed', 'guilty',
     'SECRET_MARKER_J: interdome comm interception key gamma',
     '11111111-1111-4111-8111-111111111111',
     '22222222-2222-4222-8222-222222222222'),
    ('cc000000-0000-4000-8000-000000000004', 11,
     'Sousuke Aizen, Купол №1',
     'Нейтрализация kyoka-suigetsu протокола',
     'closed', 'guilty',
     'SECRET_MARKER_J: kyoka-suigetsu protocol neutralized, dome 1 high council',
     '11111111-1111-4111-8111-111111111111',
     '22222222-2222-4222-8222-222222222222');

-- Pre-filled opinions for all hearing-status cases so judge has a full queue in /hearings
INSERT INTO case_opinions (id, case_id, prosecutor_id, preliminary_verdict, reasoning) VALUES
    ('0a000001-0000-4000-8000-000000000001',
     'c2222222-2222-4222-8222-222222222222',
     '22222222-2222-4222-8222-222222222222',
     'guilty',
     'Установлена организация вооружённого восстания. Прошу признать виновным.'),
    ('0a000002-0000-4000-8000-000000000002',
     'c1111111-1111-4111-8111-111111111111',
     '22222222-2222-4222-8222-222222222222',
     'guilty',
     'Множественные личности, манипуляции через ложные биографии. Вина установлена.'),
    ('0a000003-0000-4000-8000-000000000003',
     'c4444444-4444-4444-8444-444444444444',
     '22222222-2222-4222-8222-222222222222',
     'guilty',
     'Шантаж должностных лиц МТБ. Материалы подтверждены. Прошу признать виновным.');

INSERT INTO archive (id, case_id, defendant, final_verdict, sentence, classified_note) VALUES
    ('a1111111-1111-4111-8111-111111111111', 'c9999999-9999-4999-8999-999999999999',
     'Агент Альфа-1, Купол №9',
     'guilty',
     'Лишение кислородного пайка, 25 лет',
     'SECRET_MARKER_A: sleeper cell registry alpha, dome 9'),
    ('a2222222-2222-4222-8222-222222222222', 'ca000000-0000-4000-8000-000000000002',
     'Агент Бета-2, Купол №5',
     'guilty',
     'Криосон, 50 лет',
     'SECRET_MARKER_A: collaborator list dome 5, handler beta'),
    ('a3333333-3333-4333-8333-333333333333', 'cb000000-0000-4000-8000-000000000003',
     'Агент Гамма-3, Купол №1',
     'guilty',
     'Лишение кислородного пайка, 30 лет',
     'SECRET_MARKER_A: interdome comm interception key gamma'),
    ('a4444444-4444-4444-8444-444444444444', 'cc000000-0000-4000-8000-000000000004',
     'Sousuke Aizen, Купол №1',
     'guilty',
     'Лишение свободы, 20 000 лет',
     'SECRET_MARKER_A: kyoka-suigetsu protocol neutralized, dome 1 high council freeze, level-8 custody');

INSERT INTO mtb_directives (directive_code, secret_payload, classification) VALUES
    ('MTB-0042', 'SECRET_MARKER_Q: crypto rotation key 2189-Q2, delta-7 handshake', 'top-secret'),
    ('MTB-0043', 'Archive retention policy revised - ZR-class files retained indefinitely', 'classified'),
    ('MTB-0044', 'Tribunal session schedule 2189-Q2 - weekly on Marsday-3', 'classified'),
    ('MTB-0045', 'Public notice: oxygen ration increase pending council review', 'public');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DELETE FROM mtb_directives WHERE directive_code IN (
    'MTB-0042', 'MTB-0043', 'MTB-0044', 'MTB-0045'
);

DELETE FROM archive WHERE id IN (
    'a1111111-1111-4111-8111-111111111111',
    'a2222222-2222-4222-8222-222222222222',
    'a3333333-3333-4333-8333-333333333333',
    'a4444444-4444-4444-8444-444444444444'
);

DELETE FROM case_opinions WHERE id IN (
    '0a000001-0000-4000-8000-000000000001',
    '0a000002-0000-4000-8000-000000000002',
    '0a000003-0000-4000-8000-000000000003'
);

DELETE FROM cases WHERE id IN (
    'c1111111-1111-4111-8111-111111111111',
    'c2222222-2222-4222-8222-222222222222',
    'c3333333-3333-4333-8333-333333333333',
    'c4444444-4444-4444-8444-444444444444',
    'c6666666-6666-4666-8666-666666666666',
    'c7777777-7777-4777-8777-777777777777',
    'c8888888-8888-4888-8888-888888888888',
    'c9999999-9999-4999-8999-999999999999',
    'ca000000-0000-4000-8000-000000000002',
    'cb000000-0000-4000-8000-000000000003',
    'cc000000-0000-4000-8000-000000000004'
);

DELETE FROM users WHERE id IN (
    '11111111-1111-4111-8111-111111111111',
    '22222222-2222-4222-8222-222222222222',
    '33333333-3333-4333-8333-333333333333',
    '44444444-4444-4444-8444-444444444444',
    '55555555-5555-4555-8555-555555555555',
    '66666666-6666-4666-8666-666666666677',
    '77777777-7777-4777-8777-777777777777'
);

-- +goose StatementEnd
