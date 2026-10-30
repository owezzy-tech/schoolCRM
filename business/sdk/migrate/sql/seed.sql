INSERT INTO users (user_id, name, email, roles, password_hash, department, enabled, date_created, date_updated) VALUES
	('b7f1d86f-0f1f-4c7b-84d7-99a570f14b6f', 'Super Admin Gopher', 'superadmin@example.com', '{SUPER_ADMIN}', '$2a$10$XHwZsUFDExC0fsgpm4oHn.M9uT2kHaSTd5MO9QuuBl4nJt71BfubO', NULL, true, '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('5cf37266-3473-4006-984f-9325122678b7', 'Admin Gopher', 'admin@example.com', '{SCHOOL_ADMIN}', '$2a$10$1ggfMVZV6Js0ybvJufLRUOWHS5f6KneuP0XwwHpJ8L8ipdry9f2/a', NULL, true, '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('c41fa5d3-d61f-45f1-b054-d2c7a3704019', 'Teacher Gopher', 'teacher@example.com', '{TEACHER}', '$2a$10$OSJzEzY9lSNIK8pbN/T0e.myTTnatVrYonbvXY.aqXaS.qZRn6ELm', NULL, true, '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('45b5fbd3-755f-4379-8f07-a58d4a30fa2f', 'User Gopher', 'user@example.com', '{STUDENT}', '$2a$10$9/XASPKBbJKVfCAZKDH.UuhsuALDr5vVm6VrYA9VFR8rccK86C1hW', NULL, true, '2019-03-24 00:00:00', '2019-03-24 00:00:00')
ON CONFLICT DO NOTHING;

INSERT INTO admissions_staff_profiles (staff_profile_id, user_id, admissions_roles, is_active, date_created, date_updated) VALUES
	('3adf627a-d32a-4436-8ca0-5e76fa1593d7', 'b7f1d86f-0f1f-4c7b-84d7-99a570f14b6f', '{ADMISSIONS_ADMIN}', true, '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('ccda637b-e691-4575-b487-93538fd4c943', '5cf37266-3473-4006-984f-9325122678b7', '{ADMISSIONS_ADMIN,APPLICATION_REVIEWER,EVENT_MANAGER,REPORT_VIEWER}', true, '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('aee70594-4f18-4353-8bec-9d3fef476d03', 'c41fa5d3-d61f-45f1-b054-d2c7a3704019', '{APPLICATION_REVIEWER}', true, '2019-03-24 00:00:00', '2019-03-24 00:00:00')
ON CONFLICT DO NOTHING;

INSERT INTO users (user_id, name, email, roles, password_hash, department, enabled, date_created, date_updated) VALUES
	('f47ac10b-58cc-4372-a567-0e02b2c3d479', 'John Applicant', 'applicant@example.com', '{STUDENT}', '$2a$10$XHwZsUFDExC0fsgpm4oHn.M9uT2kHaSTd5MO9QuuBl4nJt71BfubO', NULL, true, '2019-03-24 00:00:00', '2019-03-24 00:00:00')
ON CONFLICT DO NOTHING;

UPDATE users
SET roles = '{STUDENT}'
WHERE user_id = 'f47ac10b-58cc-4372-a567-0e02b2c3d479';

INSERT INTO admissions_constituents (constituent_id, first_name, last_name, preferred_name, middle_name, suffix, date_of_birth, primary_email, primary_phone, external_sis_id, lifecycle_stage, duplicate_status, duplicate_of_id, sis_synced_at, date_created, date_updated, national_id, national_id_verified_at, national_id_verified_by_adapter, upi, upi_verified_at, upi_verified_by_adapter, kcse_index_number, kcse_index_verified_at, kcse_index_verified_by_adapter, sms_opt_in, whatsapp_opt_in, email_opt_in, notification_priority) VALUES
	('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'John', 'Applicant', NULL, NULL, NULL, '2000-01-15 00:00:00', '<applicant@example.com>', '+254712345678', NULL, 'APPLICANT', 'ACTIVE', NULL, NULL, '2019-03-24 00:00:00', '2019-03-24 00:00:00', '12345678', NULL, NULL, NULL, NULL, NULL, '12345678901', NULL, NULL, true, false, true, ARRAY['SMS', 'WHATSAPP', 'EMAIL']::TEXT[])
ON CONFLICT DO NOTHING;

-- Admissions admin fixture constituents for campaigns, communications, inquiries, and dashboard filters.
INSERT INTO admissions_constituents (constituent_id, first_name, last_name, preferred_name, middle_name, suffix, date_of_birth, primary_email, primary_phone, external_sis_id, lifecycle_stage, duplicate_status, duplicate_of_id, sis_synced_at, date_created, date_updated, national_id, national_id_verified_at, national_id_verified_by_adapter, upi, upi_verified_at, upi_verified_by_adapter, kcse_index_number, kcse_index_verified_at, kcse_index_verified_by_adapter, sms_opt_in, whatsapp_opt_in, email_opt_in, notification_priority) VALUES
	('a3000000-0000-4000-8000-000000000001', 'Sofia', 'Martinez', NULL, NULL, NULL, '2002-04-12 00:00:00', 'sofia.martinez@example.com', '+254711100001', 'ADM-CON-921', 'APPLICANT', 'ACTIVE', NULL, '2026-05-28 08:00:00', '2026-05-20 09:00:00', '2026-05-28 08:00:00', '22345678', '2026-05-20 09:05:00', 'iprs', 'UPI-ADM-921', '2026-05-20 09:06:00', 'kuccps', '20220092101', '2026-05-20 09:07:00', 'knec', true, true, true, ARRAY['EMAIL', 'SMS', 'WHATSAPP']::TEXT[]),
	('a3000000-0000-4000-8000-000000000002', 'James', 'Okoro', NULL, NULL, NULL, '2001-09-03 00:00:00', 'james.okoro@example.com', '+254711100002', 'ADM-CON-884', 'APPLICANT', 'ACTIVE', NULL, '2026-05-24 08:00:00', '2026-05-15 10:00:00', '2026-05-24 08:00:00', '22345679', NULL, NULL, 'UPI-ADM-884', NULL, NULL, '20210088401', NULL, NULL, true, false, true, ARRAY['SMS', 'EMAIL', 'WHATSAPP']::TEXT[]),
	('a3000000-0000-4000-8000-000000000003', 'Achieng', 'Otieno', NULL, NULL, NULL, '2003-02-21 00:00:00', 'achieng.otieno@example.com', '+254711100003', 'ADM-CON-924', 'INQUIRY', 'ACTIVE', NULL, NULL, '2026-05-18 11:00:00', '2026-05-18 11:00:00', NULL, NULL, NULL, NULL, NULL, NULL, '20230092401', NULL, NULL, true, true, true, ARRAY['WHATSAPP', 'SMS', 'EMAIL']::TEXT[]),
	('a3000000-0000-4000-8000-000000000004', 'Priya', 'Patel', NULL, NULL, NULL, '1998-11-30 00:00:00', 'priya.patel@example.com', '+254711100004', 'ADM-CON-773', 'APPLICANT', 'ACTIVE', NULL, '2026-05-25 08:00:00', '2026-05-10 12:00:00', '2026-05-25 08:00:00', '22345680', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, true, false, true, ARRAY['EMAIL', 'SMS', 'WHATSAPP']::TEXT[]),
	('a3000000-0000-4000-8000-000000000005', 'Liam', 'Chen', NULL, NULL, NULL, '2000-07-18 00:00:00', 'liam.chen@example.com', '+254711100005', 'ADM-CON-742', 'APPLICANT', 'ACTIVE', NULL, '2026-05-26 08:00:00', '2026-05-16 13:00:00', '2026-05-26 08:00:00', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, false, true, true, ARRAY['EMAIL', 'WHATSAPP', 'SMS']::TEXT[]),
	('a3000000-0000-4000-8000-000000000006', 'Aisha', 'Bello', NULL, NULL, NULL, '1999-05-14 00:00:00', 'aisha.bello@example.com', '+254711100006', 'ADM-CON-695', 'ADMITTED', 'ACTIVE', NULL, '2026-05-29 08:00:00', '2026-05-12 14:00:00', '2026-05-29 08:00:00', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, true, false, true, ARRAY['EMAIL', 'SMS', 'WHATSAPP']::TEXT[]),
	('a3000000-0000-4000-8000-000000000007', 'Amara', 'Ndlovu', NULL, NULL, NULL, '2002-12-01 00:00:00', 'amara.ndlovu@example.com', '+254711100007', 'ADM-CON-652', 'ADMITTED', 'ACTIVE', NULL, '2026-05-28 08:00:00', '2026-05-12 15:00:00', '2026-05-28 08:00:00', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, true, true, true, ARRAY['SMS', 'WHATSAPP', 'EMAIL']::TEXT[]),
	('a3000000-0000-4000-8000-000000000008', 'Noah', 'Williams', NULL, NULL, NULL, '2001-01-22 00:00:00', 'noah.williams@example.com', '+254711100008', 'ADM-CON-641', 'PROSPECT', 'ACTIVE', NULL, NULL, '2026-05-17 16:00:00', '2026-05-17 16:00:00', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, true, false, true, ARRAY['SMS', 'EMAIL', 'WHATSAPP']::TEXT[]),
	('a3000000-0000-4000-8000-000000000009', 'Wanjiku', 'Njoroge', NULL, NULL, NULL, '2002-08-09 00:00:00', 'wanjiku.njoroge@example.com', '+254711100009', 'ADM-CON-640', 'APPLICANT', 'ACTIVE', NULL, '2026-05-27 08:00:00', '2026-05-19 17:00:00', '2026-05-27 08:00:00', '22345681', NULL, NULL, 'UPI-ADM-640', NULL, NULL, '20220064001', NULL, NULL, true, true, true, ARRAY['SMS', 'WHATSAPP', 'EMAIL']::TEXT[]),
	('a3000000-0000-4000-8000-000000000010', 'Emma', 'Davis', NULL, NULL, NULL, '2000-03-25 00:00:00', 'emma.davis@example.com', '+254711100010', 'ADM-CON-614', 'ENROLLED', 'ACTIVE', NULL, '2026-05-26 08:00:00', '2026-05-08 18:00:00', '2026-05-26 08:00:00', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, true, false, true, ARRAY['EMAIL', 'SMS', 'WHATSAPP']::TEXT[])
ON CONFLICT DO NOTHING;

INSERT INTO admissions_programs (program_id, external_sis_id, name, code, description, degree_level, is_active, date_created, date_updated) VALUES
	('b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'KUCCPS-BCOM-2026-QA', 'Bachelor of Commerce', 'BCOM-QA', 'QA fixture commerce programme', 'BACHELOR', true, '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a23', 'SELF-BSCN-2026-QA', 'Bachelor of Science in Nursing', 'BSCN-QA', 'Self-sponsored nursing programme for private applicant flows', 'BACHELOR', true, '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('b2eebc99-9c0b-4ef8-bb6d-6bb9bd380a24', 'DIP-ICT-2026-QA', 'Diploma in Information Communication Technology', 'DICT-QA', 'Diploma programme used by TVET and diploma intake flows', 'DIPLOMA', true, '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('b3eebc99-9c0b-4ef8-bb6d-6bb9bd380a25', 'CERT-ECDE-2026-QA', 'Certificate in Early Childhood Development Education', 'CERT-ECDE-QA', 'Certificate programme for lower KNQF application settings', 'CERTIFICATE', true, '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('b4eebc99-9c0b-4ef8-bb6d-6bb9bd380a26', 'MSC-DS-2026-QA', 'Master of Science in Data Science', 'MSC-DS-QA', 'Graduate programme used by masters application flows', 'MASTERS', true, '2019-03-24 00:00:00', '2019-03-24 00:00:00')
ON CONFLICT DO NOTHING;

INSERT INTO admissions_academic_terms (academic_term_id, external_sis_id, name, code, term_type, start_date, end_date, application_start_date, application_deadline, is_active, date_created, date_updated) VALUES
	('c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'INTAKE-2026-QA', '2026 Main Intake', 'INT2026-QA', 'MAIN', '2026-01-01 00:00:00', '2026-12-31 00:00:00', '2025-09-01 00:00:00', '2026-01-15 00:00:00', true, '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a34', 'INTAKE-MAY-2026-QA', 'May 2026 Intake', 'MAY2026-QA', 'TRIMESTER', '2026-05-04 00:00:00', '2026-08-28 00:00:00', '2026-01-05 00:00:00', '2026-04-17 00:00:00', true, '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('c2eebc99-9c0b-4ef8-bb6d-6bb9bd380a35', 'INTAKE-SEP-2026-QA', 'September 2026 Intake', 'SEP2026-QA', 'MAIN', '2026-09-07 00:00:00', '2026-12-18 00:00:00', '2026-04-01 00:00:00', '2026-08-21 00:00:00', true, '2019-03-24 00:00:00', '2019-03-24 00:00:00')
ON CONFLICT DO NOTHING;

INSERT INTO admissions_application_form_templates (form_template_id, program_id, academic_term_id, application_type, name, description, version, required_fields, checklist_items, is_active, priority, date_created, date_updated) VALUES
	('d1eebc99-9c0b-4ef8-bb6d-6bb9bd380a45', 'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'KUCCPS_PLACEMENT', 'KUCCPS placement undergraduate form', 'Default active form for placed undergraduate applicants.', 1, '[{"fieldName":"kuccpsPlacementId","fieldType":"TEXT","required":true,"displayOrder":10,"validation":"^[A-Z0-9-]{6,40}$"},{"fieldName":"kcseIndexNumber","fieldType":"TEXT","required":true,"displayOrder":20,"validation":"^[0-9]{11}$"},{"fieldName":"nationalId","fieldType":"TEXT","required":true,"displayOrder":30,"validation":"^[0-9]{7,8}$"}]'::JSONB, '[{"itemKey":"national_id","documentName":"National ID or Passport","description":"Government-issued identity document used for verification.","required":true,"displayOrder":10},{"itemKey":"kcse_result_slip","documentName":"KCSE Result Slip","description":"KNEC result slip or certificate matching the KCSE index number.","required":true,"displayOrder":20},{"itemKey":"kuccps_admission_letter","documentName":"KUCCPS Placement Letter","description":"Downloaded KUCCPS placement or admission letter.","required":true,"displayOrder":30}]'::JSONB, true, 10, '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('d2eebc99-9c0b-4ef8-bb6d-6bb9bd380a46', 'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a23', 'c2eebc99-9c0b-4ef8-bb6d-6bb9bd380a35', 'SELF_SPONSORED_UNDERGRAD', 'Self-sponsored undergraduate form', 'Default form for self-sponsored undergraduate applicants.', 1, '[{"fieldName":"kcseIndexNumber","fieldType":"TEXT","required":true,"displayOrder":10,"validation":"^[0-9]{11}$"},{"fieldName":"sponsorshipSource","fieldType":"SELECT","required":true,"displayOrder":20,"validation":"SELF|PARENT_GUARDIAN|EMPLOYER|SCHOLARSHIP"},{"fieldName":"countyCode","fieldType":"TEXT","required":true,"displayOrder":30,"validation":"^[0-9]{1,2}$"}]'::JSONB, '[{"itemKey":"national_id","documentName":"National ID or Passport","description":"Government-issued identity document used for verification.","required":true,"displayOrder":10},{"itemKey":"kcse_certificate","documentName":"KCSE Certificate","description":"Final KNEC certificate or official result slip.","required":true,"displayOrder":20},{"itemKey":"passport_photo","documentName":"Passport Photo","description":"Recent passport-size photo for student records.","required":false,"displayOrder":30}]'::JSONB, true, 20, '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('d3eebc99-9c0b-4ef8-bb6d-6bb9bd380a47', 'b2eebc99-9c0b-4ef8-bb6d-6bb9bd380a24', 'c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a34', 'DIPLOMA', 'Diploma intake form', 'Default form for diploma and TVET-style admissions.', 1, '[{"fieldName":"kcseIndexNumber","fieldType":"TEXT","required":true,"displayOrder":10,"validation":"^[0-9]{11}$"},{"fieldName":"highestQualification","fieldType":"SELECT","required":true,"displayOrder":20,"validation":"KCSE|CERTIFICATE|ARTISAN|OTHER"},{"fieldName":"preferredStudyMode","fieldType":"SELECT","required":true,"displayOrder":30,"validation":"FULL_TIME|PART_TIME|EVENING|ONLINE"}]'::JSONB, '[{"itemKey":"kcse_result_slip","documentName":"KCSE Result Slip","description":"KNEC result slip or certificate.","required":true,"displayOrder":10},{"itemKey":"previous_certificate","documentName":"Previous Certificate","description":"Certificate or transcript for applicants with prior qualifications.","required":false,"displayOrder":20}]'::JSONB, true, 30, '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('d4eebc99-9c0b-4ef8-bb6d-6bb9bd380a48', 'b4eebc99-9c0b-4ef8-bb6d-6bb9bd380a26', 'c2eebc99-9c0b-4ef8-bb6d-6bb9bd380a35', 'MASTERS', 'Masters application form', 'Default form for graduate school applicants.', 1, '[{"fieldName":"undergraduateInstitution","fieldType":"TEXT","required":true,"displayOrder":10},{"fieldName":"undergraduateGpa","fieldType":"NUMBER","required":true,"displayOrder":20,"validation":"^([0-3](\\.[0-9]{1,2})?|4(\\.0{1,2})?)$"},{"fieldName":"researchInterest","fieldType":"TEXTAREA","required":true,"displayOrder":30}]'::JSONB, '[{"itemKey":"degree_certificate","documentName":"Degree Certificate","description":"Certified copy of undergraduate degree certificate.","required":true,"displayOrder":10},{"itemKey":"academic_transcripts","documentName":"Academic Transcripts","description":"Official transcripts for all undergraduate study.","required":true,"displayOrder":20},{"itemKey":"referee_letter","documentName":"Referee Letter","description":"Academic or professional referee letter.","required":false,"displayOrder":30}]'::JSONB, true, 40, '2019-03-24 00:00:00', '2019-03-24 00:00:00')
ON CONFLICT DO NOTHING;

INSERT INTO admissions_custom_field_definitions (custom_field_definition_id, owner, field_key, label, description, data_type, is_required, options, validation, is_searchable, is_reportable, is_importable, is_exportable, display_order, is_active, date_created, date_updated) VALUES
	('e1eebc99-9c0b-4ef8-bb6d-6bb9bd380a51', 'CONSTITUENT', 'county_code', 'County of Residence', 'Kenya county code used for admissions reporting and regional outreach.', 'SELECT', true, ARRAY['1','16','22','30','32','37','47']::TEXT[], '^[0-9]{1,2}$', true, true, true, true, 10, true, '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('e2eebc99-9c0b-4ef8-bb6d-6bb9bd380a52', 'CONSTITUENT', 'preferred_notification_channel', 'Preferred Notification Channel', 'Applicant communication preference for admissions updates.', 'SELECT', false, ARRAY['SMS','WHATSAPP','EMAIL']::TEXT[], NULL, true, true, true, true, 20, true, '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('e3eebc99-9c0b-4ef8-bb6d-6bb9bd380a53', 'APPLICATION', 'sponsorship_source', 'Sponsorship Source', 'Expected fee sponsorship source for non-KUCCPS applicants.', 'SELECT', true, ARRAY['SELF','PARENT_GUARDIAN','EMPLOYER','SCHOLARSHIP']::TEXT[], NULL, true, true, true, true, 10, true, '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('e4eebc99-9c0b-4ef8-bb6d-6bb9bd380a54', 'APPLICATION', 'preferred_study_mode', 'Preferred Study Mode', 'Applicant-selected delivery mode for timetable planning.', 'SELECT', true, ARRAY['FULL_TIME','PART_TIME','EVENING','ONLINE']::TEXT[], NULL, true, true, true, true, 20, true, '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('e5eebc99-9c0b-4ef8-bb6d-6bb9bd380a55', 'APPLICATION', 'research_interest', 'Research Interest', 'Graduate applicant research interest statement.', 'TEXTAREA', false, ARRAY[]::TEXT[], NULL, false, true, true, true, 30, true, '2019-03-24 00:00:00', '2019-03-24 00:00:00')
ON CONFLICT DO NOTHING;

INSERT INTO admissions_import_batches (import_batch_id, source, file_type, target, status, file_name, storage_key, uploaded_by_id, total_rows, valid_rows, invalid_rows, duplicate_rows, field_mapping, invalid_report_key, validation_summary, committed_at, date_created, date_updated) VALUES
	('f1eebc99-9c0b-4ef8-bb6d-6bb9bd380a61', 'MANUAL_UPLOAD', 'CSV', 'CONSTITUENTS', 'PREVIEWED', 'admissions-constituents-template.csv', 'seed/imports/admissions/constituents-template.csv', '5cf37266-3473-4006-984f-9325122678b7', 0, 0, 0, 0, '{"first_name":"First Name","last_name":"Last Name","primary_email":"Email","primary_phone":"Phone","national_id":"National ID","kcse_index_number":"KCSE Index Number","county_code":"County Code","preferred_notification_channel":"Preferred Notification Channel"}'::JSONB, NULL, 'Seeded blank constituent import template for admin settings mapping preview.', NULL, '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('f2eebc99-9c0b-4ef8-bb6d-6bb9bd380a62', 'MANUAL_UPLOAD', 'XLSX', 'APPLICATIONS', 'PREVIEWED', 'admissions-applications-template.xlsx', 'seed/imports/admissions/applications-template.xlsx', '5cf37266-3473-4006-984f-9325122678b7', 0, 0, 0, 0, '{"application_type":"Application Type","program_code":"Programme Code","term_code":"Intake Code","kcse_index_number":"KCSE Index Number","sponsorship_source":"Sponsorship Source","preferred_study_mode":"Preferred Study Mode"}'::JSONB, NULL, 'Seeded blank application import template for application settings mapping preview.', NULL, '2019-03-24 00:00:00', '2019-03-24 00:00:00')
ON CONFLICT DO NOTHING;

INSERT INTO admissions_applications (application_id, constituent_id, program_id, academic_term_id, application_type, status, assigned_reviewer_id, submitted_at, date_created, date_updated, kuccps_placement, kcse_result) VALUES
	('d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a44', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'KUCCPS_PLACEMENT', 'SUBMITTED', 'c41fa5d3-d61f-45f1-b054-d2c7a3704019', '2026-01-02 09:00:00', '2019-03-24 00:00:00', '2019-03-24 00:00:00', '{}'::JSONB, '{}'::JSONB)
ON CONFLICT DO NOTHING;

INSERT INTO admissions_applications (application_id, constituent_id, program_id, academic_term_id, application_type, status, assigned_reviewer_id, submitted_at, date_created, date_updated, kuccps_placement, kcse_result) VALUES
	('d3000000-0000-4000-8000-000000000001', 'a3000000-0000-4000-8000-000000000001', 'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'KUCCPS_PLACEMENT', 'IN_REVIEW', 'c41fa5d3-d61f-45f1-b054-d2c7a3704019', '2026-05-20 09:30:00', '2026-05-20 09:00:00', '2026-05-28 08:00:00', '{"placementId":"KUCCPS-921","programmeCode":"BCOM-QA"}'::JSONB, '{"indexNumber":"20220092101","meanGrade":"B+"}'::JSONB),
	('d3000000-0000-4000-8000-000000000002', 'a3000000-0000-4000-8000-000000000002', 'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'DIPLOMA', 'AWAITING_DOCUMENTS', 'c41fa5d3-d61f-45f1-b054-d2c7a3704019', '2026-05-15 10:30:00', '2026-05-15 10:00:00', '2026-05-24 08:00:00', '{}'::JSONB, '{"indexNumber":"20210088401","meanGrade":"B"}'::JSONB),
	('d3000000-0000-4000-8000-000000000003', 'a3000000-0000-4000-8000-000000000004', 'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'MASTERS', 'READY_FOR_REVIEW', '5cf37266-3473-4006-984f-9325122678b7', '2026-05-10 12:30:00', '2026-05-10 12:00:00', '2026-05-25 08:00:00', '{}'::JSONB, '{}'::JSONB),
	('d3000000-0000-4000-8000-000000000004', 'a3000000-0000-4000-8000-000000000005', 'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'SELF_SPONSORED_UNDERGRAD', 'SUBMITTED', 'c41fa5d3-d61f-45f1-b054-d2c7a3704019', '2026-05-16 13:30:00', '2026-05-16 13:00:00', '2026-05-26 08:00:00', '{}'::JSONB, '{}'::JSONB),
	('d3000000-0000-4000-8000-000000000005', 'a3000000-0000-4000-8000-000000000006', 'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'KUCCPS_PLACEMENT', 'ADMITTED', 'c41fa5d3-d61f-45f1-b054-d2c7a3704019', '2026-05-12 14:30:00', '2026-05-12 14:00:00', '2026-05-29 08:00:00', '{}'::JSONB, '{}'::JSONB),
	('d3000000-0000-4000-8000-000000000006', 'a3000000-0000-4000-8000-000000000007', 'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'SELF_SPONSORED_UNDERGRAD', 'ADMITTED', '5cf37266-3473-4006-984f-9325122678b7', '2026-05-12 15:30:00', '2026-05-12 15:00:00', '2026-05-28 08:00:00', '{}'::JSONB, '{}'::JSONB),
	('d3000000-0000-4000-8000-000000000007', 'a3000000-0000-4000-8000-000000000009', 'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'KUCCPS_PLACEMENT', 'DECISION_PENDING', 'c41fa5d3-d61f-45f1-b054-d2c7a3704019', '2026-05-19 17:30:00', '2026-05-19 17:00:00', '2026-05-27 08:00:00', '{"placementId":"KUCCPS-640","programmeCode":"BCOM-QA"}'::JSONB, '{"indexNumber":"20220064001","meanGrade":"A-"}'::JSONB),
	('d3000000-0000-4000-8000-000000000008', 'a3000000-0000-4000-8000-000000000010', 'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'DIPLOMA', 'ENROLLED', 'c41fa5d3-d61f-45f1-b054-d2c7a3704019', '2026-05-08 18:30:00', '2026-05-08 18:00:00', '2026-05-26 08:00:00', '{}'::JSONB, '{}'::JSONB)
ON CONFLICT DO NOTHING;

UPDATE admissions_applications
SET
	status = 'IN_REVIEW',
	assigned_reviewer_id = 'c41fa5d3-d61f-45f1-b054-d2c7a3704019',
	submitted_at = COALESCE(submitted_at, '2026-01-02 09:00:00'),
	date_updated = '2026-01-04 11:00:00',
	kuccps_placement = COALESCE(kuccps_placement, '{"placementID":"KUCCPS-2026-0001","admissionNumber":"ADM-2026-0001","institutionCode":"UON","programmeCode":"UON-BCOM","programmeName":"Bachelor of Commerce","placementYear":2026,"clusterCode":"CL04","clusterPoints":42.7,"weightedPointsNote":"QA fixture placement from KUCCPS import"}'::JSONB),
	kcse_result = COALESCE(kcse_result, '{"indexNumber":"12345678901","examYear":2024,"subjects":[{"subjectCode":"101","grade":"B+","points":10},{"subjectCode":"102","grade":"A-","points":11},{"subjectCode":"121","grade":"B","points":9}],"meanGrade":"B+","meanPoints":70}'::JSONB)
WHERE application_id = 'd0eebc99-9c0b-4ef8-bb6d-6bb9bd380a44';

INSERT INTO admissions_applicant_profiles (applicant_profile_id, user_id, constituent_id, is_active, date_created, date_updated) VALUES
	('f0eebc99-9c0b-4ef8-bb6d-6bb9bd380a66', 'f47ac10b-58cc-4372-a567-0e02b2c3d479', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', true, '2019-03-24 00:00:00', '2019-03-24 00:00:00')
ON CONFLICT DO NOTHING;

INSERT INTO admissions_application_transitions (application_transition_id, application_id, from_status, to_status, actor_id, reason, note, metadata, date_created) VALUES
	('d0eebc99-9c0b-4ef8-bb6d-6bb9bd381001', 'd0eebc99-9c0b-4ef8-bb6d-6bb9bd380a44', 'DRAFT', 'SUBMITTED', 'f47ac10b-58cc-4372-a567-0e02b2c3d479', 'Applicant submitted KUCCPS placement application', 'Portal QA seed submission by John Applicant.', '{"channel":"applicant_portal","fixture":"schoolCRM-dpj.2"}'::JSONB, '2026-01-02 09:00:00'),
	('d0eebc99-9c0b-4ef8-bb6d-6bb9bd381002', 'd0eebc99-9c0b-4ef8-bb6d-6bb9bd380a44', 'SUBMITTED', 'AWAITING_DOCUMENTS', '5cf37266-3473-4006-984f-9325122678b7', 'Required admissions evidence requested', 'Checklist generated for KCSE result, national ID, placement letter, and passport photo.', '{"channel":"admin_review","fixture":"schoolCRM-dpj.2"}'::JSONB, '2026-01-02 10:00:00'),
	('d0eebc99-9c0b-4ef8-bb6d-6bb9bd381003', 'd0eebc99-9c0b-4ef8-bb6d-6bb9bd380a44', 'AWAITING_DOCUMENTS', 'READY_FOR_REVIEW', 'f47ac10b-58cc-4372-a567-0e02b2c3d479', 'Applicant uploaded required evidence', 'All required documents are present for reviewer intake.', '{"channel":"applicant_portal","uploadedDocuments":3,"fixture":"schoolCRM-dpj.2"}'::JSONB, '2026-01-03 14:30:00'),
	('d0eebc99-9c0b-4ef8-bb6d-6bb9bd381004', 'd0eebc99-9c0b-4ef8-bb6d-6bb9bd380a44', 'READY_FOR_REVIEW', 'IN_REVIEW', 'c41fa5d3-d61f-45f1-b054-d2c7a3704019', 'Reviewer started eligibility review', 'Teacher reviewer is validating KUCCPS placement and KCSE evidence.', '{"channel":"admin_review","fixture":"schoolCRM-dpj.2"}'::JSONB, '2026-01-04 11:00:00')
ON CONFLICT DO NOTHING;

INSERT INTO admissions_checklist_items (checklist_item_id, application_id, item_key, document_name, description, is_required, status, display_order, date_created, date_updated) VALUES
	('d0eebc99-9c0b-4ef8-bb6d-6bb9bd382001', 'd0eebc99-9c0b-4ef8-bb6d-6bb9bd380a44', 'kcse-result-slip', 'KCSE Result Slip', 'Certified KCSE result slip showing index number and mean grade.', true, 'ACCEPTED', 1, '2026-01-02 10:00:00', '2026-01-04 10:15:00'),
	('d0eebc99-9c0b-4ef8-bb6d-6bb9bd382002', 'd0eebc99-9c0b-4ef8-bb6d-6bb9bd380a44', 'national-id', 'National ID or Birth Certificate', 'Identity document used for admissions identity verification.', true, 'PENDING_REVIEW', 2, '2026-01-02 10:00:00', '2026-01-03 14:30:00'),
	('d0eebc99-9c0b-4ef8-bb6d-6bb9bd382003', 'd0eebc99-9c0b-4ef8-bb6d-6bb9bd380a44', 'kuccps-placement-letter', 'KUCCPS Placement Letter', 'KUCCPS placement letter for the 2026 intake.', true, 'ACCEPTED', 3, '2026-01-02 10:00:00', '2026-01-04 10:45:00'),
	('d0eebc99-9c0b-4ef8-bb6d-6bb9bd382004', 'd0eebc99-9c0b-4ef8-bb6d-6bb9bd380a44', 'passport-photo', 'Passport Photo', 'Recent passport-size photo for student record creation.', false, 'UPLOADED', 4, '2026-01-02 10:00:00', '2026-01-03 14:35:00')
ON CONFLICT DO NOTHING;

INSERT INTO admissions_documents (document_id, application_id, checklist_item_id, file_name, content_type, size_bytes, storage_key, status, reviewer_id, reviewer_notes, uploaded_by_id, uploaded_at, reviewed_at, date_created, date_updated) VALUES
	('d0eebc99-9c0b-4ef8-bb6d-6bb9bd383001', 'd0eebc99-9c0b-4ef8-bb6d-6bb9bd380a44', 'd0eebc99-9c0b-4ef8-bb6d-6bb9bd382001', 'john-applicant-kcse-result-slip.pdf', 'application/pdf', 348912, 'filepond/admissions/d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a44/kcse-result-slip.pdf', 'ACCEPTED', 'c41fa5d3-d61f-45f1-b054-d2c7a3704019', 'KCSE index and grades match submitted application data.', 'f47ac10b-58cc-4372-a567-0e02b2c3d479', '2026-01-03 14:05:00', '2026-01-04 10:15:00', '2026-01-03 14:05:00', '2026-01-04 10:15:00'),
	('d0eebc99-9c0b-4ef8-bb6d-6bb9bd383002', 'd0eebc99-9c0b-4ef8-bb6d-6bb9bd380a44', 'd0eebc99-9c0b-4ef8-bb6d-6bb9bd382002', 'john-applicant-national-id.jpg', 'image/jpeg', 812304, 'filepond/admissions/d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a44/national-id.jpg', 'PENDING_REVIEW', NULL, NULL, 'f47ac10b-58cc-4372-a567-0e02b2c3d479', '2026-01-03 14:12:00', NULL, '2026-01-03 14:12:00', '2026-01-03 14:12:00'),
	('d0eebc99-9c0b-4ef8-bb6d-6bb9bd383003', 'd0eebc99-9c0b-4ef8-bb6d-6bb9bd380a44', 'd0eebc99-9c0b-4ef8-bb6d-6bb9bd382003', 'john-applicant-kuccps-placement-letter.pdf', 'application/pdf', 421876, 'filepond/admissions/d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a44/kuccps-placement-letter.pdf', 'ACCEPTED', 'c41fa5d3-d61f-45f1-b054-d2c7a3704019', 'Placement letter confirms Bachelor of Commerce admission for 2026 intake.', 'f47ac10b-58cc-4372-a567-0e02b2c3d479', '2026-01-03 14:20:00', '2026-01-04 10:45:00', '2026-01-03 14:20:00', '2026-01-04 10:45:00'),
	('d0eebc99-9c0b-4ef8-bb6d-6bb9bd383004', 'd0eebc99-9c0b-4ef8-bb6d-6bb9bd380a44', 'd0eebc99-9c0b-4ef8-bb6d-6bb9bd382004', 'john-applicant-passport-photo.png', 'image/png', 156204, 'filepond/admissions/d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a44/passport-photo.png', 'UPLOADED', NULL, NULL, 'f47ac10b-58cc-4372-a567-0e02b2c3d479', '2026-01-03 14:35:00', NULL, '2026-01-03 14:35:00', '2026-01-03 14:35:00')
ON CONFLICT DO NOTHING;

INSERT INTO admissions_events (
	event_id,
	title,
	event_type,
	status,
	description,
	start_time,
	end_time,
	location,
	is_virtual,
	capacity,
	registration_deadline,
	auto_confirmation_enabled,
	auto_reminder_enabled,
	date_created,
	date_updated
) VALUES
	('e1000000-0000-4000-8000-000000000001', 'Nairobi Open Day', 'open-day', 'upcoming', 'Prospective students and families tour campus, meet admissions staff, and attend programme showcases.', '2026-06-20 08:00:00', '2026-06-20 15:00:00', 'Nairobi Main Campus', false, 600, '2026-06-19 23:59:00', true, true, '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('e1000000-0000-4000-8000-000000000002', 'Engineering Webinar', 'webinar', 'upcoming', 'Virtual admissions webinar covering engineering pathways, fees, and scholarship timelines.', '2026-06-12 17:00:00', '2026-06-12 18:30:00', 'Zoom', true, 450, '2026-06-11 23:59:00', true, true, '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('e1000000-0000-4000-8000-000000000003', 'MBA Information Session', 'info-session', 'completed', 'Completed postgraduate information evening for working professionals evaluating MBA applications.', '2026-06-02 15:00:00', '2026-06-02 17:00:00', 'CBD Learning Centre', false, 120, '2026-06-01 23:59:00', true, false, '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('e1000000-0000-4000-8000-000000000004', 'Architecture Studio Tour', 'campus-tour', 'live', 'Guided studio and fabrication lab tour for architecture prospects currently on campus.', '2026-06-05 09:00:00', '2026-06-05 11:00:00', 'Design Block', false, 40, '2026-06-04 18:00:00', false, false, '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('e1000000-0000-4000-8000-000000000005', 'County Career Fair Booth', 'fair', 'cancelled', 'Regional fair appearance cancelled due to venue changes.', '2026-06-28 09:00:00', '2026-06-28 16:00:00', 'KICC Nairobi', false, 1000, '2026-06-27 23:59:00', false, false, '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('e1000000-0000-4000-8000-000000000006', 'September Intake Planning', 'open-day', 'draft', 'Draft event being prepared for the September intake recruitment cycle.', '2026-09-10 09:00:00', '2026-09-10 14:00:00', 'TBD', false, 800, '2026-09-08 23:59:00', true, true, '2019-03-24 00:00:00', '2019-03-24 00:00:00')
ON CONFLICT DO NOTHING;

INSERT INTO admissions_event_registrations (
	event_registration_id,
	event_id,
	constituent_id,
	first_name,
	last_name,
	email,
	phone,
	status,
	match_status,
	source,
	registered_at,
	checked_in_at,
	checked_in_by_id,
	date_created,
	date_updated
) VALUES
	('e2000000-0000-4000-8000-000000000001', 'e1000000-0000-4000-8000-000000000001', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'John', 'Applicant', 'applicant@example.com', '+254712345678', 'registered', 'matched', 'portal', '2026-05-28 10:15:00', NULL, NULL, '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('e2000000-0000-4000-8000-000000000002', 'e1000000-0000-4000-8000-000000000001', NULL, 'Sofia', 'Mwangi', 'sofia.mwangi@example.com', '+254711000001', 'registered', 'new-prospect', 'campaign', '2026-05-29 08:45:00', NULL, NULL, '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('e2000000-0000-4000-8000-000000000003', 'e1000000-0000-4000-8000-000000000002', NULL, 'Liam', 'Otieno', 'liam.otieno@example.com', '+254711000002', 'registered', 'needs-review', 'staff', '2026-05-30 13:10:00', NULL, NULL, '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('e2000000-0000-4000-8000-000000000004', 'e1000000-0000-4000-8000-000000000003', NULL, 'Priya', 'Patel', 'priya.patel@example.com', '+254711000003', 'checked-in', 'matched', 'staff', '2026-05-25 09:30:00', '2026-06-02 14:52:00', '5cf37266-3473-4006-984f-9325122678b7', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('e2000000-0000-4000-8000-000000000005', 'e1000000-0000-4000-8000-000000000004', NULL, 'Achieng', 'Okoro', 'achieng.okoro@example.com', '+254711000004', 'checked-in', 'new-prospect', 'portal', '2026-06-01 16:20:00', '2026-06-05 09:10:00', '5cf37266-3473-4006-984f-9325122678b7', '2019-03-24 00:00:00', '2019-03-24 00:00:00')
ON CONFLICT DO NOTHING;

-- Admissions admin inquiry queue fixtures covering new, contacted, converted, and closed filters.
INSERT INTO admissions_inquiries (inquiry_id, constituent_id, first_name, last_name, date_of_birth, primary_email, primary_phone, program_of_interest, term_of_interest, source, utm_source, utm_medium, utm_campaign, message, status, date_created, date_updated) VALUES
	('e3000000-0000-4000-8000-000000000001', 'a3000000-0000-4000-8000-000000000003', 'Achieng', 'Otieno', '2003-02-21 00:00:00', 'achieng.otieno@example.com', '+254711100003', 'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'portal', 'google', 'cpc', 'fall-open-house', 'Interested in September intake scholarships and document requirements.', 'NEW', '2026-05-18 11:00:00', '2026-05-18 11:00:00'),
	('e3000000-0000-4000-8000-000000000002', 'a3000000-0000-4000-8000-000000000008', 'Noah', 'Williams', '2001-01-22 00:00:00', 'noah.williams@example.com', '+254711100008', NULL, 'c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'campaign', 'sms', 'broadcast', 'campus-tour-reminder', 'Asked for campus tour directions and admissions requirements.', 'CONTACTED', '2026-05-17 16:00:00', '2026-05-27 17:45:00'),
	('e3000000-0000-4000-8000-000000000003', 'a3000000-0000-4000-8000-000000000001', 'Sofia', 'Martinez', '2002-04-12 00:00:00', 'sofia.martinez@example.com', '+254711100001', 'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'event', 'open-house', 'qr', 'fall-open-house', 'Converted after Nairobi Open Day registration follow-up.', 'CONVERTED', '2026-05-20 09:00:00', '2026-05-20 09:30:00'),
	('e3000000-0000-4000-8000-000000000004', 'a3000000-0000-4000-8000-000000000010', 'Emma', 'Davis', '2000-03-25 00:00:00', 'emma.davis@example.com', '+254711100010', 'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', NULL, 'phone', NULL, NULL, NULL, 'Parent financial aid questions resolved after enrolment.', 'CLOSED', '2026-05-08 18:00:00', '2026-05-26 11:20:00')
ON CONFLICT DO NOTHING;

-- Admissions admin campaign fixtures. IDs are deterministic for API wiring and E2E assertions.
INSERT INTO admissions_campaigns (campaign_id, name, status, channel, audience_name, template_name, message_preview, segment, metrics, starts_at, ends_at, created_by_id, date_created, date_updated) VALUES
	('c3000000-0000-4000-8000-000000000001', 'Fall 2026 Open House Invite', 'ACTIVE', 'EMAIL', 'Warm prospects near campus', 'Open house invitation', 'Join our admissions team for a guided campus experience and program Q&A.', '{"applicationTypes":["KUCCPS_PLACEMENT","SELF_SPONSORED_UNDERGRAD"],"applicationStatuses":["DRAFT","SUBMITTED"],"academicTerms":["September 2026"],"programs":["B.Sc. Computer Science","Bachelor of Commerce"],"eventAttendance":"ANY","leadScoreBands":["WARM","HOT"],"recruiters":["Maya Schultz","Andre Park"],"territories":["Nairobi","Kiambu"]}'::JSONB, '{"audienceSize":2450,"sent":2450,"delivered":2402,"opened":1081,"clicked":294,"bounced":48,"replied":0}'::JSONB, '2026-05-01 00:00:00', '2026-06-06 23:59:00', '5cf37266-3473-4006-984f-9325122678b7', '2026-04-24 09:00:00', '2026-05-18 09:00:00'),
	('c3000000-0000-4000-8000-000000000002', 'Missing Documents Reminder', 'ACTIVE', 'SMS', 'Awaiting documents', 'Checklist nudge SMS', 'Your application is almost complete. Upload your remaining documents today.', '{"applicationTypes":["KUCCPS_PLACEMENT","DIPLOMA","MASTERS"],"applicationStatuses":["AWAITING_DOCUMENTS"],"academicTerms":["September 2026"],"programs":["All active programs"],"eventAttendance":"ANY","leadScoreBands":["HOT","READY_TO_APPLY"],"recruiters":["All recruiters"],"territories":["All territories"]}'::JSONB, '{"audienceSize":342,"sent":342,"delivered":333,"opened":0,"clicked":96,"bounced":9,"replied":27}'::JSONB, '2026-05-15 00:00:00', '2026-08-01 23:59:00', '5cf37266-3473-4006-984f-9325122678b7', '2026-05-10 09:00:00', '2026-05-18 08:10:00'),
	('c3000000-0000-4000-8000-000000000003', 'Financial Aid Deadline', 'DRAFT', 'EMAIL', 'Admitted students without FAFSA', 'Aid deadline reminder', 'Financial aid deadlines are approaching. Complete your next steps before July 1.', '{"applicationTypes":["KUCCPS_PLACEMENT","SELF_SPONSORED_UNDERGRAD"],"applicationStatuses":["ADMITTED"],"academicTerms":["September 2026"],"programs":["All active programs"],"eventAttendance":"ATTENDED","leadScoreBands":["READY_TO_APPLY"],"recruiters":["Maya Schultz"],"territories":["Nairobi","Mombasa"]}'::JSONB, '{"audienceSize":890,"sent":0,"delivered":0,"opened":0,"clicked":0,"bounced":0,"replied":0}'::JSONB, '2026-06-15 00:00:00', '2026-07-01 23:59:00', 'b7f1d86f-0f1f-4c7b-84d7-99a570f14b6f', '2026-05-22 14:30:00', '2026-05-22 14:30:00'),
	('c3000000-0000-4000-8000-000000000004', 'Spring 2026 Yield Campaign', 'COMPLETED', 'EMAIL', 'Admitted not enrolled', 'Yield story sequence', 'Hear from current students and confirm your enrollment intent before the deadline.', '{"applicationTypes":["KUCCPS_PLACEMENT","DIPLOMA","MASTERS"],"applicationStatuses":["ADMITTED"],"academicTerms":["January 2026"],"programs":["All active programs"],"eventAttendance":"REGISTERED","leadScoreBands":["HOT","READY_TO_APPLY"],"recruiters":["All recruiters"],"territories":["All territories"]}'::JSONB, '{"audienceSize":1200,"sent":1200,"delivered":1182,"opened":733,"clicked":216,"bounced":18,"replied":48}'::JSONB, '2026-03-01 00:00:00', '2026-05-01 17:00:00', NULL, '2026-02-12 09:00:00', '2026-05-01 17:00:00')
ON CONFLICT DO NOTHING;

INSERT INTO admissions_campaign_audit_events (campaign_audit_event_id, campaign_id, actor_name, action, occurred_at, date_created) VALUES
	('c3100000-0000-4000-8000-000000000001', 'c3000000-0000-4000-8000-000000000001', 'Maya Schultz', 'Activated campaign', '2026-05-01 09:15:00', '2026-05-01 09:15:00'),
	('c3100000-0000-4000-8000-000000000002', 'c3000000-0000-4000-8000-000000000001', 'System', 'Queued first send batch', '2026-05-01 09:20:00', '2026-05-01 09:20:00'),
	('c3100000-0000-4000-8000-000000000003', 'c3000000-0000-4000-8000-000000000002', 'System', 'Audience refreshed from checklist status', '2026-05-18 07:00:00', '2026-05-18 07:00:00'),
	('c3100000-0000-4000-8000-000000000004', 'c3000000-0000-4000-8000-000000000002', 'Andre Park', 'Approved SMS send', '2026-05-18 08:10:00', '2026-05-18 08:10:00'),
	('c3100000-0000-4000-8000-000000000005', 'c3000000-0000-4000-8000-000000000003', 'Maya Schultz', 'Drafted audience and message', '2026-05-22 14:30:00', '2026-05-22 14:30:00'),
	('c3100000-0000-4000-8000-000000000006', 'c3000000-0000-4000-8000-000000000004', 'System', 'Completed campaign schedule', '2026-05-01 17:00:00', '2026-05-01 17:00:00')
ON CONFLICT DO NOTHING;

INSERT INTO admissions_communications (communication_id, external_message_id, channel, direction, constituent_id, application_id, campaign_id, recipient_sender, recipient_initials, subject, preview, status, provider, owner_name, outcome, duration, occurred_at, provider_payload, date_created, date_updated) VALUES
	('c3200000-0000-4000-8000-000000000001', 'MSG-8492', 'EMAIL', 'OUTBOUND', 'a3000000-0000-4000-8000-000000000001', 'd3000000-0000-4000-8000-000000000001', 'c3000000-0000-4000-8000-000000000001', 'Sofia Martinez', 'SM', 'Application status update', 'Your Fall 2026 application has moved into first review with the admissions team.', 'OPENED', 'SendGrid', 'Maya Reynolds', NULL, NULL, '2026-06-05 10:30:00', '{"providerMessageId":"sg-8492","openedAt":"2026-06-05T10:42:00Z"}'::JSONB, '2026-06-05 10:30:00', '2026-06-05 10:42:00'),
	('c3200000-0000-4000-8000-000000000002', 'MSG-8491', 'SMS', 'OUTBOUND', 'a3000000-0000-4000-8000-000000000002', 'd3000000-0000-4000-8000-000000000002', 'c3000000-0000-4000-8000-000000000002', 'James Okoro', 'JO', 'Missing transcript reminder', 'Reminder: upload your corrected transcript by May 31 to keep your application on schedule.', 'DELIVERED', 'Twilio', 'Admissions Automation', NULL, NULL, '2026-06-05 09:15:00', '{"providerMessageId":"tw-8491","deliveredAt":"2026-06-05T09:15:42Z"}'::JSONB, '2026-06-05 09:15:00', '2026-06-05 09:15:42'),
	('c3200000-0000-4000-8000-000000000003', 'MSG-8490-WA', 'WHATSAPP', 'OUTBOUND', 'a3000000-0000-4000-8000-000000000003', NULL, 'c3000000-0000-4000-8000-000000000002', 'Achieng Otieno', 'AO', 'Document verification reminder', 'WhatsApp template reminder: please confirm your KCSE result slip and National ID details before review closes.', 'DELIVERED', 'WhatsApp Cloud', 'Admissions Automation', NULL, NULL, '2026-06-05 08:05:00', '{"providerMessageId":"wa-8490","template":"document_verification"}'::JSONB, '2026-06-05 08:05:00', '2026-06-05 08:05:31'),
	('c3200000-0000-4000-8000-000000000004', 'MSG-8490', 'PHONE_CALL', 'OUTBOUND', 'a3000000-0000-4000-8000-000000000004', 'd3000000-0000-4000-8000-000000000003', NULL, 'Priya Patel', 'PP', 'Interview confirmation call', 'Confirmed virtual interview time and reviewed scholarship documentation expectations.', 'LOGGED', NULL, 'Noah Williams', 'Interview confirmed', '12 min', '2026-06-04 16:40:00', NULL, '2026-06-04 16:40:00', '2026-06-04 16:52:00'),
	('c3200000-0000-4000-8000-000000000005', 'MSG-8489', 'EMAIL', 'INBOUND', 'a3000000-0000-4000-8000-000000000005', 'd3000000-0000-4000-8000-000000000004', NULL, 'Liam Chen', 'LC', 'Counselor evaluation question', 'Applicant replied with counselor contact details for the final evaluation request.', 'REPLIED', 'SendGrid', 'Maya Reynolds', NULL, NULL, '2026-06-04 13:22:00', '{"providerMessageId":"sg-8489","inbound":true}'::JSONB, '2026-06-04 13:22:00', '2026-06-04 13:22:00'),
	('c3200000-0000-4000-8000-000000000006', 'MSG-8488', 'EMAIL', 'OUTBOUND', 'a3000000-0000-4000-8000-000000000006', 'd3000000-0000-4000-8000-000000000005', 'c3000000-0000-4000-8000-000000000004', 'Aisha Bello', 'AB', 'Welcome to the Fall 2026 cohort', 'Admission decision package failed to send because the recipient mailbox bounced.', 'BOUNCED', 'SendGrid', 'Admissions Automation', NULL, NULL, '2026-05-29 08:05:00', '{"providerMessageId":"sg-8488","bounceCode":"mailbox_unavailable"}'::JSONB, '2026-05-29 08:05:00', '2026-05-29 08:06:00'),
	('c3200000-0000-4000-8000-000000000007', 'MSG-8487', 'NOTIFICATION', 'OUTBOUND', 'a3000000-0000-4000-8000-000000000007', 'd3000000-0000-4000-8000-000000000006', NULL, 'Amara Ndlovu', 'AN', 'Decision posted in applicant portal', 'In-app decision notification was created when the reviewer published the final decision.', 'DELIVERED', NULL, 'System', NULL, NULL, '2026-05-28 15:15:00', NULL, '2026-05-28 15:15:00', '2026-05-28 15:15:00'),
	('c3200000-0000-4000-8000-000000000008', 'MSG-8486', 'SMS', 'OUTBOUND', 'a3000000-0000-4000-8000-000000000008', NULL, 'c3000000-0000-4000-8000-000000000001', 'Noah Williams', 'NW', 'Campus tour reminder', 'Reminder text for Saturday campus tour failed after provider rate limiting.', 'FAILED', 'Twilio', 'Admissions Automation', 'Retry scheduled', NULL, '2026-05-27 17:45:00', '{"providerMessageId":"tw-8486","errorCode":"rate_limited"}'::JSONB, '2026-05-27 17:45:00', '2026-05-27 17:46:00'),
	('c3200000-0000-4000-8000-000000000009', 'MSG-8486-WA', 'WHATSAPP', 'INBOUND', 'a3000000-0000-4000-8000-000000000009', 'd3000000-0000-4000-8000-000000000007', NULL, 'Wanjiku Njoroge', 'WN', 'Applicant replied on WhatsApp', 'Applicant confirmed SMS is preferred for urgent alerts and WhatsApp is available for document nudges.', 'REPLIED', 'WhatsApp Cloud', 'Owen Adirah', 'Preference confirmed', NULL, '2026-05-27 16:20:00', '{"providerMessageId":"wa-8486","inbound":true}'::JSONB, '2026-05-27 16:20:00', '2026-05-27 16:20:00'),
	('c3200000-0000-4000-8000-000000000010', 'MSG-8485', 'PHONE_CALL', 'INBOUND', 'a3000000-0000-4000-8000-000000000010', 'd3000000-0000-4000-8000-000000000008', NULL, 'Emma Davis', 'ED', 'Financial aid follow-up', 'Parent called to confirm financial aid document deadlines and next steps.', 'LOGGED', NULL, 'Owen Adirah', 'Follow-up email requested', '18 min', '2026-05-26 11:20:00', NULL, '2026-05-26 11:20:00', '2026-05-26 11:38:00')
ON CONFLICT DO NOTHING;

INSERT INTO admissions_sync_jobs (sync_job_id, name, status, direction, started_at, completed_at, records_pulled, records_pushed, events_requeued, failure_reason, retryable, created_by_id, date_created, date_updated, adapter, operation, attempt_count, max_attempts, next_retry_at, external_ref, external_receipt_id, error_code, error_detail, last_error_at) VALUES
	('c3300000-0000-4000-8000-000000000001', 'KUCCPS placement pull', 'SUCCEEDED', 'INBOUND', '2026-06-05 06:00:00', '2026-06-05 06:03:00', 128, 0, 0, NULL, false, '5cf37266-3473-4006-984f-9325122678b7', '2026-06-05 06:00:00', '2026-06-05 06:03:00', 'kuccps', 'KUCCPS_PLACEMENT_PULL', 1, 3, NULL, 'KUCCPS-BATCH-20260605', NULL, NULL, NULL, NULL),
	('c3300000-0000-4000-8000-000000000002', 'KNEC result verification', 'RUNNING', 'INBOUND', '2026-06-05 07:45:00', NULL, 42, 0, 0, NULL, false, '5cf37266-3473-4006-984f-9325122678b7', '2026-06-05 07:45:00', '2026-06-05 07:45:00', 'knec', 'KNEC_RESULT_VERIFICATION', 1, 3, NULL, 'KNEC-BATCH-20260605', NULL, NULL, NULL, NULL),
	('c3300000-0000-4000-8000-000000000003', 'Twilio SMS delivery retry', 'RETRY_READY', 'OUTBOUND', '2026-05-27 17:45:00', NULL, 0, 31, 1, 'Provider rate limit exceeded', true, '5cf37266-3473-4006-984f-9325122678b7', '2026-05-27 17:45:00', '2026-05-27 17:46:00', 'celcom_africa', 'SMS_OUTBOUND', 2, 3, '2026-05-27 18:15:00', 'SMS-BATCH-8486', NULL, 'RATE_LIMITED', 'Provider rate limit exceeded for tour reminder batch.', '2026-05-27 17:46:00'),
	('c3300000-0000-4000-8000-000000000004', 'WhatsApp Cloud delivery webhook', 'FAILED', 'INBOUND', '2026-05-27 16:20:00', '2026-05-27 16:21:00', 0, 0, 0, 'Webhook signature validation failed', true, NULL, '2026-05-27 16:20:00', '2026-05-27 16:21:00', 'whatsapp_cloud', 'WHATSAPP_WEBHOOK_INBOUND', 3, 3, NULL, 'WA-WEBHOOK-8486', NULL, 'INVALID_SIGNATURE', 'Webhook signature validation failed for inbound callback replay.', '2026-05-27 16:21:00')
ON CONFLICT DO NOTHING;

INSERT INTO admissions_sync_events (sync_event_id, sync_job_id, event_type, status, direction, resource_type, resource_id, payload_hash, attempts, next_retry_at, failure_reason, audit_message, date_created, date_updated, adapter, operation, max_attempts, external_ref, external_receipt_id, error_code, error_detail, last_error_at) VALUES
	('c3400000-0000-4000-8000-000000000001', 'c3300000-0000-4000-8000-000000000001', 'BATCH_ENROLLMENT_PULL', 'SUCCEEDED', 'INBOUND', 'application', 'd3000000-0000-4000-8000-000000000001', 'seed-kuccps-placement-001', 1, NULL, NULL, 'KUCCPS placement pulled for Sofia Martinez.', '2026-06-05 06:01:00', '2026-06-05 06:03:00', 'kuccps', 'KUCCPS_PLACEMENT_PULL', 3, 'KUCCPS-BATCH-20260605', NULL, NULL, NULL, NULL),
	('c3400000-0000-4000-8000-000000000002', 'c3300000-0000-4000-8000-000000000002', 'BATCH_PERSON_MATCHES_PULL', 'PROCESSING', 'INBOUND', 'constituent', 'a3000000-0000-4000-8000-000000000009', 'seed-knec-verification-001', 1, NULL, NULL, 'KNEC result verification running for Wanjiku Njoroge.', '2026-06-05 07:46:00', '2026-06-05 07:46:00', 'knec', 'KNEC_RESULT_VERIFICATION', 3, 'KNEC-BATCH-20260605', NULL, NULL, NULL, NULL),
	('c3400000-0000-4000-8000-000000000003', 'c3300000-0000-4000-8000-000000000003', 'SMS_OUTBOUND', 'RETRY_READY', 'OUTBOUND', 'communication', 'c3200000-0000-4000-8000-000000000008', 'seed-sms-failed-8486', 2, '2026-05-27 18:15:00', 'Provider rate limit exceeded', 'SMS campus tour reminder queued for retry.', '2026-05-27 17:45:00', '2026-05-27 17:46:00', 'celcom_africa', 'SMS_OUTBOUND', 3, 'SMS-BATCH-8486', NULL, 'RATE_LIMITED', 'Provider rate limit exceeded for tour reminder batch.', '2026-05-27 17:46:00'),
	('c3400000-0000-4000-8000-000000000004', 'c3300000-0000-4000-8000-000000000004', 'WHATSAPP_WEBHOOK_INBOUND', 'FAILED', 'INBOUND', 'communication', 'c3200000-0000-4000-8000-000000000009', 'seed-whatsapp-inbound-8486', 3, NULL, 'Webhook signature validation failed', 'WhatsApp inbound webhook failed signature validation.', '2026-05-27 16:20:00', '2026-05-27 16:21:00', 'whatsapp_cloud', 'WHATSAPP_WEBHOOK_INBOUND', 3, 'WA-WEBHOOK-8486', NULL, 'INVALID_SIGNATURE', 'Webhook signature validation failed for inbound callback replay.', '2026-05-27 16:21:00')
ON CONFLICT DO NOTHING;

INSERT INTO audit (id, obj_id, obj_domain, obj_name, actor_id, action, data, message, timestamp) VALUES
	('c3500000-0000-4000-8000-000000000001', 'a3000000-0000-4000-8000-000000000001', 'admissions_constituent', 'Sofia Martinez', '5cf37266-3473-4006-984f-9325122678b7', 'CREATE', '{"source":"seed","lifecycleStage":"APPLICANT"}'::JSONB, 'Seeded applicant profile for admin constituent list.', '2026-05-20 09:00:00'),
	('c3500000-0000-4000-8000-000000000002', 'd3000000-0000-4000-8000-000000000001', 'admissions_application', 'Sofia Martinez - Bachelor of Commerce', 'c41fa5d3-d61f-45f1-b054-d2c7a3704019', 'TRANSITION', '{"from":"SUBMITTED","to":"IN_REVIEW"}'::JSONB, 'Application moved into first review.', '2026-05-28 08:00:00'),
	('c3500000-0000-4000-8000-000000000003', 'c3000000-0000-4000-8000-000000000001', 'admissions_campaign', 'Fall 2026 Open House Invite', '5cf37266-3473-4006-984f-9325122678b7', 'UPDATE', '{"status":"ACTIVE","channel":"EMAIL"}'::JSONB, 'Campaign activated and first batch queued.', '2026-05-01 09:15:00'),
	('c3500000-0000-4000-8000-000000000004', 'c3200000-0000-4000-8000-000000000008', 'admissions_communication', 'Campus tour reminder', '5cf37266-3473-4006-984f-9325122678b7', 'SYNC_PUSH', '{"channel":"SMS","status":"FAILED","provider":"Twilio"}'::JSONB, 'SMS provider rate limit created retry-ready sync event.', '2026-05-27 17:46:00'),
	('c3500000-0000-4000-8000-000000000005', 'c3300000-0000-4000-8000-000000000001', 'admissions_sync_job', 'KUCCPS placement pull', '5cf37266-3473-4006-984f-9325122678b7', 'SYNC_PULL', '{"adapter":"kuccps","status":"SUCCEEDED","recordsPulled":128}'::JSONB, 'KUCCPS placement sync completed successfully.', '2026-06-05 06:03:00'),
	('c3500000-0000-4000-8000-000000000006', 'e3000000-0000-4000-8000-000000000001', 'admissions_inquiry', 'Achieng Otieno inquiry', 'b7f1d86f-0f1f-4c7b-84d7-99a570f14b6f', 'CREATE', '{"source":"portal","status":"NEW"}'::JSONB, 'Portal inquiry created for September intake scholarship follow-up.', '2026-05-18 11:00:00')
ON CONFLICT DO NOTHING;

INSERT INTO admissions_lead_score_rules (lead_score_rule_id, name, description, criteria, points, is_active, priority, date_created, date_updated) VALUES
	('38c7b57f-0c95-4dd7-9a08-66df08c20f1d', 'Applicant lifecycle stage', 'Applicants are warmer than prospects.', '[{"Field":"lifecycle_stage","Operator":"EQ","Values":["APPLICANT"]}]', 30, true, 10, '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('62d0d3e7-58e5-4d97-865a-7a7db50c9e61', 'Submitted application', 'Submitted applications show high intent.', '[{"Field":"application_status","Operator":"IN","Values":["SUBMITTED","READY_FOR_REVIEW","IN_REVIEW"]}]', 40, true, 20, '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('c424b3dd-5e18-4118-9448-8876acdc42bd', 'Transfer or graduate interest', 'Transfer and graduate prospects often have urgent timelines.', '[{"Field":"application_type","Operator":"IN","Values":["DIPLOMA","MASTERS"]}]', 15, true, 30, '2019-03-24 00:00:00', '2019-03-24 00:00:00')
ON CONFLICT DO NOTHING;

-- Kenya reference seed data adapted from njoguamos/kenya-demographics-units, cross-referenced against HDX COD-AB Kenya research notes.
INSERT INTO counties (code, name, date_created, date_updated) VALUES
	('1', 'Mombasa', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('2', 'Kwale', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('3', 'Kilifi', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('4', 'Tana River', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('5', 'Lamu', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('6', 'Taita Taveta', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('7', 'Garissa', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('8', 'Wajir', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('9', 'Mandera', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('10', 'Marsabit', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('11', 'Isiolo', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('12', 'Meru', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('13', 'Tharaka Nithi', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('14', 'Embu', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('15', 'Kitui', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('16', 'Machakos', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('17', 'Makueni', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('18', 'Nyandarua', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('19', 'Nyeri', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('20', 'Kirinyaga', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('21', 'Murang''a', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('22', 'Kiambu', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('23', 'Turkana', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('24', 'West Pokot', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('25', 'Samburu', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('26', 'Trans Nzoia', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('27', 'Uasin Gishu', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('28', 'Elgeyo Marakwet', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('29', 'Nandi', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('30', 'Baringo', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('31', 'Laikipia', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('32', 'Nakuru', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('33', 'Narok', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('34', 'Kajiado', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('35', 'Kericho', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('36', 'Bomet', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('37', 'Kakamega', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('38', 'Vihiga', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('39', 'Bungoma', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('40', 'Busia', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('41', 'Siaya', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('42', 'Kisumu', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('43', 'Homa Bay', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('44', 'Migori', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('45', 'Kisii', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('46', 'Nyamira', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('47', 'Nairobi', '2019-03-24 00:00:00', '2019-03-24 00:00:00')
ON CONFLICT DO NOTHING;

INSERT INTO sub_counties (code, county_code, name, date_created, date_updated) VALUES
	('1', '1', 'Changamwe', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('2', '1', 'Jomvu', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('3', '1', 'Kisauni', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('4', '1', 'Nyali', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('5', '1', 'Likoni', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('6', '1', 'Mvita', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('7', '2', 'Msambweni', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('8', '2', 'Lunga Lunga', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('9', '2', 'Kwale', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('10', '2', 'Kinango', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('11', '3', 'Bahari(Kilifi)', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('12', '3', 'Kilifi South', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('13', '3', 'Kaloleni', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('14', '3', 'Rabai', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('15', '3', 'Ganze', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('16', '3', 'Malindi', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('17', '3', 'Magarini', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('18', '4', 'Tana Delta', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('19', '4', 'Tana River', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('20', '4', 'Bura(Tana North)', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('21', '5', 'Lamu East', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('22', '5', 'Lamu West', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('23', '6', 'Taveta', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('24', '6', 'Wundanyi(Taita)', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('25', '6', 'Mwatate', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('26', '6', 'Voi', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('27', '7', 'Hulugho', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('28', '7', 'Ijara', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('29', '7', 'Balambala', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('30', '7', 'Lagdera', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('31', '7', 'Dadaab', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('32', '7', 'Fafi', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('33', '7', 'Ijara', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('34', '8', 'Wajir North', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('35', '8', 'Wajir East', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('36', '8', 'Buna', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('37', '8', 'Habaswein', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('38', '8', 'Tarbaj', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('39', '8', 'Wajir West', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('40', '8', 'Eldas', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('41', '8', 'Wajir South', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('42', '9', 'Mandera West', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('43', '9', 'Banisa', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('44', '9', 'Mandera North', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('45', '9', 'Mandera Central', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('46', '9', 'Mandera East', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('47', '9', 'Lafey', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('48', '10', 'Moyale', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('49', '10', 'Marsabit', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('50', '10', 'Horr North', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('51', '10', 'Loiyangalani', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('52', '10', 'Chalbi', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('53', '10', 'Sololo', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('54', '10', 'Marsabit South(Laisamis)', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('55', '11', 'Garbatula', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('56', '11', 'Isiolo', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('57', '11', 'Merti', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('58', '12', 'Igembe South', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('59', '12', 'Igembe Central', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('60', '12', 'Igembe North', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('61', '12', 'Imenti South', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('62', '12', 'Imenti North', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('63', '12', 'Meru Central', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('64', '12', 'Tigania Central', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('65', '12', 'Tigania East', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('66', '12', 'Tigania West', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('67', '12', 'Buuri', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('68', '13', 'Meru South', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('69', '13', 'Maara', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('70', '13', 'Tharaka South', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('71', '13', 'Tharaka North', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('72', '14', 'Embu East', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('73', '14', 'Embu North', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('74', '14', 'Embu West', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('75', '14', 'Mbeere North', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('76', '14', 'Mbeere South', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('77', '15', 'Mwingi Central', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('78', '15', 'Mwingi West', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('79', '15', 'Mwingi East', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('80', '15', 'Kitui West', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('81', '15', 'Kitui Central', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('82', '15', 'Nzambani', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('83', '15', 'Mutomo', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('84', '15', 'Mutomo', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('85', '15', 'Ikutha', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('86', '15', 'Katulani', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('87', '15', 'Kisasi', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('88', '15', 'Kyuso', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('89', '15', 'Lower Yatta', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('90', '15', 'Matinyani', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('91', '15', 'Mumoni', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('92', '15', 'Mutito', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('93', '16', 'Masinga', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('94', '16', 'Yatta', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('95', '16', 'Kangundo', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('96', '16', 'Matungulu', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('97', '16', 'Kathiani', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('98', '16', 'Athi River', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('99', '16', 'Machakos', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('100', '16', 'Mwala', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('101', '17', 'Kibwezi', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('102', '17', 'Kilungu', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('103', '17', 'Makindu', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('104', '17', 'Makueni', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('105', '17', 'Mbooni West', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('106', '17', 'Mbooni East', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('107', '17', 'Kathonzweni', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('108', '17', 'Mukaa', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('109', '17', 'Nzaui', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('110', '18', 'Kinangop', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('111', '18', 'Kipipiri', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('112', '18', 'Mirangine', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('113', '18', 'Nyandarua Central', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('114', '18', 'Nyandarua North', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('115', '18', 'Nyandarua South', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('116', '18', 'Nyandarua West', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('117', '19', 'Tetu', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('118', '19', 'Kieni East', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('119', '19', 'Kieni West', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('120', '19', 'Mathira East', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('121', '19', 'Mathira West', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('122', '19', 'Mukurwe-ini', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('123', '19', 'Nyeri Central', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('124', '19', 'Nyeri South', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('125', '20', 'Kirinyaga Central', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('126', '20', 'Kirinyaga East', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('127', '20', 'Kirinyaga West', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('128', '20', 'Mwea East', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('129', '20', 'Mwea West', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('130', '21', 'Kangema', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('131', '21', 'Mathioya', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('132', '21', 'Kahuro', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('133', '21', 'Kigumo', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('134', '21', 'Murang''a East', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('135', '21', 'Kandara', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('136', '21', 'Gatanga', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('137', '21', 'Murang''a South', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('138', '22', 'Gatundu South', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('139', '22', 'Gatundu North', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('140', '22', 'Juja', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('141', '22', 'Ruiru', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('142', '22', 'Githunguri', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('143', '22', 'Kiambu', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('144', '22', 'Kiambaa', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('145', '22', 'Kabete', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('146', '22', 'Kikuyu', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('147', '22', 'Limuru', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('148', '22', 'Lari', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('149', '22', 'Thika East', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('150', '22', 'Thika West', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('151', '23', 'Turkana North', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('152', '23', 'Turkana West', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('153', '23', 'Turkana Central', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('154', '23', 'Loima', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('155', '23', 'Turkana South', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('156', '23', 'Turkana East', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('157', '23', 'Kibish', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('158', '24', 'Kipkomo', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('159', '24', 'Pokot Central', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('160', '24', 'Pokot North', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('161', '24', 'Pokot South', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('162', '24', 'West Pokot', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('163', '25', 'Samburu Central', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('164', '25', 'Samburu North', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('165', '25', 'Samburu East', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('166', '26', 'Kwanza', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('167', '26', 'Endebes', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('168', '26', 'Transzoia East', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('169', '26', 'Kiminini', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('170', '26', 'Transzoia West', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('171', '27', 'Soy', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('172', '27', 'Wareng', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('173', '27', 'Moiben', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('174', '27', 'Eldoret West', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('175', '27', 'Eldoret East', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('176', '27', 'Kesses', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('177', '28', 'Marakwet East', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('178', '28', 'Marakwet West', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('179', '28', 'Keiyo', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('180', '28', 'Keiyo South', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('181', '29', 'Tinderet', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('182', '29', 'Nandi Central', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('183', '29', 'Nandi East', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('184', '29', 'Chesumei', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('185', '29', 'Nandi North', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('186', '29', 'Nandi South', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('187', '30', 'East Pokot', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('188', '30', 'Baringo North', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('189', '30', 'Baringo Central', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('190', '30', 'Koibatek', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('191', '30', 'Marigat', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('192', '30', 'Mogotio', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('193', '31', 'Laikipia West', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('194', '31', 'Laikipia East', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('195', '31', 'Laikipia North', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('196', '31', 'Laikipia Central', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('197', '31', 'Nyahururu', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('198', '32', 'Molo', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('199', '32', 'Njoro', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('200', '32', 'Naivasha', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('201', '32', 'Gilgil', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('202', '32', 'Kuresoi', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('203', '32', 'Kuresoi North', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('204', '32', 'Subukia', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('205', '32', 'Rongai', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('206', '32', 'Nakuru', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('207', '32', 'Nakuru West', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('208', '32', 'Nakuru North', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('209', '33', 'Narok North', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('210', '33', 'Narok East', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('211', '33', 'Narok South', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('212', '33', 'Narok West', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('213', '33', 'Transmara East', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('214', '33', 'Transmara West', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('215', '34', 'Kajiado North', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('216', '34', 'Kajiado Central', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('217', '34', 'Kajiado West', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('218', '34', 'Isinya', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('219', '34', 'Loitokitok', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('220', '34', 'Mashuuru', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('221', '35', 'Kipkelion', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('222', '35', 'Kericho', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('223', '35', 'Londiani', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('224', '35', 'Bureti', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('225', '35', 'Belgut', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('226', '35', 'Sigowei/Soin', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('227', '36', 'Sotik', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('228', '36', 'Chepalungu', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('229', '36', 'Bomet East', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('230', '36', 'Bomet', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('231', '36', 'Konoin', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('232', '37', 'Lugari', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('233', '37', 'Matete', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('234', '37', 'Likuyani', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('235', '37', 'Kakamega Central', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('236', '37', 'Kakamega East', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('237', '37', 'Navakholo', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('238', '37', 'Mumias', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('239', '37', 'Mumias East', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('240', '37', 'Matungu', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('241', '37', 'Butere', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('242', '37', 'Khwisero', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('243', '37', 'Kakamega North', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('244', '37', 'Kakamega South', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('245', '38', 'Vihiga', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('246', '38', 'Sabatia', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('247', '38', 'Hamisi', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('248', '38', 'Luanda', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('249', '38', 'Emuhaya', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('250', '39', 'Mt. Elgon', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('251', '39', 'Bungoma Central', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('252', '39', 'Bungoma East', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('253', '39', 'Bungoma North', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('254', '39', 'Bungoma South', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('255', '39', 'Bungoma West', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('256', '39', 'Webuye West', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('257', '39', 'Cheptais', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('258', '39', 'Kimilili', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('259', '39', 'Bumula', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('260', '40', 'Teso North', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('261', '40', 'Teso South', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('262', '40', 'Nambale', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('263', '40', 'Bunyala', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('264', '40', 'Busia', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('265', '40', 'Butula', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('266', '40', 'Samia', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('267', '41', 'Ugenya', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('268', '41', 'Ugunja', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('269', '41', 'Siaya', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('270', '41', 'Gem', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('271', '41', 'Bondo', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('272', '41', 'Rarieda', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('273', '42', 'Kisumu East', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('274', '42', 'Kisumu West', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('275', '42', 'Kisumu Central', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('276', '42', 'Seme', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('277', '42', 'Nyando', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('278', '42', 'Muhoroni', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('279', '42', 'Nyakach', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('280', '43', 'Mbita', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('281', '43', 'Rachuonyo East', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('282', '43', 'Rachuonyo North', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('283', '43', 'Rachuonyo South', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('284', '43', 'Homa Bay', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('285', '43', 'Ndhiwa', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('286', '43', 'Rangwe', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('287', '43', 'Suba', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('288', '44', 'Rongo', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('289', '44', 'Awendo', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('290', '44', 'Migori', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('291', '44', 'Suna West', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('292', '44', 'Uriri', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('293', '44', 'Nyatike', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('294', '44', 'Kuria West', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('295', '44', 'Kuria East', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('296', '45', 'Gucha', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('297', '45', 'Gucha South', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('298', '45', 'Kenyenya', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('299', '45', 'Kisii Central', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('300', '45', 'Kisii South', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('301', '45', 'Kitutu Central', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('302', '45', 'Marani', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('303', '45', 'Masaba south', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('304', '45', 'Nyamache', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('305', '45', 'Sameta', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('306', '46', 'Manga', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('307', '46', 'Masaba North', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('308', '46', 'Borabu', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('309', '46', 'Nyamira North', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('310', '46', 'Nyamira South', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('311', '47', 'Westlands', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('312', '47', 'Dagoretti', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('313', '47', 'Langata', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('314', '47', 'Kibra', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('315', '47', 'Kasarani', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('316', '47', 'Embakasi', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('317', '47', 'Makadara', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('318', '47', 'Kamukunji', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('319', '47', 'Starehe', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('320', '47', 'Mathare', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('321', '47', 'Njiru', '2019-03-24 00:00:00', '2019-03-24 00:00:00')
ON CONFLICT DO NOTHING;

-- Kenya ward seed slice adapted from HDX COD-AB/IEBC ward references. Full national ward import remains source-versioned operational data.
INSERT INTO wards (code, county_code, sub_county_code, name, date_created, date_updated) VALUES
	('W001', '1', '1', 'Port Reitz', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('W002', '1', '1', 'Kipevu', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('W003', '1', '1', 'Airport', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('W004', '1', '2', 'Jomvu Kuu', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('W005', '1', '2', 'Miritini', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('W006', '47', '312', 'Parklands/Highridge', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('W007', '47', '312', 'Kangemi', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('W008', '47', '312', 'Mountain View', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('W009', '47', '319', 'Nairobi Central', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('W010', '47', '319', 'Ngara', '2019-03-24 00:00:00', '2019-03-24 00:00:00')
ON CONFLICT DO NOTHING;

-- Kenya university seed slice aligned with Commission for University Education institution categories.
INSERT INTO universities (code, name, institution_type, date_created, date_updated) VALUES
	('UON', 'University of Nairobi', 'PUBLIC_UNIVERSITY', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('KU', 'Kenyatta University', 'PUBLIC_UNIVERSITY', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('JKUAT', 'Jomo Kenyatta University of Agriculture and Technology', 'PUBLIC_UNIVERSITY', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('MKU', 'Mount Kenya University', 'PRIVATE_UNIVERSITY', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('TUK', 'Technical University of Kenya', 'PUBLIC_UNIVERSITY', '2019-03-24 00:00:00', '2019-03-24 00:00:00')
ON CONFLICT DO NOTHING;

-- KUCCPS programme cluster seed slice for admissions lookup and routing tests.
INSERT INTO programme_clusters (code, name, description, date_created, date_updated) VALUES
	('CL01', 'Medicine and Health Sciences', 'Programmes whose placement cluster emphasizes biology, chemistry, mathematics or physics.', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('CL02', 'Engineering and Technology', 'Programmes whose placement cluster emphasizes mathematics, physics and technical sciences.', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('CL03', 'Education and Humanities', 'Programmes whose placement cluster emphasizes languages, humanities and teaching subjects.', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('CL04', 'Business and Social Sciences', 'Programmes whose placement cluster emphasizes business, economics and social sciences.', '2019-03-24 00:00:00', '2019-03-24 00:00:00')
ON CONFLICT DO NOTHING;

-- Kenya National Qualifications Framework levels 1-10.
INSERT INTO knqf_levels (code, level, name, descriptor, qualification, date_created, date_updated) VALUES
	('KNQF-1', 1, 'KNQF Level 1', 'Basic introductory knowledge and skills.', 'National Vocational Certificate I', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('KNQF-2', 2, 'KNQF Level 2', 'Basic operational knowledge and skills.', 'National Vocational Certificate II', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('KNQF-3', 3, 'KNQF Level 3', 'Intermediate operational knowledge and skills.', 'Artisan Certificate', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('KNQF-4', 4, 'KNQF Level 4', 'Technical knowledge and supervised practice.', 'Craft Certificate', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('KNQF-5', 5, 'KNQF Level 5', 'Specialized technical knowledge and practice.', 'Diploma', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('KNQF-6', 6, 'KNQF Level 6', 'Advanced technical or professional knowledge.', 'Higher Diploma', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('KNQF-7', 7, 'KNQF Level 7', 'Broad professional knowledge and analytical skills.', 'Bachelor Degree', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('KNQF-8', 8, 'KNQF Level 8', 'Advanced professional and research knowledge.', 'Postgraduate Diploma', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('KNQF-9', 9, 'KNQF Level 9', 'Specialized research and professional leadership.', 'Masters Degree', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('KNQF-10', 10, 'KNQF Level 10', 'Original research and highest-level expertise.', 'Doctoral Degree', '2019-03-24 00:00:00', '2019-03-24 00:00:00')
ON CONFLICT DO NOTHING;

-- Kenya programme seed slice aligned with KUCCPS/CUE-backed public catalog patterns.
INSERT INTO programmes (code, university_code, cluster_code, knqf_level_code, name, award_type, date_created, date_updated) VALUES
	('UON-BMED', 'UON', 'CL01', 'KNQF-7', 'Bachelor of Medicine and Bachelor of Surgery', 'BACHELOR', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('UON-BCOM', 'UON', 'CL04', 'KNQF-7', 'Bachelor of Commerce', 'BACHELOR', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('JKUAT-BENG-CIVIL', 'JKUAT', 'CL02', 'KNQF-7', 'Bachelor of Science in Civil Engineering', 'BACHELOR', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('KU-BED-ARTS', 'KU', 'CL03', 'KNQF-7', 'Bachelor of Education Arts', 'BACHELOR', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('TUK-DIP-ICT', 'TUK', 'CL02', 'KNQF-5', 'Diploma in Information Communication Technology', 'DIPLOMA', '2019-03-24 00:00:00', '2019-03-24 00:00:00')
ON CONFLICT DO NOTHING;
