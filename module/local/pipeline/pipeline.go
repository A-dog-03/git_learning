package pipeline

import (
	"fmt"

	"../../../module"

	"../../../module/stub"

	"../../../log"
)

// logger 代表日志记录器。
var logger = log.DLogger()

// myPipeline 代表条目处理管道的实现类型。
type myPipeline struct {
	// stub.ModuleInternal 代表组件基础实例。
	stub.ModuleInternal
	// itemProcessors 代表条目处理函数的列表。
	itemProcessors []module.ProcessItem
	// failFast 代表处理是否需要快速失败。
	failFast bool
}

// New 用于创建一个条目处理管道实例。
func New(
	mid module.MID,
	// 条目处理函数列表
	itemProcessors []module.ProcessItem,
	// 得分计算函数
	scoreCalculator module.CalculateScore) (module.Pipeline, error) {
	// 创建基础模型
	moduleBase, err := stub.NewModuleInternal(mid, scoreCalculator)
	if err != nil {
		return nil, err
	}
	if itemProcessors == nil {
		return nil, genParameterError("nil item processor list")
	}
	if len(itemProcessors) == 0 {
		return nil, genParameterError("empty item processor list")
	}
	var innerProcessors []module.ProcessItem
	for i, pipeline := range itemProcessors {
		if pipeline == nil {
			err := genParameterError(fmt.Sprintf("nil item processor[%d]", i))
			return nil, err
		}
		innerProcessors = append(innerProcessors, pipeline)
	}
	return &myPipeline{
		ModuleInternal: moduleBase,
		itemProcessors: innerProcessors,
	}, nil
}

// 返回条目条目处理函数列表
func (pipeline *myPipeline) ItemProcessors() []module.ProcessItem {
	processors := make([]module.ProcessItem, len(pipeline.itemProcessors))
	copy(processors, pipeline.itemProcessors)
	return processors
}

func (pipeline *myPipeline) Send(item module.Item) []error {
	// 实时处理计数加一
	pipeline.ModuleInternal.IncrHandlingNumber()
	// 实时处理计数减一
	defer pipeline.ModuleInternal.DecrHandlingNumber()
	// 调用计数加一
	pipeline.ModuleInternal.IncrCalledCount()
	// 条目不为空即接收
	var errs []error
	if item == nil {
		err := genParameterError("nil item")
		errs = append(errs, err)
		return errs
	}
	// 接收计数加一
	pipeline.ModuleInternal.IncrAcceptedCount()
	logger.Infof("Process item %+v... \n", item)
	// 处理条目
	var currentItem = item
	for _, processor := range pipeline.itemProcessors {
		processedItem, err := processor(currentItem)
		if err != nil {
			errs = append(errs, err)
			// 如果是快速失败就立即返回
			if pipeline.failFast {
				break
			}
		}
		if processedItem != nil {
			currentItem = processedItem
		}
	}
	// 没有错误完成计数加一
	if len(errs) == 0 {
		pipeline.ModuleInternal.IncrCompletedCount()
	}
	return errs
}

func (pipeline *myPipeline) FailFast() bool {
	return pipeline.failFast
}

func (pipeline *myPipeline) SetFailFast(failFast bool) {
	pipeline.failFast = failFast
}

// extraSummaryStruct 代表条目处理管道实例额外信息的摘要类型。
type extraSummaryStruct struct {
	FailFast        bool `json:"fail_fast"`
	ProcessorNumber int  `json:"processor_number"`
}

// 返回解析器摘要的详细信息
func (pipeline *myPipeline) Summary() module.SummaryStruct {
	summary := pipeline.ModuleInternal.Summary()
	summary.Extra = extraSummaryStruct{
		FailFast:        pipeline.failFast,
		ProcessorNumber: len(pipeline.itemProcessors),
	}
	return summary
}
