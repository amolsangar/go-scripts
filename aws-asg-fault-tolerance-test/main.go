package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
)

const (
	minSleepTime = 300  // min time to sleep between instance termination
	maxSleepTime = 600 // max time to sleep between instance termination
	loopCount     = 10   // no of iterations to test
)

func main() {
	fmt.Println("***************************")
	fmt.Println("ZVAULT FAULT TOLERANCE TEST")
	fmt.Println("***************************")

	// Random seeding for sleep and instance selection
	rand.Seed(time.Now().UnixNano())

	// GET PROFILE NAME
	fmt.Print("ENTER PROFILE NAME: ")
	reader := bufio.NewReader(os.Stdin)
	profile, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("An error occured while reading profile. Please try again", err)
		return
	}
	profile = strings.TrimSuffix(profile, "\n")

	// GET AUTOSCALING GROUP NAME
	fmt.Print("ENTER GROUP NAME: ")
	reader2 := bufio.NewReader(os.Stdin)
	asgGroup, err := reader2.ReadString('\n')
	if err != nil {
		fmt.Println("An error occured while reading profile. Please try again", err)
		return
	}
	asgGroup = strings.TrimSuffix(asgGroup, "\n")
	fmt.Println()

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedConfigProfile(profile))
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	asgClient := autoscaling.NewFromConfig(cfg)

	// Infinite loop
	// Removes random no of randomly selected EC2 instances from ASG
	count := loopCount
	for {
		fmt.Println("**************************************")

		// Get instances by group name
		input := &autoscaling.DescribeAutoScalingGroupsInput{AutoScalingGroupNames: []string{asgGroup}}
		result, err := asgClient.DescribeAutoScalingGroups(context.TODO(), input)
		if err != nil {
			fmt.Println("Got an error retrieving information about your AutoScaling Group EC2 instances:")
			fmt.Println(err)
			return
		}
		asgInstanceIdMap := make(map[string][]string)
		for _, r := range result.AutoScalingGroups[0].Instances {
			asgInstanceIdMap[asgGroup] = append(asgInstanceIdMap[asgGroup], *r.InstanceId)
		}
		if len(asgInstanceIdMap) == 0 {
			fmt.Println("Autoscaling group not present. Please check the group name again!\n")
			return
		}

		// Print filtered group instances
		for asgName, instances := range asgInstanceIdMap {
			fmt.Println("AutoScaling Group Name:", asgName)
			for idx, id := range instances {
				fmt.Printf("Instance %d: %v\n", idx, id)
			}
		}

		noOfInstances := len(asgInstanceIdMap[asgGroup])
		randomNoOfIterations := rand.Intn(noOfInstances-1) + 1 // [1,noOfInstances)
		fmt.Println("\nNo of instances scheduled to be deleted:", randomNoOfIterations)
		fmt.Println()

		for i := 0; i < randomNoOfIterations; i++ {
			// Randomly selecting an instance for termination
			randomIndex := rand.Intn(noOfInstances)
			randomInstanceId := asgInstanceIdMap[asgGroup][randomIndex]

			fmt.Println("Terminating Instance Id:", randomInstanceId)
			inputAsgTerminateInstance := &autoscaling.TerminateInstanceInAutoScalingGroupInput{
				InstanceId:                     aws.String(randomInstanceId),
				ShouldDecrementDesiredCapacity: aws.Bool(false),
			}
			response, err := asgClient.TerminateInstanceInAutoScalingGroup(context.TODO(), inputAsgTerminateInstance)
			if err != nil {
				fmt.Println("Got an error while terminating Amazon EC2 instance:")
				fmt.Println(err)
				return
			}

			// Printing response
			fmt.Println("Response:")
			s, _ := json.MarshalIndent(*response, "", "    ")
			fmt.Printf("%s\n\n", s)

			// Remove the selected instance from map
			slice1 := []string{}
			slice2 := []string{}
			slice1 = asgInstanceIdMap[asgGroup][:randomIndex]
			slice2 = asgInstanceIdMap[asgGroup][randomIndex+1:]
			newSlice := append(slice1, slice2...)
			asgInstanceIdMap[asgGroup] = newSlice
			noOfInstances -= 1
		}

		count--
		if count == 0 {
			fmt.Println("Exiting")
			break
		}

		// Sleep for random time
		r := rand.Intn(maxSleepTime-minSleepTime) + minSleepTime
		sleepSecs := time.Duration(r) * time.Second
		fmt.Println("Sleeping for", sleepSecs)
		for i := time.Duration(0); i < sleepSecs; i += time.Second {
			fmt.Print(".")
			time.Sleep(time.Second * 2)
		}
		fmt.Println("\n")
	}
}
