drop database if exists finepoints;
create database  finepoints character set utf8;

use finepoints



drop table if exists t_user;
create table t_user
(
   user_id              int primary key auto_increment,
   user_name            varchar(50),
   chat_id              varchar(50),
   join_date            timestamp
);

drop table if exists t_user_detail;
create table t_user_detail
(
   detail_id            int primary key auto_increment,
   user_id              int,
   sex                  varchar(6),
   age                  int,
   work_age             int,
   high_edu             varchar(30),
   tags                 varchar(50),
   create_time          timestamp
);

alter table t_user_detail add constraint FK_Reference_3 foreign key (user_id)
      references t_user (user_id) on delete restrict on update restrict;



drop table if exists t_enterprise;

create table t_enterprise
(
   et_id                int primary key auto_increment,
   et_name              varchar(100),
   credit_id            varchar(50),
   business_lic         varchar(100),
   phone                varchar(20),
   login_id             varchar(30),
   et_email             varchar(40),
   passwd               varchar(30),
   et_logo              varchar(100),
   et_url               varchar(100),
   et_info              varchar(500),
   et_scale             varchar(30),           
   create_time          timestamp
);

create unique index UK_enterprise on t_enterprise (phone);

drop table if exists t_enterprise_hr;

create table t_enterprise_hr
(
   rp_hr_id             int primary key auto_increment,
   et_id                int,
   user_id              int,
   remark               varchar(50),     
   create_time          timestamp
);



drop table if exists t_job_type;

create table t_job_type(
    job_type_id int primary key auto_increment,
    job_type_name varchar(30) not null,
    parent_jt_id int,
    create_time timestamp
);

drop table t_job_deliver;
drop table t_job_info ;

create table t_job_info
(
   job_id               int primary key auto_increment,
   et_id                int,
   job_name             varchar(100),
   sal_scope            varchar(30),
   high_edu             varchar(30),
   experience           varchar(30),
   job_type_id          int,
   work_city            varchar(30),
   work_addr            varchar(100),
   job_property         varchar(30),
   is_check             bool,
   is_headhunter        bool,
   job_detail           text,
   status               char(1),
   create_time          timestamp
);

alter table t_job_info add constraint FK_Reference_1 foreign key (et_id)
      references t_enterprise(et_id) on delete restrict on update restrict;




create table t_job_deliver
(
   jd_id                int primary key auto_increment,
   user_id              int,
   job_id               int,
   sal_scope            int,
   create_time          timestamp,
   similar              int,
   comment              varchar(200),
   status               char(1)
);

alter table t_job_deliver add constraint FK_Reference_2 foreign key (user_id)
      references t_user (user_id) on delete restrict on update restrict;

alter table t_job_deliver add constraint FK_Reference_4 foreign key (job_id)
      references t_job_info (job_id) on delete restrict on update restrict;


-- 参数表
drop table  if exists  t_params;
create table t_params(
    param_id int primary key auto_increment,
	table_name varchar(30),
    column_name varchar(30) not null,
    column_value varchar(30),
    param_value  varchar(30),
    remark varchar(50)
);
delete from t_params;

-- 公司规模参数

insert into t_params(param_id,table_name,column_name,column_value,param_value,remark) 
values(1,'t_enterprise','et_scale','0-19','1',"系统初始化");
insert into t_params(param_id,table_name,column_name,column_value,param_value,remark) 
values(2,'t_enterprise','et_scale','20-99','2',"系统初始化");
insert into t_params(param_id,table_name,column_name,column_value,param_value,remark) 
values(3,'t_enterprise','et_scale','100-499','3',"系统初始化");
insert into t_params(param_id,table_name,column_name,column_value,param_value,remark) 
values(4,'t_enterprise','et_scale','500-999','4',"系统初始化");
insert into t_params(param_id,table_name,column_name,column_value,param_value,remark)
values(5,'t_enterprise','et_scale','1000-9999','5',"系统初始化");
insert into t_params(param_id,table_name,column_name,column_value,param_value,remark)
values(6,'t_enterprise','et_scale','10000-','6',"系统初始化");

-- 学历参数

insert into t_params(param_id,table_name,column_name,column_value,param_value,remark)
values(7,'','high_edu','高中以下','1',"系统初始化");
insert into t_params(param_id,table_name,column_name,column_value,param_value,remark)
values(8,'','high_edu','高中','2',"系统初始化");
insert into t_params(param_id,table_name,column_name,column_value,param_value,remark)
values(9,'','high_edu','大专','3',"系统初始化");
insert into t_params(param_id,table_name,column_name,column_value,param_value,remark)
values(10,'','high_edu','大学本科','4',"系统初始化");
insert into t_params(param_id,table_name,column_name,column_value,param_value,remark)
values(11,'','high_edu','硕士','5',"系统初始化");
insert into t_params(param_id,table_name,column_name,column_value,param_value,remark)
values(12,'','high_edu','博士','6',"系统初始化");
insert into t_params(param_id,table_name,column_name,column_value,param_value,remark)
values(13,'','high_edu','博士后','7',"系统初始化");

-- status状态
insert into t_params(param_id,table_name,column_name,column_value,param_value,remark)
values(14,'t_job_deliver','status','0','未查看',"系统初始化");

insert into t_params(param_id,table_name,column_name,column_value,param_value,remark)
values(15,'t_job_deliver','status','1','已查看',"系统初始化");

insert into t_params(param_id,table_name,column_name,column_value,param_value,remark)
values(16,'t_job_info','status','0','在用',"系统初始化");

insert into t_params(param_id,table_name,column_name,column_value,param_value,remark)
values(17,'t_job_info','status','1','已下架',"系统初始化");

-- 工作年限
insert into t_params(param_id,table_name,column_name,column_value,param_value,remark)
values(18,'','work_age','应届毕业生','1',"系统初始化");
insert into t_params(param_id,table_name,column_name,column_value,param_value,remark)
values(19,'','work_age','3年及以下','2',"系统初始化");
insert into t_params(param_id,table_name,column_name,column_value,param_value,remark)
values(20,'','work_age','3-5年','3',"系统初始化");
insert into t_params(param_id,table_name,column_name,column_value,param_value,remark)
values(21,'','work_age','5-10年','4',"系统初始化");
insert into t_params(param_id,table_name,column_name,column_value,param_value,remark)
values(22,'','work_age','10年以上','5',"系统初始化");

-- job类型参数
insert into t_job_type(job_type_id,job_type_name,parent_jt_id) values(1,"区块链",0);

