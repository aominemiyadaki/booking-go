INSERT INTO booking_bubuka.appointment_statuses (name) VALUES 
('в ожидании'),
('подтвержден'),
('отменен')
ON CONFLICT (name) DO NOTHING;
