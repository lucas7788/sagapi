INSERT INTO `tbl_category`(`Id`, `NameZh`, `NameEn`, `Icon`, `State`) VALUES (1, '所有', 'All', 'https://dev4.sagamarket.io/img/icons/Agriculture.svg', 1);
INSERT INTO `tbl_category`(`Id`, `NameZh`, `NameEn`, `Icon`, `State`) VALUES (2, '商业', 'Business', 'https://dev4.sagamarket.io/img/icons/Business.svg', 1);
INSERT INTO `tbl_category`(`Id`, `NameZh`, `NameEn`, `Icon`, `State`) VALUES (3, '常识', 'General', 'https://dev4.sagamarket.io/img/icons/General.svg', 1);
INSERT INTO `tbl_category`(`Id`, `NameZh`, `NameEn`, `Icon`, `State`) VALUES (4, '新冠专题', 'COVID-19', 'https://dev4.sagamarket.io/img/icons/COVID.svg', 1);
INSERT INTO `tbl_category`(`Id`, `NameZh`, `NameEn`, `Icon`, `State`) VALUES (5, '其它', 'Other', 'https://dev4.sagamarket.io/img/icons/Other.svg', 1);
INSERT INTO `tbl_category`(`Id`, `NameZh`, `NameEn`, `Icon`, `State`) VALUES (6, '科教', 'Science', 'https://dev4.sagamarket.io/img/icons/Science.svg', 1);
INSERT INTO `tbl_category`(`Id`, `NameZh`, `NameEn`, `Icon`, `State`) VALUES (7, '金融', 'Finance', 'https://dev4.sagamarket.io/img/icons/Finance.svg', 1);
INSERT INTO `tbl_category`(`Id`, `NameZh`, `NameEn`, `Icon`, `State`) VALUES (8, '娱乐', 'Entertainment', 'https://dev4.sagamarket.io/img/icons/Entertainment.svg', 1);
INSERT INTO `tbl_category`(`Id`, `NameZh`, `NameEn`, `Icon`, `State`) VALUES (9, '物联网', 'Internet of Things', 'https://dev4.sagamarket.io/img/icons/Internet.svg', 1);
INSERT INTO `tbl_category`(`Id`, `NameZh`, `NameEn`, `Icon`, `State`) VALUES (10, '地理', 'Geography', 'https://dev4.sagamarket.io/img/icons/Geography.svg', 1);
INSERT INTO `tbl_category`(`Id`, `NameZh`, `NameEn`, `Icon`, `State`) VALUES (11, '医疗', 'Medicine', 'https://dev4.sagamarket.io/img/icons/Medicine.svg', 1);
INSERT INTO `tbl_category`(`Id`, `NameZh`, `NameEn`, `Icon`, `State`) VALUES (12, '生活', 'Lifestyle', 'https://dev4.sagamarket.io/img/icons/Lifestyle.svg', 1);
INSERT INTO `tbl_category`(`Id`, `NameZh`, `NameEn`, `Icon`, `State`) VALUES (13, '社交', 'Social', 'https://dev4.sagamarket.io/img/icons/Social.svg', 1);
INSERT INTO `tbl_category`(`Id`, `NameZh`, `NameEn`, `Icon`, `State`) VALUES (14, '气象', 'Climate', 'https://dev4.sagamarket.io/img/icons/Climate.svg', 1);
INSERT INTO `tbl_category`(`Id`, `NameZh`, `NameEn`, `Icon`, `State`) VALUES (15, '城市', 'Demography', 'https://dev4.sagamarket.io/img/icons/Demography.svg', 1);
INSERT INTO `tbl_category`(`Id`, `NameZh`, `NameEn`, `Icon`, `State`) VALUES (16, '出行', 'Mobility', 'https://dev4.sagamarket.io/img/icons/Mobility.svg', 1);

INSERT INTO `tbl_tag`(`Id`, `Name`, `CategoryId`, `State`, `CreateTime`) VALUES (1, 'Agriculture', 1, 1, '2020-03-13 10:16:22');
INSERT INTO `tbl_tag`(`Id`, `Name`, `CategoryId`, `State`, `CreateTime`) VALUES (2, 'Business', 2, 1, '2020-03-13 10:16:22');
INSERT INTO `tbl_tag`(`Id`, `Name`, `CategoryId`, `State`, `CreateTime`) VALUES (3, 'General', 3, 1, '2020-03-13 10:16:22');
INSERT INTO `tbl_tag`(`Id`, `Name`, `CategoryId`, `State`, `CreateTime`) VALUES (4, 'COVID-19', 4, 1, '2020-03-13 10:16:22');
INSERT INTO `tbl_tag`(`Id`, `Name`, `CategoryId`, `State`, `CreateTime`) VALUES (5, 'Other', 5, 1, '2020-03-13 10:16:22');
INSERT INTO `tbl_tag`(`Id`, `Name`, `CategoryId`, `State`, `CreateTime`) VALUES (6, 'Science', 6, 1, '2020-03-13 10:16:22');
INSERT INTO `tbl_tag`(`Id`, `Name`, `CategoryId`, `State`, `CreateTime`) VALUES (7, 'Finance', 7, 1, '2020-03-13 10:16:22');
INSERT INTO `tbl_tag`(`Id`, `Name`, `CategoryId`, `State`, `CreateTime`) VALUES (8, 'Entertainment', 8, 1, '2020-03-13 10:16:22');
INSERT INTO `tbl_tag`(`Id`, `Name`, `CategoryId`, `State`, `CreateTime`) VALUES (9, 'Internet of Things', 9, 1, '2020-03-13 10:16:22');
INSERT INTO `tbl_tag`(`Id`, `Name`, `CategoryId`, `State`, `CreateTime`) VALUES (10, 'Geography', 10, 1, '2020-03-13 10:16:22');
INSERT INTO `tbl_tag`(`Id`, `Name`, `CategoryId`, `State`, `CreateTime`) VALUES (11, 'Medicine', 11, 1, '2020-03-13 10:16:22');
INSERT INTO `tbl_tag`(`Id`, `Name`, `CategoryId`, `State`, `CreateTime`) VALUES (12, 'Lifestyle', 12, 1, '2020-03-13 10:16:22');
INSERT INTO `tbl_tag`(`Id`, `Name`, `CategoryId`, `State`, `CreateTime`) VALUES (13, 'Social', 13, 1, '2020-03-13 10:16:22');
INSERT INTO `tbl_tag`(`Id`, `Name`, `CategoryId`, `State`, `CreateTime`) VALUES (14, 'Climate', 14, 1, '2020-03-13 10:16:22');
INSERT INTO `tbl_tag`(`Id`, `Name`, `CategoryId`, `State`, `CreateTime`) VALUES (15, 'Demography', 15, 1, '2020-03-13 10:16:22');
INSERT INTO `tbl_tag`(`Id`, `Name`, `CategoryId`, `State`, `CreateTime`) VALUES (16, 'Mobility', 16, 1, '2020-03-13 10:16:22');

insert into tbl_api_basic_info (Coin,Icon, Title, ApiProvider, ApiSagaUrlKey, ApiUrl, Price, ApiDesc,Specifications,Popularity,Delay,SuccessRate,InvokeFrequency,ApiState,RequestType, Mark, ResponseParam, ResponseExample, DataDesc, DataSource,ApplicationScenario) values ('ONG','https://dev4.sagamarket.io/img/icons/Geography.svg','Daily astronomical picture','https://api.nasa.gov/planetary/apod','sagaurl_dae5dcb0-90fd-11ea-8fa6-27687839068c', 'https://dev3.sagamarket.io/api/v1/nasa/apod/:apikey', '0', 'Get daily astronomical information',1,100,0,100,0,1,'GET','mark','ResponseParam', '{"copyright":"Juan Filas","date":"2020-04-20","explanation":"To some, it looks like a giant chicken running across the sky. To others, it looks like a gaseous nebula where star formation takes place. Cataloged as IC 2944, the Running Chicken Nebula spans about 100 light years and lies about 6,000 light years away toward the constellation of the Centaur (Centaurus).  The featured image, shown in scientifically assigned colors, was captured recently in a 12-hour exposure. The star cluster Collinder 249 is visible embedded in the nebula''s glowing gas.  Although difficult to discern here, several dark molecular clouds with distinct shapes can be found inside the nebula.","hdurl":"https://apod.nasa.gov/apod/image/2004/IC2944_Filas_3320.jpg","media_type":"image","service_version":"v1","title":"IC 2944: The Running Chicken Nebula","url":"https://apod.nasa.gov/apod/image/2004/IC2944_Filas_960.jpg"}','Daily astronomical picture','nasa','');
insert into tbl_api_basic_info (Coin,Icon, Title, ApiProvider, ApiSagaUrlKey, ApiUrl, Price, ApiDesc,Specifications,Popularity,Delay,SuccessRate,InvokeFrequency,ApiState,RequestType, Mark, ResponseParam, ResponseExample, DataDesc, DataSource,ApplicationScenario) values ('ONG','https://dev4.sagamarket.io/img/icons/Geography.svg','Near Earth Asteroid information','https://api.nasa.gov/neo/rest/v1/feed','sagaurl_bb747994-90fe-11ea-b9de-efc42e86bd80', 'https://dev3.sagamarket.io/api/v1/nasa/feed/:startdate/:enddate/:apikey', '0', 'Near Earth Asteroid information',1,100,0,100,0,1,'GET','mark','ResponseParam','{"links":{"next":"http://www.neowsapp.com/rest/v1/feed?start_date=2015-09-08&end_date=2015-09-08&detailed=false&api_key=DEMO_KEY","prev":"http://www.neowsapp.com/rest/v1/feed?start_date=2015-09-06&end_date=2015-09-06&detailed=false&api_key=DEMO_KEY","self":"http://www.neowsapp.com/rest/v1/feed?start_date=2015-09-07&end_date=2015-09-07&detailed=false&api_key=DEMO_KEY"},"element_count":12,"near_earth_objects":{"2015-09-07":[{"links":{"self":"http://www.neowsapp.com/rest/v1/neo/3726788?api_key=DEMO_KEY"},"id":"3726788","neo_reference_id":"3726788","name":"(2015 RG2)","nasa_jpl_url":"http://ssd.jpl.nasa.gov/sbdb.cgi?sstr=3726788","absolute_magnitude_h":26.7,"estimated_diameter":{"kilometers":{"estimated_diameter_min":0.0121494041,"estimated_diameter_max":0.0271668934},"meters":{"estimated_diameter_min":12.14940408,"estimated_diameter_max":27.1668934089},"miles":{"estimated_diameter_min":0.0075492874,"estimated_diameter_max":0.0168807197},"feet":{"estimated_diameter_min":39.8602508817,"estimated_diameter_max":89.1302305717}},"is_potentially_hazardous_asteroid":false,"close_approach_data":[{"close_approach_date":"2015-09-07","close_approach_date_full":"2015-Sep-07 17:58","epoch_date_close_approach":1441648680000,"relative_velocity":{"kilometers_per_second":"8.0887368746","kilometers_per_hour":"29119.4527484721","miles_per_hour":"18093.6955147381"},"miss_distance":{"astronomical":"0.0163818512","lunar":"6.3725401168","kilometers":"2450690.046176944","miles":"1522788.1820680672"},"orbiting_body":"Earth"}],"is_sentry_object":false}]}}','Near Earth Asteroid information','nasa','');
insert into tbl_api_basic_info (Coin,ApiType,Icon, Title, ApiProvider, ApiSagaUrlKey, ApiUrl, Price, ApiDesc,Specifications,Popularity,Delay,SuccessRate,InvokeFrequency,ApiState,RequestType, Mark, ResponseParam, ResponseExample, DataDesc, DataSource,ApplicationScenario,ApiKind) values ('ONG','Weather Forecast','https://dev4.sagamarket.io/img/icons/tianqi.svg','Weather Forecast','https://api.stormglass.io/v2/weather/point?params=airTemperature','sagaurl_7d07a1fa-95da-11ea-9c30-5f81df211d0b','https://dev3.sagamarket.io/api/v1/data_source/sagaurl_7d07a1fa-95da-11ea-9c30-5f81df211d0b/:apikey','0','Get weather forecast for a city by  analyzing historical data collected from worldwide weather stations',1,100,0,100,0,1,'GET','mark','ResponseParam','{"hours":[{"airTemperature":{"noaa":22.26,"sg":22.26},"time":"2020-05-02T00:00:00+00:00"}]}','airTemperature','Stormglass','',2);
insert into tbl_request_param (ApiId,ParamName,ParamType,ParamWhere,Required,Note,ValueDesc) values (2,'startDate','string',2,true,'2016-06-06','');
insert into tbl_request_param (ApiId,ParamName,ParamType,ParamWhere,Required,Note,ValueDesc) values (2,'endDate','string',2,true,'2016-06-07','(Specified time interval must not exceed 7 days)');
insert into tbl_request_param (ApiId,ParamName,ParamType,ParamWhere,Required,Note,ValueDesc) values (3,'lat','string',2,true,'31.2222222','');
insert into tbl_request_param (ApiId,ParamName,ParamType,ParamWhere,Required,Note,ValueDesc) values (3,'lng','string',2,true,'121.45','');
insert into tbl_request_param (ApiId,ParamName,ParamType,ParamWhere,Required,Note,ValueDesc) values (3,'start','string',2,true,'2020-05-07','');
insert into tbl_error_code (ErrorCode,ErrorDesc) values (40001,'inter error');
insert into tbl_error_code (ErrorCode,ErrorDesc) values (40002,'param error');
insert into tbl_error_code (ErrorCode,ErrorDesc) values (40005,'api key is nil');
insert into tbl_api_tag (ApiId, TagId,state) values (1,10,1);
insert into tbl_api_tag (ApiId, TagId,state) values (2,10,1);
insert into tbl_specifications (ApiId,Price,Amount) values (1,'0',500);
insert into tbl_specifications (ApiId,Price,Amount) values (1,'0.01',1000);
insert into tbl_specifications (ApiId,Price,Amount) values (1,'0.0075',2000);
insert into tbl_specifications (ApiId,Price,Amount) values (2,'0',500);
insert into tbl_specifications (ApiId,Price,Amount) values (2,'0.01',1000);
insert into tbl_specifications (ApiId,Price,Amount) values (2,'0.0075',2000);
insert into tbl_algorithm (AlgName, Provider, Description, Price, Coin,ResourceId,TokenHash,OwnerAddress) values ('storm gass predictor', 'Ontology','a prediction of wether','0.1','ONG','','','');
insert into tbl_env (EnvName, Provider, Description, Price, Coin,ResourceId,TokenHash,OwnerAddress) values ('normal server', 'Ontology', 'a linux server', '0.1', 'ONG','','','');
insert into tbl_api_algorithm (ApiId, AlgorithmId) values (3,1);
insert into tbl_algorithm_env (AlgorithmId, EnvId) values (1,1);
insert into tbl_tool_box (Title,ToolBoxDesc,Icon) values ('Weather Forecast','Get weather forecast for a city by analyzing historical data collected from worldwide weather stations','https://dev4.sagamarket.io/img/icons/tianqi.svg')
