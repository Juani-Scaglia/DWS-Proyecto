-- ============================================================
-- SEED DATA - Establecimientos, Eventos y Asientos de prueba
-- Ejecutar: mysql -u root -p --default-character-set=utf8mb4 dws_proyecto < database/seed.sql
-- ============================================================

SET NAMES utf8mb4;

DELETE FROM tickets;
DELETE FROM seats;
DELETE FROM events;
DELETE FROM venues;

-- ESTABLECIMIENTOS (capacidades por sector)
INSERT INTO venues (id, nombre, direccion, capacidad, cap_platea_norte, cap_platea_sur, cap_tribuna_este, cap_tribuna_oeste, cap_platea_preferencial, cap_campo, created_at) VALUES
(1, 'Estadio Mario Alberto Kempes', 'Av. Carcano s/n, Cordoba',         300, 80, 80, 50, 50, 20, 20, NOW()),
(2, 'Orfeo Superdomo',              'Av. Fuerza Aerea 5200, Cordoba',    200, 60, 60, 40, 40, 0, 0, NOW()),
(3, 'Teatro Real de Cordoba',       'Av. Velez Sarsfield 365, Cordoba',   60, 30, 30, 0, 0, 0, 0, NOW()),
(4, 'Quality Espacio',              'Av. Cruz Roja Argentina 200, Cordoba',40, 20, 20, 0, 0, 0, 0, NOW()),
(5, 'Plaza de la Musica',           'Bv. Las Heras 80, Cordoba',          150, 40, 40, 30, 30, 0, 10, NOW());

-- EVENTOS
INSERT INTO events (id, titulo, descripcion, categoria, fecha, lugar, precio, imagen, cupo_maximo, cupo_disponible, venue_id, created_at) VALUES
(1, 'Cosquin Rock 2026',
 'El festival de rock mas importante de Argentina vuelve a Cordoba.',
 'Recitales', '2026-08-15 18:00:00',
 'Estadio Mario Alberto Kempes - Av. Carcano s/n, Cordoba',
 8500.00, 'https://images.unsplash.com/photo-1459749411175-04bf5292ceea?w=800', 300, 300, 1, NOW()),

(2, 'Noche de Tango - Orquesta Tipica',
 'Una velada con los mejores exponentes del tango argentino.',
 'Teatro', '2026-07-20 21:00:00',
 'Teatro Real de Cordoba - Av. Velez Sarsfield 365, Cordoba',
 4500.00, 'https://images.unsplash.com/photo-1504680177321-2e6a879aac86?w=800', 60, 60, 3, NOW()),

(3, 'Boca Juniors vs Talleres',
 'Fecha 18 de la Liga Profesional de Futbol.',
 'Deportes', '2026-09-05 19:30:00',
 'Estadio Mario Alberto Kempes - Av. Carcano s/n, Cordoba',
 6000.00, 'https://images.unsplash.com/photo-1431324155629-1a6deb1dec8d?w=800', 300, 300, 1, NOW()),

(4, 'WOS - Gira Oscuro Extasis',
 'WOS presenta su nuevo album en vivo.',
 'Recitales', '2026-08-28 20:00:00',
 'Orfeo Superdomo - Av. Fuerza Aerea 5200, Cordoba',
 12000.00, 'https://images.unsplash.com/photo-1493225457124-a3eb161ffa5f?w=800', 200, 200, 2, NOW()),

(5, 'Avant Premiere - La Cordillera 2',
 'Proyeccion exclusiva antes del estreno nacional.',
 'Cine', '2026-07-10 20:30:00',
 'Quality Espacio - Av. Cruz Roja Argentina 200, Cordoba',
 3500.00, 'https://images.unsplash.com/photo-1489599849927-2ee91cede3ba?w=800', 40, 40, 4, NOW()),

(6, 'Trueno - Bien o Mal Tour',
 'Trueno llega a Cordoba con todos los hits.',
 'Recitales', '2026-10-12 21:00:00',
 'Plaza de la Musica - Bv. Las Heras 80, Cordoba',
 9500.00, 'https://images.unsplash.com/photo-1470229722913-7c0e2dbbafd3?w=800', 150, 150, 5, NOW()),

(7, 'Hamlet - Compania Nacional de Teatro',
 'La tragedia de Shakespeare.',
 'Teatro', '2026-08-02 20:00:00',
 'Teatro Real de Cordoba - Av. Velez Sarsfield 365, Cordoba',
 5000.00, 'https://images.unsplash.com/photo-1503095396549-807759245b35?w=800', 60, 60, 3, NOW()),

(8, 'Talleres vs River Plate',
 'Fecha 22 de la Liga Profesional.',
 'Deportes', '2026-10-25 17:00:00',
 'Estadio Mario Alberto Kempes - Av. Carcano s/n, Cordoba',
 7000.00, 'https://images.unsplash.com/photo-1522778119026-d647f0596c20?w=800', 300, 300, 1, NOW());

-- GENERAR ASIENTOS POR SECTOR
-- Para cada evento, leer los sectores de su venue y generar asientos
DELIMITER //
CREATE PROCEDURE IF NOT EXISTS generar_asientos_sectores()
BEGIN
    DECLARE done INT DEFAULT FALSE;
    DECLARE v_event_id INT;
    DECLARE v_pn, v_ps, v_te, v_to2, v_pref, v_campo INT;
    DECLARE cur CURSOR FOR
        SELECT e.id, v.cap_platea_norte, v.cap_platea_sur, v.cap_tribuna_este,
               v.cap_tribuna_oeste, v.cap_platea_preferencial, v.cap_campo
        FROM events e JOIN venues v ON e.venue_id = v.id;
    DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;

    OPEN cur;
    read_loop: LOOP
        FETCH cur INTO v_event_id, v_pn, v_ps, v_te, v_to2, v_pref, v_campo;
        IF done THEN LEAVE read_loop; END IF;

        CALL insertar_sector(v_event_id, 'Platea Norte', v_pn);
        CALL insertar_sector(v_event_id, 'Platea Sur', v_ps);
        CALL insertar_sector(v_event_id, 'Tribuna Este', v_te);
        CALL insertar_sector(v_event_id, 'Tribuna Oeste', v_to2);
        CALL insertar_sector(v_event_id, 'Preferencial', v_pref);
        CALL insertar_sector(v_event_id, 'Campo', v_campo);
    END LOOP;
    CLOSE cur;
END //

CREATE PROCEDURE IF NOT EXISTS insertar_sector(IN p_event INT, IN p_sector VARCHAR(50), IN p_cap INT)
proc_body: BEGIN
    DECLARE v_cols INT;
    DECLARE v_filas INT;
    DECLARE f INT DEFAULT 0;
    DECLARE c INT;
    DECLARE v_count INT DEFAULT 0;

    IF p_cap <= 0 THEN
        LEAVE proc_body;
    END IF;

    SET v_cols = LEAST(CEILING(SQRT(p_cap)), 50);
    SET v_filas = CEILING(p_cap / v_cols);

    WHILE f < v_filas AND v_count < p_cap DO
        SET c = 1;
        WHILE c <= v_cols AND v_count < p_cap DO
            INSERT INTO seats (event_id, sector, fila, numero, ocupado)
            VALUES (p_event, p_sector, CHAR(65 + (f MOD 26)), c, 0);
            SET c = c + 1;
            SET v_count = v_count + 1;
        END WHILE;
        SET f = f + 1;
    END WHILE;
END //
DELIMITER ;

CALL generar_asientos_sectores();
DROP PROCEDURE IF EXISTS generar_asientos_sectores;
DROP PROCEDURE IF EXISTS insertar_sector;

-- USUARIOS ADMIN (passwords hasheadas con SHA-256)
DELETE FROM tickets;
DELETE FROM users;
INSERT INTO users (nombre, apellido, email, password, dni, rol, created_at) VALUES
('Simon',  'Factor',       '2408192@ucc.edu.ar', SHA2('Boca07',     256), '47192801', 'admin', NOW()),
('Juan',   'Scaglia',      '2413770@ucc.edu.ar', SHA2('Juanceto01', 256), '46377960', 'admin', NOW()),
('Facundo','Arribillaga',  '2411006@ucc.edu.ar', SHA2('faqarri06',  256), '46450861', 'admin', NOW());

SELECT 'Venues:' AS info, COUNT(*) AS total FROM venues
UNION ALL SELECT 'Events:', COUNT(*) FROM events
UNION ALL SELECT 'Seats:', COUNT(*) FROM seats
UNION ALL SELECT 'Users:', COUNT(*) FROM users;
