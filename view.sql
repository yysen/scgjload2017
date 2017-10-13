 create or replace view ent_info as 
 SELECT
    pripid,
    '01' as rtype,
    cast(null as varchar2(32))  as uuid,
    nvl(UNISCID,'--') as 社会信用代码,
    cast(null as varchar2(1)) as 工商业务类型,
    cast(null as varchar2(9)) as 组织机构代码,
    nvl(REGNO,'--') AS 注册号,
    ENTNAME AS 名称,
    cast(null as varchar2(4)) as 类型,
    cast(null as varchar2(1)) as 单位类型,
    cast(null as varchar2(3)) as 登记注册类型,
    cast(null as varchar2(1)) as 控股情况,
    cast(null as varchar2(2)) as 机构类型,
    cast(null as varchar2(2)) as 执行会计标准类别,
    regorg as 企业登记机关,
    cast(null as varchar2(15)) as 数据处理地代码,
    to_char(ESTDATE,'yyyyMMdd') as 开业日期,
    to_char(OPFROM,'yyyyMMdd') as 经营期限自,
    to_char(OPTO,'yyyyMMdd') as 经营期限止,
    cast(null as varchar2(6)) as 经营期限,
    DOM as 住所,
    POSTALCODE as 邮政编码,
    cast(null as varchar2(60)) as 联系电话,
    PROLOC as 生产经营地址,
    OPSCOPE as 经营范围,
    cast(REGCAP as varchar2(24)) as 注册资本,
    cast(decode(trim(REGCAPCUR),
        '澳大利亚元','36',
        '奥地利先令','40',
        '新加坡元','702',
        '人民币','156',
        '港币','344',
        '港元','344',
        '意大利里拉','380',
        '日元','392',
        '韩元','410',
        '荷兰盾','528',
        '新西兰元','554',
        '挪威克朗','578',
        '欧元','954',
        '瑞典克朗','752',
        '瑞士法郎','756',
        '英镑','826',
        '美元','840',
        '比利时法郎','56',
        '加拿大元','124',
        '丹麦克朗','208',
        '法国法郎','250',
        '德国马克','280',
        '   ') as varchar2(3)) as 货币种类,
    cast(null as varchar2(24)) as 货币金额,
    '0' as 信息操作类型,
    cast(null as varchar2(14)) as 数据修改时间,
    cast(null as varchar2(40)) AS 行政区划代码,
    cast(null as varchar2(20)) as 数据包编码,
    cast(decode(trim(REGSTATE),
          '已注销','1',
          '注销','1',
          '0' ) as varchar2(1) ) AS 是否注销,
    cast(null as varchar2(32)) as 批次号,
    cast(null as varchar2(100)) as 上级主管部门名称,
    NAME as 法定代表人,
    cast(null as varchar2(50)) AS 财务负责人,
    cast(null as varchar2(10)) as 投资人数量,
    cast(null as varchar2(10)) as 下级子公司数量,
    cast(null as varchar2(10)) as 变更时间,
    '999' as 状态,
    '0' as 是否重码,
    '0' as 人工处理结果类型,
    '0' as 是否已审核,
    '1' as 是否已推送到名录库,
    cast(null as varchar2(14)) as 数据上传时间,
    cast(null as varchar2(14)) as 审核时间
 FROM hz_e_baseinfo;
   / 
    create or replace view 法人信息表 as
      select
        pripid,
       '02' as rtype,
       cast(null as varchar2(32)) as zuuid,
        rawtohex(sys_guid())  AS uuid,
        NAME as 法人姓名,
        CERTYPE as 证件类型,
        CERNO as 证件编号,
        TELNUMBER as 固话,
        MOBTEL as 移动电话,
        EMAIL as 电子邮件
    from hz_e_pri_person;
/
    create or replace view 财务负责人 as
       select
          pripid,
          '03' as rtype,
          cast(null as varchar2(32)) as zuuid,
          rawtohex(sys_guid())  AS uuid,
          NAME as 负责人姓名,
          CERTYPE as 证件类型,
          CERNO as 证件编号,
          TEL as 固话,
          MOBTEL as 移动电话,
          EMAIL as 电子邮件
    from hz_dzhy_e_fin_leader ;  
/
    create or replace view 上级主管部门信息表 as
       select
          pripid,
         '04'  as rtype,
          cast(null as varchar2(32)) as zuuid,
          rawtohex(sys_guid())  AS uuid
    from hz_dzhy_e_fin_leader where pripid='无用';
  /  
    create or replace view  投资人信息 as
       select 
           pripid,
          '05' as rtype,
          cast(null as varchar2(32)) as zuuid,
          rawtohex(sys_guid())  AS uuid,
          INVTYPE as 投资类型,
          INV as 投资方名称,
          BLICTYPE as 证件类型,
          BLICNO as 证件编号,
          cast(null as varchar2(10)) as 证照类型,
          cast(null as varchar2(60)) as 证照号码,
          to_char(SUBCONAM) as 认缴出资额,
          cast(null as varchar2(20)) as 认缴出资时间,
          cast(null as varchar2(100)) as 认缴出资方式,
          to_char(SUBCONPROP) as 认缴出资比例,
          COUNTRY as 国籍,
          cast(null as varchar2(10)) as 币种
    from hz_e_inv;
/
    create or replace view 下级子公司信息表 as
      select
          pripid,
         '04'  as rtype,
          cast(null as varchar2(32)) as zuuid,
          rawtohex(sys_guid())  AS uuid
    from hz_dzhy_e_fin_leader where pripid='无用';
/
    create or replace view  变更信息表 as
       select 
           pripid,
          '07' as rtype,
          cast(null as varchar2(32)) as zuuid,
          rawtohex(sys_guid())  AS uuid,
          ALTITEM as 变更事项,
          ALTBE as 变更前内容,
          ALTAF as 变更后内容,
          to_char(ALTDATE,'yyyyMMdd') as 变更日期,
          cast(null as varchar2(20)) as 数据包编号
    from hz_dzhy_e_alter_recoder;
/
    create or replace view  注销信息 as
       select 
           pripid,
          '08' as rtype,
          cast(null as varchar2(32)) as zuuid,
          rawtohex(sys_guid())  AS uuid,
          to_char(CANDATE,'yyyyMMdd') as 注销日期,
          cast(null as varchar2(200)) as 注销原因
    from hz_e_cancel;
/
    create or replace view  吊销信息 as
       select 
           pripid,
          '08' as rtype,
          cast(null as varchar2(32)) as zuuid,
          rawtohex(sys_guid())  AS uuid,
          to_char(REVDATE,'yyyyMMdd') as 注销日期,
          REVBASIS as 注销原因
    from hz_e_revoke; 
/
create or replace view illegal_info as 
 SELECT
    illid,
    '01' as rtype,
    cast(null as varchar2(32))  as uuid,
    nvl(UNISCID,'--') AS 社会信用代码,
    ENTNAME AS 名称,
    nvl(REGNO,'--') AS 注册号,
    to_char(ABNTIME,'yyyyMMdd') as 列入日期,
    SERILLREA_CN as 列入原因,
    DECORG_CN as 决定机关,
    cast(null as varchar2(32)) as 批次号,
    cast(null as varchar2(6)) AS 行政区划代码,
    cast(null as varchar2(6)) AS 季度,
    cast(null as varchar2(14)) as 数据上传时间,
    cast(null as varchar2(14)) as 修改时间,
    '999' as 状态,
    '0' as 是否重码,
    '0' as 是否已审核,
    '1' as 是否已推送到名录库,
    cast(null as varchar2(14)) as 排序,
    cast(null as varchar2(14)) as 审核时间
 from hz_e_li_illdisdetail;
/
create or replace view abnormal_info as
  SELECT
    busexclist,
    '01' as rtype,
    cast(null as varchar2(32)) as uuid,
    nvl(UNISCID,'--') AS 社会信用代码,
    ENTNAME AS 名称,
    nvl(REGNO,'--') AS 注册号,
    to_char(ABNTIME,'yyyyMMdd') as 列入日期,
    SPECAUSE_CN as 列入原因,
    DECORG_CN as 决定机关,
    cast(null as varchar2(6)) AS 行政区划代码,
    cast(null as varchar2(32)) as 批次号,
    cast(null as varchar2(6)) AS 季度,
    cast(null as varchar2(14)) as 数据上传时间,
    cast(null as varchar2(14)) as 修改时间,
    '999' as 状态,
    '0' as 是否重码,
    '0' as 是否已审核,
    '1' as 是否已推送到名录库,
    cast(null as varchar2(14)) as 排序,
    cast(null as varchar2(14)) as 审核时间 
  from hz_ao_opa_detail;
/
create or replace view annualreport_info as
 SELECT
    ancheid,
    '01' as rtype,
    cast(null as varchar2(32)) as uuid,
    nvl(REGNO,'--') AS 注册号,
    nvl(UNISCID,'--') AS 社会信用代码,
    ENTNAME AS 名称,
    replace(substr(trim(ANCHEDATE),0,10),'-','') as 年报时间,
    ANCHEYEAR as 年报年度,
    cast(null as varchar2(30)) as 联系电话,
    cast(null as varchar2(100)) as 通信地址,
    cast(null as varchar2(6)) as 邮政编码,
    cast(null as varchar2(100)) as 电子邮箱,
    BUSST_CN as 经营状态,
    to_char(EMPNUM) as 从业人数,
    ENTTYPE as 企业类型,
    cast(null as varchar2(1)) as 单位类型,
    cast(null as varchar2(9)) as 组织机构,
    to_char(ASSGRO) as 资产总额,
    to_char(LIAGRO) as 负债总额,
    to_char(VENDINC) as 销售收入,
    to_char(MAIBUSINC) as 主营业务收入,
    to_char(PROGRO) as 利润总额,
    to_char(NETINC) as 净利润,
    to_char(RATGRO) as 纳税总额,
    to_char(TOTEQU) as 所有者权益合计,
    cast(null as varchar2(2000)) as 主营业务活动,
    cast(null as varchar2(10)) as 女性从业人员,
    cast(null as varchar2(1)) as 控股情况,
    cast(null as varchar2(18)) as 上级社会信用代码,
    cast(null as varchar2(9)) as 上级组织机构,
    cast(null as varchar2(30)) as 姓名,
    cast(null as varchar2(30)) as 移动电话,
    cast(null as varchar2(1)) as 是否有对外投资信息,
    cast(null as varchar2(1)) as 是否有网站网店信息,
    cast(null as varchar2(1)) as 是否有股权变更信息,
    cast(null as varchar2(32)) as 批次号,
    cast(null as varchar2(6)) as 行政区划,
    cast(null as varchar2(14)) as 数据上传时间,
    cast(null as varchar2(14)) as 修改时间,
    '999' as 状态,
    '0' as 是否已审核,
    '1' as 是否已推送到名录库,
    cast(null as varchar2(16)) as 排序,
    cast(null as varchar2(8)) as 认缴信息计数,
    cast(null as varchar2(8)) as 实缴信息计数,
    cast(null as varchar2(8)) as 投资信息计数,
    cast(null as varchar2(8)) as 网站信息计数,
    cast(null as varchar2(8)) as 股权变更计数,
    cast(null as varchar2(14)) as 审核时间,
     '0' as 社会信用代码是否重复
 from hz_an_baseinfo ;
 /
   create or replace view 企业认缴出资信息 as
      select
          ancheid,
          '03' as rtype,
          cast(null as varchar2(32)) as zuuid,
          rawtohex(sys_guid())  AS uuid
    from  hz_an_alterstockinfo where ancheid='无用';
    /
    create or replace view 企业实缴出资信息 as
      select
          ancheid,
          '03' as rtype,
          cast(null as varchar2(32)) as zuuid,
          rawtohex(sys_guid())  AS uuid
    from  hz_an_alterstockinfo where ancheid='无用';
/
    create or replace view 企业对外投资信息 as
       select
          ancheid,
          '03' as rtype,
          cast(null as varchar2(32)) as zuuid,
          rawtohex(sys_guid())  AS uuid
    from  hz_an_alterstockinfo where ancheid='无用';
/
    create or replace view 网站或网店信息 as
      select
          ancheid,
          '03' as rtype,
          cast(null as varchar2(32)) as zuuid,
          rawtohex(sys_guid())  AS uuid
    from  hz_an_alterstockinfo where ancheid='无用';
    /
    create or replace view 股权变更信息 as
      select
          ancheid,
          '06' as rtype,
          cast(null as varchar2(32)) as zuuid,
          rawtohex(sys_guid()) AS uuid,
          cast(null as varchar2(2000)) as 股权是否转让,
          INV as 股东名称,
          to_char(TRANSAMPR) as 转让前股权比例,
          to_char(TRANSAMAFT) as 转让后股权比例,
          to_char(ALTDATE,'yyyyMMdd') as 股权变更日期
    from hz_an_alterstockinfo;