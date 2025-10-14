-- +goose Up
-- +goose StatementBegin

-- Insert healthcare entitlements
INSERT INTO healthcare_entitlements (healthcare_entitlement) VALUES
('Social Security'),
('Universal Coverage'),
('Civil Servant Medical Benefit'),
('State Enterprise Employee');

-- Insert Users (3 patients, 3 doctors, 1 admin)
INSERT INTO users (id, password, first_name, last_name, gender, phone_number, role, created_at, updated_at) VALUES
-- Patients
('01920e5a-1234-7890-abcd-000000000001', '$argon2id$v=19$m=65536,t=3,p=4$KZPKCgTQjZUHK4A0aI72KQ$KoQcuUnR4HNcaO7364bX56V9KYE35S1I9LPQeQpXVbs', 'Somchai', 'Jaidee', 'male', '0812345678', 'patient', now(), now()),
('01920e5a-1234-7890-abcd-000000000002', '$argon2id$v=19$m=65536,t=3,p=4$KZPKCgTQjZUHK4A0aI72KQ$KoQcuUnR4HNcaO7364bX56V9KYE35S1I9LPQeQpXVbs', 'Sukanya', 'Sriprasert', 'female', '0823456789', 'patient', now(), now()),
('01920e5a-1234-7890-abcd-000000000003', '$argon2id$v=19$m=65536,t=3,p=4$KZPKCgTQjZUHK4A0aI72KQ$KoQcuUnR4HNcaO7364bX56V9KYE35S1I9LPQeQpXVbs', 'Nattapong', 'Wongsawat', 'male', '0834567890', 'patient', now(), now()),
-- Doctors
('01920e5a-1234-7890-abcd-000000000004', '$argon2id$v=19$m=65536,t=3,p=4$KZPKCgTQjZUHK4A0aI72KQ$KoQcuUnR4HNcaO7364bX56V9KYE35S1I9LPQeQpXVbs', 'Prasit', 'Tangkarnjanakul', 'male', '0845678901', 'doctor', now(), now()),
('01920e5a-1234-7890-abcd-000000000005', '$argon2id$v=19$m=65536,t=3,p=4$KZPKCgTQjZUHK4A0aI72KQ$KoQcuUnR4HNcaO7364bX56V9KYE35S1I9LPQeQpXVbs', 'Siriporn', 'Rattanakorn', 'female', '0856789012', 'doctor', now(), now()),
('01920e5a-1234-7890-abcd-000000000006', '$argon2id$v=19$m=65536,t=3,p=4$KZPKCgTQjZUHK4A0aI72KQ$KoQcuUnR4HNcaO7364bX56V9KYE35S1I9LPQeQpXVbs', 'Thanawat', 'Pongpanit', 'male', '0867890123', 'doctor', now(), now()),
-- Admin
('01920e5a-1234-7890-abcd-000000000007', '$argon2id$v=19$m=65536,t=3,p=4$KZPKCgTQjZUHK4A0aI72KQ$KoQcuUnR4HNcaO7364bX56V9KYE35S1I9LPQeQpXVbs', 'Admin', 'System', 'male', '0878901234', 'admin', now(), now());

-- Insert Patients
INSERT INTO patients (user_id, hospital_id, birth_date, id_card_number, address, allergies, emergency_contact, blood_type, created_at, updated_at) VALUES
('01920e5a-1234-7890-abcd-000000000001', 'HN001234', '1985-03-15', '1234567890123', '123 Sukhumvit Rd, Bangkok 10110', 'Penicillin', '0898765432', 'O+', now(), now()),
('01920e5a-1234-7890-abcd-000000000002', 'HN002345', '1990-07-22', '2345678901234', '456 Ratchadaphisek Rd, Bangkok 10400', 'None', '0887654321', 'A+', now(), now()),
('01920e5a-1234-7890-abcd-000000000003', 'HN003456', '1978-11-05', '3456789012345', '789 Phetkasem Rd, Bangkok 10160', 'Aspirin, Sulfa drugs', '0876543210', 'B+', now(), now());

-- Insert Patient Healthcare Entitlements
INSERT INTO user_healthcare_entitlement (patient_id, healthcare_entitlement) VALUES
('01920e5a-1234-7890-abcd-000000000001', 'Social Security'),
('01920e5a-1234-7890-abcd-000000000002', 'Universal Coverage'),
('01920e5a-1234-7890-abcd-000000000003', 'Civil Servant Medical Benefit');

-- Insert Doctors
INSERT INTO doctors (user_id, username, specialty, bio, years_experience, created_at, updated_at) VALUES
('01920e5a-1234-7890-abcd-000000000004', 'dr.prasit', 'Cardiology', 'Experienced cardiologist specializing in interventional procedures and heart disease management.', 15, now(), now()),
('01920e5a-1234-7890-abcd-000000000005', 'dr.siriporn', 'Pediatrics', 'Dedicated pediatrician with expertise in child development and preventive care.', 10, now(), now()),
('01920e5a-1234-7890-abcd-000000000006', 'dr.thanawat', 'Orthopedics', 'Orthopedic surgeon specializing in sports medicine and joint replacement.', 12, now(), now());

-- Insert Admin
INSERT INTO admins (user_id, username) VALUES
('01920e5a-1234-7890-abcd-000000000007', 'admin');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Delete in reverse order of dependencies
DELETE FROM admins WHERE user_id IN (
  '01920e5a-1234-7890-abcd-000000000007'
);

DELETE FROM doctors WHERE user_id IN (
  '01920e5a-1234-7890-abcd-000000000004',
  '01920e5a-1234-7890-abcd-000000000005',
  '01920e5a-1234-7890-abcd-000000000006'
);

DELETE FROM user_healthcare_entitlement WHERE patient_id IN (
  '01920e5a-1234-7890-abcd-000000000001',
  '01920e5a-1234-7890-abcd-000000000002',
  '01920e5a-1234-7890-abcd-000000000003'
);

DELETE FROM patients WHERE user_id IN (
  '01920e5a-1234-7890-abcd-000000000001',
  '01920e5a-1234-7890-abcd-000000000002',
  '01920e5a-1234-7890-abcd-000000000003'
);

DELETE FROM users WHERE id IN (
  '01920e5a-1234-7890-abcd-000000000001',
  '01920e5a-1234-7890-abcd-000000000002',
  '01920e5a-1234-7890-abcd-000000000003',
  '01920e5a-1234-7890-abcd-000000000004',
  '01920e5a-1234-7890-abcd-000000000005',
  '01920e5a-1234-7890-abcd-000000000006',
  '01920e5a-1234-7890-abcd-000000000007'
);

DELETE FROM healthcare_entitlements WHERE healthcare_entitlement IN (
  'Social Security',
  'Universal Coverage',
  'Civil Servant Medical Benefit',
  'State Enterprise Employee'
);

-- +goose StatementEnd
