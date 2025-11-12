# Mandatory-Activity-4
## Implementation of Ricart-Argawala

### How to run the program 
```markdown
Installation 

1. Clone the repository
   git clone https://github.com/alexmaryna/Mandatory-Activity-4.git
   cd Mandatory-Activity-4

2. Install dependencies
   go mod download


Running the system

1. Run the program
   go run ./node/src/main.go

Expected output (with local date and timestamps):
  2025/11/11 21:35:59 Node 1: localhost:5001
  2025/11/11 21:35:59 Node 2: localhost:5002
  2025/11/11 21:35:59 Node 3: localhost:5003
  2025/11/11 21:35:59 Node 1 starting gRPC Server on localhost:5001
  2025/11/11 21:35:59 Node 2 starting gRPC Server on localhost:5002
  2025/11/11 21:35:59 Node 3 starting gRPC Server on localhost:5003

2. Automatic test sequence
   The system automatically runs three tests:
     1. One node requests access to the Critical Section (CS)
     2. Two nodes request access at the same time
     3. All three of the nodes request access simultaneously


For example:
  2025/11/11 21:36:01 Test 1: One node can request CS
  2025/11/11 21:36:01 Node 1 entering CS
  2025/11/11 21:36:04 Node 1 exit CS
  
  2025/11/11 21:36:06 Test 2: Two nodes can request CS at the same time
  2025/11/11 21:36:07 Node 2 entering CS
  2025/11/11 21:36:09 Node 2 exit CS
  
  2025/11/11 21:36:26 Test 3: All nodes can request CS
  2025/11/11 21:36:27 Node 1 entering CS
  2025/11/11 21:36:29 Node 1 exit CS
  2025/11/11 21:36:30 Node 2 entering CS


Notes:
- All communication between the nodes uses gRPC.
- The .proto file defines the RequestAccess and ReplyAccess RPCs.
- The program includes all of the generated gRPC files, so no regeneration of them is needed.
