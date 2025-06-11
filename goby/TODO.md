记得把鉴权中间件加回去

记得把hash密码改回去

SELECT NOW();
#查看时区
show variables like '%zone%';
select @@time_zone;
#修改mysql全局时区为北京时间
set global time_zone = '+8:00';
#修改当前会话时区
set time_zone = '+8:00';
#立即生效
flush privileges;





deleted的





db.Debug()



添加address