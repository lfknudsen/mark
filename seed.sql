INSERT INTO folders (version, name, desc, date_added, date_modified, date_last_used, guid, icon)
VALUES ('1.0.0', 'Folder 001', 'Description 001', 1234556789,
        123456789, 123456789, '9c76112d-357d-4dd0-a6b9-2fe1be65aba1', ''),
       ('1.0.0', 'Folder 002', 'Description 002', 1234556789,
        123456789, 123456789, '3303ec6d-d5c0-474e-8f33-234f97b52a9d', ''),
       ('1.0.0', 'Folder 003', 'Description 003', 1234556789,
        123456789, 123456789, '118c0d34-25fb-460b-9a33-6910383175fe', ''),
       ('1.0.0', 'Folder 004', 'Description 004', 1234556789,
        123456789, 123456789, 'c686b0ab-2ade-49ba-b2fa-49d84ed2446a', ''),
       ('1.0.0', 'Folder 005', 'Description 005', 1234556789,
        123456789, 123456789, '7a9108a4-dc46-4ba5-ac7c-f8a72eded6c4', ''),
       ('1.0.0', 'Folder 006', 'Description 006', 1234556789,
        123456789, 123456789, '966ca238-b09b-4842-aae6-ee603d5619f1', ''),
       ('1.0.0', 'Folder 007', 'Description 007', 1234556789,
        123456789, 123456789, '5d1c0ad1-aea0-4c08-b9a2-188da25abcca', '');

UPDATE subfolder_relations
SET parent = '3303ec6d-d5c0-474e-8f33-234f97b52a9d'
WHERE child = '118c0d34-25fb-460b-9a33-6910383175fe';
UPDATE subfolder_relations
SET parent = '3303ec6d-d5c0-474e-8f33-234f97b52a9d'
WHERE child = '5d1c0ad1-aea0-4c08-b9a2-188da25abcca';
UPDATE subfolder_relations
SET parent = 'c686b0ab-2ade-49ba-b2fa-49d84ed2446a'
WHERE child = '966ca238-b09b-4842-aae6-ee603d5619f1';
UPDATE subfolder_relations
SET parent = '966ca238-b09b-4842-aae6-ee603d5619f1'
WHERE child = '9c76112d-357d-4dd0-a6b9-2fe1be65aba1';

INSERT INTO bookmarks (version, url, name, desc, date_added, date_modified, date_last_used, guid, icon)
VALUES ('1.0.0', 'url 01', 'name 01', 'desc 01', 123456789,
        123456789, 123456789, '879c0256-efb9-4b48-acb1-9376d1408c8e', ''),
       ('1.0.0', 'url 02', 'name 02', 'desc 02', 123456789,
        123456789, 123456789, 'b95f14d9-6351-44c4-a313-d3371ef61875', ''),
       ('1.0.0', 'url 03', 'name 03', 'desc 03', 123456789,
        123456789, 123456789, 'b6bfe958-44ab-4853-a018-a4406f8d652a', ''),
       ('1.0.0', 'url 04', 'name 04', 'desc 04', 123456789,
        123456789, 123456789, '62a0a8d3-223e-4035-9a27-3bbc8b23745d', ''),
       ('1.0.0', 'url 05', 'name 05', 'desc 05', 123456789,
        123456789, 123456789, 'a941717a-a791-42c5-a4f2-a643547e3258', ''),
       ('1.0.0', 'url 06', 'name 06', 'desc 06', 123456789,
        123456789, 123456789, '801b3663-c65d-40ac-bf46-89cc1dc5a033', ''),
       ('1.0.0', 'url 07', 'name 07', 'desc 07', 123456789,
        123456789, 123456789, '5e48b58d-98d3-40e7-a06d-ec9256bdb82f', ''),
       ('1.0.0', 'url 08', 'name 08', 'desc 08', 123456789,
        123456789, 123456789, '52722b5a-dd18-4c3f-804f-cfb36028fbab', ''),
       ('1.0.0', 'url 09', 'name 09', 'desc 09', 123456789,
        123456789, 123456789, '0d029441-e46c-49a0-b079-dc4288e64cd9', '');

UPDATE child_relations
SET parent = '3303ec6d-d5c0-474e-8f33-234f97b52a9d'
WHERE child = 'b95f14d9-6351-44c4-a313-d3371ef61875';
UPDATE child_relations
SET parent = '3303ec6d-d5c0-474e-8f33-234f97b52a9d'
WHERE child = '5e48b58d-98d3-40e7-a06d-ec9256bdb82f';
UPDATE child_relations
SET parent = '3303ec6d-d5c0-474e-8f33-234f97b52a9d'
WHERE child = '0d029441-e46c-49a0-b079-dc4288e64cd9';
UPDATE child_relations
SET parent = '3303ec6d-d5c0-474e-8f33-234f97b52a9d'
WHERE child = '52722b5a-dd18-4c3f-804f-cfb36028fbab';
UPDATE child_relations
SET parent = '118c0d34-25fb-460b-9a33-6910383175fe'
WHERE child = '879c0256-efb9-4b48-acb1-9376d1408c8e';
UPDATE child_relations
SET parent = '5d1c0ad1-aea0-4c08-b9a2-188da25abcca'
WHERE child = 'b6bfe958-44ab-4853-a018-a4406f8d652a';