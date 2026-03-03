-- Script de Poblado Completo para MediTrack - VERSIÓN MEJORADA
-- Fecha: 2025-02-05
-- Descripción: Poblado completo con relaciones coherentes entre todas las tablas
--
-- MEJORAS IMPLEMENTADAS:
-- - Códigos de insumo ampliados y categorizados por tipo de cirugía (54 códigos)
-- - 50 lotes distribuidos entre Bodega Central (40) y Bodega Consignación (10)
-- - Lotes con mejor distribución y relación con cirugías específicas
-- - Generación automática de insumos médicos con QR codes únicos
-- - Surgery_typical_supply completo para cirugías principales (120+ relaciones)
-- - 6 solicitudes de ejemplo en diferentes estados del flujo
-- - 58 transferencias a pabellones con diferentes antigüedades:
--   * 25 transferencias >48 horas (PENDIENTES RETORNO) - Todas al Pabellón 1
--   * 15 transferencias <8 horas (VÁLIDAS)
--   * 18 transferencias >72 horas (CRÍTICAS - RETORNO URGENTE) - 15 al Pabellón 1, resto distribuidas
--   * NOTA: El Pabellón 1 tiene 70+ items distribuidos en múltiples lotes para probar:
--     - Stock 0 (rojo): Batch 8 y 17 consumidos totalmente
--     - Stock bajo 1-4 (naranja): Batch 5 (1 unidad), Batch 15 (2 unidades), Batch 3 (3 unidades), Batch 12 (4 unidades)
--     - Stock 5-9 (amarillo): Batch 10 (6 unidades), Batch 18 (7 unidades), Batch 2 (8 unidades)
--     - Stock 10+ (verde): Batch 23 (11 unidades), Batch 1 (12 unidades), Batch 9 (15 unidades)
-- - Inventario activo en pabellones (pavilion_inventory_summary actualizado)
-- - Historial de transferencias documentado
-- - Configuración de 17 proveedores con diferentes políticas de vencimiento
-- - Reportes de gestión de retornos a bodega
--
-- DATOS PARA GESTIÓN DE RETORNOS:
-- Los insumos transferidos a pabellones hace más de 8 horas laborales
-- aparecerán automáticamente como "Pendientes de Retorno" en la vista
-- de gestión de retornos a bodega.

-- ============================================================================
-- CENTROS MÉDICOS Y PABELLONES (NO MODIFICAR)
-- ============================================================================
INSERT INTO medical_center (id, name, address, phone, email) VALUES
(1, 'Centro Médico Principal', 'Av. Principal 123', '+56 2 2345 6789', 'info@centromedico.cl'),
(2, 'Clínica Norte', 'Calle Norte 456', '987654321', 'norte@meditrack.com')
ON CONFLICT (id) DO NOTHING;

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

-- ============================================================================
-- BODEGAS
-- ============================================================================
INSERT INTO store (id, name, type, medical_center_id) VALUES
(1, 'Bodega Central', 'central', 1),
(2, 'Bodega Consignación', 'consignacion', 2)
ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- CÓDIGOS DE INSUMOS - AMPLIADO Y CATEGORIZADO
-- ============================================================================
INSERT INTO supply_code (code, name, code_supplier, critical_stock) VALUES
-- Insumos generales (1000-1099)
(1001, 'Guantes Quirúrgicos Estériles', 5001, 100),
(1002, 'Mascarillas N95', 5002, 50),
(1003, 'Jeringas 10ml', 5003, 80),
(1004, 'Agujas 21G', 5004, 100),
(1005, 'Gasas Estériles 10x10cm', 5005, 200),
(1006, 'Alcohol Gel 500ml', 5006, 30),
(1007, 'Batas Quirúrgicas Estériles', 5007, 40),
(1008, 'Campos Quirúrgicos Estériles', 5008, 60),
(1009, 'Guantes de Examinación No Estériles', 5009, 150),
(1010, 'Mascarillas Quirúrgicas Estándar', 5010, 200),

-- Suturas y material de cierre (1100-1199)
(1101, 'Sutura Vicryl 2-0', 5011, 40),
(1102, 'Sutura Nylon 3-0', 5012, 35),
(1103, 'Sutura Prolene 4-0', 5013, 30),
(1104, 'Grapas Quirúrgicas', 5014, 25),
(1105, 'Steri-Strips', 5015, 50),
(1106, 'Pegamento Tisular', 5016, 20),

-- Material para anestesia (1200-1299)
(1201, 'Tubo Endotraqueal 7.0', 5017, 30),
(1202, 'Tubo Endotraqueal 7.5', 5018, 30),
(1203, 'Mascarilla Laríngea #4', 5019, 25),
(1204, 'Cánula de Guedel #3', 5020, 40),
(1205, 'Sonda Nelaton #14', 5021, 50),

-- Catéteres y drenajes (1300-1399)
(1301, 'Catéter Venoso Central', 5022, 15),
(1302, 'Catéter Foley #16', 5023, 30),
(1303, 'Drenaje Blake #19', 5024, 20),
(1304, 'Catéter Venoso Periférico #20', 5025, 80),
(1305, 'Catéter Venoso Periférico #22', 5026, 80),

-- Material específico cardiovascular (1400-1499)
(1401, 'Bisturí Eléctrico Desechable', 5027, 15),
(1402, 'Compresas Hemostáticas', 5028, 25),
(1403, 'Ligaduras Vasculares', 5029, 30),

-- Material ortopédico (1500-1599)
(1501, 'Yeso Sintético 10cm', 5030, 20),
(1502, 'Yeso Sintético 15cm', 5031, 20),
(1503, 'Vendas Elásticas 10cm', 5032, 40),
(1504, 'Tornillos Ortopédicos 3.5mm', 5033, 15),
(1505, 'Placas Ortopédicas Rectas', 5034, 10),

-- Material oftalmológico (1600-1699)
(1601, 'Lente Intraocular', 5035, 15),
(1602, 'Viscoelástico Oftálmico', 5036, 20),
(1603, 'Sutura Nylon 10-0', 5037, 10),

-- Material urológico (1700-1799)
(1701, 'Sonda Vesical #18', 5038, 25),
(1702, 'Uréter Doble J #6', 5039, 15),
(1703, 'Balón Dilatación Uretral', 5040, 10),

-- Material ginecológico (1800-1899)
(1801, 'Espéculo Vaginal Desechable', 5041, 50),
(1802, 'Pinzas de Biopsia Cervical', 5042, 15),
(1803, 'Histerómetro Desechable', 5043, 20),

-- Material neurocirugía (1900-1999)
(1901, 'Craneotomo Desechable', 5044, 5),
(1902, 'Válvula Derivación Ventriculoperitoneal', 5045, 8),
(1903, 'Hemostáticos Neurológicos', 5046, 12)
ON CONFLICT (code) DO NOTHING;

-- ============================================================================
-- ESPECIALIDADES MÉDICAS
-- ============================================================================
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

-- ============================================================================
-- CIRUGÍAS (NO MODIFICAR - MANTENER TAL CUAL)
-- ============================================================================
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

-- ============================================================================
-- PROVEEDORES (debe ir antes de lotes por la FK)
-- ============================================================================
INSERT INTO supplier_config (id, supplier_name, notes) VALUES
(1,  'Proveedor Uno',            'Proveedor general de insumos'),
(2,  'Proveedor Dos',            'Proveedor general de insumos'),
(3,  'Proveedor Tres',           'Proveedor especializado'),
(4,  'Proveedor Cuatro',         'Proveedor general de insumos'),
(5,  'Proveedor Cinco',          'Proveedor de suturas'),
(6,  'Proveedor Seis',           'Proveedor de material de anestesia'),
(7,  'Proveedor Siete',          'Proveedor de catéteres y drenajes'),
(8,  'Proveedor Ocho',           'Proveedor de material ortopédico'),
(9,  'Proveedor Nueve',          'Proveedor de material oftalmológico'),
(10, 'Proveedor Diez',           'Proveedor de material urológico'),
(11, 'Proveedor Once',           'Proveedor de material ginecológico'),
(12, 'Proveedor Doce',           'Proveedor de material de neurocirugía'),
(13, 'Proveedor Consignación A', 'Bodega consignación'),
(14, 'Proveedor Consignación B', 'Bodega consignación'),
(15, 'Proveedor Consignación C', 'Bodega consignación'),
(16, 'Proveedor Consignación D', 'Bodega consignación - material especializado'),
(17, 'Proveedor Consignación E', 'Bodega consignación'),
(18, 'Proveedor Test A',         'Proveedor de prueba QA'),
(19, 'Proveedor Test B',         'Proveedor de prueba QA'),
(20, 'Proveedor Test C',         'Proveedor de prueba QA'),
(21, 'Proveedor Test D',         'Proveedor de prueba QA'),
(22, 'Proveedor Test E',         'Proveedor de prueba QA'),
(23, 'Proveedor Test F',         'Proveedor de prueba QA'),
(24, 'Proveedor Test G',         'Proveedor de prueba QA'),
(25, 'Proveedor Norte A',        'Proveedor zona norte QA'),
(26, 'Proveedor Norte B',        'Proveedor zona norte QA'),
(27, 'Proveedor Norte C',        'Proveedor zona norte QA'),
(28, 'Proveedor Norte D',        'Proveedor zona norte QA'),
(29, 'Proveedor Norte E',        'Proveedor zona norte QA'),
(30, 'Proveedor Norte F',        'Proveedor zona norte QA')
ON CONFLICT (supplier_name) DO NOTHING;

-- ============================================================================
-- LOTES - AMPLIADO Y DIVERSIFICADO
-- ============================================================================
-- Lotes para Bodega Central (ID 1)
INSERT INTO batch (id, expiration_date, amount, supplier_id, store_id, qr_code, surgery_id, location_type, location_id) VALUES
-- Insumos generales
(1, '2026-12-31', 150, 1, 1, 'BATCH_1001_001', NULL, 'store', 1),
(2, '2026-01-15', 200, 1, 1, 'BATCH_1002_001', NULL, 'store', 1), -- VENCIDO (hace 1 mes)
(3, '2026-03-10', 120, 2, 1, 'BATCH_1003_001', NULL, 'store', 1), -- PRÓXIMO A VENCIMIENTO (23 días)
(4, '2026-08-30', 150, 2, 1, 'BATCH_1004_001', NULL, 'store', 1),
(5, '2027-06-15', 300, 3, 1, 'BATCH_1005_001', NULL, 'store', 1),
(6, '2026-11-20', 10, 4, 1, 'BATCH_1006_001', NULL, 'store', 1), -- STOCK BAJO (10 unidades, transferir 9)
(7, '2027-01-10', 10, 1, 1, 'BATCH_1007_001', NULL, 'store', 1), -- STOCK BAJO (10 unidades, transferir 8)
(8, '2026-09-25', 100, 2, 1, 'BATCH_1008_001', NULL, 'store', 1),
-- Suturas (para cirugías generales y plásticas)
(9, '2028-12-31', 60, 5, 1, 'BATCH_1101_001', 1, 'store', 1),
(10, '2028-10-20', 50, 5, 1, 'BATCH_1102_001', 3, 'store', 1),
(11, '2028-11-15', 45, 5, 1, 'BATCH_1103_001', 26, 'store', 1),
-- Material anestesia
(12, '2027-07-30', 40, 6, 1, 'BATCH_1201_001', NULL, 'store', 1),
(13, '2027-07-30', 40, 6, 1, 'BATCH_1202_001', NULL, 'store', 1),
(14, '2027-09-15', 35, 6, 1, 'BATCH_1203_001', NULL, 'store', 1),
-- Catéteres y drenajes
(15, '2027-02-28', 25, 7, 1, 'BATCH_1301_001', 218, 'store', 1),
(16, '2027-04-10', 45, 7, 1, 'BATCH_1302_001', NULL, 'store', 1),
(17, '2027-05-20', 30, 7, 1, 'BATCH_1303_001', 34, 'store', 1),
-- Material ortopédico
(18, '2028-01-15', 30, 8, 1, 'BATCH_1501_001', 5, 'store', 1),
(19, '2028-01-15', 30, 8, 1, 'BATCH_1502_001', 74, 'store', 1),
(20, '2027-12-10', 60, 8, 1, 'BATCH_1503_001', 75, 'store', 1),
(21, '2029-06-30', 20, 8, 1, 'BATCH_1504_001', 133, 'store', 1),
(22, '2029-06-30', 15, 8, 1, 'BATCH_1505_001', 137, 'store', 1),
-- Material oftalmológico
(23, '2027-08-20', 20, 9, 1, 'BATCH_1601_001', 68, 'store', 1),
(24, '2027-10-05', 30, 9, 1, 'BATCH_1602_001', 68, 'store', 1),
(25, '2028-03-15', 15, 9, 1, 'BATCH_1603_001', 64, 'store', 1),
-- Material urológico
(26, '2027-11-25', 35, 10, 1, 'BATCH_1701_001', 11, 'store', 1),
(27, '2027-12-30', 20, 10, 1, 'BATCH_1702_001', 12, 'store', 1),
-- Material ginecológico
(28, '2027-06-10', 70, 11, 1, 'BATCH_1801_001', NULL, 'store', 1),
(29, '2028-02-20', 20, 11, 1, 'BATCH_1802_001', 40, 'store', 1),
(30, '2027-09-30', 25, 11, 1, 'BATCH_1803_001', 92, 'store', 1),
-- Material neurocirugía
(31, '2028-05-15', 8, 12, 1, 'BATCH_1901_001', 43, 'store', 1),
(32, '2028-08-20', 20, 12, 1, 'BATCH_1902_001', 213, 'store', 1), -- STOCK BAJO (20 unidades, transferir 17),
(33, '2027-12-15', 15, 12, 1, 'BATCH_1903_001', 31, 'store', 1),

-- Lotes para Bodega Consignación (ID 2)
(34, '2026-11-30', 100, 13, 2, 'BATCH_1001_002', NULL, 'store', 2),
(35, '2027-02-28', 150, 13, 2, 'BATCH_1002_002', NULL, 'store', 2),
(36, '2027-04-15', 80, 14, 2, 'BATCH_1003_002', NULL, 'store', 2),
(37, '2027-05-20', 100, 14, 2, 'BATCH_1004_002', NULL, 'store', 2),
(38, '2027-07-10', 200, 15, 2, 'BATCH_1005_002', NULL, 'store', 2),
(39, '2027-01-25', 40, 15, 2, 'BATCH_1007_002', NULL, 'store', 2),
(40, '2028-03-30', 35, 16, 2, 'BATCH_1101_002', 115, 'store', 2),
(41, '2028-03-30', 30, 16, 2, 'BATCH_1102_002', 116, 'store', 2),
(42, '2027-10-20', 25, 17, 2, 'BATCH_1301_002', NULL, 'store', 2),
(43, '2027-11-15', 35, 17, 2, 'BATCH_1302_002', NULL, 'store', 2),

-- Lotes adicionales para Bodega Central (antes en urgencias, ahora redistribuidos)
(44, '2026-09-30', 80, 1, 1, 'BATCH_1001_004', NULL, 'store', 1),
(45, '2026-10-15', 100, 1, 1, 'BATCH_1002_004', NULL, 'store', 1),
(46, '2027-01-20', 60, 2, 1, 'BATCH_1003_004', NULL, 'store', 1),
(47, '2026-12-05', 80, 3, 1, 'BATCH_1009_001', NULL, 'store', 1),
(48, '2027-03-10', 120, 3, 1, 'BATCH_1010_001', NULL, 'store', 1),
(49, '2027-02-15', 50, 4, 1, 'BATCH_1304_001', NULL, 'store', 1),
(50, '2027-02-15', 50, 4, 1, 'BATCH_1305_001', NULL, 'store', 1)
ON CONFLICT (id) DO NOTHING;

-- ============================================================================
-- INSUMOS MÉDICOS - GENERACIÓN MASIVA CON QR ÚNICOS
-- ============================================================================
-- Función auxiliar para generar insumos por lote
DO $$
DECLARE
    batch_record RECORD;
    i INTEGER;
    qr_prefix TEXT;
    current_supply_code INTEGER;
BEGIN
    -- Obtener el código de insumo correspondiente a cada lote
    FOR batch_record IN 
        SELECT 
            b.id as batch_id,
            b.amount,
            b.store_id,
            b.qr_code as batch_qr,
            CASE 
                -- Mapeo de lotes a códigos de insumo basado en el QR del lote
                WHEN b.qr_code LIKE '%1001%' THEN 1001
                WHEN b.qr_code LIKE '%1002%' THEN 1002
                WHEN b.qr_code LIKE '%1003%' THEN 1003
                WHEN b.qr_code LIKE '%1004%' THEN 1004
                WHEN b.qr_code LIKE '%1005%' THEN 1005
                WHEN b.qr_code LIKE '%1006%' THEN 1006
                WHEN b.qr_code LIKE '%1007%' THEN 1007
                WHEN b.qr_code LIKE '%1008%' THEN 1008
                WHEN b.qr_code LIKE '%1009%' THEN 1009
                WHEN b.qr_code LIKE '%1010%' THEN 1010
                WHEN b.qr_code LIKE '%1101%' THEN 1101
                WHEN b.qr_code LIKE '%1102%' THEN 1102
                WHEN b.qr_code LIKE '%1103%' THEN 1103
                WHEN b.qr_code LIKE '%1201%' THEN 1201
                WHEN b.qr_code LIKE '%1202%' THEN 1202
                WHEN b.qr_code LIKE '%1203%' THEN 1203
                WHEN b.qr_code LIKE '%1301%' THEN 1301
                WHEN b.qr_code LIKE '%1302%' THEN 1302
                WHEN b.qr_code LIKE '%1303%' THEN 1303
                WHEN b.qr_code LIKE '%1304%' THEN 1304
                WHEN b.qr_code LIKE '%1305%' THEN 1305
                WHEN b.qr_code LIKE '%1501%' THEN 1501
                WHEN b.qr_code LIKE '%1502%' THEN 1502
                WHEN b.qr_code LIKE '%1503%' THEN 1503
                WHEN b.qr_code LIKE '%1504%' THEN 1504
                WHEN b.qr_code LIKE '%1505%' THEN 1505
                WHEN b.qr_code LIKE '%1601%' THEN 1601
                WHEN b.qr_code LIKE '%1602%' THEN 1602
                WHEN b.qr_code LIKE '%1603%' THEN 1603
                WHEN b.qr_code LIKE '%1701%' THEN 1701
                WHEN b.qr_code LIKE '%1702%' THEN 1702
                WHEN b.qr_code LIKE '%1801%' THEN 1801
                WHEN b.qr_code LIKE '%1802%' THEN 1802
                WHEN b.qr_code LIKE '%1803%' THEN 1803
                WHEN b.qr_code LIKE '%1901%' THEN 1901
                WHEN b.qr_code LIKE '%1902%' THEN 1902
                WHEN b.qr_code LIKE '%1903%' THEN 1903
                ELSE NULL
            END as supply_code
        FROM batch b
        WHERE b.id <= 50
    LOOP
        -- Generar insumos para este lote
        IF batch_record.supply_code IS NOT NULL THEN
            FOR i IN 1..batch_record.amount LOOP
                qr_prefix := 'SUPPLY_' || batch_record.batch_id || '_' || LPAD(i::TEXT, 4, '0');
                
                INSERT INTO medical_supply (
                    code, 
                    batch_id, 
                    qr_code, 
                    status, 
                    location_type, 
                    location_id
                ) VALUES (
                    batch_record.supply_code,
                    batch_record.batch_id,
                    qr_prefix,
                    'disponible',
                    'store',
                    batch_record.store_id
                ) ON CONFLICT (qr_code) DO NOTHING;
            END LOOP;
        END IF;
    END LOOP;
END $$;

-- ============================================================================
-- USUARIOS (NO MODIFICAR - MANTENER TAL CUAL)
-- ============================================================================
INSERT INTO "user" (rut, name, email, password, role, medical_center_id, pavilion_id, specialty_id, is_active, created_at, updated_at) VALUES 
('12345678-9', 'Administrador del Sistema', 'admin@meditrack.com', '$2a$10$NA3QLOvkwhpcs.X4KxjONObslo1LreYA6qAzdQcqxRrD4ktjBrpmO', 'admin', 1, NULL, NULL, true, EXTRACT(EPOCH FROM NOW()), EXTRACT(EPOCH FROM NOW())),
('87654321-0', 'Usuario Pabellón', 'pabellon@meditrack.com', '$2a$10$NA3QLOvkwhpcs.X4KxjONObslo1LreYA6qAzdQcqxRrD4ktjBrpmO', 'pabellón', 1, NULL, NULL, true, EXTRACT(EPOCH FROM NOW()), EXTRACT(EPOCH FROM NOW())),
('11111111-1', 'Encargado Bodega', 'bodegacentral@meditrack.com', '$2a$10$NA3QLOvkwhpcs.X4KxjONObslo1LreYA6qAzdQcqxRrD4ktjBrpmO', 'encargado de bodega', 1, NULL, NULL, true, EXTRACT(EPOCH FROM NOW()), EXTRACT(EPOCH FROM NOW())),
('22222222-2', 'María González', 'enfermera@meditrack.com', '$2a$10$NA3QLOvkwhpcs.X4KxjONObslo1LreYA6qAzdQcqxRrD4ktjBrpmO', 'enfermera', 1, NULL, NULL, true, EXTRACT(EPOCH FROM NOW()), EXTRACT(EPOCH FROM NOW())),
('33333333-3', 'Dr. Carlos Pérez', 'doctor@meditrack.com', '$2a$10$NA3QLOvkwhpcs.X4KxjONObslo1LreYA6qAzdQcqxRrD4ktjBrpmO', 'doctor', 1, NULL, 1, true, EXTRACT(EPOCH FROM NOW()), EXTRACT(EPOCH FROM NOW())),
('SYSTEM-INIT', 'Sistema de Inicialización', 'system@meditrack.com', '$2a$10$NA3QLOvkwhpcs.X4KxjONObslo1LreYA6qAzdQcqxRrD4ktjBrpmO', 'admin', 1, NULL, NULL, true, EXTRACT(EPOCH FROM NOW()), EXTRACT(EPOCH FROM NOW())),
('44444444-4', 'Pavedad', 'pavedad@meditrack.com', '$2a$10$NA3QLOvkwhpcs.X4KxjONObslo1LreYA6qAzdQcqxRrD4ktjBrpmO', 'pavedad', 1, NULL, NULL, true, EXTRACT(EPOCH FROM NOW()), EXTRACT(EPOCH FROM NOW())),
('12121212-1', 'Usuario Consignación', 'bodegaconsignacion@meditrack.com', '$2a$10$NA3QLOvkwhpcs.X4KxjONObslo1LreYA6qAzdQcqxRrD4ktjBrpmO', 'encargado de bodega', 1, NULL, NULL, true, EXTRACT(EPOCH FROM NOW()), EXTRACT(EPOCH FROM NOW())),
('55555555-5', 'Dr. Ana Martínez', 'ana.martinez@meditrack.com', '$2a$10$NA3QLOvkwhpcs.X4KxjONObslo1LreYA6qAzdQcqxRrD4ktjBrpmO', 'doctor', 1, NULL, 2, true, EXTRACT(EPOCH FROM NOW()), EXTRACT(EPOCH FROM NOW())),
('66666666-6', 'Dr. Roberto Silva', 'roberto.silva@meditrack.com', '$2a$10$NA3QLOvkwhpcs.X4KxjONObslo1LreYA6qAzdQcqxRrD4ktjBrpmO', 'doctor', 1, NULL, 4, true, EXTRACT(EPOCH FROM NOW()), EXTRACT(EPOCH FROM NOW())),
('77777777-7', 'Dra. Laura Torres', 'laura.torres@meditrack.com', '$2a$10$NA3QLOvkwhpcs.X4KxjONObslo1LreYA6qAzdQcqxRrD4ktjBrpmO', 'doctor', 1, NULL, 8, true, EXTRACT(EPOCH FROM NOW()), EXTRACT(EPOCH FROM NOW())),
('88888888-8', 'Dr. Pedro Ramírez', 'pedro.ramirez@meditrack.com', '$2a$10$NA3QLOvkwhpcs.X4KxjONObslo1LreYA6qAzdQcqxRrD4ktjBrpmO', 'doctor', 1, NULL, 3, true, EXTRACT(EPOCH FROM NOW()), EXTRACT(EPOCH FROM NOW())),
('99999999-9', 'Dra. Carmen López', 'carmen.lopez@meditrack.com', '$2a$10$NA3QLOvkwhpcs.X4KxjONObslo1LreYA6qAzdQcqxRrD4ktjBrpmO', 'doctor', 1, NULL, 5, true, EXTRACT(EPOCH FROM NOW()), EXTRACT(EPOCH FROM NOW())),
('10101010-1', 'Dr. Miguel Ángel Rojas', 'miguel.rojas@meditrack.com', '$2a$10$NA3QLOvkwhpcs.X4KxjONObslo1LreYA6qAzdQcqxRrD4ktjBrpmO', 'doctor', 1, NULL, 7, true, EXTRACT(EPOCH FROM NOW()), EXTRACT(EPOCH FROM NOW())),
('20202020-2', 'Dra. Patricia Vega', 'patricia.vega@meditrack.com', '$2a$10$NA3QLOvkwhpcs.X4KxjONObslo1LreYA6qAzdQcqxRrD4ktjBrpmO', 'doctor', 2, NULL, 10, true, EXTRACT(EPOCH FROM NOW()), EXTRACT(EPOCH FROM NOW())),
('30303030-3', 'Dr. Fernando Castro', 'fernando.castro@meditrack.com', '$2a$10$NA3QLOvkwhpcs.X4KxjONObslo1LreYA6qAzdQcqxRrD4ktjBrpmO', 'doctor', 2, NULL, 6, true, EXTRACT(EPOCH FROM NOW()), EXTRACT(EPOCH FROM NOW())),
('PABELLON1', 'Pabellón 1', 'pabellon1@meditrack.com', '$2a$10$NA3QLOvkwhpcs.X4KxjONObslo1LreYA6qAzdQcqxRrD4ktjBrpmO', 'pabellón', 1, 1, NULL, true, EXTRACT(EPOCH FROM NOW()), EXTRACT(EPOCH FROM NOW()))
ON CONFLICT (rut) DO NOTHING;

-- ============================================================================
-- TRANSFERENCIAS A PABELLONES E INVENTARIO
-- NOTA IMPORTANTE: Para pruebas completas, el Pabellón 1 tiene 70+ items distribuidos en:
--   - Stock 0 (rojo): 2 lotes consumidos totalmente
--   - Stock bajo 1-4 (naranja): 4 lotes con cantidades variables
--   - Stock 5-9 (amarillo): 3 lotes con cantidades medias
--   - Stock 10+ (verde): 3 lotes con cantidades altas
-- Esto permite probar todos los casos de visualización de stock y la paginación
-- ============================================================================

-- ============================================================================
-- TRANSFERENCIAS SEED (DATOS ESTÁTICOS DE PRUEBA)
-- ============================================================================
-- Conjunto simplificado de transferencias con datos estáticos para pruebas
-- Se usa el prefijo TRANS-SEED- para evitar colisiones con códigos de producción

-- NOTA: Los medical_supply_id se obtienen de los QR codes. Los QR tienen la forma SUPPLY_X_YYYY
-- donde X es el batch_id y YYYY es el número secuencial con padding de 4 dígitos

-- Transferencias completadas antiguas (más de 3 días) - Pabellón 1
INSERT INTO supply_transfer (
    transfer_code, qr_code, medical_supply_id, origin_type, origin_id,
    destination_type, destination_id, sent_by, sent_by_name,
    picked_up_by, picked_up_by_name, picked_up_date,
    received_by, received_by_name, status, transfer_reason,
    send_date, receive_date, notes
) 
SELECT 
    'TRANS-SEED-OLD-0001', 'SUPPLY_1_0001', ms.id, 'store', 1, 'pavilion', 1,
    '11111111-1', 'Encargado Bodega', '22222222-2', 'María González',
    NOW() - INTERVAL '5 days', '22222222-2', 'María González', 'recibido',
    'Transferencia antigua para pruebas', NOW() - INTERVAL '5 days',
    NOW() - INTERVAL '5 days' + INTERVAL '1 hour', 'Insumo transferido hace 5 días'
FROM medical_supply ms WHERE ms.qr_code = 'SUPPLY_1_0001'
UNION ALL
SELECT 
    'TRANS-SEED-OLD-0002', 'SUPPLY_1_0002', ms.id, 'store', 1, 'pavilion', 1,
    '11111111-1', 'Encargado Bodega', '22222222-2', 'María González',
    NOW() - INTERVAL '5 days', '22222222-2', 'María González', 'recibido',
    'Transferencia antigua para pruebas', NOW() - INTERVAL '5 days',
    NOW() - INTERVAL '5 days' + INTERVAL '1 hour', 'Insumo transferido hace 5 días'
FROM medical_supply ms WHERE ms.qr_code = 'SUPPLY_1_0002'
UNION ALL
SELECT 
    'TRANS-SEED-OLD-0003', 'SUPPLY_2_0001', ms.id, 'store', 1, 'pavilion', 1,
    '11111111-1', 'Encargado Bodega', '22222222-2', 'María González',
    NOW() - INTERVAL '4 days', '22222222-2', 'María González', 'recibido',
    'Transferencia antigua para pruebas', NOW() - INTERVAL '4 days',
    NOW() - INTERVAL '4 days' + INTERVAL '1 hour', 'Insumo transferido hace 4 días'
FROM medical_supply ms WHERE ms.qr_code = 'SUPPLY_2_0001'
UNION ALL
SELECT 
    'TRANS-SEED-OLD-0004', 'SUPPLY_2_0002', ms.id, 'store', 1, 'pavilion', 2,
    '11111111-1', 'Encargado Bodega', '22222222-2', 'María González',
    NOW() - INTERVAL '3 days', '22222222-2', 'María González', 'recibido',
    'Transferencia antigua para pruebas', NOW() - INTERVAL '3 days',
    NOW() - INTERVAL '3 days' + INTERVAL '1 hour', 'Insumo transferido hace 3 días'
FROM medical_supply ms WHERE ms.qr_code = 'SUPPLY_2_0002'
ON CONFLICT (transfer_code) DO NOTHING;

-- Transferencias recientes (menos de 8 horas) - Válidas
INSERT INTO supply_transfer (
    transfer_code, qr_code, medical_supply_id, origin_type, origin_id,
    destination_type, destination_id, sent_by, sent_by_name,
    picked_up_by, picked_up_by_name, picked_up_date,
    received_by, received_by_name, status, transfer_reason,
    send_date, receive_date, notes
) 
SELECT 
    'TRANS-SEED-REC-0001', 'SUPPLY_3_0001', ms.id, 'store', 1, 'pavilion', 1,
    '11111111-1', 'Encargado Bodega', '22222222-2', 'María González',
    NOW() - INTERVAL '2 hours', '22222222-2', 'María González', 'recibido',
    'Transferencia reciente', NOW() - INTERVAL '2 hours',
    NOW() - INTERVAL '1 hour', 'Insumo transferido recientemente'
FROM medical_supply ms WHERE ms.qr_code = 'SUPPLY_3_0001'
UNION ALL
SELECT 
    'TRANS-SEED-REC-0002', 'SUPPLY_3_0002', ms.id, 'store', 1, 'pavilion', 2,
    '11111111-1', 'Encargado Bodega', '22222222-2', 'María González',
    NOW() - INTERVAL '4 hours', '22222222-2', 'María González', 'recibido',
    'Transferencia reciente', NOW() - INTERVAL '4 hours',
    NOW() - INTERVAL '3 hours', 'Insumo transferido recientemente'
FROM medical_supply ms WHERE ms.qr_code = 'SUPPLY_3_0002'
UNION ALL
SELECT 
    'TRANS-SEED-REC-0003', 'SUPPLY_4_0001', ms.id, 'store', 1, 'pavilion', 3,
    '11111111-1', 'Encargado Bodega', '22222222-2', 'María González',
    NOW() - INTERVAL '6 hours', '22222222-2', 'María González', 'recibido',
    'Transferencia reciente', NOW() - INTERVAL '6 hours',
    NOW() - INTERVAL '5 hours', 'Insumo transferido recientemente'
FROM medical_supply ms WHERE ms.qr_code = 'SUPPLY_4_0001'
ON CONFLICT (transfer_code) DO NOTHING;

-- Transferencias pendientes
INSERT INTO supply_transfer (
    transfer_code, qr_code, medical_supply_id, origin_type, origin_id,
    destination_type, destination_id, sent_by, sent_by_name,
    status, transfer_reason, send_date, notes
) 
SELECT 
    'TRANS-SEED-PEND-0001', 'SUPPLY_5_0001', ms.id, 'store', 1, 'pavilion', 1,
    '11111111-1', 'Encargado Bodega', 'pendiente',
    'Transferencia pendiente de recogida', NOW() - INTERVAL '1 hour',
    'Esperando recogida'
FROM medical_supply ms WHERE ms.qr_code = 'SUPPLY_5_0001'
UNION ALL
SELECT 
    'TRANS-SEED-PEND-0002', 'SUPPLY_5_0002', ms.id, 'store', 1, 'pavilion', 2,
    '11111111-1', 'Encargado Bodega', 'pendiente',
    'Transferencia pendiente de recogida', NOW() - INTERVAL '30 minutes',
    'Esperando recogida'
FROM medical_supply ms WHERE ms.qr_code = 'SUPPLY_5_0002'
ON CONFLICT (transfer_code) DO NOTHING;

-- Transferencias en tránsito
INSERT INTO supply_transfer (
    transfer_code, qr_code, medical_supply_id, origin_type, origin_id,
    destination_type, destination_id, sent_by, sent_by_name,
    picked_up_by, picked_up_by_name, picked_up_date,
    status, transfer_reason, send_date, notes
) 
SELECT 
    'TRANS-SEED-TRAN-0001', 'SUPPLY_6_0001', ms.id, 'store', 1, 'pavilion', 1,
    '11111111-1', 'Encargado Bodega', '22222222-2', 'María González',
    NOW() - INTERVAL '45 minutes', 'en_transito',
    'Transferencia en camino', NOW() - INTERVAL '1 hour',
    'En camino al pabellón'
FROM medical_supply ms WHERE ms.qr_code = 'SUPPLY_6_0001'
UNION ALL
SELECT 
    'TRANS-SEED-TRAN-0002', 'SUPPLY_6_0002', ms.id, 'store', 1, 'pavilion', 3,
    '11111111-1', 'Encargado Bodega', '22222222-2', 'María González',
    NOW() - INTERVAL '20 minutes', 'en_transito',
    'Transferencia en camino', NOW() - INTERVAL '30 minutes',
    'En camino al pabellón'
FROM medical_supply ms WHERE ms.qr_code = 'SUPPLY_6_0002'
ON CONFLICT (transfer_code) DO NOTHING;

-- ============================================================================
-- ACTUALIZAR ESTADOS DE INSUMOS TRANSFERIDOS A PABELLONES
-- IMPORTANTE: Se deshabilita el trigger para poder establecer updated_at manualmente
-- ============================================================================

-- Deshabilitar el trigger que actualiza automáticamente updated_at
ALTER TABLE medical_supply DISABLE TRIGGER trg_update_medical_supply_updated_at;

-- Actualizar estados de insumos transferidos
UPDATE medical_supply ms
SET 
    status = 'recepcionado',
    location_type = 'pavilion',
    location_id = st.destination_id,
    in_transit = false,
    transfer_date = st.receive_date,
    transferred_by = st.received_by,
    updated_at = st.receive_date  -- CRÍTICO: Usar la fecha de recepción de la transferencia
FROM supply_transfer st
WHERE ms.id = st.medical_supply_id
AND st.status = 'recibido';

-- Rehabilitar el trigger
ALTER TABLE medical_supply ENABLE TRIGGER trg_update_medical_supply_updated_at;

-- Crear historial de transferencias a pabellones
INSERT INTO supply_history (
    date_time,
    status,
    destination_type,
    destination_id,
    medical_supply_id,
    user_rut,
    notes,
    location,
    origin_type,
    origin_id,
    confirmed_by,
    confirmation_date
)
SELECT 
    st.send_date,
    'en_camino_a_pabellon',
    st.destination_type,
    st.destination_id,
    st.medical_supply_id,
    st.sent_by,
    'Insumo enviado a pabellón: ' || st.notes,
    'En tránsito desde bodega a pabellón ID: ' || st.destination_id,
    st.origin_type,
    st.origin_id,
    NULL,
    NULL
FROM supply_transfer st
WHERE NOT EXISTS (
    SELECT 1 FROM supply_history sh 
    WHERE sh.medical_supply_id = st.medical_supply_id 
    AND sh.status = 'en_camino_a_pabellon'
);

INSERT INTO supply_history (
    date_time,
    status,
    destination_type,
    destination_id,
    medical_supply_id,
    user_rut,
    notes,
    location,
    origin_type,
    origin_id,
    confirmed_by,
    confirmation_date
)
SELECT 
    st.receive_date,
    'recepcionado',
    st.destination_type,
    st.destination_id,
    st.medical_supply_id,
    st.received_by,
    'Insumo recepcionado en pabellón. ' || 
    CASE 
        WHEN EXTRACT(EPOCH FROM (NOW() - st.receive_date)) / 3600 > 8 
        THEN 'CRÍTICO: Más de 8 horas sin consumir - RETORNO REQUERIDO'
        ELSE 'Recepción normal'
    END,
    p.name || ' (ID: ' || p.id || ')',
    st.origin_type,
    st.origin_id,
    st.received_by,
    st.receive_date
FROM supply_transfer st
JOIN pavilion p ON p.id = st.destination_id
WHERE st.status = 'recibido'
AND NOT EXISTS (
    SELECT 1 FROM supply_history sh 
    WHERE sh.medical_supply_id = st.medical_supply_id 
    AND sh.status = 'recepcionado'
);

-- Crear/actualizar resumen de inventario de pabellones
INSERT INTO pavilion_inventory_summary (
    pavilion_id,
    batch_id,
    supply_code,
    total_received,
    current_available,
    total_consumed,
    total_returned,
    last_received_date,
    created_at,
    updated_at
)
SELECT 
    destination_id AS pavilion_id,
    batch_id,
    supply_code,
    total_received,
    total_received AS current_available,
    0 AS total_consumed,
    0 AS total_returned,
    last_received_date,
    NOW() AS created_at,
    NOW() AS updated_at
FROM (
    SELECT 
        st.destination_id,
        ms.batch_id,
        ms.code AS supply_code,
        COUNT(*) AS total_received,
        MAX(st.receive_date) AS last_received_date
    FROM supply_transfer st
    JOIN medical_supply ms ON ms.id = st.medical_supply_id
    WHERE st.status = 'recibido'
    GROUP BY st.destination_id, ms.batch_id, ms.code
) summary
ON CONFLICT (pavilion_id, batch_id) 
DO UPDATE SET
    total_received = pavilion_inventory_summary.total_received + EXCLUDED.total_received,
    current_available = pavilion_inventory_summary.current_available + EXCLUDED.current_available,
    last_received_date = GREATEST(pavilion_inventory_summary.last_received_date, EXCLUDED.last_received_date),
    updated_at = NOW();

-- Ya no es necesario simular transfers para consumo - el script queda más simple


-- ============================================================================
-- RESUMEN DE INVENTARIO POR BODEGA (CREAR ANTES DE ACTUALIZAR)
-- ============================================================================
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

-- Actualizar inventario de bodegas (decrementar por transferencias)
UPDATE store_inventory_summary sis
SET 
    current_in_store = sis.current_in_store - transferred.count,
    total_transferred_out = sis.total_transferred_out + transferred.count,
    last_transfer_out_date = transferred.last_transfer,
    updated_at = NOW()
FROM (
    SELECT 
        b.store_id,
        ms.batch_id,
        COUNT(*) as count,
        MAX(st.send_date) as last_transfer
    FROM supply_transfer st
    JOIN medical_supply ms ON ms.id = st.medical_supply_id
    JOIN batch b ON b.id = ms.batch_id
    WHERE st.status = 'recibido'
    GROUP BY b.store_id, ms.batch_id
) transferred
WHERE sis.store_id = transferred.store_id
AND sis.batch_id = transferred.batch_id;

-- ============================================================================
-- INSUMOS TÍPICOS POR CIRUGÍA - COMPLETO Y COHERENTE
-- ============================================================================
INSERT INTO surgery_typical_supply (surgery_id, supply_code, typical_quantity, is_required, notes) VALUES
-- Cirugías Plásticas (1, 2)
(1, 1001, 50, TRUE, 'Guantes quirúrgicos estériles'),
(1, 1002, 10, TRUE, 'Mascarillas N95 para protección'),
(1, 1005, 30, TRUE, 'Gasas estériles'),
(1, 1007, 5, TRUE, 'Batas quirúrgicas'),
(1, 1008, 3, TRUE, 'Campos quirúrgicos'),
(1, 1101, 10, TRUE, 'Suturas Vicryl para cierre'),
(1, 1102, 5, FALSE, 'Suturas Nylon para piel'),
(2, 1001, 40, TRUE, 'Guantes quirúrgicos'),
(2, 1005, 25, TRUE, 'Gasas estériles'),
(2, 1101, 8, TRUE, 'Suturas Vicryl'),

-- Cirugía General (3, 4, 26)
(3, 1001, 30, TRUE, 'Guantes quirúrgicos'),
(3, 1005, 20, TRUE, 'Gasas estériles'),
(3, 1101, 5, TRUE, 'Suturas Vicryl'),
(3, 1102, 3, TRUE, 'Suturas Nylon'),
(4, 1001, 60, TRUE, 'Guantes quirúrgicos'),
(4, 1005, 40, TRUE, 'Gasas estériles'),
(4, 1007, 6, TRUE, 'Batas quirúrgicas'),
(4, 1101, 15, TRUE, 'Suturas Vicryl'),
(4, 1103, 10, TRUE, 'Suturas Prolene'),
(4, 1504, 4, TRUE, 'Tornillos ortopédicos'),
(26, 1001, 45, TRUE, 'Guantes quirúrgicos'),
(26, 1005, 30, TRUE, 'Gasas estériles'),
(26, 1101, 8, TRUE, 'Suturas Vicryl'),
(26, 1401, 2, TRUE, 'Bisturí eléctrico'),

-- Traumatología (5, 74, 75, 133, 137)
(5, 1001, 50, TRUE, 'Guantes quirúrgicos'),
(5, 1003, 5, FALSE, 'Jeringas para anestesia local'),
(5, 1004, 10, TRUE, 'Agujas estériles'),
(5, 1005, 35, TRUE, 'Gasas para limpieza'),
(5, 1101, 10, TRUE, 'Suturas Vicryl'),
(5, 1504, 6, TRUE, 'Tornillos ortopédicos'),
(5, 1505, 2, TRUE, 'Placas ortopédicas'),
(74, 1001, 45, TRUE, 'Guantes quirúrgicos'),
(74, 1005, 30, TRUE, 'Gasas estériles'),
(74, 1504, 8, TRUE, 'Tornillos ortopédicos'),
(74, 1505, 3, TRUE, 'Placas ortopédicas'),
(75, 1001, 40, TRUE, 'Guantes quirúrgicos'),
(75, 1005, 25, TRUE, 'Gasas estériles'),
(75, 1504, 6, TRUE, 'Tornillos ortopédicos'),
(75, 1505, 2, TRUE, 'Placas ortopédicas'),
(133, 1001, 50, TRUE, 'Guantes quirúrgicos'),
(133, 1005, 35, TRUE, 'Gasas estériles'),
(133, 1504, 10, TRUE, 'Tornillos ortopédicos'),
(133, 1505, 3, TRUE, 'Placas ortopédicas'),
(137, 1001, 55, TRUE, 'Guantes quirúrgicos'),
(137, 1005, 40, TRUE, 'Gasas estériles'),
(137, 1504, 12, TRUE, 'Tornillos ortopédicos'),
(137, 1505, 4, TRUE, 'Placas ortopédicas'),

-- Colecistectomía (34, 35)
(34, 1001, 40, TRUE, 'Guantes quirúrgicos'),
(34, 1002, 5, TRUE, 'Mascarillas'),
(34, 1005, 30, TRUE, 'Gasas estériles'),
(34, 1007, 4, TRUE, 'Batas quirúrgicas'),
(34, 1101, 8, TRUE, 'Suturas Vicryl'),
(34, 1303, 2, FALSE, 'Drenaje Blake'),
(35, 1001, 35, TRUE, 'Guantes quirúrgicos'),
(35, 1005, 25, TRUE, 'Gasas estériles'),
(35, 1007, 3, TRUE, 'Batas quirúrgicas'),
(35, 1101, 6, TRUE, 'Suturas Vicryl'),

-- Oftalmología (64, 68)
(64, 1001, 20, TRUE, 'Guantes quirúrgicos'),
(64, 1005, 15, TRUE, 'Gasas estériles'),
(64, 1603, 4, TRUE, 'Suturas Nylon 10-0'),
(68, 1001, 15, TRUE, 'Guantes quirúrgicos'),
(68, 1005, 10, TRUE, 'Gasas estériles'),
(68, 1601, 1, TRUE, 'Lente intraocular'),
(68, 1602, 2, TRUE, 'Viscoelástico oftálmico'),

-- Urología (11, 12)
(11, 1001, 50, TRUE, 'Guantes quirúrgicos'),
(11, 1005, 35, TRUE, 'Gasas estériles'),
(11, 1007, 5, TRUE, 'Batas quirúrgicas'),
(11, 1302, 2, TRUE, 'Catéter Foley'),
(11, 1701, 3, TRUE, 'Sonda vesical'),
(12, 1001, 40, TRUE, 'Guantes quirúrgicos'),
(12, 1005, 30, TRUE, 'Gasas estériles'),
(12, 1302, 2, TRUE, 'Catéter Foley'),
(12, 1701, 2, TRUE, 'Sonda vesical'),

-- Ginecología (40, 92, 93)
(40, 1001, 25, TRUE, 'Guantes quirúrgicos'),
(40, 1005, 15, TRUE, 'Gasas estériles'),
(40, 1801, 2, TRUE, 'Espéculo vaginal'),
(40, 1802, 1, TRUE, 'Pinzas de biopsia cervical'),
(92, 1001, 45, TRUE, 'Guantes quirúrgicos'),
(92, 1005, 35, TRUE, 'Gasas estériles'),
(92, 1007, 4, TRUE, 'Batas quirúrgicas'),
(92, 1101, 10, TRUE, 'Suturas Vicryl'),
(93, 1001, 60, TRUE, 'Guantes quirúrgicos'),
(93, 1005, 45, TRUE, 'Gasas estériles'),
(93, 1007, 6, TRUE, 'Batas quirúrgicas'),
(93, 1101, 15, TRUE, 'Suturas Vicryl'),

-- Neurocirugía (31, 43, 213)
(31, 1001, 70, TRUE, 'Guantes quirúrgicos'),
(31, 1002, 8, TRUE, 'Mascarillas N95'),
(31, 1005, 50, TRUE, 'Gasas estériles'),
(31, 1007, 8, TRUE, 'Batas quirúrgicas'),
(31, 1903, 3, TRUE, 'Hemostáticos neurológicos'),
(43, 1001, 80, TRUE, 'Guantes quirúrgicos'),
(43, 1005, 60, TRUE, 'Gasas estériles'),
(43, 1007, 10, TRUE, 'Batas quirúrgicas'),
(43, 1901, 1, TRUE, 'Craneotomo desechable'),
(43, 1903, 4, TRUE, 'Hemostáticos neurológicos'),
(213, 1001, 50, TRUE, 'Guantes quirúrgicos'),
(213, 1005, 40, TRUE, 'Gasas estériles'),
(213, 1902, 1, TRUE, 'Válvula derivación ventriculoperitoneal'),

-- Cardiología (218)
(218, 1001, 70, TRUE, 'Guantes quirúrgicos'),
(218, 1002, 10, TRUE, 'Mascarillas N95'),
(218, 1005, 60, TRUE, 'Gasas estériles'),
(218, 1007, 8, TRUE, 'Batas quirúrgicas'),
(218, 1301, 2, TRUE, 'Catéter venoso central'),
(218, 1401, 3, TRUE, 'Bisturí eléctrico'),
(218, 1402, 5, TRUE, 'Compresas hemostáticas'),
(218, 1403, 10, TRUE, 'Ligaduras vasculares'),

-- Cirugías Plásticas adicionales (115, 116)
(115, 1001, 55, TRUE, 'Guantes quirúrgicos'),
(115, 1002, 8, TRUE, 'Mascarillas N95'),
(115, 1005, 40, TRUE, 'Gasas estériles'),
(115, 1007, 6, TRUE, 'Batas quirúrgicas'),
(115, 1101, 12, TRUE, 'Suturas Vicryl'),
(115, 1102, 8, TRUE, 'Suturas Nylon'),
(116, 1001, 50, TRUE, 'Guantes quirúrgicos'),
(116, 1005, 35, TRUE, 'Gasas estériles'),
(116, 1007, 5, TRUE, 'Batas quirúrgicas'),
(116, 1101, 10, TRUE, 'Suturas Vicryl')
ON CONFLICT (surgery_id, supply_code) DO NOTHING;

-- ============================================================================
-- HISTORIAL INICIAL DE LOTES
-- ============================================================================
INSERT INTO batch_history (date_time, change_details, previous_values, new_values, user_name, batch_id, user_rut, batch_number) VALUES
('2025-01-15 09:00:00', 'Lote creado - Bodega Central', NULL, '{"expiration_date": "2026-12-31", "amount": 150, "supplier": "Proveedor Uno"}', 'Administrador del Sistema', 1, '12345678-9', 1),
('2025-01-15 09:05:00', 'Lote creado - Bodega Central', NULL, '{"expiration_date": "2026-10-15", "amount": 200, "supplier": "Proveedor Uno"}', 'Encargado Bodega', 2, '11111111-1', 2),
('2025-01-16 10:30:00', 'Lote creado - Bodega Central', NULL, '{"expiration_date": "2027-03-20", "amount": 120, "supplier": "Proveedor Dos"}', 'Encargado Bodega', 3, '11111111-1', 3),
('2025-01-16 11:00:00', 'Lote creado - Bodega Central', NULL, '{"expiration_date": "2026-08-30", "amount": 150, "supplier": "Proveedor Dos"}', 'Encargado Bodega', 4, '11111111-1', 4),
('2025-01-17 08:15:00', 'Lote creado - Bodega Consignación', NULL, '{"expiration_date": "2026-11-30", "amount": 100, "supplier": "Proveedor Consignación A"}', 'Usuario Consignación', 34, '12121212-1', 34),
('2025-01-17 08:20:00', 'Lote creado - Bodega Consignación', NULL, '{"expiration_date": "2027-02-28", "amount": 150, "supplier": "Proveedor Consignación A"}', 'Usuario Consignación', 35, '12121212-1', 35),
('2025-01-18 09:45:00', 'Lote creado - Bodega Central', NULL, '{"expiration_date": "2026-09-30", "amount": 80, "supplier": "Proveedor Uno"}', 'Encargado Bodega', 44, '11111111-1', 44),
('2025-01-18 10:00:00', 'Lote creado - Bodega Central', NULL, '{"expiration_date": "2026-10-15", "amount": 100, "supplier": "Proveedor Uno"}', 'Encargado Bodega', 45, '11111111-1', 45)
ON CONFLICT DO NOTHING;

-- ============================================================================
-- SOLICITUDES DE INSUMOS - DIVERSOS ESTADOS
-- ============================================================================
-- Solicitud 1: Pendiente Pavedad (Estado inicial)
INSERT INTO supply_request (
    request_number, pavilion_id, requested_by, requested_by_name,
    request_date, surgery_datetime, status, notes, medical_center_id,
    surgeon_id, surgeon_name, surgery_id, specialty_id
) VALUES (
    'SOL-20250205090000',
    1,
    '33333333-3',
    'Dr. Carlos Pérez',
    NOW() - INTERVAL '2 hours',
    NOW() + INTERVAL '47 hours',
    'pendiente_pavedad',
    'Cirugía de reducción abierta programada para pasado mañana. Requiere material ortopédico.',
    1,
    '33333333-3',
    'Dr. Carlos Pérez',
    5,
    (SELECT id FROM medical_specialty WHERE code = 'TRAUMA')
) ON CONFLICT DO NOTHING;

INSERT INTO supply_request_item (supply_request_id, supply_code, supply_name, quantity_requested, is_pediatric) VALUES 
((SELECT id FROM supply_request WHERE request_number = 'SOL-20250205090000'), 1001, 'Guantes Quirúrgicos Estériles', 50, FALSE),
((SELECT id FROM supply_request WHERE request_number = 'SOL-20250205090000'), 1005, 'Gasas Estériles 10x10cm', 35, FALSE),
((SELECT id FROM supply_request WHERE request_number = 'SOL-20250205090000'), 1504, 'Tornillos Ortopédicos 3.5mm', 6, FALSE),
((SELECT id FROM supply_request WHERE request_number = 'SOL-20250205090000'), 1505, 'Placas Ortopédicas Rectas', 2, FALSE)
ON CONFLICT DO NOTHING;

-- Solicitud 2: Asignado a bodega
INSERT INTO supply_request (
    request_number, pavilion_id, requested_by, requested_by_name,
    request_date, surgery_datetime, status, notes, medical_center_id,
    surgeon_id, surgeon_name, surgery_id, specialty_id,
    assigned_to, assigned_to_name, assigned_date, assigned_by_pavedad, assigned_by_pavedad_name
) VALUES (
    'SOL-20250204140000',
    2,
    '55555555-5',
    'Dr. Ana Martínez',
    NOW() - INTERVAL '1 day',
    NOW() + INTERVAL '71 hours',
    'asignado_bodega',
    'Procedimiento cardiovascular. Pavedad asignó a Bodega Central.',
    1,
    '55555555-5',
    'Dr. Ana Martínez',
    218,
    (SELECT id FROM medical_specialty WHERE code = 'CARDIOL'),
    '11111111-1',
    'Encargado Bodega',
    NOW() - INTERVAL '20 hours',
    '44444444-4',
    'Pavedad'
) ON CONFLICT DO NOTHING;

INSERT INTO supply_request_item (supply_request_id, supply_code, supply_name, quantity_requested, is_pediatric) VALUES 
((SELECT id FROM supply_request WHERE request_number = 'SOL-20250204140000'), 1001, 'Guantes Quirúrgicos Estériles', 70, FALSE),
((SELECT id FROM supply_request WHERE request_number = 'SOL-20250204140000'), 1002, 'Mascarillas N95', 10, FALSE),
((SELECT id FROM supply_request WHERE request_number = 'SOL-20250204140000'), 1005, 'Gasas Estériles 10x10cm', 60, FALSE),
((SELECT id FROM supply_request WHERE request_number = 'SOL-20250204140000'), 1301, 'Catéter Venoso Central', 2, FALSE),
((SELECT id FROM supply_request WHERE request_number = 'SOL-20250204140000'), 1402, 'Compresas Hemostáticas', 5, FALSE)
ON CONFLICT DO NOTHING;

-- Solicitud 3: En proceso (encargado está preparando)
INSERT INTO supply_request (
    request_number, pavilion_id, requested_by, requested_by_name,
    request_date, surgery_datetime, status, notes, medical_center_id,
    surgeon_id, surgeon_name, surgery_id, specialty_id,
    assigned_to, assigned_to_name, assigned_date, assigned_by_pavedad, assigned_by_pavedad_name
) VALUES (
    'SOL-20250203100000',
    3,
    '99999999-9',
    'Dra. Carmen López',
    NOW() - INTERVAL '2 days',
    NOW() + INTERVAL '25 hours',
    'en_proceso',
    'Cirugía oftalmológica mañana. Encargado está preparando el material.',
    1,
    '99999999-9',
    'Dra. Carmen López',
    68,
    (SELECT id FROM medical_specialty WHERE code = 'OFTAL'),
    '11111111-1',
    'Encargado Bodega',
    NOW() - INTERVAL '36 hours',
    '44444444-4',
    'Pavedad'
) ON CONFLICT DO NOTHING;

INSERT INTO supply_request_item (supply_request_id, supply_code, supply_name, quantity_requested, is_pediatric, item_status, reviewed_by, reviewed_by_name, reviewed_at) VALUES 
((SELECT id FROM supply_request WHERE request_number = 'SOL-20250203100000'), 1001, 'Guantes Quirúrgicos Estériles', 15, FALSE, 'aceptado', '11111111-1', 'Encargado Bodega', NOW() - INTERVAL '12 hours'),
((SELECT id FROM supply_request WHERE request_number = 'SOL-20250203100000'), 1005, 'Gasas Estériles 10x10cm', 10, FALSE, 'aceptado', '11111111-1', 'Encargado Bodega', NOW() - INTERVAL '12 hours'),
((SELECT id FROM supply_request WHERE request_number = 'SOL-20250203100000'), 1601, 'Lente Intraocular', 1, FALSE, 'aceptado', '11111111-1', 'Encargado Bodega', NOW() - INTERVAL '12 hours'),
((SELECT id FROM supply_request WHERE request_number = 'SOL-20250203100000'), 1602, 'Viscoelástico Oftálmico', 2, FALSE, 'aceptado', '11111111-1', 'Encargado Bodega', NOW() - INTERVAL '12 hours')
ON CONFLICT DO NOTHING;

-- Solicitud 4: Aprobada (lista para retirar)
INSERT INTO supply_request (
    request_number, pavilion_id, requested_by, requested_by_name,
    request_date, surgery_datetime, status, notes, medical_center_id,
    surgeon_id, surgeon_name, surgery_id, specialty_id,
    assigned_to, assigned_to_name, assigned_date, assigned_by_pavedad, assigned_by_pavedad_name,
    approved_by, approved_by_name, approval_date
) VALUES (
    'SOL-20250202153000',
    4,
    '66666666-6',
    'Dr. Roberto Silva',
    NOW() - INTERVAL '3 days',
    NOW() + INTERVAL '12 hours',
    'aprobado',
    'Neurocirugía urgente mañana temprano. Material listo para retiro.',
    1,
    '66666666-6',
    'Dr. Roberto Silva',
    213,
    (SELECT id FROM medical_specialty WHERE code = 'NEURO'),
    '11111111-1',
    'Encargado Bodega',
    NOW() - INTERVAL '60 hours',
    '44444444-4',
    'Pavedad',
    '11111111-1',
    'Encargado Bodega',
    NOW() - INTERVAL '6 hours'
) ON CONFLICT DO NOTHING;

INSERT INTO supply_request_item (supply_request_id, supply_code, supply_name, quantity_requested, quantity_approved, is_pediatric, item_status, reviewed_by, reviewed_by_name, reviewed_at) VALUES 
((SELECT id FROM supply_request WHERE request_number = 'SOL-20250202153000'), 1001, 'Guantes Quirúrgicos Estériles', 50, 50, FALSE, 'aceptado', '11111111-1', 'Encargado Bodega', NOW() - INTERVAL '8 hours'),
((SELECT id FROM supply_request WHERE request_number = 'SOL-20250202153000'), 1005, 'Gasas Estériles 10x10cm', 40, 40, FALSE, 'aceptado', '11111111-1', 'Encargado Bodega', NOW() - INTERVAL '8 hours'),
((SELECT id FROM supply_request WHERE request_number = 'SOL-20250202153000'), 1902, 'Válvula Derivación Ventriculoperitoneal', 1, 1, FALSE, 'aceptado', '11111111-1', 'Encargado Bodega', NOW() - INTERVAL '8 hours')
ON CONFLICT DO NOTHING;

-- Solicitud 5: Completada (cirugía ya realizada)
INSERT INTO supply_request (
    request_number, pavilion_id, requested_by, requested_by_name,
    request_date, surgery_datetime, status, notes, medical_center_id,
    surgeon_id, surgeon_name, surgery_id, specialty_id,
    assigned_to, assigned_to_name, assigned_date, assigned_by_pavedad, assigned_by_pavedad_name,
    approved_by, approved_by_name, approval_date, completed_date
) VALUES (
    'SOL-20250130083000',
    5,
    '77777777-7',
    'Dra. Laura Torres',
    NOW() - INTERVAL '6 days',
    NOW() - INTERVAL '1 day',
    'completado',
    'Histerectomía realizada exitosamente. Todos los insumos fueron consumidos.',
    1,
    '77777777-7',
    'Dra. Laura Torres',
    92,
    (SELECT id FROM medical_specialty WHERE code = 'GINEC'),
    '11111111-1',
    'Encargado Bodega',
    NOW() - INTERVAL '5 days',
    '44444444-4',
    'Pavedad',
    '11111111-1',
    'Encargado Bodega',
    NOW() - INTERVAL '2 days',
    NOW() - INTERVAL '18 hours'
) ON CONFLICT DO NOTHING;

INSERT INTO supply_request_item (supply_request_id, supply_code, supply_name, quantity_requested, quantity_approved, quantity_delivered, is_pediatric, item_status, reviewed_by, reviewed_by_name, reviewed_at) VALUES 
((SELECT id FROM supply_request WHERE request_number = 'SOL-20250130083000'), 1001, 'Guantes Quirúrgicos Estériles', 45, 45, 45, FALSE, 'aceptado', '11111111-1', 'Encargado Bodega', NOW() - INTERVAL '3 days'),
((SELECT id FROM supply_request WHERE request_number = 'SOL-20250130083000'), 1005, 'Gasas Estériles 10x10cm', 35, 35, 35, FALSE, 'aceptado', '11111111-1', 'Encargado Bodega', NOW() - INTERVAL '3 days'),
((SELECT id FROM supply_request WHERE request_number = 'SOL-20250130083000'), 1007, 'Batas Quirúrgicas Estériles', 4, 4, 4, FALSE, 'aceptado', '11111111-1', 'Encargado Bodega', NOW() - INTERVAL '3 days'),
((SELECT id FROM supply_request WHERE request_number = 'SOL-20250130083000'), 1101, 'Sutura Vicryl 2-0', 10, 10, 10, FALSE, 'aceptado', '11111111-1', 'Encargado Bodega', NOW() - INTERVAL '3 days')
ON CONFLICT DO NOTHING;

-- ============================================================================
-- ASIGNACIONES QR PARA SOLICITUDES (vincular insumos con solicitudes)
-- ============================================================================

-- Asignar QR codes específicos a la solicitud completada SOL-20250130083000
-- Estos insumos se transferirán al pabellón y quedarán pendientes de retorno
DO $$
DECLARE
    v_request_id INT;
    v_item_id INT;
    v_supply_id INT;
    v_supply_qr VARCHAR(255);
    v_counter INT := 0;
BEGIN
    -- Obtener ID de la solicitud
    SELECT id INTO v_request_id 
    FROM supply_request 
    WHERE request_number = 'SOL-20250130083000';
    
    -- Para cada item de la solicitud, asignar insumos específicos
    FOR v_item_id IN 
        SELECT id FROM supply_request_item 
        WHERE supply_request_id = v_request_id 
        AND quantity_approved > 0
        LIMIT 20  -- Limitar para no sobrecargar
    LOOP
        -- Buscar un insumo disponible del código correcto
        SELECT ms.id, ms.qr_code INTO v_supply_id, v_supply_qr
        FROM medical_supply ms
        JOIN supply_request_item sri ON sri.supply_code = ms.code
        WHERE sri.id = v_item_id
        AND ms.status = 'disponible'
        AND ms.location_type = 'store'
        AND NOT EXISTS (
            SELECT 1 FROM supply_request_qr_assignment qa 
            WHERE qa.medical_supply_id = ms.id
        )
        LIMIT 1;
        
        -- Si encontró un insumo, crear la asignación
        IF v_supply_id IS NOT NULL THEN
            INSERT INTO supply_request_qr_assignment (
                supply_request_id,
                supply_request_item_id,
                medical_supply_id,
                qr_code,
                assigned_date,
                assigned_by,
                assigned_by_name,
                delivered_date,
                delivered_by,
                delivered_by_name,
                status
            ) VALUES (
                v_request_id,
                v_item_id,
                v_supply_id,
                v_supply_qr,
                NOW() - INTERVAL '6 days' + INTERVAL '4 hours',
                '11111111-1',
                'Encargado Bodega',
                NOW() - INTERVAL '5 days',
                '22222222-2',
                'María González',
                'delivered'
            ) ON CONFLICT DO NOTHING;
            
            v_counter := v_counter + 1;
        END IF;
    END LOOP;
    
    RAISE NOTICE 'Asignadas % asignaciones QR para solicitud SOL-20250130083000', v_counter;
END $$;

-- Solicitud 6: Rechazada (falta de stock)
INSERT INTO supply_request (
    request_number, pavilion_id, requested_by, requested_by_name,
    request_date, surgery_datetime, status, notes, medical_center_id,
    surgeon_id, surgeon_name, surgery_id, specialty_id,
    assigned_to, assigned_to_name, assigned_date, assigned_by_pavedad, assigned_by_pavedad_name,
    approved_by, approved_by_name, approval_date
) VALUES (
    'SOL-20250131163000',
    6,
    '88888888-8',
    'Dr. Pedro Ramírez',
    NOW() - INTERVAL '5 days',
    NOW() + INTERVAL '96 hours',
    'rechazado',
    'Solicitud rechazada por falta de stock de material específico. Se sugiere reprogramar o buscar proveedor alternativo.',
    1,
    '88888888-8',
    'Dr. Pedro Ramírez',
    36,
    (SELECT id FROM medical_specialty WHERE code = 'CIR_GEN'),
    '11111111-1',
    'Encargado Bodega',
    NOW() - INTERVAL '4 days',
    '44444444-4',
    'Pavedad',
    '11111111-1',
    'Encargado Bodega',
    NOW() - INTERVAL '3 days'
) ON CONFLICT DO NOTHING;

INSERT INTO supply_request_item (supply_request_id, supply_code, supply_name, quantity_requested, is_pediatric, item_status, item_notes, reviewed_by, reviewed_by_name, reviewed_at) VALUES 
((SELECT id FROM supply_request WHERE request_number = 'SOL-20250131163000'), 1001, 'Guantes Quirúrgicos Estériles', 60, FALSE, 'aceptado', NULL, '11111111-1', 'Encargado Bodega', NOW() - INTERVAL '3 days'),
((SELECT id FROM supply_request WHERE request_number = 'SOL-20250131163000'), 1005, 'Gasas Estériles 10x10cm', 40, FALSE, 'aceptado', NULL, '11111111-1', 'Encargado Bodega', NOW() - INTERVAL '3 days'),
((SELECT id FROM supply_request WHERE request_number = 'SOL-20250131163000'), 1401, 'Bisturí Eléctrico Desechable', 3, FALSE, 'rechazado', 'Stock insuficiente. Disponible: 1, Solicitado: 3', '11111111-1', 'Encargado Bodega', NOW() - INTERVAL '3 days')
ON CONFLICT DO NOTHING;

-- ============================================================================
-- HISTORIAL INICIAL DE INSUMOS
-- ============================================================================
INSERT INTO supply_history (
    date_time,
    status,
    destination_type,
    destination_id,
    medical_supply_id,
    user_rut,
    notes,
    location
)
SELECT 
    NOW() - INTERVAL '30 days' AS date_time,
    'disponible' AS status,
    'store' AS destination_type,
    b.store_id AS destination_id,
    ms.id AS medical_supply_id,
    'SYSTEM-INIT' AS user_rut,
    'Registro inicial - insumo ingresado a bodega desde proveedor ' || b.supplier_id AS notes,
    s.name || ' (ID: ' || s.id || ')' AS location
FROM medical_supply ms
JOIN batch b ON ms.batch_id = b.id
JOIN store s ON b.store_id = s.id
WHERE ms.status = 'disponible'
AND NOT EXISTS (
    SELECT 1 FROM supply_history sh 
    WHERE sh.medical_supply_id = ms.id
)
LIMIT 100;

-- (Los proveedores ya fueron insertados antes de los lotes por la FK)

-- ============================================================================
-- AJUSTE DE SECUENCIAS
-- ============================================================================
SELECT setval('surgery_id_seq', (SELECT COALESCE(MAX(id), 0) FROM surgery));
SELECT setval('medical_center_id_seq', (SELECT COALESCE(MAX(id), 0) FROM medical_center));
SELECT setval('pavilion_id_seq', (SELECT COALESCE(MAX(id), 0) FROM pavilion));
SELECT setval('store_id_seq', (SELECT COALESCE(MAX(id), 0) FROM store));
SELECT setval('batch_id_seq', (SELECT COALESCE(MAX(id), 0) FROM batch));
SELECT setval('medical_supply_id_seq', (SELECT COALESCE(MAX(id), 0) FROM medical_supply));
SELECT setval('supply_history_id_seq', (SELECT COALESCE(MAX(id), 0) FROM supply_history));
SELECT setval('medical_specialty_id_seq', (SELECT COALESCE(MAX(id), 0) FROM medical_specialty));
SELECT setval('surgery_typical_supply_id_seq', (SELECT COALESCE(MAX(id), 0) FROM surgery_typical_supply));
SELECT setval('supply_request_id_seq', (SELECT COALESCE(MAX(id), 0) FROM supply_request));
SELECT setval('supply_request_item_id_seq', (SELECT COALESCE(MAX(id), 0) FROM supply_request_item));
SELECT setval('store_inventory_summary_id_seq', (SELECT COALESCE(MAX(id), 0) FROM store_inventory_summary));
SELECT setval('batch_history_id_seq', (SELECT COALESCE(MAX(id), 0) FROM batch_history));

-- ============================================================================
-- REPORTES DE VERIFICACIÓN
-- ============================================================================
SELECT 'RESUMEN DE POBLACIÓN DE BASE DE DATOS' as titulo;

SELECT 'Centros Médicos:' as categoria, COUNT(*) as total FROM medical_center
UNION ALL
SELECT 'Pabellones:', COUNT(*) FROM pavilion
UNION ALL
SELECT 'Bodegas:', COUNT(*) FROM store
UNION ALL
SELECT 'Códigos de Insumo:', COUNT(*) FROM supply_code
UNION ALL
SELECT 'Especialidades Médicas:', COUNT(*) FROM medical_specialty
UNION ALL
SELECT 'Cirugías:', COUNT(*) FROM surgery
UNION ALL
SELECT 'Lotes:', COUNT(*) FROM batch
UNION ALL
SELECT 'Insumos Médicos:', COUNT(*) FROM medical_supply
UNION ALL
SELECT 'Usuarios:', COUNT(*) FROM "user"
UNION ALL
SELECT 'Insumos Típicos por Cirugía:', COUNT(*) FROM surgery_typical_supply
UNION ALL
SELECT 'Solicitudes de Insumo:', COUNT(*) FROM supply_request
UNION ALL
SELECT 'Items en Solicitudes:', COUNT(*) FROM supply_request_item
UNION ALL
SELECT 'Registros en Historial:', COUNT(*) FROM supply_history
UNION ALL
SELECT 'Resúmenes de Inventario:', COUNT(*) FROM store_inventory_summary
UNION ALL
SELECT 'Configuraciones de Proveedor:', COUNT(*) FROM supplier_config;

-- Distribución de insumos por bodega
SELECT 
    '
DISTRIBUCIÓN DE INSUMOS POR BODEGA' as titulo;
    
SELECT 
    s.name as bodega,
    COUNT(DISTINCT b.id) as total_lotes,
    COUNT(ms.id) as total_insumos,
    COUNT(DISTINCT ms.code) as tipos_insumo_diferentes
FROM store s
LEFT JOIN batch b ON s.id = b.store_id
LEFT JOIN medical_supply ms ON b.id = ms.batch_id
GROUP BY s.id, s.name
ORDER BY s.id;

-- Solicitudes por estado
SELECT 
    '
SOLICITUDES POR ESTADO' as titulo;
    
SELECT 
    status,
    COUNT(*) as total_solicitudes,
    COUNT(DISTINCT surgeon_id) as doctores_diferentes
FROM supply_request
GROUP BY status
ORDER BY 
    CASE status
        WHEN 'pendiente_pavedad' THEN 1
        WHEN 'asignado_bodega' THEN 2
        WHEN 'en_proceso' THEN 3
        WHEN 'aprobado' THEN 4
        WHEN 'completado' THEN 5
        WHEN 'rechazado' THEN 6
        ELSE 7
    END;

-- Cirugías con insumos típicos configurados
SELECT 
    '
CIRUGÍAS CON INSUMOS TÍPICOS' as titulo;
    
SELECT 
    COUNT(DISTINCT surgery_id) as cirugias_con_insumos,
    COUNT(*) as total_relaciones_insumo_cirugia,
    ROUND(AVG(typical_quantity), 2) as cantidad_promedio_por_relacion
FROM surgery_typical_supply;

-- ============================================================================
-- REPORTES DE GESTIÓN DE RETORNOS A BODEGA
-- ============================================================================
SELECT '
GESTIÓN DE RETORNOS A BODEGA' as titulo;

-- Insumos pendientes de retorno (>8 horas laborales sin consumir)
SELECT 
    '
Pendientes de Retorno:' as categoria,
    COUNT(*) as total
FROM medical_supply ms
JOIN supply_transfer st ON st.medical_supply_id = ms.id
WHERE ms.status = 'recepcionado'
AND ms.location_type = 'pavilion'
AND st.status = 'recibido'
AND EXTRACT(EPOCH FROM (NOW() - st.receive_date)) / 3600 > 8;

-- Insumos críticos (>8 horas laborales)
SELECT 
    'Críticos (>8 horas laborales):' as categoria,
    COUNT(*) as total
FROM medical_supply ms
JOIN supply_transfer st ON st.medical_supply_id = ms.id
WHERE ms.status = 'recepcionado'
AND ms.location_type = 'pavilion'
AND st.status = 'recibido'
AND EXTRACT(EPOCH FROM (NOW() - st.receive_date)) / 3600 > 8;

-- Detalle de insumos pendientes de retorno
SELECT 
    '
DETALLE DE INSUMOS PENDIENTES DE RETORNO' as titulo;

SELECT 
    p.name as pabellon,
    sc.name as insumo,
    COUNT(*) as cantidad_insumos,
    ROUND(AVG(EXTRACT(EPOCH FROM (NOW() - st.receive_date)) / 3600), 1) as horas_promedio_en_pabellon,
    MIN(st.receive_date) as primera_recepcion,
    MAX(st.receive_date) as ultima_recepcion
FROM medical_supply ms
JOIN supply_transfer st ON st.medical_supply_id = ms.id
JOIN pavilion p ON p.id = ms.location_id
JOIN supply_code sc ON sc.code = ms.code
WHERE ms.status = 'recepcionado'
AND ms.location_type = 'pavilion'
AND st.status = 'recibido'
AND EXTRACT(EPOCH FROM (NOW() - st.receive_date)) / 3600 > 8
GROUP BY p.name, sc.name
ORDER BY horas_promedio_en_pabellon DESC;

-- Insumos por urgencia de retorno
SELECT 
    '
INSUMOS POR URGENCIA DE RETORNO' as titulo;

SELECT 
    nivel_urgencia,
    COUNT(*) as cantidad_insumos,
    COUNT(DISTINCT location_id) as pabellones_afectados
FROM (
    SELECT 
        ms.location_id,
        CASE 
            WHEN EXTRACT(EPOCH FROM (NOW() - st.receive_date)) / 3600 > 24 THEN 'URGENTE (>24h)'
            WHEN EXTRACT(EPOCH FROM (NOW() - st.receive_date)) / 3600 > 12 THEN 'CRÍTICO (12-24h)'
            WHEN EXTRACT(EPOCH FROM (NOW() - st.receive_date)) / 3600 > 8 THEN 'ALTO (8-12h)'
            ELSE 'NORMAL (<8h)'
        END as nivel_urgencia,
        CASE 
            WHEN EXTRACT(EPOCH FROM (NOW() - st.receive_date)) / 3600 > 24 THEN 1
            WHEN EXTRACT(EPOCH FROM (NOW() - st.receive_date)) / 3600 > 12 THEN 2
            WHEN EXTRACT(EPOCH FROM (NOW() - st.receive_date)) / 3600 > 8 THEN 3
            ELSE 4
        END as orden_urgencia
    FROM medical_supply ms
    JOIN supply_transfer st ON st.medical_supply_id = ms.id
    WHERE ms.status = 'recepcionado'
    AND ms.location_type = 'pavilion'
    AND st.status = 'recibido'
) urgency_data
GROUP BY nivel_urgencia, orden_urgencia
ORDER BY orden_urgencia;

-- ============================================================================
-- DATOS DE PRUEBA QA: INVENTARIO COMPLETO PABELLÓN 1 Y 2
-- Cubre: paginación (>10 filas), todos los estados de stock, indicadores de
--        vencimiento (vencido / próximo / ok) y cambio de pabellón.
-- ============================================================================

-- Nuevos lotes para visualización QA (IDs 51-76, NO modificar IDs anteriores)
INSERT INTO batch (id, expiration_date, amount, supplier_id, store_id, qr_code, surgery_id, location_type, location_id) VALUES
-- Pabellón 1 – 14 lotes (>10 → activa paginación)
(51, '2025-12-01', 20, 18, 1, 'BATCH_PAB1_051', NULL, 'store', 1), -- stock 0  | VENCIDO
(52, '2027-06-15', 20, 19, 1, 'BATCH_PAB1_052', NULL, 'store', 1), -- stock 0  | bueno
(53, '2026-01-10', 15, 18, 1, 'BATCH_PAB1_053', NULL, 'store', 1), -- stock 2  | VENCIDO
(54, '2026-03-10', 15, 20, 1, 'BATCH_PAB1_054', NULL, 'store', 1), -- stock 3  | próximo ~15d
(55, '2026-03-15', 15, 20, 1, 'BATCH_PAB1_055', NULL, 'store', 1), -- stock 4  | próximo ~20d
(56, '2026-04-14', 15, 19, 1, 'BATCH_PAB1_056', NULL, 'store', 1), -- stock 6  | próximo ~50d
(57, '2027-09-20', 15, 21, 1, 'BATCH_PAB1_057', NULL, 'store', 1), -- stock 7  | bueno
(58, '2028-01-15', 15, 21, 1, 'BATCH_PAB1_058', NULL, 'store', 1), -- stock 8  | bueno
(59, '2026-03-05', 25, 22, 1, 'BATCH_PAB1_059', NULL, 'store', 1), -- stock 12 | próximo ~10d
(60, '2027-12-31', 25, 22, 1, 'BATCH_PAB1_060', NULL, 'store', 1), -- stock 15 | bueno
(61, '2028-06-30', 15, 23, 1, 'BATCH_PAB1_061', NULL, 'store', 1), -- stock 5  | bueno
(62, '2028-10-15', 15, 18, 1, 'BATCH_PAB1_062', NULL, 'store', 1), -- stock 1  | bueno
(63, '2025-08-10', 25, 23, 1, 'BATCH_PAB1_063', NULL, 'store', 1), -- stock 20 | VENCIDO
(64, '2026-03-20', 15, 24, 1, 'BATCH_PAB1_064', NULL, 'store', 1), -- stock 9  | próximo ~25d
-- Pabellón 2 – 12 lotes (cambio de pabellón + paginación)
(65, '2027-08-15', 20, 25, 1, 'BATCH_PAB2_065', NULL, 'store', 1), -- stock 15 | bueno
(66, '2027-11-20', 20, 25, 1, 'BATCH_PAB2_066', NULL, 'store', 1), -- stock 10 | bueno
(67, '2028-02-10', 15, 26, 1, 'BATCH_PAB2_067', NULL, 'store', 1), -- stock 8  | bueno
(68, '2027-05-30', 15, 26, 1, 'BATCH_PAB2_068', NULL, 'store', 1), -- stock 5  | bueno
(69, '2026-04-30', 10, 27, 1, 'BATCH_PAB2_069', NULL, 'store', 1), -- stock 3  | próximo ~66d
(70, '2025-10-15', 15, 27, 1, 'BATCH_PAB2_070', NULL, 'store', 1), -- stock 0  | VENCIDO
(71, '2028-09-01', 20, 28, 1, 'BATCH_PAB2_071', NULL, 'store', 1), -- stock 12 | bueno
(72, '2026-02-20', 10, 28, 1, 'BATCH_PAB2_072', NULL, 'store', 1), -- stock 2  | VENCIDO
(73, '2027-03-10', 15, 29, 1, 'BATCH_PAB2_073', NULL, 'store', 1), -- stock 7  | bueno
(74, '2028-12-31', 25, 29, 1, 'BATCH_PAB2_074', NULL, 'store', 1), -- stock 20 | bueno
(75, '2025-11-05', 10, 30, 1, 'BATCH_PAB2_075', NULL, 'store', 1), -- stock 4  | VENCIDO
(76, '2027-07-20', 20, 30, 1, 'BATCH_PAB2_076', NULL, 'store', 1)  -- stock 11 | bueno
ON CONFLICT (id) DO NOTHING;

-- Generar medical_supply para los nuevos lotes
DO $$
DECLARE
    r RECORD;
    i INTEGER;
    bamt INTEGER;
BEGIN
    FOR r IN (
        SELECT * FROM (VALUES
            (51,1001),(52,1002),(53,1003),(54,1004),(55,1005),(56,1006),(57,1007),
            (58,1008),(59,1009),(60,1010),(61,1101),(62,1102),(63,1103),(64,1104),
            (65,1201),(66,1202),(67,1203),(68,1204),(69,1205),(70,1301),(71,1302),
            (72,1303),(73,1304),(74,1305),(75,1401),(76,1402)
        ) AS t(bid, scode)
    ) LOOP
        SELECT amount INTO bamt FROM batch WHERE id = r.bid;
        FOR i IN 1..bamt LOOP
            INSERT INTO medical_supply (code, batch_id, qr_code, status, location_type, location_id)
            VALUES (
                r.scode, r.bid,
                'SUPPLY_' || r.bid || '_' || LPAD(i::TEXT, 4, '0'),
                'disponible', 'store', 1
            ) ON CONFLICT (qr_code) DO NOTHING;
        END LOOP;
    END LOOP;
END $$;

-- Inventario Pabellón 1 – 14 filas con todos los estados visuales:
--   stock 0 (rojo) · stock 1–4 (naranja + "Stock bajo") · stock 5–9 (amarillo) · stock 10+ (verde)
--   vencido · próximo vencimiento (≤30d) · próximo vencimiento (≤90d) · bueno
INSERT INTO pavilion_inventory_summary (
    pavilion_id, batch_id, supply_code,
    total_received, current_available, total_consumed, total_returned,
    last_received_date, created_at, updated_at
) VALUES
-- stock 0 – rojo
(1, 51, 1001, 20, 0,  20, 0, NOW() - INTERVAL '2 days',  NOW(), NOW()), -- VENCIDO
(1, 52, 1002, 20, 0,  20, 0, NOW() - INTERVAL '3 days',  NOW(), NOW()), -- bueno (agotado)
-- stock 1–4 – naranja + "⚠️ Stock bajo"
(1, 62, 1102, 15, 1,  14, 0, NOW() - INTERVAL '3 hours', NOW(), NOW()), -- bueno
(1, 53, 1003, 15, 2,  13, 0, NOW() - INTERVAL '1 day',   NOW(), NOW()), -- VENCIDO
(1, 54, 1004, 15, 3,  12, 0, NOW() - INTERVAL '6 hours', NOW(), NOW()), -- próximo 15d
(1, 55, 1005, 15, 4,  11, 0, NOW() - INTERVAL '4 hours', NOW(), NOW()), -- próximo 20d
-- stock 5–9 – amarillo
(1, 61, 1101, 15, 5,  10, 0, NOW() - INTERVAL '1 hour',  NOW(), NOW()), -- bueno
(1, 56, 1006, 15, 6,  9,  0, NOW() - INTERVAL '2 hours', NOW(), NOW()), -- próximo 50d
(1, 57, 1007, 15, 7,  8,  0, NOW() - INTERVAL '5 hours', NOW(), NOW()), -- bueno
(1, 58, 1008, 15, 8,  7,  0, NOW() - INTERVAL '1 day',   NOW(), NOW()), -- bueno
(1, 64, 1104, 15, 9,  6,  0, NOW() - INTERVAL '2 days',  NOW(), NOW()), -- próximo 25d
-- stock 10+ – verde
(1, 59, 1009, 25, 12, 13, 0, NOW() - INTERVAL '3 hours', NOW(), NOW()), -- próximo 10d
(1, 60, 1010, 25, 15, 10, 0, NOW() - INTERVAL '6 hours', NOW(), NOW()), -- bueno
(1, 63, 1103, 25, 20, 5,  0, NOW() - INTERVAL '2 days',  NOW(), NOW())  -- VENCIDO
ON CONFLICT (pavilion_id, batch_id)
DO UPDATE SET
    total_received     = EXCLUDED.total_received,
    current_available  = EXCLUDED.current_available,
    total_consumed     = EXCLUDED.total_consumed,
    total_returned     = EXCLUDED.total_returned,
    last_received_date = EXCLUDED.last_received_date,
    updated_at         = NOW();

-- Inventario Pabellón 2 – 12 filas (cambio de pabellón + paginación completa)
INSERT INTO pavilion_inventory_summary (
    pavilion_id, batch_id, supply_code,
    total_received, current_available, total_consumed, total_returned,
    last_received_date, created_at, updated_at
) VALUES
(2, 70, 1301, 15, 0,  15, 0, NOW() - INTERVAL '5 days',  NOW(), NOW()), -- stock 0  VENCIDO
(2, 72, 1303, 10, 2,  8,  0, NOW() - INTERVAL '6 hours', NOW(), NOW()), -- stock 2  VENCIDO
(2, 69, 1205, 10, 3,  7,  0, NOW() - INTERVAL '1 hour',  NOW(), NOW()), -- stock 3  próximo ~66d
(2, 75, 1401, 10, 4,  6,  0, NOW() - INTERVAL '4 hours', NOW(), NOW()), -- stock 4  VENCIDO
(2, 68, 1204, 15, 5,  10, 0, NOW() - INTERVAL '3 hours', NOW(), NOW()), -- stock 5  bueno
(2, 67, 1203, 15, 8,  7,  0, NOW() - INTERVAL '4 hours', NOW(), NOW()), -- stock 8  bueno
(2, 73, 1304, 15, 7,  8,  0, NOW() - INTERVAL '1 day',   NOW(), NOW()), -- stock 7  bueno
(2, 66, 1202, 20, 10, 10, 0, NOW() - INTERVAL '2 days',  NOW(), NOW()), -- stock 10 bueno
(2, 76, 1402, 20, 11, 9,  0, NOW() - INTERVAL '5 hours', NOW(), NOW()), -- stock 11 bueno
(2, 71, 1302, 20, 12, 8,  0, NOW() - INTERVAL '2 hours', NOW(), NOW()), -- stock 12 bueno
(2, 65, 1201, 20, 15, 5,  0, NOW() - INTERVAL '1 day',   NOW(), NOW()), -- stock 15 bueno
(2, 74, 1305, 25, 20, 5,  0, NOW() - INTERVAL '3 days',  NOW(), NOW())  -- stock 20 bueno
ON CONFLICT (pavilion_id, batch_id)
DO UPDATE SET
    total_received     = EXCLUDED.total_received,
    current_available  = EXCLUDED.current_available,
    total_consumed     = EXCLUDED.total_consumed,
    total_returned     = EXCLUDED.total_returned,
    last_received_date = EXCLUDED.last_received_date,
    updated_at         = NOW();

-- ============================================================================
-- DATOS DE PRUEBA QA: INSUMOS EN TRÁNSITO AL PABELLÓN 3
-- El backend busca medical_supply.status = 'en_camino_a_pabellon'
-- con location_type = 'pavilion' y location_id = 3.
-- Se usan QR codes de batches 7, 8 y 10 que no tienen transferencias previas.
-- ============================================================================

-- Marcar 7 insumos como en camino al pabellón 3
ALTER TABLE medical_supply DISABLE TRIGGER trg_update_medical_supply_updated_at;

UPDATE medical_supply SET
    status        = 'en_camino_a_pabellon',
    location_type = 'pavilion',
    location_id   = 3,
    in_transit    = TRUE,
    transfer_date = NOW() - INTERVAL '45 minutes',
    transferred_by = '11111111-1'
WHERE qr_code IN (
    'SUPPLY_7_0001',  -- Batas Quirúrgicas (code 1007)
    'SUPPLY_7_0002',
    'SUPPLY_7_0003',
    'SUPPLY_8_0001',  -- Campos Quirúrgicos (code 1008)
    'SUPPLY_8_0002',
    'SUPPLY_10_0001', -- Mascarillas Estándar (code 1010)
    'SUPPLY_10_0002'
);

ALTER TABLE medical_supply ENABLE TRIGGER trg_update_medical_supply_updated_at;

-- Registrar las transferencias correspondientes en supply_transfer
INSERT INTO supply_transfer (
    transfer_code, qr_code, medical_supply_id,
    origin_type, origin_id, destination_type, destination_id,
    sent_by, sent_by_name,
    picked_up_by, picked_up_by_name, picked_up_date,
    status, transfer_reason, send_date, notes
)
SELECT
    'TRANS-TRAN-PAB3-' || LPAD(ROW_NUMBER() OVER (ORDER BY ms.qr_code)::TEXT, 4, '0'),
    ms.qr_code,
    ms.id,
    'store', 1, 'pavilion', 3,
    '11111111-1', 'Encargado Bodega',
    '22222222-2', 'María González', NOW() - INTERVAL '45 minutes',
    'en_transito',
    'Insumos en camino a Pabellón 3 para pruebas QA',
    NOW() - INTERVAL '1 hour',
    'En camino al pabellón 3'
FROM medical_supply ms
WHERE ms.qr_code IN (
    'SUPPLY_7_0001','SUPPLY_7_0002','SUPPLY_7_0003',
    'SUPPLY_8_0001','SUPPLY_8_0002',
    'SUPPLY_10_0001','SUPPLY_10_0002'
)
ON CONFLICT (transfer_code) DO NOTHING;

-- Inventario actual en pabellones
SELECT 
    '
INVENTARIO ACTUAL EN PABELLONES' as titulo;

SELECT 
    pabellon,
    tipos_insumo,
    total_insumos,
    pendientes_retorno,
    insumos_validos
FROM (
    SELECT 
        p.name as pabellon,
        p.id as pabellon_id,
        COUNT(DISTINCT ms.code) as tipos_insumo,
        COUNT(*) as total_insumos,
        SUM(CASE WHEN EXTRACT(EPOCH FROM (NOW() - st.receive_date)) / 3600 > 8 THEN 1 ELSE 0 END) as pendientes_retorno,
        SUM(CASE WHEN EXTRACT(EPOCH FROM (NOW() - st.receive_date)) / 3600 <= 8 THEN 1 ELSE 0 END) as insumos_validos
    FROM medical_supply ms
    JOIN supply_transfer st ON st.medical_supply_id = ms.id
    JOIN pavilion p ON p.id = ms.location_id
    WHERE ms.status = 'recepcionado'
    AND ms.location_type = 'pavilion'
    AND st.status = 'recibido'
    GROUP BY p.id, p.name
) inventory
ORDER BY pendientes_retorno DESC, total_insumos DESC;

-- ============================================================================
-- RE-AJUSTE DE SECUENCIAS (después de todos los INSERTs con IDs explícitos)
-- ============================================================================
SELECT setval('batch_id_seq', (SELECT COALESCE(MAX(id), 0) FROM batch));
SELECT setval('medical_supply_id_seq', (SELECT COALESCE(MAX(id), 0) FROM medical_supply));
SELECT setval('supplier_config_id_seq', (SELECT COALESCE(MAX(id), 0) FROM supplier_config));