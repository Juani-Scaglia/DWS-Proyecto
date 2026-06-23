-- ============================================================
-- SEED DATA - Establecimientos, Eventos y Asientos de prueba
-- Ejecutar: mysql -u root -p --default-character-set=utf8mb4 ticket_system < database/seed.sql
-- ============================================================

SET NAMES utf8mb4;

-- Limpiar datos previos (orden por FKs)
DELETE FROM tickets;
DELETE FROM seats;
DELETE FROM events;
DELETE FROM venues;

-- ESTABLECIMIENTOS (filas y columnas se calculan automaticamente desde capacidad)
INSERT INTO venues (id, nombre, direccion, filas, columnas_por_fila, capacidad, created_at) VALUES
(1, 'Estadio Mario Alberto Kempes',   'Av. Carcano s/n, Cordoba',             13, 13, 150, NOW()),
(2, 'Orfeo Superdomo',                'Av. Fuerza Aerea 5200, Cordoba',       10, 10,  96, NOW()),
(3, 'Teatro Real de Cordoba',         'Av. Velez Sarsfield 365, Cordoba',      8,  8,  60, NOW()),
(4, 'Quality Espacio',                'Av. Cruz Roja Argentina 200, Cordoba',   7,  6,  40, NOW()),
(5, 'Plaza de la Musica',             'Bv. Las Heras 80, Cordoba',             9,  9,  70, NOW());

-- EVENTOS
INSERT INTO events (id, titulo, descripcion, categoria, fecha, lugar, precio, imagen, cupo_maximo, cupo_disponible, venue_id, created_at) VALUES
(1, 'Cosquin Rock 2026',
 'El festival de rock mas importante de Argentina vuelve a Cordoba con una grilla de artistas nacionales e internacionales.',
 'Recitales', '2026-08-15 18:00:00',
 'Estadio Mario Alberto Kempes - Av. Carcano s/n, Cordoba',
 8500.00, 'https://images.unsplash.com/photo-1459749411175-04bf5292ceea?w=800', 150, 150, 1, NOW()),

(2, 'Noche de Tango - Orquesta Tipica',
 'Una velada con los mejores exponentes del tango argentino en el historico Teatro Real.',
 'Teatro', '2026-07-20 21:00:00',
 'Teatro Real de Cordoba - Av. Velez Sarsfield 365, Cordoba',
 4500.00, 'https://images.unsplash.com/photo-1504680177321-2e6a879aac86?w=800', 60, 60, 3, NOW()),

(3, 'Boca Juniors vs Talleres',
 'Fecha 18 de la Liga Profesional de Futbol. Talleres recibe a Boca en el Kempes.',
 'Deportes', '2026-09-05 19:30:00',
 'Estadio Mario Alberto Kempes - Av. Carcano s/n, Cordoba',
 6000.00, 'https://images.unsplash.com/photo-1431324155629-1a6deb1dec8d?w=800', 150, 150, 1, NOW()),

(4, 'WOS - Gira Oscuro Extasis',
 'WOS presenta su nuevo album en vivo en el Orfeo Superdomo de Cordoba.',
 'Recitales', '2026-08-28 20:00:00',
 'Orfeo Superdomo - Av. Fuerza Aerea 5200, Cordoba',
 12000.00, 'https://images.unsplash.com/photo-1493225457124-a3eb161ffa5f?w=800', 96, 96, 2, NOW()),

(5, 'Avant Premiere - La Cordillera 2',
 'Proyeccion exclusiva antes del estreno nacional con presencia del director y elenco.',
 'Cine', '2026-07-10 20:30:00',
 'Quality Espacio - Av. Cruz Roja Argentina 200, Cordoba',
 3500.00, 'https://images.unsplash.com/photo-1489599849927-2ee91cede3ba?w=800', 40, 40, 4, NOW()),

(6, 'Trueno - Bien o Mal Tour',
 'Trueno llega a Cordoba con todos los hits de su ultimo disco.',
 'Recitales', '2026-10-12 21:00:00',
 'Plaza de la Musica - Bv. Las Heras 80, Cordoba',
 9500.00, 'https://images.unsplash.com/photo-1470229722913-7c0e2dbbafd3?w=800', 70, 70, 5, NOW()),

(7, 'Hamlet - Compania Nacional de Teatro',
 'La tragedia de Shakespeare interpretada por la Compania Nacional.',
 'Teatro', '2026-08-02 20:00:00',
 'Teatro Real de Cordoba - Av. Velez Sarsfield 365, Cordoba',
 5000.00, 'https://images.unsplash.com/photo-1503095396549-807759245b35?w=800', 60, 60, 3, NOW()),

(8, 'Talleres vs River Plate',
 'Fecha 22 de la Liga Profesional. El Matador recibe al Millonario.',
 'Deportes', '2026-10-25 17:00:00',
 'Estadio Mario Alberto Kempes - Av. Carcano s/n, Cordoba',
 7000.00, 'https://images.unsplash.com/photo-1522778119026-d647f0596c20?w=800', 150, 150, 1, NOW());

-- GENERAR ASIENTOS
DELIMITER //
CREATE PROCEDURE IF NOT EXISTS generar_asientos()
BEGIN
    DECLARE done INT DEFAULT FALSE;
    DECLARE v_event_id INT;
    DECLARE v_filas INT;
    DECLARE v_cols INT;
    DECLARE f INT;
    DECLARE c INT;
    DECLARE cur CURSOR FOR
        SELECT e.id, v.filas, v.columnas_por_fila
        FROM events e JOIN venues v ON e.venue_id = v.id;
    DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;

    OPEN cur;
    read_loop: LOOP
        FETCH cur INTO v_event_id, v_filas, v_cols;
        IF done THEN LEAVE read_loop; END IF;
        SET f = 0;
        WHILE f < v_filas DO
            SET c = 1;
            WHILE c <= v_cols DO
                INSERT INTO seats (event_id, fila, numero, ocupado)
                VALUES (v_event_id, CHAR(65 + f), c, 0);
                SET c = c + 1;
            END WHILE;
            SET f = f + 1;
        END WHILE;
    END LOOP;
    CLOSE cur;
END //
DELIMITER ;

CALL generar_asientos();
DROP PROCEDURE IF EXISTS generar_asientos;

SELECT 'Venues:' AS info, COUNT(*) AS total FROM venues
UNION ALL SELECT 'Events:', COUNT(*) FROM events
UNION ALL SELECT 'Seats:', COUNT(*) FROM seats;
