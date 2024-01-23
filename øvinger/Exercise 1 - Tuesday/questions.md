Exercise 1 - Theory questions
-----------------------------

### Concepts

What is the difference between *concurrency* and *parallelism*?
>  Parallelism is when you actually run two processes at the same time, while concurrency progresses multiple tasks concurrently, but not simultaneously by randomly selecting the next thread 

What is the difference between a *race condition* and a *data race*? 
>  Race condition is any process where the order of operations affects the outcome, data race is a type of race condition that applies to writing and reading from the same object
 
*Very* roughly - what does a *scheduler* do, and how does it do it?
> Scheduler is the system that decides which thread should be run next, and can be done in a variety of ways


### Engineering

Why would we use multiple threads? What kinds of problems do threads solve?
> When we want to to multiple things at the same time, but the processes are unrelated, it is useful to seperate the processes from eachother and we use multiple threads to accomplish this

Some languages support "fibers" (sometimes called "green threads") or "coroutines"? What are they, and why would we rather use them over threads?
> If we want to divide threads into multiple tasks, we can do this by introducing another custom scheduler for the tasks and implementing a form for "blocking".

Does creating concurrent programs make the programmer's life easier? Harder? Maybe both?
> Once the systems for doing so are implemented the code becomes easier to read and better organized, but the introduction of complicated systems to allow this can come with problems as well

What do you think is best - *shared variables* or *message passing*?
> message passing is the best when working with concurrency. 


