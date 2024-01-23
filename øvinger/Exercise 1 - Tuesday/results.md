It worked in C but not Go. not sure why
The issue occurs since the threads run concurrently and read and write without taking into account the other thread.

We use Mutex since this is more secure, and we dont need the versatility of semaphore