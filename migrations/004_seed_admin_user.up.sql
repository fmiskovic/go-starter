INSERT INTO users (id, email, full_name, date_of_birth, location, gender)
VALUES
  (
    '110cea28-b2b0-4051-9eb6-9a99e451af02',
    'john.doe@example.com', 
    'John Doe',             
    '1990-01-01',           
    'New York, USA',        
    0
  );

INSERT INTO credentials (id, user_id, username, password_hash)
VALUES
  (
    '440cea28-b2b0-4051-9eb6-9a99e451af02',
    '110cea28-b2b0-4051-9eb6-9a99e451af02',           
    'admin',             
    '$2a$14$UoQvtVdIN3DdPiRXPSzOVuQJTxcMGbHY2a8euzqjJ/VvaJVOr3TbC'    
  );

INSERT INTO roles (id, user_id, name)
VALUES
  (
    '550cea28-b2b0-4051-9eb6-9a99e451af02',
    '110cea28-b2b0-4051-9eb6-9a99e451af02',
    'ROLE_ADMIN'
  );


