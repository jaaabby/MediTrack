-- Script de Poblado Completo para MediTrack
-- Fecha: 2025-08-19
-- Descripción: Poblado completo de todas las tablas del sistema

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

-- Poblar especialidades médicas (si no existen ya por la migración)
INSERT INTO medical_specialty (name, description, code, is_active) VALUES
('Traumatología', 'Especialidad médica dedicada al diagnóstico y tratamiento de lesiones del sistema musculoesquelético', 'TRAUMA', TRUE),
('Cardiología', 'Especialidad médica dedicada al diagnóstico y tratamiento de enfermedades del corazón y del sistema circulatorio', 'CARDIOL', TRUE),
('Cirugía General', 'Especialidad médica que abarca procedimientos quirúrgicos de diversas partes del cuerpo', 'CIR_GEN', TRUE),
('Neurocirugía', 'Especialidad médica dedicada al diagnóstico y tratamiento quirúrgico de enfermedades del sistema nervioso', 'NEURO', TRUE),
('Oftalmología', 'Especialidad médica dedicada al diagnóstico y tratamiento de enfermedades de los ojos', 'OFTAL', TRUE),
('Otorrinolaringología', 'Especialidad médica dedicada al diagnóstico y tratamiento de enfermedades del oído, nariz y garganta', 'ORL', TRUE),
('Urología', 'Especialidad médica dedicada al diagnóstico y tratamiento de enfermedades del sistema urinario y genital masculino', 'URO', TRUE),
('Ginecología', 'Especialidad médica dedicada al diagnóstico y tratamiento de enfermedades del sistema reproductor femenino', 'GINEC', TRUE),
('Anestesiología', 'Especialidad médica dedicada al alivio del dolor y el cuidado del paciente durante procedimientos quirúrgicos', 'ANEST', TRUE),
('Plástica y Reconstructiva', 'Especialidad médica dedicada a la reconstrucción y mejoramiento estético de partes del cuerpo', 'PLAST', TRUE)
ON CONFLICT (name) DO NOTHING;

-- Poblar cirugías (con especialidades asignadas)
INSERT INTO surgery (id, name, duration, specialty_id) VALUES
(1, 'COLGAJOS COMPLEJOS (ABBE,MUSTARDA,CONV', 2.5, (SELECT id FROM medical_specialty WHERE code = 'PLAST')),
(2, 'HASTA 5N SUPERFICIE CORPORAL RECEPTORA', 2, (SELECT id FROM medical_specialty WHERE code = 'PLAST')),
(3, 'CIERRE MUCOSO VESTÍBULO ORAL O GINGIV', 2, (SELECT id FROM medical_specialty WHERE code = 'CIR_GEN')),
(4, 'OSTEOTOMIAS TOTALES DEL MAXILAR O MAND', 5, (SELECT id FROM medical_specialty WHERE code = 'CIR_GEN')),
(5, 'REDUCCION ABIERTA DE FRACTURAS MAXILOF', 2.5, (SELECT id FROM medical_specialty WHERE code = 'TRAUMA')),
(6, 'SECCIÓN Y/O RESECCIÓN FRENILLOS CAVIDA', 2, NULL),
(7, 'ABDOMINOPLASTIA ( REPARADORA)', 3.5, NULL),
(8, 'ABORTO RETENIDO, VACIAMIENTO DE (INCLU', 1.5, NULL),
(9, 'ABSCESO Y/O HEMATOMA DE MAMA, TRAT.QUIR.', 2, NULL),
(10, 'ADENOIDECTOMIA (PROC. AUT.)', 1.5, NULL),
(11, 'ADENOMA O CÁNCER PROSTÁTICO, RESECCIÓN E', 5, NULL),
(12, 'ADENOMA PROSTATICO, TRAT. QUIR. CUALQUIE', 2, NULL),
(13, 'AMIGDALECTOMIA C/PT', 2, NULL),
(14, 'AMIGDALECTOMIA C/S ADENOIDECTOMIA, UNI O', 2, NULL),
(15, 'AMPUTACION', 2, NULL),
(16, 'ANASTOMOSIS DEFERENTES O EPIDIDIMO-DEFER', 4, NULL),
(17, 'ANEXECTOMIA Y/O VAC. DE ABSCESO TUBO-OVA', 2, NULL),
(18, 'ANTROSTOMÍA SENO MAXILAR, CUALQUIER VÍA', 4, NULL),
(19, 'ARTRODESIS DE HOMBRO, CADERA,RODILLA, TO', 2, NULL),
(20, 'ARTROSCOPIA DIAGNOSTICA C/S BIOPSIA, C/S', 4, NULL),
(21, 'AUTO O HETEROTRANSPLANTE', 5, NULL),
(22, 'AXILO-SUPRACLAVICULAR', 2.5, NULL),
(23, 'BARTOLINOCISTONEOSTOMIA O EXTIRP. DE LA', 1, NULL),
(24, 'BIOPSIA QUIR. GANGLIONAR (CUALQUIER REGI', 2.5, NULL),
(25, 'CIRCUNCISION (INCLUYE SECCION DE FRENILL', 2, NULL),
(26, 'CIRUGÍA ABIERTA O ENDOSCÓPICA DE LESIONE', 2, NULL),
(27, 'CIRUGIA BARIATRICA BY PASS GASTRICO LAP', 4.5, NULL),
(28, 'CIRUGIA BARIATRICA MANGA GASTRICA LAP', 4, NULL),
(29, 'CISTECTOMIA PARCIAL Y/O TRAT. QUIR. DE D', 4, NULL),
(30, 'CISTOPLASTIA, PROC. COMPLETO', 5, NULL),
(31, 'COAGULACION DE NUCLEOS O VIAS ENCEFALICA', 9, NULL),
(32, 'COLANGIOENTEROANASTOMOSIS INTRAHEPATICA', 2, NULL),
(33, 'COLANGIOPANCREATOGRAFÍA RETRÓGRADA C/S P', 2, NULL),
(34, 'COLECISTECTOMIA C/S COLANGIOGRAFIA OPERA', 3, NULL),
(35, 'COLECISTECTOMIA POR VIDEOLAPAROSCOPIA, P', 2.5, NULL),
(36, 'COLECTOMIA PARCIAL O HEMICOLECTOMIA', 4, NULL),
(37, 'COLGAJO SIMPLE UNICO', 2.5, NULL),
(38, 'COLGAJOS SIMPLES 2 O MAS', 2, NULL),
(39, 'CONDILOMAS ANALES, TRAT. QUIR.(PARA ELEC', 2, NULL),
(40, 'CONIZACION Y/O AMPUTACION DEL CUELLO, DI', 1.5, NULL),
(41, 'CONTRACTURA DUPUYTREN,TRAT. QUIR.,CADA T', 2, NULL),
(42, 'CRANEOPLASTIA CON PROTESIS (NO INCLUYE E', 3, NULL),
(43, 'CRANIECTOMIAS C/S REMODELACION OSEA', 7, NULL),
(44, 'CURETAJE POR GRUPO', 2, NULL),
(45, 'DE COMPLEJIDAD MAYOR: INCLUYE REEMPLAZ', 7, NULL),
(46, 'DECORTICACION DE CUERDAS VOCALES C/MICRO', 2.5, NULL),
(47, 'DEDOS EN GATILLO,TRAT. QUIR.,CUALQUIER N', 2, NULL),
(48, 'DESCENSO TESTICULAR CON O SIN HERNIA, CU', 2, NULL),
(49, 'DISFUNCION PATELO-FEMORAL,REALINEAMIENTO', 3, NULL),
(50, 'DRENAJE PERCUTÁNEO O ENDOSCÓPICO DE HIDR', 2, NULL),
(51, 'ECTROPION, PLASTIA DE', 1, NULL),
(52, 'EMBARAZO TUBARIO, TRAT. QUIR.', 3, NULL),
(53, 'ENCEFALICOS Y DE HIPOFISIS', 8, NULL),
(54, 'ENDOPROTESIS TOTAL DE CADERA', 3, NULL),
(55, 'ENDOPROTESIS TOTAL DE HOMBRO,(CUALQUIER', 3, NULL),
(56, 'ENDOPROTESIS TOTAL DE RODILLA, (CUALQUIE', 3, NULL),
(57, 'ENTERO-ENTEROANASTOMOSIS O ENTEROCOLOANA', 3.5, NULL),
(58, 'EPIFISIODESIS (FEMUR Y/O TIBIA)', 2, NULL),
(59, 'ESCARECTOMIA  HASTA 10N SUPERFICIE CORPO', 2, NULL),
(60, 'ESCARECTOMIA HASTA 5 N SUPERFICIE CORPOR', 2, NULL),
(61, 'ESCOLIOSIS,TRAT.QUIR.,CUALQUIER VIA DE A', 6, NULL),
(62, 'ESOFAGECTOMIA TOTAL CON ESOFAGOSTOMIA, G', 4, NULL),
(63, 'ESTENOSIS LARINGOTRAQUEALES Y/O FARÍNGEA', 2.5, NULL),
(64, 'ESTRABISMO TRAT. QX. COMPLETO UNI O BIL', 2, NULL),
(65, 'EXODONCIA DE PIEZAS INCLUIDAS', 2, NULL),
(66, 'EXTIRPACIÓN DE GLÁNDULA SALIVAL SUBMANDI', 2.5, NULL),
(67, 'EXTIRPACIÓN TOTAL O PARCIAL DE LA GLÁNDU', 2, NULL),
(68, 'FACOERESIS EXTRACAPSULAR CON IMPLANTE D', 1, NULL),
(69, 'FENESTRACION, SEPTOSTOMIA O COAGULACION', 5, NULL),
(70, 'FIJACION DE COLUMNA (CERVICAL-DORSAL-LUM', 4.5, NULL),
(71, 'FISTULA ANORRECTAL, TRAT.QUIR.', 2, NULL),
(72, 'FISTULA ARTERIOVENOSA (DE BRESCIA O SIMI', 2, NULL),
(73, 'FISTULA ARTERIOVENOSA CONGENITA O TRAUMA', 2, NULL),
(74, 'FRACTURA CUELLO HUMERAL, TRAT. QUIR.', 3, NULL),
(75, 'FRACTURA DE CLAVICULA, OSTEOSINTESIS', 3, NULL),
(76, 'FRACTURA O PSEUDOARTROSIS ESCAFOIDES,TRA', 3, NULL),
(77, 'FRACTURA ROTULA: OSTEOSINTESIS O PATELEC', 3, NULL),
(78, 'FRACTURAS CONDILEAS O DE PLATILLOS TIBIA', 3.5, NULL),
(79, 'GASTRECTOMIA TOTAL O SUBTOTAL AMPLIADA (', 4, NULL),
(80, 'GASTRODUODENOSCOPIA (INCLUYE ESOFAGOSCOP', 2, NULL),
(81, 'GASTROENTEROANASTOMOSIS, CUALQUIER TECNI', 4, NULL),
(82, 'GASTROSTOMIA PERCUTANEA TECNICA SELDINGE', 2, NULL),
(83, 'GASTROTOMIA Y/O GASTROSTOMIA (PROC. AUT.', 3, NULL),
(84, 'GINECOMASTIA, CORRECCION PLASTICA', 2, NULL),
(85, 'HALLUX VALGUS O RIGIDUS,TRAT.QUIR. COMPL', 2, NULL),
(86, 'HEMATOMA, EMPIEMA O COLECCION SUBDURAL,', 4, NULL),
(87, 'HEMORROIDECTOMIA (INCLUYE OTRAS OPERA-', 2, NULL),
(88, 'HERNIA ABDOMINAL POR LAPAROTOMÍA (NO INC', 2.5, NULL),
(89, 'HERNIA INGUINAL, CRURAL, UMBILICAL, DE L', 2, NULL),
(90, 'HERNIA NUCLEO PULPOSO, ESTENORRAQUIS, AR', 3, NULL),
(91, 'HIDROCELE Y/O HEMATOCELE, INCLUYE QUISTE', 2, NULL),
(92, 'HISTERECTOMIA POR VIA VAGINAL', 2.5, NULL),
(93, 'HISTERECTOMIA RADICAL CON DISECCION PELV', 3.5, NULL),
(94, 'HISTERECTOMIA TOTAL C/INTERVENCION INCON', 3, NULL),
(95, 'HISTERECTOMIA TOTAL O AMPLIADA POR VIA A', 2.5, NULL),
(96, 'HISTEROSCOPÍA DIAGNÓSTICA (PROC. AUT.)', 2, NULL),
(97, 'IMPLANTE CATHETER RESERVOREO PARA QMT', 1.5, NULL),
(98, 'IMPLANTE FILTROS VENOSOS', 2, NULL),
(99, 'INCOMPETENCIA CERVICAL TRAT. QUIR.', 2, NULL),
(100, 'INCONTINENCIA URINARIA DE ESFUERZO O D.I', 2, NULL),
(101, 'INESTABILIDAD CRONICA DE RODILLA, RECONS', 3, NULL),
(102, 'INFILTRACION FACETARIA COLUMNA Y RADICUL', 1.5, NULL),
(103, 'INJERTOS HASTA 10% SUPERFICIE CORPORAL R', 2, NULL),
(104, 'INSTALACION CATETER TUNELIZADO DIALISIS', 1.5, NULL),
(105, 'INSTALACIÓN DE CATÉTER CON RESERVORIO SU', 1.5, NULL),
(106, 'LAMINECTOMIA DESCOMPRESIVA', 2.5, NULL),
(107, 'LAPAROTOMIA EXPLORADORA, C/S LIBERACION', 3, NULL),
(108, 'LIBERACIÓN QUIRÚRGICA DE NERVIO PERIFÉRI', 2, NULL),
(109, 'LIGAMENTO ANCHO: ABSCESOS Y/O HEMATOMAS',2.5, NULL),
(110, 'LOBECTOMIA HEPATICA (PROC. AUT.)', 4.5, NULL),
(111, 'LOBECTOMIA O BILOBECTOMIA', 4.5, NULL),
(112, 'LUXACION ACROMIO-CLAVICULAR O ESTERNO CL', 3, NULL),
(113, 'LUXACION RECIDIVANTE, TRAT. QUIR.', 2.5, NULL),
(114, 'LUXOFRACTURA TOBILLO, CUALQUIER TIPO, OS', 2.5, NULL),
(115, 'MAMOPLASTIA DE REDUCCION', 3.5, NULL),
(116, 'MAMOPLASTIA ESTETICA BILATERIAL DE AUMEN', 3, NULL),
(117, 'MASTECTOMIA PARCIAL (CUADRANTECTOMIA O S', 2, NULL),
(118, 'MASTECTOMIA RADICAL O TUMORECTOMIA C/VAC', 3, NULL),
(119, 'MEDIASTINICOS', 2, NULL),
(120, 'MENISCECTOMIA U OTRAS INTERVENCIONES POR', 2, NULL),
(121, 'MIOMECTOMIA', 2, NULL),
(122, 'NEFRECTOMÍA PARCIAL CUALQUIER VÍA Y TÉCN', 5, NULL),
(123, 'NEFRECTOMIA RADICAL POR CANCER RENAL, TR', 4.5, NULL),
(124, 'NEUROLISIS CON TECNICA MICROQUIRURGICA', 3, NULL),
(125, 'NEUROLISIS EXTERNA', 2.5, NULL),
(126, 'OOFORECTOMIA PARCIAL O TOTAL, UNI O BILA', 2.5, NULL),
(127, 'OPERACION DE SALVATAJE CADERA, COLUMNA O', 3.5, NULL),
(128, 'ORQUIDECTOMÍA AMPLIADA POR CÁNCER TESTIC', 2, NULL),
(129, 'ORTEJOS EN GARRA,TRAT.QUIR.,CUALQ.NUMERO', 2, NULL),
(130, 'OSTEOMIELITIS AGUDA HEMATOGENA, DRENAJE', 3, NULL),
(131, 'OSTEOMIELITIS CRONICA HUESOS LARGOS, LEG', 2.5, NULL),
(132, 'OSTEOMIELITIS, LIMPIEZA QUIRURGICA', 2, NULL),
(133, 'OSTEOSINTESIS DIAFISIARIA (CUALQUIER TEC', 3, NULL),
(134, 'OSTEOSINTESIS METACARPIANAS O DE FALANGE', 2, NULL),
(135, 'OSTEOSINTESIS RADIO, (CUALQUIER TECNICA)', 3, NULL),
(136, 'OSTEOSINTESIS SUPRA O INTERCONDILEA (CUA', 2.5, NULL),
(137, 'OSTEOSINTESIS TIBIO-PERONE  (CUALQUIER T', 3.5, NULL),
(138, 'OSTEOSINTESIS, FRACT.CERRADA CUBITO Y/O', 3, NULL),
(139, 'OSTEOTOMIA CORRECTORA', 2, NULL),
(140, 'OTRAS DERIVACIONES: FEMORO-FEMORAL, AXIL', 3.5, NULL),
(141, 'PABELLON DE ESTADIA MINIMA', 2, NULL),
(142, 'PANCREATECTOMIA TOTAL C/S ESPLENECTOMIA', 7, NULL),
(143, 'PARATIROIDECTOMÍA O EXPLORACIÓN DE PARAT', 2.5, NULL),
(144, 'PARATIROIDECTOMÍA O EXPLORACIÓN PARATIRO', 3.5, NULL),
(145, 'PAROTIDECTOMÍA SUPRAFACIAL CON DISECCIÓN', 3, NULL),
(146, 'PAROTIDECTOMÍA TOTAL, CON DISECCIÓN Y PR', 3, NULL),
(147, 'PIE BOT U OTRAS MALFORMACIONES CONGENITA', 2, NULL),
(148, 'PIE PLANO, TRAT. QUIR. (CUALQUIER TECNIC', 2, NULL),
(149, 'PIELOTOMIA EXPLORADORA Y/O TERAPEUTICA (', 4.5, NULL),
(150, 'PLASTIAS EN Z, HASTA 3', 2, NULL),
(151, 'PLEURODESIS POR TORACOTOMIA', 2, NULL),
(152, 'POLIDACTILIA, EXTIRPACION Y PLASTIA UN L', 2, NULL),
(153, 'PROLAPSO ANTERIOR Y/O POSTERIOR C/S TRAT', 2, NULL),
(154, 'PROLAPSO ANTERIOR Y/O POSTERIOR CON REPA', 2, NULL),
(155, 'PTERIGION Y/O PSEUDOPTERIGION O SU RECID', 1, NULL),
(156, 'PTOSIS, TRAT. QUIR.', 2, NULL),
(157, 'PUNCION LUMBAR PARA MIELOGRAMAS', 1.5, NULL),
(158, 'QUERATECTOMIA FOTORREFRACTIVA O FOTOTERA', 1, NULL),
(159, 'QUISTE SACROCOXIGEO, TRAT. QUIR.', 2, NULL),
(160, 'QUISTES SINOVIALES DE VAINAS FLEXORAS, B', 2, NULL),
(161, 'RASPADO UTERINO DIAGNOSTICO O TERAPEUT', 1, NULL),
(162, 'RECONSTITUCION  DE TRANSITO EN 2 TIEMPO', 3, NULL),
(163, 'RECONSTITUCION  TRANSITO POST OPERACION', 3.5, NULL),
(164, 'RECONSTRUCCION AREOLA Y/O PEZON C/S PLAS', 2.5, NULL),
(165, 'RECONSTRUCCION MAMARIA', 4, NULL),
(166, 'REPARACION DE FISTULA DE LCR', 3, NULL),
(167, 'REPARACION QUIR. DE VASOS ARTERIALES Y/O', 3.5, NULL),
(168, 'RESECCION DE COSTILLAS Y/O PARED COSTAL', 3, NULL),
(169, 'RESECCION DE PARED COSTAL C/PLASTIA (TOR', 2.5, NULL),
(170, 'RESECCION ENDOSCOPICA DE CANCER VESICAL', 2.5, NULL),
(171, 'RESECCIÓN PLÁSTICA DE HASTA 2 CICATRICES', 2, NULL),
(172, 'RESECCIONES SEGMENTARIAS ANATÓMI DE PULM', 3, NULL),
(173, 'RETIRO CATHETER TUNELIZADO DIALISIS', 2, NULL),
(174, 'RETIRO DE CATÉTER CON RESERVORIO SUBCUTÁ', 2, NULL),
(175, 'RETIRO DE ENDOPROTESIS U OSTEOSINTESIS I', 2, NULL),
(176, 'RETIRO DE PLACAS RECTAS O ANGULADAS', 2, NULL),
(177, 'RETIRO DE TORNILLOS, CLAVOS, AGUJAS DE O', 2, NULL),
(178, 'RINOPLASTIA Y/O SEPTOPLASTIA, CUALQUIER', 3.5, NULL),
(179, 'RIZOTOMIA (CUALQUIER TECNICA)', 2, NULL),
(180, 'RUPTURA MANGUITO ROTADORES, TRAT. QUIR.', 3.5, NULL),
(181, 'RUPTURA TENDON DE AQUILES O TIBIAL POSTE', 2, NULL),
(182, 'SAFENECTOMÍA INTERNA Y/O EXTERNA,', 2, NULL),
(183, 'SAFENECTOMIA POR RADIOFRECUENCIA', 2, NULL),
(184, 'SALPINGECTOMIA UNI O BILATERAL', 2.5, NULL),
(185, 'SECCION DE  NERVIO, REPARACION CON INJER', 3, NULL),
(186, 'SENO ESFENOIDAL, ABERTURA POR CUALQUIER', 3, NULL),
(187, 'SIMPATECTOMIA CERVICO-TORACICA', 2, NULL),
(188, 'SUPRARRENALECTOMIA UNILATERAL', 4.5, NULL),
(189, 'TENORRAFIA EXTENSORES', 2, NULL),
(190, 'TENORRAFIA EXTENSORES O TENOTOMIA DE ALA', 3, NULL),
(191, 'TENORRAFIA O INJERTOS FLEXORES', 2, NULL),
(192, 'TIMPANOPLASTIA FUNCIONAL (CUALQUIER TIPO', 3.5, NULL),
(193, 'TIROIDECTOMIA BILATERAL TOTAL', 3.5, NULL),
(194, 'TIROIDECTOMÍA TOTAL AMPLIADA (INCLUYE EX', 3, NULL),
(195, 'TORACOTOMIA EXPLORADORA, C/S BIOPSIA, C/', 2, NULL),
(196, 'TRAQUEOSTOMIA (PROC. AUT.)', 2, NULL),
(197, 'TRATAMIENTO QUIRÚRGICO DE MUCOSITIS TIMP', 2, NULL),
(198, 'TRATAMIENTO QUIRÚRGICO PÓLIPO NASAL', 3, NULL),
(199, 'TUMOR BENIGNO DE CUERDAS VOCALES, TRAT.', 2, NULL),
(200, 'TUMOR OSEO, RESECCION EN BLOQUE, C/S OST', 2, NULL),
(201, 'TUMOR Y/O QUISTE RETROPERITONEAL', 3, NULL),
(202, 'TUMORES MALIGNOS DE PROSTATA O VESICULAS', 6, NULL),
(203, 'TUMORES O QUISTES DE MEDIASTINO (ANTERIO', 3, NULL),
(204, 'TUMORES O QUISTES O LESIONES PSEUDOQUIST', 2, NULL),
(205, 'URETERO-LITOTOMIA ENDOSCOPICA C/URETEROS', 2.5, NULL),
(206, 'URETRECTOMÍA Y /O PLASTÍA ABIERTA DE UR', 4, NULL),
(207, 'URETROTOMIA INTERNA Y/O URETROLITOTOMIA', 2.5, NULL),
(208, 'VACIAMIENTO (DISECCIÓN) RADICAL CUELLO (', 3, NULL),
(209, 'VACIAMIENTO Y CURETAJE QUIRURGICO DE LES', 2.5, NULL),
(210, 'VARICOCELE UNILATERAL Y/O DENERVACIÓN CO', 2, NULL),
(211, 'VASECTOMIA BILATERAL, (PROC. AUT.) (LA V', 1, NULL),
(212, 'VENTRICULOCISTERNOSTOMIA', 4.5, NULL),
(213, 'VENTRICULOSTOMIA O INSTALACION DE DERIVA', 2.5, NULL),
(214, 'VIDEOLAPAROSCOPIA GINECOLOGICA EXPLORADO', 2.5, NULL),
(215, 'VITRECTOMIA C/RETINOTOMIA (C/S INYECCION', 3, NULL),
(216, 'VULVECTOMIA SIMPLE', 1, NULL),
(217, 'YUGULAR SIMPLE', 3, NULL),
(218, 'ANEURISMA AORTICO ABDOMINAL TRAT. QUIR.', 4, (SELECT id FROM medical_specialty WHERE code = 'CARDIOL'))
ON CONFLICT (id) DO UPDATE SET specialty_id = EXCLUDED.specialty_id;

-- Poblar lotes
INSERT INTO batch (id, expiration_date, amount, supplier, store_id, qr_code, surgery_id, location_type, location_id) VALUES
(1, '2026-12-31', 10, 'Proveedor Uno', 1, 'BATCH_1_1', 1, 'store', 1),
(2, '2025-12-31', 5, 'Proveedor Dos', 2, 'BATCH_2_1', 2, 'store', 2),
(3, '2026-06-30', 15, 'Proveedor Tres', 1, 'BATCH_3_1', 5, 'store', 1),
(4, '2025-10-15', 8, 'Proveedor Cuatro', 2, 'BATCH_4_1', 10, 'store', 2)
ON CONFLICT (id) DO NOTHING;

-- Poblar insumos médicos - TODOS EN ESTADO DISPONIBLE
INSERT INTO medical_supply (code, batch_id, qr_code, status, location_type, location_id) VALUES
(1001, 1, 'SUPPLY_1_1', 'disponible', 'store', 1),
(1001, 1, 'SUPPLY_2_1', 'disponible', 'store', 1),
(1001, 1, 'SUPPLY_3_1', 'disponible', 'store', 1),
(1001, 1, 'SUPPLY_4_1', 'disponible', 'store', 1),
(1001, 1, 'SUPPLY_5_1', 'disponible', 'store', 1),
(1001, 1, 'SUPPLY_6_1', 'disponible', 'store', 1),
(1001, 1, 'SUPPLY_7_1', 'disponible', 'store', 1),
(1001, 1, 'SUPPLY_8_1', 'disponible', 'store', 1),
(1001, 1, 'SUPPLY_9_1', 'disponible', 'store', 1),
(1001, 1, 'SUPPLY_10_1', 'disponible', 'store', 1),
(1002, 2, 'SUPPLY_11_1', 'disponible', 'store', 2),
(1002, 2, 'SUPPLY_12_1', 'disponible', 'store', 2),
(1002, 2, 'SUPPLY_13_1', 'disponible', 'store', 2),
(1002, 2, 'SUPPLY_14_1', 'disponible', 'store', 2),
(1002, 2, 'SUPPLY_15_1', 'disponible', 'store', 2),
(1003, 3, 'SUPPLY_16_1', 'disponible', 'store', 1),
(1003, 3, 'SUPPLY_17_1', 'disponible', 'store', 1),
(1003, 3, 'SUPPLY_18_1', 'disponible', 'store', 1),
(1003, 3, 'SUPPLY_19_1', 'disponible', 'store', 1),
(1003, 3, 'SUPPLY_20_1', 'disponible', 'store', 1),
(1004, 4, 'SUPPLY_21_1', 'disponible', 'store', 2),
(1004, 4, 'SUPPLY_22_1', 'disponible', 'store', 2),
(1004, 4, 'SUPPLY_23_1', 'disponible', 'store', 2),
(1004, 4, 'SUPPLY_24_1', 'disponible', 'store', 2),
(1004, 4, 'SUPPLY_25_1', 'disponible', 'store', 2)
ON CONFLICT (qr_code) DO NOTHING;

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
    '$2a$10$NA3QLOvkwhpcs.X4KxjONObslo1LreYA6qAzdQcqxRrD4ktjBrpmO',
    'admin',
    1,
    true,
    EXTRACT(EPOCH FROM NOW()),
    EXTRACT(EPOCH FROM NOW())
) ON CONFLICT (rut) DO NOTHING;

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
    '$2a$10$NA3QLOvkwhpcs.X4KxjONObslo1LreYA6qAzdQcqxRrD4ktjBrpmO',
    'pabellón',
    1,
    true,
    EXTRACT(EPOCH FROM NOW()),
    EXTRACT(EPOCH FROM NOW())
) ON CONFLICT (rut) DO NOTHING;

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
    '$2a$10$NA3QLOvkwhpcs.X4KxjONObslo1LreYA6qAzdQcqxRrD4ktjBrpmO',
    'encargado de bodega',
    1,
    true,
    EXTRACT(EPOCH FROM NOW()),
    EXTRACT(EPOCH FROM NOW())
) ON CONFLICT (rut) DO NOTHING;

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
    '$2a$10$NA3QLOvkwhpcs.X4KxjONObslo1LreYA6qAzdQcqxRrD4ktjBrpmO',
    'enfermera',
    1,
    true,
    EXTRACT(EPOCH FROM NOW()),
    EXTRACT(EPOCH FROM NOW())
) ON CONFLICT (rut) DO NOTHING;

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
    '$2a$10$NA3QLOvkwhpcs.X4KxjONObslo1LreYA6qAzdQcqxRrD4ktjBrpmO',
    'doctor',
    1,
    true,
    EXTRACT(EPOCH FROM NOW()),
    EXTRACT(EPOCH FROM NOW())
) ON CONFLICT (rut) DO NOTHING;

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
    '$2a$10$NA3QLOvkwhpcs.X4KxjONObslo1LreYA6qAzdQcqxRrD4ktjBrpmO',
    'admin',
    1,
    true,
    EXTRACT(EPOCH FROM NOW()),
    EXTRACT(EPOCH FROM NOW())
) ON CONFLICT (rut) DO NOTHING;

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
    '44444444-4',
    'Pavedad',
    'pavedad@meditrack.com',
    '$2a$10$NA3QLOvkwhpcs.X4KxjONObslo1LreYA6qAzdQcqxRrD4ktjBrpmO', -- Reemplazar con contraseña hasheada real
    'pavedad',
    1,
    true,
    EXTRACT(EPOCH FROM NOW()),
    EXTRACT(EPOCH FROM NOW())
)
ON CONFLICT (rut) DO NOTHING;

INSERT INTO batch_history (date_time, change_details, previous_values, new_values, user_name, batch_id, user_rut, batch_number) VALUES
('2025-08-16 10:00:00', 'Lote creado', NULL, '{"expiration_date": "2026-12-31", "amount": 10, "supplier": "Proveedor Uno", "store_id": 1}', 'Administrador del Sistema', 1, '12345678-9', 1),
('2025-08-16 11:00:00', 'Lote creado', NULL, '{"expiration_date": "2025-12-31", "amount": 5, "supplier": "Proveedor Dos", "store_id": 2}', 'Usuario Pabellón', 2, '87654321-0', 2),
('2025-08-16 12:00:00', 'Lote creado', NULL, '{"expiration_date": "2026-06-30", "amount": 15, "supplier": "Proveedor Tres", "store_id": 1}', 'Encargado Bodega', 3, '11111111-1', 3),
('2025-08-16 13:00:00', 'Lote creado', NULL, '{"expiration_date": "2025-10-15", "amount": 8, "supplier": "Proveedor Cuatro", "store_id": 2}', 'Administrador del Sistema', 4, '12345678-9', 4),
('2025-08-16 14:00:00', 'Cantidad actualizada', '{"amount": 10}', '{"amount": 8}', 'Encargado Bodega', 1, '11111111-1', 1)
ON CONFLICT DO NOTHING;

-- Poblar información extendida de doctores
INSERT INTO doctor_info (
    user_rut, medical_license, specialization, specialty_id, 
    years_of_experience, phone, is_available
) VALUES (
    '33333333-3', 
    'LIC-33333333', 
    'Cirugía General', 
    (SELECT id FROM medical_specialty WHERE code = 'CIR_GEN'),
    15, 
    '+56 9 1234 5678', 
    TRUE
) ON CONFLICT (user_rut) DO NOTHING;

-- Poblar insumos típicos por cirugía (ejemplos)
-- Para cirugías generales - guantes y gasas típicamente requeridos
-- Para cirugías plásticas - guantes, mascarillas y gasas
-- Para traumatología - jeringas, agujas y gasas
INSERT INTO surgery_typical_supply (surgery_id, supply_code, typical_quantity, is_required, notes) VALUES
((SELECT id FROM surgery WHERE name LIKE 'CIRUGÍA ABIERTA%' LIMIT 1), 1001, 50, TRUE, 'Guantes estériles requeridos'),
((SELECT id FROM surgery WHERE name LIKE 'CIRUGÍA ABIERTA%' LIMIT 1), 1005, 20, TRUE, 'Gasas estériles requeridas'),
((SELECT id FROM surgery WHERE name LIKE 'COLGAJOS%' LIMIT 1), 1001, 30, TRUE, 'Guantes estériles'),
((SELECT id FROM surgery WHERE name LIKE 'COLGAJOS%' LIMIT 1), 1002, 10, TRUE, 'Mascarillas N95'),
((SELECT id FROM surgery WHERE name LIKE 'COLGAJOS%' LIMIT 1), 1005, 15, TRUE, 'Gasas estériles'),
((SELECT id FROM surgery WHERE name LIKE 'REDUCCION ABIERTA%' LIMIT 1), 1003, 5, FALSE, 'Jeringas para anestesia'),
((SELECT id FROM surgery WHERE name LIKE 'REDUCCION ABIERTA%' LIMIT 1), 1004, 10, TRUE, 'Agujas estériles'),
((SELECT id FROM surgery WHERE name LIKE 'REDUCCION ABIERTA%' LIMIT 1), 1005, 25, TRUE, 'Gasas para limpieza')
ON CONFLICT (surgery_id, supply_code) DO NOTHING;

INSERT INTO supply_request (
    request_number, pavilion_id, requested_by, requested_by_name,
    request_date, surgery_datetime, status, notes, medical_center_id,
    surgeon_id, surgeon_name, surgery_id, specialty_id
) VALUES (
    'SOL-20250120140000', 
    1, 
    '33333333-3', 
    'Dr. Carlos Pérez',
    NOW() - INTERVAL '1 hour', 
    NOW() + INTERVAL '23 hours', 
    'pendiente_pavedad', 
    'Solicitud de prueba para implementación de trazabilidad QR', 
    1,
    '33333333-3', 
    'Dr. Carlos Pérez', 
    (SELECT id FROM surgery WHERE name LIKE 'CIRUGÍA ABIERTA%' LIMIT 1),
    (SELECT id FROM medical_specialty WHERE code = 'CIR_GEN')
);

-- Insertar items de ejemplo
INSERT INTO supply_request_item (supply_request_id, supply_code, supply_name, quantity_requested, is_pediatric) VALUES 
(1, 1001, 'Guantes', 50, FALSE),
(1, 1002, 'Mascarillas', 100, TRUE);

SELECT 'ID máximo actual en batch:' as info, COALESCE(MAX(id), 0) as max_id FROM batch;

SELECT setval('batch_id_seq', COALESCE((SELECT MAX(id) FROM batch), 0) + 1, false);

SELECT 'Secuencia después del reset:' as info, last_value, is_called FROM batch_id_seq;

SELECT 'QR codes duplicados:' as info, qr_code, COUNT(*) as count 
FROM batch 
WHERE qr_code IS NOT NULL 
GROUP BY qr_code 
HAVING COUNT(*) > 1;

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
    NOW() - INTERVAL '30 days' AS date_time,
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

SELECT 
    'Registros de historial agregados:' as info,
    COUNT(*) as total_records
FROM supply_history 
WHERE user_rut = 'SYSTEM-INIT';

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

INSERT INTO store_inventory_summary (
    store_id,
    batch_id,
    supply_code,
    surgery_id,
    original_amount,
    current_in_store,
    total_transferred_out,
    total_returned_in,
    total_consumed_in_store,
    created_at,
    updated_at
)
SELECT 
    b.store_id,
    b.id as batch_id,
    ms.code as supply_code,
    b.surgery_id,
    b.amount as original_amount,
    b.amount as current_in_store,
    0 as total_transferred_out,
    0 as total_returned_in,
    0 as total_consumed_in_store,
    NOW() as created_at,
    NOW() as updated_at
FROM batch b
JOIN medical_supply ms ON ms.batch_id = b.id
WHERE b.location_type = 'store'
GROUP BY b.store_id, b.id, ms.code, b.surgery_id, b.amount
ON CONFLICT (batch_id) DO NOTHING;

SELECT 
    'Registros de resumen de bodega:' as info,
    COUNT(*) as total_records
FROM store_inventory_summary;

SELECT setval('surgery_id_seq', (SELECT COALESCE(MAX(id), 0) FROM surgery));

SELECT setval('medical_center_id_seq', (SELECT COALESCE(MAX(id), 0) FROM medical_center));

SELECT setval('pavilion_id_seq', (SELECT COALESCE(MAX(id), 0) FROM pavilion));

SELECT setval('store_id_seq', (SELECT COALESCE(MAX(id), 0) FROM store));

SELECT setval('batch_id_seq', (SELECT COALESCE(MAX(id), 0) FROM batch));

SELECT setval('medical_supply_id_seq', (SELECT COALESCE(MAX(id), 0) FROM medical_supply));

SELECT setval('supply_history_id_seq', (SELECT COALESCE(MAX(id), 0) FROM supply_history));

DO $$
DECLARE
    max_id INTEGER;
BEGIN
    IF EXISTS (SELECT 1 FROM pg_class WHERE relname = 'supply_transfer') THEN
        SELECT COALESCE(MAX(id), 0) INTO max_id FROM supply_transfer;
        IF max_id > 0 THEN
            PERFORM setval('supply_transfer_id_seq', max_id);
        END IF;
    END IF;
END $$;

DO $$
DECLARE
    max_id INTEGER;
BEGIN
    IF EXISTS (SELECT 1 FROM pg_class WHERE relname = 'store_inventory_summary') THEN
        SELECT COALESCE(MAX(id), 0) INTO max_id FROM store_inventory_summary;
        IF max_id > 0 THEN
            PERFORM setval('store_inventory_summary_id_seq', max_id);
        END IF;
    END IF;
END $$;

DO $$
DECLARE
    max_id INTEGER;
BEGIN
    IF EXISTS (SELECT 1 FROM pg_class WHERE relname = 'pavilion_inventory_summary') THEN
        SELECT COALESCE(MAX(id), 0) INTO max_id FROM pavilion_inventory_summary;
        IF max_id > 0 THEN
            PERFORM setval('pavilion_inventory_summary_id_seq', max_id);
        END IF;
    END IF;
END $$;

DO $$
DECLARE
    max_id INTEGER;
BEGIN
    IF EXISTS (SELECT 1 FROM pg_class WHERE relname = 'batch_history') THEN
        SELECT COALESCE(MAX(id), 0) INTO max_id FROM batch_history;
        IF max_id > 0 THEN
            PERFORM setval('batch_history_id_seq', max_id);
        END IF;
    END IF;
END $$;

DO $$
DECLARE
    max_id INTEGER;
BEGIN
    IF EXISTS (SELECT 1 FROM pg_class WHERE relname = 'supply_request') THEN
        SELECT COALESCE(MAX(id), 0) INTO max_id FROM supply_request;
        IF max_id > 0 THEN
            PERFORM setval('supply_request_id_seq', max_id);
        END IF;
    END IF;
END $$;

DO $$
DECLARE
    max_id INTEGER;
BEGIN
    IF EXISTS (SELECT 1 FROM pg_class WHERE relname = 'supply_request_item') THEN
        SELECT COALESCE(MAX(id), 0) INTO max_id FROM supply_request_item;
        IF max_id > 0 THEN
            PERFORM setval('supply_request_item_id_seq', max_id);
        END IF;
    END IF;
END $$;

SELECT 'surgery_id_seq' as sequence_name, last_value FROM surgery_id_seq
UNION ALL
SELECT 'medical_center_id_seq', last_value FROM medical_center_id_seq
UNION ALL
SELECT 'pavilion_id_seq', last_value FROM pavilion_id_seq
UNION ALL
SELECT 'store_id_seq', last_value FROM store_id_seq
UNION ALL
SELECT 'batch_id_seq', last_value FROM batch_id_seq
UNION ALL
SELECT 'medical_supply_id_seq', last_value FROM medical_supply_id_seq
UNION ALL
SELECT 'supply_history_id_seq', last_value FROM supply_history_id_seq
UNION ALL
SELECT 'medical_specialty_id_seq', last_value FROM medical_specialty_id_seq
UNION ALL
SELECT 'surgery_typical_supply_id_seq', last_value FROM surgery_typical_supply_id_seq
ORDER BY sequence_name;

-- Ajustar secuencias de configuración médica si es necesario
SELECT setval('medical_specialty_id_seq', (SELECT COALESCE(MAX(id), 0) FROM medical_specialty));
SELECT setval('surgery_typical_supply_id_seq', (SELECT COALESCE(MAX(id), 0) FROM surgery_typical_supply));