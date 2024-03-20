package main

import (
	"fmt"
)

func main() {
	a := "goroutine 122 [running]:\nruntime/debug.Stack()\n\t/usr/local/go/src/runtime/debug/stack.go:24 +0x5e\nc360-matchGID/internal/handler.(*MessageHandler).ProcessingReceiveMessages.func1.1()\n\t/app/internal/handler/message.go:42 +0x4b\npanic({0x9533e0?, 0xe3d210?})\n\t/usr/local/go/src/runtime/panic.go:770 +0x132\nc360-matchGID/internal/models.ResultRecord.IsEmptyCustomer(...)\n\t/app/internal/models/resultRecord.go:24\nc360-matchGID/internal/services/gid_service.(*Service).addCustomerID(0x9eab21?, 0x9?)\n\t/app/internal/services/gid_service/private.go:41 +0x45\nc360-matchGID/internal/services/gid_service.(*Service).SmartParsingMessage(0xc0003f4b40, 0xc0005c6000)\n\t/app/internal/services/gid_service/public.go:39 +0x159\nc360-matchGID/internal/handler.runSmartParsingMessage(...)\n\t/app/internal/handler/message.go:77\nc360-matchGID/internal/handler.(*MessageHandler).ProcessingReceiveMessages.func1(0xc00039c780)\n\t/app/internal/handler/message.go:66 +0x1a6\nc360-matchGID/pkg/broker.Consumer.ConsumeClaim({0xc000282120?}, {0xb03cb0, 0xc000282120}, {0xb01120?, 0xc00042e0f0?})\n\t/app/pkg/broker/kafka.go:49 +0x74\ngithub.com/IBM/sarama.(*consumerGroupSession).consume(0xc000282120, {0xc0005008c0, 0x10}, 0x0)\n\t/go/pkg/mod/github.com/!i!b!m/sarama@v1.42.1/consumer_group.go:949 +0x21b\ngithub.com/IBM/sarama.newConsumerGroupSession.func2({0xc0005008c0?, 0x0?}, 0x0?)\n\t/go/pkg/mod/github.com/!i!b!m/sarama@v1.42.1/consumer_group.go:874 +0x70\ncreated by github.com/IBM/sarama.newConsumerGroupSession in goroutine 89\n\t/go/pkg/mod/github.com/!i!b!m/sarama@v1.42.1/consumer_group.go:866 +0x456\n"
	fmt.Println(a)
}
