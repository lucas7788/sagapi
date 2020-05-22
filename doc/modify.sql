alter table tbl_api_basic_info modify column `ApiType` varchar(255) NOT NULL DEFAULT '';
alter table tbl_api_basic_info modify column `ApiProvider` varchar(255) NOT NULL DEFAULT '';
alter table tbl_api_basic_info modify column `ApiDesc` varchar(255) NOT NULL DEFAULT '';
alter table tbl_api_basic_info modify column `InvokeFrequency` bigint(20) NOT NULL DEFAULT 0;

alter table tbl_api_basic_info add column  `ApiSagaUrlKey` varchar(100) NOT NULL after ApiProvider;
alter table tbl_api_basic_info add column  `ErrorDesc` varchar(255) NOT NULL DEFAULT '' after ApiDesc;


alter table tbl_api_basic_info add column  `ApiState` int(11) NOT NULL DEFAULT 0 after InvokeFrequency;
alter table tbl_api_basic_info add column  `RequestType` varchar(20) NOT NULL after ApiState;
alter table tbl_api_basic_info add column  `Mark` varchar(100) NOT NULL DEFAULT '' after RequestType;
alter table tbl_api_basic_info add column  `ResponseParam` varchar(255) NOT NULL DEFAULT '' after Mark;
alter table tbl_api_basic_info add column  `ResponseExample` varchar(2000) NOT NULL DEFAULT '' after ResponseParam;
alter table tbl_api_basic_info add column  `DataDesc` varchar(255) NOT NULL DEFAULT '' after ResponseExample;
alter table tbl_api_basic_info add column  `DataSource` varchar(255) NOT NULL DEFAULT '' after DataDesc;
alter table tbl_api_basic_info add column  `ApplicationScenario` varchar(255) NOT NULL DEFAULT '' after DataSource;
alter table tbl_api_basic_info add column  `ApiKind` int(11) NOT NULL DEFAULT 1 after ApplicationScenario;
alter table tbl_api_basic_info add column  `OntId` varchar(50) NOT NULL DEFAULT '' after ApiKind;
alter table tbl_api_basic_info add column  `Author` varchar(50) NOT NULL DEFAULT '' after OntId;
alter table tbl_api_basic_info add column  `ResourceId` varchar(255) NOT NULL DEFAULT '' after Author;
alter table tbl_api_basic_info add column  `TokenHash` char(255) NOT NULL DEFAULT '' after ResourceId;
alter table tbl_api_basic_info add column  `OwnerAddress` varchar(255) NOT NULL DEFAULT '' after TokenHash;

alter table tbl_api_basic_info add key(ApiState);
alter table tbl_api_basic_info add key(OntId);
alter table tbl_api_basic_info add key(Author);
alter table tbl_api_basic_info add key(ApiKind);


alter table tbl_api_basic_info add unique(ApiProvider);
alter table tbl_api_basic_info add unique(ApiSagaUrlKey);

ALTER TABLE tbl_api_basic_info AUTO_INCREMENT = 1;

insert into tbl_api_basic_info (Coin,Icon, Title, ApiProvider, ApiSagaUrlKey, ApiUrl, Price, ApiDesc,Specifications,Popularity,Delay,SuccessRate,InvokeFrequency,ApiState,RequestType, Mark, ResponseParam, ResponseExample, DataDesc, DataSource,ApplicationScenario) values ('ONG','https://www.sagamarket.io/img/icons/Geography.svg','Daily astronomical picture','https://api.nasa.gov/planetary/apod','sagaurl_dae5dcb0-90fd-11ea-8fa6-27687839068c', 'https://api.sagamarket.io/api/v1/nasa/apod/:apikey', '0', 'Get daily astronomical information',1,100,0,100,0,1,'GET','mark','ResponseParam', '{"copyright":"Juan Filas","date":"2020-04-20","explanation":"To some, it looks like a giant chicken running across the sky. To others, it looks like a gaseous nebula where star formation takes place. Cataloged as IC 2944, the Running Chicken Nebula spans about 100 light years and lies about 6,000 light years away toward the constellation of the Centaur (Centaurus).  The featured image, shown in scientifically assigned colors, was captured recently in a 12-hour exposure. The star cluster Collinder 249 is visible embedded in the nebula''s glowing gas.  Although difficult to discern here, several dark molecular clouds with distinct shapes can be found inside the nebula.","hdurl":"https://apod.nasa.gov/apod/image/2004/IC2944_Filas_3320.jpg","media_type":"image","service_version":"v1","title":"IC 2944: The Running Chicken Nebula","url":"https://apod.nasa.gov/apod/image/2004/IC2944_Filas_960.jpg"}','Daily astronomical picture','nasa','');
insert into tbl_api_basic_info (Coin,Icon, Title, ApiProvider, ApiSagaUrlKey, ApiUrl, Price, ApiDesc,Specifications,Popularity,Delay,SuccessRate,InvokeFrequency,ApiState,RequestType, Mark, ResponseParam, ResponseExample, DataDesc, DataSource,ApplicationScenario) values ('ONG','https://www.sagamarket.io/img/icons/Geography.svg','Near Earth Asteroid information','https://api.nasa.gov/neo/rest/v1/feed','sagaurl_bb747994-90fe-11ea-b9de-efc42e86bd80', 'https://api.sagamarket.io/api/v1/nasa/feed/:startdate/:enddate/:apikey', '0', 'Near Earth Asteroid information',1,100,0,100,0,1,'GET','mark','ResponseParam','{"links":{"next":"http://www.neowsapp.com/rest/v1/feed?start_date=2015-09-08&end_date=2015-09-08&detailed=false&api_key=DEMO_KEY","prev":"http://www.neowsapp.com/rest/v1/feed?start_date=2015-09-06&end_date=2015-09-06&detailed=false&api_key=DEMO_KEY","self":"http://www.neowsapp.com/rest/v1/feed?start_date=2015-09-07&end_date=2015-09-07&detailed=false&api_key=DEMO_KEY"},"element_count":12,"near_earth_objects":{"2015-09-07":[{"links":{"self":"http://www.neowsapp.com/rest/v1/neo/3726788?api_key=DEMO_KEY"},"id":"3726788","neo_reference_id":"3726788","name":"(2015 RG2)","nasa_jpl_url":"http://ssd.jpl.nasa.gov/sbdb.cgi?sstr=3726788","absolute_magnitude_h":26.7,"estimated_diameter":{"kilometers":{"estimated_diameter_min":0.0121494041,"estimated_diameter_max":0.0271668934},"meters":{"estimated_diameter_min":12.14940408,"estimated_diameter_max":27.1668934089},"miles":{"estimated_diameter_min":0.0075492874,"estimated_diameter_max":0.0168807197},"feet":{"estimated_diameter_min":39.8602508817,"estimated_diameter_max":89.1302305717}},"is_potentially_hazardous_asteroid":false,"close_approach_data":[{"close_approach_date":"2015-09-07","close_approach_date_full":"2015-Sep-07 17:58","epoch_date_close_approach":1441648680000,"relative_velocity":{"kilometers_per_second":"8.0887368746","kilometers_per_hour":"29119.4527484721","miles_per_hour":"18093.6955147381"},"miss_distance":{"astronomical":"0.0163818512","lunar":"6.3725401168","kilometers":"2450690.046176944","miles":"1522788.1820680672"},"orbiting_body":"Earth"}],"is_sentry_object":false}]}}','Near Earth Asteroid information','nasa','');
insert into tbl_api_basic_info (Coin,ApiType,Icon, Title, ApiProvider, ApiSagaUrlKey, ApiUrl, Price, ApiDesc,Specifications,Popularity,Delay,SuccessRate,InvokeFrequency,ApiState,RequestType, Mark, ResponseParam, ResponseExample, DataDesc, DataSource,ApplicationScenario,ApiKind,ResourceId,TokenHash,OwnerAddress) values ('ONG','Weather Forecast','https://www.sagamarket.io/img/icons/tianqi.svg','Global Weather','https://api.stormglass.io/v2/weather/point?params=airTemperature','sagaurl_7d07a1fa-95da-11ea-9c30-5f81df211d0b','https://api.sagamarket.io/api/v1/data_source/sagaurl_7d07a1fa-95da-11ea-9c30-5f81df211d0b/:apikey','0','Fetch weather data for any coordinate on the globe',1,100,0,100,0,1,'GET','mark','ResponseParam','{"hours":[{"airTemperature":{"noaa":22.26,"sg":22.26},"time":"2020-05-02T00:00:00+00:00"}]}',' Global marine weather as well as weather for land and lakes','Storm Glass','',2,'539cedd7-63c9-4b1b-8ae3-87d026e311e2','997a019b98e9847a0c3343bae2d7ad8d931bf784e11ad736539702b661b7f163','');





alter table tbl_order modify column `OrderTime` bigint(20) NOT NULL DEFAULT 0 COMMENT '下单时间';
alter table tbl_order modify column `PayTime` bigint(20) NOT NULL DEFAULT 0 COMMENT '支付时间';
alter table tbl_order add column `OrderKind` int(11) NOT NULL after Coin;
alter table tbl_order add column `Request` varchar(4095) NOT NULL COMMENT '币种' after OrderKind;
alter table tbl_order add column `Result` varchar(4095) NOT NULL COMMENT '币种' after Request;
alter table tbl_order modify column `OrderStatus` tinyint(4) NOT NULL DEFAULT 0;
alter table tbl_order modify column `Amount` varchar(255) NOT NULL DEFAULT '';



alter table tbl_api_key modify column `RequestLimit` bigint(20) NOT NULL DEFAULT 0;
alter table tbl_api_key modify column `UsedNum` bigint(20) NOT NULL DEFAULT 0;
alter table tbl_api_key add key(`ApiId`);
alter table tbl_api_key add constraint `tbl_api_key_ibfk_2` FOREIGN KEY (`ApiId`) REFERENCES `tbl_api_basic_info` (`ApiId`);



alter table tbl_api_test_key add column `OrderId` varchar(20) NOT NULL DEFAULT 'TST_ORDER' after ApiId;
alter table tbl_api_test_key modify column `RequestLimit` bigint(20) NOT NULL DEFAULT 0;
alter table tbl_api_test_key modify column `UsedNum` bigint(20) NOT NULL DEFAULT 0;
alter table tbl_api_test_key add key(OntId);


alter table tbl_qr_code modify column `Exp` bigint(20) NOT NULL DEFAULT 0;
alter table tbl_qr_code add column `ContractType` varchar(10) NOT NULL DEFAULT '' after QrCodeDesc;
