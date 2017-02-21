package main

import "testing"

func packetIDGenerator(c chan packetID) {
	for i := 0; i < 1000; i++ {
		c <- newPacketID()
	}
}

func TestUniquePacketIDs(T *testing.T) {
	//This will simulate a bunch of id generations
	//to see if we get any duplicates within a short amount of time
	//which is deemed a reasonable expected upper limit.
	aggregationChan := make(chan packetID, 1000) //make buffer size 1000 just in case
	go packetIDGenerator(aggregationChan)
	go packetIDGenerator(aggregationChan)
	var idArray []packetID
	for id := range aggregationChan {
		for _, otherID := range idArray {
			if id == otherID {
				T.Error("Encountered a duplicate ID")
			}
		}
		idArray = append(idArray, id)
		if len(idArray) == 2000 {
			break
		}
	}
}
