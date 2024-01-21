package worker

import (
	"context"
	"log"
	"os"

	"github.com/camunda/zeebe/clients/go/v8/pkg/entities"
	"github.com/camunda/zeebe/clients/go/v8/pkg/worker"
	"github.com/camunda/zeebe/clients/go/v8/pkg/zbc"
	ZBMessage "github.com/hkm15022001/Supply-Chain-Event-Management/internal/service/zeebe/message"
)

// RunLongShipFinish to start this worker
func RunLongShipFinish() {
	client, err := zbc.NewClient(&zbc.ClientConfig{
		GatewayAddress:         os.Getenv("BROKER_ADDRESS"),
		UsePlaintextConnection: true,
	})
	if err != nil {
		panic(err)
	}
	go client.NewJobWorker().JobType("long_ship_finish").Handler(handleJobLongShipFinish).Open()
}

func handleJobLongShipFinish(client worker.JobClient, job entities.Job) {
	log.Print("Start job LongshipFinish")
	jobKey := job.GetKey()

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		// failed to handle job as we require the variables
		failJob(client, job)
		return
	}
	var uintLongShipID uint
	longShipID, ok := variables["long_ship_id"].(float64)
	if ok == true {
		uintLongShipID = uint(longShipID)
	} else {
		failJob(client, job)
		return
	}
	if err = ZBMessage.LongShipFinished(uintLongShipID); err != nil {
		failJob(client, job)
		return
	}

	request, err := client.NewCompleteJobCommand().JobKey(jobKey).VariablesFromMap(variables)
	if err != nil {
		// failed to set the updated variables
		failJob(client, job)
		return
	}

	log.Println("Complete job", jobKey, "of type", job.Type)
	log.Println("Sending finish long ship id:", uintLongShipID)

	ctx := context.Background()
	_, err = request.Send(ctx)
	if err != nil {
		panic(err)
	}

	log.Println("Successfully completed job")
}

func failJob(client worker.JobClient, job entities.Job) {
	log.Println("Failed to complete job", job.GetKey())

	ctx := context.Background()
	_, err := client.NewFailJobCommand().JobKey(job.GetKey()).Retries(job.Retries - 1).Send(ctx)
	if err != nil {
		panic(err)
	}
}
