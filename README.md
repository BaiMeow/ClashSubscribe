# 我就算只有一个节点我也要CLASH订阅!!!

## 使用

### 初始化

请务必配合<https://github.com/Jrohy/trojan>使用,只照顾trojan

且将所有的trojan指向同一个数据库

此外你还要自行创建一个数据库来储存所有服务器的地址

```sql
CREATE TABLE `servers`  (
  `name` varchar(255),
  `addr` varchar(255),
  `port` bigint(255),
  `area` varchar(255)
)
```

大概就这样，也请你手动录入一下数据

|字段|解释|
|--|--|
|name|名称，显示在clash中|
|addr|服务器ip地址|
|port|服务器端口|
|area|服务器地域，暂时没有用处|

### 订阅

默认25001端口，已经写死了，请使用nginx进行代理，并以此套上https

必须要开https，原因从下文显然可见

订阅地址格式为`https://yourhost.com/clash?passwd=yourpasswd`

yourhost是代理后的域名

推荐加一个clash目录

yourpasswd是trojan那边的passwd，数据库统一后各个节点的密码也是统一的
