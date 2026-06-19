-- Script de datos iniciales para ticket_system
-- Las tablas son creadas automáticamente por GORM AutoMigrate al iniciar el backend.
-- Pasos:
--   1. Arrancar el backend una vez (para que GORM cree las tablas)
--   2. Ejecutar: mysql -u root -p ticket_system < seed.sql

USE ticket_system;

-- Eventos de prueba
INSERT INTO events (titulo, descripcion, categoria, fecha, lugar, precio, cupo_maximo, cupo_disponible, created_at) VALUES
('Coldplay - Music of the Spheres',     'Gira mundial del grupo británico Coldplay.',             'Recitales', '2025-10-15 21:00:00', 'Estadio Monumental, Buenos Aires',  45000.00, 50000, 50000, NOW()),
('Wos en el Luna Park',                 'Show del rapero argentino Wos.',                         'Recitales', '2025-09-20 21:00:00', 'Luna Park, Buenos Aires',            8500.00,  5000,  5000,  NOW()),
('Romeo y Julieta',                     'Obra clásica de Shakespeare por el Teatro Colón.',       'Teatro',    '2025-08-10 20:00:00', 'Teatro Colón, Buenos Aires',         3500.00,  800,   800,   NOW()),
('El método Gronholm',                  'Thriller psicológico de Jordi Galcerán.',                'Teatro',    '2025-08-25 20:30:00', 'Teatro Broadway, Buenos Aires',      2800.00,  400,   400,   NOW()),
('River vs Boca - Superclásico',        'Fecha 18 del Torneo Apertura.',                         'Deportes',  '2025-11-02 16:00:00', 'Estadio Monumental, Buenos Aires',  25000.00, 84567, 84567, NOW()),
('Argentina vs Brasil - Eliminatorias', 'Eliminatorias Sudamericanas para el Mundial 2026.',      'Deportes',  '2025-11-15 21:00:00', 'Estadio Monumental, Buenos Aires',  35000.00, 84567, 84567, NOW()),
('Festival Lollapalooza Argentina',     'Festival multigenero con artistas internacionales.',     'Recitales', '2026-03-20 14:00:00', 'Hipódromo de San Isidro',            55000.00, 70000, 70000, NOW()),
('Stand Up: Malena Guinzburg',          'Show de stand up comedy de Malena Guinzburg.',           'Teatro',    '2025-09-05 21:30:00', 'Teatro Astros, Buenos Aires',        1800.00,  300,   300,   NOW());

-- Usuarios admin (contraseñas hasheadas con SHA-256, igual que el backend)
INSERT INTO users (nombre, apellido, email, password, dni, rol, created_at) VALUES
('Simon',  'Admin', '2408192@ucc.edu.ar', SHA2('Boca07',     256), '47192801', 'admin', NOW()),
('Juan',   'Admin', '2413770@ucc.edu.ar', SHA2('Juanceto01', 256), '46377960', 'admin', NOW()),
('Facu',   'Admin', '2411006@ucc.edu.ar', SHA2('faqarri06',  256), '46450861', 'admin', NOW());
