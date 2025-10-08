-- Script de Poblado Completo para MediTrack
-- Fecha: 2025-08-19
-- Descripción: Poblado completo de todas las tablas del sistema

-- ============================================
-- POBLADO DE TABLAS BASE
-- ============================================

-- Poblar centros médicos
INSERT INTO medical_center (id, name, address, phone, email) VALUES
(1, 'Centro Médico Principal', 'Av. Principal 123', '+56 2 2345 6789', 'info@centromedico.cl'),
(2, 'Clínica Norte', 'Calle Norte 456', '987654321', 'norte@meditrack.com')
ON CONFLICT (id) DO NOTHING;

-- Poblar pabellones
INSERT INTO pavilion (id, name, medical_center_id) VALUES
(1, 'Pabellón Central 01', 1),
(2, 'Pabellón Central 02', 1),
(3, 'Pabellón Central 03', 1),
(4, 'Pabellón Central 04', 1),
(5, 'Pabellón Central 05', 1),
(6, 'Pabellón Central 06', 1),
(7, 'Pabellón Central 07', 1),
(8, 'Pabellón Central 08', 1),
(9, 'Pabellón Central 09', 1),
(10, 'Pabellón Central 10', 1),
(11, 'Pabellón Central 11', 1),
(12, 'Pabellón Central 12', 1),
(13, 'Pabellón Central 13', 1),
(14, 'Pabellón Central 14', 1),
(15, 'Pabellón Central 15', 1),
(16, 'Pabellón Central 16', 1),
(17, 'Pabellón Central 17', 1),
(18, 'Pabellón Central 18', 1),
(19, 'Pabellón Central 19', 1),
(20, 'Pabellón Central 20', 1),
(21, 'Pabellón Central 21', 1),
(22, 'Pabellón Central 22', 1),
(23, 'Pabellón Central 23', 1),
(24, 'Pabellón Central 24', 1),
(25, 'Pabellón Central 25', 1)
ON CONFLICT (id) DO NOTHING;

-- Poblar bodegas
INSERT INTO store (id, name, type, medical_center_id) VALUES
(1, 'Bodega Principal', 'central', 1),
(2, 'Bodega Secundaria', 'secundaria', 2)
ON CONFLICT (id) DO NOTHING;

-- Poblar códigos de insumos
INSERT INTO supply_code (code, name, code_supplier, critical_stock) VALUES
(1001, 'Guantes', 5001, 5),
(1002, 'Mascarillas', 5002, 10),
(1003, 'Jeringas', 5003, 10),
(1004, 'Agujas', 5004, 10),
(1005, 'Gasas', 5005, 10)
ON CONFLICT (code) DO NOTHING;

-- Poblar lotes
INSERT INTO batch (id, expiration_date, amount, supplier, store_id, qr_code) VALUES
(1, '2026-12-31', 10, 'Proveedor Uno', 1, 'BATCH_1_1'),
(2, '2025-12-31', 5, 'Proveedor Dos', 2, 'BATCH_2_1'),
(3, '2026-06-30', 15, 'Proveedor Tres', 1, 'BATCH_3_1'),
(4, '2025-10-15', 8, 'Proveedor Cuatro', 2, 'BATCH_4_1')
ON CONFLICT (id) DO NOTHING;

-- Poblar insumos médicos - TODOS EN ESTADO DISPONIBLE
INSERT INTO medical_supply (code, batch_id, qr_code, status) VALUES
-- Guantes del lote 1
(1001, 1, 'SUPPLY_1_1', 'disponible'),
(1001, 1, 'SUPPLY_2_1', 'disponible'),
(1001, 1, 'SUPPLY_3_1', 'disponible'),
(1001, 1, 'SUPPLY_4_1', 'disponible'),
(1001, 1, 'SUPPLY_5_1', 'disponible'),
(1001, 1, 'SUPPLY_6_1', 'disponible'),
(1001, 1, 'SUPPLY_7_1', 'disponible'),
(1001, 1, 'SUPPLY_8_1', 'disponible'),
(1001, 1, 'SUPPLY_9_1', 'disponible'),
(1001, 1, 'SUPPLY_10_1', 'disponible'),
-- Mascarillas del lote 2
(1002, 2, 'SUPPLY_11_1', 'disponible'),
(1002, 2, 'SUPPLY_12_1', 'disponible'),
(1002, 2, 'SUPPLY_13_1', 'disponible'),
(1002, 2, 'SUPPLY_14_1', 'disponible'),
(1002, 2, 'SUPPLY_15_1', 'disponible'),
-- Jeringas del lote 3
(1003, 3, 'SUPPLY_16_1', 'disponible'),
(1003, 3, 'SUPPLY_17_1', 'disponible'),
(1003, 3, 'SUPPLY_18_1', 'disponible'),
(1003, 3, 'SUPPLY_19_1', 'disponible'),
(1003, 3, 'SUPPLY_20_1', 'disponible'),
-- Agujas del lote 4
(1004, 4, 'SUPPLY_21_1', 'disponible'),
(1004, 4, 'SUPPLY_22_1', 'disponible'),
(1004, 4, 'SUPPLY_23_1', 'disponible'),
(1004, 4, 'SUPPLY_24_1', 'disponible'),
(1004, 4, 'SUPPLY_25_1', 'disponible')
ON CONFLICT (qr_code) DO NOTHING;

-- Poblar cirugías
INSERT INTO surgery (id, name, duration) VALUES
(1, 'COLGAJOS COMPLEJOS (ABBE,MUSTARDA,CONV', 2.5),
(2, 'HASTA 5N SUPERFICIE CORPORAL RECEPTORA', 2),
(3, 'CIERRE MUCOSO VESTÍBULO ORAL O GINGIV', 2),
(4, 'OSTEOTOMIAS TOTALES DEL MAXILAR O MAND', 5),
(5, 'REDUCCION ABIERTA DE FRACTURAS MAXILOF', 2.5),
(6, 'SECCIÓN Y/O RESECCIÓN FRENILLOS CAVIDA', 2),
(7, 'ABDOMINOPLASTIA ( REPARADORA)', 3.5),
(8, 'ABORTO RETENIDO, VACIAMIENTO DE (INCLU', 1.5),
(9, 'ABSCESO Y/O HEMATOMA DE MAMA, TRAT.QUIR.', 2),
(10, 'ADENOIDECTOMIA (PROC. AUT.)', 1.5),
(11, 'ADENOMA O CÁNCER PROSTÁTICO, RESECCIÓN E', 5),
(12, 'ADENOMA PROSTATICO, TRAT. QUIR. CUALQUIE', 2),
(13, 'AMIGDALECTOMIA C/PT', 2),
(14, 'AMIGDALECTOMIA C/S ADENOIDECTOMIA, UNI O', 2),
(15, 'AMPUTACION', 2),
(16, 'ANASTOMOSIS DEFERENTES O EPIDIDIMO-DEFER', 4),
(17, 'ANEXECTOMIA Y/O VAC. DE ABSCESO TUBO-OVA', 2),
(18, 'ANTROSTOMÍA SENO MAXILAR, CUALQUIER VÍA', 4),
(19, 'ARTRODESIS DE HOMBRO, CADERA,RODILLA, TO', 2),
(20, 'ARTROSCOPIA DIAGNOSTICA C/S BIOPSIA, C/S', 4),
(21, 'AUTO O HETEROTRANSPLANTE', 5),
(22, 'AXILO-SUPRACLAVICULAR', 2.5),
(23, 'BARTOLINOCISTONEOSTOMIA O EXTIRP. DE LA', 1),
(24, 'BIOPSIA QUIR. GANGLIONAR (CUALQUIER REGI', 2.5),
(25, 'CIRCUNCISION (INCLUYE SECCION DE FRENILL', 2),
(26, 'CIRUGÍA ABIERTA O ENDOSCÓPICA DE LESIONE', 2),
(27, 'CIRUGIA BARIATRICA BY PASS GASTRICO LAP', 4.5),
(28, 'CIRUGIA BARIATRICA MANGA GASTRICA LAP', 4),
(29, 'CISTECTOMIA PARCIAL Y/O TRAT. QUIR. DE D', 4),
(30, 'CISTOPLASTIA, PROC. COMPLETO', 5),
(31, 'COAGULACION DE NUCLEOS O VIAS ENCEFALICA', 9),
(32, 'COLANGIOENTEROANASTOMOSIS INTRAHEPATICA', 2),
(33, 'COLANGIOPANCREATOGRAFÍA RETRÓGRADA C/S P', 2),
(34, 'COLECISTECTOMIA C/S COLANGIOGRAFIA OPERA', 3),
(35, 'COLECISTECTOMIA POR VIDEOLAPAROSCOPIA, P', 2.5),
(36, 'COLECTOMIA PARCIAL O HEMICOLECTOMIA', 4),
(37, 'COLGAJO SIMPLE UNICO', 2.5),
(38, 'COLGAJOS SIMPLES 2 O MAS', 2),
(39, 'CONDILOMAS ANALES, TRAT. QUIR.(PARA ELEC', 2),
(40, 'CONIZACION Y/O AMPUTACION DEL CUELLO, DI', 1.5),
(41, 'CONTRACTURA DUPUYTREN,TRAT. QUIR.,CADA T', 2),
(42, 'CRANEOPLASTIA CON PROTESIS (NO INCLUYE E', 3),
(43, 'CRANIECTOMIAS C/S REMODELACION OSEA', 7),
(44, 'CURETAJE POR GRUPO', 2),
(45, 'DE COMPLEJIDAD MAYOR: INCLUYE REEMPLAZ', 7),
(46, 'DECORTICACION DE CUERDAS VOCALES C/MICRO', 2.5),
(47, 'DEDOS EN GATILLO,TRAT. QUIR.,CUALQUIER N', 2),
(48, 'DESCENSO TESTICULAR CON O SIN HERNIA, CU', 2),
(49, 'DISFUNCION PATELO-FEMORAL,REALINEAMIENTO', 3),
(50, 'DRENAJE PERCUTÁNEO O ENDOSCÓPICO DE HIDR', 2),
(51, 'ECTROPION, PLASTIA DE', 1),
(52, 'EMBARAZO TUBARIO, TRAT. QUIR.', 3),
(53, 'ENCEFALICOS Y DE HIPOFISIS', 8),
(54, 'ENDOPROTESIS TOTAL DE CADERA', 3),
(55, 'ENDOPROTESIS TOTAL DE HOMBRO,(CUALQUIER', 3),
(56, 'ENDOPROTESIS TOTAL DE RODILLA, (CUALQUIE', 3),
(57, 'ENTERO-ENTEROANASTOMOSIS O ENTEROCOLOANA', 3.5),
(58, 'EPIFISIODESIS (FEMUR Y/O TIBIA)', 2),
(59, 'ESCARECTOMIA  HASTA 10N SUPERFICIE CORPO', 2),
(60, 'ESCARECTOMIA HASTA 5 N SUPERFICIE CORPOR', 2),
(61, 'ESCOLIOSIS,TRAT.QUIR.,CUALQUIER VIA DE A', 6),
(62, 'ESOFAGECTOMIA TOTAL CON ESOFAGOSTOMIA, G', 4),
(63, 'ESTENOSIS LARINGOTRAQUEALES Y/O FARÍNGEA', 2.5),
(64, 'ESTRABISMO TRAT. QX. COMPLETO UNI O BIL', 2),
(65, 'EXODONCIA DE PIEZAS INCLUIDAS', 2),
(66, 'EXTIRPACIÓN DE GLÁNDULA SALIVAL SUBMANDI', 2.5),
(67, 'EXTIRPACIÓN TOTAL O PARCIAL DE LA GLÁNDU', 2),
(68, 'FACOERESIS EXTRACAPSULAR CON IMPLANTE D', 1),
(69, 'FENESTRACION, SEPTOSTOMIA O COAGULACION', 5),
(70, 'FIJACION DE COLUMNA (CERVICAL-DORSAL-LUM', 4.5),
(71, 'FISTULA ANORRECTAL, TRAT.QUIR.', 2),
(72, 'FISTULA ARTERIOVENOSA (DE BRESCIA O SIMI', 2),
(73, 'FISTULA ARTERIOVENOSA CONGENITA O TRAUMA', 2),
(74, 'FRACTURA CUELLO HUMERAL, TRAT. QUIR.', 3),
(75, 'FRACTURA DE CLAVICULA, OSTEOSINTESIS', 3),
(76, 'FRACTURA O PSEUDOARTROSIS ESCAFOIDES,TRA', 3),
(77, 'FRACTURA ROTULA: OSTEOSINTESIS O PATELEC', 3),
(78, 'FRACTURAS CONDILEAS O DE PLATILLOS TIBIA', 3.5),
(79, 'GASTRECTOMIA TOTAL O SUBTOTAL AMPLIADA (', 4),
(80, 'GASTRODUODENOSCOPIA (INCLUYE ESOFAGOSCOP', 2),
(81, 'GASTROENTEROANASTOMOSIS, CUALQUIER TECNI', 4),
(82, 'GASTROSTOMIA PERCUTANEA TECNICA SELDINGE', 2),
(83, 'GASTROTOMIA Y/O GASTROSTOMIA (PROC. AUT.', 3),
(84, 'GINECOMASTIA, CORRECCION PLASTICA', 2),
(85, 'HALLUX VALGUS O RIGIDUS,TRAT.QUIR. COMPL', 2),
(86, 'HEMATOMA, EMPIEMA O COLECCION SUBDURAL,', 4),
(87, 'HEMORROIDECTOMIA (INCLUYE OTRAS OPERA-', 2),
(88, 'HERNIA ABDOMINAL POR LAPAROTOMÍA (NO INC', 2.5),
(89, 'HERNIA INGUINAL, CRURAL, UMBILICAL, DE L', 2),
(90, 'HERNIA NUCLEO PULPOSO, ESTENORRAQUIS, AR', 3),
(91, 'HIDROCELE Y/O HEMATOCELE, INCLUYE QUISTE', 2),
(92, 'HISTERECTOMIA POR VIA VAGINAL', 2.5),
(93, 'HISTERECTOMIA RADICAL CON DISECCION PELV', 3.5),
(94, 'HISTERECTOMIA TOTAL C/INTERVENCION INCON', 3),
(95, 'HISTERECTOMIA TOTAL O AMPLIADA POR VIA A', 2.5),
(96, 'HISTEROSCOPÍA DIAGNÓSTICA (PROC. AUT.)', 2),
(97, 'IMPLANTE CATHETER RESERVOREO PARA QMT', 1.5),
(98, 'IMPLANTE FILTROS VENOSOS', 2),
(99, 'INCOMPETENCIA CERVICAL TRAT. QUIR.', 2),
(100, 'INCONTINENCIA URINARIA DE ESFUERZO O D.I', 2),
(101, 'INESTABILIDAD CRONICA DE RODILLA, RECONS', 3),
(102, 'INFILTRACION FACETARIA COLUMNA Y RADICUL', 1.5),
(103, 'INJERTOS HASTA 10% SUPERFICIE CORPORAL R', 2),
(104, 'INSTALACION CATETER TUNELIZADO DIALISIS', 1.5),
(105, 'INSTALACIÓN DE CATÉTER CON RESERVORIO SU', 1.5),
(106, 'LAMINECTOMIA DESCOMPRESIVA', 2.5),
(107, 'LAPAROTOMIA EXPLORADORA, C/S LIBERACION', 3),
(108, 'LIBERACIÓN QUIRÚRGICA DE NERVIO PERIFÉRI', 2),
(109, 'LIGAMENTO ANCHO: ABSCESOS Y/O HEMATOMAS',2.5),
(110, 'LOBECTOMIA HEPATICA (PROC. AUT.)', 4.5),
(111, 'LOBECTOMIA O BILOBECTOMIA', 4.5),
(112, 'LUXACION ACROMIO-CLAVICULAR O ESTERNO CL', 3),
(113, 'LUXACION RECIDIVANTE, TRAT. QUIR.', 2.5),
(114, 'LUXOFRACTURA TOBILLO, CUALQUIER TIPO, OS', 2.5),
(115, 'MAMOPLASTIA DE REDUCCION', 3.5),
(116, 'MAMOPLASTIA ESTETICA BILATERIAL DE AUMEN', 3),
(117, 'MASTECTOMIA PARCIAL (CUADRANTECTOMIA O S', 2),
(118, 'MASTECTOMIA RADICAL O TUMORECTOMIA C/VAC', 3),
(119, 'MEDIASTINICOS', 2),
(120, 'MENISCECTOMIA U OTRAS INTERVENCIONES POR', 2),
(121, 'MIOMECTOMIA', 2),
(122, 'NEFRECTOMÍA PARCIAL CUALQUIER VÍA Y TÉCN', 5),
(123, 'NEFRECTOMIA RADICAL POR CANCER RENAL, TR', 4.5),
(124, 'NEUROLISIS CON TECNICA MICROQUIRURGICA', 3),
(125, 'NEUROLISIS EXTERNA', 2.5),
(126, 'OOFORECTOMIA PARCIAL O TOTAL, UNI O BILA', 2.5),
(127, 'OPERACION DE SALVATAJE CADERA, COLUMNA O', 3.5),
(128, 'ORQUIDECTOMÍA AMPLIADA POR CÁNCER TESTIC', 2),
(129, 'ORTEJOS EN GARRA,TRAT.QUIR.,CUALQ.NUMERO', 2),
(130, 'OSTEOMIELITIS AGUDA HEMATOGENA, DRENAJE', 3),
(131, 'OSTEOMIELITIS CRONICA HUESOS LARGOS, LEG', 2.5),
(132, 'OSTEOMIELITIS, LIMPIEZA QUIRURGICA', 2),
(133, 'OSTEOSINTESIS DIAFISIARIA (CUALQUIER TEC', 3),
(134, 'OSTEOSINTESIS METACARPIANAS O DE FALANGE', 2),
(135, 'OSTEOSINTESIS RADIO, (CUALQUIER TECNICA)', 3),
(136, 'OSTEOSINTESIS SUPRA O INTERCONDILEA (CUA', 2.5),
(137, 'OSTEOSINTESIS TIBIO-PERONE  (CUALQUIER T', 3.5),
(138, 'OSTEOSINTESIS, FRACT.CERRADA CUBITO Y/O', 3),
(139, 'OSTEOTOMIA CORRECTORA', 2),
(140, 'OTRAS DERIVACIONES: FEMORO-FEMORAL, AXIL', 3.5),
(141, 'PABELLON DE ESTADIA MINIMA', 2),
(142, 'PANCREATECTOMIA TOTAL C/S ESPLENECTOMIA', 7),
(143, 'PARATIROIDECTOMÍA O EXPLORACIÓN DE PARAT', 2.5),
(144, 'PARATIROIDECTOMÍA O EXPLORACIÓN PARATIRO', 3.5),
(145, 'PAROTIDECTOMÍA SUPRAFACIAL CON DISECCIÓN', 3),
(146, 'PAROTIDECTOMÍA TOTAL, CON DISECCIÓN Y PR', 3),
(147, 'PIE BOT U OTRAS MALFORMACIONES CONGENITA', 2),
(148, 'PIE PLANO, TRAT. QUIR. (CUALQUIER TECNIC', 2),
(149, 'PIELOTOMIA EXPLORADORA Y/O TERAPEUTICA (', 4.5),
(150, 'PLASTIAS EN Z, HASTA 3', 2),
(151, 'PLEURODESIS POR TORACOTOMIA', 2),
(152, 'POLIDACTILIA, EXTIRPACION Y PLASTIA UN L', 2),
(153, 'PROLAPSO ANTERIOR Y/O POSTERIOR C/S TRAT', 2),
(154, 'PROLAPSO ANTERIOR Y/O POSTERIOR CON REPA', 2),
(155, 'PTERIGION Y/O PSEUDOPTERIGION O SU RECID', 1),
(156, 'PTOSIS, TRAT. QUIR.', 2),
(157, 'PUNCION LUMBAR PARA MIELOGRAMAS', 1.5),
(158, 'QUERATECTOMIA FOTORREFRACTIVA O FOTOTERA', 1),
(159, 'QUISTE SACROCOXIGEO, TRAT. QUIR.', 2),
(160, 'QUISTES SINOVIALES DE VAINAS FLEXORAS, B', 2),
(161, 'RASPADO UTERINO DIAGNOSTICO O TERAPEUT', 1),
(162, 'RECONSTITUCION  DE TRANSITO EN 2 TIEMPO', 3),
(163, 'RECONSTITUCION  TRANSITO POST OPERACION', 3.5),
(164, 'RECONSTRUCCION AREOLA Y/O PEZON C/S PLAS', 2.5),
(165, 'RECONSTRUCCION MAMARIA', 4),
(166, 'REPARACION DE FISTULA DE LCR', 3),
(167, 'REPARACION QUIR. DE VASOS ARTERIALES Y/O', 3.5),
(168, 'RESECCION DE COSTILLAS Y/O PARED COSTAL', 3),
(169, 'RESECCION DE PARED COSTAL C/PLASTIA (TOR', 2.5),
(170, 'RESECCION ENDOSCOPICA DE CANCER VESICAL', 2.5),
(171, 'RESECCIÓN PLÁSTICA DE HASTA 2 CICATRICES', 2),
(172, 'RESECCIONES SEGMENTARIAS ANATÓMI DE PULM', 3),
(173, 'RETIRO CATHETER TUNELIZADO DIALISIS', 2),
(174, 'RETIRO DE CATÉTER CON RESERVORIO SUBCUTÁ', 2),
(175, 'RETIRO DE ENDOPROTESIS U OSTEOSINTESIS I', 2),
(176, 'RETIRO DE PLACAS RECTAS O ANGULADAS', 2),
(177, 'RETIRO DE TORNILLOS, CLAVOS, AGUJAS DE O', 2),
(178, 'RINOPLASTIA Y/O SEPTOPLASTIA, CUALQUIER', 3.5),
(179, 'RIZOTOMIA (CUALQUIER TECNICA)', 2),
(180, 'RUPTURA MANGUITO ROTADORES, TRAT. QUIR.', 3.5),
(181, 'RUPTURA TENDON DE AQUILES O TIBIAL POSTE', 2),
(182, 'SAFENECTOMÍA INTERNA Y/O EXTERNA,', 2),
(183, 'SAFENECTOMIA POR RADIOFRECUENCIA', 2),
(184, 'SALPINGECTOMIA UNI O BILATERAL', 2.5),
(185, 'SECCION DE  NERVIO, REPARACION CON INJER', 3),
(186, 'SENO ESFENOIDAL, ABERTURA POR CUALQUIER', 3),
(187, 'SIMPATECTOMIA CERVICO-TORACICA', 2),
(188, 'SUPRARRENALECTOMIA UNILATERAL', 4.5),
(189, 'TENORRAFIA EXTENSORES', 2),
(190, 'TENORRAFIA EXTENSORES O TENOTOMIA DE ALA', 3),
(191, 'TENORRAFIA O INJERTOS FLEXORES', 2),
(192, 'TIMPANOPLASTIA FUNCIONAL (CUALQUIER TIPO', 3.5),
(193, 'TIROIDECTOMIA BILATERAL TOTAL', 3.5),
(194, 'TIROIDECTOMÍA TOTAL AMPLIADA (INCLUYE EX', 3),
(195, 'TORACOTOMIA EXPLORADORA, C/S BIOPSIA, C/', 2),
(196, 'TRAQUEOSTOMIA (PROC. AUT.)', 2),
(197, 'TRATAMIENTO QUIRÚRGICO DE MUCOSITIS TIMP', 2),
(198, 'TRATAMIENTO QUIRÚRGICO PÓLIPO NASAL', 3),
(199, 'TUMOR BENIGNO DE CUERDAS VOCALES, TRAT.', 2),
(200, 'TUMOR OSEO, RESECCION EN BLOQUE, C/S OST', 2),
(201, 'TUMOR Y/O QUISTE RETROPERITONEAL', 3),
(202, 'TUMORES MALIGNOS DE PROSTATA O VESICULAS', 6),
(203, 'TUMORES O QUISTES DE MEDIASTINO (ANTERIO', 3),
(204, 'TUMORES O QUISTES O LESIONES PSEUDOQUIST', 2),
(205, 'URETERO-LITOTOMIA ENDOSCOPICA C/URETEROS', 2.5),
(206, 'URETRECTOMÍA Y /O PLASTÍA ABIERTA DE UR', 4),
(207, 'URETROTOMIA INTERNA Y/O URETROLITOTOMIA', 2.5),
(208, 'VACIAMIENTO (DISECCIÓN) RADICAL CUELLO (', 3),
(209, 'VACIAMIENTO Y CURETAJE QUIRURGICO DE LES', 2.5),
(210, 'VARICOCELE UNILATERAL Y/O DENERVACIÓN CO', 2),
(211, 'VASECTOMIA BILATERAL, (PROC. AUT.) (LA V', 1),
(212, 'VENTRICULOCISTERNOSTOMIA', 4.5),
(213, 'VENTRICULOSTOMIA O INSTALACION DE DERIVA', 2.5),
(214, 'VIDEOLAPAROSCOPIA GINECOLOGICA EXPLORADO', 2.5),
(215, 'VITRECTOMIA C/RETINOTOMIA (C/S INYECCION', 3),
(216, 'VULVECTOMIA SIMPLE', 1),
(217, 'YUGULAR SIMPLE', 3),
(218, 'ANEURISMA AORTICO ABDOMINAL TRAT. QUIR.', 4);


-- ============================================
-- POBLADO DE USUARIOS
-- ============================================

-- Script para insertar usuarios del sistema
-- Contraseña: admin123 (hasheada con bcrypt)
-- Nota: Las contraseñas se almacenan como hash, no como texto plano

-- Usuario administrador del sistema
INSERT INTO "user" (
    rut, 
    name, 
    email, 
    password, 
    role, 
    medical_center_id, 
    is_active, 
    created_at, 
    updated_at
) VALUES (
    '12345678-9',
    'Administrador del Sistema',
    'admin@meditrack.com',
    '$2a$10$NA3QLOvkwhpcs.X4KxjONObslo1LreYA6qAzdQcqxRrD4ktjBrpmO', -- admin123 hasheada con bcrypt
    'admin',
    1,
    true,
    EXTRACT(EPOCH FROM NOW()),
    EXTRACT(EPOCH FROM NOW())
) ON CONFLICT (rut) DO NOTHING;

-- Usuario de pabellón de ejemplo
INSERT INTO "user" (
    rut, 
    name, 
    email, 
    password, 
    role, 
    medical_center_id, 
    is_active, 
    created_at, 
    updated_at
) VALUES (
    '87654321-0',
    'Usuario Pabellón',
    'pabellon@meditrack.com',
    '$2a$10$NA3QLOvkwhpcs.X4KxjONObslo1LreYA6qAzdQcqxRrD4ktjBrpmO', -- admin123 hasheada con bcrypt
    'pabellón',
    1,
    true,
    EXTRACT(EPOCH FROM NOW()),
    EXTRACT(EPOCH FROM NOW())
) ON CONFLICT (rut) DO NOTHING;

-- Usuario encargado de bodega de ejemplo
INSERT INTO "user" (
    rut, 
    name, 
    email, 
    password, 
    role, 
    medical_center_id, 
    is_active, 
    created_at, 
    updated_at
) VALUES (
    '11111111-1',
    'Encargado Bodega',
    'bodega@meditrack.com',
    '$2a$10$NA3QLOvkwhpcs.X4KxjONObslo1LreYA6qAzdQcqxRrD4ktjBrpmO', -- admin123 hasheada con bcrypt
    'encargado de bodega',
    1,
    true,
    EXTRACT(EPOCH FROM NOW()),
    EXTRACT(EPOCH FROM NOW())
) ON CONFLICT (rut) DO NOTHING;

-- Usuario enfermera de ejemplo
INSERT INTO "user" (
    rut, 
    name, 
    email, 
    password, 
    role, 
    medical_center_id, 
    is_active, 
    created_at, 
    updated_at
) VALUES (
    '22222222-2',
    'María González',
    'enfermera@meditrack.com',
    '$2a$10$NA3QLOvkwhpcs.X4KxjONObslo1LreYA6qAzdQcqxRrD4ktjBrpmO', -- admin123 hasheada con bcrypt
    'enfermera',
    1,
    true,
    EXTRACT(EPOCH FROM NOW()),
    EXTRACT(EPOCH FROM NOW())
) ON CONFLICT (rut) DO NOTHING;

-- Usuario doctor de ejemplo
INSERT INTO "user" (
    rut, 
    name, 
    email, 
    password, 
    role, 
    medical_center_id, 
    is_active, 
    created_at, 
    updated_at
) VALUES (
    '33333333-3',
    'Dr. Carlos Pérez',
    'doctor@meditrack.com',
    '$2a$10$NA3QLOvkwhpcs.X4KxjONObslo1LreYA6qAzdQcqxRrD4ktjBrpmO', -- admin123 hasheada con bcrypt
    'doctor',
    1,
    true,
    EXTRACT(EPOCH FROM NOW()),
    EXTRACT(EPOCH FROM NOW())
) ON CONFLICT (rut) DO NOTHING;

-- Usuario especial del sistema para operaciones automáticas
INSERT INTO "user" (
    rut, 
    name, 
    email, 
    password, 
    role, 
    medical_center_id, 
    is_active, 
    created_at, 
    updated_at
) VALUES (
    'SYSTEM-INIT',
    'Sistema de Inicialización',
    'system@meditrack.com',
    '$2a$10$NA3QLOvkwhpcs.X4KxjONObslo1LreYA6qAzdQcqxRrD4ktjBrpmO', -- admin123 hasheada con bcrypt
    'admin',
    1,
    true,
    EXTRACT(EPOCH FROM NOW()),
    EXTRACT(EPOCH FROM NOW())
) ON CONFLICT (rut) DO NOTHING;

-- ============================================
-- POBLADO DE HISTORIALES
-- ============================================

-- NOTA: No se incluye historial de insumos para mantener todos los insumos en estado inicial "disponible"

-- Poblar historial de lotes
INSERT INTO batch_history (date_time, change_details, previous_values, new_values, user_name, batch_id, user_rut, batch_number) VALUES
('2025-08-16 10:00:00', 'Lote creado', NULL, '{"expiration_date": "2026-12-31", "amount": 10, "supplier": "Proveedor Uno", "store_id": 1}', 'Administrador del Sistema', 1, '12345678-9', 1),
('2025-08-16 11:00:00', 'Lote creado', NULL, '{"expiration_date": "2025-12-31", "amount": 5, "supplier": "Proveedor Dos", "store_id": 2}', 'Usuario Pabellón', 2, '87654321-0', 2),
('2025-08-16 12:00:00', 'Lote creado', NULL, '{"expiration_date": "2026-06-30", "amount": 15, "supplier": "Proveedor Tres", "store_id": 1}', 'Encargado Bodega', 3, '11111111-1', 3),
('2025-08-16 13:00:00', 'Lote creado', NULL, '{"expiration_date": "2025-10-15", "amount": 8, "supplier": "Proveedor Cuatro", "store_id": 2}', 'Administrador del Sistema', 4, '12345678-9', 4),
('2025-08-16 14:00:00', 'Cantidad actualizada', '{"amount": 10}', '{"amount": 8}', 'Encargado Bodega', 1, '11111111-1', 1)
ON CONFLICT DO NOTHING;

-- Insertar solicitud de ejemplo
INSERT INTO supply_request (
    request_number, pavilion_id, requested_by, requested_by_name,
    request_date, status, priority, notes, medical_center_id
) VALUES (
    'SOL-20250120140000', 1, '12345678-9', 'Juan Pérez',
    NOW() - INTERVAL '1 hour', 'pending', 'normal',
    'Solicitud de prueba para implementación de trazabilidad QR', 1
);

-- Insertar items de ejemplo
INSERT INTO supply_request_item (
    supply_request_id, supply_code, supply_name, quantity_requested,
    specifications, is_pediatric, size, urgency_level
) VALUES 
(1, 1001, 'Guantes', 50, 'Talla M, látex libre', FALSE, 'M', 'normal'),
(1, 1002, 'Mascarillas', 100, 'N95, uso pediátrico', TRUE, 'Pediatric', 'high');

-- Script para arreglar el problema de secuencia en la tabla batch
-- Este script resetea la secuencia de auto-incremento al valor correcto

-- Verificar el ID máximo actual en la tabla batch
SELECT 'ID máximo actual en batch:' as info, COALESCE(MAX(id), 0) as max_id FROM batch;

-- Resetear la secuencia al valor correcto
-- La secuencia debe ser mayor que el ID máximo existente
SELECT setval('batch_id_seq', COALESCE((SELECT MAX(id) FROM batch), 0) + 1, false);

-- Verificar el estado de la secuencia después del reset
SELECT 'Secuencia después del reset:' as info, last_value, is_called FROM batch_id_seq;

-- También verificar que no hay problemas con QR codes duplicados
SELECT 'QR codes duplicados:' as info, qr_code, COUNT(*) as count 
FROM batch 
WHERE qr_code IS NOT NULL 
GROUP BY qr_code 
HAVING COUNT(*) > 1;

-- Migración para poblar correctamente la tabla supply_history
-- Fecha: 2025-01-20
-- Descripción: Agregar registros de historial para mostrar que los insumos estuvieron en bodega

-- Insertar historial inicial para todos los insumos disponibles
-- Esto simula que todos los insumos estuvieron inicialmente en bodega antes de estar disponibles

INSERT INTO supply_history (
    date_time,
    status,
    destination_type,
    destination_id,
    medical_supply_id,
    user_rut,
    notes
)
SELECT 
    NOW() - INTERVAL '30 days' AS date_time,  -- Hace 30 días
    'disponible' AS status,
    'store' AS destination_type,
    b.store_id AS destination_id,
    ms.id AS medical_supply_id,
    'SYSTEM-INIT' AS user_rut,
    'Registro inicial - insumo ingresado a bodega desde proveedor' AS notes
FROM medical_supply ms
JOIN batch b ON ms.batch_id = b.id
WHERE ms.status = 'disponible'
AND NOT EXISTS (
    SELECT 1 FROM supply_history sh 
    WHERE sh.medical_supply_id = ms.id
);

-- Verificar los registros insertados
SELECT 
    'Registros de historial agregados:' as info,
    COUNT(*) as total_records
FROM supply_history 
WHERE user_rut = 'SYSTEM-INIT';

-- Mostrar ejemplo de registros insertados
SELECT 
    'Ejemplo de registros insertados:' as info,
    ms.qr_code,
    sc.name as supply_name,
    s.name as store_name,
    sh.date_time,
    sh.status,
    sh.notes
FROM supply_history sh
JOIN medical_supply ms ON sh.medical_supply_id = ms.id
JOIN supply_code sc ON ms.code = sc.code
JOIN store s ON sh.destination_id = s.id
WHERE sh.user_rut = 'SYSTEM-INIT'
LIMIT 5;