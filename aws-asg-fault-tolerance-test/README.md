## Fault Tolerance Test

### Removes random no of randomly selected EC2 instances from an ASG

- Run locally after logging into the desired profile through AWS CLI
### Inputs example: 
- ```Profile Name``` : sandbox
- ```Group Name``` : sandbox-zvault

### Variables 
- ```min_sleep_time``` : min time to sleep between instance termination
- ```max_sleep_time``` : max time to sleep between instance termination
- ```loop_count``` : no of iterations to test
