# NATS Golang example JetStream usage with Group Queue

$ natscli stream info
? Select a Stream MyJStream
Information for Stream MyJStream created 2022-05-27T10:11:57Z

Configuration:

             Subjects: MyJSTopic.*
     Acknowledgements: true
            Retention: File - WorkQueue
             Replicas: 1
       Discard Policy: Old
     Duplicate Window: 2m0s
    Allows Msg Delete: true
         Allows Purge: true
       Allows Rollups: false
     Maximum Messages: unlimited
        Maximum Bytes: unlimited
          Maximum Age: unlimited
 Maximum Message Size: unlimited
    Maximum Consumers: unlimited


State:

             Messages: 411
                Bytes: 46 KiB
             FirstSeq: 1 @ 2022-05-27T10:29:59 UTC
              LastSeq: 411 @ 2022-05-27T10:55:03 UTC
     Active Consumers: 1


? Select a Consumer MyQueueConsumer
Information for Consumer MyJStream > MyQueueConsumer created 2022-05-27T10:12:55Z

Configuration:

        Durable Name: MyQueueConsumer
           Pull Mode: true
      Filter Subject: MyJSTopic.*
      Deliver Policy: All
          Ack Policy: Explicit
            Ack Wait: 30s
       Replay Policy: Instant
     Max Ack Pending: 1,000
   Max Waiting Pulls: 512

State:

   Last Delivered Message: Consumer sequence: 4,568 Stream sequence: 411 Last delivery: 1.75s ago
     Acknowledgment floor: Consumer sequence: 0 Stream sequence: 0
         Outstanding Acks: 411 out of maximum 1,000
     Redelivered Messages: 411
     Unprocessed Messages: 0
            Waiting Pulls: 2 of maximum 512
? Select a Consumer MyQueueConsumer
Information for Consumer MyJStream > MyQueueConsumer created 2022-05-27T10:12:55Z

Configuration:

        Durable Name: MyQueueConsumer
           Pull Mode: true
      Filter Subject: MyJSTopic.*
      Deliver Policy: All
          Ack Policy: Explicit
            Ack Wait: 30s
       Replay Policy: Instant
     Max Ack Pending: 1,000
   Max Waiting Pulls: 512

State:

   Last Delivered Message: Consumer sequence: 4,568 Stream sequence: 411 Last delivery: 1.75s ago
     Acknowledgment floor: Consumer sequence: 0 Stream sequence: 0
         Outstanding Acks: 411 out of maximum 1,000
     Redelivered Messages: 411
     Unprocessed Messages: 0
            Waiting Pulls: 2 of maximum 512
