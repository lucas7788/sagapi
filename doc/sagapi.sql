DROP TABLE IF EXISTS `tbl_api_basic_info`;
create table tbl_api_basic_info
(
 ApiId int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
 ApiLogo varchar(100) not null default '' COMMENT '',
 ApiName varchar(100) not null  default '' COMMENT '',
 ApiProvider varchar(100) not null default '' COMMENT '',
 ApiUrl varchar(100) not null  default '' COMMENT '',
 ApiPrice varchar(100) not null  default '' COMMENT '',
 ApiDesc varchar(100) not null  default '' COMMENT '',
 Specifications int(11) not null  default 0 COMMENT '规格',
 Popularity int(11) not null default 0 COMMENT '流行度',
 Delay int(11) not null default 0 COMMENT '',
 SuccessRate int(11) not null default 0 COMMENT '',
 InvokeFrequency int(11) not null default 0 COMMENT '',
 PRIMARY KEY (ApiId)
);

DROP TABLE IF EXISTS `tbl_api_detail_info`;
create table tbl_api_detail_info
(
 Id int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
 ApiId int(11) not null,
 Mark varchar(100) not null default '' COMMENT '',
 ResponseParam varchar(100) not null default ''  COMMENT '',
 ResponseExample varchar(100) not null default ''  COMMENT '',
 DataDesc varchar(100) not null default '' COMMENT '',
 DataSource varchar(100) not null default ''  COMMENT '',
 ApplicationScenario varchar(100) not null default '' COMMENT '',
 PRIMARY KEY (Id),
 foreign key(ApiId) references tbl_api_basic_info(ApiId)
);
