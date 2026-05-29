INSERT INTO users (user_id, name, email, roles, password_hash, department, enabled, date_created, date_updated) VALUES
	('b7f1d86f-0f1f-4c7b-84d7-99a570f14b6f', 'Super Admin Gopher', 'superadmin@example.com', '{SUPER_ADMIN}', '$2a$10$XHwZsUFDExC0fsgpm4oHn.M9uT2kHaSTd5MO9QuuBl4nJt71BfubO', NULL, true, '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('5cf37266-3473-4006-984f-9325122678b7', 'Admin Gopher', 'admin@example.com', '{SCHOOL_ADMIN}', '$2a$10$1ggfMVZV6Js0ybvJufLRUOWHS5f6KneuP0XwwHpJ8L8ipdry9f2/a', NULL, true, '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('c41fa5d3-d61f-45f1-b054-d2c7a3704019', 'Teacher Gopher', 'teacher@example.com', '{TEACHER}', '$2a$10$OSJzEzY9lSNIK8pbN/T0e.myTTnatVrYonbvXY.aqXaS.qZRn6ELm', NULL, true, '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('45b5fbd3-755f-4379-8f07-a58d4a30fa2f', 'User Gopher', 'user@example.com', '{STUDENT}', '$2a$10$9/XASPKBbJKVfCAZKDH.UuhsuALDr5vVm6VrYA9VFR8rccK86C1hW', NULL, true, '2019-03-24 00:00:00', '2019-03-24 00:00:00')
ON CONFLICT DO NOTHING;

INSERT INTO admissions_staff_profiles (staff_profile_id, user_id, admissions_roles, is_active, date_created, date_updated) VALUES
	('3adf627a-d32a-4436-8ca0-5e76fa1593d7', 'b7f1d86f-0f1f-4c7b-84d7-99a570f14b6f', '{ADMISSIONS_ADMIN}', true, '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('ccda637b-e691-4575-b487-93538fd4c943', '5cf37266-3473-4006-984f-9325122678b7', '{ADMISSIONS_ADMIN,APPLICATION_REVIEWER,REPORT_VIEWER}', true, '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('aee70594-4f18-4353-8bec-9d3fef476d03', 'c41fa5d3-d61f-45f1-b054-d2c7a3704019', '{APPLICATION_REVIEWER}', true, '2019-03-24 00:00:00', '2019-03-24 00:00:00')
ON CONFLICT DO NOTHING;
