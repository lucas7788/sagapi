DROP TABLE IF EXISTS `tbl_api_test_key`;
DROP TABLE IF EXISTS `tbl_qr_code`;
DROP TABLE IF EXISTS `tbl_api_key`;
DROP TABLE IF EXISTS `tbl_order`;
DROP TABLE IF EXISTS `tbl_error_code`;
DROP TABLE IF EXISTS `tbl_request_param`;
DROP TABLE IF EXISTS `tbl_specifications`;
DROP TABLE IF EXISTS `tbl_api_tag`;
DROP TABLE IF EXISTS `tbl_tag`;
DROP TABLE IF EXISTS `tbl_category`;
DROP TABLE IF EXISTS `tbl_api_basic_info`;

create table tbl_api_basic_info
(
 ApiId INT NOT NULL AUTO_INCREMENT COMMENT '主键',
 Coin  varchar(10) NOT NULL DEFAULT '' COMMENT '',
 ApiType varchar(10) NOT NULL DEFAULT '' COMMENT '',
 Icon text NOT NULL COMMENT '',
 Title varchar(100) NOT NULL  DEFAULT '' COMMENT '',
 ApiProvider varchar(255) unique NOT NULL DEFAULT '' COMMENT '',
 ApiSagaUrlKey varchar(100) unique NOT NULL COMMENT '',
 ApiUrl varchar(255) NOT NULL  DEFAULT '' COMMENT '',
 Price varchar(100) NOT NULL  DEFAULT '' COMMENT '',
 ApiDesc varchar(255) NOT NULL  DEFAULT '' COMMENT '',
 ErrorDesc varchar(255) NOT NULL  DEFAULT '' COMMENT '',
 Specifications INT NOT NULL  DEFAULT 0 COMMENT '规格',
 Popularity INT NOT NULL DEFAULT 0 COMMENT '流行度',
 Delay INT NOT NULL DEFAULT 0 COMMENT '',
 SuccessRate INT NOT NULL DEFAULT 0 COMMENT '',
 InvokeFrequency BIGINT NOT NULL DEFAULT 0 COMMENT '',
 ApiState INT NOT NULL DEFAULT 0 COMMENT '',
 RequestType varchar(20) NOT NULL COMMENT '',
 Mark varchar(100) NOT NULL DEFAULT '' COMMENT '',
 ResponseParam varchar(255) NOT NULL DEFAULT ''  COMMENT '',
 ResponseExample varchar(2000) NOT NULL DEFAULT ''  COMMENT '',
 DataDesc varchar(255) NOT NULL DEFAULT '' COMMENT '',
 DataSource varchar(255) NOT NULL DEFAULT ''  COMMENT '',
 ApplicationScenario varchar(255) NOT NULL DEFAULT '' COMMENT '',
 CreateTime TIMESTAMP DEFAULT current_timestamp,
 PRIMARY KEY (ApiId),
 INDEX(Price),
 INDEX(Title),
 INDEX(ApiDesc),
 INDEX(ApiState)
)DEFAULT charset=utf8;

create table tbl_category
(
 Id INT NOT NULL AUTO_INCREMENT COMMENT '主键',
 NameZh varchar(100) UNIQUE NOT NULL DEFAULT '' COMMENT '',
 NameEn varchar(100) UNIQUE NOT NULL  DEFAULT '' COMMENT '',
 Icon varchar(255) NOT NULL  DEFAULT '' COMMENT '',
 State TINYINT NOT NULL DEFAULT 1 COMMENT '0:delete, 1:active',
 PRIMARY KEY (Id)
)DEFAULT charset=utf8;

create table tbl_tag
(
 Id INT NOT NULL AUTO_INCREMENT COMMENT '主键',
 Name varchar(50) UNIQUE NOT NULL DEFAULT '' COMMENT '',
 CategoryId INT NOT NULL COMMENT '',
 State TINYINT NOT NULL DEFAULT 1 COMMENT '0:delete, 1:active',
 CreateTime TIMESTAMP DEFAULT current_timestamp,
 PRIMARY KEY (Id),
 INDEX(CategoryId),
 CONSTRAINT FK_CategoryId FOREIGN KEY (CategoryId) REFERENCES tbl_category(Id)
)DEFAULT charset=utf8;

create table tbl_api_tag
(
 Id INT NOT NULL AUTO_INCREMENT COMMENT '主键',
 ApiId INT NOT NULL COMMENT '',
 TagId INT NOT NULL COMMENT '',
 State TINYINT NOT NULL DEFAULT 1 COMMENT '0:delete, 1:active',
 CreateTime TIMESTAMP DEFAULT current_timestamp,
 PRIMARY KEY (id),
 INDEX(TagId),
 INDEX(ApiId),
 CONSTRAINT FK_ApiId FOREIGN KEY (ApiId) REFERENCES tbl_api_basic_info(ApiId),
 CONSTRAINT FK_TagId FOREIGN KEY (TagId) REFERENCES tbl_tag(Id)
)DEFAULT charset=utf8;

create table tbl_specifications
(
 Id INT NOT NULL AUTO_INCREMENT COMMENT '主键',
 ApiId INT NOT NULL,
 Price  varchar(50) NOT NULL DEFAULT '' COMMENT '',
 Amount BIGINT NOT NULL DEFAULT 0,
 PRIMARY KEY (Id),
 CONSTRAINT FK_specifications_id FOREIGN KEY (ApiId) REFERENCES tbl_api_basic_info(ApiId)
)DEFAULT charset=utf8;

create table tbl_request_param (
  Id INT NOT NULL AUTO_INCREMENT COMMENT '主键',
  ApiId INT NOT NULL,
  ParamName varchar(50) NOT NULL DEFAULT '',
  Required  TINYINT NOT NULL,
  ParamWhere INT NOT NULL DEFAULT 0,
  ParamType varchar(10) NOT NULL DEFAULT '',
  Note varchar(255) NOT NULL DEFAULT '',
  ValueDesc varchar(50) NOT NULL DEFAULT '',
  PRIMARY KEY (Id),
  CONSTRAINT FK_request_param_id FOREIGN KEY (ApiId) REFERENCES tbl_api_basic_info(ApiId),
  INDEX(ApiId)
)DEFAULT charset=utf8;

create table tbl_error_code (
  Id INT NOT NULL AUTO_INCREMENT COMMENT '主键',
  ApiId INT NOT NULL,
  ErrorCode INT NOT NULL,
  ErrorDesc varchar(50) NOT NULL DEFAULT '',
  PRIMARY KEY (Id),
  CONSTRAINT FK_error_code_id FOREIGN KEY (ApiId) REFERENCES tbl_api_basic_info(ApiId),
  INDEX(ApiId)
)DEFAULT charset=utf8;


create table tbl_order (
  OrderId varchar(100) unique NOT NULL COMMENT '',
  Title varchar(100) NOT NULL COMMENT '',
  ProductName varchar(50) NOT NULL DEFAULT '' COMMENT '',
  OrderType varchar(50) NOT NULL DEFAULT ''  COMMENT '',
  OrderTime BIGINT NOT NULL DEFAULT 0 COMMENT '下单时间',
  PayTime  BIGINT NOT NULL DEFAULT 0  COMMENT '支付时间',
  OrderStatus TINYINT NOT NULL DEFAULT 0,
  Amount varchar(255) NOT NULL DEFAULT '' COMMENT '',
  OntId varchar(50) NOT NULL DEFAULT '' COMMENT '用户ontid',
  UserName varchar(50) NOT NULL DEFAULT '' COMMENT '',
  TxHash varchar(100) NOT NULL  DEFAULT '' COMMENT '支付交易hash',
  Price varchar(50) NOT NULL DEFAULT ''  COMMENT '',
  ApiId INT NOT NULL COMMENT '',
  ApiUrl varchar(255) NOT NULL  DEFAULT '' COMMENT '',
  SpecificationsId INT NOT NULL COMMENT '规格',
  Coin varchar(20) NOT NULL COMMENT '币种',
  PRIMARY KEY (OrderId),
  CONSTRAINT FK_tbl_order_id FOREIGN KEY (ApiId) REFERENCES tbl_api_basic_info(ApiId),
  INDEX(OntId)
)DEFAULT charset=utf8;


create table tbl_api_key (
  Id INT NOT NULL AUTO_INCREMENT COMMENT '主键',
  ApiKey varchar(50) unique NOT NULL  DEFAULT '',
  ApiId INT NOT NULL,
  OrderId varchar(100) unique NOT NULL COMMENT '',
  RequestLimit BIGINT NOT NULL DEFAULT 0,
  UsedNum BIGINT NOT NULL DEFAULT 0,
  OntId varchar(50) NOT NULL DEFAULT '',
  PRIMARY KEY (Id),
  foreign key(OrderId) references tbl_order(OrderId),
  foreign key(ApiId) references tbl_api_basic_info(ApiId),
  INDEX(ApiKey),
  INDEX(OntId)
)DEFAULT charset=utf8;

create table tbl_api_test_key (
  Id INT NOT NULL AUTO_INCREMENT COMMENT '主键',
  ApiKey varchar(50) unique NOT NULL  DEFAULT '',
  ApiId INT NOT NULL,
  RequestLimit BIGINT NOT NULL DEFAULT 0,
  UsedNum BIGINT NOT NULL DEFAULT 0,
  OntId varchar(50) NOT NULL DEFAULT '',
  PRIMARY KEY (Id),
  foreign key(ApiId) references tbl_api_basic_info(ApiId),
  INDEX(ApiId),
  INDEX(ApiKey),
  INDEX(OntId)
) DEFAULT charset=utf8;

CREATE TABLE `tbl_qr_code` (
  Id INT NOT NULL AUTO_INCREMENT COMMENT '主键',
  QrCodeId varchar(100) unique NOT NULL DEFAULT '',
  Ver varchar(50) NOT NULL DEFAULT '',
  OrderId varchar(100) NOT NULL DEFAULT '' ,
  Requester varchar(50) NOT NULL DEFAULT '',
  Signature varchar(200) NOT NULL DEFAULT '',
  Signer varchar(50) NOT NULL DEFAULT '',
  QrCodeData text,
  Callback varchar(400) NOT NULL DEFAULT '',
  Exp BIGINT NOT NULL DEFAULT 0,
  Chain varchar(50) NOT NULL DEFAULT '',
  QrCodeDesc varchar(100) NOT NULL DEFAULT '',
  PRIMARY KEY (Id),
  foreign key(OrderId) references tbl_order(OrderId),
  INDEX(QrCodeId)
)DEFAULT charset=utf8;
