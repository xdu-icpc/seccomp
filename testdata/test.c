#include <pthread.h>
#include <string.h>
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>

void *thread_routine(void *useless)
{
	fprintf(stderr, "this is a thread\n");
	return useless;
}

int main()
{
	pthread_t thread;
	pid_t pid;

	int ret = pthread_create(&thread, NULL, &thread_routine, NULL);
	if (ret) {
		fprintf(stderr, "pthread_create: %s", strerror(ret));
		exit(127);
	}

	ret = pthread_join(thread, NULL);
	if (ret) {
		fprintf(stderr, "pthread_join: %s", strerror(ret));
		exit(127);
	}

	pid = fork();
	if (pid == -1) {
		perror("fork");
		exit(2);
	}

	return 0;
}
