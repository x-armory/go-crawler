# go-crawler
get data easily to custom struct using annotations, 

```text
main（需要你自己完成）
    配置：Crawlers
    启动所有Crawler

Crawler
    配置：Name、Cron、RequestGenerator、DataUnmarshaler、DataProcessor、Notification
    方法：
        Start()
            RequestGenerator.Gen()
            DataUnmarshaler.unmarshal(request, dataType, this.DataProcessor.Validate())
            DataProcessor.Process(data)
            Notification.Send()

RequestGenerator
    配置：GenParameters方法（返回Method、URL、Values、Header）
    方法：
        Gen()，调用抽象GenParameters()，生成http.Request
    实现类：
        TimeRequestGenerator
            配置：上次同步时间，配置滚动周期，是否忽略周末（默认忽略），最早数据时间、Offset（重复同步N天内的数据，有更新）
            方法：
                GenParameters()，缓存上次同步时间，调用timeUtil.NextPeriod()，生成请求参数
        TreeRequestGenerator
            配置：根节点URL，递归层数，最大执行数（默认无限），根据字段注解，决定下一层URL字段和下一层数据结构类型
            方法：
                GenParameters()，缓存上次处理的数据，递归生成下次请求需要的Method、URL、Values、Header
        ListRequestGenerator
            配置：Method、URL、Header、Values列表
            方法：
                GenParameters()，根据配置循环返回Method、URL、Header、Values

DataUnmarshaler
    配置：文件类型（zip xls xlsx csv txt html）、unmarshal方法
    方法：
        unmarshal(*http.Request, *Type, func validate())
    实现类：
        XpathUnmarshaler
        XlsUnmarshaler
        XlsxUnmarshaler
        CsvUnmarshaler

 DataProcessor
    配置：Dao
    方法：
        Validate() 校验数据是否有效
        Process()
        GenReport()

 Notification
    方法：
        Send()
    实现类：
        DingTalk
        WeiChart
```