# limit-rate

利用go-echo，简单测试并实现了常见的限流算法：
1.固定窗口
2.滑动窗口
3.漏桶算法
4.令牌桶算法
5.滑动日志算法

# 用法
如果想要测试不同的算法实现，只需要修改configs/config.yaml中的配置即可。
现在使用`make main`一键启动

# 限流算法
文件目录：internal/limiter
- fixed_window
- sliding_window
- leaky_bucket
- token_bucket
- sliding_log
