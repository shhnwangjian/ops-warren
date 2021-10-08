package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

/*
prometheus.MustRegister() 函数来将指标注册到全局默认注册中，此外我们还可以使用 prometheus.NewRegistry() 函数来创建和使用自己的非全局的注册表。

为什么我们还需要自定义注册表呢？
1、全局变量通常不利于维护和测试，创建一个非全局的注册表，并明确地将其传递给程序中需要注册指标的地方，这也一种更加推荐的做法。
2、全局默认注册表包括一组默认的指标，我们有时候可能希望除了自定义的指标之外，不希望暴露其他的指标。

创建的一个自定义的注册表：
1、首先使用 prometheus.NewRegistry() 函数创建我们自己的注册表对象。
2、使用自定义注册表对象上面的 MustRegister() 来注册 guage 指标，而不是调用 prometheus.MustRegister() 函数（这会使用全局默认的注册表）。
3、如果我们希望在自定义注册表中也有进程和 Go 运行时相关的指标，我们可以通过实例化 Collector 收集器来添加他们。
4、 最后在暴露指标的时候必须通过调用 promhttp.HandleFor() 函数来创建一个专门针对我们自定义注册表的 HTTP 处理器，
为了同时暴露前面示例中的 promhttp_* 相关的指标，我们还需要在 promhttp.HandlerOpts 配置对象的 Registry 字段中传递我们的注册表对象。
 */
func main() {
	// 创建一个自定义的注册表
	registry := prometheus.NewRegistry()
	// 可选: 添加 process 和 Go 运行时指标到我们自定义的注册表中
	// prometheus.NewProcessCollector 和prometheus.NewGoCollector 已废弃
	//registry.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
	//registry.MustRegister(prometheus.NewGoCollector())
	registry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	registry.MustRegister(collectors.NewGoCollector())

	// 创建一个简单呃 gauge 指标。
	temp := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ops_custom_registry_gauge",
		Help: "custom registry(gauge)",
	})

	// 使用我们自定义的注册表注册 gauge
	registry.MustRegister(temp)

	// 设置 gague 的值为 911
	temp.Set(911)

	// 暴露自定义指标
	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{Registry: registry}))
	if err := http.ListenAndServe(":9099", nil); err !=nil{
		log.Fatal(err)
	}
}

