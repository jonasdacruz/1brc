# 1brc

## Personal Message

This is my personal approach using Go for solving the [1 Billion Rows Challenge](https://www.morling.dev/blog/one-billion-row-challenge/) launched early 2024. I'm not an expert in Go, actually I'm not an expert in any language anymore :cry: and I'm not coding much last years :sob:.

Anyway, I think this is an excellent exercise for anyone to face hard programing problems that nowadays it is not a very common situation for most developers due to the crazy amount of abstractions between the computer and the programmer.

Ok, I'm just finished to bla bla bla bla... let's code.

## Apporoach

### Basics

The 1 Billion Rows file must be created with a python script using the following command

```
python ./scripts/create_measurements.py 1000000000
```

It could be a good idea to generate a smaller one for testing, I didn't do it and I regret it.

### The Problem

Due to the really simple nature of the problem functioanlly speaking, the only complexity here is the actual challenge about try to process as fast as possible.

So, the most complicated thing to me is to try to keep the whole computer actually computing effectively as much as possible, and at the same time being careful with overlapping memory access.

Last words, this is 100% concurrency problem, so after being able to use the whole computing power, the memory bottleneck will be the biggest problem.

### My first approach (2024-05-30)

I alway start [drawing main ideas](https://link.excalidraw.com/readonly/3Yp8PTdbYENO8o0pyTQQ?darkMode=true) to drive the coding work, I usualy don't code anything until I have a clear idea of what I want to do.

Then my second hobby is to break down the problem in smaller problems using as much sub-routines as possible in order to isolate responsibilities and problems.

So, my de docomposition was:

- Open the file and calculate byte positions to create several and consistent blocks of bytes to be processed in parallel.
- Process each block to find specific lines.
- Process each line to calculate values of each block.
- Merge values of each block to have the final result.

you can see this approach in the v2 folder

## Fork it!

Don't hesitate to fork this repo to improve or change totally the approach.
