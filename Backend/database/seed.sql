-- ============================================================
-- SEED DATA — Establecimientos, Eventos y Asientos de prueba
-- Ejecutar: mysql -u root -p ticket_system < database/seed.sql
-- ============================================================

-- Limpiar datos previos (orden por FKs)
DELETE FROM tickets;
DELETE FROM seats;
DELETE FROM events;
DELETE FROM venues;

-- ──────────────────────────────────────────────
-- ESTABLECIMIENTOS (venues)
-- ──────────────────────────────────────────────

INSERT INTO venues (id, nombre, direccion, filas, columnas_por_fila, capacidad, created_at) VALUES
(1, 'Estadio Mario Alberto Kempes',   'Av. Cárcano s/n, Córdoba',             10, 15, 150, NOW()),
(2, 'Orfeo Superdomo',                'Av. Fuerza Aérea 5200, Córdoba',        8, 12,  96, NOW()),
(3, 'Teatro Real de Córdoba',         'Av. Vélez Sarsfield 365, Córdoba',      6, 10,  60, NOW()),
(4, 'Quality Espacio',                'Av. Cruz Roja Argentina 200, Córdoba',   5,  8,  40, NOW()),
(5, 'Plaza de la Música',             'Bv. Las Heras 80, Córdoba',             7, 10,  70, NOW());

-- ──────────────────────────────────────────────
-- EVENTOS
-- ──────────────────────────────────────────────

INSERT INTO events (id, titulo, descripcion, categoria, fecha, lugar, precio, cupo_maximo, cupo_disponible, venue_id, created_at) VALUES
(1,
 'Cosquín Rock 2026',
 'El festival de rock más importante de Argentina vuelve a Córdoba con una grilla de artistas nacionales e internacionales.',
 'Recitales',
 '2026-08-15 18:00:00',
 'Estadio Mario Alberto Kempes - Av. Cárcano s/n, Córdoba',
 8500.00, 150, 150, 1, NOW()),

(2,
 'Noche de Tango — Orquesta Típica',
 'Una velada con los mejores exponentes del tango argentino en el histórico Teatro Real.',
 'Teatro',
 '2026-07-20 21:00:00',
 'Teatro Real de Córdoba - Av. Vélez Sarsfield 365, Córdoba',
 4500.00, 60, 60, 3, NOW()),

(3,
 'Boca Juniors vs Talleres',
 'Fecha 18 de la Liga Profesional de Fútbol. Talleres recibe a Boca en el Kempes.',
 'Deportes',
 '2026-09-05 19:30:00',
 'Estadio Mario Alberto Kempes - Av. Cárcano s/n, Córdoba',
 6000.00, 150, 150, 1, NOW()),

(4,
 'WOS — Gira Oscuro Éxtasis',
 'WOS presenta su nuevo álbum en vivo en el Orfeo Superdomo de Córdoba.',
 'Recitales',
 '2026-08-28 20:00:00',
 'Orfeo Superdomo - Av. Fuerza Aérea 5200, Córdoba',
 12000.00, 96, 96, 2, NOW()),

(5,
 'Avant Première — La Cordillera 2',
 'Proyección exclusiva antes del estreno nacional con presencia del director y elenco.',
 'Cine',
 '2026-07-10 20:30:00',
 'Quality Espacio - Av. Cruz Roja Argentina 200, Córdoba',
 3500.00, 40, 40, 4, NOW()),

(6,
 'Trueno — Bien o Mal Tour',
 'Trueno llega a Córdoba con todos los hits de su último disco.',
 'Recitales',
 '2026-10-12 21:00:00',
 'Plaza de la Música - Bv. Las Heras 80, Córdoba',
 9500.00, 70, 70, 5, NOW()),

(7,
 'Hamlet — Compañía Nacional de Teatro',
 'La tragedia de Shakespeare interpretada por la Compañía Nacional.',
 'Teatro',
 '2026-08-02 20:00:00',
 'Teatro Real de Córdoba - Av. Vélez Sarsfield 365, Córdoba',
 5000.00, 60, 60, 3, NOW()),

(8,
 'Talleres vs River Plate',
 'Fecha 22 de la Liga Profesional. El Matador recibe al Millonario.',
 'Deportes',
 '2026-10-25 17:00:00',
 'Estadio Mario Alberto Kempes - Av. Cárcano s/n, Córdoba',
 7000.00, 150, 150, 1, NOW());

-- ──────────────────────────────────────────────
-- ASIENTOS (seats) — generados por evento
-- Filas: A, B, C... según venue.filas
-- Números: 1..venue.columnas_por_fila
-- ──────────────────────────────────────────────

-- Procedimiento para generar asientos automáticamente
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
        FROM events e
        JOIN venues v ON e.venue_id = v.id;

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

-- Verificación
SELECT 'Venues:' AS info, COUNT(*) AS total FROM venues
UNION ALL
SELECT 'Events:', COUNT(*) FROM events
UNION ALL
SELECT 'Seats:', COUNT(*) FROM seats;
